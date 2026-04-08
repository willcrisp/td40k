package models

import (
	"strconv"
	"strings"
)

// ConvertDatasheetModelToUnit converts Wahapedia WhDatasheetModel and WhDatasheet
// to a playable BaseUnit with stats parsed from string values and footprint
// derived from base_size string.
func ConvertDatasheetModelToUnit(
	datasheet *WhDatasheet,
	model *WhDatasheetModel,
) BaseUnit {
	// Parse stats from string values
	stats := parseUnitStats(model)

	// Derive footprint from base size string
	footprint := parseFootprint(model.BaseSize)

	// Build keywords from faction and datasheet metadata
	keywords := []string{}
	if datasheet.Role != "" {
		keywords = append(keywords, strings.ToLower(datasheet.Role))
	}
	// Additional keywords could be extracted from abilities, but keeping it simple

	return NewBaseUnit(
		model.Name,
		datasheet.FactionID,
		keywords,
		stats,
		footprint,
	)
}

// parseUnitStats converts string stat values from Wahapedia to UnitStats integers.
// Handles values like "6", "4", "3+", "4++", etc.
func parseUnitStats(model *WhDatasheetModel) UnitStats {
	return UnitStats{
		Movement:         parseStatInt(model.M),
		Toughness:        parseStatInt(model.T),
		Save:             parseSaveValue(model.SV),
		InvulnerableSave: parseSaveValue(model.InvSV),
		Wounds:           parseStatInt(model.W),
		Leadership:       parseSaveValue(model.LD),
		ObjectiveControl: parseStatInt(model.OC),
	}
}

// parseStatInt converts simple numeric stat strings like "6", "4", "2"
func parseStatInt(s string) int {
	if s == "" || s == "-" {
		return 0
	}
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return v
}

// parseSaveValue handles save notation like "3+", "4++", "5+++"
// Returns just the numeric threshold (3 for "3+", 4 for "4++", etc.)
func parseSaveValue(s string) int {
	if s == "" || s == "-" {
		return 0
	}

	// Remove all non-digit characters
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)

	v, err := strconv.Atoi(digits)
	if err != nil {
		return 0
	}
	return v
}

// parseFootprint derives Footprint from Wahapedia base_size string.
// Typical values: "25mm", "32mm", "40mm", "60mm", "Hull", "N/A", etc.
// Returns:
// - For round bases (25/32/40/60mm): Footprint{X: mm, Y: mm, HasBase: true}
// - For oval/rectangular: Footprint{X, Y, HasBase: true} if dimensions available
// - For vehicles (Hull): Footprint{0, 0, HasBase: false} (dimensions set elsewhere)
// - Default: 32mm round base if unparseable
func parseFootprint(baseSize string) Footprint {
	if baseSize == "" || baseSize == "-" || baseSize == "N/A" {
		// Default to 32mm round base
		return Footprint{X: 32, Y: 32, HasBase: true}
	}

	baseSize = strings.TrimSpace(strings.ToLower(baseSize))

	// Handle standard round bases
	switch baseSize {
	case "25mm":
		return Footprint{X: 25, Y: 25, HasBase: true}
	case "32mm":
		return Footprint{X: 32, Y: 32, HasBase: true}
	case "40mm":
		return Footprint{X: 40, Y: 40, HasBase: true}
	case "50mm":
		return Footprint{X: 50, Y: 50, HasBase: true}
	case "60mm":
		return Footprint{X: 60, Y: 60, HasBase: true}
	case "hull", "hull-mounted":
		// Vehicles measured hull-to-hull; dimensions not determined by base size alone
		return Footprint{X: 0, Y: 0, HasBase: false}
	}

	// Try to parse custom dimensions like "50x50mm", "80x60mm"
	if strings.Contains(baseSize, "x") {
		parts := strings.Split(baseSize, "x")
		if len(parts) == 2 {
			x := parseFootprintDim(parts[0])
			y := parseFootprintDim(parts[1])
			if x > 0 && y > 0 {
				return Footprint{X: float64(x), Y: float64(y), HasBase: true}
			}
		}
	}

	// Try to parse single dimension like "32 mm" or just "32"
	if dim := parseFootprintDim(baseSize); dim > 0 {
		return Footprint{X: float64(dim), Y: float64(dim), HasBase: true}
	}

	// Default fallback
	return Footprint{X: 32, Y: 32, HasBase: true}
}

// parseFootprintDim extracts a numeric dimension from strings like "32mm", "50", "80 mm"
func parseFootprintDim(s string) int {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "mm")
	s = strings.TrimSpace(s)

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}
