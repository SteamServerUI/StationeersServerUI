package backupsv2

import (
	"StationeersServerUI/src/config"
	"fmt"
	"path/filepath"
	"time"
)

// Global manager instance (consider making this non-global in future)
var manager *BackupManager

// InitBackupManager initializes the backup system
func InitBackupManager() error {
	backupDir := filepath.Join("./saves/", config.WorldName, "backup")
	safeBackupDir := filepath.Join("./saves/", config.WorldName, "Safebackups")

	manager = NewBackupManager(BackupConfig{
		WorldName:     config.WorldName,
		BackupDir:     backupDir,
		SafeBackupDir: safeBackupDir,
		WaitTime:      30 * time.Second,
	})

	return manager.Start()
}

// GetBackups returns a list of available backups
// limit: number of recent backups to return (0 for all)
func GetBackups(limit int) ([]BackupGroup, error) {
	if manager == nil {
		return nil, fmt.Errorf("backup manager not initialized")
	}

	backups, err := manager.ListBackups(limit)
	if err != nil {
		return nil, err
	}

	if limit > 0 && limit < len(backups) {
		backups = backups[:limit]
	}

	return backups, nil
}

// RestoreBackupIndex restores a specific backup
func RestoreBackupIndex(index int) error {
	if manager == nil {
		return fmt.Errorf("backup manager not initialized")
	}
	return manager.RestoreBackup(index)
}
