package discordbot

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/src/backupmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/src/gamemgr"
	"github.com/JacksonTheMaster/StationeersServerUI/src/logger"

	"github.com/bwmarrin/discordgo"
)

type commandHandler func(*discordgo.Session, *discordgo.InteractionCreate, EmbedData) error

// Command handlers map
var handlers = map[string]commandHandler{
	"start":        handleStart,
	"stop":         handleStop,
	"status":       handleStatus,
	"help":         handleHelp,
	"restore":      handleRestore,
	"list":         handleList,
	"bansteamid":   handleBan,
	"unbansteamid": handleUnban,
}

// Check channel and handle initial validation
func listenToSlashCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand || i.ChannelID != config.ControlChannelID {
		respond(s, i, EmbedData{
			Title: "Wrong Channel", Description: "Commands must be sent to the configured control channel",
			Color: 0xFF0000, Fields: []EmbedField{{Name: "Accepted Channel", Value: fmt.Sprintf("<#%s>", config.ControlChannelID), Inline: true}},
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
	isRunning := gamemgr.InternalIsServerRunning()
	data.Title = "🎮 Server Status"
	data.Description = "Current process state for the Stationeers game server.\n*Note: 'Started' indicates a running process was found, but not necessarily fully operational.*"
	data.Color = map[bool]int{true: 0x00FF00, false: 0xFF0000}[isRunning]
	data.Fields = []EmbedField{
		{Name: "Status:", Value: map[bool]string{true: "🟢 Started", false: "🔴 Stopped"}[isRunning], Inline: true},
		{Name: "Checked:", Value: time.Now().Format("15:04:05 MST"), Inline: true},
	}
	return respond(s, i, data)
}

func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title, data.Description, data.Color = "Command Help", "Available Commands:", 0x1E90FF
	data.Fields = []EmbedField{
		{Name: "/start", Value: "Starts the server"},
		{Name: "/stop", Value: "Stops the server"},
		{Name: "/restore <index>", Value: "Restores a backup"},
		{Name: "/list [limit]", Value: "Lists recent backups (default: 5)"},
		{Name: "/help", Value: "Shows this help"},
		{Name: "/bansteamid <SteamID>", Value: "Bans a player"},
		{Name: "/unbansteamid <SteamID>", Value: "Unbans a player"},
	}
	return respond(s, i, data)
}

func handleRestore(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	index, err := strconv.Atoi(i.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		data.Title, data.Description = "Restore Failed", "Invalid index provided"
		data.Fields = []EmbedField{{Name: "Error", Value: "Please provide a valid number", Inline: true}}
		return respond(s, i, data)
	}
	data.Title, data.Description, data.Color = "Backup Restore", fmt.Sprintf("Restoring backup #%d...", index), 0xFFA500
	data.Fields = []EmbedField{{Name: "Status", Value: "🕛 Recieved", Inline: true}}
	if err := respond(s, i, data); err != nil {
		return err
	}
	gamemgr.InternalStopServer()
	if err := backupmgr.GlobalBackupManager.RestoreBackup(index); err != nil {
		SendMessageToControlChannel(fmt.Sprintf("❌Failed to restore backup %d: %v", index, err))
		SendMessageToStatusChannel("⚠️Restore command failed")
		return nil
	}
	SendMessageToControlChannel(fmt.Sprintf("✅Backup %d restored, Starting Server...", index))
	time.Sleep(5 * time.Second)
	gamemgr.InternalStartServer()
	return nil
}

func handleList(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	limit := 5
	if len(i.ApplicationCommandData().Options) > 0 {
		limitStr := i.ApplicationCommandData().Options[0].StringValue()
		if strings.ToLower(limitStr) == "all" {
			limit = 0
		} else if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		} else {
			data.Title, data.Description = "List Failed", "Invalid limit provided"
			data.Fields = []EmbedField{{Name: "Error", Value: "Use a number or 'all'", Inline: true}}
			return respond(s, i, data)
		}
	}

	backups, err := backupmgr.GlobalBackupManager.ListBackups(limit)
	if err != nil {
		data.Title, data.Description = "List Failed", "Error fetching backups"
		data.Fields = []EmbedField{{Name: "Error", Value: "Failed to fetch backup list", Inline: true}}
		return respond(s, i, data)
	}
	if len(backups) == 0 {
		data.Title, data.Description, data.Color = "Backup List", "No backups found", 0xFFD700
		return respond(s, i, data)
	}

	sort.Slice(backups, func(i, j int) bool { return backups[i].ModTime.After(backups[j].ModTime) })
	batchSize := 20
	embeds := []*discordgo.MessageEmbed{}
	for start := 0; start < len(backups); start += batchSize {
		end := start + batchSize
		if end > len(backups) {
			end = len(backups)
		}
		fields := make([]EmbedField, end-start)
		for j, b := range backups[start:end] {
			fields[j] = EmbedField{Name: fmt.Sprintf("📂 Backup #%d", b.Index), Value: b.ModTime.Format("January 2, 2006, 3:04 PM")}
		}
		embeds = append(embeds, generateEmbed(EmbedData{
			Title: "📜 Backup Archives", Description: fmt.Sprintf("Showing %d-%d of %d backups", start+1, end, len(backups)),
			Color: 0xFFD700, Fields: fields,
		}))
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embeds[0]}},
	}); err != nil {
		return err
	}
	for _, embed := range embeds[1:] {
		time.Sleep(500 * time.Millisecond)
		s.ChannelMessageSendEmbed(i.ChannelID, embed)
	}
	return nil
}

func handleBan(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	return handleBanUnban(s, i, data, banSteamID, "Banned", "Ban Failed", 0xFF0000)
}

func handleUnban(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	return handleBanUnban(s, i, data, unbanSteamID, "Unbanned", "Unban Failed", 0x00FF00)
}

func handleBanUnban(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData, fn func(string) error, successTitle, failTitle string, color int) error {
	if len(i.ApplicationCommandData().Options) == 0 {
		data.Title, data.Description = failTitle, "No SteamID provided"
		data.Fields = []EmbedField{{Name: "Error", Value: "Please provide a SteamID", Inline: true}}
		return respond(s, i, data)
	}
	steamID := i.ApplicationCommandData().Options[0].StringValue()
	if err := fn(steamID); err != nil {
		data.Title, data.Description = failTitle, fmt.Sprintf("Could not %s SteamID %s", strings.ToLower(failTitle[:len(failTitle)-6]), steamID)
		data.Fields = []EmbedField{{Name: "Error", Value: err.Error(), Inline: true}}
		return respond(s, i, data)
	}
	data.Title, data.Description, data.Color = successTitle, fmt.Sprintf("SteamID %s has been %s", steamID, strings.ToLower(successTitle)), color
	data.Fields = []EmbedField{{Name: "Status", Value: "✅ Completed", Inline: true}}
	return respond(s, i, data)
}
