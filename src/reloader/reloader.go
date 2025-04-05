// reloader.go
package reloader

import (
	"StationeersServerUI/src/backupmgr"
	"StationeersServerUI/src/config"
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
	// If Discord is enabled, start the Discord bot
	if config.IsDiscordEnabled {
		go discordbot.InitializeDiscordBot()
	}
}
