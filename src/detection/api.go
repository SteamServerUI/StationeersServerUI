package detection

import (
	"fmt"
	"net/http"
	"sync"
)

// This file contains the public API for the detection package

var (
	detectionClients   = make(map[chan string]bool)
	detectionClientsMu sync.Mutex
)

// Start initializes the detector and returns it
func Start() *Detector {
	return NewDetector()
}

// AddHandler is a convenient method to register a handler for an event type
func AddHandler(detector *Detector, eventType EventType, handler Handler) {
	detector.RegisterHandler(eventType, handler)
}

// ProcessLog is a convenient method to process a log message
func ProcessLog(detector *Detector, logMessage string) {
	detector.ProcessLogMessage(logMessage)
}

// GetPlayers returns the currently connected players
func GetPlayers(detector *Detector) map[string]string {
	return detector.GetConnectedPlayers()
}

func StartDetectionEventStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Create a new channel for this client
		messageChan := make(chan string)

		// Register this client
		detectionClientsMu.Lock()
		detectionClients[messageChan] = true
		detectionClientsMu.Unlock()

		// Remove client when connection is closed
		notify := r.Context().Done()
		go func() {
			<-notify
			detectionClientsMu.Lock()
			delete(detectionClients, messageChan)
			detectionClientsMu.Unlock()
			close(messageChan)
		}()

		// Stream events to client
		for msg := range messageChan {
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush()
		}
	}
}

// BroadcastDetectionEvent sends an event to all connected clients
func BroadcastDetectionEvent(message string) {
	detectionClientsMu.Lock()
	defer detectionClientsMu.Unlock()

	for client := range detectionClients {
		select {
		case client <- message:
		default:
			// If the client's channel is full, move on
		}
	}
}
