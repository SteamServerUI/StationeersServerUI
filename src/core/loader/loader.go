// loader.go
package loader

import (
	"embed"
	"os"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/discord/discordbot"
	"github.com/SteamServerUI/SteamServerUI/v7/src/localization"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/backupmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/setup"
	"github.com/SteamServerUI/SteamServerUI/v7/src/setup/update"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamcmd"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/plugins"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/telemetry"
)

// only call this once at startup
func InitBackend(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	ReloadConfig()
	ReloadRunfile()
	ReloadBepInEx()
	ReloadBackupMgr()
	ReloadLocalizer()
	ReloadAppInfoPoller()
	ReloadDiscordBot()
	InitDetector()
	StartIsGameServerRunningCheck()
	telemetry.InitTelemetry()
}

// use this to reload backend at runtime
func ReloadBackend() {

	logger.Core.Info("Reloading backend...")
	ReloadConfig()
	ReloadRunfile()
	ReloadBepInEx()
	ReloadBackupMgr()
	ReloadLocalizer()
	ReloadAppInfoPoller()
	PrintConfigDetails()
	plugins.ManagePlugins()
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

func ReloadBepInEx() {
	if config.GetIsBepInExEnabled() {
		setup.CheckAndInstallBepInEx()
		if config.GetIsSSCMEnabled() {
			setup.CheckAndDownloadSSCM()
		}
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

func StartIsGameServerRunningCheck() {
	gamemgr.StartIsGameServerRunningCheck()
}

func ReloadAppInfoPoller() {
	steamcmd.AppInfoPoller()
}

func ReloadBackupMgr() {
	backupmgr.InitBackupMgr()
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
		logger.Main.Info("If you want to continue anyway, run SSUI with the --NoSanityCheck flag, but be aware there may be Dragons ahead.")
		logger.Main.Info("This is not recommended nor supported and may cause unexpected behavior, including potential data loss!")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
}
