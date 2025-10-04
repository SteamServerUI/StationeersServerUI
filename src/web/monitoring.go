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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200 OK
		http.Error(w, "Failed to respond with Game Server status", http.StatusInternalServerError)
		return
	}
}
