package db

import (
	"context"

	"github.com/willcrisp/blueprint/internal/models"
)

func ListNotes() ([]models.Note, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT n.id, n.player_id, p.username, n.content, n.created_at
		FROM notes n
		JOIN players p ON p.id = n.player_id
		ORDER BY n.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.PlayerID, &n.Username, &n.Content, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if notes == nil {
		notes = []models.Note{}
	}
	return notes, rows.Err()
}

func CreateNote(id, playerID, content string) (*models.Note, error) {
	row := Pool.QueryRow(context.Background(), `
		INSERT INTO notes (id, player_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, player_id, (SELECT username FROM players WHERE id = $2), content, created_at
	`, id, playerID, content)

	var n models.Note
	if err := row.Scan(&n.ID, &n.PlayerID, &n.Username, &n.Content, &n.CreatedAt); err != nil {
		return nil, err
	}
	return &n, nil
}

// DeleteNote removes a note only if it belongs to playerID.
func DeleteNote(id, playerID string) error {
	tag, err := Pool.Exec(context.Background(),
		`DELETE FROM notes WHERE id = $1 AND player_id = $2`, id, playerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errNotFound
	}
	return nil
}

var errNotFound = &notFoundError{}

type notFoundError struct{}

func (e *notFoundError) Error() string { return "not found" }

func IsNotFound(err error) bool {
	_, ok := err.(*notFoundError)
	return ok
}
