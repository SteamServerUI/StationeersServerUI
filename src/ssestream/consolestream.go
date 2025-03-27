package ssestream

import (
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

		// Ensure the response writer supports flushing, not turning the server into a Potato until the client disconnects again (this took me a while to figure out)
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Create a new buffered channel for this client
		messageChan := make(chan string, 2000)

		// Register this client
		consoleClientsMu.Lock()
		consoleClients[messageChan] = true
		consoleClientsMu.Unlock()

		// Send initial connection event
		_, err := fmt.Fprintf(w, "data: Console log stream connected\n\n")
		if err != nil {
			consoleClientsMu.Lock()
			delete(consoleClients, messageChan)
			consoleClientsMu.Unlock()
			close(messageChan)
			log.Printf("%s🖥️ [Console/LogStream] ⚠️ Failed to send initial log message: %v%s", colorRed, err, colorReset)
			return
		}
		flusher.Flush()

		// Handle client disconnection
		notify := r.Context().Done()

		// Start a goroutine to stream messages
		go func(ch chan string) {
			defer func() {
				// Cleanup when the goroutine exits
				consoleClientsMu.Lock()
				delete(consoleClients, ch)
				consoleClientsMu.Unlock()
				close(ch)
				log.Printf("%s🖥️ [Console/LogStream] 👋 Client disconnected, channel closed%s", colorYellow, colorReset)
			}()

			for {
				select {
				case msg := <-ch:
					// Send the message to the client
					_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
					if err != nil {
						log.Printf("%s🖥️ [Console/LogStream] ❌ Failed to send log to client: %v%s", colorRed, err, colorReset)
						return
					}
					flusher.Flush()
				case <-notify:
					// Client disconnected
					return
				}
			}
		}(messageChan)

		// Keep the handler alive without blocking
		<-notify // Wait for the client to disconnect, but don’t block other handlers
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
			//if config.IsDebugMode {
			//	fmt.Println("Broadcasting to console clients:", message)
			//}
		case <-time.After(time.Second):
			// Timeout - client might be slow; log and keep the client
			log.Printf("%s🖥️ [Console/LogStream] ⏳ Timeout sending message after 1s: %q%s", colorMagenta, message, colorReset)
		}
	}
}
