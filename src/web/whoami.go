package web

import (
	"encoding/json"
	"net/http"
)

func WhoAmIHandler(w http.ResponseWriter, r *http.Request) {

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"username": "SSUI", "accessLevel": "SSUI-Admin"}); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
