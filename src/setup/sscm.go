package setup

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

var installMutex sync.Mutex

func CheckAndDownloadSSCM() {
	SCCMPluginDir := config.SCCMPluginDir

	requiredDirs := []string{SCCMPluginDir}

	// Set branch
	if config.Branch == "release" || config.Branch == "Release" {
		downloadBranch = "main"
	} else {
		downloadBranch = config.Branch
	}
	logger.Install.Info("Using branch: " + downloadBranch)

	// Define file mappings
	files := map[string]string{
		SCCMPluginDir + "SSCM.dll": fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/SSCM/SSCM.dll", downloadBranch),
		SCCMPluginDir + "SSCM.pdb": fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/SSCM/SSCM.pdb", downloadBranch),
	}

	// Check if the directory exists
	if _, err := os.Stat(SCCMPluginDir); os.IsNotExist(err) {
		logger.Install.Warn("⚠️SCCMdir does not exist. Creating it...")

		// Create directories
		for _, dir := range requiredDirs {
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
		config.ConfigMu.Lock()
		config.IsFirstTimeSetup = true
		config.ConfigMu.Unlock()
		downloadAllFiles(files)
	} else {
		// Directory exists
		config.ConfigMu.Lock()
		config.IsFirstTimeSetup = false
		config.ConfigMu.Unlock()
		logger.Install.Info(fmt.Sprintf("IsUpdateEnabled: %v", config.IsUpdateEnabled))
		logger.Install.Info(fmt.Sprintf("IsFirstTimeSetup: %v", config.IsFirstTimeSetup))
		if config.IsUpdateEnabled {
			logger.Install.Info("🔍Validating SSCM files for updates...")
			if config.Branch == "release" || config.Branch == "Release" {
				downloadBranch = "main"
				updateFilesIfDifferent(files)
			} else {
				downloadBranch = config.Branch
				updateFilesIfDifferent(files)
			}
		} else {
			logger.Install.Info("♻️Folder SSCM already exists. Updates disabled, skipping validation.")
		}
	}
}

func CheckAndInstallBepInEx() {
	// Ensure thread safety
	installMutex.Lock()
	defer installMutex.Unlock()

	logger.Install.Info("Checking for BepInEx installation...")

	// Check if BepInEx is already installed
	if _, err := os.Stat("BepInEx"); err == nil {
		logger.Install.Info("BepInEx folder already exists, skipping installation")
		return
	}

	// Determine the URL based on platform
	var url string
	if runtime.GOOS == "windows" {
		url = "https://github.com/BepInEx/BepInEx/releases/download/v5.4.23.2/BepInEx_win_x64_5.4.23.2.zip"
		logger.Install.Info("Detected Windows platform, using Windows BepInEx package")
	} else {
		url = "https://github.com/BepInEx/BepInEx/releases/download/v5.4.23.2/BepInEx_linux_x64_5.4.23.2.zip"
		logger.Install.Info("Detected non-Windows platform, using Linux BepInEx package")
	}

	// Download and install BepInEx
	if err := downloadAndInstallBepInEx(url); err != nil {
		logger.Install.Error(fmt.Sprintf("❌Failed to install BepInEx: %v", err))
	} else {
		logger.Install.Info("✅BepInEx installed successfully")
	}
}

// downloadAndInstallBepInEx downloads the BepInEx zip and extracts it to the current directory
func downloadAndInstallBepInEx(url string) error {
	// Create a temporary file to store the downloaded zip
	tempFile, err := os.CreateTemp("", "bepinex_*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the zip file when done

	// Download the BepInEx zip file
	logger.Install.Info("📥Downloading BepInEx from: " + url)
	err = downloadFile(tempFile.Name(), url)
	if err != nil {
		return fmt.Errorf("failed to download BepInEx: %w", err)
	}

	// Get file info for the zip
	fileInfo, err := tempFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Reopen the file for reading
	tempFile.Close()
	zipFile, err := os.Open(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open downloaded zip: %w", err)
	}
	defer zipFile.Close()

	// Extract the zip file to the current directory
	logger.Install.Info("📦Extracting BepInEx to current directory")
	err = unzip(zipFile, fileInfo.Size(), ".")
	if err != nil {
		return fmt.Errorf("failed to extract BepInEx: %w", err)
	}

	// Clean up changelog.txt if it exists
	if _, err := os.Stat("changelog.txt"); err == nil {
		logger.Install.Debug("🗑️Removing changelog.txt")
		if err := os.Remove("changelog.txt"); err != nil {
			logger.Install.Warn(fmt.Sprintf("⚠️Failed to remove changelog.txt: %v", err))
		}
	}

	if runtime.GOOS == "linux" {
		// make sure run_bepinex.sh is executable
		if err := os.Chmod("./run_bepinex.sh", os.ModePerm); err != nil {
			logger.Install.Warn(fmt.Sprintf("⚠️Failed to make run_bepinex.sh executable: %v", err))
		}
	}

	return nil
}

func InstallSSCM() {
	logger.Install.Info("🕑Installing SSCM...")

	CheckAndInstallBepInEx()
	CheckAndDownloadSSCM()

	// Enable SSCM
	config.ConfigMu.Lock()
	config.IsSSCMEnabled = true
	config.ConfigMu.Unlock()

	logger.Install.Info("✅SSCM enabled")
}
