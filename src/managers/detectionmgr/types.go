// types.go
package detectionmgr

import (
	"regexp"
	"sync"
)

// EventType defines the type of event detected
type EventType string

const (
	EventServerReady       EventType = "SERVER_READY"
	EventServerStarting    EventType = "SERVER_STARTING"
	EventServerError       EventType = "SERVER_ERROR"
	EventPlayerConnecting  EventType = "PLAYER_CONNECTING"
	EventPlayerReady       EventType = "PLAYER_READY"
	EventPlayerDisconnect  EventType = "PLAYER_DISCONNECT"
	EventWorldSaved        EventType = "WORLD_SAVED"
	EventException         EventType = "EXCEPTION"
	EventSettingsChanged   EventType = "SETTINGS_CHANGED"
	EventServerHosted      EventType = "SERVER_HOSTED"
	EventNewGameStarted    EventType = "NEW_GAME_STARTED"
	EventVersionExtracted  EventType = "VERSION_EXTRACTED"
	EventServerRunning     EventType = "SERVER_RUNNING"
	EventGameManagerReady  EventType = "GAME_MANAGER_READY"
	EventSessionStarting   EventType = "SESSION_STARTING"
	EventSessionRegistered EventType = "SESSION_REGISTERED"
	EventCustomDetection   EventType = "CUSTOM_DETECTION"
)

type Detector struct {
	handlers         map[EventType][]Handler
	connectedPlayers map[string]string // SteamID -> Username
	customPatterns   []CustomPattern
}

type CustomPattern struct {
	Pattern     *regexp.Regexp
	EventType   EventType
	MessageTmpl string
	IsRegex     bool
	Keyword     string
}

// Event represents a detected event from server logs
type Event struct {
	Type          EventType
	Message       string
	RawLog        string
	Timestamp     string
	PlayerInfo    *PlayerInfo
	ExceptionInfo *ExceptionInfo
}

// PlayerInfo contains information about a player
type PlayerInfo struct {
	Username string
	SteamID  string
}

// ExceptionInfo contains information about a server exception
type ExceptionInfo struct {
	StackTrace string
}

// Handler is a function that handles detected events
type Handler func(event Event)

var (
	serverStateHandler   func(EventType)
	serverStateHandlerMu sync.RWMutex
)

func SetServerStateHandler(handler func(EventType)) {
	serverStateHandlerMu.Lock()
	defer serverStateHandlerMu.Unlock()
	serverStateHandler = handler
}

func handleServerStateEvent(eventType EventType) {
	serverStateHandlerMu.RLock()
	handler := serverStateHandler
	serverStateHandlerMu.RUnlock()
	if handler != nil {
		handler(eventType)
	}
}
