CREATE TABLE IF NOT EXISTS notes (
    id         TEXT PRIMARY KEY,
    player_id  TEXT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notes_player ON notes(player_id);
CREATE INDEX IF NOT EXISTS idx_notes_created ON notes(created_at DESC);

CREATE OR REPLACE FUNCTION notify_notes_update()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        PERFORM pg_notify('notes_updates', json_build_object(
            'op',  'delete',
            'id',  OLD.id
        )::text);
        RETURN OLD;
    ELSE
        PERFORM pg_notify('notes_updates', json_build_object(
            'op',         'insert',
            'id',         NEW.id,
            'player_id',  NEW.player_id,
            'username',   (SELECT username FROM players WHERE id = NEW.player_id),
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
