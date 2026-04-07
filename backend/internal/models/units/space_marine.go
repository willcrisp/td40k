package units

import "github.com/willcrisp/td40k/internal/models"

type SpaceMarine struct {
	models.BaseUnit
}

func NewSpaceMarine() *SpaceMarine {
	return &SpaceMarine{
		BaseUnit: models.NewBaseUnit(
			"Space Marine",
			"Adeptus Astartes",
			[]string{"Infantry", "Imperium", "Adeptus Astartes"},
			models.UnitStats{
				Movement: 6, Toughness: 4, Save: 3,
				Wounds: 2, Leadership: 6, ObjectiveControl: 2,
			},
			models.Footprint{X: 32, Y: 32, HasBase: true},
		),
	}
}
