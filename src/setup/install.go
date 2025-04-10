package setup

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/loader"
	"StationeersServerUI/src/logger"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var downloadBranch string // Holds the branch to download from

// Install performs the entire installation process and ensures the server waits for it to complete
func Install(wg *sync.WaitGroup) {
	defer wg.Done()             // Signal that installation is complete
	time.Sleep(1 * time.Second) // Small pause for effect

	loader.ReloadConfig()

	// Step 0:  Check for updates
	if err := UpdateExecutable(); err != nil {
		logger.Install.Error("❌Update check went sideways: " + err.Error())
	}

	// Step 1: Check and download the UIMod folder contents
	logger.Install.Info("🔄Checking UIMod folder contents...")
	CheckAndDownloadUIMod()
	logger.Install.Info("✅UIMod folder setup complete.")
	time.Sleep(1 * time.Second)

	// Step 2: Check for Blacklist.txt and create it if it doesn't exist
	logger.Install.Info("🔄Checking for Blacklist.txt...")
	checkAndCreateBlacklist()
	logger.Install.Info("✅Blacklist.txt verified or created.")
	time.Sleep(1 * time.Second)

	// Step 3: Install and run SteamCMD
	logger.Install.Info("🔄Installing and running SteamCMD...")
	InstallAndRunSteamCMD()
	logger.Install.Warn("🙏Thank you for using StationeersServerUI!")
	logger.Install.Info("✅Setup complete!")
}

func CheckAndDownloadUIMod() {
	workingDir := "./UIMod/"
	loginDir := "./UIMod/login/"
	detectionmanagerDir := "./UIMod/detectionmanager/"

	// Set branch
	if config.Branch == "Release" || config.Branch == "release" {
		downloadBranch = "main"
	} else {
		downloadBranch = config.Branch
	}
	logger.Install.Info("Using branch: " + downloadBranch)

	// Define file mappings
	files := map[string]string{
		workingDir + "apiinfo.html":                   fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/apiinfo.html", downloadBranch),
		workingDir + "config.html":                    fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/config.html", downloadBranch),
		workingDir + "config.json":                    fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/config.json", downloadBranch),
		workingDir + "index.html":                     fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/index.html", downloadBranch),
		workingDir + "script.js":                      fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/script.js", downloadBranch),
		workingDir + "stationeers.png":                fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/stationeers.png", downloadBranch),
		workingDir + "style.css":                      fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/style.css", downloadBranch),
		workingDir + "favicon.ico":                    fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/favicon.ico", downloadBranch),
		loginDir + "login.css":                        fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/login/login.css", downloadBranch),
		loginDir + "login.js":                         fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/login/login.js", downloadBranch),
		loginDir + "login.html":                       fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/login/login.html", downloadBranch),
		detectionmanagerDir + "detectionmanager.js":   fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/detectionmanager/detectionmanager.js", downloadBranch),
		detectionmanagerDir + "detectionmanager.html": fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/detectionmanager/detectionmanager.html", downloadBranch),
		detectionmanagerDir + "detectionmanager.css":  fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/detectionmanager/detectionmanager.css", downloadBranch),
	}

	// Check if the directory exists
	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		logger.Install.Warn("⚠️Folder ./UIMod does not exist. Creating it...")

		// Create directories
		for _, dir := range []string{workingDir, loginDir, detectionmanagerDir} {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err := os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					logger.Install.Error("❌Error creating folder: " + err.Error())
					return
				}
				logger.Install.Warn("⚠️Created folder: " + dir)
			}
		}

		// Initial download
		config.IsFirstTimeSetup = true
		downloadAllFiles(files)
	} else {
		// Directory exists
		config.IsFirstTimeSetup = false
		logger.Install.Info(fmt.Sprintf("IsUpdateEnabled: %v", config.IsUpdateEnabled))
		logger.Install.Info(fmt.Sprintf("IsFirstTimeSetup: %v", config.IsFirstTimeSetup))
		if config.IsUpdateEnabled {
			logger.Install.Info("🔍Validating UIMod files for updates...")
			if config.Branch == " Release" || config.Branch == "Release" {
				downloadBranch = "main"
				updateFilesIfDifferent(files)
			} else {
				downloadBranch = config.Branch
				updateFilesIfDifferent(files)
			}
		} else {
			logger.Install.Info("♻️Folder ./UIMod already exists. Updates disabled, skipping validation.")
		}
	}
}

// downloadAllFiles downloads all files in the provided map
func downloadAllFiles(files map[string]string) {
	for filepath, url := range files {
		fileName := filepath[strings.LastIndex(filepath, "/")+1:]
		logger.Install.Info("Downloading " + fileName + "...")
		err := downloadFileWithProgress(filepath, url)
		if err != nil {
			logger.Install.Error("❌Error downloading " + fileName + " (setup may be incomplete): " + err.Error())
			return
		}
		logger.Install.Info("✅Downloaded " + fileName + " successfully from branch " + downloadBranch)
	}
	logger.Install.Info("✅All files downloaded successfully.")
}

// updateFilesIfDifferent checks for differences and updates files if necessary
func updateFilesIfDifferent(files map[string]string) {
	for filepath, url := range files {
		fileName := filepath[strings.LastIndex(filepath, "/")+1:]
		if fileName == "config.json" {
			continue
		}
		// Check if local file exists
		localData, err := os.ReadFile(filepath)
		if err != nil && !os.IsNotExist(err) {
			logger.Install.Error("❌Error reading local file " + fileName + ": " + err.Error())
			continue
		}

		// Fetch remote content
		resp, err := http.Get(url)
		if err != nil {
			logger.Install.Error("❌Error fetching remote file " + fileName + ": " + err.Error())
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Install.Info("Using branch: " + downloadBranch)
			logger.Install.Info("URL: " + url)
			logger.Install.Error("❌Bad status for " + fileName + ": " + resp.Status)
			continue
		}

		remoteData, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Install.Error("❌Error reading remote file " + fileName + ": " + err.Error())
			continue
		}

		// Compare local and remote content
		if os.IsNotExist(err) || !bytes.Equal(localData, remoteData) {
			logger.Install.Info("🔄Updating " + fileName + " due to differences...")
			err = downloadFileWithProgress(filepath, url)
			if err != nil {
				logger.Install.Error("❌Error updating " + fileName + ": " + err.Error())
				continue
			}
			logger.Install.Info("✅Updated " + fileName + " successfully from branch " + downloadBranch)
		} else {
			logger.Install.Info("✅" + fileName + " is up-to-date.")
		}
	}
	logger.Install.Info("✅File validation and update check complete.")
}

// checkAndCreateBlacklist ensures Blacklist.txt exists in the root directory
func checkAndCreateBlacklist() {
	blacklistFile := "./Blacklist.txt"

	// Check if Blacklist.txt exists
	if _, err := os.Stat(blacklistFile); os.IsNotExist(err) {
		// Create Blacklist.txt file with a dummy steamID64 so the gameserver doesn't fail reading this file, as it would not be the expected format if it was empty.
		perm := os.FileMode(0644) // Still works cross-platform
		err := os.WriteFile(blacklistFile, []byte("76561197960265728"), perm)
		if err != nil {
			logger.Install.Error("❌Error creating Blacklist.txt: " + err.Error())
			return
		}

		logger.Install.Info("✅Created Blacklist.txt with dummy steamID64.")
	} else {
		logger.Install.Info("♻️Blacklist.txt already exists. Skipping creation.")
	}
}

// downloadFileWithProgress downloads a file from the given URL and saves it to the given filepath with progress indication
func downloadFileWithProgress(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Get the total size for progress reporting
	size := resp.ContentLength

	// Create a counter for tracking progress
	counter := &writeCounter{
		Total: size,
	}

	// Write the body to file with progress tracking
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

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
