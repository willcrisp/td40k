package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// GameUnit represents a unit deployed on the board
type GameUnit struct {
	ID            string `json:"id"`
	RoomID        string `json:"room_id"`
	FactionID     string `json:"faction_id"`
	DatasheetID   string `json:"datasheet_id"`
	ModelName     string `json:"model_name"`
	NameOnBoard   *string `json:"name_on_board"`
	ModelCount    int    `json:"model_count"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	FacingDegrees int    `json:"facing_degrees"`
	Status        string `json:"status"`
	Wounds        int    `json:"wounds"`
	OwnerPlayerID string `json:"owner_player_id"`
	CreatedBy     string `json:"created_by"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CreateGameUnit inserts a new unit on the board and logs the event
func CreateGameUnit(
	roomID, factionID, datasheetID, modelName string,
	x, y float64,
	facingDegrees int,
	modelCount int,
	ownerPlayerID, createdBy string,
	nameOnBoard *string,
) (*GameUnit, error) {
	unitID := uuid.New().String()

	_, err := Pool.Exec(context.Background(), `
		INSERT INTO game_units
			(id, room_id, faction_id, datasheet_id, model_name, name_on_board,
			 model_count, x, y, facing_degrees, owner_player_id, created_by,
			 status, wounds)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			'alive', 0)
	`, unitID, roomID, factionID, datasheetID, modelName, nameOnBoard,
		modelCount, x, y, facingDegrees, ownerPlayerID, createdBy)
	if err != nil {
		return nil, err
	}

	// Log the event
	payload := json.RawMessage(fmt.Sprintf(
		`{"unit_id":"%s","x":%f,"y":%f,"facing_degrees":%d,"datasheet_id":"%s"}`,
		unitID, x, y, facingDegrees, datasheetID))
	_ = LogEvent(roomID, &createdBy, "unit_placed", payload)

	return GetGameUnit(unitID)
}

// GetGameUnit retrieves a single unit by ID
func GetGameUnit(unitID string) (*GameUnit, error) {
	row := Pool.QueryRow(context.Background(), `
		SELECT id, room_id, faction_id, datasheet_id, model_name, name_on_board,
		       model_count, x, y, facing_degrees, status, wounds,
		       owner_player_id, created_by, created_at, updated_at
		FROM game_units WHERE id = $1
	`, unitID)

	var u GameUnit
	if err := row.Scan(
		&u.ID, &u.RoomID, &u.FactionID, &u.DatasheetID, &u.ModelName,
		&u.NameOnBoard, &u.ModelCount, &u.X, &u.Y, &u.FacingDegrees,
		&u.Status, &u.Wounds, &u.OwnerPlayerID, &u.CreatedBy,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetRoomUnits retrieves all units in a room
func GetRoomUnits(roomID string) ([]GameUnit, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT id, room_id, faction_id, datasheet_id, model_name, name_on_board,
		       model_count, x, y, facing_degrees, status, wounds,
		       owner_player_id, created_by, created_at, updated_at
		FROM game_units WHERE room_id = $1 ORDER BY created_at
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []GameUnit
	for rows.Next() {
		var u GameUnit
		if err := rows.Scan(
			&u.ID, &u.RoomID, &u.FactionID, &u.DatasheetID, &u.ModelName,
			&u.NameOnBoard, &u.ModelCount, &u.X, &u.Y, &u.FacingDegrees,
			&u.Status, &u.Wounds, &u.OwnerPlayerID, &u.CreatedBy,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

// GetPlayerUnits retrieves all units owned by a player in a room
func GetPlayerUnits(roomID, playerID string) ([]GameUnit, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT id, room_id, faction_id, datasheet_id, model_name, name_on_board,
		       model_count, x, y, facing_degrees, status, wounds,
		       owner_player_id, created_by, created_at, updated_at
		FROM game_units
		WHERE room_id = $1 AND owner_player_id = $2
		ORDER BY created_at
	`, roomID, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []GameUnit
	for rows.Next() {
		var u GameUnit
		if err := rows.Scan(
			&u.ID, &u.RoomID, &u.FactionID, &u.DatasheetID, &u.ModelName,
			&u.NameOnBoard, &u.ModelCount, &u.X, &u.Y, &u.FacingDegrees,
			&u.Status, &u.Wounds, &u.OwnerPlayerID, &u.CreatedBy,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

// UpdateUnitPosition moves a unit to a new x/y and updates facing
func UpdateUnitPosition(
	roomID, unitID string,
	x, y float64,
	facingDegrees int,
	playerID string,
) (*GameUnit, error) {
	// Verify unit exists and player owns it or is GM
	unit, err := GetGameUnit(unitID)
	if err != nil {
		return nil, err
	}
	room, err := GetRoom(roomID)
	if err != nil {
		return nil, err
	}
	if unit.OwnerPlayerID != playerID && room.GameMasterID != playerID {
		return nil, fmt.Errorf("unauthorized: not unit owner or game master")
	}

	_, err = Pool.Exec(context.Background(), `
		UPDATE game_units
		SET x = $1, y = $2, facing_degrees = $3
		WHERE id = $4
	`, x, y, facingDegrees, unitID)
	if err != nil {
		return nil, err
	}

	// Log the event
	payload := json.RawMessage(fmt.Sprintf(
		`{"unit_id":"%s","x":%f,"y":%f,"facing_degrees":%d}`,
		unitID, x, y, facingDegrees))
	_ = LogEvent(roomID, &playerID, "unit_moved", payload)

	return GetGameUnit(unitID)
}

// UpdateUnitStatus changes a unit's status (alive/in_reserves/dead)
func UpdateUnitStatus(
	roomID, unitID, status string,
	playerID string,
) (*GameUnit, error) {
	// Verify unit exists and player owns it or is GM
	unit, err := GetGameUnit(unitID)
	if err != nil {
		return nil, err
	}
	room, err := GetRoom(roomID)
	if err != nil {
		return nil, err
	}
	if unit.OwnerPlayerID != playerID && room.GameMasterID != playerID {
		return nil, fmt.Errorf("unauthorized: not unit owner or game master")
	}

	// Validate status
	if status != "alive" && status != "in_reserves" && status != "dead" {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	_, err = Pool.Exec(context.Background(), `
		UPDATE game_units SET status = $1 WHERE id = $2
	`, status, unitID)
	if err != nil {
		return nil, err
	}

	// Log the event
	payload := json.RawMessage(fmt.Sprintf(
		`{"unit_id":"%s","status":"%s"}`, unitID, status))
	_ = LogEvent(roomID, &playerID, "unit_status_changed", payload)

	return GetGameUnit(unitID)
}

// ApplyWoundsToUnit increments the wound counter
func ApplyWoundsToUnit(
	roomID, unitID string,
	amount int,
	playerID string,
) (*GameUnit, error) {
	// Verify unit exists and player owns it or is GM
	unit, err := GetGameUnit(unitID)
	if err != nil {
		return nil, err
	}
	room, err := GetRoom(roomID)
	if err != nil {
		return nil, err
	}
	if unit.OwnerPlayerID != playerID && room.GameMasterID != playerID {
		return nil, fmt.Errorf("unauthorized: not unit owner or game master")
	}

	if amount < 0 {
		return nil, fmt.Errorf("wound amount must be non-negative")
	}

	newWounds := unit.Wounds + amount
	_, err = Pool.Exec(context.Background(), `
		UPDATE game_units SET wounds = $1 WHERE id = $2
	`, newWounds, unitID)
	if err != nil {
		return nil, err
	}

	// Log the event
	payload := json.RawMessage(fmt.Sprintf(
		`{"unit_id":"%s","wounds":%d}`, unitID, newWounds))
	_ = LogEvent(roomID, &playerID, "unit_wounded", payload)

	return GetGameUnit(unitID)
}

// DeleteGameUnit removes a unit from the board
func DeleteGameUnit(
	roomID, unitID string,
	playerID string,
) error {
	// Verify unit exists and player owns it or is GM
	unit, err := GetGameUnit(unitID)
	if err != nil {
		return err
	}
	room, err := GetRoom(roomID)
	if err != nil {
		return err
	}
	if unit.OwnerPlayerID != playerID && room.GameMasterID != playerID {
		return fmt.Errorf("unauthorized: not unit owner or game master")
	}

	_, err = Pool.Exec(context.Background(), `
		DELETE FROM game_units WHERE id = $1
	`, unitID)
	if err != nil {
		return err
	}

	// Log the event
	payload := json.RawMessage(fmt.Sprintf(`{"unit_id":"%s"}`, unitID))
	_ = LogEvent(roomID, &playerID, "unit_removed", payload)

	return nil
}
