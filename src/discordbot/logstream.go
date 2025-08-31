package discordbot

import (
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// PassLogMessageToDiscordLogBuffer is called from the detection module to add a log message to the buffer.
func PassLogStreamToDiscordLogBuffer(logMessage string) {
	LogMessageBuffer += logMessage + "\n"
	if len(LogMessageBuffer) >= config.GetDiscordCharBufferSize() && config.GetIsDiscordEnabled() {
		flushLogBufferToDiscord()
	}
}

// FlushLogBufferToDiscord flushes the log buffer to Discord periodically with a configurable "DiscordCharBufferSize" character limit per message.
func flushLogBufferToDiscord() {
	if len(LogMessageBuffer) == 0 {
		return // No messages to send
	}
	if !config.GetIsDiscordEnabled() || config.DiscordSession == nil {
		return
	}

	discordMaxMessageLength := config.DiscordCharBufferSize

	message := LogMessageBuffer

	for len(message) > 0 {
		// Determine how much of the message we can send
		chunkSize := discordMaxMessageLength
		if len(message) < discordMaxMessageLength {
			chunkSize = len(message)
		}

		// Send the chunk to Discord
		_, err := config.DiscordSession.ChannelMessageSend(config.GetLogChannelID(), message[:chunkSize])
		if err != nil {
			logger.Discord.Error("Error sending log to Discord: " + err.Error())
			break
		}

		// Move to the next chunk
		message = message[chunkSize:]
	}

	// Clear the buffer after sending
	LogMessageBuffer = ""
}
