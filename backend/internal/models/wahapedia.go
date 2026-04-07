package models

type WhFaction struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type WhDatasheet struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	FactionID string `json:"faction_id"`
}

type WhDatasheetModel struct {
	DatasheetID string `json:"datasheet_id"`
	Name        string `json:"name"`
	M           string `json:"m"`
	T           string `json:"t"`
	SV          string `json:"sv"`
	W           string `json:"w"`
	LD          string `json:"ld"`
	OC          string `json:"oc"`
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
	Name        string `json:"name"`
	Description string `json:"description"`
}

type WhSyncResult struct {
	Changed        bool     `json:"changed"`
	UpdatedSources []string `json:"updated_sources,omitempty"`
}
