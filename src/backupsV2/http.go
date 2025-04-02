package backupsv2

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// HTTP handler for restoring backups
func RestoreBackup(w http.ResponseWriter, r *http.Request) {
	indexStr := r.URL.Query().Get("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid backup index", http.StatusBadRequest)
		return
	}

	if err := RestoreBackupIndex(index); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Backup restored successfully"))
}

// HTTP handler for listing backups
func ListBackups(w http.ResponseWriter, r *http.Request) {
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

	backups, err := GetBackups(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(backups)
}
