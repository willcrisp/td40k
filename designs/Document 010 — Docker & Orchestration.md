Document 010 — Docker & Orchestration

Purpose


Define all Docker, Compose, and Justfile configuration to run the full stack locally.


---

.env.example

	POSTGRES_USER=w40k
	POSTGRES_PASSWORD=w40k
	POSTGRES_DB=w40k
	POSTGRES_DSN=postgres://w40k:w40k@db:5432/w40k?sslmode=disable
	PORT=8080
	VITE_API_BASE_URL=http://localhost:8080
	VITE_WS_BASE_URL=ws://localhost:8080

Copy to .env before running:


	cp .env.example .env


---

docker-compose.yml

	services:
	  db:
	    image: postgres:16-alpine
	    restart: unless-stopped
	    environment:
	      POSTGRES_USER: ${POSTGRES_USER}
	      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
	      POSTGRES_DB: ${POSTGRES_DB}
	    volumes:
	      - pgdata:/var/lib/postgresql/data
	    ports:
	      - "5432:5432"
	    healthcheck:
	      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER}"]
	      interval: 5s
	      timeout: 5s
	      retries: 10
	
	  backend:
	    build:
	      context: ./backend
	      dockerfile: Dockerfile
	    restart: unless-stopped
	    environment:
	      POSTGRES_DSN: ${POSTGRES_DSN}
	      PORT: ${PORT}
	    ports:
	      - "8080:8080"
	    volumes:
	      # Mount db migrations so the Go binary can read them at runtime
	      - ./db:/db:ro
	    depends_on:
	      db:
	        condition: service_healthy
	
	  frontend:
	    build:
	      context: ./frontend
	      dockerfile: Dockerfile
	    restart: unless-stopped
	    environment:
	      VITE_API_BASE_URL: ${VITE_API_BASE_URL}
	      VITE_WS_BASE_URL: ${VITE_WS_BASE_URL}
	    ports:
	      - "5173:5173"
	    depends_on:
	      - backend
	
	volumes:
	  pgdata:


---

backend/Dockerfile

	FROM golang:1.22-alpine AS builder
	WORKDIR /app
	COPY go.mod go.sum ./
	RUN go mod download
	COPY . .
	RUN go build -o server ./cmd/server
	
	FROM alpine:3.19
	WORKDIR /app
	COPY --from=builder /app/server .
	EXPOSE 8080
	CMD ["./server"]


---

frontend/Dockerfile

	FROM oven/bun:1 AS base
	WORKDIR /app
	COPY package.json bun.lockb ./
	RUN bun install --frozen-lockfile
	COPY . .
	EXPOSE 5173
	CMD ["bun", "run", "dev", "--host"]


---

Justfile

	# Load .env
	set dotenv-load := true
	
	# ── Main commands ─────────────────────────────────────────────────────────────
	
	# Build and start all containers in background
	up:
	    docker compose up --build -d
	
	# Stop and remove containers
	down:
	    docker compose down
	
	# Stream all logs
	logs:
	    docker compose logs -f
	
	# Stream backend logs only
	logs-backend:
	    docker compose logs -f backend
	
	# Stream frontend logs only
	logs-frontend:
	    docker compose logs -f frontend
	
	# Hard reset: destroy containers, wipe db volume, rebuild
	reset:
	    docker compose down -v --remove-orphans
	    docker compose up --build -d
	
	# ── Dev helpers ───────────────────────────────────────────────────────────────
	
	# Open a psql shell in the db container
	db:
	    docker compose exec db psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}
	
	# Run go tests in backend
	test-backend:
	    cd backend && go test ./...
	
	# Format frontend
	fmt-frontend:
	    cd frontend && bun run format
	
	# Install frontend dependencies
	install-frontend:
	    cd frontend && bun install
	
	# Tail only WebSocket-related log lines
	logs-ws:
	    docker compose logs -f backend | grep --line-buffered '\[ws\]\|\[listen\]'