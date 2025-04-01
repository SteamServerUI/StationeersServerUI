package discord

import (
	"StationeersServerUI/src/config"
	"fmt"
	"strings"
)

/*
Unused Discord functions from v2 (not used in v4):
*/

func updateConnectedPlayersMessage(channelID string) {
	content := formatConnectedPlayers()
	sendAndEditMessageInConnectedPlayersChannel(channelID, content)
}

func sendAndEditMessageInConnectedPlayersChannel(channelID, message string) {
	if !config.IsDiscordEnabled {
		return
	}
	if config.DiscordSession == nil {
		fmt.Println("[DISCORD] Discord Error: Discord is enabled but session is not initialized")
		return
	}
	//only clear messages if we are on release branch
	if config.Branch == "release" || config.Branch == "Release" {
		clearMessagesAboveLastN(config.ControlChannelID, 1)
	}
	if config.ConnectedPlayersMessageID == "" {
		// Send a new message if there's no existing message to edit
		msg, err := config.DiscordSession.ChannelMessageSend(channelID, message)
		if err != nil {
			fmt.Printf("[DISCORD] Error sending message to channel %s: %v\n", channelID, err)
		} else {
			config.ConnectedPlayersMessageID = msg.ID
			fmt.Printf("Sent message to channel %s: %s\n", channelID, message)
		}
	} else {
		// Edit the existing message
		_, err := config.DiscordSession.ChannelMessageEdit(channelID, config.ConnectedPlayersMessageID, message)
		if err != nil {
			fmt.Printf("[DISCORD]Error editing message in channel %s: %v\n", channelID, err)
		} else {
			fmt.Printf("[DISCORD] Updated message in channel %s: %s\n", channelID, message)
		}
	}
}

func formatConnectedPlayers() string {
	if len(config.ConnectedPlayers) == 0 {
		return "No players are currently connected."
	}

	var sb strings.Builder
	sb.WriteString("**Connected Players:**\n")
	sb.WriteString("```\n")
	sb.WriteString("Username              | Steam ID\n")
	sb.WriteString("----------------------|------------------------\n")

	for steamID, username := range config.ConnectedPlayers {
		sb.WriteString(fmt.Sprintf("%-20s | %s\n", username, steamID))
	}

	sb.WriteString("```")
	return sb.String()
}
