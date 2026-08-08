package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/backupmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/systeminfo"
)

type overviewPlayer struct {
	SteamID string `json:"steamId"`
	Name    string `json:"name"`
}

func HandleOverview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeV3Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use GET for the overview")
		return
	}
	stats, err := systeminfo.RefreshCachedStats()
	if err != nil {
		writeV3Error(w, http.StatusInternalServerError, "system_stats_unavailable", "System statistics are unavailable")
		return
	}
	connected := detectionmgr.GetPlayers(detectionmgr.GetDetector())
	players := make([]overviewPlayer, 0, len(connected))
	for steamID, name := range connected {
		players = append(players, overviewPlayer{SteamID: steamID, Name: name})
	}
	data := map[string]any{
		"capturedAt": time.Now().UTC(),
		"server": map[string]any{
			"isRunning": gamemgr.InternalIsServerRunning(),
			"id":        gamemgr.GameServerUUID.String(),
		},
		"system":  stats,
		"players": players,
		"backup": map[string]any{
			"systemReady":   backupmgr.IsSystemReady(),
			"isRunning":     backupmgr.IsBackupRunning(),
			"isLoopRunning": backupmgr.IsLoopRunning(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
}

func writeV3Error(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": code, "message": message}})
}
