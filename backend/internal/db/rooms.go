package db

import (
	"context"
	"fmt"

	"github.com/willcrisp/td40k/internal/models"
)

func CreateRoom(r models.Room) error {
	_, err := Pool.Exec(context.Background(), `
		INSERT INTO rooms
		    (id, name, status, game_master_id, battle_round,
		     active_player, current_phase)
		VALUES ($1,$2,'lobby',$3,1,'attacker','command')
	`, r.ID, r.Name, r.GameMasterID)
	return err
}

func GetRoom(id string) (*models.Room, error) {
	row := Pool.QueryRow(context.Background(), `
		SELECT id, name, status, game_master_id, attacker_id, defender_id,
		       battle_round, active_player, current_phase,
		       winner, created_at, updated_at
		FROM rooms WHERE id = $1
	`, id)
	var r models.Room
	if err := row.Scan(
		&r.ID, &r.Name, &r.Status, &r.GameMasterID,
		&r.AttackerID, &r.DefenderID,
		&r.BattleRound, &r.ActivePlayer, &r.CurrentPhase,
		&r.Winner, &r.CreatedAt, &r.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &r, nil
}

func SetRoomAttacker(roomID, playerID string) error {
	res, err := Pool.Exec(context.Background(), `
		UPDATE rooms SET attacker_id = $1
		WHERE id = $2 AND attacker_id IS NULL
	`, playerID, roomID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("attacker slot already taken")
	}
	return nil
}

func SetRoomDefender(roomID, playerID string) error {
	res, err := Pool.Exec(context.Background(), `
		UPDATE rooms SET defender_id = $1
		WHERE id = $2 AND defender_id IS NULL
	`, playerID, roomID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("defender slot already taken")
	}
	return nil
}

func StartRoom(roomID string) error {
	_, err := Pool.Exec(context.Background(), `
		UPDATE rooms SET status = 'active'
		WHERE id = $1
		  AND attacker_id IS NOT NULL
		  AND defender_id IS NOT NULL
		  AND status = 'lobby'
	`, roomID)
	return err
}

func UpdateRoomPhase(
	roomID, phase, activePlayer string,
	battleRound int,
	winner *string,
	status string,
) error {
	_, err := Pool.Exec(context.Background(), `
		UPDATE rooms
		SET current_phase  = $1,
		    active_player  = $2,
		    battle_round   = $3,
		    winner         = $4,
		    status         = $5
		WHERE id = $6
	`, phase, activePlayer, battleRound, winner, status, roomID)
	return err
}

func CloseRoom(roomID string) error {
	_, err := Pool.Exec(context.Background(), `
		UPDATE rooms SET status = 'closed' WHERE id = $1
	`, roomID)
	return err
}

func LogEvent(
	roomID string,
	playerID *string,
	eventType string,
	payload []byte,
) error {
	_, err := Pool.Exec(context.Background(), `
		INSERT INTO room_events (room_id, player_id, event_type, payload)
		VALUES ($1, $2, $3, $4)
	`, roomID, playerID, eventType, payload)
	return err
}
