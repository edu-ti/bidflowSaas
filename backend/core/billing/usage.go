package billing

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// IncrementUsage forces an atomic upscale mapping resources to SaaS tenant capabilities with rollover bounds.
func IncrementUsage(ctx context.Context, db *sql.DB, tenantID uuid.UUID, resource string) error {
	now := time.Now()
	// Next period resets in 1 month seamlessly simulating billing cycle resets natively directly from Postgres metrics.
	nextPeriodEnd := now.AddDate(0, 1, 0)
	
	query := `
		INSERT INTO usage_tracking (tenant_id, resource, current_usage, period_start, period_end)
		VALUES ($1, $2, 1, $3, $4)
		ON CONFLICT (tenant_id, resource) DO UPDATE SET
			current_usage = CASE 
				WHEN usage_tracking.period_end < EXCLUDED.period_start THEN 1 
				ELSE usage_tracking.current_usage + 1 
			END,
			period_start = CASE 
				WHEN usage_tracking.period_end < EXCLUDED.period_start THEN EXCLUDED.period_start 
				ELSE usage_tracking.period_start 
			END,
			period_end = CASE 
				WHEN usage_tracking.period_end < EXCLUDED.period_start THEN EXCLUDED.period_end 
				ELSE usage_tracking.period_end 
			END
	`
	
	_, err := db.ExecContext(ctx, query, tenantID, resource, now, nextPeriodEnd)
	if err != nil {
		slog.Error("Failed to increment SaaS billing usage tracking", "tenant", tenantID, "resource", resource, "err", err)
	}
	return err
}
