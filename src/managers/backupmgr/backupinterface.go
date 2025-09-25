// backupinterface.go
package backupmgr

import (
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// GlobalBackupManager is the singleton instance of the backup manager
var GlobalBackupManager *BackupManager

// Track all HTTP handlers that need updating when manager changes
var activeHTTPHandlers []*HTTPHandler

// InitGlobalBackupManager initializes the global backup manager instance
func InitGlobalBackupManager(config BackupConfig) error {
	if GlobalBackupManager != nil {
		logger.Backup.Debug("Shutting down global backup manager")
		GlobalBackupManager.Shutdown()
	}

	GlobalBackupManager = NewBackupManager(config)

	// Update all active HTTP handlers with the new manager
	for _, handler := range activeHTTPHandlers {
		handler.manager = GlobalBackupManager
	}

	// Start the backup manager in a goroutine to avoid blocking
	go func() {
		if err := GlobalBackupManager.Start(); err != nil {
			logger.Backup.Error("Failed to start global backup manager: " + err.Error())
		}
	}()

	// Return immediately, initialization will complete in the background
	return nil
}

// RegisterHTTPHandler registers an HTTP handler to be updated when the manager changes
func RegisterHTTPHandler(handler *HTTPHandler) {
	activeHTTPHandlers = append(activeHTTPHandlers, handler)
}

// GetBackupConfig returns a properly configured BackupConfig
func GetBackupConfig() BackupConfig {

	return BackupConfig{
		WorldName:     config.GetSaveName(),
		BackupDir:     config.GetConfiguredBackupDir(),
		SafeBackupDir: config.GetConfiguredSafeBackupDir(),
		WaitTime:      30 * time.Second, // not sure why we are not using config.BackupWaitTime here, but ill not touch it in this commit (config rework)
		RetentionPolicy: RetentionPolicy{
			KeepLastN:       config.GetBackupKeepLastN(),
			KeepDailyFor:    config.GetBackupKeepDailyFor(),
			KeepWeeklyFor:   config.GetBackupKeepWeeklyFor(),
			KeepMonthlyFor:  config.GetBackupKeepMonthlyFor(),
			CleanupInterval: config.GetBackupCleanupInterval(),
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
