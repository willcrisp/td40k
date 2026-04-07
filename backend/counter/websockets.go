package counter

import (
	"context"
	"log"
	"time"
)

// ListenForUpdates encapsulates the specific Postgres trigger listen loop.
// This is the true bridge between your Database and your Websocket Clients!
// It runs entirely asynchronously in an infinite Goroutine loop.
func (m *Module) ListenForUpdates() {
	// 1. We acquire a dedicated, persistent physical TCP connection to PostgreSQL
	// We MUST hold a direct line open to listen for pushed notifications.
	conn, err := m.Repo.DB.Acquire(context.Background())
	if err != nil {
		log.Printf("[Counter] Critical Listener acquire error: %v", err)
		return
	}
	// Guarantee the connection returns to the pool if this function ever crashes
	defer conn.Release()

	// 2. Execute the Postgres LISTEN command.
	// 'counter_updates' is the literal channel name that our Postgres TRIGGER uses in `pg_notify`.
	_, err = conn.Exec(context.Background(), "LISTEN counter_updates")
	if err != nil {
		log.Printf("[Counter] Failed to execute Postgres LISTEN: %v", err)
		return
	}
	
	log.Println("[Counter] Successfully subscribed to Postgres Pub/Sub channel: 'counter_updates'")

	// 3. Enter an infinite blocking loop.
	for {
		// WaitForNotification puts the thread to sleep until Postgres actively pushes data down the pipe.
		// THIS is what eliminates standard API polling! The database pushes to Go instantly.
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			log.Printf("[Counter] Notification wait error, retrying... %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 4. We caught an event! `notification.Payload` contains the pure JSON string emitted by the SQL trigger.
		// We blindly pass this string directly to our generic Hub broadcaster, which blasts it
		// perfectly raw to all connected Vue browser clients instantly over WebSockets.
		m.Broadcast([]byte(notification.Payload))
	}
}
