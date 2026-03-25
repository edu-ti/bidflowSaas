package main

import (
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

func main() {
	config.LoadEnvFile()

	env := config.GetEnv()
	cfg, err := config.Load(env)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	slog.Info("Starting SaaS Platform", "mode", cfg.Environment)

	// In real app, load URI from config, for now hardcoded default
	dbURI := os.Getenv("DATABASE_URL")
	if dbURI == "" {
		dbURI = "postgres://postgres:postgres@localhost:5432/bidflow?sslmode=disable"
	}

	database, err := db.Connect(dbURI)
	if err != nil {
		slog.Error("Failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Run migrations
	workingDir, _ := os.Getwd()
	migrationsPath := filepath.Join(workingDir, "database", "migrations")
	if err := database.RunMigrations(migrationsPath, "bidflow"); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
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

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server listening on port", "port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		slog.Error("Server shutdown", "error", err)
		os.Exit(1)
	}
}
