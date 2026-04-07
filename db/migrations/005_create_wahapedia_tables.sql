-- Tracks content hash per CSV source for change detection
CREATE TABLE IF NOT EXISTS wahapedia_sync (
    source_name TEXT PRIMARY KEY,
    content_hash TEXT NOT NULL,
    synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wh_factions (
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    link TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS wh_datasheets (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    faction_id TEXT NOT NULL REFERENCES wh_factions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_wh_datasheets_faction
    ON wh_datasheets(faction_id);

CREATE TABLE IF NOT EXISTS wh_datasheet_models (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    m            TEXT NOT NULL DEFAULT '',
    t            TEXT NOT NULL DEFAULT '',
    sv           TEXT NOT NULL DEFAULT '',
    w            TEXT NOT NULL DEFAULT '',
    ld           TEXT NOT NULL DEFAULT '',
    oc           TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (datasheet_id, name)
);

CREATE TABLE IF NOT EXISTS wh_datasheet_weapons (
    id           SERIAL PRIMARY KEY,
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    type         TEXT NOT NULL DEFAULT '',
    "range"      TEXT NOT NULL DEFAULT '',
    a            TEXT NOT NULL DEFAULT '',
    bs           TEXT NOT NULL DEFAULT '',
    s            TEXT NOT NULL DEFAULT '',
    ap           TEXT NOT NULL DEFAULT '',
    d            TEXT NOT NULL DEFAULT '',
    abilities    TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_wh_weapons_datasheet
    ON wh_datasheet_weapons(datasheet_id);

CREATE TABLE IF NOT EXISTS wh_datasheet_abilities (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (datasheet_id, name)
);
