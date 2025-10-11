package configchanger

import (
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

func SaveConfig(cfg *config.JsonConfig, reloadBackend ...bool) error {
	err := config.SaveConfigToFile(cfg)
	if err != nil {
		logger.Core.Error("Failed to save config: " + err.Error())
		return err
	}
	// Call ReloadBackend by default, unless reloadBackend is explicitly false
	if len(reloadBackend) == 0 || reloadBackend[0] {
		loader.ReloadBackend()
	}
	return nil
}
