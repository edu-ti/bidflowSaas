package radar

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "net/http"
    "strconv"
    "github.com/google/uuid"
)

// RegisterRadarRoutes registers HTTP handlers for the Radar module.
func RegisterRadarRoutes(mux *http.ServeMux, svc *Service) {
    mux.HandleFunc("/api/v1/radar", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        handleListEditais(w, r, svc)
    })
    mux.HandleFunc("/api/v1/radar/follow/", func(w http.ResponseWriter, r *http.Request) {
        // Expect URL pattern: /api/v1/radar/follow/{id}
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        // Extract ID after the prefix.
        idStr := r.URL.Path[len("/api/v1/radar/follow/"):]
        editalID, err := uuid.Parse(idStr)
        if err != nil {
            http.Error(w, "invalid edital id", http.StatusBadRequest)
            return
        }
        // In a real request, tenant ID would be derived from auth context; using placeholder.
        tenantID := uuid.Nil
        if err := svc.FollowEdital(r.Context(), tenantID, editalID); err != nil {
            slog.Error("follow failed", "err", err)
            http.Error(w, "failed to follow", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    })
    mux.HandleFunc("/api/v1/radar/", func(w http.ResponseWriter, r *http.Request) {
        // Expect URL pattern: /api/v1/radar/{id}
        if r.Method != http.MethodGet {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        idStr := r.URL.Path[len("/api/v1/radar/"):]
        editalID, err := uuid.Parse(idStr)
        if err != nil {
            http.Error(w, "invalid edital id", http.StatusBadRequest)
            return
        }
        // Retrieve single edital – service does not have GetEdital yet, use DB query directly.
        // For brevity, we reuse ListEditais with limit=1 and filter.
        // In production, implement a dedicated GetEdital method.
        // Here we just respond with not implemented.
        http.Error(w, "not implemented", http.StatusNotImplemented)
    })
}

// handleListEditais handles GET /api/v1/radar with pagination.
func handleListEditais(w http.ResponseWriter, r *http.Request, svc *Service) {
    q := r.URL.Query()
    page, _ := strconv.Atoi(q.Get("page"))
    if page < 1 {
        page = 1
    }
    limit, _ := strconv.Atoi(q.Get("limit"))
    if limit < 1 {
        limit = 20
    }
    // Placeholder tenant ID.
    tenantID := uuid.Nil
    editais, err := svc.ListEditais(r.Context(), tenantID, page, limit)
    if err != nil {
        slog.Error("list editais failed", "err", err)
        http.Error(w, "failed to list", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(editais); err != nil {
        slog.Error("encode response failed", "err", err)
    }
}
