package workers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"

	"lastsaas/core/billing"
)

// Task names
const (
	TypeAIAnalysis = "ai:analyze"
)

type AIAnalysisPayload struct {
	TenantID       string `json:"tenant_id"`
	EditalID       string `json:"edital_id"`
	DocumentText   string `json:"document_text"`
	CompanyProfile string `json:"company_profile"`
	PastHistory    string `json:"past_history"`
}

type FastAPIResponse struct {
	Score             int      `json:"score"`
	Decision          string   `json:"decision"`
	RiskLevel         string   `json:"risk_level"`
	Summary           string   `json:"summary"`
	DetailedReasoning string   `json:"detailed_reasoning"`
	ProposalDraft     string   `json:"proposal_draft"`
	Strengths         []string `json:"strengths"`
	Risks             []string `json:"risks"`
	Requirements      []string `json:"requirements"`
	EstimatedValue    float64  `json:"estimated_value"`
}

type EmbedRequest struct {
	DocumentText string `json:"document_text"`
}

type EmbedResponse struct {
	Embedding []float64 `json:"embedding"`
}

func NewAIAnalysisTask(tenantID, editalID, text, profile, history string) (*asynq.Task, error) {
	payload, err := json.Marshal(AIAnalysisPayload{
		TenantID:       tenantID,
		EditalID:       editalID,
		DocumentText:   text,
		CompanyProfile: profile,
		PastHistory:    history,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAIAnalysis, payload, asynq.MaxRetry(3), asynq.Timeout(5*time.Minute)), nil
}

func HandleAIAnalysisTask(db *sql.DB) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p AIAnalysisPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}

		tId, _ := uuid.Parse(p.TenantID)
		eId, _ := uuid.Parse(p.EditalID)

		_, _ = db.ExecContext(ctx, `
			UPDATE ai_analysis_results 
			SET status = 'processing', updated_at = NOW() 
			WHERE tenant_id = $1 AND edital_id = $2
		`, tId, eId)

		client := &http.Client{Timeout: 90 * time.Second}

		// 1. Generate text hash for the raw normalized document text string (first 10,000 chars roughly simulating object/requirements extracted)
		hashLength := len(p.DocumentText)
		if hashLength > 10000 {
			hashLength = 10000
		}
		hashSum := sha256.Sum256([]byte(p.DocumentText[:hashLength]))
		textHash := hex.EncodeToString(hashSum[:])

		var similarCases []map[string]interface{}
		var newEmbedding []float64

		// 2. Obtain Embedding securely falling back gracefully without halting base logic!
		embedReqBody, _ := json.Marshal(EmbedRequest{DocumentText: p.DocumentText[:hashLength]})
		embedReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8000/embed", bytes.NewBuffer(embedReqBody))
		embedReq.Header.Set("Content-Type", "application/json")
		
		embedResp, err := client.Do(embedReq)
		if err == nil && embedResp.StatusCode == 200 {
			defer embedResp.Body.Close()
			bodyBytes, _ := io.ReadAll(embedResp.Body)
			var er EmbedResponse
			if json.Unmarshal(bodyBytes, &er) == nil {
				newEmbedding = er.Embedding
			}
		} else {
			slog.Warn("Embedding extraction safely failed, continuing without similar context inject mapping", "error", err)
		}

		// 3. High-performance Similar Extraction (PgVector < 0.25 distance limits > 0.75 strict cosine match limit)
		if len(newEmbedding) > 0 {
			rows, searchErr := db.QueryContext(ctx, `
				SELECT f.real_result, a.score, a.decision, a.summary, (1 - (e.embedding <=> $1)) as similarity
				FROM ai_embeddings e
				JOIN ai_analysis_results a ON e.tenant_id = a.tenant_id AND e.edital_id = a.edital_id
				JOIN ai_feedback f ON e.tenant_id = f.tenant_id AND e.edital_id = f.edital_id
				WHERE e.tenant_id = $2 AND e.edital_id != $3 AND (e.embedding <=> $1) < 0.25
				ORDER BY e.embedding <=> $1 ASC
				LIMIT 3
			`, pq.Array(newEmbedding), tId, eId)

			if searchErr == nil {
				defer rows.Close()
				for rows.Next() {
					c := make(map[string]interface{})
					var sim float64
					var outcome, decision, summary sql.NullString
					var s int
					if err := rows.Scan(&outcome, &s, &decision, &summary, &sim); err == nil {
						c["outcome"] = outcome.String
						c["score"] = s
						c["decision"] = decision.String
						c["summary"] = summary.String
						c["similarity"] = sim
						similarCases = append(similarCases, c)
					}
				}
			}
		}

		// 4. Final FastAPI Analyze Execution Pipeline
		fullPayload := map[string]interface{}{
			"tenant_id": p.TenantID,
			"edital_id": p.EditalID,
			"document_text": p.DocumentText,
			"company_profile": p.CompanyProfile,
			"past_history": p.PastHistory,
			"similar_cases": similarCases,
		}

		reqBody, _ := json.Marshal(fullPayload)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8000/analyze", bytes.NewBuffer(reqBody))
		if err != nil {
			return markFailed(ctx, db, tId, eId, "Failed formatting FastAPI request")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return markFailed(ctx, db, tId, eId, err.Error())
		}
		defer resp.Body.Close()

		respBytes, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return markFailed(ctx, db, tId, eId, string(respBytes))
		}

		var aiResp FastAPIResponse
		if err := json.Unmarshal(respBytes, &aiResp); err != nil {
			return markFailed(ctx, db, tId, eId, "Failed unmarshaling python output")
		}

		structuredJsonBytes, _ := json.Marshal(map[string]interface{}{
			"strengths":    aiResp.Strengths,
			"risks":        aiResp.Risks,
			"requirements": aiResp.Requirements,
		})

		_, err = db.ExecContext(ctx, `
			UPDATE ai_analysis_results
			SET status = 'completed',
			    score = $1, decision = $2, risk_level = $3, summary = $4,
			    detailed_reasoning = $5, proposal_draft = $6,
			    structured_data = $7, raw_response = $8,
			    updated_at = NOW()
			WHERE tenant_id = $9 AND edital_id = $10
		`, aiResp.Score, aiResp.Decision, aiResp.RiskLevel, aiResp.Summary,
			aiResp.DetailedReasoning, aiResp.ProposalDraft,
			string(structuredJsonBytes), string(respBytes), tId, eId)

		if err != nil {
			return err
		}

		// Save Semantic Embedding if acquired
		if len(newEmbedding) > 0 {
			_, _ = db.ExecContext(ctx, `
				INSERT INTO ai_embeddings (tenant_id, edital_id, embedding, text_hash, created_at)
				VALUES ($1, $2, $3, $4, NOW())
				ON CONFLICT (tenant_id, edital_id) DO UPDATE SET embedding = EXCLUDED.embedding, text_hash = EXCLUDED.text_hash
			`, tId, eId, pq.Array(newEmbedding), textHash)
		}

		_ = billing.IncrementUsage(ctx, db, tId, "ai_analysis")
		return nil
	}
}

func markFailed(ctx context.Context, db *sql.DB, tenantID, editalID uuid.UUID, errorMsg string) error {
	_, _ = db.ExecContext(ctx, `
		UPDATE ai_analysis_results 
		SET status = 'failed', error = $1, updated_at = NOW() 
		WHERE tenant_id = $2 AND edital_id = $3
	`, errorMsg, tenantID, editalID)
	return errors.New(errorMsg)
}
