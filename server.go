// Package main is the entry point for StationeersServerUI, a tool for managing Stationeers servers.
// It coordinates setup, configuration, logging, resource loading, and a web-based UI.
//
// The server initializes by running the setup process, loading resources, and starting a web server.
// Key functionality is provided by the following subpackages:
//   - src/config: Manages server configuration.
//   - src/gamemgr: Manages process management.
//   - src/loader: Handles resource loading and detection initialization.
//   - src/logger: Provides logging utilities.
//   - src/setup: Performs initial server setup.
//   - src/web: Runs the web-based user interface.
//   - src/discordbot: Handles Discord bot functionality.
//   - src/backupmgr: Manages backups of the server's world.
//   - src/detectionmgr: Manages event detection and processing.
//   - src/setup: Performs initial server setup.
//   - src/ssestream: Handles Server-Sent Events (SSE) streaming.
//   - src/security: Handles security-related tasks.
//
// For detailed documentation, see the subpackages or the project Wiki on GitHub.
package main

import (
	"embed"
	"fmt"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/socketapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/cli"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/setup"
)

//go:embed SSUI/onboard_bundled
var v1uiFS embed.FS

func main() {
	var wg sync.WaitGroup
	logger.ConfigureConsole()
	loader.ParseFlags()
	loader.HandleSanityCheckFlag()
	loader.SanityCheck(&wg)
	wg.Wait()
	logger.Main.Info("Initializing resources...")
	loader.InitVirtFS(v1uiFS)
	logger.Install.Info("Starting setup...")
	loader.ReloadConfig() // Load the config file before starting the setup process
	loader.HandleFlags(&wg)
	wg.Wait()
	setup.Install(&wg)
	wg.Wait()
	setupSecret, err := security.PrepareIdentitySetup(time.Now())
	if err != nil {
		logger.Security.Error("Failed to prepare owner setup: " + err.Error())
		return
	}
	if setupSecret != "" {
		fmt.Println("S7UI owner setup secret (valid for 30 minutes): " + setupSecret)
	}
	logger.Main.Debug("Initializing Backend...")
	loader.InitBackend(&wg)
	wg.Wait()
	logger.Main.Debug("Initializing after start tasks...")
	loader.AfterStartComplete(&wg)
	wg.Wait()
	logger.Main.Debug("Starting socket server...")
	socketapi.StartSocketServer(&wg)
	logger.Main.Debug("Starting webserver...")
	api.StartWebServer(&wg)
	logger.Main.Debug("Initializing SSUICLI...")
	cli.StartConsole(&wg)
	wg.Wait()
}
