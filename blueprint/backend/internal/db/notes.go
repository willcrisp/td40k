package db

import (
	"context"

	"github.com/willcrisp/blueprint/internal/models"
)

func ListNotes() ([]models.Note, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT n.id, n.user_id, u.username, n.content, n.created_at
		FROM notes n
		JOIN users u ON u.id = n.user_id
		ORDER BY n.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Username, &n.Content, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if notes == nil {
		notes = []models.Note{}
	}
	return notes, rows.Err()
}

func CreateNote(id, userID, content string) (*models.Note, error) {
	row := Pool.QueryRow(context.Background(), `
		INSERT INTO notes (id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, (SELECT username FROM users WHERE id = $2), content, created_at
	`, id, userID, content)

	var n models.Note
	if err := row.Scan(&n.ID, &n.UserID, &n.Username, &n.Content, &n.CreatedAt); err != nil {
		return nil, err
	}
	return &n, nil
}

// DeleteNote removes a note only if it belongs to userID.
func DeleteNote(id, userID string) error {
	tag, err := Pool.Exec(context.Background(),
		`DELETE FROM notes WHERE id = $1 AND user_id = $2`, id, userID)
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
