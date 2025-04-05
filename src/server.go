package main

import (
	"StationeersServerUI/src/loader"
	"StationeersServerUI/src/logger"
	"StationeersServerUI/src/setup"
	"StationeersServerUI/src/web"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	logger.Main.Install("Starting setup...")

	// Start the installation process and wait for it to complete
	wg.Add(1)
	go setup.Install(&wg)

	// Wait for the installation to finish before starting the rest of the server
	wg.Wait()

	// Load config,discordbot, backupmgr and detectionmgr using the loader package
	loader.ReloadAll()
	loader.InitDetector()

	web.StartWebServer(&wg)
}
