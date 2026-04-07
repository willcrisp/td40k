package db

import (
	"context"

	"github.com/willcrisp/td40k/internal/models"
)

func UpsertPlayer(p models.Player) error {
	_, err := Pool.Exec(context.Background(), `
		INSERT INTO players (id, nickname, last_seen)
		VALUES ($1, $2, NOW())
		ON CONFLICT (id) DO UPDATE
			SET nickname  = EXCLUDED.nickname,
			    last_seen = NOW()
	`, p.ID, p.Nickname)
	return err
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
