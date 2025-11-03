package backupmgr

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

const filetimeEpochOffset = 116444736000000000 // difference between 1601 and 1970 in 100-ns units

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

// getBackupSaveFiles retrieves all backup save files from the safe backup directory
func (m *BackupManager) getBackupSaveFiles() ([]BackupSaveFile, error) {
	var saves []BackupSaveFile

	err := filepath.WalkDir(m.config.SafeBackupDir, func(path string, de os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !de.IsDir() {
			// Process the save file
			filename := de.Name()

			// Skip invalid backup files
			if !isValidBackupFile(filename) {
				return nil
			}

			// Get the full path
			fullPath := filepath.Join(m.config.SafeBackupDir, filename)

			// Get the save time from the file
			// Unzip the save file and open the world_meta.xml file inside
			r, err := zip.OpenReader(fullPath)
			if err != nil {
				return err
			}
			defer r.Close()
			worldMetadata, err := r.Open("world_meta.xml")
			if err != nil {
				return err
			}
			defer worldMetadata.Close()
			// Read the world_meta.xml file content using the XML library
			type WorldMeta struct {
				SaveTime int64 `xml:"DateTime"`
			}
			var meta WorldMeta
			decoder := xml.NewDecoder(worldMetadata)
			if err := decoder.Decode(&meta); err != nil {
				return err
			}

			// Convert FILETIME (100-ns intervals) → Unix time (seconds + nanoseconds)
			ns := (meta.SaveTime - filetimeEpochOffset) * 100
			saveTime := time.Unix(0, ns)

			// Add the backup save file info to the list
			saves = append(saves, BackupSaveFile{
				SaveFile: fullPath,
				SaveTime: saveTime,
			})
		}
		return nil
	})

	// Handle errors
	if err != nil {
		// if the error contains no such file or directory, return nil but return a custom string intsted 	of the error
		if strings.Contains(err.Error(), "no such file or directory") || strings.Contains(err.Error(), "The system cannot find the file specified") {
			return nil, fmt.Errorf("save dir doesn't seem to exist (yet). Try starting the gameserver and click ↻ once it's up. If the Save folder exists and you still get this error, verify the 'Use New Terrain and Save System' setting. Detailed Error: %w", err)
		}
		return nil, fmt.Errorf("failed to handle safe backup dir: %w", err)
	}

	// Sort saves by save time ascending
	sort.Slice(saves, func(i, j int) bool {
		return saves[i].SaveTime.Before(saves[j].SaveTime)
	})
	// Add the index to each save
	for i := range saves {
		saves[i].Index = i
	}

	return saves, nil
}

// deleteBackupGroup removes all files in a backup group
func (m *BackupManager) deleteBackupGroup(saveFile BackupSaveFile) {
	if err := os.Remove(saveFile.SaveFile); err != nil {
		logger.Backup.Error("Failed to delete backup file " + saveFile.SaveFile + ": " + err.Error())
	}
}
