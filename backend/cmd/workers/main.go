package main

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/lib/pq"

	"lastsaas/workers"
)

func main() {
	dbURL := getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/bidflow?sslmode=disable")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")

	slog.Info("Starting BidFlow workers", "redis", redisAddr)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("Database ping failed", "error", err)
		os.Exit(1)
	}

	processor := workers.NewProcessor(redisAddr, db)

	slog.Info("Workers started – listening for tasks on queues: critical, default, low")
	if err := processor.Start(); err != nil {
		slog.Error("Worker processor failed", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
