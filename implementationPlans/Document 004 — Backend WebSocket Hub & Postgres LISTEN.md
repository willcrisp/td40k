Document 004 — Backend: WebSocket Hub & Postgres LISTEN

Purpose


Implement the real-time pipeline: WebSocket client management, room-filtered broadcasting, and the Postgres LISTEN worker.


---

backend/internal/ws/client.go

	package ws
	
	import (
	    "log"
	    "time"
	
	    "github.com/gorilla/websocket"
	)
	
	const (
	    writeWait  = 10 * time.Second
	    pongWait   = 60 * time.Second
	    pingPeriod = (pongWait * 9) / 10
	    maxMsgSize = 512
	)
	
	type Client struct {
	    hub      *Hub
	    conn     *websocket.Conn
	    send     chan []byte
	    RoomID   string
	    PlayerID string
	}
	
	func (c *Client) writePump() {
	    ticker := time.NewTicker(pingPeriod)
	    defer func() {
	        ticker.Stop()
	        c.conn.Close()
	    }()
	    for {
	        select {
	        case msg, ok := <-c.send:
	            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	            if !ok {
	                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	                return
	            }
	            if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
	                return
	            }
	        case <-ticker.C:
	            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
	                return
	            }
	        }
	    }
	}
	
	func (c *Client) readPump() {
	    defer func() {
	        c.hub.unregister <- c
	        c.conn.Close()
	    }()
	    c.conn.SetReadLimit(maxMsgSize)
	    c.conn.SetReadDeadline(time.Now().Add(pongWait))
	    c.conn.SetPongHandler(func(string) error {
	        c.conn.SetReadDeadline(time.Now().Add(pongWait))
	        return nil
	    })
	    for {
	        // We discard all inbound messages — clients do not send WS messages
	        if _, _, err := c.conn.ReadMessage(); err != nil {
	            if websocket.IsUnexpectedCloseError(
	                err,
	                websocket.CloseGoingAway,
	                websocket.CloseAbnormalClosure,
	            ) {
	                log.Printf("[ws] unexpected close: %v", err)
	            }
	            break
	        }
	    }
	}


---

backend/internal/ws/hub.go

	package ws
	
	import (
	    "encoding/json"
	    "log"
	    "net/http"
	
	    "github.com/gorilla/websocket"
	    "github.com/yourorg/w40k/internal/models"
	)
	
	var upgrader = websocket.Upgrader{
	    CheckOrigin: func(r *http.Request) bool { return true },
	}
	
	type Hub struct {
	    clients    map[*Client]bool
	    broadcast  chan *models.RoomStatePaylod
	    register   chan *Client
	    unregister chan *Client
	}
	
	func NewHub() *Hub {
	    return &Hub{
	        clients:    make(map[*Client]bool),
	        broadcast:  make(chan *models.RoomStatePaylod, 256),
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
	
	// Broadcast sends a payload to all clients in the given room.
	func (h *Hub) Broadcast(payload *models.RoomStatePaylod) {
	    h.broadcast <- payload
	}
	
	// ServeWS upgrades an HTTP connection to WebSocket and registers the client.
	func ServeWS(hub *Hub) http.HandlerFunc {
	    return func(w http.ResponseWriter, r *http.Request) {
	        roomID   := r.URL.Query().Get("room_id")
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


---

backend/internal/listen/listener.go

	package listen
	
	import (
	    "context"
	    "encoding/json"
	    "log"
	    "time"
	
	    "github.com/jackc/pgx/v5"
	    "github.com/yourorg/w40k/internal/models"
	    "github.com/yourorg/w40k/internal/ws"
	)
	
	func StartListener(dsn string, hub *ws.Hub) {
	    for {
	        if err := listen(dsn, hub); err != nil {
	            log.Printf("[listen] error: %v — reconnecting in 5s", err)
	            time.Sleep(5 * time.Second)
	        }
	    }
	}
	
	func listen(dsn string, hub *ws.Hub) error {
	    conn, err := pgx.Connect(context.Background(), dsn)
	    if err != nil {
	        return err
	    }
	    defer conn.Close(context.Background())
	
	    if _, err := conn.Exec(
	        context.Background(), "LISTEN room_updates",
	    ); err != nil {
	        return err
	    }
	    log.Println("[listen] listening on channel: room_updates")
	
	    for {
	        notification, err := conn.WaitForNotification(context.Background())
	        if err != nil {
	            return err
	        }
	        var payload models.RoomStatePaylod
	        if err := json.Unmarshal(
	            []byte(notification.Payload), &payload,
	        ); err != nil {
	            log.Printf("[listen] unmarshal error: %v", err)
	            continue
	        }
	        hub.Broadcast(&payload)
	    }
	}