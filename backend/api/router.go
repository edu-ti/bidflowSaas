package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"lastsaas/api/middleware"
	"lastsaas/core/auth"
	"lastsaas/core/billing"
	"lastsaas/core/tenant"
	"lastsaas/modules/crm"
	"lastsaas/modules/licitacoes"
	"lastsaas/modules/ai"
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

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Public routes
	v1.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Guarded routes require JWT
	guarded := v1.PathPrefix("").Subrouter()
	guarded.Use(middleware.RequireAuth(cfg.JWTSecret))

	// Tenant scoped routes require Tenant-ID header
	tenantRoutes := guarded.PathPrefix("/t").Subrouter()
	tenantRoutes.Use(middleware.RequireTenant)
    
	// Sub-modules
	crmRouter := tenantRoutes.PathPrefix("/crm").Subrouter()
	crmRouter.Use(middleware.RequireModule(cfg.BillingSvc, "crm"))
	
	if cfg.CRMHandler != nil {
		crmRouter.HandleFunc("/customers", cfg.CRMHandler.ListCustomers).Methods("GET")
		crmRouter.HandleFunc("/customers", cfg.CRMHandler.CreateCustomer).Methods("POST")
		crmRouter.HandleFunc("/opportunities", cfg.CRMHandler.ListOpportunities).Methods("GET")
	}

	licitacoesRouter := tenantRoutes.PathPrefix("/licitacoes").Subrouter()
	licitacoesRouter.Use(middleware.RequireModule(cfg.BillingSvc, "licitacoes"))

	if cfg.LicHandler != nil {
		licitacoesRouter.HandleFunc("/editais", cfg.LicHandler.ListEditais).Methods("GET")
		licitacoesRouter.HandleFunc("/editais", cfg.LicHandler.CreateEdital).Methods("POST")
	}

	aiRouter := tenantRoutes.PathPrefix("/ai").Subrouter()
	aiRouter.Use(middleware.RequireModule(cfg.BillingSvc, "ai"))
	if cfg.AIHandler != nil {
		aiRouter.HandleFunc("/process-edital", cfg.AIHandler.ProcessEdital).Methods("POST")
	}

	return r
}
