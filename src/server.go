package main

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/core"
	"StationeersServerUI/src/detection"
	"StationeersServerUI/src/discord"
	"StationeersServerUI/src/install"
	"StationeersServerUI/src/ssestream"
	"StationeersServerUI/src/tlsconfig"
	"StationeersServerUI/src/ui"
	"fmt"
	"net/http"
	"net/http/pprof"
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
	go install.Install(&wg)

	// Wait for the installation to finish before starting the rest of the server
	wg.Wait()

	fmt.Println(string(colorGreen), "Setup complete!", string(colorReset))

	fmt.Println(string(colorBlue), "Reloading configuration", string(colorReset))
	config.LoadConfig()

	// Initialize the detection module
	fmt.Println(string(colorBlue), "Initializing detection module...", string(colorReset))
	detector := detection.Start()
	detection.RegisterDefaultHandlers(detector)
	fmt.Println(string(colorGreen), "Detection module ready!", string(colorReset))

	// If Discord is enabled, start the Discord bot
	if config.IsDiscordEnabled {
		fmt.Println(string(colorGreen), "Starting Discord bot...", string(colorReset))
		go discord.StartDiscordBot()
	}

	go detection.StreamLogs(detector) // Pass the detector to the log stream function

	fmt.Println(string(colorBlue), "Starting API services...", string(colorReset))
	go core.StartBackupCleanupRoutine()
	go core.WatchBackupDir()

	// Set up handlers with auth middleware
	mux := http.NewServeMux() // Use a mux to apply middleware globally

	// Unprotected auth routes
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./UIMod/login/login.html")
	})
	mux.HandleFunc("/login/login.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "./UIMod/login/login.js")
	})
	mux.HandleFunc("/login/login.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./UIMod/login/login.css")
	})
	mux.HandleFunc("/auth/login", tlsconfig.LoginHandler) // Token issuer
	mux.HandleFunc("/auth/logout", tlsconfig.LogoutHandler)

	// Protected routes (wrapped with middleware)
	protectedMux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./UIMod"))
	protectedMux.Handle("/static/", http.StripPrefix("/static/", fs))
	protectedMux.HandleFunc("/", ui.ServeIndex)
	protectedMux.HandleFunc("/start", core.StartServer)
	protectedMux.HandleFunc("/stop", core.StopServer)
	protectedMux.HandleFunc("/console", core.GetLogOutput)
	protectedMux.HandleFunc("/backups", core.ListBackups)
	protectedMux.HandleFunc("/restore", core.RestoreBackup)
	protectedMux.HandleFunc("/config", core.HandleConfigJSON)
	protectedMux.HandleFunc("/saveconfigasjson", core.SaveConfigJSON)
	protectedMux.HandleFunc("/events", ssestream.StartDetectionEventStream())

	// Apply middleware only to protected routes
	mux.Handle("/", tlsconfig.AuthMiddleware(protectedMux)) // Wrap protected routes under root

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(string(colorYellow), "Starting the HTTP server on port 8443...", string(colorReset))
		fmt.Println(string(colorGreen), "UI available at: https://0.0.0.0:8443 or https://localhost:8443", string(colorReset))
		if config.IsFirstTimeSetup {
			fmt.Println(string(colorMagenta), "For first time Setup, follow the instructions on:", string(colorReset))
			fmt.Println(string(colorMagenta), "https://github.com/JacksonTheMaster/StationeersServerUI/wiki/First-Time-Setup", string(colorReset))
			fmt.Println(string(colorMagenta), "Or just copy your save folder to /Saves and edit the save file name from the UI (Config Page)", string(colorReset))
		}
		// Ensure TLS certs are ready
		if err := tlsconfig.EnsureTLSCerts(); err != nil {
			fmt.Printf(string(colorRed)+"Error setting up TLS certificates: %v\n"+string(colorReset), err)
			//os.Exit(1)
		}
		err := http.ListenAndServeTLS("0.0.0.0:8443", config.TLSCertPath, config.TLSKeyPath, mux)
		if err != nil {
			fmt.Printf(string(colorRed)+"Error starting HTTPS server: %v\n"+string(colorReset), err)
		}
	}()

	// Start the pprof server if debug mode is enabled (HTTP/1.1)
	if config.IsDebugMode {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pprofMux := http.NewServeMux()
			// Register pprof handler
			pprofMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			fmt.Println(string(colorRed), "⚠️Starting pprof server on :6060/debug/pprof", string(colorReset))
			err := http.ListenAndServe("0.0.0.0:6060", pprofMux)
			if err != nil {
				fmt.Printf(string(colorRed)+"Error starting pprof server: %v\n"+string(colorReset), err)
			}
		}()
	}

	// Wait for both servers to be running
	wg.Wait()
}
