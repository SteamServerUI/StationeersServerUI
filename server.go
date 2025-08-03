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
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/loader"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/terminal"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/web"
)

//go:embed UIMod/onboard_bundled
var v1uiFS embed.FS

func main() {
	var wg sync.WaitGroup
	logger.Main.Install("Starting setup...")
	loader.ReloadConfig() // Load the config file before starting the setup process
	// Start the installation process and wait for it to complete
	wg.Add(1)
	go setup.Install(&wg)

	// Wait for the installation to finish before starting the rest of the server
	wg.Wait()

	// Load config,discordbot, backupmgr and detectionmgr using the loader package
	loader.InitVirtFS(v1uiFS)
	loader.ReloadAll()
	loader.InitDetector()
	loader.AfterStartComplete()

	terminal.StartConsole(&wg)

	web.StartWebServer(&wg)
}
