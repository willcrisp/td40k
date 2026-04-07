package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub encapsulates all standard WebSocket connection logic (Upgrading, Mapping, State)
type Hub struct {
	// clients maps connected WebSockets and is protected by a mutex locking thread
	clients      map[*websocket.Conn]bool
	clientsMutex sync.Mutex
	upgrader     websocket.Upgrader
}

// NewHub constructs a clean Websocket Hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// BroadcastMessage iterates over active TCP sockets mapping back to user browsers,
// attempting to blind-fire raw JSON byte payloads to them.
func (h *Hub) BroadcastMessage(message []byte) {
	h.clientsMutex.Lock()
	defer h.clientsMutex.Unlock()

	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("[WSS] Send failed. Client likely crashed/closed. Evicting.")
			client.Close()
			delete(h.clients, client)
		}
	}
}

// HandleWebSocket serves as the standard HTTP-to-WebSocket upgrade pathway
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WSS] Upgrade failure: %v", err)
		return
	}

	h.clientsMutex.Lock()
	h.clients[conn] = true
	h.clientsMutex.Unlock()
	log.Printf("[WSS] New Client Connect. Total: %d", len(h.clients))

	// According to spec, the client application technically acts as a Heartbeat.
	// We MUST constantly attempt to 'read' the socket to detect if the user closed their tab.
	go func() {
		defer func() {
			h.clientsMutex.Lock()
			delete(h.clients, conn)
			h.clientsMutex.Unlock()
			conn.Close()
			log.Printf("[WSS] Client Disconnect. Total: %d", len(h.clients))
		}()

		// Infinite drop loop. We wait for ReadMessage to error out indicating tab close.
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}
