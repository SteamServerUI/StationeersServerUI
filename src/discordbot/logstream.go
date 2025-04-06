package discordbot

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/logger"
)

// PassLogMessageToDiscordLogBuffer is called from the detection module to add a log message to the buffer.
func PassLogStreamToDiscordLogBuffer(logMessage string) {
	config.LogMessageBuffer += logMessage + "\n"
	if len(config.LogMessageBuffer) >= config.DiscordCharBufferSize && config.IsDiscordEnabled {
		flushLogBufferToDiscord()
	}
}

// FlushLogBufferToDiscord flushes the log buffer to Discord periodically with a configurable "DiscordCharBufferSize" character limit per message.
func flushLogBufferToDiscord() {
	if len(config.LogMessageBuffer) == 0 {
		return // No messages to send
	}
	if !config.IsDiscordEnabled || config.DiscordSession == nil {
		return
	}

	discordMaxMessageLength := config.DiscordCharBufferSize

	message := config.LogMessageBuffer

	for len(message) > 0 {
		// Determine how much of the message we can send
		chunkSize := discordMaxMessageLength
		if len(message) < discordMaxMessageLength {
			chunkSize = len(message)
		}

		// Send the chunk to Discord
		_, err := config.DiscordSession.ChannelMessageSend(config.LogChannelID, message[:chunkSize])
		if err != nil {
			logger.Discord.Error("Error sending log to Discord: " + err.Error())
			break
		}

		// Move to the next chunk
		message = message[chunkSize:]
	}

	// Clear the buffer after sending
	config.LogMessageBuffer = ""
}
