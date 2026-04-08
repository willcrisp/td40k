-- Create game_units table for deployed units on the board
CREATE TABLE game_units (
    id              TEXT PRIMARY KEY,
    room_id         TEXT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    faction_id      TEXT NOT NULL,
    datasheet_id    TEXT NOT NULL REFERENCES wh_datasheets(id),
    model_name      TEXT NOT NULL,
    name_on_board   TEXT,
    model_count     INT NOT NULL DEFAULT 1,
    x               FLOAT NOT NULL,
    y               FLOAT NOT NULL,
    facing_degrees  INT NOT NULL DEFAULT 0,
    status          TEXT NOT NULL DEFAULT 'alive'
                        CHECK (status IN ('alive', 'in_reserves', 'dead')),
    wounds          INT NOT NULL DEFAULT 0,
    owner_player_id TEXT NOT NULL REFERENCES players(id),
    created_by      TEXT NOT NULL REFERENCES players(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_game_units_room_id ON game_units(room_id);
CREATE INDEX idx_game_units_owner ON game_units(owner_player_id);
CREATE INDEX idx_game_units_created_by ON game_units(created_by);

-- Create unit_roster table for tracking available units per game
CREATE TABLE unit_roster (
    id              TEXT PRIMARY KEY,
    room_id         TEXT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    player_id       TEXT NOT NULL REFERENCES players(id),
    datasheet_id    TEXT NOT NULL REFERENCES wh_datasheets(id),
    model_name      TEXT NOT NULL,
    quantity        INT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_unit_roster_room_id ON unit_roster(room_id);
CREATE INDEX idx_unit_roster_player ON unit_roster(player_id);
