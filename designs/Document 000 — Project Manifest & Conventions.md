Document 000 — Project Manifest & Conventions

Purpose


This document defines the project-wide conventions, folder layout, and rules every subsequent document must follow. The AI agent must read and internalize this document before executing any other.


---

Absolute Rules

1. Never deviate from the folder structure defined in this document. If a file path is specified in any document, use it exactly.

2. PostgreSQL is the single source of truth. No in-memory state is authoritative. All game logic is persisted before a response is sent.

3. WebSocket events are broadcast-only from the server. Clients do not send WebSocket messages. All mutations go through HTTP.

4. All Go handlers must return JSON. Never return plain text errors except for HTTP 500 fallbacks.

5. All Vue components are typed with TypeScript. No any types permitted.

6. PrimeVue 4 is the only permitted component library. Do not install other UI libraries. Do not write ad-hoc CSS beyond what is strictly necessary for canvas layout.

7. Prettier enforces formatting. Print width: 80. Use the config defined in 005.


---

Technology Versions

Technology	Version
Go	1.22
PostgreSQL	16
Vue	3 (Composition API only)
Vite	5
PrimeVue	4
Pinia	2
TypeScript	5
Bun	latest
go-chi	v5
pgx	v5
gorilla/websocket	v1

---

Repository Layout

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


---

Naming Conventions

Go

- Files: snake_case.go

- Exported types/functions: PascalCase

- Unexported: camelCase

- Database query functions live in internal/db/, named after the entity they query

- Handler functions named Handle{Action}{Resource} e.g. HandleCreateRoom

TypeScript / Vue

- Store files: use{Name}Store.ts

- Component files: PascalCase.vue

- Types/interfaces: PascalCase

- Composables: use{Name}.ts

- API functions in lib/api.ts: named api{Action}{Resource} e.g. apiCreateRoom


---

Environment Variables


Defined in .env (local) and injected via Docker Compose. The .env.example must always be kept in sync.


Variable	Used By	Example
POSTGRES_DSN	Backend	postgres://user:pass@db:5432/w40k
PORT	Backend	8080
VITE_API_BASE_URL	Frontend	http://localhost:8080
VITE_WS_BASE_URL	Frontend	ws://localhost:8080
POSTGRES_USER	Docker	w40k
POSTGRES_PASSWORD	Docker	w40k
POSTGRES_DB	Docker	w40k

---

HTTP Conventions

- All API responses are application/json

- Success responses use 200 or 201

- All error responses follow this shape:

	{ "error": "human readable message" }



- The X-Player-ID header carries the client's UUID on every request. Middleware extracts it and attaches it to the request context.