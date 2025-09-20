// Package main is the entry point for StationeersServerUI, a tool for managing Stationeers servers.
// It coordinates setup, configuration, logging, resource loading, and a web-based UI.
//
// The server initializes by running the setup process, loading resources, and starting a web server.
// Key functionality is provided by the following subpackages:
//   - src/config: Manages server configuration.
//   - src/configchanger: Handles configuration changes.
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
	"os"
	"sync"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/cli"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/loader"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/web"
)

//go:embed UIMod/onboard_bundled
var v1uiFS embed.FS

func main() {
	var wg sync.WaitGroup
	logger.ConfigureConsole()
	loader.SetupWorkingDir()
	if err := loader.SanityCheck(); err != nil {
		logger.Main.Error("Sanity check failed, exiting in 10 secconds: " + err.Error())
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	logger.Main.Debug("Initializing resources...")
	loader.InitVirtFS(v1uiFS)
	logger.Install.Info("Starting setup...")
	loader.ReloadConfig() // Load the config file before starting the setup process
	loader.LoadCmdArgs()
	setup.Install(&wg)
	wg.Wait()
	logger.Main.Debug("Initializing Backend...")
	loader.InitBackend(&wg)
	wg.Wait()
	logger.Main.Debug("Initializing after start tasks...")
	loader.AfterStartComplete(&wg)
	wg.Wait()
	logger.Main.Debug("Starting webserver...")
	web.StartWebServer(&wg)
	logger.Main.Debug("Initializing SSUICLI...")
	cli.StartConsole(&wg)
	wg.Wait()
}
