# Agent Instructions for `postgresTest` Project

Welcome, Agent. This is a full-stack real-time counter application utilizing PostgreSQL's `LISTEN/NOTIFY` mechanism and WebSockets to act as a single source of truth across all clients.

## Tech Stack Overview
- **Frontend:** Vue 3, Vite, PrimeVue 4, **Pinia** (state management), TypeScript. We use `bun` as the package manager instead of `npm`.
- **Backend:** Go 1.22, `go-chi` (router), `pgx` (Postgres driver), Gorilla WebSockets.
- **Database:** PostgreSQL 16
- **Tooling:** Docker & Docker Compose for orchestration. `just` for command running.

## Architectural Principles
- **Real-Time Reactivity:** The application relies on real-time state synchronization via WebSockets. The source of truth is always PostgreSQL. Modifications trigger `NOTIFY`, which the Go backend picks up via `LISTEN` and broadcasts to clients via WebSockets.
- **Single Responsibility Principle (SRP):** Keep components modularized and decoupled.
  - Frontend: Business logic lives in **Pinia stores** (`src/stores/`), not in components. `useWebSocketStore.ts` manages the singleton WebSocket connection; `useCounterStore.ts` owns counter state and actions. Components consume stores via `useCounterStore()` + `storeToRefs()`.
  - Backend: Use dedicated packages for state models, database handling, and websocket broadcasting instead of putting everything in `main.go`.

## Key Commands
You can automate operations using `just` instead of native docker commands.
- `just up`: Build containers from scratch and start them in the background (equivalent to `docker compose up --build -d`). 
- `just down`: Stop all containers.
- `just logs`: Follow the live logs for all services (Useful for observing WebSocket activity).
- `just reset`: Hard Reset, destroys containers, **wipes the database volume**, and rebuilds.

When making backend module changes (e.g. standard Go module issues), make sure you check module resolution with `go mod tidy` in the `backend/` directory. 

## Important Guidelines
- Do not add standard `REST` endpoints unless explicitly requested; our primary dynamic mechanism uses WebSocket broadcast triggered by Postgres events.
- Never write ad-hoc CSS for things PrimeVue provides.
- Ensure to include proper types in TypeScript, and maintain strict typing.
- Always check Docker Compose logs (`just logs`) when troubleshooting connectivity or database setup issues.
