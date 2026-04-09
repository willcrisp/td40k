package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/willcrisp/blueprint/internal/db"
	"github.com/willcrisp/blueprint/internal/handlers"
	mw "github.com/willcrisp/blueprint/internal/middleware"
	"github.com/willcrisp/blueprint/internal/listen"
	"github.com/willcrisp/blueprint/internal/ws"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "POSTGRES_DSN is required")
		os.Exit(1)
	}
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	if err := db.Init(dsn); err != nil {
		fmt.Fprintf(os.Stderr, "db init: %v\n", err)
		os.Exit(1)
	}
	if err := db.RunMigrations(); err != nil {
		fmt.Fprintf(os.Stderr, "migrations: %v\n", err)
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
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Public
	r.Get("/ws", hub.ServeWS)
	r.Post("/api/auth/register", handlers.HandleRegister)
	r.Post("/api/auth/login", handlers.HandleLogin)

	// Protected
	r.Group(func(r chi.Router) {
		r.Use(mw.RequireAuth(jwtSecret))
		r.Get("/api/counter", handlers.HandleGetCounter)
		r.Post("/api/counter/increment", handlers.HandleIncrementCounter)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("listening on :%s\n", port)
	http.ListenAndServe(":"+port, r)
}
