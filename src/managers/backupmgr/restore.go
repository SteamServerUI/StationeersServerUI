package backupmgr

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// RestoreBackup restores a backup with the given index
func (m *BackupManager) RestoreBackup(index int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	logger.Backup.Infof("Restoring backup with index %d", index)

	saves, err := m.getBackupSaveFiles()
	if err != nil {
		return fmt.Errorf("failed to get backup groups: %w", err)
	}

	if index < 0 || index >= len(saves) {
		return fmt.Errorf("invalid backup index %d", index)
	}
	targetSave := saves[index]

	backupFile := targetSave.SaveFile
	destFile := filepath.Join("./saves/"+m.config.WorldName, m.config.WorldName+".save")

	// Backup current save before restoring
	tmpfile, err := os.CreateTemp("", "ssui_restore_bak")
	if err != nil {
		return fmt.Errorf("failed to create temp backup file for current save: %w", err)
	}
	defer tmpfile.Close()
	// Note we do not defer removal of tmpfile here, as we may need it on multiple failures

	// Copy current save to temp backup file
	if err := copyFile(destFile, tmpfile.Name()); err != nil {
		os.Remove(tmpfile.Name())
		return fmt.Errorf("failed to backup current save before restore: %w", err)
	}

	// Copy the backup file to the destination
	if err := copyFile(backupFile, destFile); err != nil {
		// Restore the original save from temp backup on failure
		err2 := copyFile(tmpfile.Name(), destFile)
		if err2 != nil {
			return fmt.Errorf("failed to restore backup file: %w; additionally failed to restore original save: %w. Main save backup located at: %s", err, err2, tmpfile.Name())
		}
		os.Remove(tmpfile.Name())
		return fmt.Errorf("failed to restore backup file: %w", err)
	}
	os.Remove(tmpfile.Name())

	logger.Backup.Infof("Backup with index %d restored successfully", index)

	return nil
}
