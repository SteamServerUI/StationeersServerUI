package backupmgr

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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
		destFile := filepath.Join("./saves/"+m.config.WorldName, m.config.WorldName+".save")

		// Before restore, check if we have existing .save files in the root saves/WorldName dir
		saveDir := filepath.Join("./saves/", m.config.WorldName)
		files, err := os.ReadDir(saveDir)
		if err != nil {
			return fmt.Errorf("failed to read save directory %s: %w", saveDir, err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasSuffix(file.Name(), ".save") {
				existingFile := filepath.Join(saveDir, file.Name())
				// Move existing .save file to SafeBackupDir with timestamp to avoid overwrites
				timestamp := time.Now().Format("2006-01-02_15-04-05")
				savedPreviousHeadSaveFilePath := filepath.Join(m.config.SafeBackupDir, fmt.Sprintf("%s_%s_%s", "oldHeadSaveBackup", timestamp, file.Name()))
				if err := os.Rename(existingFile, savedPreviousHeadSaveFilePath); err != nil {
					return fmt.Errorf("failed to move existing HEAD .save file %s to %s: %w", existingFile, savedPreviousHeadSaveFilePath, err)
				}
				logger.Backup.Info("Moved previous HEAD .save file to: " + savedPreviousHeadSaveFilePath)
			}
		}

		// Create temp directory for mod time shenenigans (https://discordapp.com/channels/276525882049429515/392080751648178188/1407157281606336602)
		tempDir := filepath.Join("./saves", m.config.WorldName, "tmp")
		if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create temp directory %s: %w", tempDir, err)
		}
		defer os.RemoveAll(tempDir)

		// Extract .save (zip) file to tempDir
		r, err := zip.OpenReader(backupFile)
		if err != nil {
			return fmt.Errorf("failed to open zip reader for %s: %w", backupFile, err)
		}
		defer r.Close()

		for _, f := range r.File {
			path := filepath.Join(tempDir, f.Name)
			if f.FileInfo().IsDir() {
				if err := os.MkdirAll(path, f.Mode()); err != nil {
					m.revertRestore(restoredFiles)
					return fmt.Errorf("failed to create directory %s: %w", path, err)
				}
				continue
			}

			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				m.revertRestore(restoredFiles)
				return fmt.Errorf("failed to create parent directory for %s: %w", path, err)
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				m.revertRestore(restoredFiles)
				return fmt.Errorf("failed to create file %s: %w", path, err)
			}

			rc, err := f.Open()
			if err != nil {
				outFile.Close()
				m.revertRestore(restoredFiles)
				return fmt.Errorf("failed to open file in zip %s: %w", f.Name, err)
			}

			if _, err := io.Copy(outFile, rc); err != nil {
				rc.Close()
				outFile.Close()
				m.revertRestore(restoredFiles)
				return fmt.Errorf("failed to extract file %s: %w", path, err)
			}
			rc.Close()
			outFile.Close()
		}

		// Modify timestamps of extracted files to current system time
		now := time.Now()
		if err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			return os.Chtimes(path, now, now)
		}); err != nil {
			m.revertRestore(restoredFiles)
			return fmt.Errorf("failed to modify timestamps in %s: %w", tempDir, err)
		}

		// Create new .save (zip) file at destFile with updated timestamps
		dest, err := os.Create(destFile)
		if err != nil {
			m.revertRestore(restoredFiles)
			return fmt.Errorf("failed to create destination .save file %s: %w", destFile, err)
		}
		defer dest.Close()

		w := zip.NewWriter(dest)
		defer w.Close()

		if err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			relPath, err := filepath.Rel(tempDir, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}
			relPath = filepath.ToSlash(relPath)

			// Create zip entry with current system timestamp
			fw, err := w.CreateHeader(&zip.FileHeader{
				Name:     relPath,
				Method:   zip.Deflate,
				Modified: now,
			})
			if err != nil {
				return fmt.Errorf("failed to create zip entry %s: %w", relPath, err)
			}

			srcFile, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer srcFile.Close()

			if _, err := io.Copy(fw, srcFile); err != nil {
				return fmt.Errorf("failed to write file %s to zip: %w", relPath, err)
			}
			return nil
		}); err != nil {
			m.revertRestore(restoredFiles)
			return fmt.Errorf("failed to restore .save file %s: %w", backupFile, err)
		}
		restoredFiles[destFile] = backupFile
		return nil // restore and mod time shenanigans successful, no need to return an error
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
			destFile := filepath.Join("./saves/"+m.config.WorldName, file.destName)

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
