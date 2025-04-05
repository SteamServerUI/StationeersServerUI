// loader.go
package loader

import (
	"StationeersServerUI/src/backupmgr"
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/detectionmgr"
	"StationeersServerUI/src/discordbot"
	"log"
)

func ReloadAll() {
	ReloadConfig()
	ReloadBackupManager()
	ReloadDiscordBot()
}

func ReloadConfig() {
	if _, err := config.LoadConfig(); err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}
}

func ReloadBackupManager() {
	if err := backupmgr.ReloadBackupManagerFromConfig(); err != nil {
		log.Printf("Failed to reload backup manager: %v", err)
		return
	}
}

func ReloadDiscordBot() {
	if config.IsDiscordEnabled {
		go discordbot.InitializeDiscordBot()
	}
}

// The detector should NOT be reloaded, as it is a singleton. Instead, dynamic changes come in via the custom detections manager.
func InitDetector() {
	detector := detectionmgr.Start()
	detectionmgr.RegisterDefaultHandlers(detector)
	detectionmgr.InitCustomDetectionsManager(detector)
	go detectionmgr.StreamLogs(detector)
}
