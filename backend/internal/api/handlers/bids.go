package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"lastsaas/api/middleware"
	// "lastsaas/internal/db/queries" // This will be imported once generated
)

type BidsHandler struct {
	db *sql.DB
}

func NewBidsHandler(db *sql.DB) *BidsHandler {
	return &BidsHandler{db: db}
}

func (h *BidsHandler) ListBids(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// dbQueries := bidsqueries.New(h.db)
	// bids, err := dbQueries.ListBids(r.Context(), tenantID)
	// if err != nil {
	// 	http.Error(w, "Failed to list bids", http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(bids)
	w.Write([]byte(`{"message": "ListBids OK. Please run sqlc generate to enable DB code."}`))
}

func (h *BidsHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// var params bidsqueries.CreateBidParams
	// ... bind request ...
	// dbQueries := bidsqueries.New(h.db)
	// bid, err := dbQueries.CreateBid(r.Context(), params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "CreateBid OK for tenant ` + tenantID.String() + `."}`))
}
