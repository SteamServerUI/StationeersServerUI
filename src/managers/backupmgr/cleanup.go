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

// sameCalendarDay returns true if two times fall on the same calendar day (year + day-of-year).
func sameCalendarDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

// updateRetentionTrackers updates the daily/weekly/monthly tracker timestamps for a kept backup.
func updateRetentionTrackers(saveTime time.Time, lastKeptDaily, lastKeptWeekly, lastKeptMonthly *time.Time) {
	// Track daily
	if lastKeptDaily.IsZero() || !sameCalendarDay(saveTime, *lastKeptDaily) {
		*lastKeptDaily = saveTime
	}
	// Track weekly
	y1, w1 := saveTime.ISOWeek()
	y2, w2 := lastKeptWeekly.ISOWeek()
	if lastKeptWeekly.IsZero() || y1 != y2 || w1 != w2 {
		*lastKeptWeekly = saveTime
	}
	// Track monthly
	if lastKeptMonthly.IsZero() || saveTime.Month() != lastKeptMonthly.Month() || saveTime.Year() != lastKeptMonthly.Year() {
		*lastKeptMonthly = saveTime
	}
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

	for i, backup := range saves {
		age := now.Sub(backup.SaveTime)

		// Always keep the most recent N backups, but also update the retention
		// trackers so the daily/weekly/monthly logic doesn't redundantly keep
		// backups for days already covered by KeepLastN.
		if i < m.config.RetentionPolicy.KeepLastN {
			updateRetentionTrackers(backup.SaveTime, &lastKeptDaily, &lastKeptWeekly, &lastKeptMonthly)
			continue
		}

		// Keep daily backups for specified duration
		// Compare full calendar day (year + day-of-year) instead of just day-of-month
		// to avoid incorrectly treating e.g. Jan 15 and Feb 15 as the "same day".
		if age < m.config.RetentionPolicy.KeepDailyFor {
			if lastKeptDaily.IsZero() || !sameCalendarDay(backup.SaveTime, lastKeptDaily) {
				lastKeptDaily = backup.SaveTime
				continue
			}
		}

		// Keep weekly backups for specified duration
		if age < m.config.RetentionPolicy.KeepWeeklyFor {
			year1, week1 := backup.SaveTime.ISOWeek()
			year2, week2 := lastKeptWeekly.ISOWeek()
			if lastKeptWeekly.IsZero() || year1 != year2 || week1 != week2 {
				lastKeptWeekly = backup.SaveTime
				continue
			}
		}

		// Keep monthly backups for specified duration
		if age < m.config.RetentionPolicy.KeepMonthlyFor {
			if lastKeptMonthly.IsZero() ||
				backup.SaveTime.Month() != lastKeptMonthly.Month() ||
				backup.SaveTime.Year() != lastKeptMonthly.Year() {
				lastKeptMonthly = backup.SaveTime
				continue
			}
		}

		// If we get here, the backup should be deleted
		m.deleteBackupGroup(backup)
	}

	return nil
}

// deleteBackupGroup removes all files in a backup group
func (m *BackupManager) deleteBackupGroup(saveFile BackupSaveFile) {
	if err := os.Remove(saveFile.SaveFile); err != nil {
		logger.Backup.Error("Failed to delete backup file " + saveFile.SaveFile + ": " + err.Error())
	}
}
