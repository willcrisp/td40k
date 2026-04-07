package counter

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CounterResponse serves as the core JSON schema boundary between the Database and the Browser.
// The `json` tags tell Go's JSON encoder exactly what keys to use when sending to the Vue app.
type CounterResponse struct {
	Name      string    `json:"name,omitempty"`
	Value     int       `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IncrementRequest represents the JSON body we expect from the Vue POST request.
type IncrementRequest struct {
	Name string `json:"name"`
}

// Repository is a common architectural pattern!
// It wraps a database connection pool, hiding all SQL syntax from the rest of the application.
// Handlers will call Repository functions instead of writing SQL directly.
type Repository struct {
	DB *pgxpool.Pool
}

// GetCounter retrieves the room's current state cleanly.
// It takes a context.Context to support request timeouts and cancellation.
func (r *Repository) GetCounter(ctx context.Context, name string) (CounterResponse, error) {
	var resp CounterResponse
	
	// QueryRow executes a single row SQL statement safely using parameterized inputs ($1)
	// to completely prevent SQL Injection attacks.
	err := r.DB.QueryRow(ctx,
		`SELECT name, value, updated_at FROM counter WHERE name = $1`, name,
	).Scan(&resp.Name, &resp.Value, &resp.UpdatedAt)

	if err != nil {
		// pgx.ErrNoRows means the query executed successfully but found absolutely nothing!
		if errors.Is(err, pgx.ErrNoRows) {
			// Because we support dynamic rooms on the fly, if it doesn't exist, we don't throw an error.
			// We gracefully fake a starting configuration of "0" and pass it to the frontend!
			return CounterResponse{Name: name, Value: 0, UpdatedAt: time.Now()}, nil
		}
		// If it's a real database error (like connection loss), return it
		return resp, err
	}
	return resp, nil
}

// IncrementCounter runs our complex Postgres UPSERT magic.
func (r *Repository) IncrementCounter(ctx context.Context, name string) (CounterResponse, error) {
	var resp CounterResponse
	
	// 'UPSERT' pattern: INSERT ... ON CONFLICT DO UPDATE
	// 1. It blindly attempts to INSERT a brand new row starting at 1.
	// 2. If a room with that `name` already exists, a UNIQUE constraint conflict occurs!
	// 3. PostgreSQL instantly catches the conflict and pivots to the DO UPDATE clause,
	//    incrementing the existing value safely and atomically.
	// 4. RETURNING allows us to grab the final updated state instantly in the same network call.
	err := r.DB.QueryRow(ctx,
		`INSERT INTO counter (name, value) VALUES ($1, 1)
		 ON CONFLICT (name) DO UPDATE SET value = counter.value + 1, updated_at = NOW()
		 RETURNING name, value, updated_at`, name,
	).Scan(&resp.Name, &resp.Value, &resp.UpdatedAt)

	return resp, err
}
