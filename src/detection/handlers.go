// handlers.go
package detection

import (
	"StationeersServerUI/src/discord"
	"StationeersServerUI/src/ssestream"
	"fmt"
	"strings"
	"time"
)

/*
Event Handler Subsystem
- Defines default handling logic for detected events
- Formats and routes event notifications to:
  - Terminal output with ANSI coloring
  - SSE stream for web UI
*/

// DefaultHandlers returns a map of event types to default handlers
func DefaultHandlers() map[EventType]Handler {
	return map[EventType]Handler{

		EventCustomDetection: func(event Event) {
			message := fmt.Sprintf("🎮 [Custom Detection] %s", event.Message)
			fmt.Printf("%s%s%s%s\n", colorMagenta, message, colorReset, colorReset)
			ssestream.BroadcastDetectionEvent(message)
			discord.SendMessageToStatusChannel(message)
		},

		EventServerReady: func(event Event) {
			message := "🎮 [Gameserver] 🔔 Server is ready to connect!"
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorGreen, colorReset)
			ssestream.BroadcastDetectionEvent(message)
			discord.SendMessageToStatusChannel(message)
		},
		EventServerStarting: func(event Event) {
			message := "🎮 [Gameserver] 🕑 Server is starting up..."
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorYellow, colorReset)
			ssestream.BroadcastDetectionEvent(message)
			discord.SendMessageToStatusChannel(message)
		},
		EventServerError: func(event Event) {
			message := "🎮 [Gameserver] ⚠️ Server error detected"
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorRed, colorReset)
			ssestream.BroadcastDetectionEvent(message)
			discord.SendMessageToStatusChannel(message)
		},
		EventSettingsChanged: func(event Event) {
			message := fmt.Sprintf("🎮 [Gameserver] ⚙️ %s", event.Message)
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorYellow, colorReset)
			ssestream.BroadcastDetectionEvent(message)
			discord.SendMessageToStatusChannel(message)
		},
		EventServerHosted: func(event Event) {
			message := fmt.Sprintf("🎮 [Gameserver] 🌐 %s", event.Message)
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorGreen, colorReset)
			ssestream.BroadcastDetectionEvent(message)
			discord.SendMessageToStatusChannel(message)
		},
		EventNewGameStarted: func(event Event) {
			message := fmt.Sprintf("🎮 [Gameserver] 🎲 %s", event.Message)
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorGreen, colorReset)
			ssestream.BroadcastDetectionEvent(message)
			discord.SendMessageToStatusChannel(message)
		},
		EventServerRunning: func(event Event) {
			message := "🎮 [Gameserver] ✅ Server process has started!"
			fmt.Printf("%s%s%s%s\n", colorCyan, message, colorGreen, colorReset)
			ssestream.BroadcastDetectionEvent(message)
			discord.SendMessageToStatusChannel(message)
		},
		EventPlayerConnecting: func(event Event) {
			if event.PlayerInfo != nil {
				message := fmt.Sprintf("🎮 [Gameserver] 🔄 Player %s (SteamID: %s) is connecting...",
					event.PlayerInfo.Username, event.PlayerInfo.SteamID)
				fmt.Printf("%s%s%s%s%s%s%s\n",
					colorCyan, colorBlue, message, colorMagenta,
					colorBlue, event.PlayerInfo.SteamID, colorReset)
				ssestream.BroadcastDetectionEvent(message)
				discord.SendMessageToStatusChannel(message)
			}
		},
		EventPlayerReady: func(event Event) {
			if event.PlayerInfo != nil {
				message := fmt.Sprintf("🎮 [Gameserver] ✅ Player %s (SteamID: %s) is ready!",
					event.PlayerInfo.Username, event.PlayerInfo.SteamID)
				fmt.Printf("%s%s%s%s%s%s%s\n",
					colorCyan, colorGreen, message, colorMagenta,
					colorGreen, event.PlayerInfo.SteamID, colorReset)
				ssestream.BroadcastDetectionEvent(message)
				discord.SendMessageToStatusChannel(message)
			}
		},
		EventPlayerDisconnect: func(event Event) {
			if event.PlayerInfo != nil {
				message := fmt.Sprintf("🎮 [Gameserver] 👋 Player %s disconnected",
					event.PlayerInfo.Username)
				fmt.Printf("%s%s%s%s%s%s\n",
					colorCyan, colorYellow, message, colorMagenta,
					colorYellow, colorReset)
				ssestream.BroadcastDetectionEvent(message)
				discord.SendMessageToStatusChannel(message)
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
				ssestream.BroadcastDetectionEvent(message)
				discord.SendMessageToSavesChannel(message)
			}
		},
		EventException: func(event Event) {
			// Initial alert message
			alertMessage := "🎮 [Gameserver] 🚨 Exception detected!"
			fmt.Printf("%s%s%s%s\n", colorCyan, alertMessage, colorRed, colorReset)
			ssestream.BroadcastDetectionEvent(alertMessage)
			discord.SendUntrackedMessageToErrorChannel(alertMessage)

			if event.ExceptionInfo != nil && len(event.ExceptionInfo.StackTrace) > 0 {
				// Format stack trace as a single-line string for SSE compatibility
				stackTrace := strings.ReplaceAll(event.ExceptionInfo.StackTrace, "\n", " | ")
				detailedMessage := fmt.Sprintf("Exception Details: Stack Trace: %s", stackTrace)

				// Console output with original formatting
				fmt.Printf("%sException Details:\nStack Trace:\n%s%s%s%s\n",
					colorYellow, event.ExceptionInfo.StackTrace, colorReset, colorRed, colorReset)

				// Broadcast UI-friendly version
				ssestream.BroadcastDetectionEvent(detailedMessage)
				discord.SendUntrackedMessageToErrorChannel(detailedMessage)
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
