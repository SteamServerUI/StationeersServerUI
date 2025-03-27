package ssestream

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

// StartDetectionEventStream creates and returns an HTTP handler for SSE event streaming
func StartDetectionEventStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Ensure the response writer supports flushing
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Create a new buffered channel for this client
		messageChan := make(chan string, 2000)

		// Register this client
		eventClientsMu.Lock()
		eventClients[messageChan] = true
		eventClientsMu.Unlock()

		// Send initial connection event
		_, err := fmt.Fprintf(w, "data: UI event stream connected\n\n")
		if err != nil {
			eventClientsMu.Lock()
			delete(eventClients, messageChan)
			eventClientsMu.Unlock()
			close(messageChan)
			log.Printf("%s🖥️ [UI/EventStream] ⚠️ Failed to send initial event: %v%s", colorRed, err, colorReset)
			return
		}
		flusher.Flush()

		// Handle client disconnection
		notify := r.Context().Done()

		// Start a goroutine to stream messages
		go func(ch chan string) {
			defer func() {
				// Cleanup when the goroutine exits
				eventClientsMu.Lock()
				delete(eventClients, ch)
				eventClientsMu.Unlock()
				close(ch)
				log.Printf("%s🖥️ [UI/EventStream] 👋 Client disconnected, channel closed%s", colorYellow, colorReset)
			}()

			for {
				select {
				case msg := <-ch:
					// Send the message to the client
					_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
					if err != nil {
						log.Printf("%s🖥️ [UI/EventStream] ❌ Failed to send event to client: %v%s", colorRed, err, colorReset)
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
		<-notify // Wait for the client to disconnect, but don't block other handlers
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
