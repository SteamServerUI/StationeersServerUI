// backupinterface.go
package backupmgr

import (
	"StationeersServerUI/src/config"
	"time"
)

// GlobalBackupManager is the singleton instance of the backup manager
var GlobalBackupManager *BackupManager

// InitGlobalBackupManager initializes the global backup manager instance
func InitGlobalBackupManager(config BackupConfig) error {

	if GlobalBackupManager != nil {
		GlobalBackupManager.Shutdown()
	}

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

// ReloadBackupManagerFromConfig reloads the global backup manager with the current config. This should be called whenever the config is changed.
func ReloadBackupManagerFromConfig() error {
	// Create a new backupManager config from the global config
	backupConfig := GetBackupConfig()

	// Reinitialize the global backup manager with the new config
	return InitGlobalBackupManager(backupConfig)
}
