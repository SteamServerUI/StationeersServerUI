package loader

import (
	"fmt"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/steamserverui/runfile"
)

// used via Runfile Gallery
func InitRunfile(game string) error {
	// Validate input
	game = strings.TrimSpace(game)
	if game == "" {
		return fmt.Errorf("game cannot be empty")
	}

	logger.Runfile.Info("Updating runfile game to " + game)
	logger.Runfile.Info("Stopping server if running")
	gamemgr.InternalStopServer()
	config.SetRunfileGame(game)

	if err := ReloadRunfile(); err != nil {
		return err
	}

	logger.Runfile.Info("Running SteamCMD, this may take a while...")
	//steammgr.RunSteamCMD()
	logger.Runfile.Warn("Steamcmd for runfile not implemented yet")
	logger.Runfile.Info("Runfile game updated to " + game)

	return nil
}

// used to only reload runfile into memory. Can be triggered from v1 UI -> Runfile Reset terminal
func ReloadRunfile() error {
	if err := runfile.LoadRunfile(config.GetRunfileGame(), config.GetRunFilesFolder()); err != nil {
		logger.Runfile.Warn("Failed to reload runfile: " + err.Error())
		return err
	}
	logger.Runfile.Info("Runfile reloaded successfully")
	return nil
}
