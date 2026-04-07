package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/willcrisp/td40k/internal/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *models.RoomStatePayload
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *models.RoomStatePayload, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("[ws] client joined room=%s player=%s",
				client.RoomID, client.PlayerID)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case payload := <-h.broadcast:
			msg, err := json.Marshal(models.RoomStateEvent{
				Event:   "room_state",
				Payload: *payload,
			})
			if err != nil {
				continue
			}
			for client := range h.clients {
				if client.RoomID != payload.RoomID {
					continue
				}
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Broadcast sends a room state payload to all clients in the given room.
func (h *Hub) Broadcast(payload *models.RoomStatePayload) {
	h.broadcast <- payload
}

// ServeWS upgrades an HTTP connection to WebSocket and registers the client.
func ServeWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("room_id")
		playerID := r.URL.Query().Get("player_id")
		if roomID == "" || playerID == "" {
			http.Error(w, "room_id and player_id required", http.StatusBadRequest)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[ws] upgrade error: %v", err)
			return
		}
		client := &Client{
			hub:      hub,
			conn:     conn,
			send:     make(chan []byte, 256),
			RoomID:   roomID,
			PlayerID: playerID,
		}
		hub.register <- client
		go client.writePump()
		go client.readPump()
	}
}
