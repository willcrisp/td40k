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
