package discord

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/core"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// listenToDiscordMessages triggers when any message is sent in the bots scope, regardless of the channel. IF the message is sent in the control channel, and the author is not the bot, process it.
func listenToDiscordMessages(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || m.ChannelID != config.ControlChannelID {
		return
	}

	content := strings.TrimSpace(m.Content)

	switch {
	case strings.HasPrefix(content, "!start"):
		core.InternalStartServer()
		s.ChannelMessageSend(m.ChannelID, "🕛Server is starting...")
		SendMessageToStatusChannel("🕛Start command received from Server Controller, Server is Starting...")

	case strings.HasPrefix(content, "!stop"):
		core.InternalStopServer()
		s.ChannelMessageSend(m.ChannelID, "🕛Server is stopping...")
		SendMessageToStatusChannel("🕛Stop command received from Server Controller, flatlining Server in 5 Seconds...")

	case strings.HasPrefix(content, "!restore"):
		SendMessageToStatusChannel("⚠️Restore command received, flatlining and restoring Server in 5 Seconds. Server will come back online in about 60 Seconds.")
		handleRestoreCommand(content)

	case strings.HasPrefix(content, "!list"):
		handleListCommand(content)

	case strings.HasPrefix(content, "!update"):
		handleUpdateCommand(s, m.ChannelID)

	case strings.HasPrefix(content, "!help"):
		handleHelpCommand(s, m.ChannelID)

	case strings.HasPrefix(content, "!ban"):
		handleBanCommand(s, m.ChannelID, content)

	case strings.HasPrefix(content, "!unban"):
		handleUnbanCommand(s, m.ChannelID, content)

	case strings.HasPrefix(content, "!validate"):
		handleValidateCommand(s, m.ChannelID)
	default:
		// Optionally handle unrecognized commands or ignore them
	}
}

// listenToDiscordReactions triggers when any reaction is added to any message. IF the reaction was added to a controled message, process it.
func listenToDiscordReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Ignore bot's own reactions
	if r.UserID == s.State.User.ID {
		return
	}

	// Check if the reaction was added to the control message for server control
	if r.MessageID == config.ControlMessageID {
		handleControlReactions(s, r)
		return
	}

	// Check if the reaction was added to the last exception message (not used in v4.3)
	if r.MessageID == config.ExceptionMessageID {
		handleExceptionReactions(s, r)
		return
	}
	// Optionally, we could add more message-specific handlers here for other features
}
