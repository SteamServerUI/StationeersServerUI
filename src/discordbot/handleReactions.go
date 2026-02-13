package discordbot

import (
	"github.com/bwmarrin/discordgo"
)

// listenToDiscordReactions triggers when any reaction is added to any message. IF the reaction was added to a controled message, process it.
func listenToDiscordReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Ignore bot's own reactions
	if r.UserID == s.State.User.ID {
		return
	}

	// Check if the reaction was added to the control message for server control
	if r.MessageID == ControlMessageID {
		handleControlReactions(s, r)
		return
	}
	// Optionally, we could add more message-specific handlers here for other features
}
