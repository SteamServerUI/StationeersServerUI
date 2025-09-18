// loader.go
package loader

import (
	"embed"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/discordbot"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/localization"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/backupmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/detectionmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup/update"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/steamcmd"
)

func SetupWorkingDir() error {
	if runtime.GOOS == "windows" {
		// For now Windows doesn't have symlinking issues so we'll just let is use the current working directory
		return nil
	}
	if runtime.GOOS == "linux" {
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
		if cwd != dirPath {
			logger.Core.Info("Changing working directory to " + dirPath)
			err = os.Chdir(dirPath)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

// only call this once at startup
func InitBackend(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	ReloadConfig()
	ReloadSSCM()
	ReloadBackupManager()
	ReloadLocalizer()
	ReloadAppInfoPoller()
	ReloadDiscordBot()
	InitDetector()
}

// use this to reload backend at runtime
func ReloadBackend() {

	logger.Core.Info("Reloading backend...")
	ReloadConfig()
	ReloadSSCM()
	ReloadBackupManager()
	ReloadLocalizer()
	ReloadAppInfoPoller()
	PrintConfigDetails()
	logger.Core.Info("Backend reload done!")
}

// should ideally not be called standalone, if feasable, call ReloadBackend instead
func ReloadConfig() {
	if _, err := config.LoadConfig(); err != nil {
		logger.Core.Error("Failed to load config: " + err.Error())
		return
	}
	logger.Core.Info("Config loaded successfully")

}

func ReloadSSCM() {
	if config.GetIsSSCMEnabled() {
		setup.InstallSSCM()
	}
}

func ReloadBackupManager() {
	if err := backupmgr.ReloadBackupManagerFromConfig(); err != nil {
		logger.Backup.Error("Failed to reload backup manager: " + err.Error())
		return
	}
	logger.Backup.Info("Backup manager reloaded successfully")
}

func ReloadDiscordBot() {
	if config.GetIsDiscordEnabled() {
		go discordbot.InitializeDiscordBot()
		logger.Discord.Info("Discord bot reloaded successfully")
	}
}

// The detector should NOT be reloaded, as it is a singleton. Instead, dynamic changes come in via the custom detections manager.
func InitDetector() {
	detector := detectionmgr.Start()
	detectionmgr.RegisterDefaultHandlers(detector)
	detectionmgr.InitCustomDetectionsManager(detector)
	go detectionmgr.StreamLogs(detector)
	logger.Detection.Info("Detector loaded successfully")
}

func RestartBackend() {
	update.RestartMySelf()
}

func ReloadLocalizer() {
	localization.ReloadLocalizer()
}

func ReloadAppInfoPoller() {
	steamcmd.AppInfoPoller()
}

// InitBundler initialized the onboard bundled assets for the web UI
func InitVirtFS(v1uiFS embed.FS) {
	config.SetV1UIFS(v1uiFS)
}
