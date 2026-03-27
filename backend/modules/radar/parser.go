package radar

import (
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "html"
    "log/slog"
    "strings"
    "time"
    "github.com/google/uuid"
)

// NormalizeEdital converts a RawEdital into a RadarEdital with cleaned fields and deterministic external_id.
func NormalizeEdital(raw RawEdital) (*RadarEdital, error) {
    // Assume raw.Data contains fields: title, description, value, deadline (RFC3339).
    var payload struct {
        Title       string  `json:"title"`
        Description string  `json:"description"`
        Value       float64 `json:"value"`
        Deadline    string  `json:"deadline"`
    }
    if err := json.Unmarshal(raw.Data, &payload); err != nil {
        return nil, fmt.Errorf("failed to unmarshal raw edital data: %w", err)
    }

    // Clean text: strip HTML, trim, lowercase for consistency.
    cleanTitle := strings.TrimSpace(html.UnescapeString(payload.Title))
    cleanDesc := strings.TrimSpace(html.UnescapeString(payload.Description))
    cleanTitle = strings.ToLower(cleanTitle)
    cleanDesc = strings.ToLower(cleanDesc)

    // Parse deadline.
    deadline, err := time.Parse(time.RFC3339, payload.Deadline)
    if err != nil {
        return nil, fmt.Errorf("invalid deadline format: %w", err)
    }

    // Generate external_id as SHA256 of concatenated normalized fields.
    hashInput := fmt.Sprintf("%s|%s|%f|%s", cleanTitle, cleanDesc, payload.Value, deadline.Format(time.RFC3339))
    externalID := fmt.Sprintf("%x", sha256.Sum256([]byte(hashInput)))

    edital := &RadarEdital{
        ID:          uuid.New(),
        ExternalID:  externalID,
        Title:       cleanTitle,
        Description: cleanDesc,
        Value:       payload.Value,
        Deadline:    deadline,
        RawJSON:     raw.Data,
        // TenantID will be set by the caller (service.IngestEdital).
        CreatedAt:   time.Now(),
    }

    slog.Info("Normalized edital", "external_id", externalID, "source", raw.ID)
    return edital, nil
}
