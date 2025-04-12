// ssemanager.go
package ssestream

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// The SSE blocking issue is NOT related to the backend; the API handles 200 clients per channel fine.
// The issue lies in the JS frontend, where browsers cap HTTP/1.1 SSE streams at ~6 concurrent connections per domain.
// See RFC 2616 (HTTP/1.1), Section 8.1.4: https://datatracker.ietf.org/doc/html/rfc2616#section-8.1.4
// SSE spec (WHATWG HTML): https://html.spec.whatwg.org/multipage/server-sent-events.html#server-sent-events
// Workaround: Use HTTP/2 (RFC 7540, Section 5.1.2) for up to 100 streams: https://datatracker.ietf.org/doc/html/rfc7540#section-5.1.2
// I spent way too much time on this thinking it would be the backend blocking requests. It's not.

// Client represents a connected SSE client
type Client struct {
	messages chan string
	lastSeen time.Time
}

// SSEManager manages Server-Sent Event streams
type SSEManager struct {
	clients    map[*Client]bool
	clientsMu  sync.RWMutex
	maxClients int
	maxBuffer  int
}

// NewSSEManager creates a new SSE stream manager
func NewSSEManager(maxClients, maxBuffer int) *SSEManager {
	return &SSEManager{
		clients:    make(map[*Client]bool),
		maxClients: maxClients,
		maxBuffer:  maxBuffer,
	}
}

// CreateStreamHandler creates an HTTP handler for SSE streaming
func (m *SSEManager) CreateStreamHandler(streamType string) http.HandlerFunc {
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

		// Check maximum client limit
		m.clientsMu.Lock()
		if len(m.clients) >= m.maxClients {
			m.clientsMu.Unlock()
			http.Error(w, "Too many clients", http.StatusServiceUnavailable)
			return
		}

		// Create a new client
		client := &Client{
			messages: make(chan string, m.maxBuffer),
			lastSeen: time.Now(),
		}
		m.clients[client] = true
		m.clientsMu.Unlock()

		// Send initial connection event
		_, err := fmt.Fprintf(w, "data: %s Stream Connected\n\n", streamType)
		if err != nil {
			m.removeClient(client)
			logger.SSE.Error(" ⚠️ Failed to send initial message: " + err.Error())
			return
		}
		flusher.Flush()

		// Handle client disconnection
		notify := r.Context().Done()

		// Start streaming messages
		go m.streamMessages(w, flusher, client, streamType, notify)

		// Wait for client disconnection
		<-notify
	}
}

// streamMessages handles sending messages to a specific client
func (m *SSEManager) streamMessages(
	w http.ResponseWriter,
	flusher http.Flusher,
	client *Client,
	streamType string,
	notify <-chan struct{},
) {
	defer m.removeClient(client)

	for {
		select {
		case msg := <-client.messages:
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				logger.SSE.Error(" ❌ Failed to send message: " + err.Error())
				return
			}
			flusher.Flush()

		case <-notify:
			return
		}
	}
}

// Broadcast sends a message to all clients with a non-blocking approach
func (m *SSEManager) Broadcast(message string) {
	m.clientsMu.RLock()
	defer m.clientsMu.RUnlock()

	for client := range m.clients {
		select {
		case client.messages <- message:
			// Message sent successfully
		default:
			// Client channel is full, log and skip
			logger.SSE.Warn("⏳ Message dropped for slow client")
		}
	}
}

func (m *SSEManager) AddInternalSubscriber() chan string {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()

	client := &Client{
		messages: make(chan string, m.maxBuffer),
		lastSeen: time.Now(),
	}
	m.clients[client] = true
	return client.messages
}

// removeClient safely removes a client from the manager
func (m *SSEManager) removeClient(client *Client) {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()

	delete(m.clients, client)
	close(client.messages)
}
