package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"lastsaas/api"
	"lastsaas/core/auth"
	"lastsaas/core/billing"
	"lastsaas/core/db"
	"lastsaas/core/tenant"
	"lastsaas/internal/config"
	"lastsaas/modules/ai"
	"lastsaas/modules/crm"
	"lastsaas/modules/licitacoes"
)

type App struct {
	cfg      *config.Config
	database *db.PostgresDB
	router   http.Handler
}

// New creates a new App instance and initializes all dependencies
func New(cfg *config.Config) (*App, error) {
	// Initialize database
	dbURI := os.Getenv("DATABASE_URL")
	if dbURI == "" {
		dbURI = "postgres://postgres:postgres@localhost:5432/bidflow?sslmode=disable"
	}

	database, err := db.Connect(dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Run migrations
	workingDir, _ := os.Getwd()
	migrationsPath := filepath.Join(workingDir, "database", "migrations")
	if err := database.RunMigrations(migrationsPath, "bidflow"); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize Services
	authSvc := auth.NewService(database.Client)
	tenantSvc := tenant.NewService(database.Client)
	billingSvc := billing.NewService(database.Client)
	crmSvc := crm.NewService(database.Client)
	licSvc := licitacoes.NewService(database.Client)
	aiSvc := ai.NewService("localhost:6379") // default redis

	// Initialize Handlers
	crmHandler := crm.NewHandler(crmSvc)
	licHandler := licitacoes.NewHandler(licSvc)
	aiHandler := ai.NewHandler(aiSvc)

	// Initialize Router
	routerCfg := api.RouterConfig{
		DB:         database.Client,
		JWTSecret:  cfg.JWT.AccessSecret,
		AuthSvc:    authSvc,
		TenantSvc:  tenantSvc,
		BillingSvc: billingSvc,
		CRMHandler: crmHandler,
		LicHandler: licHandler,
		AIHandler:  aiHandler,
	}

	router := api.NewRouter(routerCfg)

	return &App{
		cfg:      cfg,
		database: database,
		router:   router,
	}, nil
}

// Run starts the application server
func (a *App) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", a.cfg.Server.Port)
	}
	if port == "0" {
		port = "8080"
	}

	slog.Info("Server listening on port", "port", port)
	
	// Ensure cleanup when server stops
	defer a.database.Close()
	
	if err := http.ListenAndServe(":"+port, a.router); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}
	
	return nil
}
