package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willcrisp/td40k/internal/db"
	mw "github.com/willcrisp/td40k/internal/middleware"
	"github.com/willcrisp/td40k/internal/models"
)

func HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		jsonError(w, "name required", http.StatusBadRequest)
		return
	}
	roomID := generateRoomID()
	if err := db.CreateRoom(models.Room{
		ID: roomID, Name: body.Name, GameMasterID: playerID,
	}); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	db.LogEvent(roomID, &playerID, "game_created", nil)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": roomID})
}

func HandleGetRoom(w http.ResponseWriter, r *http.Request) {
	room, err := db.GetRoom(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "room not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(room)
}

func HandleJoinRoom(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "id")
	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Role != "attacker" && body.Role != "defender" {
		jsonError(w, "role must be attacker or defender", http.StatusBadRequest)
		return
	}
	var err error
	if body.Role == "attacker" {
		err = db.SetRoomAttacker(roomID, playerID)
	} else {
		err = db.SetRoomDefender(roomID, playerID)
	}
	if err != nil {
		jsonError(w, err.Error(), http.StatusConflict)
		return
	}
	db.LogEvent(roomID, &playerID, "player_joined",
		jsonBytes(map[string]string{"role": body.Role}))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandleStartGame(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "id")
	room, err := db.GetRoom(roomID)
	if err != nil {
		jsonError(w, "room not found", http.StatusNotFound)
		return
	}
	if room.GameMasterID != playerID {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}
	if room.AttackerID == nil || room.DefenderID == nil {
		jsonError(w, "both roles must be filled before starting", http.StatusBadRequest)
		return
	}
	if err := db.StartRoom(roomID); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	db.LogEvent(roomID, &playerID, "game_started", nil)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandlePhaseNext(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "id")
	room, err := db.GetRoom(roomID)
	if err != nil {
		jsonError(w, "room not found", http.StatusNotFound)
		return
	}
	if room.GameMasterID != playerID {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}
	if room.Status != "active" {
		jsonError(w, "game is not active", http.StatusBadRequest)
		return
	}
	phase, activePlayer, round, gameOver := advancePhase(room)
	status := "active"
	var winner *string
	if gameOver {
		status = "finished"
	}
	if err := db.UpdateRoomPhase(
		roomID, phase, activePlayer, round, winner, status,
	); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	db.LogEvent(roomID, &playerID, "phase_advanced",
		jsonBytes(map[string]any{
			"phase": phase, "round": round,
			"active_player": activePlayer,
		}))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandlePhasePrev(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "id")
	room, err := db.GetRoom(roomID)
	if err != nil {
		jsonError(w, "room not found", http.StatusNotFound)
		return
	}
	if room.GameMasterID != playerID {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}
	if room.Status != "active" {
		jsonError(w, "game is not active", http.StatusBadRequest)
		return
	}
	phase, activePlayer, round := retreatPhase(room)
	if err := db.UpdateRoomPhase(
		roomID, phase, activePlayer, round, room.Winner, "active",
	); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	db.LogEvent(roomID, &playerID, "phase_retreated",
		jsonBytes(map[string]any{
			"phase": phase, "round": round,
			"active_player": activePlayer,
		}))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandleCloseRoom(w http.ResponseWriter, r *http.Request) {
	playerID := mw.GetPlayerID(r)
	roomID := chi.URLParam(r, "id")
	room, err := db.GetRoom(roomID)
	if err != nil {
		jsonError(w, "room not found", http.StatusNotFound)
		return
	}
	if room.GameMasterID != playerID {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := db.CloseRoom(roomID); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	db.LogEvent(roomID, &playerID, "game_closed", nil)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ── Phase logic ───────────────────────────────────────────────────────────────

var phases = []string{"command", "movement", "shooting", "charge", "fight"}

func indexOf(slice []string, val string) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

func advancePhase(room *models.Room) (
	phase, activePlayer string, battleRound int, gameOver bool,
) {
	idx := indexOf(phases, room.CurrentPhase)
	if idx < len(phases)-1 {
		return phases[idx+1], room.ActivePlayer, room.BattleRound, false
	}
	if room.ActivePlayer == "attacker" {
		return "command", "defender", room.BattleRound, false
	}
	if room.BattleRound >= 5 {
		return "fight", "defender", 5, true
	}
	return "command", "attacker", room.BattleRound + 1, false
}

func retreatPhase(room *models.Room) (phase, activePlayer string, battleRound int) {
	idx := indexOf(phases, room.CurrentPhase)
	if idx > 0 {
		return phases[idx-1], room.ActivePlayer, room.BattleRound
	}
	if room.ActivePlayer == "defender" {
		return "fight", "attacker", room.BattleRound
	}
	if room.BattleRound > 1 {
		return "fight", "defender", room.BattleRound - 1
	}
	return room.CurrentPhase, room.ActivePlayer, room.BattleRound
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func jsonBytes(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func generateRoomID() string {
	adjectives := []string{
		"iron", "flame", "void", "storm", "blood",
		"dark", "grim", "holy", "chaos", "death",
	}
	nouns := []string{
		"wolf", "eagle", "fist", "blade", "skull",
		"angel", "guard", "titan", "raven", "lance",
	}
	id := uuid.New().String()[:4]
	adj := adjectives[int(id[0])%len(adjectives)]
	noun := nouns[int(id[1])%len(nouns)]
	return fmt.Sprintf("%s-%s-%s", adj, noun, id)
}
