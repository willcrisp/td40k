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
	"Factions":              "Factions.csv",
	"Datasheets":            "Datasheets.csv",
	"Datasheets_models":     "Datasheets_models.csv",
	"Datasheets_weapons":    "Datasheets_weapons.csv",
	"Datasheets_abilities":  "Datasheets_abilities.csv",
}

const wahapediaBaseURL = "https://wahapedia.ru/wh40k10ed/data/"

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
		out = append(out, models.WhDatasheet{
			ID:        id,
			Name:      field(row, cols, "name"),
			FactionID: field(row, cols, "faction_id"),
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
			DatasheetID: dsID,
			Name:        field(row, cols, "name"),
			M:           field(row, cols, "M"),
			T:           field(row, cols, "T"),
			SV:          field(row, cols, "Sv"),
			W:           field(row, cols, "W"),
			LD:          field(row, cols, "Ld"),
			OC:          field(row, cols, "OC"),
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
			Name:        field(row, cols, "name"),
			Description: field(row, cols, "description"),
		})
	}
	return out, nil
}

func HandleSyncWahapedia(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	datasheets, err := parseDatasheets(fetched["Datasheets"].body)
	if err != nil {
		jsonError(w, "parse datasheets: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build a set of valid faction IDs for filtering orphan datasheets
	factionSet := make(map[string]struct{}, len(factions))
	for _, f := range factions {
		factionSet[f.ID] = struct{}{}
	}
	var validDatasheets []models.WhDatasheet
	for _, d := range datasheets {
		if _, ok := factionSet[d.FactionID]; ok {
			validDatasheets = append(validDatasheets, d)
		} else {
			log.Printf("[wahapedia] skipping datasheet %s: unknown faction_id %s", d.ID, d.FactionID)
		}
	}

	// Build a set of valid datasheet IDs for filtering orphan children
	dsSet := make(map[string]struct{}, len(validDatasheets))
	for _, d := range validDatasheets {
		dsSet[d.ID] = struct{}{}
	}

	dsModels, err := parseDatasheetModels(fetched["Datasheets_models"].body)
	if err != nil {
		jsonError(w, "parse models: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var validModels []models.WhDatasheetModel
	for _, m := range dsModels {
		if _, ok := dsSet[m.DatasheetID]; ok {
			validModels = append(validModels, m)
		}
	}

	weapons, err := parseDatasheetWeapons(fetched["Datasheets_weapons"].body)
	if err != nil {
		jsonError(w, "parse weapons: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var validWeapons []models.WhDatasheetWeapon
	for _, wp := range weapons {
		if _, ok := dsSet[wp.DatasheetID]; ok {
			validWeapons = append(validWeapons, wp)
		}
	}

	abilities, err := parseDatasheetAbilities(fetched["Datasheets_abilities"].body)
	if err != nil {
		jsonError(w, "parse abilities: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var validAbilities []models.WhDatasheetAbility
	for _, a := range abilities {
		if _, ok := dsSet[a.DatasheetID]; ok {
			validAbilities = append(validAbilities, a)
		}
	}

	// Store everything atomically
	if err := db.SyncWahapediaData(
		factions, validDatasheets, validModels, validWeapons, validAbilities, hashes,
	); err != nil {
		jsonError(w, "sync db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[wahapedia] sync complete: %d factions, %d datasheets, %d models, %d weapons, %d abilities",
		len(factions), len(validDatasheets), len(validModels), len(validWeapons), len(validAbilities))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.WhSyncResult{
		Changed:        true,
		UpdatedSources: changed,
	})
}
