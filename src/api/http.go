package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/localization"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/commandmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamcmd"
)

// StartServer HTTP handler
func StartServer(w http.ResponseWriter, r *http.Request) {
	logger.API.Debug("Received start request from API")
	if err := gamemgr.InternalStartServer(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.API.Error("Error starting server: " + err.Error())
		return
	}
	fmt.Fprint(w, localization.GetString("BackendText_ServerStarted"))
	logger.API.Info("Server started.")
}

// StopServer HTTP handler
func StopServer(w http.ResponseWriter, r *http.Request) {
	logger.API.Debug("Received stop request from API")
	if err := gamemgr.InternalStopServer(); err != nil {
		if err.Error() == "server not running" {
			fmt.Fprint(w, localization.GetString("BackendText_ServerNotRunningOrAlreadyStopped"))
			logger.API.Warn("Server not running or was already stopped")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.API.Error("Error stopping server: " + err.Error())
		return
	}
	detectionmgr.ClearPlayers(detectionmgr.GetDetector())
	fmt.Fprint(w, localization.GetString("BackendText_ServerStopped"))
	logger.API.Info("Server stopped.")
}

func GetGameServerRunState(w http.ResponseWriter, r *http.Request) {
	runState := gamemgr.InternalIsServerRunning()
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

// CommandHandler handles POST requests to execute commands via commandmgr.
// Expects a command in the request body. Returns 204 on success or error details.
func CommandHandler(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read command from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	command := strings.TrimSpace(string(body))

	// Validate command
	if command == "" {
		http.Error(w, "Command cannot be empty", http.StatusBadRequest)
		return
	}

	// Execute command via commandmgr
	if err := commandmgr.WriteCommand(command); err != nil {
		switch err {
		case os.ErrNotExist:
			http.Error(w, "Command file path not configured", http.StatusInternalServerError)
		case os.ErrInvalid:
			http.Error(w, "Invalid command", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to write command: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Success: return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func HandleIsSSCMEnabled(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if SSCM is enabled
	if !config.GetIsSSCMEnabled() {
		http.Error(w, "SSCM is disabled", http.StatusForbidden)
		return
	}

	// Success: return 200 OK
	w.WriteHeader(http.StatusOK)
}

// run SteamCMD from API, but only allow once every 5 minutes to "kinda" prevent concurrent executions although that woluldnt hurn.
// If the user has a 5mbit connection, I cannot help them anyways.
func HandleRunSteamCMD(w http.ResponseWriter, r *http.Request) {

	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.Core.Info("Running SteamCMD")
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
