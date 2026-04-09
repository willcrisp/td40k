package db

import (
	"context"

	"github.com/willcrisp/blueprint/internal/models"
)

func CreatePlayer(id, username, passwordHash string) (*models.Player, error) {
	row := Pool.QueryRow(context.Background(), `
		INSERT INTO players (id, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, is_admin, created_at
	`, id, username, passwordHash)

	var p models.Player
	if err := row.Scan(&p.ID, &p.Username, &p.IsAdmin, &p.CreatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}

func GetPlayerByUsername(username string) (id, passwordHash string, isAdmin bool, err error) {
	row := Pool.QueryRow(context.Background(), `
		SELECT id, password_hash, is_admin FROM players WHERE username = $1
	`, username)
	err = row.Scan(&id, &passwordHash, &isAdmin)
	return
}
