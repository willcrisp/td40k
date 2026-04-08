package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/willcrisp/td40k/internal/models"
)

// CreatePlayer inserts a new player with a server-generated UUID and returns
// the created player. Returns an error wrapping "username taken" if the
// username already exists.
func CreatePlayer(username, nickname, passwordHash string) (models.Player, error) {
	var p models.Player
	err := Pool.QueryRow(context.Background(), `
		INSERT INTO players (id, username, nickname, password_hash, is_admin)
		VALUES (gen_random_uuid()::text, $1, $2, $3, false)
		RETURNING id, username, nickname, is_admin, created_at, last_seen
	`, username, nickname, passwordHash).Scan(
		&p.ID, &p.Username, &p.Nickname, &p.IsAdmin, &p.CreatedAt, &p.LastSeen,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return p, fmt.Errorf("username taken")
		}
		return p, err
	}
	return p, nil
}

// GetPlayerByUsername fetches the player row and its bcrypt hash for login.
// Returns an error wrapping "not found" if no matching username exists.
func GetPlayerByUsername(username string) (models.Player, string, error) {
	var p models.Player
	var hash string
	err := Pool.QueryRow(context.Background(), `
		SELECT id, username, nickname, is_admin, password_hash, created_at, last_seen
		FROM players
		WHERE username = $1
	`, username).Scan(
		&p.ID, &p.Username, &p.Nickname, &p.IsAdmin, &hash, &p.CreatedAt, &p.LastSeen,
	)
	if err != nil {
		return p, "", fmt.Errorf("not found")
	}
	return p, hash, nil
}

// GetPlayerByID fetches a player by UUID.
func GetPlayerByID(playerID string) (models.Player, error) {
	var p models.Player
	err := Pool.QueryRow(context.Background(), `
		SELECT id, username, nickname, is_admin, created_at, last_seen
		FROM players
		WHERE id = $1
	`, playerID).Scan(
		&p.ID, &p.Username, &p.Nickname, &p.IsAdmin, &p.CreatedAt, &p.LastSeen,
	)
	if err != nil {
		return p, fmt.Errorf("not found")
	}
	return p, nil
}

func GetPlayerGames(playerID string) (
	[]models.OwnedGameSummary,
	[]models.JoinedGameSummary,
	error,
) {
	ownedRows, err := Pool.Query(context.Background(), `
		SELECT id, name, status, battle_round, active_player,
		       current_phase, attacker_id, defender_id, created_at
		FROM rooms
		WHERE game_master_id = $1
		  AND status != 'closed'
		ORDER BY created_at DESC
	`, playerID)
	if err != nil {
		return nil, nil, err
	}
	defer ownedRows.Close()

	var owned []models.OwnedGameSummary
	for ownedRows.Next() {
		var g models.OwnedGameSummary
		if err := ownedRows.Scan(
			&g.ID, &g.Name, &g.Status, &g.BattleRound,
			&g.ActivePlayer, &g.CurrentPhase,
			&g.AttackerID, &g.DefenderID, &g.CreatedAt,
		); err != nil {
			return nil, nil, err
		}
		owned = append(owned, g)
	}

	joinedRows, err := Pool.Query(context.Background(), `
		SELECT
		    r.id, r.name, r.status, r.battle_round,
		    r.current_phase, r.created_at,
		    CASE
		        WHEN r.attacker_id = $1 THEN 'attacker'
		        WHEN r.defender_id = $1 THEN 'defender'
		    END AS role
		FROM rooms r
		WHERE (r.attacker_id = $1 OR r.defender_id = $1)
		  AND r.game_master_id != $1
		  AND r.status != 'closed'
		ORDER BY r.created_at DESC
	`, playerID)
	if err != nil {
		return nil, nil, err
	}
	defer joinedRows.Close()

	var joined []models.JoinedGameSummary
	for joinedRows.Next() {
		var g models.JoinedGameSummary
		if err := joinedRows.Scan(
			&g.ID, &g.Name, &g.Status, &g.BattleRound,
			&g.CurrentPhase, &g.CreatedAt, &g.Role,
		); err != nil {
			return nil, nil, err
		}
		joined = append(joined, g)
	}

	return owned, joined, nil
}

// isUniqueViolation returns true for PostgreSQL unique constraint errors (23505).
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
