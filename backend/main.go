package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Requires standard library SQL bridge for Goose engine
	"github.com/pressly/goose/v3"

	"counterapp/counter"
	"counterapp/websocket"
)

func main() {
	// ==========================================
	// 1. BOOTSTRAP DATABASES
	// ==========================================
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL must be supplied")
	}

	var db *pgxpool.Pool
	var err error
	for i := 0; i < 10; i++ {
		db, err = pgxpool.New(context.Background(), dsn)
		if err == nil {
			if pingErr := db.Ping(context.Background()); pingErr == nil {
				break
			}
		}
		log.Printf("[System] Waiting for postgres... %d/10", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Fatal Database Initialization Error: %v", err)
	}
	defer db.Close()
	log.Println("[System] Database connected beautifully.")

	runMigrations(dsn)

	// ==========================================
	// 2. ORCHESTRATE DOMAIN MODULES
	// ==========================================
	
	// Create the Global Websocket communication structure
	hub := websocket.NewHub()

	// Spin up the Counter application logic
	counterModule := counter.NewModule(db, hub.BroadcastMessage)
	go counterModule.ListenForUpdates() // Deploy background triggers for this domain logic

	// ==========================================
	// 3. LAUNCH HTTP ENGINE
	// ==========================================
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
	}))

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Get("/api/ws", hub.HandleWebSocket)

	counterModule.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("[System] Node spinning up fully prepared on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server abort: %v", err)
	}
}

// runMigrations is explicitly instructed to search specific modules, retrieve their natively embedded 
// raw SQL blocks, and apply them using their own tracked versioning tables to fully abide by SRP!
func runMigrations(dsn string) {
	// Standard SQL connection needed for Goose Library.
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open standard SQL connection for Goose: %v", err)
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Goose failed setting dialect: %v", err)
	}

	log.Println("[Goose] Deploying Counter Module schemas...")
	
	// Ensure the counter domain uses an isolated version tracking table (goose isolates by standard name instead)
	goose.SetTableName("goose_db_version_counter")
	
	// Pass the natively Compiled embedded file-system containing all physical migrations specific to `counter/migrations/*`
	goose.SetBaseFS(counter.MigrationsFS)
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatalf("Goose fatal deploy error: %v", err)
	}
	log.Println("[Goose] Counter schemas completed.")
}
