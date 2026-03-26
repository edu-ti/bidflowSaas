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

	// API Versioning
	api := r.PathPrefix("/api/v1").Subrouter()

	// Public routes
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Secured routes (Requires Auth)
	secured := api.PathPrefix("").Subrouter()
	secured.Use(middleware.RequireAuth(cfg.JWTSecret))

	// Tenant-scoped routes (Requires Tenant context extracted and validated)
	tenantRoutes := secured.PathPrefix("/t").Subrouter()
	tenantRoutes.Use(middleware.RequireTenant(cfg.DB))

	// CRM Module Isolation
	crmRouter := tenantRoutes.PathPrefix("/crm").Subrouter()
	crmRouter.Use(middleware.RequireModule(cfg.BillingSvc, "crm"))
	if cfg.CRMHandler != nil {
		crmRouter.HandleFunc("/customers", cfg.CRMHandler.ListCustomers).Methods(http.MethodGet)
		crmRouter.HandleFunc("/customers", cfg.CRMHandler.CreateCustomer).Methods(http.MethodPost)
		crmRouter.HandleFunc("/opportunities", cfg.CRMHandler.ListOpportunities).Methods(http.MethodGet)
	}

	// Licitacoes Module Isolation
	licitacoesRouter := tenantRoutes.PathPrefix("/licitacoes").Subrouter()
	licitacoesRouter.Use(middleware.RequireModule(cfg.BillingSvc, "licitacoes"))
	if cfg.LicHandler != nil {
		licitacoesRouter.HandleFunc("/editais", cfg.LicHandler.ListEditais).Methods(http.MethodGet)
		licitacoesRouter.HandleFunc("/editais", cfg.LicHandler.CreateEdital).Methods(http.MethodPost)
	}

	// AI Module Isolation
	aiRouter := tenantRoutes.PathPrefix("/ai").Subrouter()
	aiRouter.Use(middleware.RequireModule(cfg.BillingSvc, "ai"))
	if cfg.AIHandler != nil {
		aiRouter.HandleFunc("/process-edital", cfg.AIHandler.ProcessEdital).Methods(http.MethodPost)
	}

	return r
}

