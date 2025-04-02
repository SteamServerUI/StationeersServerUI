// interface.go
package backupsv2

import (
	"StationeersServerUI/src/config"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var manager *BackupManager

func InitBackupManager() {

	backupDir := filepath.Join("./saves/", config.WorldName, "backup")
	safeBackupDir := filepath.Join("./saves/", config.WorldName, "Safebackups")
	backupManager := NewBackupManager(BackupConfig{
		WorldName:     config.WorldName,
		BackupDir:     backupDir,
		SafeBackupDir: safeBackupDir,
		WaitTime:      30 * time.Second,
	})
	err := backupManager.Start()
	if err != nil {
		fmt.Println("[BACKUPS] Failed to start backups manager: " + err.Error())
	}
}

// HTTP handler for listing backups
func ListBackups(w http.ResponseWriter, r *http.Request) {
	backups, err := manager.ListBackups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(backups)
}

// HTTP handler for restoring backups
func RestoreBackup(w http.ResponseWriter, r *http.Request) {
	indexStr := r.URL.Query().Get("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid backup index", http.StatusBadRequest)
		return
	}

	if err := manager.RestoreBackup(index); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Backup restored successfully"))
}
