// start.go
package web

import (
	"log"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/security"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

type WebServerLogger struct{}

func (cl *WebServerLogger) Write(p []byte) (n int, err error) {
	// Redirect HTTP server logs (like TLS handshake errors) to logger.API.Debug
	logger.API.Debug(string(p))
	return len(p), nil
}

func StartWebServer(wg *sync.WaitGroup) {

	logger.Web.Info("Starting API services...")
	mux, protectedMux := SetupRoutes()

	// Apply middleware only to protected routes
	mux.Handle("/", AuthMiddleware(protectedMux)) // Wrap protected routes under root

	httpLogger := log.New(&WebServerLogger{}, "", 0)
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
			Addr:     "0.0.0.0:" + config.GetSSUIWebPort(),
			Handler:  mux,
			ErrorLog: httpLogger,
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
			logger.Web.Warn("⚠️Starting pprof server on :6060/debug/pprof")
			err := http.ListenAndServe("0.0.0.0:6060", pprofMux)
			if err != nil {
				logger.Web.Error("Error starting pprof server: " + err.Error())
			}
		})
	}
}
