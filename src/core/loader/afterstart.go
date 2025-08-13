package loader

import (
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup"
)

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
	if config.AutoStartServerOnStartup {
		logger.Core.Info("AutoStartServerOnStartup is enabled, starting server...")
		gamemgr.InternalStartServer()
	}
	setup.SetupAutostartScripts()
}
