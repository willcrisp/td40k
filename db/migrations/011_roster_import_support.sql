ALTER TABLE unit_roster
  ADD COLUMN IF NOT EXISTS faction_id TEXT
    REFERENCES wh_factions(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS points INT NOT NULL DEFAULT 0;

ALTER TABLE unit_roster
  ADD CONSTRAINT uq_roster_entry
  UNIQUE (room_id, player_id, datasheet_id, model_name);
