package setup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// githubRelease represents the structure of a GitHub release response
type githubRelease struct {
	TagName    string `json:"tag_name"`
	Prerelease bool   `json:"prerelease"`
	Assets     []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

// Version holds semantic version components
type Version struct {
	Major int
	Minor int
	Patch int
}

// UpdateExecutable checks for and applies the latest release from GitHub
func UpdateExecutable() error {
	if !config.IsUpdateEnabled {
		logger.Install.Warn("⚠️ Update check is disabled. Skipping update check.")
		return nil
	}

	if config.Branch != "release" {
		logger.Install.Warn("⚠️ You are running a development build. Skipping update check.")
		return nil
	}

	logger.Install.Info("🕵️ Querying GitHub API for the latest release...")
	latestRelease, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("❌ Failed to fetch latest release: %v", err)
	}

	// Parse current and latest versions
	currentVer, err := parseVersion(config.Version)
	if err != nil {
		return fmt.Errorf("❌ Failed to parse current version %s: %v", config.Version, err)
	}
	latestVer, err := parseVersion(latestRelease.TagName)
	if err != nil {
		return fmt.Errorf("❌ Failed to parse latest version %s: %v", latestRelease.TagName, err)
	}

	logger.Install.Info(fmt.Sprintf("Current version: %s, Latest version: %s", config.Version, latestRelease.TagName))

	// Check pre-release status
	if latestRelease.Prerelease && !config.AllowPrereleaseUpdates {
		logger.Install.Warn(fmt.Sprintf("⚠️ Latest version %s is a pre-release. Enable 'AllowPrerelease' in config to update.", latestRelease.TagName))
		return nil
	}

	// Check if we should update
	updateReason, shouldUpdate := shouldUpdate(currentVer, latestVer)
	if !shouldUpdate {
		switch updateReason {
		case "up-to-date":
			logger.Install.Info("🎉 No update needed: you’re already on the latest version.")
		case "major-update":
			logger.Install.Warn(fmt.Sprintf("⚠️ Latest version %s is a major update from %s. Major Updates include Breaking changes in this project. Read the release notes and backup your Server folder before updating. Enable 'AllowMajorUpdates' in config to proceed.", latestRelease.TagName, config.Version))
		}
		return nil
	}

	// Proceed with update
	expectedExt := ".exe"
	if runtime.GOOS != "windows" {
		expectedExt = ".x86_64"
	}
	expectedExe := fmt.Sprintf("StationeersServerControl%s%s", latestRelease.TagName, expectedExt)

	// Find the asset
	var downloadURL string
	for _, asset := range latestRelease.Assets {
		if asset.Name == expectedExe {
			downloadURL = asset.URL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("❌ No matching asset found for %s", expectedExe)
	}

	// Download and replace
	logger.Install.Info(fmt.Sprintf("📡 Updating from %s to %s...", config.Version, latestRelease.TagName))
	if err := downloadNewExecutable(expectedExe, downloadURL); err != nil {
		logger.Install.Warn(fmt.Sprintf("⚠️ Update failed: %v. Keeping version %s.", err, config.Version))
		return err
	}

	// Set executable permissions on Linux
	if runtime.GOOS != "windows" {
		if err := os.Chmod(expectedExe, 0755); err != nil {
			logger.Install.Warn(fmt.Sprintf("⚠️ Update failed: couldn’t make %s executable: %v. Keeping version %s.", expectedExe, err, config.Version))
			return err
		}
	}

	// Launch the new executable and exit
	logger.Install.Info("🚀 Launching the new version and retiring the old one...")
	if runtime.GOOS == "windows" {
		if err := runAndExit(expectedExe); err != nil {
			logger.Install.Warn(fmt.Sprintf("⚠️ Update failed: couldn’t launch %s: %v. Keeping version %s.", expectedExe, err, config.Version))
			return err
		}
	}
	if runtime.GOOS == "linux" {
		if err := runAndExitLinux(expectedExe); err != nil {
			logger.Install.Warn(fmt.Sprintf("⚠️ Update failed: couldn’t launch %s: %v. Keeping version %s.", expectedExe, err, config.Version))
			return err
		}
	}

	return nil
}

func RestartMySelf() {
	currentExe, err := os.Executable()
	if err != nil {
		logger.Install.Warn(fmt.Sprintf("⚠️ Restart failed: couldn’t get current executable path: %v. Keeping version %s.", err, config.Version))
		return
	}

	if runtime.GOOS == "windows" {
		if err := runAndExit(currentExe); err != nil {
			logger.Install.Warn(fmt.Sprintf("⚠️ Restart failed: couldn’t launch %s: %v. Keeping version %s.", currentExe, err, config.Version))
			return
		}
	}
	if runtime.GOOS == "linux" {
		if err := runAndExitLinux(currentExe); err != nil {
			logger.Install.Warn(fmt.Sprintf("⚠️ Restart failed: couldn’t launch %s: %v. Keeping version %s.", currentExe, err, config.Version))
			return
		}
	}
}

// parseVersion parses a version string (e.g., "4.6.10") into a Version struct and tries to handle a few culprits too
func parseVersion(v string) (Version, error) {
	v = strings.TrimPrefix(v, "v")
	if idx := strings.Index(v, "-"); idx != -1 {
		v = v[:idx]
	}

	var ver Version
	_, err := fmt.Sscanf(v, "%d.%d.%d", &ver.Major, &ver.Minor, &ver.Patch)
	if err != nil {
		return Version{}, fmt.Errorf("no valid X.Y.Z in tag: %s", v)
	}
	return ver, nil
}

// shouldUpdate determines if an update should proceed, returning reason if not
func shouldUpdate(current, latest Version) (string, bool) {
	// Check if already up-to-date or older
	if latest.Major < current.Major ||
		(latest.Major == current.Major && latest.Minor < current.Minor) ||
		(latest.Major == current.Major && latest.Minor == current.Minor && latest.Patch <= current.Patch) {
		return "up-to-date", false
	}

	// Check if it’s a major update and not allowed
	if current.Major != latest.Major && !config.AllowMajorUpdates {
		return "major-update", false
	}

	return "", true
}

// getLatestRelease fetches the latest release info from GitHub API
func getLatestRelease() (*githubRelease, error) {
	resp, err := http.Get("https://api.github.com/repos/JacksonTheMaster/StationeersServerUI/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response from GitHub API: %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub API response: %v", err)
	}
	return &release, nil
}

// downloadNewExecutable downloads the new executable with a progress bar
func downloadNewExecutable(filename, url string) error {
	// Use a temp file to avoid partial downloads
	tmpFile := filename + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile) // Clean up .tmp on any failure after creation

	// Download from GitHub
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return fmt.Errorf("failed to fetch %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out.Close()
		return fmt.Errorf("bad response from download: %s", resp.Status)
	}

	// Show progress
	counter := &writeCounter{Total: resp.ContentLength}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		out.Close()
		return fmt.Errorf("failed to write download to file: %v", err)
	}

	// Explicitly close the file before renaming
	if err := out.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %v", err)
	}

	// Rename temp file to final name
	if err := os.Rename(tmpFile, filename); err != nil {
		return fmt.Errorf("failed to rename temp file to %s: %v", filename, err)
	}

	logger.Install.Info("✅ Downloaded " + filename)
	return nil
}

// runAndExit launches the new executable and terminates the current process
func runAndExit(newExe string) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(newExe)
	if err != nil {
		return fmt.Errorf("❌ Couldn’t resolve path to %s: %v", newExe, err)
	}

	// Prepare the new process
	cmd := exec.Command(absPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set SysProcAttr based on OS using the OS-specific implementation
	setSysProcAttr(cmd)

	// Start the new process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("❌ Failed to start new executable: %v", err)
	}

	// Exit gracefully
	logger.Install.Warn("✨ New version’s live! Catch you on the flip side!")
	time.Sleep(500 * time.Millisecond) // Dramatic pause
	os.Exit(0)
	return nil
}

func runAndExitLinux(newExe string) error {
	absPath, err := filepath.Abs(newExe)
	if err != nil {
		return fmt.Errorf("❌ Couldn’t resolve path to %s: %v", newExe, err)
	}

	// Use syscall.Exec to replace the current process
	logger.Install.Warn("✨ New version’s live! Catch you on the flip side!")
	time.Sleep(500 * time.Millisecond)

	// Replace the current process with the new executable
	err = syscall.Exec(absPath, []string{absPath}, os.Environ())
	if err != nil {
		return fmt.Errorf("❌ Failed to exec new executable: %v", err)
	}

	// This line is never reached if Exec succeeds
	return nil
}

// writeCounter tracks download progress
type writeCounter struct {
	Total int64
	count int64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.count += int64(n)
	wc.printProgress()
	return n, nil
}

func (wc *writeCounter) printProgress() {
	// If we don't know the total size, just show downloaded bytes
	if wc.Total <= 0 {
		logger.Backup.Info(fmt.Sprintf("\r%s downloaded", bytesToHuman(wc.count)))
		return
	}

	// Calculate percentage with bounds checking
	percent := float64(wc.count) / float64(wc.Total) * 100
	if percent > 100 {
		percent = 100
	}

	// Create simple progress bar
	width := 20
	complete := int(percent / 100 * float64(width))

	progressBar := "["
	for i := 0; i < width; i++ {
		if i < complete {
			progressBar += "="
		} else if i == complete && complete < width {
			progressBar += ">"
		} else {
			progressBar += " "
		}
	}
	progressBar += "]"

	// Print progress and erase to end of line
	logger.Backup.Info(fmt.Sprintf("\r%s %.1f%% (%s/%s)",
		progressBar,
		percent,
		bytesToHuman(wc.count),
		bytesToHuman(wc.Total)))
}

// bytesToHuman converts bytes to human readable format
func bytesToHuman(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
