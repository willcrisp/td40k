package db

import (
	"context"

	"github.com/willcrisp/blueprint/internal/models"
)

func GetCounter() (*models.CounterState, error) {
	row := Pool.QueryRow(context.Background(), `SELECT value FROM counter WHERE id = 1`)
	var c models.CounterState
	if err := row.Scan(&c.Value); err != nil {
		return nil, err
	}
	return &c, nil
}

func IncrementCounter() (*models.CounterState, error) {
	row := Pool.QueryRow(context.Background(),
		`UPDATE counter SET value = value + 1 WHERE id = 1 RETURNING value`)
	var c models.CounterState
	if err := row.Scan(&c.Value); err != nil {
		return nil, err
	}
	return &c, nil
}
