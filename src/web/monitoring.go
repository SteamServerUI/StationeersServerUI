package web

import (
	"encoding/json"
	"net/http"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

func HandleMonitorStatus(w http.ResponseWriter, r *http.Request) {
	runState := config.GetIsGameServerRunning()
	response := map[string]interface{}{
		"isRunning": runState,
	}
	w.Header().Set("Content-Type", "application/json")

	if !runState {
		w.WriteHeader(http.StatusServiceUnavailable) // 503
	} else {
		w.WriteHeader(http.StatusOK) // 200
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to respond with Game Server status", http.StatusInternalServerError)
		return
	}
}
