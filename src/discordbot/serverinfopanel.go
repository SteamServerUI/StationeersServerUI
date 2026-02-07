package discordbot

import (
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/bwmarrin/discordgo"
)

// ServerInfoMessageID stores the ID of the server info panel message
var ServerInfoMessageID string

// Custom IDs for button interactions
const (
	ButtonGetPassword       = "ssui_get_password"
	ButtonDownloadBackupPfx = "ssui_download_backup_" // Prefix for download backup button. Has NOTHING TO DO WITH SERVERINFOPANEL, but this file felt like the best place to put it for now
)

// sendServerInfoPanel sends the server info panel with interactive buttons
func sendServerInfoPanel() {
	if !config.GetIsDiscordEnabled() {
		return
	}

	channelID := config.GetServerInfoPanelChannelID()
	if channelID == "" {
		logger.Discord.Debug("Server info panel channel ID not configured, skipping panel")
		return
	}

	serverName := config.GetServerName()

	// Create an embed for the server info panel
	embed := &discordgo.MessageEmbed{
		Title:       "🎮 Server Information",
		Description: "Your Stationeers server is running as **" + serverName + "**",
		Color:       0x5865F2, // Discord blurple color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Server Name",
				Value:  serverName,
				Inline: true,
			},
			{
				Name:   "SSUI Version",
				Value:  config.GetVersion(),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Build message components (buttons in an action row)
	var components []discordgo.MessageComponent

	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "🔑 Get Current Server Password",
				Style:    discordgo.PrimaryButton,
				CustomID: ButtonGetPassword,
			},
		},
	})

	// Send the message with embed and buttons
	msg, err := config.DiscordSession.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		logger.Discord.Error("Error sending server info panel: " + err.Error())
		return
	}

	clearMessagesAboveLastN(channelID, 1) // Clear all old server info panel messages
	ServerInfoMessageID = msg.ID
	logger.Discord.Info("Server info panel sent successfully")
}

// handleServerInfoButtonInteraction handles button interactions from the server info panel
func handleServerInfoButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Only handle component interactions (buttons)
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// Get the custom ID of the button that was clicked
	customID := i.MessageComponentData().CustomID

	switch customID {
	case ButtonGetPassword:
		handleGetPasswordButton(s, i)
	default:
		// Not our button, ignore
		return
	}
}

// handleGetPasswordButton sends the current server password as an ephemeral message (only visible to the user)
func handleGetPasswordButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Get the current server password from config
	password := config.GetServerPassword()

	// Build the embed based on whether password is set
	var embed *discordgo.MessageEmbed
	if password == "" {
		embed = &discordgo.MessageEmbed{
			Title:       "🔓 No Password Set",
			Description: "No password is currently configured for this server.",
			Color:       0xFFA500, // Orange
			Footer: &discordgo.MessageEmbedFooter{
				Text: "This message will disappear in 30 seconds",
			},
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Title:       "🔑 Current Server Password",
			Description: "Use this password to connect to the server.",
			Color:       0x00FF00, // Green
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Password",
					Value:  "```" + password + "```",
					Inline: false,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "🔒 Only visible to you • Disappears in 30 seconds",
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	// Respond with an ephemeral message (only visible to the user who clicked)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Discord.Error("Error responding to get password button: " + err.Error())
		return
	}

	go func() {
		time.Sleep(30 * time.Second)
		err := s.InteractionResponseDelete(i.Interaction)
		if err != nil {
			logger.Discord.Debug("Could not delete ephemeral password message: " + err.Error())
		}
	}()
}
