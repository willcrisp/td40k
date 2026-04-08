package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/willcrisp/td40k/internal/db"
	"github.com/willcrisp/td40k/internal/middleware"
	"github.com/willcrisp/td40k/internal/models"
)

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func stripHTML(s string) string {
	s = htmlTagRe.ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	return strings.TrimSpace(s)
}

// csvSources maps logical names to filenames on Wahapedia.
var csvSources = map[string]string{
	"Factions":                       "Factions.csv",
	"Source":                         "Source.csv",
	"Datasheets":                     "Datasheets.csv",
	"Datasheets_abilities":           "Datasheets_abilities.csv",
	"Datasheets_keywords":            "Datasheets_keywords.csv",
	"Datasheets_models":              "Datasheets_models.csv",
	"Datasheets_options":             "Datasheets_options.csv",
	"Datasheets_wargear":             "Datasheets_wargear.csv",
	"Datasheets_unit_composition":    "Datasheets_unit_composition.csv",
	"Datasheets_models_cost":         "Datasheets_models_cost.csv",
	"Datasheets_stratagems":          "Datasheets_stratagems.csv",
	"Datasheets_enhancements":        "Datasheets_enhancements.csv",
	"Datasheets_detachment_abilities": "Datasheets_detachment_abilities.csv",
	"Datasheets_leader":              "Datasheets_leader.csv",
	"Stratagems":                     "Stratagems.csv",
	"Abilities":                      "Abilities.csv",
	"Enhancements":                   "Enhancements.csv",
	"Detachment_abilities":           "Detachment_abilities.csv",
	"Detachments":                    "Detachments.csv",
}

const wahapediaBaseURL = "http://wahapedia.ru/wh40k10ed/"

func fetchCSV(ctx context.Context, filename string) ([]byte, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wahapediaBaseURL+filename, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("fetch %s: status %d", filename, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	sum := sha256.Sum256(body)
	hash := hex.EncodeToString(sum[:])
	return body, hash, nil
}

// colIndex builds a header-name → column-index map from the first CSV row.
func colIndex(headers []string) map[string]int {
	m := make(map[string]int, len(headers))
	for i, h := range headers {
		m[strings.TrimSpace(h)] = i
	}
	return m
}

// field safely returns a stripped field value by column name, or "" if missing.
func field(row []string, cols map[string]int, name string) string {
	idx, ok := cols[name]
	if !ok || idx >= len(row) {
		return ""
	}
	return stripHTML(row[idx])
}

func newPipeReader(data []byte) *csv.Reader {
	// Strip UTF-8 BOM if present
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = '|'
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	return r
}

func parseFactions(data []byte) ([]models.WhFaction, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("factions CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhFaction
	for _, row := range rows[1:] {
		id := field(row, cols, "id")
		if id == "" {
			continue
		}
		out = append(out, models.WhFaction{
			ID:   id,
			Name: field(row, cols, "name"),
			Link: field(row, cols, "link"),
		})
	}
	return out, nil
}

func parseDatasheets(data []byte) ([]models.WhDatasheet, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheets CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheet
	for _, row := range rows[1:] {
		id := field(row, cols, "id")
		if id == "" {
			continue
		}
		virtual := field(row, cols, "virtual") == "true"
		out = append(out, models.WhDatasheet{
			ID:                 id,
			Name:               field(row, cols, "name"),
			FactionID:          field(row, cols, "faction_id"),
			SourceID:           field(row, cols, "source_id"),
			Legend:             field(row, cols, "legend"),
			Role:               field(row, cols, "role"),
			Loadout:            field(row, cols, "loadout"),
			Transport:          field(row, cols, "transport"),
			Virtual:            virtual,
			LeaderHead:         field(row, cols, "leader_head"),
			LeaderFooter:       field(row, cols, "leader_footer"),
			DamagedW:           field(row, cols, "damaged_w"),
			DamagedDescription: field(row, cols, "damaged_description"),
			Link:               field(row, cols, "link"),
		})
	}
	return out, nil
}

func parseDatasheetModels(data []byte) ([]models.WhDatasheetModel, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_models CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetModel
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		if dsID == "" {
			continue
		}
		out = append(out, models.WhDatasheetModel{
			DatasheetID:  dsID,
			Name:         field(row, cols, "name"),
			M:            field(row, cols, "M"),
			T:            field(row, cols, "T"),
			SV:           field(row, cols, "Sv"),
			InvSV:        field(row, cols, "inv_sv"),
			InvSVDescr:   field(row, cols, "inv_sv_descr"),
			W:            field(row, cols, "W"),
			LD:           field(row, cols, "Ld"),
			OC:           field(row, cols, "OC"),
			BaseSize:     field(row, cols, "base_size"),
			BaseSizeDescr: field(row, cols, "base_size_descr"),
		})
	}
	return out, nil
}

func parseDatasheetWeapons(data []byte) ([]models.WhDatasheetWeapon, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_weapons CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetWeapon
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		if dsID == "" {
			continue
		}
		out = append(out, models.WhDatasheetWeapon{
			DatasheetID: dsID,
			Name:        field(row, cols, "name"),
			Type:        field(row, cols, "type"),
			Range:       field(row, cols, "Range"),
			A:           field(row, cols, "A"),
			BS:          field(row, cols, "BS"),
			S:           field(row, cols, "S"),
			AP:          field(row, cols, "AP"),
			D:           field(row, cols, "D"),
			Abilities:   field(row, cols, "abilities"),
		})
	}
	return out, nil
}

func parseDatasheetAbilities(data []byte) ([]models.WhDatasheetAbility, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_abilities CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetAbility
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		if dsID == "" {
			continue
		}
		out = append(out, models.WhDatasheetAbility{
			DatasheetID: dsID,
			Line:        field(row, cols, "line"),
			AbilityID:   field(row, cols, "ability_id"),
			Model:       field(row, cols, "model"),
			Name:        field(row, cols, "name"),
			Description: field(row, cols, "description"),
			Type:        field(row, cols, "type"),
			Parameter:   field(row, cols, "parameter"),
		})
	}
	return out, nil
}

func parseSources(data []byte) ([]models.WhSource, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("source CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhSource
	for _, row := range rows[1:] {
		id := field(row, cols, "id")
		if id == "" {
			continue
		}
		out = append(out, models.WhSource{
			ID:         id,
			Name:       field(row, cols, "name"),
			Type:       field(row, cols, "type"),
			Edition:    field(row, cols, "edition"),
			Version:    field(row, cols, "version"),
			ErrataDate: field(row, cols, "errata_date"),
			ErrataLink: field(row, cols, "errata_link"),
		})
	}
	return out, nil
}

func parseStratagems(data []byte) ([]models.WhStratagem, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("stratagems CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhStratagem
	for _, row := range rows[1:] {
		id := field(row, cols, "id")
		if id == "" {
			continue
		}
		out = append(out, models.WhStratagem{
			ID:           id,
			FactionID:    field(row, cols, "faction_id"),
			Name:         field(row, cols, "name"),
			Type:         field(row, cols, "type"),
			CPCost:       field(row, cols, "cp_cost"),
			Legend:       field(row, cols, "legend"),
			Turn:         field(row, cols, "turn"),
			Phase:        field(row, cols, "phase"),
			Description:  field(row, cols, "description"),
			Detachment:   field(row, cols, "detachment"),
			DetachmentID: field(row, cols, "detachment_id"),
		})
	}
	return out, nil
}

func parseAbilities(data []byte) ([]models.WhAbility, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("abilities CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhAbility
	for _, row := range rows[1:] {
		id := field(row, cols, "id")
		if id == "" {
			continue
		}
		out = append(out, models.WhAbility{
			ID:          id,
			FactionID:   field(row, cols, "faction_id"),
			Name:        field(row, cols, "name"),
			Legend:      field(row, cols, "legend"),
			Description: field(row, cols, "description"),
		})
	}
	return out, nil
}

func parseEnhancements(data []byte) ([]models.WhEnhancement, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("enhancements CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhEnhancement
	for _, row := range rows[1:] {
		id := field(row, cols, "id")
		if id == "" {
			continue
		}
		out = append(out, models.WhEnhancement{
			ID:           id,
			FactionID:    field(row, cols, "faction_id"),
			Name:         field(row, cols, "name"),
			Legend:       field(row, cols, "legend"),
			Description:  field(row, cols, "description"),
			Cost:         field(row, cols, "cost"),
			Detachment:   field(row, cols, "detachment"),
			DetachmentID: field(row, cols, "detachment_id"),
		})
	}
	return out, nil
}

func parseDetachments(data []byte) ([]models.WhDetachment, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("detachments CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDetachment
	for _, row := range rows[1:] {
		id := field(row, cols, "id")
		if id == "" {
			continue
		}
		out = append(out, models.WhDetachment{
			ID:        id,
			FactionID: field(row, cols, "faction_id"),
			Name:      field(row, cols, "name"),
			Legend:    field(row, cols, "legend"),
			Type:      field(row, cols, "type"),
		})
	}
	return out, nil
}

func parseDetachmentAbilities(data []byte) ([]models.WhDetachmentAbility, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("detachment_abilities CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDetachmentAbility
	for _, row := range rows[1:] {
		id := field(row, cols, "id")
		if id == "" {
			continue
		}
		out = append(out, models.WhDetachmentAbility{
			ID:           id,
			FactionID:    field(row, cols, "faction_id"),
			Name:         field(row, cols, "name"),
			Legend:       field(row, cols, "legend"),
			Description:  field(row, cols, "description"),
			Detachment:   field(row, cols, "detachment"),
			DetachmentID: field(row, cols, "detachment_id"),
		})
	}
	return out, nil
}

func parseDatasheetKeywords(data []byte) ([]models.WhDatasheetKeyword, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_keywords CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetKeyword
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		if dsID == "" {
			continue
		}
		isFaction := field(row, cols, "is_faction_keyword") == "true"
		out = append(out, models.WhDatasheetKeyword{
			DatasheetID:      dsID,
			Line:             field(row, cols, "line"),
			Keyword:          field(row, cols, "keyword"),
			Model:            field(row, cols, "model"),
			IsFactionKeyword: isFaction,
		})
	}
	return out, nil
}

func parseDatasheetOptions(data []byte) ([]models.WhDatasheetOption, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_options CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetOption
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		if dsID == "" {
			continue
		}
		out = append(out, models.WhDatasheetOption{
			DatasheetID: dsID,
			Line:        field(row, cols, "line"),
			Button:      field(row, cols, "button"),
			Description: field(row, cols, "description"),
		})
	}
	return out, nil
}

func parseDatasheetWargear(data []byte) ([]models.WhDatasheetWargear, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_wargear CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetWargear
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		if dsID == "" {
			continue
		}
		out = append(out, models.WhDatasheetWargear{
			DatasheetID:   dsID,
			Line:          field(row, cols, "line"),
			LineInWargear: field(row, cols, "line_in_wargear"),
			Dice:          field(row, cols, "dice"),
			Name:          field(row, cols, "name"),
			Description:   field(row, cols, "description"),
			Range:         field(row, cols, "range"),
			Type:          field(row, cols, "type"),
			A:             field(row, cols, "A"),
			BSWS:          field(row, cols, "BS_WS"),
			S:             field(row, cols, "S"),
			AP:            field(row, cols, "AP"),
			D:             field(row, cols, "D"),
		})
	}
	return out, nil
}

func parseDatasheetUnitComposition(data []byte) ([]models.WhDatasheetUnitComposition, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_unit_composition CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetUnitComposition
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		if dsID == "" {
			continue
		}
		out = append(out, models.WhDatasheetUnitComposition{
			DatasheetID: dsID,
			Line:        field(row, cols, "line"),
			Description: field(row, cols, "description"),
		})
	}
	return out, nil
}

func parseDatasheetModelsCost(data []byte) ([]models.WhDatasheetModelCost, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_models_cost CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetModelCost
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		if dsID == "" {
			continue
		}
		out = append(out, models.WhDatasheetModelCost{
			DatasheetID: dsID,
			Line:        field(row, cols, "line"),
			Description: field(row, cols, "description"),
			Cost:        field(row, cols, "cost"),
		})
	}
	return out, nil
}

func parseDatasheetStratagems(data []byte) ([]models.WhDatasheetStratagem, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_stratagems CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetStratagem
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		stratID := field(row, cols, "stratagem_id")
		if dsID == "" || stratID == "" {
			continue
		}
		out = append(out, models.WhDatasheetStratagem{
			DatasheetID: dsID,
			StratagemID: stratID,
		})
	}
	return out, nil
}

func parseDatasheetEnhancements(data []byte) ([]models.WhDatasheetEnhancement, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_enhancements CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetEnhancement
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		enhID := field(row, cols, "enhancement_id")
		if dsID == "" || enhID == "" {
			continue
		}
		out = append(out, models.WhDatasheetEnhancement{
			DatasheetID:   dsID,
			EnhancementID: enhID,
		})
	}
	return out, nil
}

func parseDatasheetDetachmentAbilities(data []byte) ([]models.WhDatasheetDetachmentAbility, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_detachment_abilities CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetDetachmentAbility
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		detchAbilID := field(row, cols, "detachment_ability_id")
		if dsID == "" || detchAbilID == "" {
			continue
		}
		out = append(out, models.WhDatasheetDetachmentAbility{
			DatasheetID:         dsID,
			DetachmentAbilityID: detchAbilID,
		})
	}
	return out, nil
}

func parseDatasheetLeaders(data []byte) ([]models.WhDatasheetLeader, error) {
	r := newPipeReader(data)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, fmt.Errorf("datasheet_leader CSV: no header row")
	}
	cols := colIndex(rows[0])
	var out []models.WhDatasheetLeader
	for _, row := range rows[1:] {
		dsID := field(row, cols, "datasheet_id")
		attachID := field(row, cols, "attached_datasheet_id")
		if dsID == "" || attachID == "" {
			continue
		}
		out = append(out, models.WhDatasheetLeader{
			DatasheetID:        dsID,
			AttachedDatasheetID: attachID,
		})
	}
	return out, nil
}

func HandleSyncWahapedia(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check if player is admin
	playerID := middleware.GetPlayerID(r)
	player, err := db.GetPlayerByID(playerID)
	if err != nil || !player.IsAdmin {
		jsonError(w, "admin access required", http.StatusForbidden)
		return
	}

	// Fetch all CSVs and compute hashes
	type fetchResult struct {
		body []byte
		hash string
	}
	fetched := make(map[string]fetchResult, len(csvSources))
	hashes := make(map[string]string, len(csvSources))

	for name, filename := range csvSources {
		body, hash, err := fetchCSV(ctx, filename)
		if err != nil {
			jsonError(w, fmt.Sprintf("fetch %s failed: %v", name, err), http.StatusBadGateway)
			return
		}
		fetched[name] = fetchResult{body: body, hash: hash}
		hashes[name] = hash
	}

	// Compare hashes to detect changes
	var changed []string
	for name, hash := range hashes {
		stored, err := db.GetWahapediaHash(name)
		if err != nil {
			jsonError(w, "db error checking hash", http.StatusInternalServerError)
			return
		}
		if stored != hash {
			changed = append(changed, name)
		}
	}

	if len(changed) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.WhSyncResult{Changed: false})
		return
	}

	// Parse all CSVs
	factions, err := parseFactions(fetched["Factions"].body)
	if err != nil {
		jsonError(w, "parse factions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sources, err := parseSources(fetched["Source"].body)
	if err != nil {
		jsonError(w, "parse sources: "+err.Error(), http.StatusInternalServerError)
		return
	}

	datasheets, err := parseDatasheets(fetched["Datasheets"].body)
	if err != nil {
		jsonError(w, "parse datasheets: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stratagems, err := parseStratagems(fetched["Stratagems"].body)
	if err != nil {
		jsonError(w, "parse stratagems: "+err.Error(), http.StatusInternalServerError)
		return
	}

	abilities, err := parseAbilities(fetched["Abilities"].body)
	if err != nil {
		jsonError(w, "parse abilities: "+err.Error(), http.StatusInternalServerError)
		return
	}

	enhancements, err := parseEnhancements(fetched["Enhancements"].body)
	if err != nil {
		jsonError(w, "parse enhancements: "+err.Error(), http.StatusInternalServerError)
		return
	}

	detachments, err := parseDetachments(fetched["Detachments"].body)
	if err != nil {
		jsonError(w, "parse detachments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	detachmentAbilities, err := parseDetachmentAbilities(fetched["Detachment_abilities"].body)
	if err != nil {
		jsonError(w, "parse detachment_abilities: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsModels, err := parseDatasheetModels(fetched["Datasheets_models"].body)
	if err != nil {
		jsonError(w, "parse models: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsAbilities, err := parseDatasheetAbilities(fetched["Datasheets_abilities"].body)
	if err != nil {
		jsonError(w, "parse datasheet_abilities: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsKeywords, err := parseDatasheetKeywords(fetched["Datasheets_keywords"].body)
	if err != nil {
		jsonError(w, "parse datasheet_keywords: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsOptions, err := parseDatasheetOptions(fetched["Datasheets_options"].body)
	if err != nil {
		jsonError(w, "parse datasheet_options: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsWargear, err := parseDatasheetWargear(fetched["Datasheets_wargear"].body)
	if err != nil {
		jsonError(w, "parse datasheet_wargear: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsUnitComposition, err := parseDatasheetUnitComposition(fetched["Datasheets_unit_composition"].body)
	if err != nil {
		jsonError(w, "parse datasheet_unit_composition: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsModelsCost, err := parseDatasheetModelsCost(fetched["Datasheets_models_cost"].body)
	if err != nil {
		jsonError(w, "parse datasheet_models_cost: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsStratagems, err := parseDatasheetStratagems(fetched["Datasheets_stratagems"].body)
	if err != nil {
		jsonError(w, "parse datasheet_stratagems: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsEnhancements, err := parseDatasheetEnhancements(fetched["Datasheets_enhancements"].body)
	if err != nil {
		jsonError(w, "parse datasheet_enhancements: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsDetachmentAbilities, err := parseDatasheetDetachmentAbilities(fetched["Datasheets_detachment_abilities"].body)
	if err != nil {
		jsonError(w, "parse datasheet_detachment_abilities: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dsLeaders, err := parseDatasheetLeaders(fetched["Datasheets_leader"].body)
	if err != nil {
		jsonError(w, "parse datasheet_leader: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build sets for validation/filtering
	factionSet := make(map[string]struct{}, len(factions))
	for _, f := range factions {
		factionSet[f.ID] = struct{}{}
	}

	sourceSet := make(map[string]struct{}, len(sources))
	for _, s := range sources {
		sourceSet[s.ID] = struct{}{}
	}

	stratagemSet := make(map[string]struct{}, len(stratagems))
	enhancementSet := make(map[string]struct{}, len(enhancements))
	detachmentAbilitySet := make(map[string]struct{}, len(detachmentAbilities))

	// Filter datasheets by valid faction
	var validDatasheets []models.WhDatasheet
	for _, d := range datasheets {
		if _, ok := factionSet[d.FactionID]; ok {
			validDatasheets = append(validDatasheets, d)
		} else {
			log.Printf("[wahapedia] skipping datasheet %s: unknown faction_id %s", d.ID, d.FactionID)
		}
	}

	// Build set of valid datasheet IDs
	dsSet := make(map[string]struct{}, len(validDatasheets))
	for _, d := range validDatasheets {
		dsSet[d.ID] = struct{}{}
	}

	// Filter global entities first, then populate their sets
	var validStratagems []models.WhStratagem
	for _, s := range stratagems {
		if _, ok := factionSet[s.FactionID]; ok {
			validStratagems = append(validStratagems, s)
			stratagemSet[s.ID] = struct{}{}
		}
	}

	var validAbilities []models.WhAbility
	for _, a := range abilities {
		if _, ok := factionSet[a.FactionID]; ok {
			validAbilities = append(validAbilities, a)
		}
	}

	var validEnhancements []models.WhEnhancement
	for _, e := range enhancements {
		if _, ok := factionSet[e.FactionID]; ok {
			validEnhancements = append(validEnhancements, e)
			enhancementSet[e.ID] = struct{}{}
		}
	}

	var validDetachments []models.WhDetachment
	for _, d := range detachments {
		if _, ok := factionSet[d.FactionID]; ok {
			validDetachments = append(validDetachments, d)
		}
	}

	var validDetachmentAbilities []models.WhDetachmentAbility
	for _, da := range detachmentAbilities {
		if _, ok := factionSet[da.FactionID]; ok {
			validDetachmentAbilities = append(validDetachmentAbilities, da)
			detachmentAbilitySet[da.ID] = struct{}{}
		}
	}

	// Filter datasheet children by valid datasheet
	var validModels []models.WhDatasheetModel
	for _, m := range dsModels {
		if _, ok := dsSet[m.DatasheetID]; ok {
			validModels = append(validModels, m)
		}
	}

	var validDsAbilities []models.WhDatasheetAbility
	for _, a := range dsAbilities {
		if _, ok := dsSet[a.DatasheetID]; ok {
			validDsAbilities = append(validDsAbilities, a)
		}
	}

	var validDsKeywords []models.WhDatasheetKeyword
	for _, k := range dsKeywords {
		if _, ok := dsSet[k.DatasheetID]; ok {
			validDsKeywords = append(validDsKeywords, k)
		}
	}

	var validDsOptions []models.WhDatasheetOption
	for _, o := range dsOptions {
		if _, ok := dsSet[o.DatasheetID]; ok {
			validDsOptions = append(validDsOptions, o)
		}
	}

	var validDsWargear []models.WhDatasheetWargear
	for _, w := range dsWargear {
		if _, ok := dsSet[w.DatasheetID]; ok {
			validDsWargear = append(validDsWargear, w)
		}
	}

	var validDsUnitComposition []models.WhDatasheetUnitComposition
	for _, uc := range dsUnitComposition {
		if _, ok := dsSet[uc.DatasheetID]; ok {
			validDsUnitComposition = append(validDsUnitComposition, uc)
		}
	}

	var validDsModelsCost []models.WhDatasheetModelCost
	for _, mc := range dsModelsCost {
		if _, ok := dsSet[mc.DatasheetID]; ok {
			validDsModelsCost = append(validDsModelsCost, mc)
		}
	}

	var validDsStratagems []models.WhDatasheetStratagem
	for _, ds := range dsStratagems {
		if _, okDs := dsSet[ds.DatasheetID]; okDs {
			if _, okStrat := stratagemSet[ds.StratagemID]; okStrat {
				validDsStratagems = append(validDsStratagems, ds)
			}
		}
	}

	var validDsEnhancements []models.WhDatasheetEnhancement
	for _, de := range dsEnhancements {
		if _, okDs := dsSet[de.DatasheetID]; okDs {
			if _, okEnh := enhancementSet[de.EnhancementID]; okEnh {
				validDsEnhancements = append(validDsEnhancements, de)
			}
		}
	}

	var validDsDetachmentAbilities []models.WhDatasheetDetachmentAbility
	for _, da := range dsDetachmentAbilities {
		if _, okDs := dsSet[da.DatasheetID]; okDs {
			if _, okDa := detachmentAbilitySet[da.DetachmentAbilityID]; okDa {
				validDsDetachmentAbilities = append(validDsDetachmentAbilities, da)
			}
		}
	}

	var validDsLeaders []models.WhDatasheetLeader
	for _, l := range dsLeaders {
		if _, ok1 := dsSet[l.DatasheetID]; ok1 {
			if _, ok2 := dsSet[l.AttachedDatasheetID]; ok2 {
				validDsLeaders = append(validDsLeaders, l)
			}
		}
	}

	// Store everything atomically
	if err := db.SyncWahapediaData(
		factions, sources, validDatasheets, validModels, validDsAbilities,
		validDsKeywords, validDsOptions, validDsWargear, validDsUnitComposition,
		validDsModelsCost, validDsStratagems, validDsEnhancements,
		validDsDetachmentAbilities, validDsLeaders, validStratagems,
		validAbilities, validEnhancements, validDetachmentAbilities,
		validDetachments, hashes,
	); err != nil {
		jsonError(w, "sync db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[wahapedia] sync complete: %d factions, %d sources, %d datasheets, %d models, %d abilities, %d stratagems, %d enhancements, %d detachments",
		len(factions), len(sources), len(validDatasheets), len(validModels), len(validAbilities), len(validStratagems), len(validEnhancements), len(validDetachments))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.WhSyncResult{
		Changed:        true,
		UpdatedSources: changed,
	})
}
