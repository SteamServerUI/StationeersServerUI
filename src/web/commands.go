package web

import (
	"encoding/json"
	"net/http"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/commandmgr"
)

type CommandRequest struct {
	Command string `json:"command"`
}

// CommandResponse represents the JSON response structure.
type CommandResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// CommandHandler handles POST requests to execute commands via commandmgr.
func HandleCommand(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}

	if !config.IsSSCMEnabled {
		sendErrorResponse(w, http.StatusForbidden, "SSCM is disabled, cannot execute commands")
		return
	}

	// Parse JSON request body
	var req CommandRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validate command
	if req.Command == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Command cannot be empty")
		return
	}

	// Write command using commandmgr
	if err := commandmgr.WriteCommand(req.Command); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to write command: "+err.Error())
		return
	}

	// Send success response
	sendSuccessResponse(w, "Command passed to server")
}

// sendErrorResponse sends a JSON error response with the given status code and message.
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := CommandResponse{
		Status:  "error",
		Message: message,
	}
	json.NewEncoder(w).Encode(resp)
}

// sendSuccessResponse sends a JSON success response.
func sendSuccessResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := CommandResponse{
		Status:  "success",
		Message: message,
	}
	json.NewEncoder(w).Encode(resp)
}
