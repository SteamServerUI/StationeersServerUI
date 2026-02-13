package discordbot

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/bwmarrin/discordgo"
)

var (
	connectedPlayersMessageID string // connectedPlayersMessageID tracks the message ID for editing the connected players message
	playersMutex              sync.Mutex
)

// sendConnectedPlayersPanel sends the initial "no players" embed on startup
func sendConnectedPlayersPanel() {
	if !config.GetIsDiscordEnabled() {
		return
	}

	channelID := config.GetConnectionListChannelID()
	if channelID == "" {
		logger.Discord.Debug("Connection list channel ID not configured, skipping panel")
		return
	}

	embed := buildConnectedPlayersEmbed(nil)
	sendOrEditConnectedPlayersEmbed(channelID, embed)
	clearMessagesAboveLastN(channelID, 1)
	logger.Discord.Info("Connected players panel sent successfully")
}

func AddToConnectedPlayers(username, steamID string, connectionTime time.Time, players map[string]string) {
	if !config.GetIsDiscordEnabled() || config.DiscordSession == nil {
		logger.Discord.Debug("Discord not enabled or session not initialized")
		return
	}
	channelID := config.GetConnectionListChannelID()
	if channelID == "" {
		return
	}
	embed := buildConnectedPlayersEmbed(players)
	sendOrEditConnectedPlayersEmbed(channelID, embed)
}

func RemoveFromConnectedPlayers(steamID string, players map[string]string) {
	if !config.GetIsDiscordEnabled() || config.DiscordSession == nil {
		logger.Discord.Debug("Discord not enabled or session not initialized")
		return
	}
	channelID := config.GetConnectionListChannelID()
	if channelID == "" {
		return
	}
	embed := buildConnectedPlayersEmbed(players)
	sendOrEditConnectedPlayersEmbed(channelID, embed)
}

// buildConnectedPlayersEmbed creates a Discord embed for the connected players panel
func buildConnectedPlayersEmbed(players map[string]string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:     "👥 Connected Players",
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Last updated",
		},
	}

	if len(players) == 0 {
		embed.Description = "No players are currently connected."
		embed.Color = 0x95A5A6 // Grey
		return embed
	}

	embed.Color = 0x2ECC71 // Green

	// Build a clean row-based player list in the description
	var lines strings.Builder
	fmt.Fprintf(&lines, "**%d** player(s) online, click opens Steam profile\n\n", len(players))
	for steamID, username := range players {
		fmt.Fprintf(&lines, "👤 [%s](https://steamcommunity.com/profiles/%s/)\n", username, steamID)
	}
	embed.Description = lines.String()

	return embed
}

// sendOrEditConnectedPlayersEmbed sends a new embed or edits the existing one
func sendOrEditConnectedPlayersEmbed(channelID string, embed *discordgo.MessageEmbed) {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	if connectedPlayersMessageID == "" {
		// Send a new message with embed
		msg, err := config.DiscordSession.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{embed},
		})
		if err != nil {
			logger.Discord.Error("Error sending connected players embed to channel " + channelID + ": " + err.Error())
			return
		}
		connectedPlayersMessageID = msg.ID
		logger.Discord.Debug("Sent connected players embed to channel " + channelID)
	} else {
		// Edit the existing message with the updated embed
		embeds := []*discordgo.MessageEmbed{embed}
		content := ""
		_, err := config.DiscordSession.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: channelID,
			ID:      connectedPlayersMessageID,
			Content: &content,
			Embeds:  &embeds,
		})
		if err != nil {
			logger.Discord.Error("Error editing connected players embed in channel " + channelID + ": " + err.Error())
			// If editing fails (e.g., message deleted), reset and try sending a new one
			connectedPlayersMessageID = ""
			msg, err := config.DiscordSession.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{embed},
			})
			if err != nil {
				logger.Discord.Error("Error sending fallback connected players embed to channel " + channelID + ": " + err.Error())
			} else {
				connectedPlayersMessageID = msg.ID
				logger.Discord.Debug("Sent new connected players embed after edit failure to channel " + channelID)
			}
		}
	}
}
