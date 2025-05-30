package backupmgr

import (
	"fmt"
	"os"
	"path/filepath"
)

// RestoreBackup restores a backup with the given index
func (m *BackupManager) RestoreBackup(index int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	files := []struct {
		backupName    string
		backupNameAlt string
		destName      string
	}{
		{fmt.Sprintf("world_meta(%d).xml", index), fmt.Sprintf("world_meta(%d)_AutoSave.xml", index), "world_meta.xml"},
		{fmt.Sprintf("world(%d).xml", index), fmt.Sprintf("world(%d)_AutoSave.xml", index), "world.xml"},
		{fmt.Sprintf("world(%d).bin", index), fmt.Sprintf("world(%d)_AutoSave.bin", index), "world.bin"},
	}

	restoredFiles := make(map[string]string)

	for _, file := range files {
		backupFile := filepath.Join(m.config.SafeBackupDir, file.backupName)
		destFile := filepath.Join(m.config.BackupDir, file.destName)

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
