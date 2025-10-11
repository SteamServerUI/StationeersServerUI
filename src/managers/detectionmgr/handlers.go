// handlers.go
package detectionmgr

import (
	"fmt"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/core/ssestream"
	"github.com/SteamServerUI/SteamServerUI/v7/src/discord/discordbot"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
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
			logger.Detection.Info(message)
			ssestream.BroadcastDetectionEvent(message)
			discordbot.SendMessageToStatusChannel(message)
		},

		EventServerReady: func(event Event) {
			message := "🎮 [Gameserver] 🔔 Server is ready to connect!"
			logger.Detection.Info(message)
			ssestream.BroadcastDetectionEvent(message)
			discordbot.SendMessageToStatusChannel(message)
		},
		EventServerStarting: func(event Event) {
			message := "🎮 [Gameserver] 🕑 Server is starting up..."
			logger.Detection.Info(message)
			ssestream.BroadcastDetectionEvent(message)
			discordbot.SendMessageToStatusChannel(message)
		},
		EventServerError: func(event Event) {
			message := "🎮 [Gameserver] ⚠️ Server error detected"
			logger.Detection.Info(message)
			ssestream.BroadcastDetectionEvent(message)
			discordbot.SendMessageToStatusChannel(message)
		},
		EventSettingsChanged: func(event Event) {
			message := fmt.Sprintf("🎮 [Gameserver] ⚙️ %s", event.Message)
			logger.Detection.Info(message)
			ssestream.BroadcastDetectionEvent(message)
			discordbot.SendMessageToStatusChannel(message)
		},
		EventServerHosted: func(event Event) {
			message := fmt.Sprintf("🎮 [Gameserver] 🌐 %s", event.Message)
			logger.Detection.Info(message)
			ssestream.BroadcastDetectionEvent(message)
			discordbot.SendMessageToStatusChannel(message)
		},
		EventNewGameStarted: func(event Event) {
			message := fmt.Sprintf("🎮 [Gameserver] 🎲 %s", event.Message)
			logger.Detection.Info(message)
			ssestream.BroadcastDetectionEvent(message)
			discordbot.SendMessageToStatusChannel(message)
		},
		EventVersionExtracted: func(event Event) {
			message := fmt.Sprintf("🎮 [Gameserver] 📦 Version %s detected", event.Message)
			logger.Detection.Info(message)
			ssestream.BroadcastDetectionEvent(message)
			discordbot.SendMessageToStatusChannel(message)
		},
		EventServerRunning: func(event Event) {
			message := "🎮 [Gameserver] ✅ Server process has started!"
			logger.Detection.Info(message)
			ssestream.BroadcastDetectionEvent(message)
			discordbot.SendMessageToStatusChannel(message)
		},
		EventPlayerConnecting: func(event Event) {
			if event.PlayerInfo != nil {
				message := fmt.Sprintf("🎮 [Gameserver] 🔄 Player %s (SteamID: %s) is connecting...",
					event.PlayerInfo.Username, event.PlayerInfo.SteamID)
				logger.Detection.Info(message)
				ssestream.BroadcastDetectionEvent(message)
				discordbot.SendMessageToStatusChannel(message)
			}
		},
		EventPlayerReady: func(event Event) {
			if event.PlayerInfo != nil {
				message := fmt.Sprintf("🎮 [Gameserver] ✅ Player %s (SteamID: %s) is ready!",
					event.PlayerInfo.Username, event.PlayerInfo.SteamID)
				logger.Detection.Info(message)
				ssestream.BroadcastDetectionEvent(message)
				discordbot.SendMessageToStatusChannel(message)
			}
		},
		EventPlayerDisconnect: func(event Event) {
			if event.PlayerInfo != nil {
				message := fmt.Sprintf("🎮 [Gameserver] 👋 Player %s disconnected",
					event.PlayerInfo.Username)
				logger.Detection.Info(message)
				ssestream.BroadcastDetectionEvent(message)
				discordbot.SendMessageToStatusChannel(message)
			}
		},
		EventWorldSaved: func(event Event) {
			if event.BackupInfo != nil {
				timeStr := time.Now().UTC().Format(time.RFC3339)
				message := fmt.Sprintf("🎮 [Gameserver] 💾 World Saved: BackupIndex: %s UTC Time: %s",
					event.BackupInfo.BackupIndex, timeStr)
				logger.Detection.Info(message)
				ssestream.BroadcastDetectionEvent(message)
				discordbot.SendMessageToSavesChannel(message)
			}
		},
		EventException: func(event Event) {
			// Initial alert message
			alertMessage := "🎮 [Gameserver] 🚨 Exception detected!"
			logger.Detection.Info(alertMessage)
			ssestream.BroadcastDetectionEvent(alertMessage)
			discordbot.SendUntrackedMessageToErrorChannel(alertMessage)

			if event.ExceptionInfo != nil && len(event.ExceptionInfo.StackTrace) > 0 {
				// Format stack trace as a single-line string for SSE compatibility
				stackTrace := strings.ReplaceAll(event.ExceptionInfo.StackTrace, "\n", " | ")
				message := fmt.Sprintf("Exception Details: Stack Trace: %s", stackTrace)

				logger.Detection.Info(message)
				ssestream.BroadcastDetectionEvent(message)
				discordbot.SendUntrackedMessageToErrorChannel(message)
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
