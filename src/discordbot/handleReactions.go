package discordbot

import (
	"fmt"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"

	"github.com/bwmarrin/discordgo"
)

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

func handleControlReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// handleControlReactions - Handles reactions for server control actions
	var actionMessage string

	switch r.Emoji.Name {
	case "▶️": // Start action
		gamemgr.InternalStartServer()
		actionMessage = "🕛Server is Starting..."
	case "⏹️": // Stop action
		gamemgr.InternalStopServer()
		actionMessage = "🛑Server is Stopping..."
	case "♻️": // Restart action
		actionMessage = "♻️Server is restarting..."
		go func() {
			// Perform stop operation
			gamemgr.InternalStopServer()

			// Non-blocking delay using channel and goroutine
			delayChan := make(chan bool)
			go func() {
				time.Sleep(5 * time.Second)
				delayChan <- true
			}()

			// Wait for delay to complete
			<-delayChan

			// Start server after delay
			gamemgr.InternalStartServer()
		}()

	default:
		logger.Discord.Debug("Unknown reaction: " + r.Emoji.Name)
		return
	}

	// Get the user who triggered the action
	user, err := s.User(r.UserID)
	if err != nil {
		logger.Discord.Error("Error fetching user details: " + err.Error())
		return
	}
	username := user.Username

	// Send the action message to the control channel
	SendMessageToStatusChannel(fmt.Sprintf("%s triggered by %s.", actionMessage, username))

	// Remove the reaction after processing
	err = s.MessageReactionRemove(config.ControlPanelChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
	if err != nil {
		logger.Discord.Error("Error removing reaction: " + err.Error())
	}
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
	err = s.MessageReactionRemove(config.ErrorChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
	if err != nil {
		logger.Discord.Error("Error removing reaction: " + err.Error())
	}
}
