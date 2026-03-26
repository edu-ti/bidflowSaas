package middleware

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// In a real application, you might cache roles or embed them into the JWT claims.
// For strict RBAC, we query the DB to ensure roles haven't been forcefully revoked.
func RequireRole(db *sql.DB, requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userVal := r.Context().Value(UserIDKey)
			userID, ok := userVal.(uuid.UUID)
			if !ok {
				slog.Error("Unsafe User type assertion failure during RBAC")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "invalid_user_context", "message": "invalid user type in context"}`))
				return
			}

			tenantIDStr := GetTenantID(r.Context())
			tenantID, err := uuid.Parse(tenantIDStr)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "invalid_tenant", "message": "Tenant not found or inactive"}`))
				return
			}

			var role string
			err = db.QueryRowContext(r.Context(), "SELECT role FROM tenant_members WHERE user_id = $1 AND tenant_id = $2", userID, tenantID).Scan(&role)
			if err != nil || role != requiredRole {
				if role == "admin" && requiredRole != "admin" {
					// admin typically overrides other roles but let's strictly adhere or log
					// for this example, we'll allow admin to access anything required below it.
				} else if role != requiredRole {
					slog.Warn("Unauthorized role access attempt", "user", userID, "tenant", tenantID, "required", requiredRole, "actual", role)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{"error": "forbidden", "message": "Insufficient role permissions"}`))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

