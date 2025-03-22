package install

import (
	"StationeersServerUI/src/config"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// Color codes for terminal
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Verbose mode flag - true if not Release branch
var verbose = config.Branch != "Release"

// logVerbose prints a message only if verbose mode is enabled.
func logVerbose(message string) {
	if verbose {
		fmt.Print(message)
	}
}

// logError prints an error message regardless of verbose mode.
func logError(message string) {
	fmt.Print(ColorRed + message + ColorReset)
}

// logSuccess prints a success message only if verbose mode is enabled.
func logSuccess(message string) {
	if verbose {
		fmt.Print(ColorGreen + message + ColorReset)
	}
}

// InstallAndRunSteamCMD installs and runs SteamCMD based on the platform (Windows/Linux).
// It automatically detects the OS and calls the appropriate installation function.
func InstallAndRunSteamCMD() {
	if runtime.GOOS == "windows" {
		installSteamCMDWindows()
	} else if runtime.GOOS == "linux" {
		installSteamCMDLinux()
	} else {
		logError("❌ SteamCMD installation is not supported on this OS.\n")
		return
	}
}

func installSteamCMD(platform string, steamCMDDir string, downloadURL string, extractFunc ExtractorFunc) {
	// Check if SteamCMD is already installed
	if _, err := os.Stat(steamCMDDir); os.IsNotExist(err) {
		logVerbose(ColorYellow + "⚠️ SteamCMD not found for " + platform + ", downloading...\n" + ColorReset)

		// Create SteamCMD directory
		if err := createSteamCMDDirectory(steamCMDDir); err != nil {
			logError("❌ Error creating SteamCMD directory: " + err.Error() + "\n")
			return
		}

		// Ensure cleanup on failure
		success := false
		defer func() {
			if !success {
				logVerbose(ColorYellow + "⚠️ Cleaning up due to failure...\n" + ColorReset)
				os.RemoveAll(steamCMDDir)
			}
		}()

		// Install required libraries
		if err := installRequiredLibraries(); err != nil {
			logError("❌ Error installing required libraries: " + err.Error() + "\n")
			return
		}

		// Download and extract SteamCMD
		if err := downloadAndExtractSteamCMD(downloadURL, steamCMDDir, extractFunc); err != nil {
			logError("❌ " + err.Error() + "\n")
			return
		}

		// Set executable permissions for SteamCMD files
		if err := setExecutablePermissions(steamCMDDir); err != nil {
			logError("❌ Error setting executable permissions: " + err.Error() + "\n")
			return
		}

		// Verify the steamcmd binary
		if err := verifySteamCMDBinary(steamCMDDir); err != nil {
			logError("❌ " + err.Error() + "\n")
			return
		}

		// Mark installation as successful
		success = true
		logSuccess("✅ SteamCMD installed successfully.\n")
	} else {
		logVerbose("✅ SteamCMD is already installed.\n")
	}

	// Run SteamCMD
	runSteamCMD(steamCMDDir)
}

// installSteamCMDLinux downloads and installs SteamCMD on Linux.
func installSteamCMDLinux() {
	installSteamCMD("Linux", SteamCMDLinuxDir, SteamCMDLinuxURL, untarWrapper)
}

// installSteamCMDWindows downloads and installs SteamCMD on Windows.
func installSteamCMDWindows() {
	installSteamCMD("Windows", SteamCMDWindowsDir, SteamCMDWindowsURL, unzip)
}

// runSteamCMD runs the SteamCMD command to update the game.
func runSteamCMD(steamCMDDir string) {
	currentDir, err := os.Getwd()
	if err != nil {
		logError("❌ Error getting current working directory: " + err.Error() + "\n")
		return
	}
	logVerbose("✅ Current working directory: " + currentDir + "\n")

	// Build SteamCMD command
	cmd := buildSteamCMDCommand(steamCMDDir, currentDir)

	// Set output to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	logVerbose(ColorBlue + "🕑 Running SteamCMD...\n" + ColorReset)
	err = cmd.Run()
	if err != nil {
		logError("❌ Error running SteamCMD: " + err.Error() + "\n")
		return
	}
	logSuccess("✅ SteamCMD executed successfully.\n")
}

// buildSteamCMDCommand constructs the SteamCMD command based on the OS.
func buildSteamCMDCommand(steamCMDDir, currentDir string) *exec.Cmd {
	//print the config.GameBranch and config.GameServerAppID
	fmt.Println(config.GameBranch)
	fmt.Println(config.GameServerAppID)
	if runtime.GOOS == "windows" {
		return exec.Command(filepath.Join(steamCMDDir, "steamcmd.exe"), "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", config.GameServerAppID, "-beta", config.GameBranch, "validate", "+quit")
	}
	return exec.Command(filepath.Join(steamCMDDir, "steamcmd.sh"), "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", config.GameServerAppID, "-beta", config.GameBranch, "validate", "+quit")
	//print the steamcmd command

}
