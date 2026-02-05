package loader

import (
	"fmt"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/cli/dashboard"
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
	// deactivated for now, as we are working on a new way to handle this
	//setup.SetupAutostartScripts()
	discordrpc.StartDiscordRPC()

	go func() {
		time.Sleep(500 * time.Millisecond)
		printStartupMessage()

		if config.GetIsFirstTimeSetup() {
			printFirstTimeSetupMessage()

			// Only prompt if we're in an interactive terminal
			if IsInteractiveTerminal() {
				time.Sleep(1000 * time.Millisecond)
				choice := PromptSetupTrackChoice()
				switch choice {
				case "cli":
					// Run the fancy Tea-based setup wizard
					if err := dashboard.RunSetup(); err != nil {
						logger.Core.Error("Setup wizard error: " + err.Error())
					}
				case "web":
					fmt.Printf("\n  🌐 Web UI available at: https://localhost:8443 (default) or https://<server-ip>:%s\n", config.GetSSUIWebPort())
					fmt.Println("  Open the Web UI in your browser to complete setup.")
					fmt.Println()
				case "skip":
					fmt.Println("\n  ⚠ Setup skipped. You can configure SSUI later via the Web UI")
					fmt.Println()
				}
			}
			// Non-interactive: simply let them use web UI (startup message has URL)
		}
	}()
}
