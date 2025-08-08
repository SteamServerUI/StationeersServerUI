// loader.go
package loader

import (
	"embed"
	"fmt"
	"strconv"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/backupmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/detectionmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/discordbot"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/localization"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup"
)

// only call this once at startup
func InitBackend() {
	ReloadConfig()
	ReloadSSCM()
	ReloadBackupManager()
	ReloadLocalizer()
	ReloadDiscordBot()
}

// use this to reload backend at runtime
func ReloadBackend() {

	logger.Core.Info("Reloading backend...")
	ReloadConfig()
	ReloadSSCM()
	ReloadBackupManager()
	ReloadLocalizer()
	PrintConfigDetails()
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
	if config.IsSSCMEnabled {
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
	if config.IsDiscordEnabled {
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

func PrintConfigDetails() {
	logger.Config.Debug("Gameserver config values loaded")
	logger.Config.Debug("---- GENERAL CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("Branch: %s", config.Branch))
	logger.Config.Debug(fmt.Sprintf("GameBranch: %s", config.GameBranch))
	logger.Config.Debug("IsDiscordEnabled: " + strconv.FormatBool(config.IsDiscordEnabled))
	logger.Config.Debug("IsCleanupEnabled: " + strconv.FormatBool(config.IsCleanupEnabled))
	logger.Config.Debug("IsDebugMode (pprof Server): " + strconv.FormatBool(config.IsDebugMode))
	logger.Config.Debug("IsFirstTimeSetup: " + strconv.FormatBool(config.IsFirstTimeSetup))

	logger.Config.Debug("---- DISCORD CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("BlackListFilePath: %s", config.BlackListFilePath))
	logger.Config.Debug(fmt.Sprintf("ConnectionListChannelID: %s", config.ConnectionListChannelID))
	logger.Config.Debug(fmt.Sprintf("ControlChannelID: %s", config.ControlChannelID))
	logger.Config.Debug(fmt.Sprintf("ControlPanelChannelID: %s", config.ControlPanelChannelID))
	logger.Config.Debug(fmt.Sprintf("DiscordCharBufferSize: %d", config.DiscordCharBufferSize))
	logger.Config.Debug(fmt.Sprintf("DiscordToken: %s", config.DiscordToken))
	logger.Config.Debug(fmt.Sprintf("ErrorChannelID: %s", config.ErrorChannelID))
	logger.Config.Debug(fmt.Sprintf("IsDiscordEnabled: %v", config.IsDiscordEnabled))
	logger.Config.Debug(fmt.Sprintf("LogChannelID: %s", config.LogChannelID))
	logger.Config.Debug(fmt.Sprintf("LogMessageBuffer: %s", config.LogMessageBuffer))
	logger.Config.Debug(fmt.Sprintf("SaveChannelID: %s", config.SaveChannelID))
	logger.Config.Debug(fmt.Sprintf("StatusChannelID: %s", config.StatusChannelID))

	logger.Config.Debug("---- BACKUP CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("BackupKeepLastN: %d", config.BackupKeepLastN))
	logger.Config.Debug(fmt.Sprintf("BackupKeepDailyFor: %s", config.BackupKeepDailyFor))
	logger.Config.Debug(fmt.Sprintf("BackupKeepWeeklyFor: %s", config.BackupKeepWeeklyFor))
	logger.Config.Debug(fmt.Sprintf("BackupKeepMonthlyFor: %s", config.BackupKeepMonthlyFor))
	logger.Config.Debug(fmt.Sprintf("BackupCleanupInterval: %s", config.BackupCleanupInterval))
	logger.Config.Debug(fmt.Sprintf("ConfiguredBackupDir: %s", config.ConfiguredBackupDir))
	logger.Config.Debug(fmt.Sprintf("ConfiguredSafeBackupDir: %s", config.ConfiguredSafeBackupDir))
	logger.Config.Debug(fmt.Sprintf("BackupWaitTime: %s", config.BackupWaitTime))

	logger.Config.Debug("---- AUTHENTICATION CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("AuthTokenLifetime: %d", config.AuthTokenLifetime))
	logger.Config.Debug(fmt.Sprintf("JwtKey: %s", config.JwtKey))

	logger.Config.Debug("---- MISC CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("Branch: %s", config.Branch))
	logger.Config.Debug(fmt.Sprintf("GameServerAppID: %s", config.GameServerAppID))
	logger.Config.Debug(fmt.Sprintf("Version: %s", config.Version))
	logger.Config.Debug(fmt.Sprintf("IsNewTerrainAndSaveSystem: %v", config.IsNewTerrainAndSaveSystem))

	logger.Config.Debug("----  UPDATER CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("AllowPrereleaseUpdates: %v", config.AllowPrereleaseUpdates))
	logger.Config.Debug(fmt.Sprintf("AllowMajorUpdates: %v", config.AllowMajorUpdates))
	logger.Config.Debug(fmt.Sprintf("IsUpdateEnabled: %v", config.IsUpdateEnabled))

	logger.Config.Debug("----  SSCM CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("SSCMFilePath: %s", config.SSCMFilePath))
	logger.Config.Debug(fmt.Sprintf("IsSSCMEnabled: %v", config.IsSSCMEnabled))
}

func RestartBackend() {
	setup.RestartMySelf()
}

// InitBundler initialized the onboard bundled assets for the web UI
func InitVirtFS(v1uiFS embed.FS) {
	config.SetV1UIFS(v1uiFS)
}

// this is a Hack, but it works for now. Ideally, move the getter setter logic from SteamServerUI to StationeersServerUI, but not feasible at the moment.
func SaveConfig(cfg *config.JsonConfig, reloadBackend ...bool) error {
	err := config.SaveConfig(cfg)
	if err != nil {
		logger.Core.Error("Failed to save config: " + err.Error())
		return err
	}
	// Call ReloadBackend by default, unless reloadBackend is explicitly false
	if len(reloadBackend) == 0 || reloadBackend[0] {
		ReloadBackend()
	}
	return nil
}

func AfterStartComplete() {
	existingConfig, err := config.LoadConfig()
	if err != nil {
		logger.Core.Error("AfterStartComplete: Failed to Load config: " + err.Error())
	}
	err = SaveConfig(existingConfig, false) // save config, but explicitly DONT reload backend since config is already loaded
	if err != nil {
		logger.Core.Error("AfterStartComplete: Failed to save config: " + err.Error())
	}
	err = setup.CleanUpOldUIModFolderFiles()
	if err != nil {
		logger.Core.Error("AfterStartComplete: Failed to clean up old pre-v5.5 UI mod folder files: " + err.Error())
	}
	err = setup.CleanUpOldExecutables()
	if err != nil {
		logger.Core.Error("AfterStartComplete: Failed to clean up old executables: " + err.Error())
	}
}

func ReloadLocalizer() {
	localization.ReloadLocalizer()
}
