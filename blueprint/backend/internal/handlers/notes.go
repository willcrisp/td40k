package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willcrisp/blueprint/internal/db"
	mw "github.com/willcrisp/blueprint/internal/middleware"
)

func HandleListNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := db.ListNotes()
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func HandleCreateNote(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Content == "" {
		jsonError(w, "content is required", http.StatusBadRequest)
		return
	}

	userID := mw.GetUserID(r)
	note, err := db.CreateNote(uuid.NewString(), userID, body.Content)
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func HandleDeleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := mw.GetUserID(r)

	if err := db.DeleteNote(id, userID); err != nil {
		if db.IsNotFound(err) {
			jsonError(w, "note not found or not yours", http.StatusNotFound)
			return
		}
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
