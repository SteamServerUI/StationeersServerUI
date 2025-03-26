// sse/console.go
package ssestream

import (
	"fmt"
	"net/http"
	"sync"
)

type EventManager struct {
	clients   []chan string
	clientsMu sync.Mutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		clients: make([]chan string, 0),
	}
}

func (em *EventManager) Broadcast(message string) {
	em.clientsMu.Lock()
	defer em.clientsMu.Unlock()

	for _, clientChan := range em.clients {
		select {
		case clientChan <- message:
		default: // Non-blocking send to avoid hanging on slow clients
		}
	}
}

func (em *EventManager) HandleEvents(w http.ResponseWriter, r *http.Request) {
	// Create a new channel for this client
	clientChan := make(chan string, 100) // Buffered channel to prevent blocking

	// Register the client
	em.clientsMu.Lock()
	em.clients = append(em.clients, clientChan)
	em.clientsMu.Unlock()

	// Ensure the channel is removed when the client disconnects
	defer func() {
		em.clientsMu.Lock()
		for i, ch := range em.clients {
			if ch == clientChan {
				em.clients = append(em.clients[:i], em.clients[i+1:]...)
				break
			}
		}
		em.clientsMu.Unlock()
		close(clientChan)
	}()

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Write data to the client as it comes in
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case msg, ok := <-clientChan:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
