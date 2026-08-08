package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/middleware"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/activity"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
)

type bufferedResponse struct {
	header      http.Header
	body        bytes.Buffer
	status      int
	wroteHeader bool
}

func newBufferedResponse() *bufferedResponse {
	return &bufferedResponse{header: make(http.Header), status: http.StatusOK}
}

func (response *bufferedResponse) Header() http.Header { return response.header }
func (response *bufferedResponse) WriteHeader(status int) {
	if response.wroteHeader {
		return
	}
	response.status = status
	response.wroteHeader = true
}
func (response *bufferedResponse) Write(data []byte) (int, error) { return response.body.Write(data) }

func v3JSONBoundary(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buffered := newBufferedResponse()
		next(buffered, r)
		for key, values := range buffered.header {
			if strings.EqualFold(key, "Content-Length") {
				continue
			}
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		body := buffered.body.Bytes()
		contentType := buffered.header.Get("Content-Type")
		if buffered.status == http.StatusNoContent {
			w.WriteHeader(buffered.status)
		} else if strings.Contains(contentType, "application/json") && json.Valid(body) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(buffered.status)
			_, _ = w.Write(body)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(buffered.status)
			message := strings.TrimSpace(string(body))
			if buffered.status >= 400 {
				_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": "request_failed", "message": message}})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]string{"message": message}})
			}
		}
		if changesState(r.Method) && buffered.status >= 200 && buffered.status < 400 {
			recordRequestActivity(r)
		}
	}
}

func HandleCapabilities(w http.ResponseWriter, r *http.Request) {
	principal, _ := middleware.PrincipalFromContext(r.Context())
	permissions := make([]string, 0, len(principal.Permissions))
	for _, permission := range security.AllPermissions {
		if principal.Permissions[permission] {
			permissions = append(permissions, permission)
		}
	}
	data := map[string]any{
		"apiVersion": "v3",
		"backend": map[string]any{
			"version":             config.GetVersion(),
			"oneServerPerBackend": true,
		},
		"features":     map[string]bool{"plugins": config.GetPluginsEnabled()},
		"capabilities": []string{"activity.history", "activity.stream", "auth.groups", "auth.sessions", "auth.tokens", "desktop.remote", "settings.typed"},
		"permissions":  permissions,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
}

func recordRequestActivity(r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		return
	}
	activity.Record(activity.Event{
		Type:      "operation",
		Tone:      "success",
		Message:   activityLabel(r.Method, r.URL.Path),
		Detail:    r.Method + " " + r.URL.Path,
		ActorID:   principal.UserID,
		ActorName: principal.Username,
	})
}

func activityLabel(method, path string) string {
	labels := map[string]string{
		"/api/v3/server/start":         "Server start requested",
		"/api/v3/server/stop":          "Server stop requested",
		"/api/v3/server/restart":       "Server restart requested",
		"/api/v3/backup/create":        "Backup creation accepted",
		"/api/v3/backup/restore":       "Backup restore requested",
		"/api/v3/settings":             "Setting updated",
		"/api/v3/files/content":        "File saved",
		"/api/v3/runfile/args/update":  "Game setting updated",
		"/api/v3/runfile/save":         "Game configuration saved",
		"/api/v3/gallery/select":       "Game definition installed",
		"/api/v3/plugingallery/select": "Plugin installed",
		"/api/v3/loader/reloadbackend": "Backend reload requested",
		"/api/v3/loader/reloadrunfile": "Game configuration reloaded",
	}
	if label := labels[path]; label != "" {
		return label
	}
	return method + " request completed"
}

func changesState(method string) bool {
	return method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions
}
