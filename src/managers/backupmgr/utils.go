package backupmgr

import (
	"io"
	"os"
	"strings"
)

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return destination.Sync()
}

// isValidBackupFile checks if a filename is a valid backup file
func isValidBackupFile(filename string) bool {
	return strings.HasSuffix(filename, ".save")
}
