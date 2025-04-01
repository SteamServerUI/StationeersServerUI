package discord

import (
	"fmt"
	"time"

	"StationeersServerUI/src/config"
	"StationeersServerUI/src/core"

	"github.com/bwmarrin/discordgo"
)

func handleControlReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// handleControlReactions - Handles reactions for server control actions
	var actionMessage string

	switch r.Emoji.Name {
	case "▶️": // Start action
		core.InternalStartServer()
		actionMessage = "🕛Server is Starting..."
	case "⏹️": // Stop action
		core.InternalStopServer()
		actionMessage = "🛑Server is Stopping..."
	case "♻️": // Restart action
		actionMessage = "♻️Server is restarting..."
		core.InternalStopServer()
		//sleep 5 sec
		time.Sleep(5 * time.Second)
		core.InternalStartServer()

	default:
		fmt.Println("Unknown reaction:", r.Emoji.Name)
		return
	}

	// Get the user who triggered the action
	user, err := s.User(r.UserID)
	if err != nil {
		fmt.Printf("Error fetching user details: %v\n", err)
		return
	}
	username := user.Username

	// Send the action message to the control channel
	SendMessageToStatusChannel(fmt.Sprintf("%s triggered by %s.", actionMessage, username))

	// Remove the reaction after processing
	err = s.MessageReactionRemove(config.ControlPanelChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
	if err != nil {
		fmt.Printf("Error removing reaction: %v\n", err)
	}
}

// v4 FIXED, Unused in v4.3
func handleExceptionReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	var actionMessage string

	switch r.Emoji.Name {
	case "♻️": // Stop server action due to exception
		actionMessage = "🛑 Server is manually restarting due to critical exception."
		core.InternalStopServer()
		//sleep 5 sec
		time.Sleep(5 * time.Second)
		core.InternalStartServer()

	default:
		fmt.Println("Unknown reaction:", r.Emoji.Name)
		return
	}

	// Get the user who triggered the action
	user, err := s.User(r.UserID)
	if err != nil {
		fmt.Printf("Error fetching user details: %v\n", err)
		return
	}
	username := user.Username

	// Send the action message to the error channel
	sendMessageToErrorChannel(fmt.Sprintf("%s triggered by %s.", actionMessage, username))

	// Remove the reaction after processing
	err = s.MessageReactionRemove(config.ErrorChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
	if err != nil {
		fmt.Printf("Error removing reaction: %v\n", err)
	}
}
