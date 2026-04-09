package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/willcrisp/blueprint/internal/db"
	"github.com/willcrisp/blueprint/internal/handlers"
	"github.com/willcrisp/blueprint/internal/listen"
	mw "github.com/willcrisp/blueprint/internal/middleware"
	"github.com/willcrisp/blueprint/internal/ws"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		slog.Error("POSTGRES_DSN is required")
		os.Exit(1)
	}
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	if err := db.Init(dsn); err != nil {
		slog.Error("db init failed", "err", err)
		os.Exit(1)
	}
	if err := db.RunMigrations(); err != nil {
		slog.Error("migrations failed", "err", err)
		os.Exit(1)
	}

	hub := ws.NewHub()
	go hub.Run()
	go listen.StartListener(dsn, hub)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Public
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	r.Get("/ws", hub.ServeWS)
	r.Post("/api/auth/register", handlers.HandleRegister)
	r.Post("/api/auth/login", handlers.HandleLogin)

	// Protected
	r.Group(func(r chi.Router) {
		r.Use(mw.RequireAuth(jwtSecret))
		r.Get("/api/counter", handlers.HandleGetCounter)
		r.Post("/api/counter/increment", handlers.HandleIncrementCounter)
		r.Get("/api/notes", handlers.HandleListNotes)
		r.Post("/api/notes", handlers.HandleCreateNote)
		r.Delete("/api/notes/{id}", handlers.HandleDeleteNote)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("server listening", "port", port)
	http.ListenAndServe(":"+port, r)
}
