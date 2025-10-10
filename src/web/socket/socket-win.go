//go:build windows
// +build windows

package socket

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/web"
	"github.com/microsoft/go-winio"
)

const pipePath = `\\.\pipe\ssui-api`

func StartSocketServer(wg *sync.WaitGroup) {
	logger.Socket.Info("Starting named pipe server...")

	// Set up routes
	mux, protectedMux := web.SetupRoutes()
	mux.Handle("/", protectedMux)

	// Create named pipe listener
	listener, err := winio.ListenPipe(pipePath, nil)
	if err != nil {
		logger.Socket.Error("Error starting named pipe server: " + err.Error())
		return
	}

	// Create HTTP server
	server := &http.Server{
		Handler:  mux,
		ErrorLog: log.New(&web.WebServerLogger{}, "", 0),
	}

	// Start server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Socket.Info("Named pipe server running at " + pipePath)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger.Socket.Error("Named pipe server error: " + err.Error())
		}
	}()

	// Handle graceful shutdown
	go func() {
		<-context.Background().Done()
		logger.Socket.Info("Shutting down named pipe server...")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Socket.Error("Error shutting down named pipe server: " + err.Error())
		}
		listener.Close() // Ensure the pipe is closed
	}()
}
