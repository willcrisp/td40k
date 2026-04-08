package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/willcrisp/td40k/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type claims struct {
	PlayerID string `json:"player_id"`
	jwt.RegisteredClaims
}

func issueToken(playerID string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	c := claims{
		PlayerID: playerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(secret)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if body.Username == "" || body.Nickname == "" || body.Password == "" {
		jsonError(w, "username, nickname and password are required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	player, err := db.CreatePlayer(body.Username, body.Nickname, string(hash))
	if err != nil {
		if err.Error() == "username taken" {
			jsonError(w, "username already taken", http.StatusConflict)
			return
		}
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}

	token, err := issueToken(player.ID)
	if err != nil {
		jsonError(w, "could not issue token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"token":     token,
		"player_id": player.ID,
		"username":  player.Username,
		"nickname":  player.Nickname,
		"is_admin":  player.IsAdmin,
	})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if body.Username == "" || body.Password == "" {
		jsonError(w, "username and password are required", http.StatusBadRequest)
		return
	}

	player, hash, err := db.GetPlayerByUsername(body.Username)
	if err != nil {
		// Don't reveal whether username exists
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password)); err != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := issueToken(player.ID)
	if err != nil {
		jsonError(w, "could not issue token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"token":     token,
		"player_id": player.ID,
		"username":  player.Username,
		"nickname":  player.Nickname,
		"is_admin":  player.IsAdmin,
	})
}
