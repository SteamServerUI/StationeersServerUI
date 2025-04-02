package backupsv2

import (
	"io"
	"os"
	"regexp"
	"strconv"
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

// parseBackupIndex extracts the backup index from a filename
func parseBackupIndex(filename string) int {
	re := regexp.MustCompile(`\((\d+)\)`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return -1
	}

	index, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1
	}

	return index
}

func isValidBackupFile(filename string) bool {
	return strings.Contains(filename, "world") &&
		(strings.HasSuffix(filename, ".bin") ||
			strings.HasSuffix(filename, ".xml"))
}
