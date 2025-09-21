package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func runSanityCheck() error {

	if runtime.GOOS == "windows" {
		return nil
	}

	// Check if running as root (UID 0)
	if os.Geteuid() == 0 {
		// Check if running inside a container
		if !IsInsideContainer() {
			return fmt.Errorf("root: SSUI should not be run as root")
		}
	}

	// Get the current executable path from /proc/self/exe
	exePath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return err
	}
	// Get the directory path of the executable
	dirPath := filepath.Dir(exePath)
	// Change the working directory to the executable's directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if cwd != dirPath && !strings.Contains(dirPath, "/tmp") {
		err = os.Chdir(dirPath)
		if err != nil {
			return err
		}
	}

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
	//	return fmt.Errorf("steamcmd package is installed")
	//}

	return nil
}
