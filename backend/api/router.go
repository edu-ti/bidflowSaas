package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"lastsaas/api/middleware"
	"lastsaas/core/auth"
	"lastsaas/core/billing"
	"lastsaas/core/tenant"
	"lastsaas/modules/ai"
	"lastsaas/modules/crm"
	"lastsaas/modules/licitacoes"
)

type RouterConfig struct {
	DB         *sql.DB
	JWTSecret  string
	AuthSvc    *auth.Service
	TenantSvc  *tenant.Service
	BillingSvc *billing.Service
	CRMHandler *crm.Handler
	LicHandler *licitacoes.Handler
	AIHandler  *ai.Handler
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := mux.NewRouter()

	// 1. Enforce API Versioning
	api := r.PathPrefix("/api/v1").Subrouter()

	// Rate Limiting on global API (e.g. 100 requests per second burst)
	api.Use(middleware.RateLimit(100.0, 50))

	// Public routes
	api.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
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

	// CRM Module Isolation - flat subgroup with strict exact Auth -> Tenant -> Module order
	crmRouter := api.PathPrefix("/crm").Subrouter()
	crmRouter.Use(authMiddleware)
	crmRouter.Use(tenantMiddleware)
	crmRouter.Use(middleware.RequireModule(cfg.DB, "crm"))
	crm.RegisterRoutes(crmRouter, cfg.CRMHandler)

	// Licitacoes Module Isolation
	licitacoesRouter := api.PathPrefix("/licitacoes").Subrouter()
	licitacoesRouter.Use(authMiddleware)
	licitacoesRouter.Use(tenantMiddleware)
	licitacoesRouter.Use(middleware.RequireModule(cfg.DB, "licitacoes"))
	licitacoes.RegisterRoutes(licitacoesRouter, cfg.LicHandler)

	// AI Module Isolation
	aiRouter := api.PathPrefix("/ai").Subrouter()
	aiRouter.Use(authMiddleware)
	aiRouter.Use(tenantMiddleware)
	aiRouter.Use(middleware.RequireModule(cfg.DB, "ai"))
	ai.RegisterRoutes(aiRouter, cfg.AIHandler)

	return r
}





