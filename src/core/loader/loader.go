// loader.go
package loader

import (
	"embed"
	"os"
	"sync"
	"time"

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

func SanityCheck(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	err := runSanityCheck()
	if err != nil {
		logger.Main.Error("Sanity check failed, exiting in 10 secconds: " + err.Error())
		logger.Main.Info("If you want to continue anyway, run SSUI with the --noSanityCheck flag, but be aware there may be Dragons ahead.")
		logger.Main.Info("This is not recommended nor supported and may cause unexpected behavior, including potential data loss!")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
}
