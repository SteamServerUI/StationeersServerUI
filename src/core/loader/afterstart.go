package loader

import (
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/discordrpc"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup"
)

func AfterStartComplete() {

	config.SetSaveConfig() // Save config after startup through setters
	err := setup.CleanUpOldUIModFolderFiles()
	if err != nil {
		logger.Core.Error("AfterStartComplete: Failed to clean up old pre-v5.5 UI mod folder files: " + err.Error())
	}
	err = setup.CleanUpOldExecutables()
	if err != nil {
		logger.Core.Error("AfterStartComplete: Failed to clean up old executables: " + err.Error())
	}
	if config.GetAutoStartServerOnStartup() {
		logger.Core.Info("AutoStartServerOnStartup is enabled, starting server...")
		gamemgr.InternalStartServer()
	}
	setup.SetupAutostartScripts()
	discordrpc.StartDiscordRPC()
}
