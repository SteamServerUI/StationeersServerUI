package legacyapi

import (
	"encoding/json"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"
)

func StartServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeActionError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use POST to start the server")
		return
	}
	if err := gamemgr.InternalStartServer(); err != nil {
		logger.API.Error("Error starting server: " + err.Error())
		writeActionError(w, http.StatusConflict, "server_start_failed", err.Error())
		return
	}
	logger.API.Info("Server started.")
	writeAction(w, http.StatusAccepted, "start", "Server started")
}

func StopServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeActionError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use POST to stop the server")
		return
	}
	if err := gamemgr.InternalStopServer(); err != nil {
		if err.Error() == "server not running" {
			writeAction(w, http.StatusOK, "stop", "Server was already stopped")
			return
		}
		logger.API.Error("Error stopping server: " + err.Error())
		writeActionError(w, http.StatusConflict, "server_stop_failed", err.Error())
		return
	}
	detectionmgr.ClearPlayers(detectionmgr.GetDetector())
	logger.API.Info("Server stopped.")
	writeAction(w, http.StatusAccepted, "stop", "Server stopped")
}

func writeAction(w http.ResponseWriter, status int, action, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{
		"action": action, "accepted": true, "message": message,
	}})
}

func writeActionError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": code, "message": message}})
}
