package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/api"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/managers/detectionmgr"
	"github.com/google/uuid"
)

type createDetectionRequest struct {
	Name            string                       `json:"name"`
	Type            string                       `json:"type"`
	Pattern         string                       `json:"pattern"`
	Actions         []detectionmgr.RoutineAction `json:"actions"`
	CooldownSeconds int                          `json:"cooldownSeconds"`
}

func ListDetections(w http.ResponseWriter, _ *http.Request) {
	api.WriteData(w, http.StatusOK, map[string]any{"detections": detectionmgr.GetCustomDetections()})
}

func CreateDetection(w http.ResponseWriter, r *http.Request) {
	var request createDetectionRequest
	if err := api.DecodeJSON(w, r, &request); err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	request.Type = strings.ToLower(strings.TrimSpace(request.Type))
	detection := detectionmgr.CustomDetection{
		ID:              uuid.NewString(),
		Name:            strings.TrimSpace(request.Name),
		Type:            request.Type,
		Pattern:         strings.TrimSpace(request.Pattern),
		Actions:         request.Actions,
		CooldownSeconds: request.CooldownSeconds,
	}
	for i := range detection.Actions {
		detection.Actions[i].Type = strings.TrimSpace(detection.Actions[i].Type)
		detection.Actions[i].Value = strings.TrimSpace(detection.Actions[i].Value)
		detection.Actions[i].Text = strings.TrimSpace(detection.Actions[i].Text)
	}
	if err := detectionmgr.ValidateDetection(detection); err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid_detection", err.Error())
		return
	}
	if err := detectionmgr.AddCustomDetection(detection); err != nil {
		api.WriteError(w, http.StatusInternalServerError, "detection_save_failed", "Could not save the detection")
		return
	}
	api.WriteData(w, http.StatusCreated, map[string]any{"detection": detection})
}

func DeleteDetection(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		api.WriteError(w, http.StatusBadRequest, "missing_detection_id", "Detection ID is required")
		return
	}
	if err := detectionmgr.RemoveCustomDetection(id); err != nil {
		if errors.Is(err, detectionmgr.ErrDetectionNotFound) {
			api.WriteError(w, http.StatusNotFound, "detection_not_found", "Detection not found")
			return
		}
		api.WriteError(w, http.StatusInternalServerError, "detection_delete_failed", "Could not delete the detection")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
