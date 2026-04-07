package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/willcrisp/td40k/internal/db"
	mw "github.com/willcrisp/td40k/internal/middleware"
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

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(corsMiddleware)

	// Protected routes — player auth required
	r.Group(func(r chi.Router) {
		r.Use(mw.RequirePlayerID)
		// handlers registered in Document 003
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Player-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
