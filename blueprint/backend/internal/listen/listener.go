package listen

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/willcrisp/blueprint/internal/models"
	"github.com/willcrisp/blueprint/internal/ws"
)

func StartListener(dsn string, hub *ws.Hub) {
	for {
		if err := listen(dsn, hub); err != nil {
			fmt.Printf("listener error: %v — reconnecting in 5s\n", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func listen(dsn string, hub *ws.Hub) error {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close(context.Background())

	if _, err := conn.Exec(context.Background(), "LISTEN counter_updates"); err != nil {
		return fmt.Errorf("LISTEN: %w", err)
	}

	fmt.Println("listener: watching counter_updates")

	for {
		notification, err := conn.WaitForNotification(context.Background())
		if err != nil {
			return fmt.Errorf("WaitForNotification: %w", err)
		}

		var payload models.CounterState
		if err := json.Unmarshal([]byte(notification.Payload), &payload); err != nil {
			fmt.Printf("listener: bad payload: %v\n", err)
			continue
		}

		msg := models.WsMessage{
			Event:   "counter_update",
			Payload: payload,
		}
		b, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		hub.Broadcast <- b
	}
}
