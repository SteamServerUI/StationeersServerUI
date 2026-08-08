//go:build windows
// +build windows

package socketapi

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/google/uuid"
	"github.com/microsoft/go-winio"
)

var pipePath string

func StartSocketServer(wg *sync.WaitGroup) {

	pipePath = getPipePath()
	logger.Socket.Info("Starting named pipe server...")

	// Set up routes
	mux, httpAPIMux := api.SetupV7APIRoutes()
	api.SetupSocketAPIRoutes(httpAPIMux)
	mux.Handle("/", httpAPIMux)

	// Create named pipe listener
	listener, err := winio.ListenPipe(pipePath, nil)
	if err != nil {
		logger.Socket.Error("Error starting named pipe server: " + err.Error())
		return
	}

	// Create HTTP server
	server := &http.Server{
		Handler:  mux,
		ErrorLog: log.New(&api.APIServerLogger{}, "", 0),
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

func getPipePath() string {
	uuid := uuid.New()
	identifier := `\\.\pipe\ssui-` + uuid.String() + `\`

	// Get the target directory path
	dirPath := filepath.Join(config.GetSSUIFolder(), "plugins", "sockets")

	// Create the directory if it doesn't exist
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		logger.Socket.Error("Error creating plugin sockets directory, plugins might not work: " + err.Error())
	}

	filePath := filepath.Join(dirPath, "pipename.identifier")
	err = os.WriteFile(filePath, []byte(identifier), 0644)
	if err != nil {
		logger.Socket.Error("Error writing pipename identifier file, plugins might not work: " + err.Error())
	}

	// Return the identifier for the SSUI socket
	return identifier + "ssui"
}
