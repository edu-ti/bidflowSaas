package licitacoes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"lastsaas/api/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListEditais(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())

	editais, err := h.svc.ListEditais(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "Failed to list editais", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(editais)
}

type CreateEditalRequest struct {
	Number            string    `json:"number"`
	Agency            string    `json:"agency"`
	ObjectDescription string    `json:"object_description"`
	OpeningDate       time.Time `json:"opening_date"`
}

func (h *Handler) CreateEdital(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	var req CreateEditalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	edital, err := h.svc.CreateEdital(r.Context(), tenantID, req.Number, req.Agency, req.ObjectDescription, req.OpeningDate, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		http.Error(w, "Failed to create edital", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(edital)
}
