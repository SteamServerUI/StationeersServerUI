package backupmgr

import (
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time" // used by buildSaveFileIndexMap
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

// buildSaveFileIndexMap pre-computes indexes for all .save files in O(n) time
// Returns a map of filename -> index (newest file gets highest index)
func buildSaveFileIndexMap(files []os.DirEntry) map[string]int {
	// Collect .save files with their mod times
	var saveFiles []struct {
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
		saveFiles = append(saveFiles, struct {
			name    string
			modTime time.Time
		}{file.Name(), info.ModTime()})
	}

	// Sort newest first (single sort, not per-file)
	sort.Slice(saveFiles, func(i, j int) bool {
		return saveFiles[i].modTime.After(saveFiles[j].modTime)
	})

	// Build index map: newest gets highest index
	indexMap := make(map[string]int, len(saveFiles))
	for i, f := range saveFiles {
		indexMap[f.name] = len(saveFiles) - i
	}

	return indexMap
}

// parseBackupIndex extracts the backup index from a filename or uses pre-computed map for .save files
func parseBackupIndex(filename string, saveIndexMap map[string]int) int {
	// Try to extract index from old format (e.g., world(1).xml)
	re := regexp.MustCompile(`\((\d+)\)`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) >= 2 {
		index, err := strconv.Atoi(matches[1])
		if err == nil {
			return index
		}
	}

	// For .save files, use pre-computed index map
	if strings.HasSuffix(filename, ".save") {
		if index, ok := saveIndexMap[filename]; ok {
			return index
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
