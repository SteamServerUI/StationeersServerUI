package backupsv2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// HTTPHandler provides HTTP endpoints for backup operations
type HTTPHandler struct {
	manager *BackupManager
}

// NewHTTPHandler creates a new HTTP handler for backups
func NewHTTPHandler(manager *BackupManager) *HTTPHandler {
	return &HTTPHandler{manager: manager}
}

// ListBackupsHandler handles requests to list available backups
func (h *HTTPHandler) ListBackupsHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	var limit int
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	backups, err := h.manager.ListBackups(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if classic mode is requested
	mode := r.URL.Query().Get("mode")
	if mode == "classic" {
		// Format the response in the classic format
		classicResponses := make([]string, 0, len(backups))
		for _, backup := range backups {
			// Format according to classic view: "BackupIndex: X, Created: DD.MM.YYYY HH:MM:SS"
			classicLine := fmt.Sprintf("BackupIndex: %d, Created: %s",
				backup.Index,
				backup.ModTime.Format("02.01.2006 15:04:05"))
			classicResponses = append(classicResponses, classicLine)
		}

		// Return plain text response for classic mode
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strings.Join(classicResponses, "\n")))
		return
	}

	// Default JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(backups)
}

// RestoreBackupHandler handles requests to restore a backup
func (h *HTTPHandler) RestoreBackupHandler(w http.ResponseWriter, r *http.Request) {
	indexStr := r.URL.Query().Get("index")
	if indexStr == "" {
		http.Error(w, "index parameter is required", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "invalid index parameter", http.StatusBadRequest)
		return
	}

	if err := h.manager.RestoreBackup(index); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Backup restored successfully"))
}
