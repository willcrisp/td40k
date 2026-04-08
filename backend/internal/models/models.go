package models

import "time"

type Player struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

type Room struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	GameMasterID string    `json:"game_master_id"`
	AttackerID   *string   `json:"attacker_id"`
	DefenderID   *string   `json:"defender_id"`
	BattleRound  int       `json:"battle_round"`
	ActivePlayer string    `json:"active_player"`
	CurrentPhase string    `json:"current_phase"`
	Winner       *string   `json:"winner"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RoomEvent struct {
	ID         int64      `json:"id"`
	RoomID     string     `json:"room_id"`
	PlayerID   *string    `json:"player_id"`
	EventType  string     `json:"event_type"`
	Payload    []byte     `json:"payload"`
	OccurredAt time.Time  `json:"occurred_at"`
}

// OwnedGameSummary is returned in GET /api/players/:id/games
type OwnedGameSummary struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	BattleRound  int       `json:"battle_round"`
	ActivePlayer string    `json:"active_player"`
	CurrentPhase string    `json:"current_phase"`
	AttackerID   *string   `json:"attacker_id"`
	DefenderID   *string   `json:"defender_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// JoinedGameSummary is returned in GET /api/players/:id/games
type JoinedGameSummary struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	Role         string    `json:"role"`
	BattleRound  int       `json:"battle_round"`
	CurrentPhase string    `json:"current_phase"`
	CreatedAt    time.Time `json:"created_at"`
}

// RoomStateEvent is the WebSocket broadcast payload
type RoomStateEvent struct {
	Event   string          `json:"event"` // always "room_state"
	Payload RoomStatePayload `json:"payload"`
}

type RoomStatePayload struct {
	RoomID       string  `json:"room_id"`
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	BattleRound  int     `json:"battle_round"`
	ActivePlayer string  `json:"active_player"`
	CurrentPhase string  `json:"current_phase"`
	Winner       *string `json:"winner"`
	AttackerID   *string `json:"attacker_id"`
	DefenderID   *string `json:"defender_id"`
	GameMasterID string  `json:"game_master_id"`
}
