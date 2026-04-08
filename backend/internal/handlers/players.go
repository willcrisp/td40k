package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willcrisp/td40k/internal/db"
	mw "github.com/willcrisp/td40k/internal/middleware"
	"github.com/willcrisp/td40k/internal/models"
)

func HandleGetPlayerGames(w http.ResponseWriter, r *http.Request) {
	playerID := chi.URLParam(r, "id")
	callerID := mw.GetPlayerID(r)
	if playerID != callerID {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}
	owned, joined, err := db.GetPlayerGames(playerID)
	if err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	if owned == nil {
		owned = []models.OwnedGameSummary{}
	}
	if joined == nil {
		joined = []models.JoinedGameSummary{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Owned  []models.OwnedGameSummary  `json:"owned"`
		Joined []models.JoinedGameSummary `json:"joined"`
	}{
		Owned:  owned,
		Joined: joined,
	})
}
