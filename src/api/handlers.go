package api

import (
	"encoding/json"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamcmd"
)

func GetGameServerRunState(w http.ResponseWriter, r *http.Request) {
	runState := config.GetIsGameServerRunning()
	response := map[string]interface{}{
		"isRunning": runState,
		"uuid":      gamemgr.GameServerUUID.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to respond with Game Server status", http.StatusInternalServerError)
		return
	}
}

func HandleRunSteamCMD(w http.ResponseWriter, r *http.Request) {

	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	_, err := steamcmd.InstallAndRunSteamCMD()

	// Update last execution time

	// Success: return 202 Accepted and JSON
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		json.NewEncoder(w).Encode(map[string]string{"statuscode": "202", "status": "Success", "message": "SteamCMD ran successfully, gameserver files are up-to-date!"})
		return
	}
	// Failure: return 202 Accepted and JSON with the error message
	json.NewEncoder(w).Encode(map[string]string{"statuscode": "202", "status": "Failed", "message": "SteamCMD ran unsuccessfully:" + err.Error()})
}

func HandleReloadBackend(w http.ResponseWriter, r *http.Request) {
	logger.API.Debug("Received reloadbackend request from API")
	// accept only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	// Reload all loaders
	loader.ReloadBackend()

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"}); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
