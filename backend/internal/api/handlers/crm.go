package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"lastsaas/api/middleware"
)

type CRMHandler struct {
	db *sql.DB
}

func NewCRMHandler(db *sql.DB) *CRMHandler {
	return &CRMHandler{db: db}
}

func (h *CRMHandler) ListFunnels(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "ListFunnels OK for tenant ` + tenantID.String() + `."}`))
}
