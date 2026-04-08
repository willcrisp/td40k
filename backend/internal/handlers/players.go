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
	json.NewEncoder(w).Encode(map[string]any{
		"owned":  owned,
		"joined": joined,
	})
}
