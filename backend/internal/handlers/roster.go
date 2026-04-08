package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/willcrisp/td40k/internal/db"
	"github.com/willcrisp/td40k/internal/models"
	mw "github.com/willcrisp/td40k/internal/middleware"
)

// MatchedUnit is a unit from the import that was successfully matched to a
// Wahapedia datasheet.
type MatchedUnit struct {
	Name        string `json:"name"`
	DatasheetID string `json:"datasheet_id"`
	FactionID   string `json:"faction_id"`
	Quantity    int    `json:"quantity"`
	ModelCount  int    `json:"model_count"`
	Points      int    `json:"points"`
}

// ImportRosterResponse is returned by POST /roster/import.
type ImportRosterResponse struct {
	FactionName string        `json:"faction_name"`
	TotalPoints int           `json:"total_points"`
	Matched     []MatchedUnit `json:"matched"`
	Unmatched   []string      `json:"unmatched"`
}

// HandleImportRoster parses a ListForge JSON export, matches units against
// Wahapedia datasheets, and saves the result to the unit_roster table.
func HandleImportRoster(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "roomId")

	room, err := db.GetRoom(roomID)
	if err != nil {
		jsonError(w, "room not found", http.StatusNotFound)
		return
	}
	if !isInRoom(room, playerID) {
		jsonError(w, "not in room", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, "failed to read body", http.StatusBadRequest)
		return
	}

	parsed, err := parseListForge(body)
	if err != nil {
		jsonError(w, "invalid listforge json: "+err.Error(), http.StatusBadRequest)
		return
	}

	if parsed.FactionName == "" {
		jsonError(w, "no faction found in roster", http.StatusBadRequest)
		return
	}

	faction, err := db.GetFactionByName(parsed.FactionName)
	if err != nil {
		jsonError(w, "db error looking up faction", http.StatusInternalServerError)
		return
	}
	if faction == nil {
		jsonError(w, "faction not found: "+parsed.FactionName, http.StatusNotFound)
		return
	}

	datasheets, err := db.GetDatasheetsByFaction(faction.ID)
	if err != nil {
		jsonError(w, "db error fetching datasheets", http.StatusInternalServerError)
		return
	}

	matched, unmatched := matchUnits(parsed.Units, datasheets, faction.ID)

	entries := make([]db.RosterEntry, 0, len(matched))
	for _, m := range matched {
		entries = append(entries, db.RosterEntry{
			DatasheetID: m.DatasheetID,
			ModelName:   m.Name,
			Quantity:    m.Quantity,
			FactionID:   m.FactionID,
			Points:      m.Points,
		})
	}

	if _, err := db.BulkCreateRosterEntries(roomID, playerID, entries); err != nil {
		jsonError(w, "failed to save roster", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ImportRosterResponse{
		FactionName: faction.Name,
		TotalPoints: parsed.TotalPoints,
		Matched:     matched,
		Unmatched:   unmatched,
	})
}

// HandleGetRoster returns the current roster for the authenticated player.
func HandleGetRoster(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "roomId")

	room, err := db.GetRoom(roomID)
	if err != nil {
		jsonError(w, "room not found", http.StatusNotFound)
		return
	}
	if !isInRoom(room, playerID) {
		jsonError(w, "not in room", http.StatusForbidden)
		return
	}

	entries, err := db.GetRoster(roomID, playerID)
	if err != nil {
		jsonError(w, "failed to fetch roster", http.StatusInternalServerError)
		return
	}
	if entries == nil {
		entries = []db.RosterEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// HandleClearRoster deletes all roster entries for the authenticated player.
func HandleClearRoster(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "roomId")

	room, err := db.GetRoom(roomID)
	if err != nil {
		jsonError(w, "room not found", http.StatusNotFound)
		return
	}
	if !isInRoom(room, playerID) {
		jsonError(w, "not in room", http.StatusForbidden)
		return
	}

	if err := db.ClearRoster(roomID, playerID); err != nil {
		jsonError(w, "failed to clear roster", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func isInRoom(room *models.Room, playerID string) bool {
	return (room.AttackerID != nil && *room.AttackerID == playerID) ||
		(room.DefenderID != nil && *room.DefenderID == playerID) ||
		room.GameMasterID == playerID
}

// matchUnits matches parsed units against Wahapedia datasheets, returning the
// matched and unmatched unit names. Matching is done on normalised names.
func matchUnits(
	units []ParsedUnit,
	datasheets []models.WhDatasheet,
	factionID string,
) ([]MatchedUnit, []string) {
	// Build normalized name → datasheet map
	normMap := make(map[string]models.WhDatasheet, len(datasheets))
	for _, ds := range datasheets {
		normMap[NormalizeName(ds.Name)] = ds
	}

	var matched []MatchedUnit
	var unmatched []string

	for _, u := range units {
		normUnit := NormalizeName(u.Name)

		// 1. Exact normalized match
		if ds, ok := normMap[normUnit]; ok {
			matched = append(matched, MatchedUnit{
				Name:        u.Name,
				DatasheetID: ds.ID,
				FactionID:   factionID,
				Quantity:    u.Quantity,
				ModelCount:  u.ModelCount,
				Points:      u.Points,
			})
			continue
		}

		// 2. Substring match
		found := false
		for normDS, ds := range normMap {
			if strings.Contains(normDS, normUnit) || strings.Contains(normUnit, normDS) {
				matched = append(matched, MatchedUnit{
					Name:        u.Name,
					DatasheetID: ds.ID,
					FactionID:   factionID,
					Quantity:    u.Quantity,
					ModelCount:  u.ModelCount,
					Points:      u.Points,
				})
				found = true
				break
			}
		}
		if !found {
			unmatched = append(unmatched, u.Name)
		}
	}

	return matched, unmatched
}
