package runfilegallery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// GalleryRunfile represents a runfile in the gallery
type GalleryRunfile struct {
	Name          string `json:"name"`
	Filename      string `json:"filename"`
	Version       string `json:"version"`
	BackgroundURL string `json:"background_url"`
	LogoURL       string `json:"logo_url"`
	SupportedOS   string `json:"supported_os"`
	MinVersion    string `json:"min_version"`
}

// galleryCache stores the parsed and filtered runfile list
var (
	galleryCache []GalleryRunfile
	cacheMutex   sync.Mutex
)

// GetRunfileGallery fetches the list of available runfiles from GitHub Pages
func GetRunfileGallery(forceUpdate bool) ([]GalleryRunfile, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Return cached results if not forcing an update and cache is populated
	if !forceUpdate && len(galleryCache) > 0 {
		logger.Runfile.Debug("Serving runfile gallery from cache")
		return galleryCache, nil
	}

	// Fetch manifest from GitHub Pages
	const manifestURL = "https://steamserverui.github.io/runfiles/manifest.ssui"
	logger.Runfile.Info("Fetching runfile gallery from " + manifestURL)
	resp, err := http.Get(manifestURL)
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to fetch manifest: %v", err))
		return nil, fmt.Errorf("couldn't reach the gallery, network's playing hide and seek")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Runfile.Error(fmt.Sprintf("Manifest fetch failed with status: %d", resp.StatusCode))
		return nil, fmt.Errorf("gallery failed to fetch manifest from GitHub, status: %d", resp.StatusCode)
	}

	// Parse manifest
	var runfiles []GalleryRunfile
	if err := json.NewDecoder(resp.Body).Decode(&runfiles); err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to parse manifest: %v", err))
		return nil, fmt.Errorf("manifest is gibberish, can't make sense of it")
	}

	// Filter by backend version
	currentVersion := config.GetVersion()
	var filtered []GalleryRunfile
	for _, rf := range runfiles {
		if compareVersions(rf.MinVersion, currentVersion) <= 0 {
			filtered = append(filtered, rf)
		} else {
			logger.Runfile.Debug(fmt.Sprintf("Skipping runfile %s, requires version %s, current is %s", rf.Name, rf.MinVersion, currentVersion))
		}
	}

	// Update cache
	galleryCache = filtered
	logger.Runfile.Info(fmt.Sprintf("Fetched and cached %d runfiles", len(filtered)))

	if len(filtered) == 0 {
		logger.Runfile.Warn("No runfiles compatible with backend version " + currentVersion)
	}

	return filtered, nil
}

// saveRunfileToDisk downloads a runfile by identifier and saves it to RunfilesDir
func SaveRunfileToDisk(identifier string) error {
	filename := fmt.Sprintf("run%s.ssui", identifier)
	baseURL := "https://steamserverui.github.io/runfiles"
	fileURL := fmt.Sprintf("%s/%s", baseURL, filename)

	logger.Runfile.Info("Fetching runfile from " + fileURL)
	resp, err := http.Get(fileURL)
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to fetch runfile %s: %v", filename, err))
		return fmt.Errorf("couldn't grab %s, network's being a jerk", filename)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Runfile.Error(fmt.Sprintf("Runfile %s fetch failed with status: %d", filename, resp.StatusCode))
		return fmt.Errorf("%s is playing hard to get, status: %d", filename, resp.StatusCode)
	}

	saveFilePath := filepath.Join(config.GetRunfilesFolder(), filename)
	logger.Runfile.Info("Saving runfile to " + saveFilePath)

	// Create or overwrite the file
	file, err := os.Create(saveFilePath)
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to create file %s: %v", saveFilePath, err))
		return fmt.Errorf("disk's throwing a fit, can't save file")
	}
	defer file.Close()

	// Copy response body to file
	if _, err := io.Copy(file, resp.Body); err != nil {
		logger.Runfile.Error(fmt.Sprintf("Failed to save runfile %s: %v", filename, err))
		return fmt.Errorf("couldn't save %s, disk's being dramatic", filename)
	}

	logger.Runfile.Info("Successfully saved runfile " + filename)
	return nil
}

// compareVersions compares two semantic version strings (x.y.z)
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	// Ensure both versions have 3 parts
	for i := 0; i < 3; i++ {
		var n1, n2 int
		if i < len(v1Parts) {
			fmt.Sscanf(v1Parts[i], "%d", &n1)
		}
		if i < len(v2Parts) {
			fmt.Sscanf(v2Parts[i], "%d", &n2)
		}
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}
	return 0
}
