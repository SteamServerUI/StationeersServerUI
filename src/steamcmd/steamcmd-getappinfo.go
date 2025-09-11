package steamcmd

import (
	"bytes"
	"fmt"
	"maps"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
)

var (
	branches     = make(map[string]string)
	branchesLock sync.RWMutex // Protects branches map for concurrent access
)

func InitAppInfoPoller() {
	go func() {
		for {
			err := getAppInfo()
			if err != nil {
				logger.Install.Error("❌ Failed to get app info: " + err.Error() + "\n")
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}

// getAppInfo fetches the branches and their build IDs for the specified app ID using SteamCMD
// and stores them in the package-level branches map.
func getAppInfo() error {
	steamcmddir := SteamCMDLinuxDir
	executable := "steamcmd.sh"
	appid := config.GetGameServerAppID()

	if runtime.GOOS == "windows" {
		executable = "steamcmd.exe"
		steamcmddir = SteamCMDWindowsDir
	}

	// Build SteamCMD command with +app_info_update to ensure fresh data
	cmd := exec.Command(filepath.Join(steamcmddir, executable), "+login", "anonymous", "+app_info_update", "1", "+app_info_print", appid, "+quit")

	// Capture output instead of printing directly
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Log the command
	if config.GetLogLevel() == 10 {
		cmdString := strings.Join(cmd.Args, " ")
		logger.Install.Debug("🕑 Running SteamCMD for app info: " + cmdString)
	}

	// Run the command
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			logger.Install.Errorf("❌ SteamCMD app info failed (code %d): %s\n", exitErr.ExitCode(), stderr.String())
			return fmt.Errorf("SteamCMD app info failed with exit code %d: %w", exitErr.ExitCode(), err)
		}
		logger.Install.Errorf("❌ Error running SteamCMD app info: %s\n", err.Error())
		return fmt.Errorf("failed to run SteamCMD app info: %w", err)
	}

	// Extract branches and build IDs
	newBranches, err := extractBranches(stdout.String())
	if err != nil {
		logger.Install.Debug("❌ Failed to extract branches: " + err.Error() + "\n")
		return err
	}

	// Update package-level branches map
	branchesLock.Lock()
	maps.Copy(branches, newBranches)
	branchesLock.Unlock()

	currentBranch := config.GetGameBranch()
	if buildID, ok := branches[currentBranch]; ok {
		if config.GetCurrentBranchBuildID() != "" && config.GetCurrentBranchBuildID() != buildID {
			logger.Install.Info("❗New gameserver update detected!")
			if config.GetAllowAutoGameServerUpdates() {
				logger.Install.Info("🔍 Updating gameserver via SteamCMD...")
				if gamemgr.InternalIsServerRunning() {
					gamemgr.InternalStopServer()
				}
				_, err := InstallAndRunSteamCMD()
				if err != nil {
					logger.Install.Error("❌ Failed to update gameserver: " + err.Error() + "\n")
				}
				gamemgr.InternalStartServer()
			}
		}
		config.SetCurrentBranchBuildID(buildID)
	}

	return nil
}

// extractBranches uses regex to extract branch names and their build IDs from SteamCMD output.
func extractBranches(output string) (map[string]string, error) {
	// Regex to match the entire branches section, handling nested braces
	pattern := regexp.MustCompile(`"branches"\s*\{([\s\S]*?)\}\s*\}`)
	branchSection := pattern.FindStringSubmatch(output)
	if len(branchSection) < 2 {
		return nil, fmt.Errorf("branches section not found in output")
	}

	// Regex to match individual branch blocks
	branchPattern := regexp.MustCompile(`"([^"]+)"\s*\{[\s\S]*?"buildid"\s*"(\d+)"[\s\S]*?\}`)
	matches := branchPattern.FindAllStringSubmatch(branchSection[1], -1)

	branches := make(map[string]string)
	for _, match := range matches {
		if len(match) >= 3 {
			branches[match[1]] = match[2]
		}
	}

	if len(branches) == 0 {
		return nil, fmt.Errorf("no branches with build IDs found")
	}

	return branches, nil
}
