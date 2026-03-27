package ai

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"lastsaas/workers"
)

type Service struct {
	db     *sql.DB
	client *asynq.Client
	redis  *redis.Client
}

func NewService(db *sql.DB, redisAddr string) *Service {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	return &Service{db: db, client: client, redis: rdb}
}

func (s *Service) AnalyzeEdital(ctx context.Context, tenantID, editalID uuid.UUID, text string) error {
	mockProfile := "Generic Construction Firm specializing in civil B2G infrastructure."
	mockHistory := "Won 3 previous bridges, defaulted on 0."

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO ai_analysis_results (tenant_id, edital_id, status, created_at, updated_at)
		VALUES ($1, $2, 'pending', NOW(), NOW())
		ON CONFLICT (tenant_id, edital_id) DO UPDATE SET status = 'pending', error = NULL, updated_at = NOW()
	`, tenantID, editalID)

	if err != nil {
		return errors.New("failed injecting native ai pipeline tracking: " + err.Error())
	}

	task, err := workers.NewAIAnalysisTask(tenantID.String(), editalID.String(), text, mockProfile, mockHistory)
	if err != nil {
		return err
	}

	_, err = s.client.EnqueueContext(ctx, task)
	return err
}

type AIResult struct {
	Status            string              `json:"status"`
	Score             *int                `json:"score"`
	SmartScore        *int                `json:"smart_score,omitempty"`
	WinPrediction     *WinProbabilityData `json:"win_prediction,omitempty"`
	Decision          *string             `json:"decision"`
	RiskLevel         *string             `json:"risk_level"`
	Summary           *string             `json:"summary"`
	DetailedReasoning *string             `json:"detailed_reasoning"`
	ProposalDraft     *string             `json:"proposal_draft"`
	Error             *string             `json:"error"`
	Category          *string             `json:"category,omitempty"`
}

type WinProbabilityData struct {
	WinProbability float64 `json:"win_probability"`
	Confidence     float64 `json:"confidence"`
	BasedOnCases   int     `json:"based_on_cases"`
}

type AIHistoryItem struct {
	EditalID      uuid.UUID           `json:"edital_id"`
	Status        string              `json:"status"`
	Score         int                 `json:"score"`
	SmartScore    int                 `json:"smart_score"`
	WinPrediction *WinProbabilityData `json:"win_prediction"`
	Confidence    float64             `json:"confidence"`
	Decision      string              `json:"decision"`
	RiskLevel     string              `json:"risk_level"`
	RealResult    *string             `json:"real_result"`
	CreatedAt     time.Time           `json:"created_at"`
}

type AIInsights struct {
	TotalAnalyses   int     `json:"total_analyses"`
	SuccessRate     float64 `json:"success_rate"`
	AverageScore    float64 `json:"average_score"`
	EnterRatio      float64 `json:"enter_ratio"`
	TotalFeedbacks  int     `json:"total_feedbacks"`
}

func (s *Service) GetAnalysisStatus(ctx context.Context, tenantID, editalID uuid.UUID) (*AIResult, error) {
	var res AIResult
	var structuredJSON sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT status, score, decision, risk_level, summary, detailed_reasoning, proposal_draft, error, structured_data::text
		FROM ai_analysis_results
		WHERE tenant_id = $1 AND edital_id = $2
	`, tenantID, editalID).Scan(
		&res.Status, &res.Score, &res.Decision, &res.RiskLevel,
		&res.Summary, &res.DetailedReasoning, &res.ProposalDraft, &res.Error, &structuredJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("analysis not found")
		}
		return nil, err
	}

	if res.Score != nil && res.Decision != nil {
		winPred, histScore, conf := s.calculateHistoricalMetrics(ctx, tenantID, editalID, *res.Score)
		res.WinPrediction = winPred
		smartScore := calculateSmartScore(*res.Score, histScore, conf)
		res.SmartScore = &smartScore
	}

	return &res, nil
}

func (s *Service) SubmitFeedback(ctx context.Context, tenantID, editalID uuid.UUID, aiDecision, realResult string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO ai_feedback (tenant_id, edital_id, ai_decision, real_result, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (tenant_id, edital_id) DO UPDATE SET ai_decision = EXCLUDED.ai_decision, real_result = EXCLUDED.real_result
	`, tenantID, editalID, aiDecision, realResult)
	return err
}

func (s *Service) GetInsights(ctx context.Context, tenantID uuid.UUID, timeframeDays int) (*AIInsights, error) {
	cacheKey := fmt.Sprintf("tenant:%s:ai_insights:%d", tenantID.String(), timeframeDays)
	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var insights AIInsights
		if err := json.Unmarshal([]byte(val), &insights); err == nil {
			return &insights, nil
		}
	}

	timeFilter := ""
	if timeframeDays > 0 {
		timeFilter = fmt.Sprintf("AND created_at >= NOW() - INTERVAL '%d days'", timeframeDays)
	}

	var insights AIInsights
	
	err = s.db.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT COUNT(*), COALESCE(AVG(score), 0)
		FROM ai_analysis_results
		WHERE tenant_id = $1 %s
	`, timeFilter), tenantID).Scan(&insights.TotalAnalyses, &insights.AverageScore)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	err = s.db.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT COUNT(*)
		FROM ai_analysis_results
		WHERE tenant_id = $1 AND decision = 'ENTER' %s
	`, timeFilter), tenantID).Scan(&insights.EnterRatio)
	if err == nil && insights.TotalAnalyses > 0 {
		insights.EnterRatio = insights.EnterRatio / float64(insights.TotalAnalyses)
	}

	var won, totalFeedback int
	err = s.db.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT COUNT(*), SUM(CASE WHEN real_result = 'WON' THEN 1 ELSE 0 END)
		FROM ai_feedback
		WHERE tenant_id = $1 %s
	`, timeFilter), tenantID).Scan(&totalFeedback, &won)
	if err == nil && totalFeedback > 0 {
		insights.SuccessRate = float64(won) / float64(totalFeedback)
		insights.TotalFeedbacks = totalFeedback
	}

	data, _ := json.Marshal(insights)
	s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return &insights, nil
}

func (s *Service) GetHistory(ctx context.Context, tenantID uuid.UUID) ([]AIHistoryItem, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT a.edital_id, a.status, COALESCE(a.score, 0), COALESCE(a.decision, ''), COALESCE(a.risk_level, ''),
		       f.real_result, a.created_at, a.structured_data::text
		FROM ai_analysis_results a
		LEFT JOIN ai_feedback f ON a.tenant_id = f.tenant_id AND a.edital_id = f.edital_id
		WHERE a.tenant_id = $1
		ORDER BY a.created_at DESC
		LIMIT 100
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []AIHistoryItem
	for rows.Next() {
		var item AIHistoryItem
		var realRes sql.NullString
		var structJSON sql.NullString
		if err := rows.Scan(&item.EditalID, &item.Status, &item.Score, &item.Decision, &item.RiskLevel, &realRes, &item.CreatedAt, &structJSON); err != nil {
			continue
		}
		if realRes.Valid {
			item.RealResult = &realRes.String
		}
		winPred, histScore, conf := s.calculateHistoricalMetrics(ctx, tenantID, item.EditalID, item.Score)
		item.WinPrediction = winPred
		item.Confidence = conf
		item.SmartScore = calculateSmartScore(item.Score, histScore, conf)

		history = append(history, item)
	}

	return history, nil
}

func (s *Service) calculateHistoricalMetrics(ctx context.Context, tenantID uuid.UUID, editalID uuid.UUID, aiScore int) (*WinProbabilityData, int, float64) {
	aiScoreNorm := float64(aiScore) / 100.0

	var targetEmbedding string
	err := s.db.QueryRowContext(ctx, `SELECT embedding::text FROM ai_embeddings WHERE tenant_id = $1 AND edital_id = $2`, tenantID, editalID).Scan(&targetEmbedding)
	if err != nil {
		w := &WinProbabilityData{WinProbability: aiScoreNorm, Confidence: 0.0, BasedOnCases: 0}
		return w, 50, 0.0
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT 
		    (1 - (e.embedding <=> $1::vector)) as raw_similarity,
		    f.real_result,
		    EXTRACT(DAY FROM (NOW() - f.created_at)) as days_ago
		FROM ai_embeddings e
		JOIN ai_feedback f ON e.tenant_id = f.tenant_id AND e.edital_id = f.edital_id
		WHERE e.tenant_id = $2 AND e.edital_id != $3 AND (e.embedding <=> $1::vector) < 0.25
	`, targetEmbedding, tenantID, editalID)
	
	if err != nil {
		w := &WinProbabilityData{WinProbability: aiScoreNorm, Confidence: 0.0, BasedOnCases: 0}
		return w, 50, 0.0
	}
	defer rows.Close()

	var totalSimilarity float64
	var weightedScore float64
	matchCount := 0

	for rows.Next() {
		var rawSim float64
		var res string
		var daysAgo float64
		if err := rows.Scan(&rawSim, &res, &daysAgo); err == nil {
			recencyFactor := 1.0
			if daysAgo > 365 {
				recencyFactor = 0.6
			} else if daysAgo > 90 {
				recencyFactor = 0.8
			}

			finalSim := rawSim * recencyFactor

			var outcomeScore float64
			if res == "WON" {
				outcomeScore = 1.0
			} else if res == "LOST" {
				outcomeScore = 0.0
			} else {
				continue
			}

			weightedScore += (finalSim * outcomeScore)
			totalSimilarity += finalSim
			matchCount++
		}
	}

	confidence := float64(matchCount) / 20.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	var winProb float64
	var histScore int

	if totalSimilarity == 0 {
		winProb = aiScoreNorm
		histScore = 50
	} else {
		historicalWinProb := weightedScore / totalSimilarity
		histScore = int(historicalWinProb * 100.0)
		winProb = (historicalWinProb * confidence) + (aiScoreNorm * (1.0 - confidence))
	}

	w := &WinProbabilityData{
		WinProbability: winProb,
		Confidence:     confidence,
		BasedOnCases:   matchCount,
	}

	return w, histScore, confidence
}

func (s *Service) GetSimilarEditais(ctx context.Context, tenantID, editalID uuid.UUID) ([]map[string]interface{}, error) {
	var targetEmbedding string
	err := s.db.QueryRowContext(ctx, `SELECT embedding::text FROM ai_embeddings WHERE tenant_id = $1 AND edital_id = $2`, tenantID, editalID).Scan(&targetEmbedding)
	if err != nil {
		return nil, errors.New("embedding_not_found")
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT 
		    e.edital_id,
		    f.real_result, 
		    COALESCE(a.score, 0), 
		    COALESCE(a.decision, ''), 
		    COALESCE(a.summary, ''), 
		    (1 - (e.embedding <=> $1::vector)) as similarity
		FROM ai_embeddings e
		JOIN ai_analysis_results a ON e.tenant_id = a.tenant_id AND e.edital_id = a.edital_id
		LEFT JOIN ai_feedback f ON e.tenant_id = f.tenant_id AND e.edital_id = f.edital_id
		WHERE e.tenant_id = $2 AND e.edital_id != $3 AND (e.embedding <=> $1::vector) < 0.25
		ORDER BY e.embedding <=> $1::vector ASC
		LIMIT 5
	`, targetEmbedding, tenantID, editalID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var similarCases []map[string]interface{}
	for rows.Next() {
		var sim float64
		var outcome sql.NullString
		var decision, summary string
		var score int
		var id uuid.UUID
		
		if err := rows.Scan(&id, &outcome, &score, &decision, &summary, &sim); err == nil {
			c := make(map[string]interface{})
			c["edital_id"] = id
			if outcome.Valid {
				c["outcome"] = outcome.String
			} else {
				c["outcome"] = "NOT_PARTICIPATED"
			}
			c["score"] = score
			c["decision"] = decision
			c["summary"] = summary
			c["similarity"] = sim
			similarCases = append(similarCases, c)
		}
	}

	return similarCases, nil
}

func calculateSmartScore(aiScore int, historicalScore int, confidence float64) int {
	var final float64
	if confidence < 0.3 {
		final = float64(aiScore)*0.9 + float64(historicalScore)*0.1
	} else if confidence > 0.7 {
		final = float64(aiScore)*0.6 + float64(historicalScore)*0.4
	} else {
		final = float64(aiScore)*0.7 + float64(historicalScore)*0.3
	}
	return int(final)
}
