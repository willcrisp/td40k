package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const PlayerIDKey contextKey = "player_id"

// ExtractPlayerID extracts the X-Player-ID header and injects player_id
// into context. All requests are accepted; clients are responsible for
// providing a valid UUID.
func ExtractPlayerID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playerID := r.Header.Get("X-Player-ID")
		ctx := context.WithValue(r.Context(), PlayerIDKey, playerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPlayerID retrieves the authenticated player ID from the request context.
func GetPlayerID(r *http.Request) string {
	v, _ := r.Context().Value(PlayerIDKey).(string)
	return v
}
