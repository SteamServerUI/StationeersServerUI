package api

import (
	"encoding/json"
	"net/http"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/steamserverui/systeminfo"
)

func HandleGetOsStats(w http.ResponseWriter, r *http.Request) {
	logger.API.Debug("Received getOsStats request from API")
	// accept only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get cached stats or refresh if needed
	stats, err := systeminfo.RefreshCachedStats()
	if err != nil {
		logger.API.Error("Failed to get OS stats")
		http.Error(w, "Failed to get OS stats", http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		logger.API.Error("Failed to write response")
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
