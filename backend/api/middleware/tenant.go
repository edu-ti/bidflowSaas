package middleware

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const TenantKey contextKey = "tenant_id"
const UserIDKey contextKey = "user_id"

func RequireTenant(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Safe context extraction for User setup by Auth middleware
			userVal := r.Context().Value(UserIDKey)
			if userVal == nil {
				slog.Warn("Missing user context attempting tenant resolution")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "unauthorized", "message": "missing user context"}`))
				return
			}

			_, ok := userVal.(uuid.UUID)
			if !ok {
				slog.Error("Unsafe User type assertion failure")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "invalid_user_context", "message": "invalid user type in context"}`))
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"invalid_tenant", "message":"Missing authorization token"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
			if err != nil {
				slog.Warn("Invalid token format for tenant extraction")
				http.Error(w, `{"error":"invalid_tenant", "message":"Invalid token format"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				slog.Warn("Could not parse map claims for tenant extraction")
				http.Error(w, `{"error":"invalid_tenant", "message":"Invalid token claims"}`, http.StatusUnauthorized)
				return
			}

			tenantIDStr, ok := claims["tenant_id"].(string)
			if !ok || tenantIDStr == "" {
				slog.Warn("Invalid tenant attempt: Tenant ID missing in JWT")
				http.Error(w, `{"error":"invalid_tenant", "message":"Tenant not found or unauthorized"}`, http.StatusUnauthorized)
				return
			}

			tenantID, err := uuid.Parse(tenantIDStr)
			if err != nil {
				slog.Warn("Malformed Tenant ID spoofing attempt", "tenantId", tenantIDStr)
				http.Error(w, `{"error":"invalid_tenant", "message":"Malformed Tenant ID"}`, http.StatusUnauthorized)
				return
			}

			// Validate tenant in database before allowing request:
			var id uuid.UUID
			err = db.QueryRowContext(r.Context(), "SELECT id FROM tenants WHERE id = $1 AND active = true;", tenantID).Scan(&id)
			if err != nil {
				slog.Warn("Tenant validation failed or inactive mapping", "tenantId", tenantIDStr, "error", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "invalid_tenant", "message": "Tenant not found or inactive"}`))
				return
			}

			ctx := context.WithValue(r.Context(), TenantKey, tenantIDStr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth extracts the user context and sets it
func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"unauthorized", "message":"Missing token"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"unauthorized", "message":"Invalid signature"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"unauthorized", "message":"Invalid claims"}`, http.StatusUnauthorized)
				return
			}

			userIdStr, _ := claims["sub"].(string)
			userId, err := uuid.Parse(userIdStr)
			if err != nil {
				http.Error(w, `{"error":"unauthorized", "message":"Invalid user ID"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetTenantID helper function
func GetTenantID(ctx context.Context) string {
	if val, ok := ctx.Value(TenantKey).(string); ok {
		return val
	}
	return ""
}




