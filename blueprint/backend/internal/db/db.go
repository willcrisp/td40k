package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(dsn string) error {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("db ping: %w", err)
	}
	Pool = pool
	return nil
}

var migrations = []string{
	"db/migrations/001_schema.sql",
}

func RunMigrations() error {
	ctx := context.Background()

	_, err := Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	for _, path := range migrations {
		version := path[len("db/migrations/") : len("db/migrations/")+3]

		var exists bool
		err := Pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`,
			version,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", version, err)
		}
		if exists {
			continue
		}

		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", path, err)
		}

		// Execute the whole file as one statement — required for PL/pgSQL
		// functions that contain semicolons inside $$ blocks.
		if _, err := Pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("exec migration %s: %w", path, err)
		}

		if _, err := Pool.Exec(ctx,
			`INSERT INTO schema_migrations (version) VALUES ($1)`, version,
		); err != nil {
			return fmt.Errorf("record migration %s: %w", version, err)
		}

		slog.Info("applied migration", "version", version)
	}
	return nil
}
