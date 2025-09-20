package steamcmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

var steamMu sync.Mutex
var isUpdatingMu sync.Mutex

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
	if isUpdatingMu.TryLock() {
		// Successfully acquired the lock; we are not updating currently
		logger.Core.Debug("🔄 Locking isUpdatingMu for SteamCMD Update run...")
	} else {
		// already updating, return
		logger.Core.Warn("🔄 isUpdatingMu is currently locked, cannot update server using SteamCMD right now...")
		return -1, fmt.Errorf("already updating")
	}
	defer isUpdatingMu.Unlock()
	defer logger.Core.Debug("🔄 Unlocking isUpdatingMu after SteamCMD Update run...")

	if gamemgr.InternalIsServerRunning() {
		logger.Core.Warn("Server is running, stopping server first...")
		err := gamemgr.InternalStopServer()
		if err != nil {
			logger.Core.Error("Error stopping server before running Steamcmd: " + err.Error())
		}
	}
	logger.Core.Info("Running SteamCMD")

	switch runtime.GOOS {
	case "windows":
		return installSteamCMDWindows()
	case "linux":
		return installSteamCMDLinux()
	default:
		err := fmt.Errorf("SteamCMD installation is not supported on this OS")
		logger.Install.Error("❌ " + err.Error() + "\n")
		return -1, err
	}
}

// runSteamCMD runs the SteamCMD command to update the game and returns its exit status and any error.
func runSteamCMD(steamCMDDir string) (int, error) {
	if steamMu.TryLock() {
		// Successfully acquired the lock; no other func holds it
		logger.Core.Debug("🔄 Locking SteamMu for SteamCMD execution...")
	} else {
		// Another goroutine holds the lock; log and wait.
		logger.Core.Warn("🔄 SteamMu is currently locked, waiting for it to be unlocked and then continuing...")
		steamMu.Lock() // Block until steamMu becomes available, then snack it and lock it again
		logger.Core.Debug("🔄 Locking SteamMu for SteamCMD execution..")
	}
	defer steamMu.Unlock()
	defer logger.Core.Debug("🔄 Unlocking SteamMu after SteamCMD execution...")
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

	if runtime.GOOS == "linux" {
		env := os.Environ()
		// Replace or set HOME
		newEnv := make([]string, 0, len(env)+1)
		foundHome := false
		for _, e := range env {
			if !strings.HasPrefix(e, "HOME=") {
				newEnv = append(newEnv, e)
			} else {
				newEnv = append(newEnv, "HOME="+currentDir)
				foundHome = true
			}
		}
		if !foundHome {
			newEnv = append(newEnv, "HOME="+currentDir)
		}
		cmd.Env = newEnv
	}

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
	logger.Install.Info("🔍 Game Branch: " + config.GetGameBranch())
	logger.Install.Debug("🔍 Game Server App ID: " + config.GetGameServerAppID())

	if runtime.GOOS == "windows" {
		return exec.Command(filepath.Join(steamCMDDir, "steamcmd.exe"), "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", config.GetGameServerAppID(), "-beta", config.GetGameBranch(), "validate", "+quit")
	}
	return exec.Command(filepath.Join(steamCMDDir, "steamcmd.sh"), "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", config.GetGameServerAppID(), "-beta", config.GetGameBranch(), "validate", "+quit")
}
