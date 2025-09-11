package discordbot

import (
	"fmt"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/steamcmd"
	"github.com/bwmarrin/discordgo"
)

func sendControlPanel() {
	if !config.GetIsDiscordEnabled() {
		return
	}

	// Create an embed for the control panel
	embed := &discordgo.MessageEmbed{
		Title:       "🚀 SSUI Control Panel",
		Description: "Use the reactions below to manage the server:",
		Color:       0x1e90ff, // Vibrant blue color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🟢 Start",
				Value:  "Launch the gameserver",
				Inline: false,
			},
			{
				Name:   "🔴 Stop",
				Value:  "Stop the gameserver",
				Inline: false,
			},
			{
				Name:   "🔄 Restart",
				Value:  "Restart the gameserver",
				Inline: false,
			},
			{
				Name:   "♻️ Update",
				Value:  "Update the gameserver via SteamCMD",
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Send the embed message
	msg, err := config.DiscordSession.ChannelMessageSendEmbed(config.GetControlPanelChannelID(), embed)
	if err != nil {
		logger.Discord.Error("Error sending control panel embed: " + err.Error())
		return
	}

	clearMessagesAboveLastN(config.GetControlPanelChannelID(), 1) // Clear all old control panel messages

	// Add reactions (acting as buttons) to the control message
	config.DiscordSession.MessageReactionAdd(config.GetControlPanelChannelID(), msg.ID, "🟢")  // Start
	config.DiscordSession.MessageReactionAdd(config.GetControlPanelChannelID(), msg.ID, "🔴")  // Stop
	config.DiscordSession.MessageReactionAdd(config.GetControlPanelChannelID(), msg.ID, "🔄")  // Restart
	config.DiscordSession.MessageReactionAdd(config.GetControlPanelChannelID(), msg.ID, "♻️") // Update
	ControlMessageID = msg.ID
}

func handleControlReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Ignore reactions from the bot itself
	if r.UserID == s.State.User.ID {
		return
	}

	var actionMessage string

	switch r.Emoji.Name {
	case "🟢": // Start action
		gamemgr.InternalStartServer()
		actionMessage = "🟢 Server is Starting..."
	case "🔴": // Stop action
		gamemgr.InternalStopServer()
		actionMessage = "🔴 Server is Stopping..."
	case "🔄": // Restart action
		actionMessage = "🔄 Server is restarting..."
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
	case "♻️": // Update action
		actionMessage = "♻️ Server is updating, this may take a while..."
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

			_, err := steamcmd.InstallAndRunSteamCMD()

			Value := map[bool]string{true: "🟢 Success", false: "🔴 Failed"}[err == nil]
			SendMessageToStatusChannel(fmt.Sprintf("SteamCMD Update status: %s", Value))
			sendTemporaryMessage(s, config.GetControlPanelChannelID(), fmt.Sprintf("SteamCMD Update status: %s", Value), 30*time.Second)
			if err != nil {
				SendMessageToStatusChannel(fmt.Sprintf("Update failed: %v", err.Error()))
			}
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

	// Send a temporary confirmation message to the control panel channel
	sendTemporaryMessage(s, config.GetControlPanelChannelID(), actionMessage, 30*time.Second)

	// Send the action message to the status channel
	SendMessageToStatusChannel(fmt.Sprintf("%s triggered by %s.", actionMessage, username))

	// Remove the reaction after processing
	err = s.MessageReactionRemove(config.GetControlPanelChannelID(), r.MessageID, r.Emoji.APIName(), r.UserID)
	if err != nil {
		logger.Discord.Error("Error removing reaction: " + err.Error())
	}
}
