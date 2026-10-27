package web

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/api"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/diagnostics"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/logger"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/supportlink"
)

func GetSupportDiagnostics(w http.ResponseWriter, r *http.Request) {
	api.WriteData(w, http.StatusOK, diagnostics.Run())
}

func DownloadSupportPackage(w http.ResponseWriter, r *http.Request) {
	streamSupportZip(w, "ssui-support", diagnostics.WritePackage)
}

func DownloadSupportRecording(w http.ResponseWriter, r *http.Request) {
	recording := logger.SupportRecordingStatus()
	if recording.Active || recording.StoppedAt == nil {
		api.WriteError(w, http.StatusConflict, "support_recording_unavailable", "Stop a recording before downloading it")
		return
	}
	streamSupportZip(w, "ssui-support-recording", diagnostics.WriteRecordingPackage)
}

func GetSupportRecording(w http.ResponseWriter, r *http.Request) {
	api.WriteData(w, http.StatusOK, logger.SupportRecordingStatus())
}

func StartSupportRecording(w http.ResponseWriter, r *http.Request) {
	recording, err := logger.StartSupportRecording()
	if errors.Is(err, logger.ErrSupportRecordingActive) {
		api.WriteError(w, http.StatusConflict, "support_recording_active", "A support recording is already running")
		return
	}
	if err != nil {
		logger.Web.Error("Could not start support recording: " + err.Error())
		api.WriteError(w, http.StatusInternalServerError, "support_recording_failed", "Could not start the support recording")
		return
	}
	api.WriteData(w, http.StatusOK, recording)
}

func StopSupportRecording(w http.ResponseWriter, r *http.Request) {
	recording, err := diagnostics.StopRecording()
	if errors.Is(err, logger.ErrNoSupportRecording) {
		api.WriteError(w, http.StatusConflict, "support_recording_inactive", "No support recording is running")
		return
	}
	if err != nil {
		logger.Web.Warn("Support recording stopped with an error: " + err.Error())
	}
	api.WriteData(w, http.StatusOK, recording)
}

func streamSupportZip(w http.ResponseWriter, name string, write func(io.Writer) error) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s-%s.zip"`, name, time.Now().Format("20060102-150405")))
	if err := write(w); err != nil {
		logger.Web.Error("Could not write support package: " + err.Error())
	}
}

func GetSupportConnection(w http.ResponseWriter, r *http.Request) {
	api.WriteData(w, http.StatusOK, supportlink.Snapshot())
}

func OpenSupportConnection(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Minutes int `json:"minutes"`
	}
	if err := api.DecodeJSON(w, r, &body); err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	status, err := supportlink.Start(time.Duration(body.Minutes) * time.Minute)
	switch {
	case errors.Is(err, supportlink.ErrActive):
		api.WriteError(w, http.StatusConflict, "support_connection_active", "A support connection is already open")
	case errors.Is(err, supportlink.ErrInvalidDuration):
		api.WriteError(w, http.StatusBadRequest, "invalid_duration", err.Error())
	case err != nil:
		logger.Web.Error("Could not open support connection: " + err.Error())
		api.WriteError(w, http.StatusInternalServerError, "support_connection_failed", "Could not open the support connection")
	default:
		api.WriteData(w, http.StatusOK, status)
	}
}

func CloseSupportConnection(w http.ResponseWriter, r *http.Request) {
	api.WriteData(w, http.StatusOK, supportlink.Stop())
}
