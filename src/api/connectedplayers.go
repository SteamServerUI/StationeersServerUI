package api

import (
	"encoding/json"
	"net/http"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/detectionmgr"
)

// PrintConnectedPlayersHandler handles HTTP requests to list connected players.
func HandleConnectedPlayersList(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	detector := detectionmgr.GetDetector()
	players := detectionmgr.GetPlayers(detector)

	// if players is empty, return an empty list
	if len(players) == 0 {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(players); err != nil {
			http.Error(w, "Failed to encode player list", http.StatusInternalServerError)
		}
		return
	}

	// Create a comma-separated string of SteamIDs
	steamIDs := make([]string, 0, len(players))
	for steamID := range players {
		steamIDs = append(steamIDs, steamID)
	}

	// Build the response player list
	playerList := make([]map[string]map[string]string, 0, len(players))
	for steamID, username := range players {
		playerInfo := map[string]string{
			"username": username,
			"steamID":  steamID,
		}
		nestedPlayer := map[string]map[string]string{
			username: playerInfo,
		}
		playerList = append(playerList, nestedPlayer)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(playerList); err != nil {
		http.Error(w, "Failed to encode player list", http.StatusInternalServerError)
	}
}
