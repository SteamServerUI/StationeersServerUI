// backupinterface.go
package backupsv2

import (
	"StationeersServerUI/src/config"
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

	return BackupConfig{
		WorldName:     config.WorldName,
		BackupDir:     config.ConfiguredBackupDir,
		SafeBackupDir: config.ConfiguredSafeBackupDir,
		WaitTime:      30 * time.Second,
		RetentionPolicy: RetentionPolicy{
			KeepLastN:       config.BackupKeepLastN,
			KeepDailyFor:    config.BackupKeepDailyFor,
			KeepWeeklyFor:   config.BackupKeepWeeklyFor,
			KeepMonthlyFor:  config.BackupKeepMonthlyFor,
			CleanupInterval: config.BackupKeepMonthlyFor,
		},
	}
}
