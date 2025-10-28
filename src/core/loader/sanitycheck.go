package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

var containerCheckWG = sync.WaitGroup{}

func runSanityCheck() error {

	if runtime.GOOS == "windows" {
		return nil
	}

	if config.GetNoSanityCheck() {
		return nil
	}

	IsInsideContainer(&containerCheckWG)
	containerCheckWG.Wait()

	// Check if running as root (UID 0)
	if os.Geteuid() == 0 {
		// Check if running inside a container
		if !config.GetIsDockerContainer() {
			return fmt.Errorf("root: SSUI should not be run as root")
		}
	}

	err := checkAndAdjustWorkdir()

	// Check if current working directory is writable
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Try to create a temporary file to test write permissions
	testFile := filepath.Join(workDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("cannot write to working directory, please make sure your user has write permissions in %s: %w", workDir, err)
	}
	// Clean up test file
	if err := os.Remove(testFile); err != nil {
		return fmt.Errorf("failed to clean up sanity check writetest file: %w", err)
	}

	// Check if steamcmd package is installed  (requires further testing, disabled for now)
	//cmd := exec.Command("dpkg-query", "-W", "-f='${Status}'", "steamcmd")
	//output, err := cmd.CombinedOutput()
	//if err == nil && strings.Contains(string(output), "install ok installed") {
	//	return fmt.Errorf("steamcmd apt package is installed, it is not recommended to run SSUI when the apt steamcmd package is installed. Please uninstall the steamcmd package or have a look at our Docker image and try again")
	//}

	return nil
}

func checkAndAdjustWorkdir() error {

	customWorkDir := SetCustomWorkDirFlag

	if customWorkDir != "" {
		logger.Core.Info("Switching to Custom WorkDir: " + customWorkDir)
		return os.Chdir(customWorkDir)
	}

	exePath, _ := os.Readlink("/proc/self/exe")
	if exePath == "" {
		return fmt.Errorf("failed to get exepath from /proc/self/exe")
	}
	exeDir := filepath.Dir(exePath)

	cwd, _ := os.Getwd()
	if cwd != exeDir {
		logger.Core.Info("Switching to executable directory: " + exeDir)
		return os.Chdir(exeDir)
	}
	return nil
}
