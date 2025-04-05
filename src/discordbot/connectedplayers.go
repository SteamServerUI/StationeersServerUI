package discordbot

import (
	"StationeersServerUI/src/config"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	connectedPlayersMessageID string // connectedPlayersMessageID tracks the message ID for editing the connected players message
	playersMutex              sync.Mutex
)

func AddToConnectedPlayers(username, steamID string, connectionTime time.Time, players map[string]string) {
	if !config.IsDiscordEnabled || config.DiscordSession == nil {
		fmt.Println("[DISCORD] Discord not enabled or session not initialized")
		return
	}
	content := formatConnectedPlayers(players)
	sendAndEditMessageInConnectedPlayersChannel(config.ConnectionListChannelID, content)
}

func RemoveFromConnectedPlayers(steamID string, players map[string]string) {
	if !config.IsDiscordEnabled || config.DiscordSession == nil {
		fmt.Println("[DISCORD] Discord not enabled or session not initialized")
		return
	}
	content := formatConnectedPlayers(players)
	sendAndEditMessageInConnectedPlayersChannel(config.ConnectionListChannelID, content)
}

func sendAndEditMessageInConnectedPlayersChannel(channelID, message string) {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	if connectedPlayersMessageID == "" {
		// Send a new message if there's no existing one
		msg, err := config.DiscordSession.ChannelMessageSend(channelID, message)
		if err != nil {
			fmt.Printf("[DISCORD] Error sending message to channel %s: %v\n", channelID, err)
			return
		}
		connectedPlayersMessageID = msg.ID
		fmt.Printf("[DISCORD] Sent new message to channel %s with ID %s\n", channelID, msg.ID)
	} else {
		// Edit the existing message
		_, err := config.DiscordSession.ChannelMessageEdit(channelID, connectedPlayersMessageID, message)
		if err != nil {
			fmt.Printf("[DISCORD] Error editing message in channel %s: %v\n", channelID, err)
			// If editing fails (e.g., message deleted), reset and try sending a new one
			connectedPlayersMessageID = ""
			msg, err := config.DiscordSession.ChannelMessageSend(channelID, message)
			if err != nil {
				fmt.Printf("[DISCORD] Error sending fallback message to channel %s: %v\n", channelID, err)
			} else {
				connectedPlayersMessageID = msg.ID
				fmt.Printf("[DISCORD] Sent new message after edit failure to channel %s with ID %s\n", channelID, msg.ID)
			}
		} else {
			fmt.Printf("[DISCORD] Updated message in channel %s with ID %s\n", channelID, connectedPlayersMessageID)
		}
	}
}

func formatConnectedPlayers(players map[string]string) string {
	if len(players) == 0 {
		return "No players are currently connected."
	}

	var sb strings.Builder
	sb.WriteString("Connected Players:\n")
	sb.WriteString("```\n")
	sb.WriteString("Username              | Steam ID\n")
	sb.WriteString("----------------------|------------------------\n")

	for steamID, username := range players {
		sb.WriteString(fmt.Sprintf("%-20s | %s\n", username, steamID))
	}

	sb.WriteString("```")
	return sb.String()
}
