package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	mw "github.com/willcrisp/td40k/internal/middleware"
	"github.com/willcrisp/td40k/internal/db"
)

// PlaceUnitRequest is the request body for POST /api/rooms/{roomId}/units
type PlaceUnitRequest struct {
	DatasheetID   string  `json:"datasheet_id"`
	ModelName     string  `json:"model_name"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	ModelCount    int     `json:"model_count"`
	FacingDegrees int     `json:"facing_degrees"`
	NameOnBoard   *string `json:"name_on_board"`
	OwnerPlayerID *string `json:"owner_player_id"` // Only for GM
	FactionID     string  `json:"faction_id"`
}

// HandlePlaceUnit places a new unit on the board
func HandlePlaceUnit(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "roomId")

	// Verify room exists and player is in it
	room, err := db.GetRoom(roomID)
	if err != nil {
		jsonError(w, "room not found", http.StatusNotFound)
		return
	}

	// Check if player is in room
	isInRoom := (room.AttackerID != nil && *room.AttackerID == playerID) ||
		(room.DefenderID != nil && *room.DefenderID == playerID) ||
		room.GameMasterID == playerID

	if !isInRoom {
		jsonError(w, "not in room", http.StatusForbidden)
		return
	}

	var req PlaceUnitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.DatasheetID == "" || req.ModelName == "" {
		jsonError(w, "datasheet_id and model_name required", http.StatusBadRequest)
		return
	}
	if req.ModelCount < 1 {
		req.ModelCount = 1
	}

	// Determine unit owner
	ownerID := playerID
	if req.OwnerPlayerID != nil {
		// Only GM can set owner for other players
		if room.GameMasterID != playerID {
			jsonError(w, "only gm can set unit owner", http.StatusForbidden)
			return
		}
		ownerID = *req.OwnerPlayerID
	}

	// Create unit in database
	unit, err := db.CreateGameUnit(
		roomID,
		req.FactionID,
		req.DatasheetID,
		req.ModelName,
		req.X,
		req.Y,
		req.FacingDegrees,
		req.ModelCount,
		ownerID,
		playerID,
		req.NameOnBoard,
	)
	if err != nil {
		jsonError(w, "failed to create unit", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(unit)
}

// MoveUnitRequest is the request body for PATCH /api/rooms/{roomId}/units/{unitId}
type MoveUnitRequest struct {
	X             *float64 `json:"x"`
	Y             *float64 `json:"y"`
	FacingDegrees *int    `json:"facing_degrees"`
}

// HandleMoveUnit updates unit position and/or rotation
func HandleMoveUnit(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "roomId")
	unitID := chi.URLParam(r, "unitId")

	// Get current unit
	unit, err := db.GetGameUnit(unitID)
	if err != nil {
		jsonError(w, "unit not found", http.StatusNotFound)
		return
	}

	if unit.RoomID != roomID {
		jsonError(w, "unit not in this room", http.StatusBadRequest)
		return
	}

	var req MoveUnitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Use existing values if not provided
	x := unit.X
	if req.X != nil {
		x = *req.X
	}
	y := unit.Y
	if req.Y != nil {
		y = *req.Y
	}
	facing := unit.FacingDegrees
	if req.FacingDegrees != nil {
		facing = *req.FacingDegrees
	}

	// Update position
	updated, err := db.UpdateUnitPosition(roomID, unitID, x, y, facing, playerID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// WoundUnitRequest is the request body for POST /api/rooms/{roomId}/units/{unitId}/wounds
type WoundUnitRequest struct {
	Amount int `json:"amount"`
}

// HandleWoundUnit applies wounds to a unit
func HandleWoundUnit(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "roomId")
	unitID := chi.URLParam(r, "unitId")

	var req WoundUnitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount < 0 {
		jsonError(w, "amount must be non-negative", http.StatusBadRequest)
		return
	}

	unit, err := db.ApplyWoundsToUnit(roomID, unitID, req.Amount, playerID)
	if err != nil {
		code := http.StatusForbidden
		if err.Error() == "unit not found" {
			code = http.StatusNotFound
		}
		jsonError(w, err.Error(), code)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(unit)
}

// UpdateStatusRequest is the request body for POST /api/rooms/{roomId}/units/{unitId}/status
type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// HandleUpdateUnitStatus changes unit status (alive/in_reserves/dead)
func HandleUpdateUnitStatus(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "roomId")
	unitID := chi.URLParam(r, "unitId")

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	unit, err := db.UpdateUnitStatus(roomID, unitID, req.Status, playerID)
	if err != nil {
		code := http.StatusForbidden
		if err.Error() == "unit not found" {
			code = http.StatusNotFound
		}
		jsonError(w, err.Error(), code)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(unit)
}

// HandleDeleteUnit removes a unit from the board
func HandleDeleteUnit(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "roomId")
	unitID := chi.URLParam(r, "unitId")

	err := db.DeleteGameUnit(roomID, unitID, playerID)
	if err != nil {
		code := http.StatusForbidden
		if err.Error() == "unit not found" {
			code = http.StatusNotFound
		}
		jsonError(w, err.Error(), code)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// HandleGetRoomUnits retrieves all units in a room
func HandleGetRoomUnits(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")

	units, err := db.GetRoomUnits(roomID)
	if err != nil {
		jsonError(w, "failed to fetch units", http.StatusInternalServerError)
		return
	}

	if units == nil {
		units = []db.GameUnit{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(units)
}
