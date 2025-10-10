//go:build linux
// +build linux

package socket

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/web"
)

const socketPath = "/tmp/ssui-api.sock"

func StartSocketServer(wg *sync.WaitGroup) {
	logger.Socket.Info("Starting Unix socket server...")

	// Remove existing socket file if it exists
	if err := os.RemoveAll(socketPath); err != nil {
		logger.Socket.Error("Error removing existing socket: " + err.Error())
	}

	// Set up routes
	mux, protectedMux := web.SetupRoutes()
	mux.Handle("/", protectedMux)

	// Create Unix socket listener
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.Socket.Error("Error starting Unix socket server: " + err.Error())
		return
	}

	// Set socket permissions
	if err := os.Chmod(socketPath, 0666); err != nil {
		logger.Socket.Error("Error setting socket permissions: " + err.Error())
	}

	// Create HTTP server
	server := &http.Server{
		Handler:  mux,
		ErrorLog: log.New(&web.WebServerLogger{}, "", 0),
	}

	// Start server in a goroutine
	wg.Go(func() {
		logger.Socket.Info("Unix socket server running at " + socketPath)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger.Socket.Error("Unix socket server error: " + err.Error())
		}
	})

	// Handle graceful shutdown
	go func() {
		<-context.Background().Done()
		logger.Socket.Info("Shutting down Unix socket server...")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Socket.Error("Error shutting down Unix socket server: " + err.Error())
		}
		if err := os.RemoveAll(socketPath); err != nil {
			logger.Socket.Error("Error removing socket file: " + err.Error())
		}
	}()
}
