package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/willcrisp/td40k/internal/models"
)

func GetWahapediaHash(sourceName string) (string, error) {
	var hash string
	err := Pool.QueryRow(context.Background(),
		`SELECT content_hash FROM wahapedia_sync WHERE source_name = $1`,
		sourceName,
	).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return hash, err
}

func SyncWahapediaData(
	factions []models.WhFaction,
	datasheets []models.WhDatasheet,
	dsModels []models.WhDatasheetModel,
	weapons []models.WhDatasheetWeapon,
	abilities []models.WhDatasheetAbility,
	hashes map[string]string,
) error {
	ctx := context.Background()
	tx, err := Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete in FK order: children first
	for _, table := range []string{
		"wh_datasheet_abilities",
		"wh_datasheet_weapons",
		"wh_datasheet_models",
		"wh_datasheets",
		"wh_factions",
	} {
		if _, err := tx.Exec(ctx, "DELETE FROM "+table); err != nil {
			return fmt.Errorf("delete %s: %w", table, err)
		}
	}

	// Insert factions
	for _, f := range factions {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_factions (id, name, link) VALUES ($1, $2, $3)`,
			f.ID, f.Name, f.Link,
		); err != nil {
			return fmt.Errorf("insert faction %s: %w", f.ID, err)
		}
	}

	// Insert datasheets
	for _, d := range datasheets {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheets (id, name, faction_id) VALUES ($1, $2, $3)`,
			d.ID, d.Name, d.FactionID,
		); err != nil {
			return fmt.Errorf("insert datasheet %s: %w", d.ID, err)
		}
	}

	// Insert datasheet models
	for _, m := range dsModels {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_models (datasheet_id, name, m, t, sv, w, ld, oc)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 ON CONFLICT (datasheet_id, name) DO NOTHING`,
			m.DatasheetID, m.Name, m.M, m.T, m.SV, m.W, m.LD, m.OC,
		); err != nil {
			return fmt.Errorf("insert model %s/%s: %w", m.DatasheetID, m.Name, err)
		}
	}

	// Insert weapons
	for _, w := range weapons {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_weapons
			 (datasheet_id, name, type, "range", a, bs, s, ap, d, abilities)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			w.DatasheetID, w.Name, w.Type, w.Range,
			w.A, w.BS, w.S, w.AP, w.D, w.Abilities,
		); err != nil {
			return fmt.Errorf("insert weapon %s/%s: %w", w.DatasheetID, w.Name, err)
		}
	}

	// Insert abilities
	for _, a := range abilities {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_abilities (datasheet_id, name, description)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (datasheet_id, name) DO NOTHING`,
			a.DatasheetID, a.Name, a.Description,
		); err != nil {
			return fmt.Errorf("insert ability %s/%s: %w", a.DatasheetID, a.Name, err)
		}
	}

	// Upsert hashes
	for source, hash := range hashes {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wahapedia_sync (source_name, content_hash, synced_at)
			 VALUES ($1, $2, NOW())
			 ON CONFLICT (source_name) DO UPDATE
			   SET content_hash = EXCLUDED.content_hash,
			       synced_at    = NOW()`,
			source, hash,
		); err != nil {
			return fmt.Errorf("upsert hash %s: %w", source, err)
		}
	}

	return tx.Commit(ctx)
}
