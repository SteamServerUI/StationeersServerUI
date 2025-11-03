package backupmgr

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// Cleanup performs backup cleanup according to retention policy
func (m *BackupManager) Cleanup() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clean regular backup dir (keep only recent)
	if err := m.cleanBackupDir(); err != nil {
		return fmt.Errorf("backup dir cleanup failed: %w", err)
	}

	// Clean safe backup dir with retention policy
	if err := m.cleanSafeBackupDir(); err != nil {
		return fmt.Errorf("safe backup dir cleanup failed: %w", err)
	}

	return nil
}

// cleanBackupDir cleans the regular backup directory
func (m *BackupManager) cleanBackupDir() error {
	files, err := os.ReadDir(m.config.BackupDir)
	if err != nil {
		return err
	}

	now := time.Now()
	cutoff := now.Add(-24 * time.Hour) // Keep only files from last 24 hours

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullPath := filepath.Join(m.config.BackupDir, file.Name())
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(fullPath); err != nil {
				logger.Backup.Error("Failed to remove old backup " + fullPath + ": " + err.Error())
			}
		}
	}

	return nil
}

// cleanSafeBackupDir cleans the safe backup directory with retention policy
func (m *BackupManager) cleanSafeBackupDir() error {
	saves, err := m.getBackupSaveFiles()
	if err != nil {
		return err
	}

	// Sort newest first
	sort.Slice(saves, func(i, j int) bool {
		return saves[i].SaveTime.After(saves[j].SaveTime)
	})

	now := time.Now()
	var (
		lastKeptDaily   time.Time
		lastKeptWeekly  time.Time
		lastKeptMonthly time.Time
	)

	for i, group := range saves {
		age := now.Sub(group.SaveTime)

		// Always keep the most recent N backups
		if i < m.config.RetentionPolicy.KeepLastN {
			continue
		}

		// Keep daily backups for specified duration
		if age < m.config.RetentionPolicy.KeepDailyFor {
			if lastKeptDaily.IsZero() || group.SaveTime.Day() != lastKeptDaily.Day() {
				lastKeptDaily = group.SaveTime
				continue
			}
		}

		// Keep weekly backups for specified duration
		if age < m.config.RetentionPolicy.KeepWeeklyFor {
			year1, week1 := group.SaveTime.ISOWeek()
			year2, week2 := lastKeptWeekly.ISOWeek()
			if lastKeptWeekly.IsZero() || year1 != year2 || week1 != week2 {
				lastKeptWeekly = group.SaveTime
				continue
			}
		}

		// Keep monthly backups for specified duration
		if age < m.config.RetentionPolicy.KeepMonthlyFor {
			if lastKeptMonthly.IsZero() ||
				group.SaveTime.Month() != lastKeptMonthly.Month() ||
				group.SaveTime.Year() != lastKeptMonthly.Year() {
				lastKeptMonthly = group.SaveTime
				continue
			}
		}

		// If we get here, the backup should be deleted
		m.deleteBackupGroup(group)
	}

	return nil
}

// deleteBackupGroup removes all files in a backup group
func (m *BackupManager) deleteBackupGroup(saveFile BackupSaveFile) {
	if err := os.Remove(saveFile.SaveFile); err != nil {
		logger.Backup.Error("Failed to delete backup file " + saveFile.SaveFile + ": " + err.Error())
	}
}
