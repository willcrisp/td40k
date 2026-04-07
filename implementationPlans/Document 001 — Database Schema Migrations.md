Document 001 — Database: Schema & Migrations

Purpose


Create all database tables, indexes, triggers, and notification functions. Migrations are plain SQL files applied in order.


---

Migration 001 — db/migrations/001_create_players.sql

	CREATE TABLE IF NOT EXISTS players (
	    id         TEXT PRIMARY KEY,
	    nickname   TEXT NOT NULL,
	    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	    last_seen  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);


---

Migration 002 — db/migrations/002_create_rooms.sql

	CREATE TABLE IF NOT EXISTS rooms (
	    id              TEXT PRIMARY KEY,
	    name            TEXT NOT NULL,
	    status          TEXT NOT NULL DEFAULT 'lobby'
	                        CHECK (status IN ('lobby','active','finished','closed')),
	    game_master_id  TEXT NOT NULL REFERENCES players(id),
	    attacker_id     TEXT REFERENCES players(id),
	    defender_id     TEXT REFERENCES players(id),
	    battle_round    INT  NOT NULL DEFAULT 1
	                        CHECK (battle_round BETWEEN 1 AND 5),
	    active_player   TEXT NOT NULL DEFAULT 'attacker'
	                        CHECK (active_player IN ('attacker','defender')),
	    current_phase   TEXT NOT NULL DEFAULT 'command'
	                        CHECK (current_phase IN
	                            ('command','movement','shooting','charge','fight')),
	    winner          TEXT CHECK (winner IN ('attacker','defender')),
	    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_rooms_game_master
	    ON rooms(game_master_id);
	
	CREATE INDEX IF NOT EXISTS idx_rooms_attacker
	    ON rooms(attacker_id);
	
	CREATE INDEX IF NOT EXISTS idx_rooms_defender
	    ON rooms(defender_id);
	
	CREATE INDEX IF NOT EXISTS idx_rooms_status
	    ON rooms(status);


---

Migration 003 — db/migrations/003_create_room_events.sql

	CREATE TABLE IF NOT EXISTS room_events (
	    id           BIGSERIAL PRIMARY KEY,
	    room_id      TEXT        NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
	    player_id    TEXT        REFERENCES players(id),
	    event_type   TEXT        NOT NULL,
	    payload      JSONB,
	    occurred_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_room_events_room_id
	    ON room_events(room_id);
	
	CREATE INDEX IF NOT EXISTS idx_room_events_type
	    ON room_events(event_type);


---

Migration 004 — db/migrations/004_create_triggers.sql

	-- Auto-update updated_at on rooms
	CREATE OR REPLACE FUNCTION set_updated_at()
	RETURNS TRIGGER AS $$
	BEGIN
	    NEW.updated_at = NOW();
	    RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE TRIGGER rooms_set_updated_at
	BEFORE UPDATE ON rooms
	FOR EACH ROW EXECUTE FUNCTION set_updated_at();
	
	-- Broadcast room state changes via pg_notify
	CREATE OR REPLACE FUNCTION notify_room_update()
	RETURNS TRIGGER AS $$
	BEGIN
	    PERFORM pg_notify(
	        'room_updates',
	        json_build_object(
	            'room_id',       NEW.id,
	            'name',          NEW.name,
	            'status',        NEW.status,
	            'battle_round',  NEW.battle_round,
	            'active_player', NEW.active_player,
	            'current_phase', NEW.current_phase,
	            'winner',        NEW.winner,
	            'attacker_id',   NEW.attacker_id,
	            'defender_id',   NEW.defender_id,
	            'game_master_id',NEW.game_master_id
	        )::text
	    );
	    RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE TRIGGER room_state_changed
	AFTER INSERT OR UPDATE ON rooms
	FOR EACH ROW EXECUTE FUNCTION notify_room_update();


---

Migration Runner


The Go backend must run migrations on startup. Add this function to internal/db/db.go:


	func RunMigrations(pool *pgxpool.Pool) error {
	    migrations := []string{
	        "db/migrations/001_create_players.sql",
	        "db/migrations/002_create_rooms.sql",
	        "db/migrations/003_create_room_events.sql",
	        "db/migrations/004_create_triggers.sql",
	    }
	    for _, path := range migrations {
	        sql, err := os.ReadFile(path)
	        if err != nil {
	            return fmt.Errorf("read migration %s: %w", path, err)
	        }
	        if _, err := pool.Exec(context.Background(), string(sql)); err != nil {
	            return fmt.Errorf("run migration %s: %w", path, err)
	        }
	        log.Printf("[db] applied migration: %s", path)
	    }
	    return nil
	}