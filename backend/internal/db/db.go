package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(dsn string) error {
	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("db init: %w", err)
	}
	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("db ping: %w", err)
	}
	log.Println("[db] connected to postgres")
	return nil
}

func RunMigrations() error {
	// Create migrations tracking table if it doesn't exist
	if _, err := Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	migrations := []string{
		"db/migrations/001_create_players.sql",
		"db/migrations/002_create_rooms.sql",
		"db/migrations/003_create_room_events.sql",
		"db/migrations/004_create_triggers.sql",
		"db/migrations/005_create_wahapedia_tables.sql",
		"db/migrations/006_add_auth_to_players.sql",
		"db/migrations/007_expand_wahapedia_tables.sql",
		"db/migrations/008_add_admin_flag_to_players.sql",
		"db/migrations/009_create_game_units.sql",
		"db/migrations/010_add_game_units_trigger.sql",
		"db/migrations/011_roster_import_support.sql",
	}
	for _, path := range migrations {
		// Extract version from filename (e.g., "001" from "001_...")
		filename := path[len("db/migrations/"):]
		version := filename[:3] // First 3 chars are the version

		// Check if migration already applied
		var count int
		err := Pool.QueryRow(
			context.Background(),
			"SELECT COUNT(*) FROM schema_migrations WHERE version = $1",
			version,
		).Scan(&count)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", path, err)
		}

		if count > 0 {
			log.Printf("[db] skipped (already applied): %s", path)
			continue
		}

		// Read and execute migration
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", path, err)
		}
		if _, err := Pool.Exec(context.Background(), string(sql)); err != nil {
			return fmt.Errorf("run migration %s: %w", path, err)
		}

		// Record migration as applied
		if _, err := Pool.Exec(
			context.Background(),
			"INSERT INTO schema_migrations (version) VALUES ($1)",
			version,
		); err != nil {
			return fmt.Errorf("record migration %s: %w", path, err)
		}

		log.Printf("[db] applied: %s", path)
	}
	return nil
}
