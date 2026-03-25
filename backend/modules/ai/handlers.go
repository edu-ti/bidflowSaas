package ai

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"lastsaas/api/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type ProcessEditalRequest struct {
	EditalID string `json:"edital_id"`
	FileURL  string `json:"file_url"`
}

func (h *Handler) ProcessEdital(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())

	var req ProcessEditalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	editalID, err := uuid.Parse(req.EditalID)
	if err != nil {
		http.Error(w, "Invalid Edital ID", http.StatusBadRequest)
		return
	}

	if err := h.svc.EnqueueEditalProcessing(r.Context(), tenantID, editalID, req.FileURL); err != nil {
		http.Error(w, "Failed to enqueue task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "queued"})
}
