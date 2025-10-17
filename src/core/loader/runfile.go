package loader

import (
	"fmt"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamcmd"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/runfile"
)

// used via Runfile Gallery
func InitRunfile(game string) error {
	// Validate input
	game = strings.TrimSpace(game)
	if game == "" {
		return fmt.Errorf("game cannot be empty")
	}

	logger.Runfile.Debug("Updating runfile game to " + game)
	logger.Runfile.Debug("Stopping server if running")
	gamemgr.InternalStopServer()
	config.SetRunfileIdentifier(game)

	if err := ReloadRunfile(); err != nil {
		return err
	}

	logger.Runfile.Info("Runfile game updated to " + game)
	logger.Runfile.Info("Running SteamCMD, this may take a while...")
	steamcmd.InstallAndRunSteamCMD()

	return nil
}

func ReloadRunfile() error {

	if err := runfile.LoadRunfile(config.GetRunfileIdentifier(), config.GetRunFilesFolder()); err != nil {
		logger.Runfile.Warn("Failed to reload runfile: " + err.Error())
		return err
	}
	logger.Runfile.Info("Runfile reloaded successfully")
	return nil
}
