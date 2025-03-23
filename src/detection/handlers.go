package detection

import (
	"fmt"
	"time"
)

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

// DefaultHandlers returns a map of event types to default handlers
func DefaultHandlers() map[EventType]Handler {
	return map[EventType]Handler{
		EventServerReady: func(event Event) {
			fmt.Printf("%s🎮 [Gameserver] %s🔔 Server is ready to connect!%s\n",
				colorCyan, colorGreen, colorReset)
		},
		EventServerStarting: func(event Event) {
			fmt.Printf("%s🎮 [Gameserver] %s🕑 Server is starting up...%s\n",
				colorCyan, colorYellow, colorReset)
		},
		EventServerError: func(event Event) {
			fmt.Printf("%s🎮 [Gameserver] %s⚠️ Server error detected%s\n",
				colorCyan, colorRed, colorReset)
		},
		EventPlayerConnecting: func(event Event) {
			if event.PlayerInfo != nil {
				fmt.Printf("%s🎮 [Gameserver] %s🔄 Player %s%s%s (SteamID: %s) is connecting...%s\n",
					colorCyan, colorBlue, colorMagenta, event.PlayerInfo.Username,
					colorBlue, event.PlayerInfo.SteamID, colorReset)
			}
		},
		EventPlayerReady: func(event Event) {
			if event.PlayerInfo != nil {
				fmt.Printf("%s🎮 [Gameserver] %s✅ Player %s%s%s (SteamID: %s) is ready!%s\n",
					colorCyan, colorGreen, colorMagenta, event.PlayerInfo.Username,
					colorGreen, event.PlayerInfo.SteamID, colorReset)
			}
		},
		EventPlayerDisconnect: func(event Event) {
			if event.PlayerInfo != nil {
				fmt.Printf("%s🎮 [Gameserver] %s👋 Player %s%s%s disconnected%s\n",
					colorCyan, colorYellow, colorMagenta, event.PlayerInfo.Username,
					colorYellow, colorReset)
			}
		},
		EventWorldSaved: func(event Event) {
			if event.BackupInfo != nil {
				fmt.Printf("%s🎮 [Gameserver] %s💾 World Saved: %sBackupIndex: %s%s UTC Time: %s%s\n",
					colorCyan, colorGreen, colorYellow, event.BackupInfo.BackupIndex,
					colorGreen, time.Now().UTC().Format(time.RFC3339), colorReset)
			}
		},
		EventException: func(event Event) {
			fmt.Printf("%s🎮 [Gameserver] %s🚨 Exception detected!%s\n",
				colorCyan, colorRed, colorReset)
			if event.ExceptionInfo != nil && len(event.ExceptionInfo.StackTrace) > 0 {
				fmt.Printf("%sStack trace:%s\n%s%s%s\n",
					colorYellow, colorReset, colorRed,
					event.ExceptionInfo.StackTrace, colorReset)
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
