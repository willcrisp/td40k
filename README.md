# Counter App — Vue + Go + PostgreSQL

A minimal full-stack counter app running entirely in Docker.

## Stack
| Layer    | Tech                              |
|----------|-----------------------------------|
| Frontend | Vue 3, Vite, Bun, PrimeVue 4     |
| Backend  | Go 1.22, chi router, pgx          |
| Database | PostgreSQL 16                     |
| Infra    | Docker Compose, nginx             |

## Quick Start

```bash
docker compose up --build
```

| Service  | URL                        |
|----------|----------------------------|
| Frontend | http://localhost:3000       |
| Backend  | http://localhost:8080/api  |
| Postgres | localhost:5432             |

## API Endpoints

| Method | Path                      | Description          |
|--------|---------------------------|----------------------|
| GET    | /api/health               | Health check         |
| GET    | /api/counter              | Get current value    |
| POST   | /api/counter/increment    | Increment by 1       |

## Architecture Flow

This application embraces the Single Source of Truth pattern using PostgreSQL as the central circulatory system.

```text
  ┌──────────────┐         REST (Mutations)         ┌──────────────┐
  │              ├─────────────────────────────────►│              │
  │   Frontend   │                                  │              │
  │  (Vue/Vite)  │◄─────────────────────────────────┤  Go Backend  │
  │              │     WebSocket (Live Updates)     │              │
  └──────┬───────┘                                  └──────┬───────┘
         │                                                 │ ▲
         │                                                 │ │ LISTEN /
         │   (No Direct Database Link)         SQL SELECT/ │ │ NOTIFY
         │                                     UPDATE      │ │ 
         │                                                 ▼ │
         │                                          ┌──────────────┐
         └────────────────X─────────────────────────┤              │
                                                    │  PostgreSQL  │
                                                    │              │
                                                    └──────────────┘
```

### The Reactive Lifecycle
1. **Mutation:** The user clicks "Increment", firing a standard REST `POST /api/counter/increment` request.
2. **Database Update:** The Go backend executes a standard PostgreSQL `UPDATE` query. 
3. **Database Trigger:** PostgreSQL natively intercepts the row change via a trigger and issues a `NOTIFY` event payload containing the new row state.
4. **Backend Listen:** The Go backend, maintaining a persistent `LISTEN` connection, intercepts the database notification instantly.
5. **WebSocket Broadcast:** The backend parses the event and blind-fires a JSON broadcast to *every* connected Vue client over their WebSockets.
6. **Reactivity:** The Vue frontend receives the payload, compares the room ID, updates the reactive reactive state variables, and CSS-animates the change to the user.

### Code-Level Data Flow

```text
  ┌───────────────────────────────────────────────┐
  │  Client UI (Vue Component)                    │
  │  File: frontend/.../CounterCard.vue           │
  └──────────────────────┬────────────────────────┘
                         │ 1. User clicks
                         ▼    "Increment"
  ┌───────────────────────────────────────────────┐
  │  Domain Logic (Vue Composable)                │
  │  File: frontend/.../useCounter.ts             │
  └──────────────────────┬────────────────────────┘
                         │ 2. REST POST  
                         ▼    /api/counter/increment
  ┌───────────────────────────────────────────────┐
  │  API Endpoint (Go Router Handler)             │
  │  File: backend/counter/handlers.go            │
  └──────────────────────┬────────────────────────┘
                         │ 3. Call Repository
                         ▼
  ┌───────────────────────────────────────────────┐
  │  Database Model (Go UPSERT Logic)             │
  │  File: backend/counter/models.go              │
  └──────────────────────┬────────────────────────┘
                         │ 4. Execute SQL
                         ▼
  ┌───────────────────────────────────────────────┐
  │  PostgreSQL Execution & Triggers              │
  │  File: backend/counter/migrations/*.sql       │
  └──────────────────────┬────────────────────────┘
                         │ 5. Trigger fires
                         ▼    pg_notify()
  ┌───────────────────────────────────────────────┐
  │  Database Listener (Go Goroutine)             │
  │  File: backend/counter/websockets.go          │
  └──────────────────────┬────────────────────────┘
                         │ 6. Parse payload &
                         ▼    Broadcast()
  ┌───────────────────────────────────────────────┐
  │  WebSocket Broadcaster (Gorilla Hub)          │
  │  File: backend/websocket/hub.go               │
  └──────────────────────┬────────────────────────┘
                         │ 7. Push raw JSON 
                         ▼    over ws:// pipe
  ┌───────────────────────────────────────────────┐
  │  Client Network Engine (Vue Composable)       │
  │  File: frontend/.../useWebSocket.ts           │
  └──────────────────────┬────────────────────────┘
                         │ 8. Parse JSON & 
                         ▼    pass to callback
  ┌───────────────────────────────────────────────┐
  │  Domain Logic (Reactive Update)               │
  │  File: frontend/.../useCounter.ts             │
  └───────────────────────────────────────────────┘
  (Triggers instant UI re-render on CounterCard.vue)
```
