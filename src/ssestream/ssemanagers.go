// ssemanagers.go
package ssestream

import (
	"StationeersServerUI/src/config"
	"net/http"
)

// Global managers for console and event streams
var (
	ConsoleStreamManager = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
	EventStreamManager   = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
)

// StartConsoleStream creates an HTTP handler for console log SSE streaming
func StartConsoleStream() http.HandlerFunc {
	return ConsoleStreamManager.CreateStreamHandler("Console")
}

// StartDetectionEventStream creates an HTTP handler for detection event SSE streaming
func StartDetectionEventStream() http.HandlerFunc {
	return EventStreamManager.CreateStreamHandler("Event")
}

// BroadcastConsoleOutput sends log to all connected console log clients
func BroadcastConsoleOutput(message string) {
	ConsoleStreamManager.Broadcast(message)
}

// BroadcastDetectionEvent sends an event to all connected clients
func BroadcastDetectionEvent(message string) {
	EventStreamManager.Broadcast(message)
}
