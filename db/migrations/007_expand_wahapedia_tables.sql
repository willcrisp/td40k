-- Expand Wahapedia schema to support full 16-file dataset

-- Add new columns to existing wh_datasheets table
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS source_id TEXT;
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS legend TEXT DEFAULT '';
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS role TEXT DEFAULT '';
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS loadout TEXT DEFAULT '';
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS transport TEXT DEFAULT '';
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS virtual BOOLEAN DEFAULT FALSE;
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS leader_head TEXT DEFAULT '';
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS leader_footer TEXT DEFAULT '';
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS damaged_w TEXT DEFAULT '';
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS damaged_description TEXT DEFAULT '';
ALTER TABLE wh_datasheets ADD COLUMN IF NOT EXISTS link TEXT DEFAULT '';

-- Add new columns to existing wh_datasheet_models table
ALTER TABLE wh_datasheet_models ADD COLUMN IF NOT EXISTS inv_sv TEXT DEFAULT '';
ALTER TABLE wh_datasheet_models ADD COLUMN IF NOT EXISTS inv_sv_descr TEXT DEFAULT '';
ALTER TABLE wh_datasheet_models ADD COLUMN IF NOT EXISTS base_size TEXT DEFAULT '';
ALTER TABLE wh_datasheet_models ADD COLUMN IF NOT EXISTS base_size_descr TEXT DEFAULT '';

-- Add new columns to existing wh_datasheet_abilities table
ALTER TABLE wh_datasheet_abilities ADD COLUMN IF NOT EXISTS line TEXT DEFAULT '';
ALTER TABLE wh_datasheet_abilities ADD COLUMN IF NOT EXISTS ability_id TEXT DEFAULT '';
ALTER TABLE wh_datasheet_abilities ADD COLUMN IF NOT EXISTS model TEXT DEFAULT '';
ALTER TABLE wh_datasheet_abilities ADD COLUMN IF NOT EXISTS type TEXT DEFAULT '';
ALTER TABLE wh_datasheet_abilities ADD COLUMN IF NOT EXISTS parameter TEXT DEFAULT '';

-- Update the primary key of wh_datasheet_abilities to include line
ALTER TABLE wh_datasheet_abilities DROP CONSTRAINT wh_datasheet_abilities_pkey;
ALTER TABLE wh_datasheet_abilities ADD PRIMARY KEY (datasheet_id, line, name);

-- wh_sources: Rulebooks, Supplements, Indexes
CREATE TABLE IF NOT EXISTS wh_sources (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT '',
    edition TEXT NOT NULL DEFAULT '',
    version TEXT NOT NULL DEFAULT '',
    errata_date TEXT DEFAULT '',
    errata_link TEXT DEFAULT ''
);

-- wh_stratagems: Stratagems per faction/detachment
CREATE TABLE IF NOT EXISTS wh_stratagems (
    id TEXT PRIMARY KEY,
    faction_id TEXT NOT NULL REFERENCES wh_factions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT '',
    cp_cost TEXT NOT NULL DEFAULT '',
    legend TEXT DEFAULT '',
    turn TEXT DEFAULT '',
    phase TEXT DEFAULT '',
    description TEXT DEFAULT '',
    detachment TEXT DEFAULT '',
    detachment_id TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_wh_stratagems_faction
    ON wh_stratagems(faction_id);

-- wh_abilities: Shared abilities across datasheets
CREATE TABLE IF NOT EXISTS wh_abilities (
    id TEXT PRIMARY KEY,
    faction_id TEXT NOT NULL REFERENCES wh_factions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    legend TEXT DEFAULT '',
    description TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_wh_abilities_faction
    ON wh_abilities(faction_id);

-- wh_enhancements: Enhancements per faction
CREATE TABLE IF NOT EXISTS wh_enhancements (
    id TEXT PRIMARY KEY,
    faction_id TEXT NOT NULL REFERENCES wh_factions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    legend TEXT DEFAULT '',
    description TEXT DEFAULT '',
    cost TEXT DEFAULT '',
    detachment TEXT DEFAULT '',
    detachment_id TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_wh_enhancements_faction
    ON wh_enhancements(faction_id);

-- wh_detachments: Detachments per faction
CREATE TABLE IF NOT EXISTS wh_detachments (
    id TEXT PRIMARY KEY,
    faction_id TEXT NOT NULL REFERENCES wh_factions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    legend TEXT DEFAULT '',
    type TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_wh_detachments_faction
    ON wh_detachments(faction_id);

-- wh_detachment_abilities: Detachment abilities
CREATE TABLE IF NOT EXISTS wh_detachment_abilities (
    id TEXT PRIMARY KEY,
    faction_id TEXT NOT NULL REFERENCES wh_factions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    legend TEXT DEFAULT '',
    description TEXT DEFAULT '',
    detachment TEXT DEFAULT '',
    detachment_id TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_wh_detachment_abilities_faction
    ON wh_detachment_abilities(faction_id);

-- wh_datasheet_keywords: Keywords per datasheet
CREATE TABLE IF NOT EXISTS wh_datasheet_keywords (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    line TEXT NOT NULL,
    keyword TEXT NOT NULL,
    model TEXT DEFAULT '',
    is_faction_keyword BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (datasheet_id, line, keyword)
);

-- wh_datasheet_options: Wargear options per datasheet
CREATE TABLE IF NOT EXISTS wh_datasheet_options (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    line TEXT NOT NULL,
    button TEXT DEFAULT '',
    description TEXT NOT NULL,
    PRIMARY KEY (datasheet_id, line)
);

-- wh_datasheet_wargear: Replaces/extends weapon data
CREATE TABLE IF NOT EXISTS wh_datasheet_wargear (
    id SERIAL PRIMARY KEY,
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    line TEXT NOT NULL,
    line_in_wargear TEXT DEFAULT '',
    dice TEXT DEFAULT '',
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    "range" TEXT DEFAULT '',
    type TEXT NOT NULL DEFAULT '',
    a TEXT DEFAULT '',
    bs_ws TEXT DEFAULT '',
    s TEXT DEFAULT '',
    ap TEXT DEFAULT '',
    d TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_wh_wargear_datasheet
    ON wh_datasheet_wargear(datasheet_id);

-- wh_datasheet_unit_composition: Unit composition entries
CREATE TABLE IF NOT EXISTS wh_datasheet_unit_composition (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    line TEXT NOT NULL,
    description TEXT NOT NULL,
    PRIMARY KEY (datasheet_id, line)
);

-- wh_datasheet_models_cost: Model cost entries
CREATE TABLE IF NOT EXISTS wh_datasheet_models_cost (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    line TEXT NOT NULL,
    description TEXT NOT NULL,
    cost TEXT NOT NULL,
    PRIMARY KEY (datasheet_id, line)
);

-- wh_datasheet_stratagems: Junction table for datasheets and stratagems
CREATE TABLE IF NOT EXISTS wh_datasheet_stratagems (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    stratagem_id TEXT NOT NULL REFERENCES wh_stratagems(id) ON DELETE CASCADE,
    PRIMARY KEY (datasheet_id, stratagem_id)
);

-- wh_datasheet_enhancements: Junction table for datasheets and enhancements
CREATE TABLE IF NOT EXISTS wh_datasheet_enhancements (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    enhancement_id TEXT NOT NULL REFERENCES wh_enhancements(id) ON DELETE CASCADE,
    PRIMARY KEY (datasheet_id, enhancement_id)
);

-- wh_datasheet_detachment_abilities: Junction table for datasheets and detachment abilities
CREATE TABLE IF NOT EXISTS wh_datasheet_detachment_abilities (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    detachment_ability_id TEXT NOT NULL REFERENCES wh_detachment_abilities(id) ON DELETE CASCADE,
    PRIMARY KEY (datasheet_id, detachment_ability_id)
);

-- wh_datasheet_leader: Junction table for datasheet leaders
CREATE TABLE IF NOT EXISTS wh_datasheet_leader (
    datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    attached_datasheet_id TEXT NOT NULL REFERENCES wh_datasheets(id) ON DELETE CASCADE,
    PRIMARY KEY (datasheet_id, attached_datasheet_id)
);
