// types.go
package detection

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
)

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
