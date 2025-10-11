package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/gallery"
)

// response wraps API responses
type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// galleryHandler handles GET /api/v2/gallery
func galleryHandler(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Info("Handling GET /api/v2/gallery request")
	forceUpdate := strings.ToLower(r.URL.Query().Get("forceUpdate")) == "true"

	runfiles, err := gallery.GetRunfileGallery(forceUpdate)
	if err != nil {
		logger.Runfile.Error("Gallery fetch failed: " + err.Error())
		sendResponse(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	logger.Runfile.Info("Returning " + strconv.Itoa(len(runfiles)) + " runfiles from gallery")
	sendResponse(w, http.StatusOK, response{Data: runfiles})
}

// selectHandler handles POST /api/v2/gallery/select
func selectHandler(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Info("Handling POST /api/v2/gallery/select request")

	var req struct {
		Identifier string `json:"identifier"`
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

	if err := gallery.SaveRunfileToDisk(req.Identifier); err != nil {
		logger.Runfile.Error("Failed to save runfile " + req.Identifier + ": " + err.Error())
		sendResponse(w, http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	logger.Runfile.Info("Successfully saved runfile " + req.Identifier)
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
