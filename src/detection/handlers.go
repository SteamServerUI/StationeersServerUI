// Update src/detection/handlers.go
package detection

import (
	"fmt"
	"strings"
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
			message := "🎮 [Gameserver] 🔔 Server is ready to connect!"
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorGreen, colorReset)
			BroadcastDetectionEvent(message)
		},
		EventServerStarting: func(event Event) {
			message := "🎮 [Gameserver] 🕑 Server is starting up..."
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorYellow, colorReset)
			BroadcastDetectionEvent(message)
		},
		EventServerError: func(event Event) {
			message := "🎮 [Gameserver] ⚠️ Server error detected"
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorRed, colorReset)
			BroadcastDetectionEvent(message)
		},
		EventPlayerConnecting: func(event Event) {
			if event.PlayerInfo != nil {
				message := fmt.Sprintf("🎮 [Gameserver] 🔄 Player %s (SteamID: %s) is connecting...",
					event.PlayerInfo.Username, event.PlayerInfo.SteamID)
				fmt.Printf("%s%s%s%s%s%s%s\n",
					colorCyan, colorBlue, message, colorMagenta,
					colorBlue, event.PlayerInfo.SteamID, colorReset)
				BroadcastDetectionEvent(message)
			}
		},
		EventPlayerReady: func(event Event) {
			if event.PlayerInfo != nil {
				message := fmt.Sprintf("🎮 [Gameserver] ✅ Player %s (SteamID: %s) is ready!",
					event.PlayerInfo.Username, event.PlayerInfo.SteamID)
				fmt.Printf("%s%s%s%s%s%s%s\n",
					colorCyan, colorGreen, message, colorMagenta,
					colorGreen, event.PlayerInfo.SteamID, colorReset)
				BroadcastDetectionEvent(message)
			}
		},
		EventPlayerDisconnect: func(event Event) {
			if event.PlayerInfo != nil {
				message := fmt.Sprintf("🎮 [Gameserver] 👋 Player %s disconnected",
					event.PlayerInfo.Username)
				fmt.Printf("%s%s%s%s%s%s\n",
					colorCyan, colorYellow, message, colorMagenta,
					colorYellow, colorReset)
				BroadcastDetectionEvent(message)
			}
		},
		EventWorldSaved: func(event Event) {
			if event.BackupInfo != nil {
				timeStr := time.Now().UTC().Format(time.RFC3339)
				message := fmt.Sprintf("🎮 [Gameserver] 💾 World Saved: BackupIndex: %s UTC Time: %s",
					event.BackupInfo.BackupIndex, timeStr)
				fmt.Printf("%s%s%s%s%s%s%s\n",
					colorCyan, colorGreen, message, colorYellow,
					colorGreen, timeStr, colorReset)
				BroadcastDetectionEvent(message)
			}
		},
		EventException: func(event Event) {
			// Initial alert message
			alertMessage := "🎮 [Gameserver] 🚨 Exception detected!"
			fmt.Printf("%s%s%s%s\n", colorCyan, alertMessage, colorRed, colorReset)
			BroadcastDetectionEvent(alertMessage)

			if event.ExceptionInfo != nil && len(event.ExceptionInfo.StackTrace) > 0 {
				// Format stack trace as a single-line string for SSE compatibility
				stackTrace := strings.ReplaceAll(event.ExceptionInfo.StackTrace, "\n", " | ")
				detailedMessage := fmt.Sprintf("Exception Details: Stack Trace: %s", stackTrace)

				// Console output with original formatting
				fmt.Printf("%sException Details:\nStack Trace:\n%s%s%s%s\n",
					colorYellow, event.ExceptionInfo.StackTrace, colorReset, colorRed, colorReset)

				// Broadcast UI-friendly version
				BroadcastDetectionEvent(detailedMessage)
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
