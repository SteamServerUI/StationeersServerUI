package backupapi

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/backupmgr"
)

type BackupItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	SizeBytes int64     `json:"sizeBytes"`
	SizeKnown bool      `json:"sizeKnown"`
	Format    string    `json:"format"`
}

type RestoreRequest struct {
	BackupName    string `json:"backupName"`
	SkipPreBackup bool   `json:"skipPreBackup"`
}

type BackupCreateRequest struct {
	Mode string `json:"mode"`
}

func HandleBackupCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		backupError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only POST requests are allowed")
		return
	}
	var request BackupCreateRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&request); err != nil {
		backupError(w, http.StatusBadRequest, "invalid_json", "A backup mode is required")
		return
	}
	if request.Mode != "copy" && request.Mode != "tar" {
		backupError(w, http.StatusBadRequest, "invalid_backup_mode", "Backup mode must be copy or tar")
		return
	}
	if err := backupmgr.CreateBackup(request.Mode); err != nil {
		logger.Web.Error("API: Failed to create backup: " + err.Error())
		backupError(w, http.StatusConflict, "backup_failed", err.Error())
		return
	}
	writeBackupJSON(w, http.StatusAccepted, map[string]any{"data": map[string]any{"accepted": true, "mode": request.Mode}})
}

func HandleBackupList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		backupError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET requests are allowed")
		return
	}
	names, err := backupmgr.GetBackupList()
	if err != nil {
		backupError(w, http.StatusInternalServerError, "backup_list_failed", "The restore points could not be listed")
		return
	}
	backups := make([]BackupItem, 0, len(names))
	for _, name := range names {
		item := BackupItem{ID: name, Name: name, Format: backupFormat(name)}
		if info, statErr := os.Stat(filepath.Join(config.GetBackupsStoreDir(), name)); statErr == nil {
			item.CreatedAt = info.ModTime().UTC()
			if !info.IsDir() {
				item.SizeBytes = info.Size()
				item.SizeKnown = true
			}
		}
		backups = append(backups, item)
	}
	writeBackupJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"backups": backups}})
}

func HandleBackupRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		backupError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only POST requests are allowed")
		return
	}
	var request RestoreRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&request); err != nil {
		backupError(w, http.StatusBadRequest, "invalid_json", "A restore point is required")
		return
	}
	request.BackupName = strings.TrimSpace(request.BackupName)
	if request.BackupName == "" || strings.Contains(request.BackupName, "..") || strings.ContainsAny(request.BackupName, "/\\") {
		backupError(w, http.StatusBadRequest, "invalid_backup_name", "The restore point name is invalid")
		return
	}
	if err := backupmgr.RestoreBackup(request.BackupName, request.SkipPreBackup); err != nil {
		logger.Web.Error("API: Failed to restore backup: " + err.Error())
		backupError(w, http.StatusConflict, "restore_failed", err.Error())
		return
	}
	writeBackupJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"restored": true, "backupName": request.BackupName}})
}

func HandleBackupStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		backupError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET requests are allowed")
		return
	}
	writeBackupJSON(w, http.StatusOK, map[string]any{"data": map[string]any{
		"systemReady":   backupmgr.IsSystemReady(),
		"isLoopRunning": backupmgr.IsLoopRunning(),
		"isRunning":     backupmgr.IsBackupRunning(),
	}})
}

func backupFormat(name string) string {
	switch {
	case strings.HasSuffix(name, ".tar.gz"):
		return "tar.gz"
	case strings.HasSuffix(name, ".tar"):
		return "tar"
	case strings.HasPrefix(name, "snapshot_"):
		return "snapshot"
	default:
		return "copy"
	}
}

func backupError(w http.ResponseWriter, status int, code, message string) {
	writeBackupJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func writeBackupJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
