# CLAUDE.md — blueprint

A minimal full-stack starter: JWT auth, real-time WebSocket sync via PostgreSQL `LISTEN/NOTIFY`, and two REST examples (a shared counter + a shared notes list). Use this as the foundation for any Go + Vue 3 + PostgreSQL project.

---

## Absolute Rules

1. **PostgreSQL is the only source of truth.** All state changes persist to the DB before any response is sent.
2. **WebSocket is broadcast-only from the server.** Clients never send WS messages. All mutations go through HTTP REST endpoints.
3. **Go handlers always return JSON.** Never return plain text errors except for HTTP 500 fallbacks.
4. **No `any` types in TypeScript.** All Vue components and stores must be fully typed.
5. **PrimeVue 4 is the only UI library.** No ad-hoc CSS beyond layout necessities.
6. **Prettier enforces formatting.** Print width: 80, singleQuote, semi, trailingComma: es5.

---

## Tech Stack

| Layer    | Technology                  |
|----------|-----------------------------|
| Backend  | Go 1.22, chi v5, pgx v5, gorilla/websocket, golang-jwt/jwt v5, bcrypt, slog |
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
│       │   ├── players.go             # Player CRUD (CreatePlayer, GetPlayerByUsername)
│       │   ├── counter.go             # GetCounter, IncrementCounter
│       │   └── notes.go               # ListNotes, CreateNote, DeleteNote
│       ├── handlers/
│       │   ├── auth.go                # HandleRegister, HandleLogin → JWT
│       │   ├── counter.go             # HandleGetCounter, HandleIncrementCounter
│       │   └── notes.go               # HandleListNotes, HandleCreateNote, HandleDeleteNote
│       ├── middleware/
│       │   └── player_auth.go         # RequireAuth, GetPlayerID, RequireAdmin, WithAdmin
│       ├── models/
│       │   └── models.go              # Player, CounterState, Note, NoteEvent, WsMessage
│       ├── ws/
│       │   ├── hub.go                 # Gorilla WebSocket hub: broadcast to all clients
│       │   └── client.go              # Per-client writePump / readPump goroutines
│       └── listen/
│           └── listener.go            # LISTEN counter_updates + notes_updates → hub.Broadcast
├── frontend/src/
│   ├── main.ts                        # Pinia + PrimeVue (Aura) + ToastService + Router
│   ├── App.vue                        # Root: <Toast />, <RouterView />, playerStore.init()
│   ├── router/index.ts                # /auth (redirectIfAuthenticated), / (requirePlayer)
│   ├── stores/
│   │   ├── usePlayerStore.ts          # JWT + isAdmin localStorage persistence
│   │   ├── useCounterStore.ts         # value, fetchCounter(), increment(), applyUpdate()
│   │   ├── useNotesStore.ts           # notes[], fetchNotes(), createNote(), deleteNote(), applyInsert/Delete()
│   │   └── useWebSocketStore.ts       # WS singleton, exponential backoff, routes events → stores
│   ├── views/
│   │   ├── LoginView.vue              # PrimeVue Tabs: Login / Register
│   │   └── HomeView.vue               # Counter + Notes UI
│   ├── lib/api.ts                     # Axios client + JWT interceptor + 401 handler + all api* functions
│   └── types/index.ts                 # AuthResponse, CounterState, Note, NoteEvent, WsMessage, CounterUpdatePayload
├── db/migrations/
│   ├── 001_create_players.sql
│   ├── 002_create_counter.sql
│   ├── 003_create_triggers.sql        # counter trigger → pg_notify
│   ├── 004_create_notes.sql           # notes table + trigger → pg_notify
│   └── 005_add_admin_to_players.sql   # is_admin boolean column
├── docker-compose.yml
├── docker-compose.dev.yml
└── justfile
```

---

## Data Flow

```
Client (Vue)  ──── POST /api/notes ────►  Go Handler
                                               │
                                      INSERT INTO notes
                                               │
                              Trigger: pg_notify('notes_updates', {op, id, content, ...})
                                               │
                                       Go listen goroutine
                                               │
                                       WebSocket Hub → Broadcast
                                               │
Client (Vue)  ◄──── {"event":"notes_update","payload":{"op":"insert",...}} ────
                              │
                  useWebSocketStore → useNotesStore.applyInsert(note) → UI
```

Same pattern for counter (`counter_updates` channel) and for deletions (`op: "delete"`).

---

## HTTP Endpoints

| Method | Path | Auth | Handler |
|--------|------|------|---------|
| GET | /health | No | inline (status ok) |
| POST | /api/auth/register | No | HandleRegister |
| POST | /api/auth/login | No | HandleLogin |
| GET | /ws | No | hub.ServeWS |
| GET | /api/counter | Yes | HandleGetCounter |
| POST | /api/counter/increment | Yes | HandleIncrementCounter |
| GET | /api/notes | Yes | HandleListNotes |
| POST | /api/notes | Yes | HandleCreateNote |
| DELETE | /api/notes/{id} | Yes | HandleDeleteNote (own notes only) |

Auth: `Authorization: Bearer <token>`. JWT payload: `{"player_id": "uuid", "exp": ...}`.

On 401 the Axios interceptor clears localStorage and redirects to `/auth` automatically.

---

## WebSocket Events

| Event | Direction | Payload |
|-------|-----------|---------|
| `counter_update` | Server → Client | `{ value: number }` |
| `notes_update` | Server → Client | `{ op: "insert", id, player_id, username, content, created_at }` or `{ op: "delete", id }` |

WS reconnect uses exponential backoff: 1s → 2s → 4s … → 30s max.

---

## Admin Flag

Players have an `is_admin` boolean (default `false`). It's returned in the auth response and stored in localStorage. Use `RequireAdmin` middleware after `RequireAuth` to gate admin-only routes. Use `WithAdmin(r, isAdmin)` to inject the flag into the request context after a DB lookup.

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
just psql     # Open a psql shell in the running DB container
just reset    # Wipe DB volume + rebuild from scratch
```

---

## Adding a New Real-Time Resource

Notes is the template for a collection. Counter is the template for a scalar. To add a new resource:

1. **Migration** — add `NNN_description.sql` in `db/migrations/` with table + `pg_notify` trigger. Never modify existing files.
2. **DB layer** — add `backend/internal/db/yourresource.go` with List/Create/Delete functions.
3. **Model** — add your struct and a `YourEvent` struct (with `op` field for mutations) to `models/models.go`.
4. **Handler** — add `backend/internal/handlers/yourresource.go`, wire routes in `main.go`.
5. **Listener** — add your channel to the `LISTEN` list in `listen/listener.go` and a `case` in the switch.
6. **Frontend types** — add your types to `src/types/index.ts`.
7. **Store** — add `src/stores/useYourStore.ts` with `applyInsert` / `applyDelete` (or `applyUpdate` for scalars).
8. **WS router** — add a case in `useWebSocketStore` for your event name.
9. **View** — consume the store in a view or component.

---

## Naming Conventions

**Go:** `snake_case` files, `PascalCase` exports, `Handle{Action}{Resource}` handlers, one file per entity in `internal/db/`.

**TypeScript/Vue:** `use{Name}Store.ts` stores, `PascalCase.vue` components, `api{Action}{Resource}` API functions in `lib/api.ts`.
