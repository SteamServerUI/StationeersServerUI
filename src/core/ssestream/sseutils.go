// sseutils.go
package ssestream

import (
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

// Global managers for SSE streams
var (
	ConsoleStreamManager    = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
	EventStreamManager      = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
	DebugLogStreamManager   = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
	InfoLogStreamManager    = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
	WarnLogStreamManager    = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
	ErrorLogStreamManager   = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
	BackendLogStreamManager = NewSSEManager(config.MaxSSEConnections, config.SSEMessageBufferSize)
)

// BroadcastConsoleOutput sends log to all connected console log clients
func BroadcastConsoleOutput(message string) {
	ConsoleStreamManager.Broadcast(message)
}

// BroadcastDetectionEvent sends an event to all connected clients
func BroadcastDetectionEvent(message string) {
	EventStreamManager.Broadcast(message)
}

// BroadcastDebugLog sends an event to all connected clients
func BroadcastDebugLog(message string) {
	DebugLogStreamManager.Broadcast(message)
}

// BroadcastInfoLog sends an event to all connected clients
func BroadcastInfoLog(message string) {
	InfoLogStreamManager.Broadcast(message)
}

// BroadcastWarnLog sends an event to all connected clients
func BroadcastWarnLog(message string) {
	WarnLogStreamManager.Broadcast(message)
}

// BroadcastErrorLog sends an event to all connected clients
func BroadcastErrorLog(message string) {
	ErrorLogStreamManager.Broadcast(message)
}

// BroadcastInternalLog sends an event to all connected clients
func BroadcastBackendLog(message string) {
	BackendLogStreamManager.Broadcast(message)
}
