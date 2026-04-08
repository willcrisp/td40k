-- Auto-update updated_at on game_units
DROP TRIGGER IF EXISTS game_units_set_updated_at ON game_units;
CREATE TRIGGER game_units_set_updated_at
BEFORE UPDATE ON game_units
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Broadcast game units changes via pg_notify
CREATE OR REPLACE FUNCTION notify_game_units_update()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify(
        'game_units_updates',
        json_build_object(
            'room_id',       COALESCE(NEW.room_id, OLD.room_id),
            'unit_id',       COALESCE(NEW.id, OLD.id),
            'event_type',    CASE
                                WHEN TG_OP = 'DELETE' THEN 'unit_removed'
                                WHEN TG_OP = 'INSERT' THEN 'unit_placed'
                                ELSE 'unit_moved'
                             END,
            'x',             NEW.x,
            'y',             NEW.y,
            'facing_degrees',NEW.facing_degrees,
            'status',        NEW.status,
            'wounds',        NEW.wounds,
            'model_count',   NEW.model_count,
            'owner_player_id',NEW.owner_player_id
        )::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS game_units_changed ON game_units;
CREATE TRIGGER game_units_changed
AFTER INSERT OR UPDATE OR DELETE ON game_units
FOR EACH ROW EXECUTE FUNCTION notify_game_units_update();
