package main

import (
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/src/loader"
	"github.com/JacksonTheMaster/StationeersServerUI/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/src/setup"
	"github.com/JacksonTheMaster/StationeersServerUI/src/web"
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
