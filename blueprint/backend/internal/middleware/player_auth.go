package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const PlayerIDKey contextKey = "player_id"

func RequireAuth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return secret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"invalid claims"}`, http.StatusUnauthorized)
				return
			}

			playerID, ok := claims["player_id"].(string)
			if !ok {
				http.Error(w, `{"error":"invalid player_id"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), PlayerIDKey, playerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetPlayerID(r *http.Request) string {
	id, _ := r.Context().Value(PlayerIDKey).(string)
	return id
}

// RequireAdmin returns 403 if the player does not have is_admin set.
// Must be used after RequireAuth — it reads a separate is_admin context value
// set by a wrapping handler, or you can check at the DB level in your handler.
// For convenience, pair with a DB lookup in your handler if needed.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, _ := r.Context().Value(contextKey("is_admin")).(bool)
		if !isAdmin {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// WithAdmin injects is_admin into the request context. Call this after
// looking up the player from the DB if you need admin-gated routes.
func WithAdmin(r *http.Request, isAdmin bool) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contextKey("is_admin"), isAdmin))
}
