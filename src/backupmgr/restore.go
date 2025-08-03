package backupmgr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// RestoreBackup restores a backup with the given index
func (m *BackupManager) RestoreBackup(index int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	logger.Backup.Info("Restoring backup with index " + fmt.Sprintf("%d", index))

	groups, err := m.getBackupGroups()
	if err != nil {
		return fmt.Errorf("failed to get backup groups: %w", err)
	}

	var targetGroup BackupGroup
	for _, group := range groups {
		if group.Index == index {
			targetGroup = group
			break
		}
	}

	if targetGroup.Index == 0 {
		return fmt.Errorf("no backup found with index %d", index)
	}

	restoredFiles := make(map[string]string)

	// Handle .save file or old-style trio
	if targetGroup.BinFile != "" && strings.HasSuffix(targetGroup.BinFile, ".save") {
		// .save file case
		backupFile := targetGroup.BinFile
		destFile := filepath.Join("./saves/"+config.WorldName, filepath.Base(backupFile))

		if err := copyFile(backupFile, destFile); err != nil {
			m.revertRestore(restoredFiles)
			return fmt.Errorf("failed to restore .save file %s: %w", backupFile, err)
		}
		restoredFiles[destFile] = backupFile
	} else {
		// Old-style trio (world_meta.xml, world.xml, world.bin)
		files := []struct {
			backupName    string
			backupNameAlt string
			destName      string
		}{
			{fmt.Sprintf("world_meta(%d).xml", index), fmt.Sprintf("world_meta(%d)_AutoSave.xml", index), "world_meta.xml"},
			{fmt.Sprintf("world(%d).xml", index), fmt.Sprintf("world(%d)_AutoSave.xml", index), "world.xml"},
			{fmt.Sprintf("world(%d).bin", index), fmt.Sprintf("world(%d)_AutoSave.bin", index), "world.bin"},
		}

		for _, file := range files {
			backupFile := filepath.Join(m.config.SafeBackupDir, file.backupName)
			destFile := filepath.Join("./saves/"+config.WorldName, file.destName)

			if err := copyFile(backupFile, destFile); err != nil {
				// Try alternative name
				backupFileAlt := filepath.Join(m.config.SafeBackupDir, file.backupNameAlt)
				if err := copyFile(backupFileAlt, destFile); err != nil {
					m.revertRestore(restoredFiles)
					return fmt.Errorf("failed to restore %s: %w", file.backupName, err)
				}
				backupFile = backupFileAlt
			}
			restoredFiles[destFile] = backupFile
		}
	}
	logger.Backup.Debug(fmt.Sprintf("%v", restoredFiles))

	return nil
}

// revertRestore undoes a failed restore operation
func (m *BackupManager) revertRestore(restoredFiles map[string]string) {
	for destFile, backupFile := range restoredFiles {
		if err := os.Remove(destFile); err == nil {
			_ = copyFile(backupFile, destFile)
		}
	}
}
