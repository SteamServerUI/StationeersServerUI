package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
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
	steamIDsStr := strings.Join(steamIDs, ",")

	jxsnDevApiKey := config.JXSNDevApiKey
	apiURL := fmt.Sprintf("https://jxsn.dev/api/v1/steamapi/userinfo?steamids=%s&accessval=%s", steamIDsStr, jxsnDevApiKey)
	logger.Web.Debug("Fetching player details from external API: " + apiURL) // REMOVE ME
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch player details", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode the API response
	var apiResponse struct {
		Response struct {
			Players []struct {
				SteamID      string `json:"steamid"`
				PersonaName  string `json:"personaname"`
				ProfileURL   string `json:"profileurl"`
				AvatarMedium string `json:"avatarmedium"`
			} `json:"players"`
		} `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		http.Error(w, "Failed to decode API response", http.StatusInternalServerError)
		return
	}

	// Create a map to store API player details by SteamID
	apiPlayers := make(map[string]struct {
		PersonaName  string
		ProfileURL   string
		AvatarMedium string
	})
	for _, player := range apiResponse.Response.Players {
		apiPlayers[player.SteamID] = struct {
			PersonaName  string
			ProfileURL   string
			AvatarMedium string
		}{
			PersonaName:  player.PersonaName,
			ProfileURL:   player.ProfileURL,
			AvatarMedium: player.AvatarMedium,
		}
	}

	// Build the response player list
	playerList := make([]map[string]map[string]string, 0, len(players))
	for steamID, username := range players {
		playerInfo := map[string]string{
			"username": username,
			"steamID":  steamID,
		}
		// Add API details if available
		if apiPlayer, exists := apiPlayers[steamID]; exists {
			playerInfo["personaname"] = apiPlayer.PersonaName
			playerInfo["profileurl"] = apiPlayer.ProfileURL
			playerInfo["avatarmedium"] = apiPlayer.AvatarMedium
			playerInfo["steamid"] = steamID
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
