package models

import "time"

type Player struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type CounterState struct {
	Value int `json:"value"`
}

type Note struct {
	ID        string    `json:"id"`
	PlayerID  string    `json:"player_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteEvent struct {
	Op        string    `json:"op"`
	ID        string    `json:"id"`
	PlayerID  string    `json:"player_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type WsMessage struct {
	Event   string `json:"event"`
	Payload any    `json:"payload"`
}
