package backupsv2

import (
	"encoding/json"
	"net/http"
	"sort"
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
	backups, err := h.manager.ListBackups()
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

// ListBackups returns information about available backups
func (m *BackupManager) ListBackups() ([]BackupGroup, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	groups, err := m.getBackupGroups()
	if err != nil {
		return nil, err
	}

	// Sort by index (newest first)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Index > groups[j].Index
	})

	return groups, nil
}
