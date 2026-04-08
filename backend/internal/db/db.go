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
	migrations := []string{
		"db/migrations/001_create_players.sql",
		"db/migrations/002_create_rooms.sql",
		"db/migrations/003_create_room_events.sql",
		"db/migrations/004_create_triggers.sql",
		"db/migrations/005_create_wahapedia_tables.sql",
		"db/migrations/006_add_auth_to_players.sql",
		"db/migrations/007_expand_wahapedia_tables.sql",
		"db/migrations/008_add_admin_flag_to_players.sql",
	}
	for _, path := range migrations {
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", path, err)
		}
		if _, err := Pool.Exec(context.Background(), string(sql)); err != nil {
			return fmt.Errorf("run migration %s: %w", path, err)
		}
		log.Printf("[db] applied: %s", path)
	}
	return nil
}
