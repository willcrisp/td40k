package db

import (
	"context"

	"github.com/willcrisp/blueprint/internal/models"
)

func CreateUser(id, username, passwordHash string) (*models.User, error) {
	row := Pool.QueryRow(context.Background(), `
		INSERT INTO users (id, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, is_admin, created_at
	`, id, username, passwordHash)

	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.IsAdmin, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByUsername(username string) (id, passwordHash string, isAdmin bool, err error) {
	row := Pool.QueryRow(context.Background(), `
		SELECT id, password_hash, is_admin FROM users WHERE username = $1
	`, username)
	err = row.Scan(&id, &passwordHash, &isAdmin)
	return
}
