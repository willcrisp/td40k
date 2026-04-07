Document 002 — Backend: Project Structure & Shared Middleware

Purpose


Bootstrap the Go module, define all shared models, initialize the database pool, and implement middleware.


---

backend/go.mod

	module github.com/yourorg/w40k
	
	go 1.22
	
	require (
	    github.com/go-chi/chi/v5 v5.0.12
	    github.com/gorilla/websocket v1.5.1
	    github.com/jackc/pgx/v5 v5.5.5
	    github.com/google/uuid v1.6.0
	)

Run after creating this file:


	cd backend && go mod tidy


---

backend/internal/models/models.go

	package models
	
	import "time"
	
	type Player struct {
	    ID        string    `json:"id"`
	    Nickname  string    `json:"nickname"`
	    CreatedAt time.Time `json:"created_at"`
	    LastSeen  time.Time `json:"last_seen"`
	}
	
	type Room struct {
	    ID            string    `json:"id"`
	    Name          string    `json:"name"`
	    Status        string    `json:"status"`
	    GameMasterID  string    `json:"game_master_id"`
	    AttackerID    *string   `json:"attacker_id"`
	    DefenderID    *string   `json:"defender_id"`
	    BattleRound   int       `json:"battle_round"`
	    ActivePlayer  string    `json:"active_player"`
	    CurrentPhase  string    `json:"current_phase"`
	    Winner        *string   `json:"winner"`
	    CreatedAt     time.Time `json:"created_at"`
	    UpdatedAt     time.Time `json:"updated_at"`
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
	    Payload RoomStatePaylod `json:"payload"`
	}
	
	type RoomStatePaylod struct {
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


---

backend/internal/db/db.go

	package db
	
	import (
	    "context"
	    "fmt"
	    "log"
	    "os"
	
	    "github.com/jackc/pgx/v5/pgxpool"
	)
	
	var Pool *pgxpool.Pool
	
	func Init(dsn string) error {
	    var err error
	    Pool, err = pgxpool.New(context.Background(), dsn)
	    if err != nil {
	        return fmt.Errorf("db init: %w", err)
	    }
	    if err := Pool.Ping(context.Background()); err != nil {
	        return fmt.Errorf("db ping: %w", err)
	    }
	    log.Println("[db] connected to postgres")
	    return nil
	}
	
	func RunMigrations() error {
	    migrations := []string{
	        "db/migrations/001_create_players.sql",
	        "db/migrations/002_create_rooms.sql",
	        "db/migrations/003_create_room_events.sql",
	        "db/migrations/004_create_triggers.sql",
	    }
	    for _, path := range migrations {
	        sql, err := os.ReadFile(path)
	        if err != nil {
	            return fmt.Errorf("read migration %s: %w", path, err)
	        }
	        if _, err := Pool.Exec(
	            context.Background(), string(sql),
	        ); err != nil {
	            return fmt.Errorf("run migration %s: %w", path, err)
	        }
	        log.Printf("[db] applied: %s", path)
	    }
	    return nil
	}


---

backend/internal/middleware/player_auth.go

	package middleware
	
	import (
	    "context"
	    "net/http"
	)
	
	type contextKey string
	
	const PlayerIDKey contextKey = "player_id"
	
	// RequirePlayerID extracts X-Player-ID header and injects it into context.
	// Returns 401 if the header is missing or empty.
	func RequirePlayerID(next http.Handler) http.Handler {
	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	        pid := r.Header.Get("X-Player-ID")
	        if pid == "" {
	            http.Error(w, `{"error":"missing X-Player-ID header"}`, http.StatusUnauthorized)
	            return
	        }
	        ctx := context.WithValue(r.Context(), PlayerIDKey, pid)
	        next.ServeHTTP(w, r.WithContext(ctx))
	    })
	}
	
	func GetPlayerID(r *http.Request) string {
	    v, _ := r.Context().Value(PlayerIDKey).(string)
	    return v
	}


---

backend/cmd/server/main.go

	package main
	
	import (
	    "log"
	    "net/http"
	    "os"
	
	    "github.com/go-chi/chi/v5"
	    chiMiddleware "github.com/go-chi/chi/v5/middleware"
	    "github.com/yourorg/w40k/internal/db"
	    "github.com/yourorg/w40k/internal/handlers"
	    mw "github.com/yourorg/w40k/internal/middleware"
	    "github.com/yourorg/w40k/internal/ws"
	    "github.com/yourorg/w40k/internal/listen"
	)
	
	func main() {
	    dsn := os.Getenv("POSTGRES_DSN")
	    if dsn == "" {
	        log.Fatal("POSTGRES_DSN not set")
	    }
	
	    if err := db.Init(dsn); err != nil {
	        log.Fatalf("db init: %v", err)
	    }
	    if err := db.RunMigrations(); err != nil {
	        log.Fatalf("migrations: %v", err)
	    }
	
	    hub := ws.NewHub()
	    go hub.Run()
	    go listen.StartListener(dsn, hub)
	
	    r := chi.NewRouter()
	    r.Use(chiMiddleware.Logger)
	    r.Use(chiMiddleware.Recoverer)
	    r.Use(corsMiddleware)
	
	    // WebSocket (no player auth middleware — player_id via query param)
	    r.Get("/ws", ws.ServeWS(hub))
	
	    // Public
	    r.Post("/api/players", handlers.HandleUpsertPlayer)
	
	    // Protected
	    r.Group(func(r chi.Router) {
	        r.Use(mw.RequirePlayerID)
	        r.Get("/api/players/{id}/games", handlers.HandleGetPlayerGames)
	        r.Post("/api/rooms", handlers.HandleCreateRoom)
	        r.Get("/api/rooms/{id}", handlers.HandleGetRoom)
	        r.Post("/api/rooms/{id}/join", handlers.HandleJoinRoom)
	        r.Post("/api/rooms/{id}/start", handlers.HandleStartGame)
	        r.Post("/api/rooms/{id}/phase/next", handlers.HandlePhaseNext)
	        r.Post("/api/rooms/{id}/phase/prev", handlers.HandlePhasePrev)
	        r.Post("/api/rooms/{id}/close", handlers.HandleCloseRoom)
	    })
	
	    port := os.Getenv("PORT")
	    if port == "" {
	        port = "8080"
	    }
	    log.Printf("[server] listening on :%s", port)
	    log.Fatal(http.ListenAndServe(":"+port, r))
	}
	
	func corsMiddleware(next http.Handler) http.Handler {
	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	        w.Header().Set("Access-Control-Allow-Origin", "*")
	        w.Header().Set("Access-Control-Allow-Headers",
	            "Content-Type, X-Player-ID")
	        w.Header().Set("Access-Control-Allow-Methods",
	            "GET, POST, OPTIONS")
	        if r.Method == http.MethodOptions {
	            w.WriteHeader(http.StatusNoContent)
	            return
	        }
	        next.ServeHTTP(w, r)
	    })
	}