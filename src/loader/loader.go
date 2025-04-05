// loader.go
package loader

import (
	"StationeersServerUI/src/backupmgr"
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/detectionmgr"
	"StationeersServerUI/src/discordbot"
	"StationeersServerUI/src/logger"
)

func ReloadAll() {
	ReloadConfig()
	ReloadBackupManager()
	ReloadDiscordBot()
}

func ReloadConfig() {
	if _, err := config.LoadConfig(); err != nil {
		logger.Core.Error("Failed to load config: " + err.Error())
		return
	}
	logger.Core.Debug("Config reloaded successfully")
}

func ReloadBackupManager() {
	if err := backupmgr.ReloadBackupManagerFromConfig(); err != nil {
		logger.Backup.Error("Failed to reload backup manager: " + err.Error())
		return
	}
	logger.Backup.Debug("Backup manager reloaded successfully")
}

func ReloadDiscordBot() {
	if config.IsDiscordEnabled {
		go discordbot.InitializeDiscordBot()
		logger.Discord.Debug("Discord bot reloaded successfully")
	}
}

// The detector should NOT be reloaded, as it is a singleton. Instead, dynamic changes come in via the custom detections manager.
func InitDetector() {
	detector := detectionmgr.Start()
	detectionmgr.RegisterDefaultHandlers(detector)
	detectionmgr.InitCustomDetectionsManager(detector)
	go detectionmgr.StreamLogs(detector)
	logger.Detection.Debug("Detector loaded successfully")
}
