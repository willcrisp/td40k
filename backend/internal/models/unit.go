package models

import "math"

// ─── Enums ────────────────────────────────────────────────────────────────────

type FootprintShape string

const (
	FootprintShapeRound FootprintShape = "round"
	FootprintShapeOval  FootprintShape = "oval"
	FootprintShapeHull  FootprintShape = "hull"
)

type UnitStatus string

const (
	UnitStatusAlive      UnitStatus = "alive"
	UnitStatusInReserves UnitStatus = "in_reserves"
	UnitStatusDead       UnitStatus = "dead"
)

// ─── Stats ────────────────────────────────────────────────────────────────────

type UnitStats struct {
	// Movement in inches
	Movement  int `json:"movement"`
	Toughness int `json:"toughness"`
	// Armour save threshold e.g. 3 = "3+"
	Save int `json:"save"`
	// Invulnerable save e.g. 4 = "4+", 0 = none
	InvulnerableSave int `json:"invulnerable_save,omitempty"`
	// Maximum wounds — never changes after construction
	Wounds int `json:"wounds"`
	// Leadership threshold e.g. 7 = "7+"
	Leadership int `json:"leadership"`
	// Objective Control value
	ObjectiveControl int `json:"objective_control"`
}

// ─── Footprint ────────────────────────────────────────────────────────────────

type Footprint struct {
	// Width in mm
	X float64 `json:"x"`
	// Depth in mm — equal to X for round bases
	Y float64 `json:"y"`
	// False for vehicles and anything measured hull-to-hull
	HasBase bool `json:"has_base"`
}

// Shape derives the rendering shape from footprint dimensions.
func (f Footprint) Shape() FootprintShape {
	if !f.HasBase {
		return FootprintShapeHull
	}
	if f.X == f.Y {
		return FootprintShapeRound
	}
	return FootprintShapeOval
}

const mmPerInch = 25.4

// InInches returns footprint dimensions converted to inches for board-space math.
func (f Footprint) InInches() (x, y float64) {
	return f.X / mmPerInch, f.Y / mmPerInch
}

// ─── Board Position ───────────────────────────────────────────────────────────

type BoardPosition struct {
	// X position on board in inches from top-left origin
	X float64 `json:"x"`
	// Y position on board in inches from top-left origin
	Y float64 `json:"y"`
	// Facing in degrees. 0 = north (up), clockwise positive.
	// Used for line of sight, charge arcs, etc.
	Facing float64 `json:"facing"`
}

// ─── Unit Interface ───────────────────────────────────────────────────────────

type Unit interface {
	GetName() string
	GetFaction() string
	GetKeywords() []string
	GetStats() UnitStats
	GetFootprint() Footprint
	GetPosition() BoardPosition
	GetStatus() UnitStatus
	GetCurrentWounds() int
	MoveTo(x, y float64)
	FaceTo(degrees float64)
	RotateBy(degrees float64)
	ApplyWounds(amount int)
	SendToReserves()
	DeployFromReserves(x, y, facing float64)
}

// ─── Base Unit ────────────────────────────────────────────────────────────────

// BaseUnit holds all common data and behaviour.
// Embed this into concrete unit structs and override the identity fields.
type BaseUnit struct {
	Name     string    `json:"name"`
	Faction  string    `json:"faction"`
	Keywords []string  `json:"keywords"`
	Stats    UnitStats `json:"stats"`
	Footprint Footprint     `json:"footprint"`
	Position  BoardPosition `json:"position"`
	Status    UnitStatus    `json:"status"`
	// CurrentWounds tracks damage state — starts equal to Stats.Wounds
	CurrentWounds int `json:"current_wounds"`
}

// NewBaseUnit constructs a BaseUnit with wounds and status initialised.
func NewBaseUnit(
	name, faction string,
	keywords []string,
	stats UnitStats,
	footprint Footprint,
) BaseUnit {
	return BaseUnit{
		Name:          name,
		Faction:       faction,
		Keywords:      keywords,
		Stats:         stats,
		Footprint:     footprint,
		Position:      BoardPosition{},
		Status:        UnitStatusAlive,
		CurrentWounds: stats.Wounds,
	}
}

func (u *BaseUnit) GetName() string            { return u.Name }
func (u *BaseUnit) GetFaction() string         { return u.Faction }
func (u *BaseUnit) GetKeywords() []string      { return u.Keywords }
func (u *BaseUnit) GetStats() UnitStats        { return u.Stats }
func (u *BaseUnit) GetFootprint() Footprint    { return u.Footprint }
func (u *BaseUnit) GetPosition() BoardPosition { return u.Position }
func (u *BaseUnit) GetStatus() UnitStatus      { return u.Status }
func (u *BaseUnit) GetCurrentWounds() int      { return u.CurrentWounds }

// ─── Movement ─────────────────────────────────────────────────────────────────

func (u *BaseUnit) MoveTo(x, y float64) {
	u.Position.X = x
	u.Position.Y = y
}

// FaceTo sets an absolute facing direction, normalised to 0–359.
func (u *BaseUnit) FaceTo(degrees float64) {
	u.Position.Facing = math.Mod(math.Mod(degrees, 360)+360, 360)
}

// RotateBy adjusts facing by a relative amount.
func (u *BaseUnit) RotateBy(degrees float64) {
	u.FaceTo(u.Position.Facing + degrees)
}

// ─── Wounds & Status ──────────────────────────────────────────────────────────

// ApplyWounds reduces current wounds by amount.
// If wounds reach 0 the unit is marked Dead.
// Has no effect on units that are already Dead.
func (u *BaseUnit) ApplyWounds(amount int) {
	if u.Status == UnitStatusDead {
		return
	}
	u.CurrentWounds -= amount
	if u.CurrentWounds <= 0 {
		u.CurrentWounds = 0
		u.Status = UnitStatusDead
	}
}

// SendToReserves removes the unit from the board.
// Only valid for Alive units.
func (u *BaseUnit) SendToReserves() {
	if u.Status != UnitStatusAlive {
		return
	}
	u.Status = UnitStatusInReserves
}

// DeployFromReserves places the unit on the board at the given position.
// Only valid for units that are InReserves.
func (u *BaseUnit) DeployFromReserves(x, y, facing float64) {
	if u.Status != UnitStatusInReserves {
		return
	}
	u.MoveTo(x, y)
	u.FaceTo(facing)
	u.Status = UnitStatusAlive
}
