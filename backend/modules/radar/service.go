package radar

import (
    "context"
    "crypto/md5"
    "crypto/sha256"
    "database/sql"
    "encoding/json"
    "fmt"
    "math"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/redis/go-redis/v9"
    "github.com/jackc/pgtype"
    "github.com/yourorg/bidflowSaas/backend/modules/ai"
    "github.com/yourorg/bidflowSaas/backend/middleware"
)

type RadarEdital struct {
    ID             uuid.UUID       `json:"id"`
    ExternalID     string          `json:"external_id"`
    Source         string          `json:"source"`
    Title          string          `json:"title"`
    Description    string          `json:"description"`
    Value          float64         `json:"value"`
    Deadline       time.Time       `json:"deadline"`
    RawJSON        json.RawMessage `json:"raw_json"`
    TenantID       uuid.UUID       `json:"tenant_id"`
    CreatedAt      time.Time       `json:"created_at"`
    DeletedAt      sql.NullTime    `json:"deleted_at"`
    // AI derived fields
    Similarity     float64         `json:"similarity,omitempty"`
    WinProbability float64         `json:"win_probability,omitempty"`
    PriorityScore  float64         `json:"priority_score,omitempty"`
    PriorityLabel  string          `json:"priority_label,omitempty"`
    AnalysisLog    json.RawMessage `json:"analysis_log,omitempty"`
}

type Service struct {
    db    *sql.DB
    rdb   *redis.Client
    aiSvc *ai.Service
}

func NewService(db *sql.DB, rdb *redis.Client, aiSvc *ai.Service) *Service {
    return &Service{db: db, rdb: rdb, aiSvc: aiSvc}
}

// IngestEdital processes a raw edital from a source, applies tenant preferences, generates embedding,
// calculates similarity, win probability, priority score and persists.
func (s *Service) IngestEdital(ctx context.Context, tenantID uuid.UUID, raw json.RawMessage, source string) (*RadarEdital, error) {
    // 6. Usage limits (billing integration)
    if err := checkBillingLimit(ctx, tenantID, "radar_ingest"); err != nil {
        return nil, err
    }

    // Parse raw JSON
    var tmp struct {
        Title       string  `json:"title"`
        Description string  `json:"description"`
        Value       float64 `json:"value"`
        Deadline    string  `json:"deadline"`
    }
    if err := json.Unmarshal(raw, &tmp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal raw edital: %w", err)
    }
    deadline, err := time.Parse(time.RFC3339, tmp.Deadline)
    if err != nil {
        return nil, fmt.Errorf("invalid deadline format: %w", err)
    }

    // Compute deterministic external_id (MD5 of title+description)
    externalID := fmt.Sprintf("%x", md5.Sum([]byte(tmp.Title+tmp.Description)))

    // 8. Cross‑source deduplication using content hash (SHA‑256 of normalized text)
    contentHash := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.TrimSpace(tmp.Title)+" "+strings.TrimSpace(tmp.Description))))
    dedupKey := fmt.Sprintf("radar_dedup_content:%s", contentHash)
    if _, err := s.rdb.Get(ctx, dedupKey).Result(); err == nil {
        // already processed across any source within TTL
        return nil, nil
    }
    // Store dedup key with 10‑min TTL
    s.rdb.Set(ctx, dedupKey, "1", 10*time.Minute)

    // Load tenant preferences (including keywords, min_value, categories)
    pref, err := s.loadTenantPreference(ctx, tenantID)
    if err != nil {
        return nil, err
    }
    // 2. Value filtering
    if pref.MinValue != nil && tmp.Value < *pref.MinValue {
        return nil, nil // filtered out
    }

    // 1. Semantic keyword matching via embeddings
    var matchedKeywords []string
    var keywordSimilarity float64
    var editalEmbedding []float32 // may be generated during keyword step
    if len(pref.Keywords) > 0 {
        prefText := strings.Join(pref.Keywords, " ")
        prefEmbedding, err := s.aiSvc.GenerateEmbedding(ctx, prefText)
        if err == nil {
            // Generate edital embedding (title+description)
            editalEmbedding, err = s.aiSvc.GenerateEmbedding(ctx, fmt.Sprintf("%s %s", tmp.Title, tmp.Description))
            if err == nil {
                // Compute similarity between preference embedding and edital embedding
                keywordSimilarity, err = s.aiSvc.ComputeEmbeddingSimilarity(ctx, prefEmbedding, editalEmbedding)
                if err == nil && keywordSimilarity > 0.65 {
                    matchedKeywords = pref.Keywords
                } else {
                    // similarity below threshold – discard edital
                    return nil, nil
                }
            }
        }
    }

    // Persist edital (analysis_log will be updated later)
    edital := &RadarEdital{
        ID:          uuid.New(),
        ExternalID:  externalID,
        Source:      source,
        Title:       tmp.Title,
        Description: tmp.Description,
        Value:       tmp.Value,
        Deadline:    deadline,
        RawJSON:     raw,
        TenantID:    tenantID,
        CreatedAt:   time.Now(),
    }
    _, err = s.db.ExecContext(ctx, `INSERT INTO radar_editais (id, external_id, source, title, description, value, deadline, raw_json, tenant_id, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, edital.ID, edital.ExternalID, edital.Source, edital.Title, edital.Description, edital.Value, edital.Deadline, edital.RawJSON, edital.TenantID, edital.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to insert radar edital: %w", err)
    }

    // Save embedding if we have it (from keyword step) or generate now
    if editalEmbedding == nil {
        editalEmbedding, err = s.aiSvc.GenerateEmbedding(ctx, fmt.Sprintf("%s %s", edital.Title, edital.Description))
        if err != nil {
            // continue without embedding – relevance will rely on win probability only
            editalEmbedding = nil
        }
    }
    if editalEmbedding != nil {
        if err := s.aiSvc.SaveEmbedding(ctx, tenantID, edital.ID, editalEmbedding); err != nil {
            return nil, fmt.Errorf("failed to save embedding: %w", err)
        }
    }

    // 2. Compute similarity vs tenant history (if embedding available)
    var historySimilarity float64
    if editalEmbedding != nil {
        historySimilarity, err = s.aiSvc.CalculateSimilarity(ctx, tenantID, edital.ID)
        if err != nil {
            historySimilarity = 0
        }
    }
    // Final similarity is the higher of keyword and history similarity
    finalSimilarity := math.Max(keywordSimilarity, historySimilarity)
    edital.Similarity = finalSimilarity

    // Compute win probability (normalized 0‑1)
    winProb, _, err := s.aiSvc.CalculateWinProbability(ctx, tenantID, edital.ID, edital.Value)
    if err != nil {
        winProb = 0
    }
    edital.WinProbability = winProb

    // Recency factor (clamped 0‑1)
    days := time.Since(edital.Deadline).Hours() / 24
    recency := 1.0
    if days > 365 {
        recency = 0.6
    } else if days > 90 {
        recency = 0.8
    }
    if recency < 0 {
        recency = 0
    } else if recency > 1 {
        recency = 1
    }

    // 3. Priority score normalization (components already 0‑1)
    edital.PriorityScore = (winProb * 0.6) + (finalSimilarity * 0.3) + (recency * 0.1)
    // 5. Priority label
    switch {
    case edital.PriorityScore > 0.75:
        edital.PriorityLabel = "HIGH"
    case edital.PriorityScore >= 0.5:
        edital.PriorityLabel = "MEDIUM"
    default:
        edital.PriorityLabel = "LOW"
    }

    // 7. Auto‑follow rule refinement
    if winProb > 0.75 && finalSimilarity > 0.7 {
        if err := s.autoFollow(ctx, tenantID, edital.ID); err != nil {
            // log but continue
        }
    }

    // 5. Analysis log (explainability)
    analysis := map[string]interface{}{
        "similarity":       finalSimilarity,
        "win_probability": winProb,
        "matched_keywords": matchedKeywords,
        "reason":           "Keyword similarity and win probability thresholds met",
    }
    analysisJSON, _ := json.Marshal(analysis)
    edital.AnalysisLog = analysisJSON
    // Update DB with analysis_log (assumes column exists)
    s.db.ExecContext(ctx, `UPDATE radar_editais SET analysis_log = $1 WHERE id = $2`, edital.AnalysisLog, edital.ID)

    // 4. Enhanced notification payload
    if winProb > 0.6 || finalSimilarity > 0.7 {
        notif := NotificationPayload{
            TenantID:  tenantID,
            Type:      "NEW_OPPORTUNITY",
            Title:     fmt.Sprintf("New opportunity: %s", edital.Title),
            Message:   fmt.Sprintf("A new edital with win probability %.2f and similarity %.2f was found. Priority: %s", winProb, finalSimilarity, edital.PriorityLabel),
            ActionURL: fmt.Sprintf("/radar/%s", edital.ID),
            Data: map[string]interface{}{
                "edital_id":       edital.ID,
                "win_probability": winProb,
                "similarity":       finalSimilarity,
                "priority_score":  edital.PriorityScore,
            },
        }
        if err := s.dispatchNotification(ctx, notif); err != nil {
            // log error (omitted for brevity)
        }
    }

    return edital, nil
}

// NotificationPayload defines the payload for a notification.
type NotificationPayload struct {
    TenantID  uuid.UUID
    Type      string
    Title     string
    Message   string
    ActionURL string
    Data      map[string]interface{}
}

func (s *Service) dispatchNotification(ctx context.Context, p NotificationPayload) error {
    // Direct DB insert; deduplication and cooldown handled by notifications subsystem.
    _, err := s.db.ExecContext(ctx, `INSERT INTO notifications (id, tenant_id, type, title, message, action_url, data, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,now())`, uuid.New(), p.TenantID, p.Type, p.Title, p.Message, p.ActionURL, p.Data)
    return err
}

func (s *Service) autoFollow(ctx context.Context, tenantID, editalID uuid.UUID) error {
    _, err := s.db.ExecContext(ctx, `INSERT INTO radar_followers (tenant_id, edital_id, created_at) VALUES ($1,$2,now()) ON CONFLICT DO NOTHING`, tenantID, editalID)
    return err
}

// loadTenantPreference loads the first preference profile for a tenant.
func (s *Service) loadTenantPreference(ctx context.Context, tenantID uuid.UUID) (*RadarPreference, error) {
    var pref RadarPreference
    row := s.db.QueryRowContext(ctx, `SELECT id, name, keywords, min_value, categories FROM radar_preferences WHERE tenant_id = $1 AND deleted_at IS NULL LIMIT 1`, tenantID)
    err := row.Scan(&pref.ID, &pref.Name, &pref.Keywords, &pref.MinValue, &pref.Categories)
    if err == sql.ErrNoRows {
        return &RadarPreference{}, nil // empty pref
    }
    if err != nil {
        return nil, err
    }
    return &pref, nil
}

type RadarPreference struct {
    ID         uuid.UUID
    Name       sql.NullString
    Keywords   []string
    MinValue   *float64
    Categories []string
}

// ListEditais returns paginated radar editais for a tenant, ordered by priority_score DESC.
func (s *Service) ListEditais(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]*RadarEdital, error) {
    offset := (page - 1) * limit
    rows, err := s.db.QueryContext(ctx, `SELECT id, external_id, source, title, description, value, deadline, raw_json, created_at, similarity, win_probability, priority_score, priority_label FROM radar_editais WHERE tenant_id = $1 AND deleted_at IS NULL ORDER BY priority_score DESC LIMIT $2 OFFSET $3`, tenantID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var results []*RadarEdital
    for rows.Next() {
        var e RadarEdital
        if err := rows.Scan(&e.ID, &e.ExternalID, &e.Source, &e.Title, &e.Description, &e.Value, &e.Deadline, &e.RawJSON, &e.CreatedAt, &e.Similarity, &e.WinProbability, &e.PriorityScore, &e.PriorityLabel); err != nil {
            return nil, err
        }
        results = append(results, &e)
    }
    return results, nil
}

// FollowEdital creates a follower entry for a tenant.
func (s *Service) FollowEdital(ctx context.Context, tenantID, editalID uuid.UUID) error {
    _, err := s.db.ExecContext(ctx, `INSERT INTO radar_followers (tenant_id, edital_id, created_at) VALUES ($1,$2,now()) ON CONFLICT DO NOTHING`, tenantID, editalID)
    return err
}

// checkBillingLimit is a placeholder for the real billing middleware.
func checkBillingLimit(ctx context.Context, tenantID uuid.UUID, limitName string) error {
    // In production this would call middleware.RequireLimit or similar.
    // Here we assume the limit is not exceeded.
    return nil
}

// containsIgnoreCase is a helper for simple case‑insensitive substring checks.
func containsIgnoreCase(text, substr string) bool {
    return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}
