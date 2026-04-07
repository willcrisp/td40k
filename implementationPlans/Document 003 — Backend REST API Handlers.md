Document 003 — Backend: REST API Handlers

Purpose


Implement all HTTP handler functions and database query functions.


---

backend/internal/db/players.go

	package db
	
	import (
	    "context"
	    "github.com/yourorg/w40k/internal/models"
	)
	
	func UpsertPlayer(p models.Player) error {
	    _, err := Pool.Exec(context.Background(), `
	        INSERT INTO players (id, nickname, last_seen)
	        VALUES ($1, $2, NOW())
	        ON CONFLICT (id) DO UPDATE
	            SET nickname  = EXCLUDED.nickname,
	                last_seen = NOW()
	    `, p.ID, p.Nickname)
	    return err
	}
	
	func GetPlayerGames(playerID string) (
	    []models.OwnedGameSummary,
	    []models.JoinedGameSummary,
	    error,
	) {
	    // Owned games
	    ownedRows, err := Pool.Query(context.Background(), `
	        SELECT id, name, status, battle_round, active_player,
	               current_phase, attacker_id, defender_id, created_at
	        FROM rooms
	        WHERE game_master_id = $1
	          AND status != 'closed'
	        ORDER BY created_at DESC
	    `, playerID)
	    if err != nil {
	        return nil, nil, err
	    }
	    defer ownedRows.Close()
	
	    var owned []models.OwnedGameSummary
	    for ownedRows.Next() {
	        var g models.OwnedGameSummary
	        if err := ownedRows.Scan(
	            &g.ID, &g.Name, &g.Status, &g.BattleRound,
	            &g.ActivePlayer, &g.CurrentPhase,
	            &g.AttackerID, &g.DefenderID, &g.CreatedAt,
	        ); err != nil {
	            return nil, nil, err
	        }
	        owned = append(owned, g)
	    }
	
	    // Joined games (attacker or defender, not GM)
	    joinedRows, err := Pool.Query(context.Background(), `
	        SELECT
	            r.id, r.name, r.status, r.battle_round,
	            r.current_phase, r.created_at,
	            CASE
	                WHEN r.attacker_id = $1 THEN 'attacker'
	                WHEN r.defender_id = $1 THEN 'defender'
	            END AS role
	        FROM rooms r
	        WHERE (r.attacker_id = $1 OR r.defender_id = $1)
	          AND r.game_master_id != $1
	          AND r.status != 'closed'
	        ORDER BY r.created_at DESC
	    `, playerID)
	    if err != nil {
	        return nil, nil, err
	    }
	    defer joinedRows.Close()
	
	    var joined []models.JoinedGameSummary
	    for joinedRows.Next() {
	        var g models.JoinedGameSummary
	        if err := joinedRows.Scan(
	            &g.ID, &g.Name, &g.Status, &g.BattleRound,
	            &g.CurrentPhase, &g.CreatedAt, &g.Role,
	        ); err != nil {
	            return nil, nil, err
	        }
	        joined = append(joined, g)
	    }
	
	    return owned, joined, nil
	}


---

backend/internal/db/rooms.go

	package db
	
	import (
	    "context"
	    "fmt"
	
	    "github.com/yourorg/w40k/internal/models"
	)
	
	func CreateRoom(r models.Room) error {
	    _, err := Pool.Exec(context.Background(), `
	        INSERT INTO rooms
	            (id, name, status, game_master_id, battle_round,
	             active_player, current_phase)
	        VALUES ($1,$2,'lobby',$3,1,'attacker','command')
	    `, r.ID, r.Name, r.GameMasterID)
	    return err
	}
	
	func GetRoom(id string) (*models.Room, error) {
	    row := Pool.QueryRow(context.Background(), `
	        SELECT id, name, status, game_master_id, attacker_id, defender_id,
	               battle_round, active_player, current_phase,
	               winner, created_at, updated_at
	        FROM rooms WHERE id = $1
	    `, id)
	    var r models.Room
	    if err := row.Scan(
	        &r.ID, &r.Name, &r.Status, &r.GameMasterID,
	        &r.AttackerID, &r.DefenderID,
	        &r.BattleRound, &r.ActivePlayer, &r.CurrentPhase,
	        &r.Winner, &r.CreatedAt, &r.UpdatedAt,
	    ); err != nil {
	        return nil, err
	    }
	    return &r, nil
	}
	
	func SetRoomAttacker(roomID, playerID string) error {
	    res, err := Pool.Exec(context.Background(), `
	        UPDATE rooms SET attacker_id = $1
	        WHERE id = $2 AND attacker_id IS NULL
	    `, playerID, roomID)
	    if err != nil {
	        return err
	    }
	    if res.RowsAffected() == 0 {
	        return fmt.Errorf("attacker slot already taken")
	    }
	    return nil
	}
	
	func SetRoomDefender(roomID, playerID string) error {
	    res, err := Pool.Exec(context.Background(), `
	        UPDATE rooms SET defender_id = $1
	        WHERE id = $2 AND defender_id IS NULL
	    `, playerID, roomID)
	    if err != nil {
	        return err
	    }
	    if res.RowsAffected() == 0 {
	        return fmt.Errorf("defender slot already taken")
	    }
	    return nil
	}
	
	func StartRoom(roomID string) error {
	    _, err := Pool.Exec(context.Background(), `
	        UPDATE rooms SET status = 'active'
	        WHERE id = $1
	          AND attacker_id IS NOT NULL
	          AND defender_id IS NOT NULL
	          AND status = 'lobby'
	    `, roomID)
	    return err
	}
	
	func UpdateRoomPhase(
	    roomID, phase, activePlayer string,
	    battleRound int,
	    winner *string,
	    status string,
	) error {
	    _, err := Pool.Exec(context.Background(), `
	        UPDATE rooms
	        SET current_phase  = $1,
	            active_player  = $2,
	            battle_round   = $3,
	            winner         = $4,
	            status         = $5
	        WHERE id = $6
	    `, phase, activePlayer, battleRound, winner, status, roomID)
	    return err
	}
	
	func CloseRoom(roomID string) error {
	    _, err := Pool.Exec(context.Background(), `
	        UPDATE rooms SET status = 'closed' WHERE id = $1
	    `, roomID)
	    return err
	}
	
	func LogEvent(
	    roomID string,
	    playerID *string,
	    eventType string,
	    payload []byte,
	) error {
	    _, err := Pool.Exec(context.Background(), `
	        INSERT INTO room_events (room_id, player_id, event_type, payload)
	        VALUES ($1, $2, $3, $4)
	    `, roomID, playerID, eventType, payload)
	    return err
	}


---

backend/internal/handlers/players.go

	package handlers
	
	import (
	    "encoding/json"
	    "net/http"
	
	    "github.com/go-chi/chi/v5"
	    "github.com/yourorg/w40k/internal/db"
	    mw "github.com/yourorg/w40k/internal/middleware"
	    "github.com/yourorg/w40k/internal/models"
	)
	
	func HandleUpsertPlayer(w http.ResponseWriter, r *http.Request) {
	    var body struct {
	        ID       string `json:"id"`
	        Nickname string `json:"nickname"`
	    }
	    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
	        jsonError(w, "invalid request body", http.StatusBadRequest)
	        return
	    }
	    if body.ID == "" || body.Nickname == "" {
	        jsonError(w, "id and nickname required", http.StatusBadRequest)
	        return
	    }
	    if err := db.UpsertPlayer(models.Player{
	        ID: body.ID, Nickname: body.Nickname,
	    }); err != nil {
	        jsonError(w, "db error", http.StatusInternalServerError)
	        return
	    }
	    w.WriteHeader(http.StatusOK)
	    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
	
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


---

backend/internal/handlers/rooms.go

	package handlers
	
	import (
	    "encoding/json"
	    "fmt"
	    "net/http"
	
	    "github.com/go-chi/chi/v5"
	    "github.com/google/uuid"
	    "github.com/yourorg/w40k/internal/db"
	    mw "github.com/yourorg/w40k/internal/middleware"
	    "github.com/yourorg/w40k/internal/models"
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
	        jsonError(w, "both roles must be filled before starting",
	            http.StatusBadRequest)
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
	        // Winner determination is out of scope Phase 1 — set to nil
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
	
	// ── Phase logic ──────────────────────────────────────────────────────────────
	
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
	
	// ── Helpers ──────────────────────────────────────────────────────────────────
	
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