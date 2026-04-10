package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"lastsaas/api/middleware"
)

type AIHandler struct {
	db *sql.DB
}

func NewAIHandler(db *sql.DB) *AIHandler {
	return &AIHandler{db: db}
}

func (h *AIHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "GetInsights OK for tenant ` + tenantID.String() + `."}`))
}
