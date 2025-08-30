package backupmgr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/commandmgr"

	"github.com/fsnotify/fsnotify"
)

/*
The BackupManager manages backup operations. Each instance is independent with its own config and context.
Background routines (file watching and cleanup) only start when Start() is called. Multiple instances
can coexist but may conflict if configured with overlapping directories.
*/

// Initialize checks for BackupDir and waits until it exists, then ensures SafeBackupDir exists.
// It returns a channel that signals when initialization is complete or an error occurs.
func (m *BackupManager) Initialize() <-chan error {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make(chan error, 1)

	go func() {
		const timeout = 90 * time.Minute
		const pollInterval = 2500 * time.Millisecond
		deadline := time.Now().Add(timeout)

		// Wait for BackupDir to exist
		for {
			if _, err := os.Stat(m.config.BackupDir); err == nil {
				// Directory exists, proceed
				break
			} else if !os.IsNotExist(err) {
				// An error other than "not exists" occurred
				result <- fmt.Errorf("error checking backup directory %s: %v", m.config.BackupDir, err)
				return
			}

			if time.Now().After(deadline) {
				result <- fmt.Errorf("timeout waiting for backup directory %s to be created", m.config.BackupDir)
				return
			}
			logger.Backup.Debug("Backup manager waiting for save folder " + m.config.BackupDir + " to be created by Stationeers...")

			// Wait before checking again
			time.Sleep(pollInterval)
		}

		// Ensure SafeBackupDir exists, create it if it doesn't
		if err := os.MkdirAll(m.config.SafeBackupDir, os.ModePerm); err != nil {
			result <- fmt.Errorf("error creating safe backup directory %s: %v", m.config.SafeBackupDir, err)
			return
		}
		logger.Backup.Debug("Backup manager created safebackups dir successfully")

		result <- nil
	}()

	return result
}

// Start begins the backup monitoring and cleanup routines
func (m *BackupManager) Start() error {
	// Wait for initialization to complete
	logger.Backup.Debug("Backup manager is waiting for save folder initialization...")
	initResult := <-m.Initialize()
	if initResult != nil {
		return fmt.Errorf("failed to initialize backup manager: %w", initResult)
	}
	logger.Backup.Info("Backup manager started")

	// Start file watcher
	watcher, err := newFsWatcher(m.config.BackupDir)
	if err != nil {
		return fmt.Errorf("failed to create autosave watcher: %w", err)
	}
	m.watcher = watcher
	go m.watchBackups()

	if config.IsCleanupEnabled {
		go m.startCleanupRoutine()
	}

	return nil
}

// watchBackups monitors the backup directory for new files
func (m *BackupManager) watchBackups() {
	m.wg.Add(1)
	defer m.wg.Done()

	logger.Backup.Debug("Starting backup file watcher...")
	defer logger.Backup.Debug("Backup file watcher stopped")

	for {
		select {
		case <-m.ctx.Done():
			return
		case event, ok := <-m.watcher.events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				logger.Backup.Info("New backup file detected: " + event.Name)
				m.handleNewBackup(event.Name)
			}
		case err, ok := <-m.watcher.errors:
			if !ok {
				return
			}
			logger.Backup.Error("Backup watcher error: " + err.Error())
		}
	}
}

// handleNewBackup processes a newly created backup file
func (m *BackupManager) handleNewBackup(filePath string) {
	if !isValidBackupFile(filepath.Base(filePath)) {
		return
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		time.Sleep(m.config.WaitTime)

		// save the world into Head save too if SSCM is enabled
		if config.IsSSCMEnabled && config.IsNewTerrainAndSaveSystem {
			commandmgr.WriteCommand("SAVE")
			logger.Backup.Debug("HEAD Save triggered via SSCM")
		} else {
			logger.Backup.Debug("HEAD Save NOT refreshed via SSCM")
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		fileName := filepath.Base(filePath)
		relativePath, err := filepath.Rel(m.config.BackupDir, filePath)
		if err != nil {
			logger.Backup.Error("Error getting relative path for " + filePath + ": " + err.Error())
			return
		}
		dstPath := filepath.Join(m.config.SafeBackupDir, relativePath)

		if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
			logger.Backup.Error("Error creating destination dir for " + dstPath + ": " + err.Error())
			return
		}

		if err := copyFile(filePath, dstPath); err != nil {
			logger.Backup.Error("Error copying backup " + fileName + ": " + err.Error())
			return
		}

		logger.Backup.Debug("Backup successfully copied to safe location: " + dstPath)
	}()
}

// startCleanupRoutine runs periodic backup cleanup
func (m *BackupManager) startCleanupRoutine() {
	m.wg.Add(1)
	defer m.wg.Done()

	ticker := time.NewTicker(m.config.RetentionPolicy.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.Cleanup(); err != nil {
				logger.Backup.Error("Backup cleanup error: " + err.Error())
			}
		}
	}
}

// ListBackups returns information about available backups
// limit: number of recent backups to return (0 for all)
func (m *BackupManager) ListBackups(limit int) ([]BackupGroup, error) {
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

	if limit > 0 && limit < len(groups) {
		groups = groups[:limit]
	}

	return groups, nil
}

// Shutdown stops all backup operations
func (m *BackupManager) Shutdown() {
	logger.Backup.Debug("Shutting down backup manager...")

	m.mu.Lock()
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}

	if m.watcher != nil {
		m.watcher.close()
		m.watcher = nil
	}
	m.mu.Unlock()

	// Wait for all goroutines to finish
	logger.Backup.Debug("Waiting for background tasks to complete...")
	m.wg.Wait()

	logger.Backup.Debug("Backup manager shut down completely")
}

// NewBackupManager creates a new BackupManager instance
func NewBackupManager(cfg BackupConfig) *BackupManager {
	ctx, cancel := context.WithCancel(context.Background())

	if cfg.WaitTime == 0 {
		cfg.WaitTime = defaultWaitTime
	}

	return &BackupManager{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}
