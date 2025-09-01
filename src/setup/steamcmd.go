package setup

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

// ExtractorFunc is a type that represents a function for extracting archives.
// It takes an io.ReaderAt, the size of the content, and the destination directory.
type ExtractorFunc func(io.ReaderAt, int64, string) error

// Constants for repeated strings
const (
	SteamCMDLinuxURL   = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"
	SteamCMDWindowsURL = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
	SteamCMDLinuxDir   = "./steamcmd"
	SteamCMDWindowsDir = "C:\\SteamCMD"
)

// InstallAndRunSteamCMD installs and runs SteamCMD based on the platform (Windows/Linux).
// It returns the exit status of the SteamCMD execution and any error encountered.
func InstallAndRunSteamCMD() (int, error) {
	if config.GetBranch() == "indev-no-steamcmd" || config.GetIsDebugMode() {
		logger.Install.Info("🔍 Detected indev-no-steamcmd branch or debug=true, skipping SteamCMD run")
		return 0, nil
	}

	if runtime.GOOS == "windows" {
		return installSteamCMDWindows()
	} else if runtime.GOOS == "linux" {
		return installSteamCMDLinux()
	} else {
		err := fmt.Errorf("SteamCMD installation is not supported on this OS")
		logger.Install.Error("❌ " + err.Error() + "\n")
		return -1, err
	}
}

func installSteamCMD(platform string, steamCMDDir string, downloadURL string, extractFunc ExtractorFunc) (int, error) {
	// Check if SteamCMD is already installed
	if _, err := os.Stat(steamCMDDir); os.IsNotExist(err) {
		logger.Install.Warn("⚠️ SteamCMD not found for " + platform + ", downloading...\n")

		// Create SteamCMD directory
		if err := createSteamCMDDirectory(steamCMDDir); err != nil {
			logger.Install.Error("❌ Error creating SteamCMD directory: " + err.Error() + "\n")
			return -1, err
		}

		// Ensure cleanup on failure
		success := false
		defer func() {
			if !success {
				logger.Install.Warn("⚠️ Cleaning up due to failure...\n")
				os.RemoveAll(steamCMDDir)
			}
		}()

		// Install required libraries
		if err := installRequiredLibraries(); err != nil {
			logger.Install.Error("❌ Error installing required libraries: " + err.Error() + "\n")
			return -1, err
		}

		// Download and extract SteamCMD
		if err := downloadAndExtractSteamCMD(downloadURL, steamCMDDir, extractFunc); err != nil {
			logger.Install.Error("❌ " + err.Error() + "\n")
			return -1, err
		}

		// Set executable permissions for SteamCMD files
		if err := setExecutablePermissions(steamCMDDir); err != nil {
			logger.Install.Error("❌ Error setting executable permissions: " + err.Error() + "\n")
			return -1, err
		}

		// Verify the steamcmd binary
		if err := verifySteamCMDBinary(steamCMDDir); err != nil {
			logger.Install.Error("❌ " + err.Error() + "\n")
			return -1, err
		}

		// Mark installation as successful
		success = true
		logger.Install.Info("✅ SteamCMD installed successfully.\n")
	} else {
		logger.Install.Info("✅ SteamCMD is already installed.")
	}

	// Run SteamCMD and return its exit status and error
	return runSteamCMD(steamCMDDir)
}

// installSteamCMDLinux downloads and installs SteamCMD on Linux.
func installSteamCMDLinux() (int, error) {
	return installSteamCMD("Linux", SteamCMDLinuxDir, SteamCMDLinuxURL, untarWrapper)
}

// installSteamCMDWindows downloads and installs SteamCMD on Windows.
func installSteamCMDWindows() (int, error) {
	return installSteamCMD("Windows", SteamCMDWindowsDir, SteamCMDWindowsURL, unzip)
}

// runSteamCMD runs the SteamCMD command to update the game and returns its exit status and any error.
func runSteamCMD(steamCMDDir string) (int, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Install.Error("❌ Error getting current working directory: " + err.Error() + "\n")
		return -1, err
	}
	logger.Install.Debug("✅ Current working directory: " + currentDir + "\n")

	// Ensure permissions every time if we run on linux
	if runtime.GOOS != "windows" {
		if err := setExecutablePermissions(steamCMDDir); err != nil {
			logger.Install.Error("❌ Error setting executable permissions, your Steamcmd install might be broken: " + err.Error() + "\n")
			return -1, err
		}
	}

	// Build SteamCMD command
	cmd := buildSteamCMDCommand(steamCMDDir, currentDir)

	// Set output to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if config.GetLogLevel() == 10 {
		cmdString := strings.Join(cmd.Args, " ")
		logger.Install.Info("🕑 Running SteamCMD: " + cmdString)
	} else {
		logger.Install.Info("🕑 Running SteamCMD...")
	}
	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			logger.Install.Error("❌ SteamCMD exited unsuccessfully: " + err.Error() + "\n")
			return exitErr.ExitCode(), err
		}
		logger.Install.Error("❌ Error running SteamCMD: " + err.Error() + "\n")
		return -1, err
	}
	logger.Install.Info("✅ SteamCMD executed successfully.\n")
	return 0, nil
}

// buildSteamCMDCommand constructs the SteamCMD command based on the OS.
func buildSteamCMDCommand(steamCMDDir, currentDir string) *exec.Cmd {
	//print the config.GameBranch and config.GameServerAppID
	logger.Install.Info("🔍 Game Branch: " + config.GetBranch())
	logger.Install.Debug("🔍 Game Server App ID: " + config.GetGameServerAppID())

	if runtime.GOOS == "windows" {
		return exec.Command(filepath.Join(steamCMDDir, "steamcmd.exe"), "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", config.GetGameServerAppID(), "-beta", config.GetGameBranch(), "validate", "+quit")
	}
	return exec.Command(filepath.Join(steamCMDDir, "steamcmd.sh"), "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", config.GetGameServerAppID(), "-beta", config.GetGameBranch(), "validate", "+quit")
}
