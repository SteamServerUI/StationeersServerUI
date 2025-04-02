package backupsv2

import (
	"fmt"
	"path/filepath"
	"time"

	"StationeersServerUI/src/discord"

	"github.com/fsnotify/fsnotify"
)

// Start begins the backup monitoring and cleanup routines
func (m *BackupManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize backup directories: %w", err)
	}

	// Start file watcher
	watcher, err := newFsWatcher(m.config.BackupDir)
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	m.watcher = watcher

	go m.watchBackups()
	go m.startCleanupRoutine()

	return nil
}

// watchBackups monitors the backup directory for new files
func (m *BackupManager) watchBackups() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case event, ok := <-m.watcher.events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				m.handleNewBackup(event.Name)
			}
		case err, ok := <-m.watcher.errors:
			if !ok {
				return
			}
			fmt.Printf("Backup watcher error: %v\n", err)
		}
	}
}

// handleNewBackup processes a newly created backup file
func (m *BackupManager) handleNewBackup(filePath string) {
	go func() {
		time.Sleep(m.config.WaitTime)

		fileName := filepath.Base(filePath)
		dstPath := filepath.Join(m.config.SafeBackupDir, fileName)

		if err := copyFile(filePath, dstPath); err != nil {
			fmt.Printf("Error copying backup %s: %v\n", fileName, err)
			return
		}

		fmt.Printf("Backup successfully copied to safe location: %s\n", dstPath)
		discord.SendMessageToSavesChannel(fmt.Sprintf("Backup file %s copied to safe location.", dstPath))
	}()
}

// startCleanupRoutine runs periodic backup cleanup
func (m *BackupManager) startCleanupRoutine() {
	ticker := time.NewTicker(m.config.RetentionPolicy.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.Cleanup(); err != nil {
				fmt.Printf("Backup cleanup error: %v\n", err)
			}
		}
	}
}
