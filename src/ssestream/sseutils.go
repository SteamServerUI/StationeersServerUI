// sseutils.go
package ssestream

import (
	"StationeersServerUI/src/config"
)

// Global managers for console and event streams
var (
	ConsoleStreamManager = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
	EventStreamManager   = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
)

// BroadcastConsoleOutput sends log to all connected console log clients
func BroadcastConsoleOutput(message string) {
	ConsoleStreamManager.Broadcast(message)
}

// BroadcastDetectionEvent sends an event to all connected clients
func BroadcastDetectionEvent(message string) {
	EventStreamManager.Broadcast(message)
}
