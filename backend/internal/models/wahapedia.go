package models

type WhFaction struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type WhDatasheet struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	FactionID           string `json:"faction_id"`
	SourceID            string `json:"source_id"`
	Legend              string `json:"legend"`
	Role                string `json:"role"`
	Loadout             string `json:"loadout"`
	Transport           string `json:"transport"`
	Virtual             bool   `json:"virtual"`
	LeaderHead          string `json:"leader_head"`
	LeaderFooter        string `json:"leader_footer"`
	DamagedW            string `json:"damaged_w"`
	DamagedDescription  string `json:"damaged_description"`
	Link                string `json:"link"`
}

type WhDatasheetModel struct {
	DatasheetID  string `json:"datasheet_id"`
	Name         string `json:"name"`
	M            string `json:"m"`
	T            string `json:"t"`
	SV           string `json:"sv"`
	InvSV        string `json:"inv_sv"`
	InvSVDescr   string `json:"inv_sv_descr"`
	W            string `json:"w"`
	LD           string `json:"ld"`
	OC           string `json:"oc"`
	BaseSize     string `json:"base_size"`
	BaseSizeDescr string `json:"base_size_descr"`
}

type WhDatasheetWeapon struct {
	DatasheetID string `json:"datasheet_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Range       string `json:"range"`
	A           string `json:"a"`
	BS          string `json:"bs"`
	S           string `json:"s"`
	AP          string `json:"ap"`
	D           string `json:"d"`
	Abilities   string `json:"abilities"`
}

type WhDatasheetAbility struct {
	DatasheetID string `json:"datasheet_id"`
	Line        string `json:"line"`
	AbilityID   string `json:"ability_id"`
	Model       string `json:"model"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Parameter   string `json:"parameter"`
}

// Shared global ability
type WhAbility struct {
	ID          string `json:"id"`
	FactionID   string `json:"faction_id"`
	Name        string `json:"name"`
	Legend      string `json:"legend"`
	Description string `json:"description"`
}

// Source / supplement
type WhSource struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Edition     string `json:"edition"`
	Version     string `json:"version"`
	ErrataDate  string `json:"errata_date"`
	ErrataLink  string `json:"errata_link"`
}

// Stratagem
type WhStratagem struct {
	ID            string `json:"id"`
	FactionID     string `json:"faction_id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	CPCost        string `json:"cp_cost"`
	Legend        string `json:"legend"`
	Turn          string `json:"turn"`
	Phase         string `json:"phase"`
	Description   string `json:"description"`
	Detachment    string `json:"detachment"`
	DetachmentID  string `json:"detachment_id"`
}

// Enhancement
type WhEnhancement struct {
	ID           string `json:"id"`
	FactionID    string `json:"faction_id"`
	Name         string `json:"name"`
	Legend       string `json:"legend"`
	Description  string `json:"description"`
	Cost         string `json:"cost"`
	Detachment   string `json:"detachment"`
	DetachmentID string `json:"detachment_id"`
}

// Detachment
type WhDetachment struct {
	ID        string `json:"id"`
	FactionID string `json:"faction_id"`
	Name      string `json:"name"`
	Legend    string `json:"legend"`
	Type      string `json:"type"`
}

// Detachment Ability
type WhDetachmentAbility struct {
	ID           string `json:"id"`
	FactionID    string `json:"faction_id"`
	Name         string `json:"name"`
	Legend       string `json:"legend"`
	Description  string `json:"description"`
	Detachment   string `json:"detachment"`
	DetachmentID string `json:"detachment_id"`
}

// Datasheet Keywords
type WhDatasheetKeyword struct {
	DatasheetID       string `json:"datasheet_id"`
	Line              string `json:"line"`
	Keyword           string `json:"keyword"`
	Model             string `json:"model"`
	IsFactionKeyword  bool   `json:"is_faction_keyword"`
}

// Datasheet Wargear Options
type WhDatasheetOption struct {
	DatasheetID string `json:"datasheet_id"`
	Line        string `json:"line"`
	Button      string `json:"button"`
	Description string `json:"description"`
}

// Datasheet Wargear (weapons)
type WhDatasheetWargear struct {
	DatasheetID   string `json:"datasheet_id"`
	Line          string `json:"line"`
	LineInWargear string `json:"line_in_wargear"`
	Dice          string `json:"dice"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Range         string `json:"range"`
	Type          string `json:"type"`
	A             string `json:"a"`
	BSWS          string `json:"bs_ws"`
	S             string `json:"s"`
	AP            string `json:"ap"`
	D             string `json:"d"`
}

// Datasheet Unit Composition
type WhDatasheetUnitComposition struct {
	DatasheetID string `json:"datasheet_id"`
	Line        string `json:"line"`
	Description string `json:"description"`
}

// Datasheet Model Cost
type WhDatasheetModelCost struct {
	DatasheetID string `json:"datasheet_id"`
	Line        string `json:"line"`
	Description string `json:"description"`
	Cost        string `json:"cost"`
}

// Junction tables

// Datasheet Stratagem junction
type WhDatasheetStratagem struct {
	DatasheetID string `json:"datasheet_id"`
	StratagemID string `json:"stratagem_id"`
}

// Datasheet Enhancement junction
type WhDatasheetEnhancement struct {
	DatasheetID   string `json:"datasheet_id"`
	EnhancementID string `json:"enhancement_id"`
}

// Datasheet Detachment Ability junction
type WhDatasheetDetachmentAbility struct {
	DatasheetID          string `json:"datasheet_id"`
	DetachmentAbilityID  string `json:"detachment_ability_id"`
}

// Datasheet Leader junction
type WhDatasheetLeader struct {
	DatasheetID        string `json:"datasheet_id"`
	AttachedDatasheetID string `json:"attached_datasheet_id"`
}

type WhSyncResult struct {
	Changed        bool     `json:"changed"`
	UpdatedSources []string `json:"updated_sources,omitempty"`
}
