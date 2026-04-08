# CLAUDE.md — td40k

A full-stack real-time Warhammer 40K game management application. PostgreSQL is the single source of truth; all state changes persist before any response is sent, and all connected clients are updated in real-time via WebSockets driven by PostgreSQL `LISTEN/NOTIFY`.

---

## Absolute Rules

These are non-negotiable. Do not deviate under any circumstances.

1. **Folder structure is fixed.** Never create files outside the paths defined in the Repository Layout section.
2. **PostgreSQL is the only source of truth.** No in-memory state is authoritative. All game logic must persist to the database before a response is sent.
3. **WebSocket is broadcast-only from the server.** Clients never send WebSocket messages. All mutations go through HTTP REST endpoints.
4. **Go handlers always return JSON.** Never return plain text errors except for HTTP 500 fallbacks.
5. **No `any` types in TypeScript.** All Vue components and stores must be fully typed.
6. **PrimeVue 4 is the only UI library.** Do not install other component libraries. Do not write ad-hoc CSS beyond what is strictly necessary for canvas layout (PrimeVue handles everything else).
7. **Prettier enforces formatting.** Print width: 80, singleQuote, semi, trailingComma: es5.

---

## Technology Stack

| Layer      | Technology                       | Version  |
|------------|----------------------------------|----------|
| Backend    | Go                               | 1.22     |
| Backend    | chi router                       | v5       |
| Backend    | pgx (PostgreSQL driver)          | v5       |
| Backend    | gorilla/websocket                | v1       |
| Backend    | google/uuid                      | v1       |
| Backend    | golang-jwt/jwt                   | v5       |
| Backend    | golang.org/x/crypto/bcrypt      | latest   |
| Frontend   | Vue (Composition API only)       | 3        |
| Frontend   | Vite                             | 5        |
| Frontend   | TypeScript                       | 5        |
| Frontend   | Pinia                            | 3        |
| Frontend   | PrimeVue                         | 4        |
| Frontend   | Axios                            | latest   |
| Frontend   | Bun (package manager)            | latest   |
| Database   | PostgreSQL                       | 16       |
| Infra      | Docker Compose + nginx           | —        |

---

## Repository Layout

```
/
├── backend/
│   ├── cmd/server/main.go             # Entry point: routing, middleware, server start
│   ├── internal/
│   │   ├── db/
│   │   │   ├── db.go                  # pgxpool init & migration runner
│   │   │   ├── players.go             # Player CRUD queries
│   │   │   ├── rooms.go               # Room CRUD queries
│   │   │   └── wahapedia.go           # Warhammer datasheet/faction queries
│   │   ├── handlers/
│   │   │   ├── auth.go                # Register/Login with JWT token issuance
│   │   │   ├── players.go             # HTTP handlers for player endpoints
│   │   │   ├── rooms.go               # HTTP handlers for room/game lifecycle
│   │   │   └── wahapedia.go           # Warhammer CSV sync from Wahapedia.ru
│   │   ├── middleware/
│   │   │   └── player_auth.go         # JWT token verification & context injection
│   │   ├── models/
│   │   │   ├── models.go              # Domain structs: Player, Room, RoomEvent, RoomStatePayload
│   │   │   ├── unit.go                # Unit simulation: BaseUnit, UnitStats, BoardPosition
│   │   │   ├── wahapedia.go           # Warhammer models: WhFaction, WhDatasheet, WhWeapon, etc.
│   │   │   └── units/                 # Individual unit types (Space Marine, Rhino, etc.)
│   │   ├── ws/
│   │   │   ├── hub.go                 # Gorilla WebSocket hub: register/unregister/broadcast
│   │   │   └── client.go              # Individual WebSocket client pump goroutines
│   │   └── listen/
│   │       └── listener.go            # Goroutine: LISTEN on postgres channel → hub.Broadcast
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── main.ts                    # Bootstrap: Pinia, PrimeVue, Router
│   │   ├── App.vue                    # Root component
│   │   ├── router/index.ts            # Vue Router with navigation guards
│   │   ├── stores/
│   │   │   ├── usePlayerStore.ts      # Player UUID (localStorage-backed)
│   │   │   ├── useRoomStore.ts        # Active room state, roles, phase
│   │   │   ├── useGameListStore.ts    # Owned & joined game listings
│   │   │   ├── useBoardStore.ts       # Board/canvas unit state
│   │   │   └── useWebSocketStore.ts   # Singleton WebSocket with auto-reconnect
│   │   ├── views/
│   │   │   ├── HomeView.vue           # Game list, create/join
│   │   │   ├── LobbyView.vue          # Role selection (Attacker/Defender/GM)
│   │   │   └── GameView.vue           # Active game with board canvas
│   │   ├── components/
│   │   │   ├── home/
│   │   │   │   ├── CreateGameModal.vue
│   │   │   │   ├── OwnedGameCard.vue
│   │   │   │   └── JoinedGameCard.vue
│   │   │   ├── lobby/
│   │   │   │   ├── RoleSelector.vue
│   │   │   │   └── LobbyStatus.vue
│   │   │   └── game/
│   │   │       ├── BoardCanvas.vue    # HTML canvas rendering
│   │   │       ├── PhaseBar.vue       # Visual phase indicator
│   │   │       ├── PhaseController.vue
│   │   │       ├── RoundTracker.vue
│   │   │       └── GameHUD.vue
│   │   ├── composables/
│   │   │   └── useBoardControls.ts    # Canvas interaction logic
│   │   ├── types/index.ts             # All shared TypeScript types/interfaces
│   │   └── lib/api.ts                 # Axios client + all api{Action}{Resource} functions
│   ├── index.html
│   ├── vite.config.ts                 # Dev proxy to backend, @ alias
│   ├── tsconfig.json                  # Strict mode, ES2020
│   ├── .prettierrc
│   ├── package.json
│   ├── Dockerfile                     # Bun builder → nginx
│   └── Dockerfile.dev                 # Bun dev server (volume-mounted src)
├── db/migrations/
│   ├── 001_create_players.sql
│   ├── 002_create_rooms.sql
│   ├── 003_create_room_events.sql
│   └── 004_create_triggers.sql        # set_updated_at + notify_room_update via pg_notify
├── implementationPlans/               # Reference docs 000–010 + design system
├── docker-compose.yml
├── docker-compose.dev.yml             # Overlay: Vite HMR dev server
├── justfile                           # Task automation (see Key Commands)
├── agents.md                          # Original agent instructions
├── architecture.md                    # High-level design reference
└── README.md
```

---

## Architecture: How Data Flows

```
  Client (Vue)  ──── REST POST ────►  Go Handler
                                          │
                                          ▼
                                     PostgreSQL UPDATE
                                          │
                                    Trigger fires pg_notify()
                                          │
                                          ▼
                                     Go LISTEN goroutine
                                     (internal/listen/)
                                          │
                                          ▼
                                     WebSocket Hub
                                     (internal/ws/hub.go)
                                          │
                                    Broadcast JSON to all clients
                                          │
                                          ▼
  Client (Vue)  ◄── WebSocket push ──  Pinia store update → UI re-render
```

Key principle: **clients never push state, only pull it via WebSocket broadcasts triggered by their own (or others') HTTP mutations.**

---

## Authentication Flow

Players authenticate via username/password at `/api/auth/register` or `/api/auth/login`. Both endpoints return a JWT token with 7-day expiry. The token must be sent in the `Authorization: Bearer <token>` header on all protected endpoints.

**Token Payload:**
```json
{
  "player_id": "uuid",
  "exp": 1234567890,
  "iat": 1234567890
}
```

The middleware extracts the player ID and stores it in `r.Context()` for handlers to access via `context.Value(key)`.

---

## Naming Conventions

### Go
- **Files:** `snake_case.go`
- **Exported types/functions:** `PascalCase`
- **Unexported:** `camelCase`
- **Handlers:** `Handle{Action}{Resource}` — e.g. `HandleCreateRoom`, `HandlePhaseNext`
- **DB query files:** one file per entity in `internal/db/` — e.g. `rooms.go`, `players.go`

### TypeScript / Vue
- **Store files:** `use{Name}Store.ts`
- **Component files:** `PascalCase.vue`
- **Types/interfaces:** `PascalCase`
- **Composables:** `use{Name}.ts`
- **API functions:** `api{Action}{Resource}` in `lib/api.ts` — e.g. `apiCreateRoom`, `apiJoinRoom`

---

## HTTP Conventions

- All responses are `application/json`
- Success: `200 OK` or `201 Created`
- Errors: `{ "error": "human readable message" }` with appropriate 4xx/5xx status
- Authentication: JWT token in `Authorization: Bearer <token>` header on protected endpoints. Middleware (`player_auth.go`) verifies and injects player ID to `context.Context`.
- Public endpoints: `/api/auth/register`, `/api/auth/login`, `/ws` (WebSocket broadcast)

---

## Game Domain

**Roles:** Game Master (GM), Attacker, Defender

**Room status lifecycle:** `waiting` → `active` → `closed`

**Battle phases (in order, per round):**
1. Command
2. Movement
3. Shooting
4. Charge
5. Fight

Up to 5 battle rounds. The GM controls phase advancement (`PhaseNext` / `PhasePrev`). GM can also assign Attacker and Defender roles from the lobby.

---

## Warhammer Data (Wahapedia Integration)

The app syncs Warhammer 40K unit datasheets from **Wahapedia.ru** via CSV download:

- **Endpoint:** `POST /api/wahapedia/sync` (protected, JWT required)
- **Source:** `https://wahapedia.ru/wh40k10ed/data/` (pipe-delimited CSVs)
- **Datasets synced:**
  - Factions (armies)
  - Datasheets (unit profiles)
  - Datasheet Models (individual unit stats)
  - Datasheet Weapons (weapon profiles)
  - Datasheet Abilities (special rules)

**Change Detection:** SHA256 hashes are stored per CSV to avoid redundant syncs. Sync only proceeds if at least one hash differs.

**Data Storage:** All Warhammer data persists to PostgreSQL tables created by migrations. The sync is idempotent—re-running replaces old data with fresh CSV entries.

**Schema:** See `internal/models/wahapedia.go` for struct definitions (WhFaction, WhDatasheet, WhDatasheetModel, WhDatasheetWeapon, WhDatasheetAbility).

---

## Environment Variables

Defined in `.env` (gitignored). `.env.example` must stay in sync.

| Variable            | Used By  | Example                              |
|---------------------|----------|--------------------------------------|
| `POSTGRES_DSN`      | Backend  | `postgres://w40k:w40k@db:5432/w40k`  |
| `PORT`              | Backend  | `8080`                               |
| `JWT_SECRET`        | Backend  | Base64-encoded 32-byte key (run `openssl rand -base64 32`) |
| `VITE_API_BASE_URL` | Frontend | `http://localhost:8080`              |
| `VITE_WS_BASE_URL`  | Frontend | `ws://localhost:8080`                |
| `POSTGRES_USER`     | Docker   | `w40k`                               |
| `POSTGRES_PASSWORD` | Docker   | `w40k`                               |
| `POSTGRES_DB`       | Docker   | `w40k`                               |

---

## Key Commands

```bash
just up       # Build images + start all containers in background
just dev      # Start with HMR (Vite dev server on :5173, mounts src volume)
just down     # Stop all containers
just logs     # Follow logs for all services (useful for WebSocket debugging)
just reset    # DESTRUCTIVE: wipe DB volume, rebuild from scratch
```

Service URLs:

| Service        | URL                       |
|----------------|---------------------------|
| Frontend (prod) | http://localhost:3000    |
| Frontend (dev HMR) | http://localhost:5173 |
| Backend API    | http://localhost:8080/api |
| PostgreSQL     | localhost:5432            |

---

## Development Workflows

### Backend changes
```bash
cd backend
go mod tidy          # After adding/removing dependencies
# Rebuild via Docker:
just up
```

### Frontend changes (with hot reload)
```bash
just dev             # Vite HMR on port 5173
```

### Database migrations
SQL files in `db/migrations/` are applied automatically by `internal/db/db.go` on startup. Name new migrations `NNN_description.sql` (sequential integer prefix). Never modify existing migration files — add new ones.

### Debugging
```bash
just logs            # All services
docker compose logs backend   # Backend only
docker compose logs frontend  # Frontend only
```

---

## Frontend State Management Rules

- **All business logic lives in Pinia stores**, not in components.
- Components consume stores via `useXxxStore()` + `storeToRefs()` for reactivity.
- `useWebSocketStore` owns the singleton WebSocket connection and auto-reconnect logic.
- `usePlayerStore` persists the player UUID to `localStorage`.
- Components are purely presentational — no direct API calls from `.vue` files; always go through a store action or `lib/api.ts`.

---

## Testing

No automated test suite currently exists. Validation is done via:
- Manual browser testing with `just dev`
- Docker Compose integration with `just up`
- Checking logs with `just logs`

When adding tests in future: use **Vitest + Vue Test Utils** for the frontend and Go's standard `testing` package with table-driven tests for the backend.

---

## Common Pitfalls

- **Do not add REST endpoints** for state reads that should come via WebSocket broadcast. Clients stay in sync through WS, not polling.
- **Do not write ad-hoc CSS** for anything PrimeVue already provides (layout, spacing, color tokens).
- **Do not use `any` in TypeScript.** If types feel complex, model them properly in `types/index.ts`.
- **Do not forget `go mod tidy`** after adding Go dependencies. The Docker build will fail if `go.sum` is out of sync.
- **Do not commit `.env` to git.** Keep `.env.example` synced but never track actual secrets.
- **JWT tokens are time-limited.** Implement refresh token logic on the frontend if longer sessions are needed (currently hardcoded to 7 days).
- When the Postgres trigger fires, it sends the entire room row as JSON. Ensure `RoomStatePayload` in `models/models.go` stays aligned with the trigger's JSON output.
- **Wahapedia sync is a protected endpoint.** Only authenticated players can trigger data refresh. Do not expose it to unauthenticated users.
