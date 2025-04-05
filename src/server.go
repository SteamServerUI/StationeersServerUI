package main

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/detectionmgr"
	"StationeersServerUI/src/reloader"
	"StationeersServerUI/src/setup"
	"StationeersServerUI/src/web"
	"fmt"
	"sync"
)

const (
	// ANSI color codes for styling terminal output
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

func main() {
	var wg sync.WaitGroup

	fmt.Println(string(colorCyan), "Starting checks...", string(colorReset))

	// Start the installation process and wait for it to complete
	wg.Add(1)
	go setup.Install(&wg)

	// Wait for the installation to finish before starting the rest of the server
	wg.Wait()

	fmt.Println(string(colorGreen), "Setup complete!", string(colorReset))

	fmt.Println(string(colorBlue), "Reloading configuration", string(colorReset))
	config.LoadConfig()

	// Initialize the detection module
	fmt.Println(string(colorBlue), "Initializing detection module...", string(colorReset))
	detector := detectionmgr.Start()
	detectionmgr.RegisterDefaultHandlers(detector)
	detectionmgr.InitCustomDetectionsManager(detector)
	fmt.Println(string(colorGreen), "Detection module ready!", string(colorReset))

	go detectionmgr.StreamLogs(detector) // Pass the detector to the log stream function

	fmt.Println(string(colorBlue), "Starting API services...", string(colorReset))

	// Load dicord, backupmgr and detectionmgr using the reloader package
	reloader.ReloadAll()

	if config.IsCleanupEnabled {
		fmt.Println(string(colorBlue), "v2 Backup cleanup is enabled.", string(colorReset))
	}

	fmt.Println(string(colorBlue), "Global backup manager initialized.", string(colorReset))

	web.StartWebServer(&wg)
}
