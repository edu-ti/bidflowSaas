package middleware

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// RequireLimit checks exactly against plan capacities extracting tenant identity
func RequireLimit(db *sql.DB, resource string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantIDStr := GetTenantID(r.Context())
			tenantID, err := uuid.Parse(tenantIDStr)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "invalid_tenant", "message": "Tenant not found or inactive"}`))
				return
			}
			
			var limitValue int
			var currentUsage int
			
			err = db.QueryRowContext(r.Context(), `
				SELECT ul.limit_value, COALESCE(ut.current_usage, 0)
				FROM subscriptions s
				JOIN usage_limits ul ON s.plan_id = ul.plan_id AND ul.resource = $2
				LEFT JOIN usage_tracking ut ON ut.tenant_id = s.tenant_id AND ut.resource = $2
				WHERE s.tenant_id = $1 AND s.status = 'active'
			`, tenantID, resource).Scan(&limitValue, &currentUsage)
			
			if err != nil {
				slog.Warn("Limit validation missing or subscription failed", "tenant", tenantID, "resource", resource)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "limit_exceeded", "message": "You have reached your plan limit"}`))
				return
			}

			if currentUsage >= limitValue {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "limit_exceeded", "message": "You have reached your plan limit"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
