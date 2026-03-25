package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	TenantIDKey contextKey = "tenant_id"
)

func RequireTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantIDStr := r.Header.Get("X-Tenant-ID")
		if tenantIDStr == "" {
			http.Error(w, "Tenant ID required via X-Tenant-ID header", http.StatusBadRequest)
			return
		}

		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			http.Error(w, "Invalid Tenant ID", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTenantID(ctx context.Context) uuid.UUID {
	if val, ok := ctx.Value(TenantIDKey).(uuid.UUID); ok {
		return val
	}
	return uuid.Nil
}
