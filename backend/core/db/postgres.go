package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresDB struct {
	Client *sql.DB
}

func Connect(uri string) (*PostgresDB, error) {
	slog.Info("Connecting to PostgreSQL using pgx...")
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	// Wait for DB to be ready
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresDB{Client: db}, nil
}

func (db *PostgresDB) RunMigrations(migrationsPath string, dbName string) error {
	slog.Info("Running migrations", "path", migrationsPath)
	driver, err := postgres.WithInstance(db.Client, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("an error occurred while running migrate up: %w", err)
	}

	slog.Info("Migrations ran successfully")
	return nil
}

func (db *PostgresDB) Close() error {
	return db.Client.Close()
}
