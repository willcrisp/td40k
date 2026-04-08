package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/willcrisp/td40k/internal/db"
	"github.com/willcrisp/td40k/internal/handlers"
	mw "github.com/willcrisp/td40k/internal/middleware"
	"github.com/willcrisp/td40k/internal/listen"
	"github.com/willcrisp/td40k/internal/ws"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN not set")
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET not set")
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

	// WebSocket (no auth — player identity irrelevant for broadcast)
	r.Get("/ws", ws.ServeWS(hub))

	// Public auth endpoints
	r.Post("/api/auth/register", handlers.HandleRegister)
	r.Post("/api/auth/login", handlers.HandleLogin)

	// Protected — requires valid JWT
	r.Group(func(r chi.Router) {
		r.Use(mw.RequireAuth(jwtSecret))
		r.Get("/api/players/{id}/games", handlers.HandleGetPlayerGames)
		r.Post("/api/rooms", handlers.HandleCreateRoom)
		r.Get("/api/rooms/{id}", handlers.HandleGetRoom)
		r.Post("/api/rooms/{id}/join", handlers.HandleJoinRoom)
		r.Post("/api/rooms/{id}/start", handlers.HandleStartGame)
		r.Post("/api/rooms/{id}/phase/next", handlers.HandlePhaseNext)
		r.Post("/api/rooms/{id}/phase/prev", handlers.HandlePhasePrev)
		r.Post("/api/rooms/{id}/close", handlers.HandleCloseRoom)
		r.Post("/api/wahapedia/sync", handlers.HandleSyncWahapedia)
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
