// start.go
package web

import (
	"io/fs"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config/configchanger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/security"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/backupmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/detectionmgr"
)

func StartWebServer(wg *sync.WaitGroup) {

	logger.Web.Info("Starting API services...")
	// Set up handlers with auth middleware
	mux := http.NewServeMux() // Use a mux to apply middleware globally

	// Unprotected auth routes
	twoboxformAssetsFS, _ := fs.Sub(config.GetV1UIFS(), "UIMod/onboard_bundled/twoboxform")
	mux.Handle("/twoboxform/", http.StripPrefix("/twoboxform/", http.FileServer(http.FS(twoboxformAssetsFS))))
	mux.HandleFunc("/auth/login", LoginHandler) // Token issuer
	mux.HandleFunc("/auth/logout", LogoutHandler)
	mux.HandleFunc("/login", ServeTwoBoxFormTemplate)

	// Protected routes (wrapped with middleware)
	protectedMux := http.NewServeMux()

	legacyAssetsFS, _ := fs.Sub(config.GetV1UIFS(), "UIMod/onboard_bundled/assets")
	protectedMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(legacyAssetsFS))))

	protectedMux.HandleFunc("/config", ServeConfigPage)
	protectedMux.HandleFunc("/detectionmanager", ServeDetectionManager)
	protectedMux.HandleFunc("/", ServeIndex)

	// --- SVELTE UI ---
	protectedMux.HandleFunc("/v2", ServeSvelteUI)
	svelteAssetsFS, _ := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/v2/assets")
	protectedMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(svelteAssetsFS))))
	protectedMux.HandleFunc("/api/v2/loader/reloadbackend", HandleReloadAll)

	// SSE routes
	protectedMux.HandleFunc("/console", GetLogOutput)
	protectedMux.HandleFunc("/events", GetEventOutput)
	protectedMux.HandleFunc("/logs/debug", GetDebugLogOutput)
	protectedMux.HandleFunc("/logs/info", GetInfoLogOutput)
	protectedMux.HandleFunc("/logs/warn", GetWarnLogOutput)
	protectedMux.HandleFunc("/logs/error", GetErrorLogOutput)
	protectedMux.HandleFunc("/logs/backend", GetBackendLogOutput)

	// Server Control
	protectedMux.HandleFunc("/start", StartServer)
	protectedMux.HandleFunc("/stop", StopServer)
	protectedMux.HandleFunc("/api/v2/server/start", StartServer)
	protectedMux.HandleFunc("/api/v2/server/stop", StopServer)
	protectedMux.HandleFunc("/api/v2/server/status", GetGameServerRunState)
	protectedMux.HandleFunc("/api/v2/server/status/connectedplayers", HandleConnectedPlayersList)

	backupHandler := backupmgr.NewHTTPHandler(backupmgr.GlobalBackupManager)
	protectedMux.HandleFunc("/api/v2/backups", backupHandler.ListBackupsHandler)
	protectedMux.HandleFunc("/api/v2/backups/restore", backupHandler.RestoreBackupHandler)

	// Configuration
	protectedMux.HandleFunc("/saveconfigasjson", configchanger.SaveConfigForm)     // legacy, used on config page
	protectedMux.HandleFunc("/api/v2/saveconfig", configchanger.SaveConfigRestful) // used on twoboxform
	protectedMux.HandleFunc("/api/v2/SSCM/run", HandleCommand)                     // Command execution via SSCM (needs to be enable, config.IsSSCMEnabled)
	protectedMux.HandleFunc("/api/v2/SSCM/enabled", HandleIsSSCMEnabled)           // Check if SSCM is enabled
	protectedMux.HandleFunc("/api/v2/steamcmd/run", HandleRunSteamCMD)             // Run SteamCMD

	// Custom Detections
	protectedMux.HandleFunc("/api/v2/custom-detections", detectionmgr.HandleCustomDetection)
	protectedMux.HandleFunc("/api/v2/custom-detections/delete/", detectionmgr.HandleDeleteCustomDetection)
	// Authentication
	protectedMux.HandleFunc("/changeuser", ServeTwoBoxFormTemplate)
	protectedMux.HandleFunc("/api/v2/auth/adduser", RegisterUserHandler) // user registration and change password
	protectedMux.HandleFunc("/api/v2/auth/whoami", WhoAmIHandler)

	// Setup
	protectedMux.HandleFunc("/setup", ServeTwoBoxFormTemplate)
	protectedMux.HandleFunc("/api/v2/auth/setup/register", RegisterUserHandler) // user registration
	protectedMux.HandleFunc("/api/v2/auth/setup/finalize", SetupFinalizeHandler)

	// Apply middleware only to protected routes
	mux.Handle("/", AuthMiddleware(protectedMux)) // Wrap protected routes under root

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Ensure TLS certs are ready
		if err := security.EnsureTLSCerts(); err != nil {
			logger.Web.Error("Error setting up TLS certificates: " + err.Error())
			//os.Exit(1)
		}
		err := http.ListenAndServeTLS("0.0.0.0:"+config.GetSSUIWebPort(), config.GetTLSCertPath(), config.GetTLSKeyPath(), mux)
		if err != nil {
			logger.Web.Error("Error starting HTTPS server: " + err.Error())
		}
	}()

	// Start the pprof server if debug mode is enabled (HTTP/1.1)
	if config.GetIsDebugMode() { // if debug mode is enabled, start pprof server
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
}
