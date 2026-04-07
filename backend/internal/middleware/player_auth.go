package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const PlayerIDKey contextKey = "player_id"

// RequirePlayerID extracts X-Player-ID header and injects it into context.
// Returns 401 if the header is missing or empty.
func RequirePlayerID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pid := r.Header.Get("X-Player-ID")
		if pid == "" {
			http.Error(w, `{"error":"missing X-Player-ID header"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), PlayerIDKey, pid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPlayerID(r *http.Request) string {
	v, _ := r.Context().Value(PlayerIDKey).(string)
	return v
}
