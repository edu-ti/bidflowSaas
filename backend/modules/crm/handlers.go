package crm

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

func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())

	customers, err := h.svc.ListCustomers(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "Failed to list customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h *Handler) ListOpportunities(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())

	opportunities, err := h.svc.ListOpportunities(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "Failed to list opportunities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opportunities)
}

type CreateCustomerRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	DocumentID string `json:"document_id"`
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.GetTenantID(r.Context())
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	customer, err := h.svc.CreateCustomer(r.Context(), tenantID, req.Name, req.Email, req.Phone, req.DocumentID, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}
