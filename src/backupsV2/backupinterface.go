// backupinterface.go
package backupsv2

import (
	"StationeersServerUI/src/config"
	"path/filepath"
	"time"
)

// GetDefaultConfig returns a properly configured BackupConfig
func GetDefaultConfig() BackupConfig {
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
