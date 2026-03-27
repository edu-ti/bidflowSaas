package ai

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

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

type FeedbackRequest struct {
	EditalID   string `json:"edital_id"`
	AIDecision string `json:"ai_decision"`
	RealResult string `json:"real_result"`
}

func (h *Handler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := middleware.GetTenantID(r.Context())
	tenantID, _ := uuid.Parse(tenantIDStr)

	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_payload"}`, http.StatusBadRequest)
		return
	}
	
	editalID, err := uuid.Parse(req.EditalID)
	if err != nil {
		http.Error(w, `{"error":"invalid_edital_id"}`, http.StatusBadRequest)
		return
	}

	err = h.svc.SubmitFeedback(r.Context(), tenantID, editalID, req.AIDecision, req.RealResult)
	if err != nil {
		http.Error(w, `{"error":"internal_error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func (h *Handler) GetInsights(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := middleware.GetTenantID(r.Context())
	tenantID, _ := uuid.Parse(tenantIDStr)

	timeframeStr := r.URL.Query().Get("timeframe")
	timeframe := 0
	if timeframeStr != "" {
		if t, err := strconv.Atoi(timeframeStr); err == nil {
			timeframe = t
		}
	}

	insights, err := h.svc.GetInsights(r.Context(), tenantID, timeframe)
	if err != nil {
		http.Error(w, `{"error":"internal_error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insights)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := middleware.GetTenantID(r.Context())
	tenantID, _ := uuid.Parse(tenantIDStr)

	history, err := h.svc.GetHistory(r.Context(), tenantID)
	if err != nil {
		http.Error(w, `{"error":"internal_error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (h *Handler) GetSimilar(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := middleware.GetTenantID(r.Context())
	tenantID, _ := uuid.Parse(tenantIDStr)

	vars := mux.Vars(r)
	editalID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"invalid_edital_id"}`, http.StatusBadRequest)
		return
	}

	similar, err := h.svc.GetSimilarEditais(r.Context(), tenantID, editalID)
	if err != nil {
		http.Error(w, `{"error":"internal_error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(similar)
}

func RegisterRoutes(r *mux.Router, h *Handler, db *sql.DB) {

	r.Handle("/analyze", middleware.RequireLimit(db, "ai_analysis")(http.HandlerFunc(h.AnalyzeEdital))).Methods(http.MethodPost)
	r.HandleFunc("/results/{id}", h.GetResults).Methods(http.MethodGet)
	r.HandleFunc("/feedback", h.SubmitFeedback).Methods(http.MethodPost)
	r.HandleFunc("/insights", h.GetInsights).Methods(http.MethodGet)
	r.HandleFunc("/history", h.GetHistory).Methods(http.MethodGet)
	r.HandleFunc("/similar/{id}", h.GetSimilar).Methods(http.MethodGet)
}
