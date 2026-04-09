package models

import "time"

type Player struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type CounterState struct {
	Value int `json:"value"`
}

type WsMessage struct {
	Event   string `json:"event"`
	Payload any    `json:"payload"`
}
