package handlers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// bsRoster is the top-level shape of a ListForge / BattleScribe JSON export.
type bsRoster struct {
	Roster bsRosterInner `json:"roster"`
}

type bsRosterInner struct {
	Name   string    `json:"name"`
	Costs  []bsCost  `json:"costs"`
	Forces []bsForce `json:"forces"`
}

type bsForce struct {
	Name          string        `json:"name"`
	CatalogueName string        `json:"catalogueName"`
	Selections    []bsSelection `json:"selections"`
}

type bsCost struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type bsCategory struct {
	Name    string `json:"name"`
	Primary bool   `json:"primary"`
}

type bsSelection struct {
	Name       string        `json:"name"`
	Number     int           `json:"number"`
	Type       string        `json:"type"` // "unit" | "model" | "upgrade"
	Costs      []bsCost      `json:"costs"`
	Categories []bsCategory  `json:"categories"`
	Selections []bsSelection `json:"selections"`
}

// ParsedUnit is a single deployable unit extracted from the import.
type ParsedUnit struct {
	Name       string
	Points     int
	ModelCount int
	Quantity   int
}

// ParsedRoster is the structured output of parseListForge.
type ParsedRoster struct {
	ListName    string
	FactionName string
	TotalPoints int
	Units       []ParsedUnit
}

var nonAlphanumRe = regexp.MustCompile(`[^a-z0-9 ]`)

// NormalizeName strips punctuation and lowercases a name for fuzzy matching.
func NormalizeName(s string) string {
	s = strings.ToLower(s)
	s = nonAlphanumRe.ReplaceAllString(s, "")
	return strings.Join(strings.Fields(s), " ")
}

// parseListForge decodes a ListForge BattleScribe JSON export and returns the
// list of deployable units with their point costs and model counts.
func parseListForge(data []byte) (ParsedRoster, error) {
	var raw bsRoster
	if err := json.Unmarshal(data, &raw); err != nil {
		return ParsedRoster{}, fmt.Errorf("invalid JSON: %w", err)
	}

	inner := raw.Roster
	if len(inner.Forces) == 0 {
		return ParsedRoster{}, fmt.Errorf("no forces found in roster")
	}

	force := inner.Forces[0]

	totalPoints := 0
	for _, c := range inner.Costs {
		if strings.EqualFold(c.Name, "pts") {
			totalPoints = int(c.Value)
			break
		}
	}

	// De-duplicate by unit name: same name → increment Quantity, sum counts.
	type accumEntry struct {
		points     int
		modelCount int
		quantity   int
	}
	order := []string{}
	accum := map[string]*accumEntry{}

	for _, sel := range force.Selections {
		if sel.Type != "unit" && sel.Type != "model" {
			continue
		}
		// Skip configuration meta-entries (Detachment, Battle Size, etc.)
		if primaryCategory(sel) == "Configuration" {
			continue
		}

		name := strings.TrimSpace(sel.Name)
		if name == "" {
			continue
		}

		pts := 0
		for _, c := range sel.Costs {
			if strings.EqualFold(c.Name, "pts") {
				pts = int(c.Value)
				break
			}
		}

		// ModelCount: sum of sub-selection numbers where type == "model".
		// If no model sub-selections, use the selection's own Number.
		modelCount := 0
		for _, sub := range sel.Selections {
			if sub.Type == "model" {
				n := sub.Number
				if n <= 0 {
					n = 1
				}
				modelCount += n
			}
		}
		if modelCount == 0 {
			n := sel.Number
			if n <= 0 {
				n = 1
			}
			modelCount = n
		}

		if e, exists := accum[name]; exists {
			e.quantity++
			e.modelCount += modelCount
			e.points += pts
		} else {
			order = append(order, name)
			accum[name] = &accumEntry{
				points:     pts,
				modelCount: modelCount,
				quantity:   1,
			}
		}
	}

	units := make([]ParsedUnit, 0, len(order))
	for _, name := range order {
		e := accum[name]
		units = append(units, ParsedUnit{
			Name:       name,
			Points:     e.points,
			ModelCount: e.modelCount,
			Quantity:   e.quantity,
		})
	}

	return ParsedRoster{
		ListName:    inner.Name,
		FactionName: force.CatalogueName,
		TotalPoints: totalPoints,
		Units:       units,
	}, nil
}

// primaryCategory returns the name of the first primary category on a
// selection, or an empty string if there are none.
func primaryCategory(sel bsSelection) string {
	for _, cat := range sel.Categories {
		if cat.Primary {
			return cat.Name
		}
	}
	return ""
}
