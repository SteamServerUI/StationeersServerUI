// start.go
package api

import (
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/middleware"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

type APIServerLogger struct{}

func (cl *APIServerLogger) Write(p []byte) (n int, err error) {
	// Redirect HTTP server logs (like TLS handshake errors) to logger.API.Debug
	logger.API.Debug(string(p))
	return len(p), nil
}

func StartWebServer(wg *sync.WaitGroup) {

	logger.Web.Info("Starting API services...")
	mux, protectedMux := SetupV7APIRoutes()

	// if debug mode, add socket only routes to http api for easy testing
	if config.IsDebugMode {
		SetupSocketAPIRoutes(protectedMux)
	}

	// Apply middleware only to protected routes
	mux.Handle("/", middleware.IdentityMiddleware(protectedMux))

	httpLogger := log.New(&APIServerLogger{}, "", 0)
	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Ensure TLS certs are ready
		if err := security.EnsureTLSCerts(); err != nil {
			logger.Web.Error("Error setting up TLS certificates: " + err.Error())
			return
		}

		// Create an HTTP server with a custom logger
		server := &http.Server{
			Addr:              "0.0.0.0:" + config.GetBackendEndpointPort(),
			Handler:           mux,
			ErrorLog:          httpLogger,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			IdleTimeout:       2 * time.Minute,
			MaxHeaderBytes:    1 << 20,
		}

		err := server.ListenAndServeTLS(config.GetTLSCertPath(), config.GetTLSKeyPath())
		if err != nil {
			logger.Web.Error("Error starting HTTPS server: " + err.Error())
		}
	}()

	// Start the pprof server if debug mode is enabled (HTTP/1.1)
	if config.GetIsDebugMode() { // if debug mode is enabled, start pprof server
		wg.Go(func() {
			pprofMux := http.NewServeMux()
			// Register pprof handler
			pprofMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			logger.Web.Warn("Starting local pprof server on 127.0.0.1:6060/debug/pprof")
			err := http.ListenAndServe("127.0.0.1:6060", pprofMux)
			if err != nil {
				logger.Web.Error("Error starting pprof server: " + err.Error())
			}
		})
	}
}
