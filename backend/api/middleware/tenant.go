package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TenantIDKey contextKey = "tenant_id"
)

// RequireTenant returns a valid mux.MiddlewareFunc that ensures the request belongs to a valid tenant.
func RequireTenant(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract tenant ONLY from authenticated user (JWT)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"invalid_tenant", "message":"Missing authorization token"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			// Since the secret is validated in RequireAuth, we can parse unverified to read claims safely,
			// or we can expect it in the context if auth middleware passed it.
			// Let's parse unverified to pull out tenant_id without needing the secret again.
			token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
			if err != nil {
				http.Error(w, `{"error":"invalid_tenant", "message":"Invalid token format"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"invalid_tenant", "message":"Invalid token claims"}`, http.StatusUnauthorized)
				return
			}

			tenantIDStr, ok := claims["tenant_id"].(string)
			if !ok || tenantIDStr == "" {
				http.Error(w, `{"error":"invalid_tenant", "message":"Tenant ID not found in token"}`, http.StatusUnauthorized)
				return
			}

			tenantID, err := uuid.Parse(tenantIDStr)
			if err != nil {
				http.Error(w, `{"error":"invalid_tenant", "message":"Malformed Tenant ID"}`, http.StatusUnauthorized)
				return
			}

			// Validate tenant existence in database
			var exists bool
			err = db.QueryRowContext(r.Context(), "SELECT EXISTS(SELECT 1 FROM tenants WHERE id = $1)", tenantID).Scan(&exists)
			if err != nil || !exists {
				http.Error(w, `{"error":"invalid_tenant", "message":"Tenant does not exist"}`, http.StatusUnauthorized)
				return
			}

			// Inject tenant into context
			ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetTenantID extracts the Tenant ID from the context. It expects it to exist.
func GetTenantID(ctx context.Context) string {
	if val, ok := ctx.Value(TenantIDKey).(uuid.UUID); ok {
		return val.String()
	}
	return ""
}

