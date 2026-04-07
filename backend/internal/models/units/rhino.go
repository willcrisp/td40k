package units

import "github.com/willcrisp/td40k/internal/models"

type Rhino struct {
	models.BaseUnit
}

func NewRhino() *Rhino {
	return &Rhino{
		BaseUnit: models.NewBaseUnit(
			"Rhino",
			"Adeptus Astartes",
			[]string{"Vehicle", "Transport", "Imperium", "Adeptus Astartes"},
			models.UnitStats{
				Movement: 12, Toughness: 9, Save: 3,
				Wounds: 10, Leadership: 6, ObjectiveControl: 2,
			},
			// No base — hull footprint used directly for rendering
			models.Footprint{X: 105, Y: 70, HasBase: false},
		),
	}
}
