package workers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"lastsaas/core/billing"
)

// Task names
const (
	TypeAIAnalysis = "ai:analyze"
)

// AIAnalysisPayload structs what is enqueued
type AIAnalysisPayload struct {
	TenantID       string `json:"tenant_id"`
	EditalID       string `json:"edital_id"`
	DocumentText   string `json:"document_text"`
	CompanyProfile string `json:"company_profile"`
	PastHistory    string `json:"past_history"`
}

// FastAPIResponse structs the exact Pydantic schema
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

// NewAIAnalysisTask creates the async job wrapper
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
	// Configure tight boundaries preventing zombie LLM requests
	return asynq.NewTask(TypeAIAnalysis, payload, asynq.MaxRetry(3), asynq.Timeout(2*time.Minute)), nil
}

// HandleAIAnalysisTask processes the queue
func HandleAIAnalysisTask(db *sql.DB) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p AIAnalysisPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}

		tId, _ := uuid.Parse(p.TenantID)
		eId, _ := uuid.Parse(p.EditalID)

		// 1. Mark as Processing
		_, _ = db.ExecContext(ctx, `
			UPDATE ai_analysis_results 
			SET status = 'processing', updated_at = NOW() 
			WHERE tenant_id = $1 AND edital_id = $2
		`, tId, eId)

		// 2. HTTP Call to FastAPI
		client := &http.Client{Timeout: 90 * time.Second}
		reqBody, _ := json.Marshal(p)
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

		// 3. Extrapolate Python Response
		var aiResp FastAPIResponse
		if err := json.Unmarshal(respBytes, &aiResp); err != nil {
			return markFailed(ctx, db, tId, eId, "Failed unmarshaling python AI JSON output")
		}

		structuredJsonBytes, _ := json.Marshal(map[string]interface{}{
			"strengths":    aiResp.Strengths,
			"risks":        aiResp.Risks,
			"requirements": aiResp.Requirements,
		})

		// 4. Upsert success metrics and audit log the raw Prompt Output!
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
			slog.Error("Failed saving AI results to database", "err", err)
			return err
		}

		// 5. Safely Increment SaaS Billing metrics exclusively ONLY on complete exact success
		err = billing.IncrementUsage(ctx, db, tId, "ai_analysis")
		if err != nil {
			slog.Error("Failed updating SaaS AI boundaries", "err", err, "tenant", p.TenantID)
		}

		slog.Info("Successfully Processed Asynq AI Licitacao", "edital_id", p.EditalID)
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
