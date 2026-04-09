-- =============================================================================
-- Initial schema
-- =============================================================================

-- Users -----------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_admin      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Counter (single shared row) -------------------------------------------
CREATE TABLE IF NOT EXISTS counter (
    id    INT PRIMARY KEY CHECK (id = 1),
    value INT NOT NULL DEFAULT 0
);

INSERT INTO counter (id, value) VALUES (1, 0) ON CONFLICT DO NOTHING;

-- Notes -----------------------------------------------------------------
CREATE TABLE IF NOT EXISTS notes (
    id         TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notes_user    ON notes(user_id);
CREATE INDEX IF NOT EXISTS idx_notes_created ON notes(created_at DESC);

-- Triggers --------------------------------------------------------------

-- Notify all listeners when the counter changes
CREATE OR REPLACE FUNCTION notify_counter_update()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify(
        'counter_updates',
        json_build_object('value', NEW.value)::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS counter_changed ON counter;
CREATE TRIGGER counter_changed
AFTER UPDATE ON counter
FOR EACH ROW EXECUTE FUNCTION notify_counter_update();

-- Notify all listeners when a note is inserted or deleted
CREATE OR REPLACE FUNCTION notify_notes_update()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        PERFORM pg_notify('notes_updates', json_build_object(
            'op', 'delete',
            'id', OLD.id
        )::text);
        RETURN OLD;
    ELSE
        PERFORM pg_notify('notes_updates', json_build_object(
            'op',         'insert',
            'id',         NEW.id,
            'user_id',    NEW.user_id,
            'username',   (SELECT username FROM users WHERE id = NEW.user_id),
            'content',    NEW.content,
            'created_at', NEW.created_at
        )::text);
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS notes_changed ON notes;
CREATE TRIGGER notes_changed
AFTER INSERT OR DELETE ON notes
FOR EACH ROW EXECUTE FUNCTION notify_notes_update();
