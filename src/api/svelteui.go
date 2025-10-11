package api

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

var reloadMu sync.Mutex

func ServeSvelteUI(w http.ResponseWriter, r *http.Request) {
	htmlFS, err := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/v2")
	if err != nil {
		http.Error(w, "Error accessing Svelte UI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlFile, err := htmlFS.Open("index.html")
	if err != nil {
		http.Error(w, "Error reading Svelte UI: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer htmlFile.Close()

	// Stream the file content to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = io.Copy(w, htmlFile)
	if err != nil {
		http.Error(w, "Error writing Svelte UI: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleReloadAll(w http.ResponseWriter, r *http.Request) {
	logger.API.Debug("Received reloadbackend request from API")
	reloadMu.Lock()
	defer reloadMu.Unlock()
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
