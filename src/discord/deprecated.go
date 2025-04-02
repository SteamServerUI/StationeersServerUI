package discord

import (
	"StationeersServerUI/src/config"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

/*
Unused Discord functionality from v2 (not used in v4):
In v2, we had a restart button on the last exception message.
In v2, we had a connected players message in a dedicated channel that was updated every time a player connected or disconnected.
In v2, the backup system was usable from discord. This is no longer the case. I will add it back in the future.
In v2, update and validate commands were usable from discord. This is no longer the case.
	I will add some kind of Update & Validate functionality back in the future wich stops the server, runs SteamCMD, and restarts the server again.
		This would also allow to change the game branch without restarting the Software.
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

// DEPRECATED
func handleUpdateCommand(s *discordgo.Session, channelID string) {
	// Notify that the update process is starting
	s.ChannelMessageSend(channelID, "🙏Sorry, this feature has been deprecated. Server Updates are now handled automatically at Software Startup. If you are interested in bringing this feature back, please report it on the GitHub repository. We will be happy to implement it.")
}

// DEPRECATED
func handleValidateCommand(s *discordgo.Session, channelID string) {
	// Notify that the update process is starting
	s.ChannelMessageSend(channelID, "🙏Sorry, this feature has been deprecated. Server File Validation is now handled automatically at Software Startup. If you are interested in bringing this feature back, please report it on the GitHub repository. We will be happy to implement it.")
}
