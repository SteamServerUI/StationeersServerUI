// start.go
package web

import (
	"StationeersServerUI/src/backupmgr"
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/configchanger"
	"StationeersServerUI/src/detectionmgr"
	"StationeersServerUI/src/logger"
	"StationeersServerUI/src/security"
	"net/http"
	"net/http/pprof"
	"sync"
)

func StartWebServer(wg *sync.WaitGroup) {

	logger.Web.Info("Starting API services...")
	// Set up handlers with auth middleware
	mux := http.NewServeMux() // Use a mux to apply middleware globally

	// Unprotected auth routes
	mux.HandleFunc("/login/login.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "./UIMod/login/login.js")
	})
	mux.HandleFunc("/login/login.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./UIMod/login/login.css")
	})
	mux.HandleFunc("/auth/login", LoginHandler) // Token issuer
	mux.HandleFunc("/auth/logout", LogoutHandler)
	mux.HandleFunc("/login", ServeLoginTemplate)
	if config.IsFirstTimeSetup {
		mux.HandleFunc("/api/v2/auth/setup/register", RegisterUserHandler) // user registration
		mux.HandleFunc("/api/v2/auth/setup/finalize", SetupFinalizeHandler)
		mux.HandleFunc("/setup", ServeLoginTemplate)
	}

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

	// Authentication
	protectedMux.HandleFunc("/changeuser", ServeLoginTemplate)
	if !config.IsFirstTimeSetup {
		protectedMux.HandleFunc("/api/v2/auth/adduser", RegisterUserHandler) // user registration and change password
	}

	// Apply middleware only to protected routes
	mux.Handle("/", AuthMiddleware(protectedMux)) // Wrap protected routes under root

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Web.Info("Starting the HTTP server on port 8443...")
		logger.Web.Info("UI available at: https://0.0.0.0:8443 or https://localhost:8443")
		if config.IsFirstTimeSetup {
			logger.Web.Warn("For first time Setup, follow the instructions on:")
			logger.Web.Warn("https://github.com/JacksonTheMaster/StationeersServerUI/wiki/First-Time-Setup")
			logger.Web.Warn("Or just copy your save folder to /Saves and edit the save file name from the UI (Config Page)")
		}
		// Ensure TLS certs are ready
		if err := security.EnsureTLSCerts(); err != nil {
			logger.Web.Error("Error setting up TLS certificates: " + err.Error())
			//os.Exit(1)
		}
		err := http.ListenAndServeTLS("0.0.0.0:8443", config.TLSCertPath, config.TLSKeyPath, mux)
		if err != nil {
			logger.Web.Error("Error starting HTTPS server: " + err.Error())
		}
	}()

	// Start the pprof server if debug mode is enabled (HTTP/1.1)
	if config.IsDebugMode && config.LogLevel < 20 { // if debug mode is enabled and log level is lower than 20 (if this triggers LogLevel is probably 10 and probably debug, but who knows), start pprof server
		wg.Add(1)
		go func() {
			defer wg.Done()
			pprofMux := http.NewServeMux()
			// Register pprof handler
			pprofMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			logger.Web.Warn("⚠️Starting pprof server on :6060/debug/pprof")
			err := http.ListenAndServe("0.0.0.0:6060", pprofMux)
			if err != nil {
				logger.Web.Error("Error starting pprof server: " + err.Error())
			}
		}()
	}

	// Wait for both servers to be running
	wg.Wait()

}
