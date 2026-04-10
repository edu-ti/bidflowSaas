package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"lastsaas/api/middleware"
	"lastsaas/core/auth"
	"lastsaas/core/billing"
	"lastsaas/core/tenant"
	"lastsaas/internal/api/handlers"
)

type RouterConfig struct {
	DB         *sql.DB
	JWTSecret  string
	AuthSvc    *auth.Service
	TenantSvc  *tenant.Service
	BillingSvc *billing.Service
	CRMHandler *handlers.CRMHandler
	BidsHandler *handlers.BidsHandler
	AIHandler  *handlers.AIHandler
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := mux.NewRouter()

	// 1. Enforce API Versioning
	api := r.PathPrefix("/api/v1").Subrouter()

	// Public routes
	
	// Stripe Webhooks (bypass auth and general restrictive rate limits, requires external signature)
	api.HandleFunc("/webhooks/stripe", billing.HandleStripeWebhook(cfg.DB)).Methods(http.MethodPost)
	
	// Stricter rate limit solely for login (e.g. 5 reqs/sec, 5 burst)
	loginRouter := api.PathPrefix("/auth").Subrouter()
	loginRouter.Use(middleware.RateLimit(5.0, 5))
	loginRouter.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// placeholder
	}).Methods(http.MethodPost)
	
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)


	// Build shared middleware instances
	authMiddleware := middleware.RequireAuth(cfg.JWTSecret)
	tenantMiddleware := middleware.RequireTenant(cfg.DB)
	
	// General higher rate limit for standard authenticated routes
	apiRateLimit := middleware.RateLimit(100.0, 50)

	// CRM Module Isolation
	crmRouter := api.PathPrefix("/crm").Subrouter()
	crmRouter.Use(apiRateLimit)
	crmRouter.Use(authMiddleware)
	crmRouter.Use(tenantMiddleware)
	crmRouter.HandleFunc("/funnels", cfg.CRMHandler.ListFunnels).Methods(http.MethodGet)

	// Bids Module Isolation
	bidsRouter := api.PathPrefix("/bids").Subrouter()
	bidsRouter.Use(apiRateLimit)
	bidsRouter.Use(authMiddleware)
	bidsRouter.Use(tenantMiddleware)
	bidsRouter.HandleFunc("", cfg.BidsHandler.ListBids).Methods(http.MethodGet)
	bidsRouter.HandleFunc("", cfg.BidsHandler.CreateBid).Methods(http.MethodPost)

	// AI Module Isolation
	aiRouter := api.PathPrefix("/ai").Subrouter()
	aiRouter.Use(apiRateLimit)
	aiRouter.Use(authMiddleware)
	aiRouter.Use(tenantMiddleware)
	aiRouter.HandleFunc("/insights", cfg.AIHandler.GetInsights).Methods(http.MethodGet)

	return r
}






