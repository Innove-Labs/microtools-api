package handlers

import (
	"encoding/json"
	"net/http"
)

// LiveHandler handles health check requests
func LiveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Live"})
}
