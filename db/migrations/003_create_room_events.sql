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
