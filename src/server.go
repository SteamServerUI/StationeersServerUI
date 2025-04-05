package main

import (
	"StationeersServerUI/src/loader"
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

	fmt.Println(string(colorCyan), "Starting setup...", string(colorReset))

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
