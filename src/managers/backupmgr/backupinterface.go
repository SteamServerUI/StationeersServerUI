package backupmgr

import (
	"sync"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/google/uuid"
)

// GlobalBackupManager is the singleton instance of the backup manager
var GlobalBackupManager *BackupManager

// Track all HTTP handlers that need updating when manager changes
var activeHTTPHandlers []*HTTPHandler

// initMutex ensures thread-safe initialization of the global backup manager
var initMutex sync.Mutex

// InitGlobalBackupManager initializes the global backup manager instance
func InitGlobalBackupManager(config BackupConfig) error {
	// Lock to prevent concurrent initialization
	initMutex.Lock()
	defer initMutex.Unlock()

	// Shut down existing manager if it exists
	if GlobalBackupManager != nil {
		logger.Backup.Debugf("%s Previous Backup manager found. Shutting it down.", config.Identifier)
		GlobalBackupManager.Shutdown()
		GlobalBackupManager = nil // Clear the manager to avoid stale references
	}

	logger.Backup.Debugf("%s Creating a global backup manager with ID %s", config.Identifier, config.Identifier)
	manager := NewBackupManager(config)
	GlobalBackupManager = manager

	// Update all active HTTP handlers with the new manager
	for _, handler := range activeHTTPHandlers {
		handler.manager = GlobalBackupManager
	}

	// Start the backup manager in a goroutine to avoid blocking
	go func(m *BackupManager) {
		if err := m.Start(config.Identifier); err != nil {
			logger.Backup.Warnf("%s Exited: "+err.Error(), config.Identifier)
		}
	}(manager)

	logger.Backup.Infof("%s Backup manager reloaded successfully", config.Identifier)
	return nil
}

// RegisterHTTPHandler registers an HTTP handler to be updated when the manager changes
func RegisterHTTPHandler(handler *HTTPHandler) {
	activeHTTPHandlers = append(activeHTTPHandlers, handler)
}

// GetBackupConfig returns a properly configured BackupConfig
func GetBackupConfig() BackupConfig {

	uuid := uuid.New()
	bmIdentifier := "[BM" + uuid.String()[:6] + "]:"
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
		Identifier: bmIdentifier,
	}
}

// ReloadBackupManagerFromConfig reloads the global backup manager with the current config. This should be called whenever the config is changed.
func ReloadBackupManagerFromConfig() error {
	// Create a new backupManager config from the global config
	backupConfig := GetBackupConfig()

	// Reinitialize the global backup manager with the new config
	return InitGlobalBackupManager(backupConfig)
}
