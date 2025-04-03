package install

import (
	"StationeersServerUI/src/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// githubRelease represents the structure of a GitHub release response
type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

// UpdateExecutable checks for and applies the latest release from GitHub
func UpdateExecutable() error {
	// Get current executable name
	currentExe := filepath.Base(os.Args[0])

	if !config.IsUpdateEnabled {
		fmt.Println(string(colorYellow), "⚠️ Update check is disabled. Skipping update check.", string(colorReset))
		return nil
	}

	if config.Branch != "release" {
		fmt.Println(string(colorYellow), "⚠️ You are running a development build. Skipping update check.", string(colorReset))
		return nil
	}

	// Fetch latest release from GitHub
	fmt.Println(string(colorCyan), "🕵️‍♂️ Querying GitHub API for the latest release...", string(colorReset))
	latestRelease, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("❌ Failed to fetch latest release: %v", err)
	}

	// Determine expected executable name based on platform
	expectedExt := ".exe"
	if runtime.GOOS != "windows" {
		expectedExt = ".x86_64"
	}
	expectedExe := fmt.Sprintf("StationeersServerControl%s%s", latestRelease.TagName, expectedExt)
	fmt.Println(string(colorCyan), "🕵️‍♂️ Expected executable name:", expectedExe, string(colorReset))

	// Check if we're already up-to-date
	if currentExe == expectedExe {
		fmt.Println(string(colorGreen), "🎉 You’re already rocking the latest version:", currentExe, string(colorReset))
		return nil
	}

	// Find the matching asset in the release
	var downloadURL string
	for _, asset := range latestRelease.Assets {
		if asset.Name == expectedExe {
			downloadURL = asset.URL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("❌ No matching asset found for %s in latest release", expectedExe)
	}

	// Download the new executable
	fmt.Println(string(colorCyan), "📡 Found a newer version! Downloading", expectedExe, "...", string(colorReset))
	if err := downloadNewExecutable(expectedExe, downloadURL); err != nil {
		return fmt.Errorf("❌ Failed to download new executable: %v", err)
	}

	// Set executable permissions on Linux
	if runtime.GOOS != "windows" {
		if err := os.Chmod(expectedExe, 0755); err != nil {
			return fmt.Errorf("❌ Couldn’t make %s executable: %v", expectedExe, err)
		}
	}

	// Launch the new executable and exit
	fmt.Println(string(colorMagenta), "🚀 Launching the new version and retiring the old one...", string(colorReset))
	if err := runAndExit(expectedExe); err != nil {
		return fmt.Errorf("❌ Couldn’t switch to new executable: %v", err)
	}

	return nil
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

	fmt.Printf("\n✅ Downloaded %s like a champ!\n", filename)
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
	fmt.Println(string(colorYellow), "✨ New version’s live! Catch you on the flip side!", string(colorReset))
	time.Sleep(500 * time.Millisecond) // Dramatic pause
	os.Exit(0)
	return nil
}
