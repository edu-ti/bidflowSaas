package main

import (
	"log/slog"
	"os"

	"lastsaas/app"
	"lastsaas/internal/config"
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

	application, err := app.New(cfg)
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		slog.Error("Failed to run application", "error", err)
		os.Exit(1)
	}
}
