package web

import (
	"StationeersServerUI/src/backupmgr"
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/configchanger"
	"StationeersServerUI/src/detectionmgr"
	"StationeersServerUI/src/security"
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

func StartWebServer(wg *sync.WaitGroup) {
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
	mux.HandleFunc("/auth/login", security.LoginHandler) // Token issuer
	mux.HandleFunc("/auth/logout", security.LogoutHandler)

	// Protected routes (wrapped with middleware)
	protectedMux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./UIMod"))
	protectedMux.Handle("/static/", http.StripPrefix("/static/", fs))
	protectedMux.HandleFunc("/config", ServeConfigPage)
	protectedMux.HandleFunc("/detectionmanager", ServeDetectionManager)
	protectedMux.HandleFunc("/", ServeIndex)

	// v1 API routes
	protectedMux.HandleFunc("/start", StartServer)
	protectedMux.HandleFunc("/stop", StopServer)

	protectedMux.HandleFunc("/saveconfigasjson", configchanger.SaveConfigForm)

	// SSE routes
	protectedMux.HandleFunc("/console", GetLogOutput)
	protectedMux.HandleFunc("/events", GetEventOutput)

	// v2 API routes that will eventually replace the API routes above, but for now we'll keep v1 (no prefix) for compatibility.
	// Server Control
	protectedMux.HandleFunc("/api/v2/server/start", StartServer)
	protectedMux.HandleFunc("/api/v2/server/stop", StopServer)

	backupHandler := backupmgr.NewHTTPHandler(backupmgr.GlobalBackupManager)
	protectedMux.HandleFunc("/api/v2/backups", backupHandler.ListBackupsHandler)
	protectedMux.HandleFunc("/api/v2/backups/restore", backupHandler.RestoreBackupHandler)

	// Configuration
	protectedMux.HandleFunc("/api/v2/saveconfig", configchanger.SaveConfigRestful)

	// Custom Detections
	protectedMux.HandleFunc("/api/v2/custom-detections", detectionmgr.HandleCustomDetection)
	protectedMux.HandleFunc("/api/v2/custom-detections/delete/", detectionmgr.HandleDeleteCustomDetection)

	// Apply middleware only to protected routes
	mux.Handle("/", security.AuthMiddleware(protectedMux)) // Wrap protected routes under root

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
		if err := security.EnsureTLSCerts(); err != nil {
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
