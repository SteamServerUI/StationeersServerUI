package pluginsapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/gallery"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/plugins"
)

type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// PluginGalleryHandler handles GET /api/v3/plugins
func PluginGalleryHandler(w http.ResponseWriter, r *http.Request) {
	forceUpdate := strings.ToLower(r.URL.Query().Get("forceUpdate")) == "true"

	plugins, err := gallery.GetPluginGallery(forceUpdate)
	if err != nil {
		logger.Plugin.Error("Plugin gallery fetch failed: " + err.Error())
		sendResponse(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	logger.Plugin.Debug("Returning " + strconv.Itoa(len(plugins)) + " plugins from gallery")
	sendResponse(w, http.StatusOK, response{Data: plugins})
}

// PluginSelectHandler handles POST /api/v3/plugins/select
func PluginSelectHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string `json:"name"`
		Redownload bool   `json:"redownload,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Plugin.Error("Invalid request body: " + err.Error())
		sendResponse(w, http.StatusBadRequest, response{Error: "invalid JSON, check your request"})
		return
	}

	if req.Name == "" {
		logger.Plugin.Error("Missing name in request")
		sendResponse(w, http.StatusBadRequest, response{Error: "name is required"})
		return
	}

	redownload := false
	if req.Redownload {
		plugins.StopPlugin(req.Name)
		redownload = true
	}

	if err := gallery.SavePluginToDisk(req.Name, redownload); err != nil {
		logger.Plugin.Error("Failed to save plugin " + req.Name + ": " + err.Error())
		sendResponse(w, http.StatusConflict, response{Error: err.Error()})
		return
	}

	logger.Plugin.Debug("Successfully saved plugin " + req.Name)
	sendResponse(w, http.StatusOK, response{Data: "Plugin " + req.Name + " saved"})
}

// sendResponse writes a JSON response with the given status code
func sendResponse(w http.ResponseWriter, status int, resp response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Runfile.Error("Failed to encode response: " + err.Error())
	}
}
