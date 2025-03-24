// sse.go
package ui

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	// eventClients tracks all connected SSE clients
	eventClients   = make(map[chan string]bool)
	eventClientsMu sync.Mutex
)

// ANSI color codes for styling terminal output
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

// StartDetectionEventStream creates and returns an HTTP handler for SSE event streaming
func StartDetectionEventStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Create a new buffered channel for this client (buffer of 2000 events)
		messageChan := make(chan string, 2000)

		// Register this client
		eventClientsMu.Lock()
		eventClients[messageChan] = true
		eventClientsMu.Unlock()

		// Send initial connection event
		_, err := fmt.Fprintf(w, "data: UI event stream connected\n\n")
		if err != nil {
			// Early disconnect; clean up immediately
			eventClientsMu.Lock()
			delete(eventClients, messageChan)
			eventClientsMu.Unlock()
			close(messageChan)
			log.Printf("%s🖥️ [UI/EventStream] ⚠️ Failed to send initial event: %v%s", colorRed, err, colorReset)
			return
		}
		w.(http.Flusher).Flush()

		// Remove client when connection is closed
		notify := r.Context().Done()
		go func(ch chan string) {
			<-notify
			eventClientsMu.Lock()
			delete(eventClients, ch)
			eventClientsMu.Unlock()
			close(ch)
			log.Printf("%s🖥️ [UI/EventStream] 👋 Client disconnected, channel closed%s", colorYellow, colorReset)
		}(messageChan)

		// Stream events to client
		for msg := range messageChan {
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				// Client likely disconnected; rely on context cancellation to clean up
				log.Printf("%s🖥️ [UI/EventStream] ❌ Failed to send event to client: %v%s", colorRed, err, colorReset)
				return
			}
			w.(http.Flusher).Flush()
		}
	}
}

// BroadcastDetectionEvent sends an event to all connected clients
func BroadcastDetectionEvent(message string) {
	eventClientsMu.Lock()
	defer eventClientsMu.Unlock()

	for client := range eventClients {
		select {
		case client <- message:
			// Message sent successfully
		case <-time.After(time.Second):
			// Timeout - client might be slow; log and keep the client
			log.Printf("%s🖥️ [UI/EventStream] ⏳ Timeout sending message after 1s: %q%s", colorMagenta, message, colorReset)
		}
	}
}
