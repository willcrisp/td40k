package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/willcrisp/blueprint/internal/db"
)

func HandleGetCounter(w http.ResponseWriter, r *http.Request) {
	counter, err := db.GetCounter()
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counter)
}

func HandleIncrementCounter(w http.ResponseWriter, r *http.Request) {
	counter, err := db.IncrementCounter()
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counter)
}
