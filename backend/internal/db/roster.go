package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RosterEntry represents a unit in a player's imported army roster.
type RosterEntry struct {
	ID          string `json:"id"`
	RoomID      string `json:"room_id"`
	PlayerID    string `json:"player_id"`
	DatasheetID string `json:"datasheet_id"`
	ModelName   string `json:"model_name"`
	Quantity    int    `json:"quantity"`
	FactionID   string `json:"faction_id"`
	Points      int    `json:"points"`
	CreatedAt   string `json:"created_at"`
}

// BulkCreateRosterEntries replaces all roster entries for a player in a room
// with the provided entries, using a single transaction.
func BulkCreateRosterEntries(roomID, playerID string, entries []RosterEntry) ([]RosterEntry, error) {
	ctx := context.Background()
	tx, err := Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		`DELETE FROM unit_roster WHERE room_id = $1 AND player_id = $2`,
		roomID, playerID,
	); err != nil {
		return nil, fmt.Errorf("clear roster: %w", err)
	}

	var result []RosterEntry
	for _, e := range entries {
		e.ID = uuid.New().String()
		e.RoomID = roomID
		e.PlayerID = playerID
		var createdAt string
		if err := tx.QueryRow(ctx,
			`INSERT INTO unit_roster (id, room_id, player_id, datasheet_id, model_name, quantity, faction_id, points)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 RETURNING created_at`,
			e.ID, e.RoomID, e.PlayerID, e.DatasheetID, e.ModelName,
			e.Quantity, e.FactionID, e.Points,
		).Scan(&createdAt); err != nil {
			return nil, fmt.Errorf("insert roster entry %s: %w", e.ModelName, err)
		}
		e.CreatedAt = createdAt
		result = append(result, e)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return result, nil
}

// GetRoster returns all roster entries for a player in a room.
func GetRoster(roomID, playerID string) ([]RosterEntry, error) {
	rows, err := Pool.Query(context.Background(),
		`SELECT id, room_id, player_id, datasheet_id, model_name, quantity,
		        COALESCE(faction_id, ''), points, created_at
		 FROM unit_roster
		 WHERE room_id = $1 AND player_id = $2
		 ORDER BY created_at`,
		roomID, playerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []RosterEntry
	for rows.Next() {
		var e RosterEntry
		if err := rows.Scan(
			&e.ID, &e.RoomID, &e.PlayerID, &e.DatasheetID,
			&e.ModelName, &e.Quantity, &e.FactionID, &e.Points, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// ClearRoster removes all roster entries for a player in a room.
func ClearRoster(roomID, playerID string) error {
	_, err := Pool.Exec(context.Background(),
		`DELETE FROM unit_roster WHERE room_id = $1 AND player_id = $2`,
		roomID, playerID,
	)
	return err
}
