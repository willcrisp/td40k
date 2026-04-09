package listen

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/willcrisp/blueprint/internal/models"
	"github.com/willcrisp/blueprint/internal/ws"
)

func StartListener(dsn string, hub *ws.Hub) {
	backoff := 5 * time.Second
	for {
		if err := listen(dsn, hub); err != nil {
			slog.Error("listener error", "err", err, "reconnect_in", backoff)
			time.Sleep(backoff)
			if backoff < 60*time.Second {
				backoff *= 2
			}
		} else {
			backoff = 5 * time.Second
		}
	}
}

func listen(dsn string, hub *ws.Hub) error {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	for _, channel := range []string{"counter_updates", "notes_updates"} {
		if _, err := conn.Exec(context.Background(), "LISTEN "+channel); err != nil {
			return err
		}
	}

	slog.Info("listener ready", "channels", []string{"counter_updates", "notes_updates"})

	for {
		notification, err := conn.WaitForNotification(context.Background())
		if err != nil {
			return err
		}

		var msg models.WsMessage

		switch notification.Channel {
		case "counter_updates":
			var payload models.CounterState
			if err := json.Unmarshal([]byte(notification.Payload), &payload); err != nil {
				slog.Warn("bad counter payload", "err", err)
				continue
			}
			msg = models.WsMessage{Event: "counter_update", Payload: payload}

		case "notes_updates":
			var payload models.NoteEvent
			if err := json.Unmarshal([]byte(notification.Payload), &payload); err != nil {
				slog.Warn("bad notes payload", "err", err)
				continue
			}
			msg = models.WsMessage{Event: "notes_update", Payload: payload}

		default:
			continue
		}

		b, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		hub.Broadcast <- b
	}
}
