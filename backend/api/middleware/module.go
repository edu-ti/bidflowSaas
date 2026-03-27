package middleware

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

func RequireModule(db *sql.DB, moduleName string) func(http.Handler) http.Handler {
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
			
			var exists int
			err = db.QueryRowContext(r.Context(), `
				SELECT 1 
				FROM subscriptions s
				JOIN plan_modules pm ON s.plan_id = pm.plan_id
				JOIN modules m ON pm.module_id = m.id
				WHERE s.tenant_id = $1 AND m.name = $2 AND s.status = 'active'
			`, tenantID, moduleName).Scan(&exists)
			
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "module_not_enabled", "message": "Module not included in subscription plan"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}



