// backupinterface.go
package backupsv2

import (
	"StationeersServerUI/src/config"
	"path/filepath"
	"time"
)

// GlobalBackupManager is the singleton instance of the backup manager
var GlobalBackupManager *BackupManager

// InitGlobalBackupManager initializes the global backup manager instance
func InitGlobalBackupManager(config BackupConfig) error {
	GlobalBackupManager = NewBackupManager(config)
	if err := GlobalBackupManager.Initialize(); err != nil {
		return err
	}
	return GlobalBackupManager.Start()
}

// GetBackupConfig returns a properly configured BackupConfig
func GetBackupConfig() BackupConfig {
	backupDir := filepath.Join("./saves/" + config.WorldName + "/Backup")
	safeBackupDir := filepath.Join("./saves/" + config.WorldName + "/Safebackups")

	return BackupConfig{
		WorldName:     config.WorldName,
		BackupDir:     backupDir,
		SafeBackupDir: safeBackupDir,
		WaitTime:      30 * time.Second,
		RetentionPolicy: RetentionPolicy{
			KeepLastN:       10,
			KeepDailyFor:    30 * 24 * time.Hour,
			KeepWeeklyFor:   90 * 24 * time.Hour,
			KeepMonthlyFor:  365 * 24 * time.Hour,
			CleanupInterval: 1 * time.Hour,
		},
	}
}
