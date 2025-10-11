package discordbot

import (
	"fmt"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"

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

	// Check if the reaction was added to the last sent exception message for attaching restart buttons. Not used in v4.3 as nothing is sending tracked Exception messages to Discord anymore.
	//  Instead, we now only yoink the exception message to Discord without tracking it, thus there is no onfig.ExceptionMessageID set anymore. Removed as this was a rather unused feature.
	if r.MessageID == config.ExceptionMessageID {
		handleExceptionReactions(s, r)
		return
	}
	// Optionally, we could add more message-specific handlers here for other features
}

// v4 FIXED, Unused in v4.3
func handleExceptionReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	var actionMessage string

	switch r.Emoji.Name {
	case "♻️": // Stop server action due to exception
		actionMessage = "🛑 Server is manually restarting due to critical exception."
		gamemgr.InternalStopServer()
		//sleep 5 sec
		time.Sleep(5 * time.Second)
		gamemgr.InternalStartServer()

	default:
		logger.Discord.Debug("Unknown reaction: " + r.Emoji.Name)
		return
	}

	// Get the user who triggered the action
	user, err := s.User(r.UserID)
	if err != nil {
		logger.Discord.Error("Error fetching user details:\n" + err.Error())
		return
	}
	username := user.Username

	// Send the action message to the error channel
	sendMessageToErrorChannel(fmt.Sprintf("%s triggered by %s.", actionMessage, username))

	// Remove the reaction after processing
	err = s.MessageReactionRemove(config.GetErrorChannelID(), r.MessageID, r.Emoji.APIName(), r.UserID)
	if err != nil {
		logger.Discord.Error("Error removing reaction: " + err.Error())
	}
}
