package ai

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"lastsaas/api/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type AnalyzeRequest struct {
	EditalID string `json:"edital_id"`
	Text     string `json:"text"`
}

func (h *Handler) AnalyzeEdital(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := middleware.GetTenantID(r.Context())
	tenantID, _ := uuid.Parse(tenantIDStr)

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_payload"}`, http.StatusBadRequest)
		return
	}

	editalID, err := uuid.Parse(req.EditalID)
	if err != nil {
		http.Error(w, `{"error":"invalid_edital_id"}`, http.StatusBadRequest)
		return
	}

	// Dispatch asynchronous execution bypassing direct blocking limitations
	err = h.svc.AnalyzeEdital(r.Context(), tenantID, editalID, req.Text)
	if err != nil {
		http.Error(w, `{"error":"dispatch_failed","message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"pending","message":"AI analysis dispatched to asynchronous workers. Checking billing limits."}`))
}

func (h *Handler) GetResults(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := middleware.GetTenantID(r.Context())
	tenantID, _ := uuid.Parse(tenantIDStr)

	vars := mux.Vars(r)
	editalID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"invalid_edital_id"}`, http.StatusBadRequest)
		return
	}

	res, err := h.svc.GetAnalysisStatus(r.Context(), tenantID, editalID)
	if err != nil {
		http.Error(w, `{"error":"not_found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func RegisterRoutes(r *mux.Router, h *Handler, db *sql.DB) {
	// Limit enforcement intercepts exactly at the route border
	r.Handle("/analyze", middleware.RequireLimit(db, "ai_analysis")(http.HandlerFunc(h.AnalyzeEdital))).Methods(http.MethodPost)
	r.HandleFunc("/results/{id}", h.GetResults).Methods(http.MethodGet)
}

