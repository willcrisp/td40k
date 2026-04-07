package counter

import (
	"encoding/json"
	"net/http"
)

// handleGet utilizes the Repository database models entirely 
func (m *Module) handleGet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "main"
	}

	resp, err := m.Repo.GetCounter(r.Context(), name)
	if err != nil {
		http.Error(w, "Failed to fetch counter data", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleIncrement routes raw internet traffic securely to the DB upsert payload handlers
func (m *Module) handleIncrement(w http.ResponseWriter, r *http.Request) {
	var req IncrementRequest
	// Silently tolerate parse fails
	json.NewDecoder(r.Body).Decode(&req)
	
	if req.Name == "" {
		req.Name = "main"
	}

	resp, err := m.Repo.IncrementCounter(r.Context(), req.Name)
	if err != nil {
		http.Error(w, "Failed to increment", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
