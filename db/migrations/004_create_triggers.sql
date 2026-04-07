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
