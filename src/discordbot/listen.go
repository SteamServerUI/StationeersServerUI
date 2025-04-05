package discordbot

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/gamemgr"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

/*
Ban and Unban commands are currently only supported by ![command] commands.
This file contains the logic for ![command] discord commands using message content instead of proper slash commands.
This functionality is deprecated and will be removed in a future version as it has been replaced with slash commands entirely.
*/

// listenToDiscordMessages triggers when any message is sent in the bots scope, regardless of the channel. IF the message is sent in the control channel, and the author is not the bot, process it.
func listenToDiscordMessages(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || m.ChannelID != config.ControlChannelID {
		if config.IsDebugMode {
			fmt.Println("Ignoring message from", m.Author.Username)
			fmt.Println("Ignored message:", m.Content)
			fmt.Println("Message channel:", m.ChannelID)
		}
		return
	}

	// log the message if debug is enabled
	if config.IsDebugMode {
		fmt.Println("Received message:", m.Content)
		fmt.Println("Message author:", m.Author.Username)
		fmt.Println("Message channel:", m.ChannelID)
	}

	content := strings.TrimSpace(m.Content)

	switch {
	case strings.HasPrefix(content, "!start"):
		gamemgr.InternalStartServer()
		SendMessageToControlChannel("🕛Server is starting...")

	case strings.HasPrefix(content, "!stop"):
		gamemgr.InternalStopServer()
		SendMessageToControlChannel("🕛Server is stopping...")

	case strings.HasPrefix(content, "!restore"):
		handleRestoreCommand()

	case strings.HasPrefix(content, "!list"):
		handleListCommand()

	case strings.HasPrefix(content, "!update"):
		handleUpdateCommand()

	case strings.HasPrefix(content, "!help"):
		handleHelpCommand(s, m.ChannelID)

	case strings.HasPrefix(content, "!ban"):
		handleBanCommand(s, m.ChannelID, content)

	case strings.HasPrefix(content, "!unban"):
		handleUnbanCommand(s, m.ChannelID, content)

	case strings.HasPrefix(content, "!validate"):
		handleValidateCommand()
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

	// Check if the reaction was added to the last sent exception message for attaching restart buttons. Not used in v4.3 as nothing is sending tracked Exception messages to Discord anymore.
	//  Instead, we now only yoink the exception message to Discord without tracking it, thus there is no onfig.ExceptionMessageID set anymore. Removed as this was a rather unused feature.
	if r.MessageID == config.ExceptionMessageID {
		handleExceptionReactions(s, r)
		return
	}
	// Optionally, we could add more message-specific handlers here for other features
}
