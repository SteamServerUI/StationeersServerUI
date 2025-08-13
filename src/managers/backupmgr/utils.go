package backupmgr

import (
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
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

// parseBackupIndex extracts the backup index from a filename or assigns a synthetic index
func parseBackupIndex(filename string, modTime time.Time, files []os.DirEntry) int {
	// Try to extract index from old format (e.g., world(1).xml)
	re := regexp.MustCompile(`\((\d+)\)`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) >= 2 {
		index, err := strconv.Atoi(matches[1])
		if err == nil {
			return index
		}
	}

	// For .save files, assign synthetic index based on mod time (newest eq highest)
	if strings.HasSuffix(filename, ".save") {
		// Sort files by mod time to assign indexes
		var sortedFiles []struct {
			name    string
			modTime time.Time
		}
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".save") {
				continue
			}
			info, err := file.Info()
			if err != nil {
				continue
			}
			sortedFiles = append(sortedFiles, struct {
				name    string
				modTime time.Time
			}{file.Name(), info.ModTime()})
		}

		// Sort newest first
		sort.Slice(sortedFiles, func(i, j int) bool {
			return sortedFiles[i].modTime.After(sortedFiles[j].modTime)
		})

		// Find the position of the current file
		for i, f := range sortedFiles {
			if f.name == filename {
				// Assign index starting from max possible index downwards
				return len(sortedFiles) - i
			}
		}
	}

	return -1
}

// isValidBackupFile checks if a filename is a valid backup file
func isValidBackupFile(filename string) bool {
	return (strings.Contains(filename, "world") &&
		(strings.HasSuffix(filename, ".bin") ||
			strings.HasSuffix(filename, ".xml"))) ||
		strings.HasSuffix(filename, ".save")
}
