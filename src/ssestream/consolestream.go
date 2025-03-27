package ssestream

import (
	"StationeersServerUI/src/config"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	// consoleClients tracks all connected SSE clients for console logs
	consoleClients   = make(map[chan string]bool)
	consoleClientsMu sync.Mutex
)

// ANSI color codes (mirroring the ones in eventstream.go)
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

// StartConsoleStream creates and returns an HTTP handler for console log SSE streaming
func StartConsoleStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "")

		// Create a new buffered channel for this client (buffer of 2000 events)
		messageChan := make(chan string, 2000)

		// Register this client
		consoleClientsMu.Lock()
		consoleClients[messageChan] = true
		consoleClientsMu.Unlock()

		// Send initial connection event
		_, err := fmt.Fprintf(w, "data: Console log stream connected\n\n")
		if err != nil {
			// Early disconnect; clean up immediately
			consoleClientsMu.Lock()
			delete(consoleClients, messageChan)
			consoleClientsMu.Unlock()
			close(messageChan)
			log.Printf("%s🖥️ [Console/LogStream] ⚠️ Failed to send initial log message: %v%s", colorRed, err, colorReset)
			return
		}
		w.(http.Flusher).Flush()

		// Remove client when connection is closed
		notify := r.Context().Done()
		go func(ch chan string) {
			<-notify
			consoleClientsMu.Lock()
			delete(consoleClients, ch)
			consoleClientsMu.Unlock()
			close(ch)
			log.Printf("%s🖥️ [Console/LogStream] 👋 Client disconnected, channel closed%s", colorYellow, colorReset)
		}(messageChan)

		// Stream events to client
		for msg := range messageChan {
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				// Client likely disconnected; rely on context cancellation to clean up
				log.Printf("%s🖥️ [Console/LogStream] ❌ Failed to send log to client: %v%s", colorRed, err, colorReset)
				return
			}
			w.(http.Flusher).Flush()
		}
	}
}

// BroadcastConsoleOutput sends log to all connected console log clients
func BroadcastConsoleOutput(message string) {
	consoleClientsMu.Lock()
	defer consoleClientsMu.Unlock()

	for client := range consoleClients {
		select {
		case client <- message:
			// Message sent successfully
			if config.IsDebugMode {
				fmt.Println("Broadcasting to console clients:", message)
			}
		case <-time.After(time.Second):
			// Timeout - client might be slow; log and keep the client
			log.Printf("%s🖥️ [Console/LogStream] ⏳ Timeout sending message after 1s: %q%s", colorMagenta, message, colorReset)
		}
	}
}
