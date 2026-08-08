package runfileapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/gallery"
)

// response wraps API responses
type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// GalleryHandler handles GET /api/v3/gallery
func GalleryHandler(w http.ResponseWriter, r *http.Request) {
	forceUpdate := strings.ToLower(r.URL.Query().Get("forceUpdate")) == "true"

	runfiles, err := gallery.GetRunfileGallery(forceUpdate)
	if err != nil {
		logger.Runfile.Error("Gallery fetch failed: " + err.Error())
		sendResponse(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	logger.Runfile.Debug("Returning " + strconv.Itoa(len(runfiles)) + " runfiles from gallery")
	sendResponse(w, http.StatusOK, response{Data: runfiles})
}

// GallerySelectHandler handles POST /api/v3/gallery/select
func GallerySelectHandler(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Identifier string `json:"identifier"`
		Redownload bool   `json:"redownload,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Runfile.Error("Invalid request body: " + err.Error())
		sendResponse(w, http.StatusBadRequest, response{Error: "invalid JSON, check your request"})
		return
	}

	if req.Identifier == "" {
		logger.Runfile.Error("Missing identifier in request")
		sendResponse(w, http.StatusBadRequest, response{Error: "identifier is required"})
		return
	}

	redownload := false
	if req.Redownload {
		redownload = true
	}

	if err := gallery.SaveRunfileToDisk(req.Identifier, redownload); err != nil {
		logger.Runfile.Error("Failed to save runfile " + req.Identifier + ": " + err.Error())
		sendResponse(w, http.StatusConflict, response{Error: err.Error()})
		return
	}

	logger.Runfile.Debug("Successfully saved runfile " + req.Identifier)
	loader.ReloadBackend()
	sendResponse(w, http.StatusOK, response{Data: "Runfile " + req.Identifier + " saved"})
}

// sendResponse writes a JSON response with the given status code
func sendResponse(w http.ResponseWriter, status int, resp response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Runfile.Error("Failed to encode response: " + err.Error())
	}
}
