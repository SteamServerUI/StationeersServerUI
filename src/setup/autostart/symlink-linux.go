package autostart

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

func SetupBinarySymlink() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("not on linux")
	}
	// Check if we are inside a docker container, if so, don't try to create a symlink
	if IsInsideContainer() {
		logger.Install.Info("Inside a container, skipping symlink creation.")
		return fmt.Errorf("inside a container")
	}

	if err := updateSymlink(); err != nil {
		return err
	}
	return nil
}

func updateSymlink() error {
	// Get the absolute path of the current executable
	executablePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Get the directory of the executable
	executableDir := filepath.Dir(executablePath)

	// Construct the target filename (e.g., StationeersServerControl5.0.x86_64)
	currentExecutableName := "StationeersServerControlv" + config.GetVersion() + ".x86_64"

	// Construct the full path to the target executable
	targetPath := filepath.Join(executableDir, currentExecutableName)

	symlinkPath := filepath.Join(executableDir, "StationeersServerUI.lnk")

	// Remove existing symlink if it exists
	if err := os.Remove(symlinkPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Create the symlink
	return os.Symlink(targetPath, symlinkPath)
}

func IsInsideContainer() bool {
	// Check .dockerenv file (Docker-specific)
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check cgroup (works for Docker and other container runtimes)
	return isContainerFromCGroup()
}

func isContainerFromCGroup() bool {
	file, err := os.Open("/proc/1/cgroup")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check for various container runtime indicators
		if strings.Contains(line, "docker") ||
			strings.Contains(line, "containerd") ||
			strings.Contains(line, "kubepods") ||
			strings.Contains(line, "crio") ||
			strings.Contains(line, "libpod") {
			return true
		}
	}
	return false
}
