package radar

import (
    "context"
    "fmt"
    "log/slog"
    "github.com/hibiken/asynq"
    "github.com/google/uuid"
    "github.com/redis/go-redis/v9"
)

// Asynq task type for radar fetching.
const TypeRadarFetch = "radar:fetch"

// RegisterRadarTasks registers the radar fetch task handler with the Asynq server.
func RegisterRadarTasks(server *asynq.Server, svc *Service) {
    mux := asynq.NewServeMux()
    mux.HandleFunc(TypeRadarFetch, func(ctx context.Context, t *asynq.Task) error {
        return handleRadarFetch(ctx, svc)
    })
    server.Start(mux)
}

// EnqueueRadarFetch enqueues a radar fetch task (used by a cron scheduler).
func EnqueueRadarFetch(client *asynq.Client) error {
    task := asynq.NewTask(TypeRadarFetch, nil)
    _, err := client.Enqueue(task)
    return err
}

// handleRadarFetch performs the full fetch‑parse‑ingest pipeline.
func handleRadarFetch(ctx context.Context, svc *Service) error {
    for _, src := range sources {
        slog.Info("Fetching editais", "source", src.Name)
        rawList, err := src.Fetch(ctx)
        if err != nil {
            slog.Error("Failed to fetch source", "source", src.Name, "err", err)
            continue // proceed with other sources
        }
        for _, raw := range rawList {
            // Deduplication based on external ID (computed in parser).
            edital, err := src.Parse(raw)
            if err != nil {
                slog.Error("Parse error", "source", src.Name, "raw_id", raw.ID, "err", err)
                continue
            }
            // Check Redis dedup key.
            dedup, err := IsDeduped(ctx, svc.rdb, src.Name, edital.ExternalID)
            if err != nil {
                slog.Error("Redis dedup check failed", "err", err)
                continue
            }
            if dedup {
                slog.Info("Skipping duplicate edital", "source", src.Name, "external_id", edital.ExternalID)
                continue
            }
            // Record dedup marker.
            if err := RecordDedup(ctx, svc.rdb, src.Name, edital.ExternalID); err != nil {
                slog.Error("Failed to set dedup key", "err", err)
            }
            // Ingest via service (tenant ID must be resolved – for now we assume a single tenant placeholder).
            // In production, the source would provide tenant context; here we use a dummy UUID.
            tenantID := uuid.Nil // replace with actual tenant resolution logic.
            if _, err := svc.IngestEdital(ctx, tenantID, raw.Data, src.Name); err != nil {
                slog.Error("Ingest failed", "source", src.Name, "external_id", edital.ExternalID, "err", err)
                continue
            }
            slog.Info("Ingested edital", "source", src.Name, "external_id", edital.ExternalID)
        }
    }
    return nil
}

// CronRadarFetch is a helper that can be called from a scheduler (e.g., Asynq periodic task).
func CronRadarFetch(client *asynq.Client) error {
    return EnqueueRadarFetch(client)
}
