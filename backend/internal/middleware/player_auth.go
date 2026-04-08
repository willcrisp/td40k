package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const PlayerIDKey contextKey = "player_id"

type jwtClaims struct {
	PlayerID string `json:"player_id"`
	jwt.RegisteredClaims
}

// RequireAuth validates a Bearer JWT and injects the player_id into context.
// Returns 401 if the token is missing, malformed, or expired.
func RequireAuth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return secret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			c, ok := token.Claims.(*jwtClaims)
			if !ok || c.PlayerID == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), PlayerIDKey, c.PlayerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetPlayerID retrieves the authenticated player ID from the request context.
func GetPlayerID(r *http.Request) string {
	v, _ := r.Context().Value(PlayerIDKey).(string)
	return v
}
