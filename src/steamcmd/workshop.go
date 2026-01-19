package steamcmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/modding"
)

// DownloadWorkshopItems downloads all installed workshop mods using SteamCMD
func UpdateWorkshopItems() error {
	workshopHandles := modding.GetModWorkshopHandles()
	if len(workshopHandles) == 0 {
		logger.Install.Debug("ℹ️  No workshop items to download")
		return nil
	}

	fmt.Println(workshopHandles)
	return DownloadWorkshopItems(workshopHandles)
}

func DownloadWorkshopItems(workshopHandles []string) error {
	logger.Install.Infof("🔄 Downloading %d workshop items...", len(workshopHandles))

	currentDir, err := os.Getwd()
	if err != nil {
		logger.Install.Error("❌ Error getting current working directory: " + err.Error())
		return err
	}

	// Acquire lock for SteamCMD access
	if steamMu.TryLock() {
		logger.Core.Debug("🔄 Locking SteamMu for SteamCMD Workshop Downloads...")
	} else {
		logger.Core.Warn("🔄 SteamMu is currently locked, waiting for it to be unlocked and then continuing...")
		steamMu.Lock()
		logger.Core.Debug("🔄 Locking SteamMu for SteamCMD Workshop Downloads...")
	}
	defer func() {
		steamMu.Unlock()
		logger.Core.Debug("🔄 Unlocking SteamMu after SteamCMD Workshop Downloads...")
	}()

	steamcmddir := SteamCMDLinuxDir
	executable := "steamcmd.sh"

	if runtime.GOOS == "windows" {
		executable = "steamcmd.exe"
		steamcmddir = SteamCMDWindowsDir
	}

	// Download each workshop item
	for i, appID := range workshopHandles {
		logger.Install.Infof("📦 Downloading workshop item %d/%d: %s", i+1, len(workshopHandles), appID)

		// Build SteamCMD command
		cmd := exec.Command(
			filepath.Join(steamcmddir, executable),
			"+force_install_dir", "../",
			"+login", "anonymous",
			"+workshop_download_item", "544550", appID,
			"validate",
			"+quit",
		)

		// Capture output
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Set up environment for Linux
		if runtime.GOOS == "linux" {
			env := os.Environ()
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
		err := cmd.Run()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				logger.Install.Warnf("⚠️  SteamCMD workshop download failed for %s (code %d): %s", appID, exitErr.ExitCode(), stderr.String())
			} else {
				logger.Install.Warnf("⚠️  Error running SteamCMD for workshop item %s: %s", appID, err.Error())
			}
			continue // Continue with next workshop item even if this one fails
		}

		logger.Install.Debugf("✅ Successfully downloaded workshop item: %s", appID)
	}

	logger.Install.Info("✅ Workshop items download complete")

	// Copy downloaded items to mods directory
	err = copyDownloadedItemsToMods(workshopHandles)
	if err != nil {
		logger.Install.Error("❌ Error copying workshop items to mods directory: " + err.Error())
		return err
	}

	return nil
}

// copyDownloadedItemsToMods copies downloaded workshop items from the Steam directory to ./mods
func copyDownloadedItemsToMods(workshopHandles []string) error {
	// Determine the steam content directory based on OS
	var steamContentDir string
	if runtime.GOOS == "windows" {
		steamContentDir = SteamCMDWindowsDir
		// Windows SteamCMD dir is C:\SteamCMD, so workshop content is at C:\SteamCMD\steamapps\workshop\content\544550
		steamContentDir = filepath.Join(steamContentDir, "steamapps", "workshop", "content", "544550")
	} else {
		// Linux: ./steamapps/workshop/content/544550
		steamContentDir = filepath.Join(".", "steamapps", "workshop", "content", "544550")
	}

	// Ensure mods directory exists
	modsDir := "./mods"
	if err := os.MkdirAll(modsDir, 0755); err != nil {
		return fmt.Errorf("failed to create mods directory: %w", err)
	}

	logger.Install.Infof("📂 Copying %d workshop items to mods directory...", len(workshopHandles))

	// Copy each workshop item
	for i, appID := range workshopHandles {
		logger.Install.Infof("📋 Processing workshop item %d/%d: %s", i+1, len(workshopHandles), appID)

		// Source path: steamapps/workshop/content/544550/{appID}
		srcPath := filepath.Join(steamContentDir, appID)

		// Check if source directory exists
		srcInfo, err := os.Stat(srcPath)
		if err != nil || !srcInfo.IsDir() {
			logger.Install.Errorf("❌ Workshop item not found at expected path: %s (skipping)", srcPath)
			continue
		}

		// Destination path: ./mods/Workshop_{appID}
		destPath := filepath.Join(modsDir, fmt.Sprintf("Workshop_%s", appID))

		// Remove existing destination directory if it exists
		if _, err := os.Stat(destPath); err == nil {
			logger.Install.Debugf("🗑️  Removing existing directory: %s", destPath)
			if err := os.RemoveAll(destPath); err != nil {
				logger.Install.Warnf("⚠️  Failed to remove existing directory %s: %s (continuing anyway)", destPath, err.Error())
			}
		}

		// Copy the entire directory
		if err := copyDir(srcPath, destPath); err != nil {
			logger.Install.Warnf("⚠️  Failed to copy workshop item %s: %s (skipping)", appID, err.Error())
			continue
		}

		logger.Install.Debugf("✅ Successfully copied workshop item to: %s", destPath)
	}

	logger.Install.Info("✅ Workshop items copy complete")
	return nil
}

// copyDir recursively copies a directory from src to dst
func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
