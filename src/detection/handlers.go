package detection

import (
	"fmt"
	"time"
)

// DefaultHandlers returns a map of event types to default handlers
func DefaultHandlers() map[EventType]Handler {
	return map[EventType]Handler{
		EventServerReady: func(event Event) {
			fmt.Println("🔔 Server is ready to connect!")
		},
		EventServerStarting: func(event Event) {
			fmt.Println("🕑 Server is starting up...")
		},
		EventServerError: func(event Event) {
			fmt.Println("⚠️ Server error detected")
		},
		EventPlayerConnecting: func(event Event) {
			if event.PlayerInfo != nil {
				fmt.Printf("🔄 Player %s (SteamID: %s) is connecting...\n",
					event.PlayerInfo.Username, event.PlayerInfo.SteamID)
			}
		},
		EventPlayerReady: func(event Event) {
			if event.PlayerInfo != nil {
				fmt.Printf("✅ Player %s (SteamID: %s) is ready!\n",
					event.PlayerInfo.Username, event.PlayerInfo.SteamID)
			}
		},
		EventPlayerDisconnect: func(event Event) {
			if event.PlayerInfo != nil {
				fmt.Printf("👋 Player %s disconnected\n", event.PlayerInfo.Username)
			}
		},
		EventWorldSaved: func(event Event) {
			if event.BackupInfo != nil {
				fmt.Printf("💾 World Saved: BackupIndex: %s UTC Time: %s\n",
					event.BackupInfo.BackupIndex, time.Now().UTC().Format(time.RFC3339))
			}
		},
		EventException: func(event Event) {
			fmt.Println("🚨 Exception detected!")
			if event.ExceptionInfo != nil && len(event.ExceptionInfo.StackTrace) > 0 {
				fmt.Println("Stack trace:")
				fmt.Println(event.ExceptionInfo.StackTrace)
			}
		},
	}
}

// RegisterDefaultHandlers registers all default handlers with a detector
func RegisterDefaultHandlers(detector *Detector) {
	for eventType, handler := range DefaultHandlers() {
		detector.RegisterHandler(eventType, handler)
	}
}
