package backupmgr

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	groups, err := m.getBackupGroups()
	if err != nil {
		return err
	}

	// Sort newest first
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].ModTime.After(groups[j].ModTime)
	})

	now := time.Now()
	var (
		lastKeptDaily   time.Time
		lastKeptWeekly  time.Time
		lastKeptMonthly time.Time
	)

	for i, group := range groups {
		age := now.Sub(group.ModTime)

		// Always keep the most recent N backups
		if i < m.config.RetentionPolicy.KeepLastN {
			continue
		}

		// Keep daily backups for specified duration
		if age < m.config.RetentionPolicy.KeepDailyFor {
			if lastKeptDaily.IsZero() || group.ModTime.Day() != lastKeptDaily.Day() {
				lastKeptDaily = group.ModTime
				continue
			}
		}

		// Keep weekly backups for specified duration
		if age < m.config.RetentionPolicy.KeepWeeklyFor {
			year1, week1 := group.ModTime.ISOWeek()
			year2, week2 := lastKeptWeekly.ISOWeek()
			if lastKeptWeekly.IsZero() || year1 != year2 || week1 != week2 {
				lastKeptWeekly = group.ModTime
				continue
			}
		}

		// Keep monthly backups for specified duration
		if age < m.config.RetentionPolicy.KeepMonthlyFor {
			if lastKeptMonthly.IsZero() ||
				group.ModTime.Month() != lastKeptMonthly.Month() ||
				group.ModTime.Year() != lastKeptMonthly.Year() {
				lastKeptMonthly = group.ModTime
				continue
			}
		}

		// If we get here, the backup should be deleted
		m.deleteBackupGroup(group)
	}

	return nil
}

// getBackupGroups collects and groups backup files
func (m *BackupManager) getBackupGroups() ([]BackupGroup, error) {
	files, err := os.ReadDir(m.config.SafeBackupDir)
	if err != nil {
		return nil, err
	}

	groups := make(map[int]BackupGroup)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		index := parseBackupIndex(file.Name())
		if index == -1 {
			continue
		}

		fullPath := filepath.Join(m.config.SafeBackupDir, file.Name())
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		group := groups[index]
		group.Index = index
		group.ModTime = info.ModTime()

		switch {
		case strings.HasSuffix(file.Name(), ".bin"):
			group.BinFile = fullPath
		case strings.Contains(file.Name(), "world(") && strings.HasSuffix(file.Name(), ".xml"):
			group.XMLFile = fullPath
		case strings.Contains(file.Name(), "world_meta(") && strings.HasSuffix(file.Name(), ".xml"):
			group.MetaFile = fullPath
		}

		groups[index] = group
	}

	var result []BackupGroup
	for _, group := range groups {
		if group.BinFile != "" && group.XMLFile != "" && group.MetaFile != "" {
			result = append(result, group)
		}
	}

	return result, nil
}

// deleteBackupGroup removes all files in a backup group
func (m *BackupManager) deleteBackupGroup(group BackupGroup) {
	for _, file := range []string{group.BinFile, group.XMLFile, group.MetaFile} {
		if file != "" {
			if err := os.Remove(file); err != nil {
				logger.Backup.Error("Failed to delete backup file " + file + ": " + err.Error())
			}
		}
	}
}
