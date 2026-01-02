package discordbot

import (
	"fmt"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/commandmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamcmd"

	"github.com/bwmarrin/discordgo"
)

type commandHandler func(*discordgo.Session, *discordgo.InteractionCreate, EmbedData) error

// Command handlers map
var handlers = map[string]commandHandler{
	"start":   handleStart,
	"stop":    handleStop,
	"status":  handleStatus,
	"help":    handleHelp,
	"update":  handleUpdate,
	"command": handleCommand,
}

// Check channel and handle initial validation
func listenToSlashCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand || i.ChannelID != config.GetControlChannelID() {
		respond(s, i, EmbedData{
			Title: "Wrong Channel", Description: "Commands must be sent to the configured control channel",
			Color: 0xFF0000, Fields: []EmbedField{{Name: "Accepted Channel", Value: fmt.Sprintf("<#%s>", config.GetControlChannelID()), Inline: true}},
		})
		return
	}

	cmd := i.ApplicationCommandData().Name
	if handler, ok := handlers[cmd]; ok {
		data := EmbedData{Title: "Command Error", Color: 0xFF0000}
		if err := handler(s, i, data); err != nil {
			logger.Discord.Error("Error handling " + cmd + ": " + err.Error())
		}
	}
}

// Generic response function
func respond(s *discordgo.Session, i *discordgo.InteractionCreate, embed EmbedData) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{generateEmbed(embed)},
		},
	})
}

func handleStart(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title, data.Description, data.Color = "Server Control", "Starting the server...", 0x00FF00
	data.Fields = []EmbedField{{Name: "Status", Value: "🕛 Recieved", Inline: true}}
	if err := respond(s, i, data); err != nil {
		return err
	}
	gamemgr.InternalStartServer()
	SendMessageToStatusChannel("🕛Start command received, Server is Starting...")
	return nil
}

func handleStop(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title, data.Description, data.Color = "Server Control", "Stopping the server...", 0xFF0000
	data.Fields = []EmbedField{{Name: "Status", Value: "🕛 Recieved", Inline: true}}
	if err := respond(s, i, data); err != nil {
		return err
	}
	gamemgr.InternalStopServer()
	SendMessageToStatusChannel("🕛Stop command received, flatlining Server in 5 Seconds...")
	return nil
}

func handleStatus(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	isRunning := config.GetIsGameServerRunning()
	data.Title = "🎮 Server Status"
	data.Description = "Current process state for the Stationeers game server.\n*Note: 'Started' indicates a running process was found, but not necessarily fully operational.*"
	data.Color = map[bool]int{true: 0x00FF00, false: 0xFF0000}[isRunning]
	data.Fields = []EmbedField{
		{Name: "Status:", Value: map[bool]string{true: "🟢 Started", false: "🔴 Stopped"}[isRunning], Inline: true},
		{Name: "Checked:", Value: time.Now().Format("15:04:05 MST"), Inline: true},
	}
	return respond(s, i, data)
}

func handleUpdate(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	thinkingData := EmbedData{
		Title:       "🎮 Gameserver Update",
		Description: "The Backend is processing the gameserver update via SteamCMD. Please wait, this may take a while...",
		Color:       0xFFA500, // Orange color for in-progress
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{generateEmbed(thinkingData)},
		},
	})
	if err != nil {
		return err
	}
	data.Title = "🎮 Gameserver Update"
	data.Description = "Gameserver update completed."
	data.Color = 0x00FF00 // Green for completion (will adjust if error)

	_, err = steamcmd.InstallAndRunSteamCMD()

	data.Fields = []EmbedField{
		{Name: "Update Status:", Value: map[bool]string{true: "🟢 Success", false: "🔴 Failed"}[err == nil], Inline: true},
	}
	if err != nil {
		data.Color = 0xFF0000 // Red for error
		data.Fields = append(data.Fields, EmbedField{Name: "Error:", Value: err.Error(), Inline: true})
	}

	// Edit the original message with "update completed" embed
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{generateEmbed(data)},
	})
	return err
}

func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title, data.Description, data.Color = "Command Help", "Available Commands:", 0x1E90FF
	data.Fields = []EmbedField{
		{Name: "/start", Value: "Starts the server"},
		{Name: "/stop", Value: "Stops the server"},
		{Name: "/status", Value: "Gets the running status of the gameserver process"},
		{Name: "/update", Value: "Updates the gameserver via SteamCMD"},
		{Name: "/list [limit]", Value: "Lists recent backups (default: 5)"},
		{Name: "/restore <index>", Value: "Restores a backup"},
		{Name: "/bansteamid <SteamID>", Value: "Bans a player"},
		{Name: "/unbansteamid <SteamID>", Value: "Unbans a player"},
		{Name: "/command <command>", Value: "Sends a command to the gameserver console"},
		{Name: "/help", Value: "Shows this help"},
	}
	return respond(s, i, data)
}

func handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title, data.Description, data.Color = "Server Control", "Sending a command to the gameserver console...", 0x00FF00
	data.Fields = []EmbedField{{Name: "Status", Value: "❌ Failed, is the server running and SSCM enabled?", Inline: true}}
	data.Color = 0xFF0000
	if config.GetIsGameServerRunning() {
		data.Color = 0x00FF00
		err := commandmgr.WriteCommand(i.ApplicationCommandData().Options[0].StringValue())
		if err != nil {
			data.Fields = []EmbedField{{Name: "Error", Value: err.Error(), Inline: true}}
			return respond(s, i, data)
		}
		data.Fields = []EmbedField{{Name: "Status", Value: "✅ Gameserver recieved command", Inline: true}}
	}

	if err := respond(s, i, data); err != nil {
		return err
	}
	return nil
}
