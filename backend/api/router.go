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

	// 1. API Versioning
	api := r.PathPrefix("/api/v1").Subrouter()

	// Public routes
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Secured routes (Requires Auth)
	secured := api.PathPrefix("").Subrouter()
	secured.Use(middleware.RequireAuth(cfg.JWTSecret))

	// Tenant-scoped routes
	tenantRoutes := secured.PathPrefix("/t").Subrouter()
	tenantRoutes.Use(middleware.RequireTenant(cfg.DB))

	// CRM Module Isolation
	crmRouter := tenantRoutes.PathPrefix("/crm").Subrouter()
	crmRouter.Use(middleware.RequireModule(cfg.BillingSvc, "crm"))
	crm.RegisterRoutes(crmRouter, cfg.CRMHandler)

	// Licitacoes Module Isolation
	licitacoesRouter := tenantRoutes.PathPrefix("/licitacoes").Subrouter()
	licitacoesRouter.Use(middleware.RequireModule(cfg.BillingSvc, "licitacoes"))
	licitacoes.RegisterRoutes(licitacoesRouter, cfg.LicHandler)

	// AI Module Isolation
	aiRouter := tenantRoutes.PathPrefix("/ai").Subrouter()
	aiRouter.Use(middleware.RequireModule(cfg.BillingSvc, "ai"))
	ai.RegisterRoutes(aiRouter, cfg.AIHandler)

	return r
}


