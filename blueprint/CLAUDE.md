# CLAUDE.md — blueprint

A minimal full-stack starter: JWT auth, real-time WebSocket sync via PostgreSQL `LISTEN/NOTIFY`, and one REST example (a shared counter). Use this as the foundation for any Go + Vue 3 + PostgreSQL project.

---

## Absolute Rules

1. **PostgreSQL is the only source of truth.** All state changes persist to the DB before any response is sent.
2. **WebSocket is broadcast-only from the server.** Clients never send WS messages. All mutations go through HTTP REST endpoints.
3. **Go handlers always return JSON.** Never return plain text errors except for HTTP 500 fallbacks.
4. **No `any` types in TypeScript.** All Vue components and stores must be fully typed.
5. **PrimeVue 4 is the only UI library.** No ad-hoc CSS beyond canvas/layout necessities.
6. **Prettier enforces formatting.** Print width: 80, singleQuote, semi, trailingComma: es5.

---

## Tech Stack

| Layer    | Technology                  |
|----------|-----------------------------|
| Backend  | Go 1.22, chi v5, pgx v5, gorilla/websocket, golang-jwt/jwt v5, bcrypt |
| Frontend | Vue 3 (Composition API), Vite 5, TypeScript 5, Pinia, PrimeVue 4, Axios, Bun |
| Database | PostgreSQL 16               |
| Infra    | Docker Compose + nginx      |

---

## Repository Layout

```
blueprint/
├── backend/
│   ├── cmd/server/main.go             # Entry point: routing, middleware, server start
│   └── internal/
│       ├── db/
│       │   ├── db.go                  # pgxpool init & migration runner
│       │   ├── players.go             # Player CRUD
│       │   └── counter.go             # Counter read/increment
│       ├── handlers/
│       │   ├── auth.go                # Register/Login → JWT
│       │   └── counter.go             # GET /api/counter, POST /api/counter/increment
│       ├── middleware/
│       │   └── player_auth.go         # JWT verification & context injection
│       ├── models/
│       │   └── models.go              # Player, CounterState, WsMessage
│       ├── ws/
│       │   ├── hub.go                 # Gorilla WebSocket hub: broadcast to all clients
│       │   └── client.go              # Per-client writePump / readPump goroutines
│       └── listen/
│           └── listener.go            # LISTEN counter_updates → hub.Broadcast
├── frontend/src/
│   ├── main.ts                        # Pinia + PrimeVue (Aura) + Router bootstrap
│   ├── App.vue                        # Root: RouterView, playerStore.init()
│   ├── router/index.ts                # /auth (redirectIfAuthenticated), / (requirePlayer)
│   ├── stores/
│   │   ├── usePlayerStore.ts          # JWT localStorage persistence
│   │   ├── useCounterStore.ts         # value, fetchCounter(), increment(), applyUpdate()
│   │   └── useWebSocketStore.ts       # WS singleton, routes counter_update → counterStore
│   ├── views/
│   │   ├── LoginView.vue              # PrimeVue Tabs: Login / Register
│   │   └── HomeView.vue               # Counter display + Increment button
│   ├── lib/api.ts                     # Axios client + JWT interceptor + all api* functions
│   └── types/index.ts                 # AuthResponse, CounterState, WsMessage, CounterUpdatePayload
├── db/migrations/
│   ├── 001_create_players.sql
│   ├── 002_create_counter.sql
│   └── 003_create_triggers.sql
├── docker-compose.yml
├── docker-compose.dev.yml
└── justfile
```

---

## Data Flow

```
Client (Vue)  ──── POST /api/counter/increment ────►  Go Handler
                                                           │
                                                    UPDATE counter
                                                           │
                                                  Trigger: pg_notify('counter_updates', ...)
                                                           │
                                                    Go listen goroutine
                                                           │
                                                    WebSocket Hub → Broadcast
                                                           │
Client (Vue)  ◄──── {"event":"counter_update","payload":{"value":N}} ────
                           │
                  useWebSocketStore → useCounterStore.applyUpdate(N) → UI
```

---

## HTTP Endpoints

| Method | Path | Auth | Handler |
|--------|------|------|---------|
| POST | /api/auth/register | No | HandleRegister |
| POST | /api/auth/login | No | HandleLogin |
| GET | /ws | No | hub.ServeWS |
| GET | /api/counter | Yes | HandleGetCounter |
| POST | /api/counter/increment | Yes | HandleIncrementCounter |

Auth: `Authorization: Bearer <token>`. JWT payload: `{"player_id": "uuid", "exp": ...}`.

---

## Environment Variables

See `.env.example`. Key variables:

| Variable | Example |
|----------|---------|
| `POSTGRES_DSN` | `postgres://app:app@db:5432/app` |
| `JWT_SECRET` | `openssl rand -base64 32` |
| `VITE_API_BASE_URL` | `http://localhost:8080` |
| `VITE_WS_BASE_URL` | `ws://localhost:8080` |

---

## Key Commands

```bash
just up       # Build + start all containers (prod, port 3000)
just dev      # Vite HMR dev server on :5173 (volume-mounted src)
just down     # Stop containers
just logs     # Follow all service logs
just reset    # Wipe DB volume + rebuild from scratch
```

---

## Adding a New Feature

The counter is the template. To add a new real-time resource:

1. **Migration** — add `NNN_description.sql` in `db/migrations/` (never modify existing files)
2. **DB layer** — add a file in `backend/internal/db/`
3. **Handler** — add a file in `backend/internal/handlers/`, wire it in `main.go`
4. **Trigger** — `pg_notify('your_channel', json_build_object(...)::text)`
5. **Listener** — update `listen/listener.go` to handle the new channel (or add a second listener)
6. **Frontend types** — extend `src/types/index.ts`
7. **Store** — add `src/stores/useYourStore.ts` with `applyUpdate()`
8. **WS router** — add a case in `useWebSocketStore` for the new event name
9. **View** — consume the store in a view or component

---

## Naming Conventions

**Go:** `snake_case` files, `PascalCase` exports, `Handle{Action}{Resource}` handlers, one file per entity in `internal/db/`.

**TypeScript/Vue:** `use{Name}Store.ts` stores, `PascalCase.vue` components, `api{Action}{Resource}` API functions in `lib/api.ts`.
