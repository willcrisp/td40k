package listen

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/willcrisp/td40k/internal/models"
	"github.com/willcrisp/td40k/internal/ws"
)

func StartListener(dsn string, hub *ws.Hub) {
	for {
		if err := listen(dsn, hub); err != nil {
			log.Printf("[listen] error: %v — reconnecting in 5s", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func listen(dsn string, hub *ws.Hub) error {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	if _, err := conn.Exec(
		context.Background(), "LISTEN room_updates",
	); err != nil {
		return err
	}
	log.Println("[listen] listening on channel: room_updates")

	for {
		notification, err := conn.WaitForNotification(context.Background())
		if err != nil {
			return err
		}
		var payload models.RoomStatePayload
		if err := json.Unmarshal(
			[]byte(notification.Payload), &payload,
		); err != nil {
			log.Printf("[listen] unmarshal error: %v", err)
			continue
		}
		hub.Broadcast(&payload)
	}
}
