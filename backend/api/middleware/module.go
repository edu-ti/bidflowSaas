package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"lastsaas/core/billing"
)

func RequireModule(billingSvc *billing.Service, moduleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantIDStr := GetTenantID(r.Context())
			tenantID, err := uuid.Parse(tenantIDStr)
			if err != nil {
				http.Error(w, "Invalid Tenant ID", http.StatusUnauthorized)
				return
			}
			
			hasAccess, err := billingSvc.TenantHasModule(r.Context(), tenantID, moduleName)
			if err != nil {
				http.Error(w, "Internal Server Error while checking module access", http.StatusInternalServerError)
				return
			}
			if !hasAccess {
				http.Error(w, "Module not enabled in current tenant plan", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
