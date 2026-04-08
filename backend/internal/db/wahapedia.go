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
	sources []models.WhSource,
	datasheets []models.WhDatasheet,
	dsModels []models.WhDatasheetModel,
	dsAbilities []models.WhDatasheetAbility,
	dsKeywords []models.WhDatasheetKeyword,
	dsOptions []models.WhDatasheetOption,
	dsWargear []models.WhDatasheetWargear,
	dsUnitComposition []models.WhDatasheetUnitComposition,
	dsModelsCost []models.WhDatasheetModelCost,
	dsStratagems []models.WhDatasheetStratagem,
	dsEnhancements []models.WhDatasheetEnhancement,
	dsDetachmentAbilities []models.WhDatasheetDetachmentAbility,
	dsLeaders []models.WhDatasheetLeader,
	stratagems []models.WhStratagem,
	abilities []models.WhAbility,
	enhancements []models.WhEnhancement,
	detachmentAbilities []models.WhDetachmentAbility,
	detachments []models.WhDetachment,
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
		"wh_datasheet_leader",
		"wh_datasheet_detachment_abilities",
		"wh_datasheet_enhancements",
		"wh_datasheet_stratagems",
		"wh_datasheet_models_cost",
		"wh_datasheet_unit_composition",
		"wh_datasheet_wargear",
		"wh_datasheet_options",
		"wh_datasheet_keywords",
		"wh_datasheet_abilities",
		"wh_datasheet_models",
		"wh_datasheets",
		"wh_detachment_abilities",
		"wh_detachments",
		"wh_enhancements",
		"wh_stratagems",
		"wh_abilities",
		"wh_sources",
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

	// Insert sources
	for _, s := range sources {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_sources (id, name, type, edition, version, errata_date, errata_link)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)
			 ON CONFLICT (id) DO NOTHING`,
			s.ID, s.Name, s.Type, s.Edition, s.Version, s.ErrataDate, s.ErrataLink,
		); err != nil {
			return fmt.Errorf("insert source %s: %w", s.ID, err)
		}
	}

	// Insert detachments (needed before datasheets that reference them)
	for _, d := range detachments {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_detachments (id, faction_id, name, legend, type)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (id) DO NOTHING`,
			d.ID, d.FactionID, d.Name, d.Legend, d.Type,
		); err != nil {
			return fmt.Errorf("insert detachment %s: %w", d.ID, err)
		}
	}

	// Insert stratagems
	for _, s := range stratagems {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_stratagems (id, faction_id, name, type, cp_cost, legend, turn, phase, description, detachment, detachment_id)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			 ON CONFLICT (id) DO NOTHING`,
			s.ID, s.FactionID, s.Name, s.Type, s.CPCost, s.Legend, s.Turn, s.Phase, s.Description, s.Detachment, s.DetachmentID,
		); err != nil {
			return fmt.Errorf("insert stratagem %s: %w", s.ID, err)
		}
	}

	// Insert abilities
	for _, a := range abilities {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_abilities (id, faction_id, name, legend, description)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (id) DO NOTHING`,
			a.ID, a.FactionID, a.Name, a.Legend, a.Description,
		); err != nil {
			return fmt.Errorf("insert ability %s: %w", a.ID, err)
		}
	}

	// Insert enhancements
	for _, e := range enhancements {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_enhancements (id, faction_id, name, legend, description, cost, detachment, detachment_id)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 ON CONFLICT (id) DO NOTHING`,
			e.ID, e.FactionID, e.Name, e.Legend, e.Description, e.Cost, e.Detachment, e.DetachmentID,
		); err != nil {
			return fmt.Errorf("insert enhancement %s: %w", e.ID, err)
		}
	}

	// Insert detachment abilities
	for _, da := range detachmentAbilities {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_detachment_abilities (id, faction_id, name, legend, description, detachment, detachment_id)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)
			 ON CONFLICT (id) DO NOTHING`,
			da.ID, da.FactionID, da.Name, da.Legend, da.Description, da.Detachment, da.DetachmentID,
		); err != nil {
			return fmt.Errorf("insert detachment_ability %s: %w", da.ID, err)
		}
	}

	// Insert datasheets
	for _, d := range datasheets {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheets (id, name, faction_id, source_id, legend, role, loadout, transport, virtual, leader_head, leader_footer, damaged_w, damaged_description, link)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			 ON CONFLICT (id) DO NOTHING`,
			d.ID, d.Name, d.FactionID, d.SourceID, d.Legend, d.Role, d.Loadout, d.Transport, d.Virtual, d.LeaderHead, d.LeaderFooter, d.DamagedW, d.DamagedDescription, d.Link,
		); err != nil {
			return fmt.Errorf("insert datasheet %s: %w", d.ID, err)
		}
	}

	// Insert datasheet models
	for _, m := range dsModels {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_models (datasheet_id, name, m, t, sv, inv_sv, inv_sv_descr, w, ld, oc, base_size, base_size_descr)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			 ON CONFLICT (datasheet_id, name) DO NOTHING`,
			m.DatasheetID, m.Name, m.M, m.T, m.SV, m.InvSV, m.InvSVDescr, m.W, m.LD, m.OC, m.BaseSize, m.BaseSizeDescr,
		); err != nil {
			return fmt.Errorf("insert model %s/%s: %w", m.DatasheetID, m.Name, err)
		}
	}

	// Insert datasheet abilities
	for _, a := range dsAbilities {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_abilities (datasheet_id, line, ability_id, model, name, description, type, parameter)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 ON CONFLICT (datasheet_id, line, name) DO NOTHING`,
			a.DatasheetID, a.Line, a.AbilityID, a.Model, a.Name, a.Description, a.Type, a.Parameter,
		); err != nil {
			return fmt.Errorf("insert datasheet_ability %s/%s: %w", a.DatasheetID, a.Name, err)
		}
	}

	// Insert datasheet keywords
	for _, k := range dsKeywords {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_keywords (datasheet_id, line, keyword, model, is_faction_keyword)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (datasheet_id, line, keyword) DO NOTHING`,
			k.DatasheetID, k.Line, k.Keyword, k.Model, k.IsFactionKeyword,
		); err != nil {
			return fmt.Errorf("insert datasheet_keyword %s/%s: %w", k.DatasheetID, k.Keyword, err)
		}
	}

	// Insert datasheet options
	for _, o := range dsOptions {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_options (datasheet_id, line, button, description)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (datasheet_id, line) DO NOTHING`,
			o.DatasheetID, o.Line, o.Button, o.Description,
		); err != nil {
			return fmt.Errorf("insert datasheet_option %s/%s: %w", o.DatasheetID, o.Line, err)
		}
	}

	// Insert datasheet wargear
	for _, w := range dsWargear {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_wargear (datasheet_id, line, line_in_wargear, dice, name, description, "range", type, a, bs_ws, s, ap, d)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
			w.DatasheetID, w.Line, w.LineInWargear, w.Dice, w.Name, w.Description, w.Range, w.Type, w.A, w.BSWS, w.S, w.AP, w.D,
		); err != nil {
			return fmt.Errorf("insert datasheet_wargear %s/%s: %w", w.DatasheetID, w.Name, err)
		}
	}

	// Insert datasheet unit composition
	for _, uc := range dsUnitComposition {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_unit_composition (datasheet_id, line, description)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (datasheet_id, line) DO NOTHING`,
			uc.DatasheetID, uc.Line, uc.Description,
		); err != nil {
			return fmt.Errorf("insert datasheet_unit_composition %s/%s: %w", uc.DatasheetID, uc.Line, err)
		}
	}

	// Insert datasheet models cost
	for _, mc := range dsModelsCost {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_models_cost (datasheet_id, line, description, cost)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (datasheet_id, line) DO NOTHING`,
			mc.DatasheetID, mc.Line, mc.Description, mc.Cost,
		); err != nil {
			return fmt.Errorf("insert datasheet_models_cost %s/%s: %w", mc.DatasheetID, mc.Line, err)
		}
	}

	// Insert datasheet stratagems junction
	for _, ds := range dsStratagems {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_stratagems (datasheet_id, stratagem_id)
			 VALUES ($1, $2)
			 ON CONFLICT (datasheet_id, stratagem_id) DO NOTHING`,
			ds.DatasheetID, ds.StratagemID,
		); err != nil {
			return fmt.Errorf("insert datasheet_stratagem %s/%s: %w", ds.DatasheetID, ds.StratagemID, err)
		}
	}

	// Insert datasheet enhancements junction
	for _, de := range dsEnhancements {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_enhancements (datasheet_id, enhancement_id)
			 VALUES ($1, $2)
			 ON CONFLICT (datasheet_id, enhancement_id) DO NOTHING`,
			de.DatasheetID, de.EnhancementID,
		); err != nil {
			return fmt.Errorf("insert datasheet_enhancement %s/%s: %w", de.DatasheetID, de.EnhancementID, err)
		}
	}

	// Insert datasheet detachment abilities junction
	for _, da := range dsDetachmentAbilities {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_detachment_abilities (datasheet_id, detachment_ability_id)
			 VALUES ($1, $2)
			 ON CONFLICT (datasheet_id, detachment_ability_id) DO NOTHING`,
			da.DatasheetID, da.DetachmentAbilityID,
		); err != nil {
			return fmt.Errorf("insert datasheet_detachment_ability %s/%s: %w", da.DatasheetID, da.DetachmentAbilityID, err)
		}
	}

	// Insert datasheet leaders junction
	for _, l := range dsLeaders {
		if _, err := tx.Exec(ctx,
			`INSERT INTO wh_datasheet_leader (datasheet_id, attached_datasheet_id)
			 VALUES ($1, $2)
			 ON CONFLICT (datasheet_id, attached_datasheet_id) DO NOTHING`,
			l.DatasheetID, l.AttachedDatasheetID,
		); err != nil {
			return fmt.Errorf("insert datasheet_leader %s/%s: %w", l.DatasheetID, l.AttachedDatasheetID, err)
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

// GetAllDatasheets returns all available Warhammer datasheets
func GetAllDatasheets() ([]models.WhDatasheet, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT id, name, faction_id, source_id, legend, role, loadout,
		       transport, virtual, leader_head, leader_footer, damaged_w,
		       damaged_description, link
		FROM wh_datasheets
		ORDER BY faction_id, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datasheets []models.WhDatasheet
	for rows.Next() {
		var ds models.WhDatasheet
		if err := rows.Scan(
			&ds.ID, &ds.Name, &ds.FactionID, &ds.SourceID, &ds.Legend,
			&ds.Role, &ds.Loadout, &ds.Transport, &ds.Virtual,
			&ds.LeaderHead, &ds.LeaderFooter, &ds.DamagedW,
			&ds.DamagedDescription, &ds.Link,
		); err != nil {
			return nil, err
		}
		datasheets = append(datasheets, ds)
	}
	return datasheets, rows.Err()
}

// GetDatasheetModels returns all models for a specific datasheet
func GetDatasheetModels(
	datasheetID string,
) ([]models.WhDatasheetModel, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT datasheet_id, name, m, t, sv, inv_sv, inv_sv_descr,
		       w, ld, oc, base_size, base_size_descr
		FROM wh_datasheet_models
		WHERE datasheet_id = $1
		ORDER BY name
	`, datasheetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dsModels []models.WhDatasheetModel
	for rows.Next() {
		var m models.WhDatasheetModel
		if err := rows.Scan(
			&m.DatasheetID, &m.Name, &m.M, &m.T, &m.SV, &m.InvSV,
			&m.InvSVDescr, &m.W, &m.LD, &m.OC, &m.BaseSize,
			&m.BaseSizeDescr,
		); err != nil {
			return nil, err
		}
		dsModels = append(dsModels, m)
	}
	return dsModels, rows.Err()
}
