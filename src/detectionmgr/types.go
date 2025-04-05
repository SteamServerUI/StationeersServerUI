// types.go
package detectionmgr

import "regexp"

const (
	// ANSI color codes for styling terminal output
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

// EventType defines the type of event detected
type EventType string

const (
	EventServerReady      EventType = "SERVER_READY"
	EventServerStarting   EventType = "SERVER_STARTING"
	EventServerError      EventType = "SERVER_ERROR"
	EventPlayerConnecting EventType = "PLAYER_CONNECTING"
	EventPlayerReady      EventType = "PLAYER_READY"
	EventPlayerDisconnect EventType = "PLAYER_DISCONNECT"
	EventWorldSaved       EventType = "WORLD_SAVED"
	EventException        EventType = "EXCEPTION"
	EventSettingsChanged  EventType = "SETTINGS_CHANGED"
	EventServerHosted     EventType = "SERVER_HOSTED"
	EventNewGameStarted   EventType = "NEW_GAME_STARTED"
	EventServerRunning    EventType = "SERVER_RUNNING"
	EventCustomDetection  EventType = "CUSTOM_DETECTION"
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
	BackupInfo    *BackupInfo
	ExceptionInfo *ExceptionInfo
}

// PlayerInfo contains information about a player
type PlayerInfo struct {
	Username string
	SteamID  string
}

// BackupInfo contains information about a world save/backup
type BackupInfo struct {
	BackupIndex string
}

// ExceptionInfo contains information about a server exception
type ExceptionInfo struct {
	StackTrace string
}

// Handler is a function that handles detected events
type Handler func(event Event)
