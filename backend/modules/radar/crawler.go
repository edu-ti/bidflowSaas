package radar

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/google/uuid"
    "log/slog"
)

// RawEdital represents the raw JSON payload fetched from a source.
type RawEdital struct {
    ID   string          `json:"id"`
    Data json.RawMessage `json:"data"`
}

// SourceConfig defines a procurement source.
type SourceConfig struct {
    Name  string
    Fetch func(ctx context.Context) ([]RawEdital, error)
    Parse func(RawEdital) (*RadarEdital, error)
}

var sources []SourceConfig

func init() {
    // Register PNCP and ComprasNet sources.
    sources = []SourceConfig{
        {
            Name:  "PNCP",
            Fetch: fetchPNCP,
            Parse: parsePNCP,
        },
        {
            Name:  "ComprasNet",
            Fetch: fetchComprasNet,
            Parse: parseComprasNet,
        },
    }
}

// httpClient returns a client with a 5‑second timeout.
func httpClient() *http.Client {
    return &http.Client{Timeout: 5 * time.Second}
}

// fetchWithRetry performs an HTTP GET with up to 2 retries.
func fetchWithRetry(ctx context.Context, url string) ([]byte, error) {
    var lastErr error
    for i := 0; i < 3; i++ { // initial try + 2 retries
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        if err != nil {
            return nil, err
        }
        resp, err := httpClient().Do(req)
        if err != nil {
            lastErr = err
            continue
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
            lastErr = fmt.Errorf("unexpected status %d", resp.StatusCode)
            continue
        }
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            lastErr = err
            continue
        }
        return body, nil
    }
    return nil, fmt.Errorf("failed after retries: %w", lastErr)
}

// fetchPNCP fetches raw editais from the PNCP source.
func fetchPNCP(ctx context.Context) ([]RawEdital, error) {
    // Placeholder URL – replace with real endpoint.
    url := "https://api.pncp.gov.br/editais"
    body, err := fetchWithRetry(ctx, url)
    if err != nil {
        return nil, err
    }
    var list []RawEdital
    if err := json.Unmarshal(body, &list); err != nil {
        return nil, err
    }
    return list, nil
}

// fetchComprasNet fetches raw editais from the ComprasNet source.
func fetchComprasNet(ctx context.Context) ([]RawEdital, error) {
    url := "https://api.comprasnet.gov.br/editais"
    body, err := fetchWithRetry(ctx, url)
    if err != nil {
        return nil, err
    }
    var list []RawEdital
    if err := json.Unmarshal(body, &list); err != nil {
        return nil, err
    }
    return list, nil
}

// parsePNCP parses a RawEdital from PNCP into a RadarEdital.
func parsePNCP(raw RawEdital) (*RadarEdital, error) {
    // The raw data format may differ; delegate to shared parser.
    return parseGeneric(raw)
}

// parseComprasNet parses a RawEdital from ComprasNet.
func parseComprasNet(raw RawEdital) (*RadarEdital, error) {
    return parseGeneric(raw)
}

// parseGeneric normalizes the raw JSON using the parser module.
func parseGeneric(raw RawEdital) (*RadarEdital, error) {
    // Use the parser utilities to normalize fields.
    // The parser.go file provides NormalizeEdital.
    return NormalizeEdital(raw)
}

// DedupKey returns the Redis key used for cross‑source deduplication.
func DedupKey(source, externalID string) string {
    return fmt.Sprintf("radar_dedup:%s:%s", source, externalID)
}

// RecordDedup stores a deduplication marker with a 10‑minute TTL.
func RecordDedup(ctx context.Context, rdb *redis.Client, source, externalID string) error {
    key := DedupKey(source, externalID)
    return rdb.Set(ctx, key, "1", 10*time.Minute).Err()
}

// IsDeduped checks whether the given edital was already processed.
func IsDeduped(ctx context.Context, rdb *redis.Client, source, externalID string) (bool, error) {
    key := DedupKey(source, externalID)
    _, err := rdb.Get(ctx, key).Result()
    if err == redis.Nil {
        return false, nil
    }
    return err == nil, err
}

// LogInfo logs structured information.
func LogInfo(msg string, attrs ...any) {
    slog.Info(msg, attrs...)
}
