package ai

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"lastsaas/workers"
)

type Service struct {
	db     *sql.DB
	client *asynq.Client
}

func NewService(db *sql.DB, redisAddr string) *Service {
	// Normally injected structurally, initializing transient client for pipeline bridging
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &Service{db: db, client: client}
}

// AnalyzeEdital generates the initial tracking entity and fires the background event
func (s *Service) AnalyzeEdital(ctx context.Context, tenantID, editalID uuid.UUID, text string) error {
	// 1. Pre-fetch Contextually Aware Tenant Profiles!
	// Here we normally query `SELECT company_profile, past_history FROM tenant_settings WHERE id = $1`
	mockProfile := "Generic Construction Firm specializing in civil B2G infrastructure."
	mockHistory := "Won 3 previous bridges, defaulted on 0."

	// 2. Create foundational pending tracking log mapped explicitly to SaaS UI limits
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO ai_analysis_results (tenant_id, edital_id, status, created_at, updated_at)
		VALUES ($1, $2, 'pending', NOW(), NOW())
		ON CONFLICT (tenant_id, edital_id) DO UPDATE SET status = 'pending', error = NULL, updated_at = NOW()
	`, tenantID, editalID)

	if err != nil {
		return errors.New("failed injecting native ai pipeline tracking: " + err.Error())
	}

	// 3. Dispatch safe background async job ensuring frontend speed isn't blocked by generic Python HTTP sockets
	task, err := workers.NewAIAnalysisTask(tenantID.String(), editalID.String(), text, mockProfile, mockHistory)
	if err != nil {
		return err
	}

	_, err = s.client.EnqueueContext(ctx, task)
	return err
}

type AIResult struct {
	Status            string  `json:"status"`
	Score             *int    `json:"score"`
	Decision          *string `json:"decision"`
	RiskLevel         *string `json:"risk_level"`
	Summary           *string `json:"summary"`
	DetailedReasoning *string `json:"detailed_reasoning"`
	ProposalDraft     *string `json:"proposal_draft"`
	Error             *string `json:"error"`
}

func (s *Service) GetAnalysisStatus(ctx context.Context, tenantID, editalID uuid.UUID) (*AIResult, error) {
	var res AIResult
	err := s.db.QueryRowContext(ctx, `
		SELECT status, score, decision, risk_level, summary, detailed_reasoning, proposal_draft, error
		FROM ai_analysis_results
		WHERE tenant_id = $1 AND edital_id = $2
	`, tenantID, editalID).Scan(
		&res.Status, &res.Score, &res.Decision, &res.RiskLevel,
		&res.Summary, &res.DetailedReasoning, &res.ProposalDraft, &res.Error,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("analysis not found")
		}
		return nil, err
	}
	return &res, nil
}
