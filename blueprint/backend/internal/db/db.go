package db

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	"db/migrations/001_create_players.sql",
	"db/migrations/002_create_counter.sql",
	"db/migrations/003_create_triggers.sql",
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

		statements := strings.Split(string(sql), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := Pool.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("exec migration %s: %w", path, err)
			}
		}

		if _, err := Pool.Exec(ctx,
			`INSERT INTO schema_migrations (version) VALUES ($1)`, version,
		); err != nil {
			return fmt.Errorf("record migration %s: %w", version, err)
		}

		fmt.Printf("applied migration %s\n", version)
	}
	return nil
}
