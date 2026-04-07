# Agent Instructions for `td40k` Project

Welcome, Agent. This is a full-stack real-time gaming application utilizing PostgreSQL's `LISTEN/NOTIFY` mechanism and WebSockets to act as a single source of truth across all clients.

## Absolute Rules

These rules are non-negotiable. Do not deviate from them under any circumstances:

1. **Folder Structure:** Never deviate from the folder structure defined below. Use exact file paths specified in any document.

2. **PostgreSQL is Source of Truth:** No in-memory state is authoritative. All game logic is persisted before a response is sent.

3. **WebSocket Protocol:** WebSocket events are broadcast-only from the server. Clients do not send WebSocket messages. All mutations go through HTTP.

4. **Go Response Format:** All Go handlers must return JSON. Never return plain text errors except for HTTP 500 fallbacks.

5. **TypeScript Typing:** All Vue components are typed with TypeScript. No `any` types permitted.

6. **UI Library:** PrimeVue 4 is the only permitted component library. Do not install other UI libraries. Do not write ad-hoc CSS beyond what is strictly necessary for canvas layout.

7. **Code Formatting:** Prettier enforces formatting. Print width: 80. Use the config defined in Document 005.

## Technology Stack

| Technology | Version |
|-----------|---------|
| Go | 1.22 |
| PostgreSQL | 16 |
| Vue | 3 (Composition API only) |
| Vite | 5 |
| PrimeVue | 4 |
| Pinia | 2 |
| TypeScript | 5 |
| Bun | latest |
| go-chi | v5 |
| pgx | v5 |
| gorilla/websocket | v1 |

## Repository Layout

```
/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── db/
│   │   │   ├── db.go              # pgx pool init
│   │   │   ├── players.go         # player queries
│   │   │   └── rooms.go           # room queries
│   │   ├── handlers/
│   │   │   ├── players.go
│   │   │   └── rooms.go
│   │   ├── middleware/
│   │   │   └── player_auth.go
│   │   ├── models/
│   │   │   └── models.go
│   │   ├── ws/
│   │   │   ├── hub.go
│   │   │   └── client.go
│   │   └── listen/
│   │       └── listener.go
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── main.ts
│   │   ├── App.vue
│   │   ├── router/
│   │   │   └── index.ts
│   │   ├── stores/
│   │   │   ├── usePlayerStore.ts
│   │   │   ├── useGameListStore.ts
│   │   │   ├── useRoomStore.ts
│   │   │   ├── useBoardStore.ts
│   │   │   └── useWebSocketStore.ts
│   │   ├── views/
│   │   │   ├── HomeView.vue
│   │   │   ├── LobbyView.vue
│   │   │   └── GameView.vue
│   │   ├── components/
│   │   │   ├── home/
│   │   │   │   ├── CreateGameModal.vue
│   │   │   │   ├── OwnedGameCard.vue
│   │   │   │   └── JoinedGameCard.vue
│   │   │   ├── lobby/
│   │   │   │   ├── RoleSelector.vue
│   │   │   │   └── LobbyStatus.vue
│   │   │   └── game/
│   │   │       ├── BoardCanvas.vue
│   │   │       ├── PhaseBar.vue
│   │   │       ├── PhaseController.vue
│   │   │       ├── RoundTracker.vue
│   │   │       └── GameHUD.vue
│   │   ├── composables/
│   │   │   └── useBoardControls.ts
│   │   ├── types/
│   │   │   └── index.ts
│   │   └── lib/
│   │       └── api.ts
│   ├── index.html
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── package.json
├── db/
│   └── migrations/
│       ├── 001_create_players.sql
│       ├── 002_create_rooms.sql
│       ├── 003_create_room_events.sql
│       └── 004_create_triggers.sql
├── docker-compose.yml
├── Justfile
└── .env.example
```

## Naming Conventions

### Go
- **Files:** snake_case.go
- **Exported types/functions:** PascalCase
- **Unexported:** camelCase
- **Database queries:** Live in `internal/db/`, named after the entity they query
- **Handlers:** Named `Handle{Action}{Resource}` e.g. `HandleCreateRoom`

### TypeScript / Vue
- **Store files:** `use{Name}Store.ts`
- **Component files:** `PascalCase.vue`
- **Types/interfaces:** PascalCase
- **Composables:** `use{Name}.ts`
- **API functions:** In `lib/api.ts`, named `api{Action}{Resource}` e.g. `apiCreateRoom`

## Environment Variables

Variables are defined in `.env` (local) and injected via Docker Compose. The `.env.example` must always be kept in sync.

| Variable | Used By | Example |
|----------|---------|---------|
| POSTGRES_DSN | Backend | postgres://user:pass@db:5432/w40k |
| PORT | Backend | 8080 |
| VITE_API_BASE_URL | Frontend | http://localhost:8080 |
| VITE_WS_BASE_URL | Frontend | ws://localhost:8080 |
| POSTGRES_USER | Docker | w40k |
| POSTGRES_PASSWORD | Docker | w40k |
| POSTGRES_DB | Docker | w40k |

## HTTP Conventions

- All API responses are `application/json`
- Success responses use 200 or 201
- All error responses follow this shape: `{ "error": "human readable message" }`
- The `X-Player-ID` header carries the client's UUID on every request. Middleware extracts it and attaches it to the request context.

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
