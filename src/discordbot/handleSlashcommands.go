package discordbot

import (
	"StationeersServerUI/src/backupmgr"
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/gamemgr"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// listenToSlashCommands handles slash command interactions
func listenToSlashCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if i.ChannelID != config.ControlChannelID {
		embed := generateEmbed(EmbedData{
			Title:       "Wrong Channel",
			Description: "Commands must be sent to the configured control channel",
			Color:       0xFF0000, // Red for error
			Fields: []EmbedField{
				{Name: "Accepted Channel", Value: fmt.Sprintf("<#%s>", config.ControlChannelID), Inline: true},
			},
		})
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to wrong channel command: %v\n", err)
		}
		return
	}

	cmdName := i.ApplicationCommandData().Name

	switch cmdName {
	case "start":
		embed := generateEmbed(EmbedData{
			Title:       "Server Control",
			Description: "Starting the server...",
			Color:       0x00FF00, // Green
			Fields: []EmbedField{
				{Name: "Status", Value: "🕛 Recieved", Inline: true},
			},
		})
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /start: %v\n", err)
			return
		}
		gamemgr.InternalStartServer()
		SendMessageToStatusChannel("🕛Start command received from Server Controller, Server is Starting...")

	case "stop":
		embed := generateEmbed(EmbedData{
			Title:       "Server Control",
			Description: "Stopping the server...",
			Color:       0xFF0000, // Red
			Fields: []EmbedField{
				{Name: "Status", Value: "🕛 Recieved", Inline: true},
			},
		})
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /stop: %v\n", err)
			return
		}
		gamemgr.InternalStopServer()
		SendMessageToStatusChannel("🕛Stop command received from Server Controller, flatlining Server in 5 Seconds...")

	case "status":
		isRunning := gamemgr.InternalIsServerRunning()
		embed := generateEmbed(EmbedData{
			Title:       "🎮 Server Status",
			Description: "Current process state for the Stationeers game server.\n*Note: 'Started' indicates a running game server process was found,\n but not necessarily that the Server configuration is correct \n or that the server is available to join.*",
			Color:       map[bool]int{true: 0x00FF00, false: 0xFF0000}[isRunning], // Green if running, Red if not
			Fields: []EmbedField{
				{Name: "Status:", Value: map[bool]string{true: "🟢 Started", false: "🔴 Stopped"}[isRunning], Inline: true},
				{Name: "Checked:", Value: time.Now().Format("15:04:05 MST"), Inline: true},
			},
		})
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /start: %v\n", err)
			return
		}

	case "help":
		embed := generateEmbed(EmbedData{
			Title:       "Command Help",
			Description: "Available Commands:",
			Color:       0x1E90FF, // Blue
			Fields: []EmbedField{
				{Name: "/start", Value: "Starts the server", Inline: false},
				{Name: "/stop", Value: "Stops the server", Inline: false},
				{Name: "/restore <index>", Value: "Restores a backup at the specified index", Inline: false},
				{Name: "/list [limit]", Value: "Lists recent backups (default: 5)", Inline: false},
				{Name: "/help", Value: "Shows this help message", Inline: false},
				{Name: "/bansteamid <SteamID>", Value: "Bans a player by SteamID", Inline: false},
				{Name: "/bansteamid <SteamID>", Value: "Unbans a player by SteamID", Inline: false},
			},
		})
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /help: %v\n", err)
		}

	case "restore":
		options := i.ApplicationCommandData().Options
		indexStr := options[0].StringValue()
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			embed := generateEmbed(EmbedData{
				Title:       "Restore Failed",
				Description: "Invalid index provided",
				Color:       0xFF0000,
				Fields: []EmbedField{
					{Name: "Error", Value: "Please provide a valid number", Inline: true},
				},
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			SendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
			return
		}

		embed := generateEmbed(EmbedData{
			Title:       "Backup Restore",
			Description: fmt.Sprintf("Restoring backup #%d...", index),
			Color:       0xFFA500, // Orange
			Fields: []EmbedField{
				{Name: "Status", Value: "🕛 Recieved", Inline: true},
			},
		})
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /restore: %v\n", err)
			return
		}

		gamemgr.InternalStopServer()
		err = backupmgr.GlobalBackupManager.RestoreBackup(index)
		if err != nil {
			SendMessageToControlChannel(fmt.Sprintf("❌Failed to restore backup at index %d: %v", index, err))
			SendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
			return
		}
		SendMessageToControlChannel(fmt.Sprintf("✅Backup %d restored successfully, Starting Server...", index))
		time.Sleep(5 * time.Second)
		gamemgr.InternalStartServer()

	case "list":
		var limit int = 5 // Default
		options := i.ApplicationCommandData().Options
		if len(options) > 0 {
			limitStr := options[0].StringValue()
			if strings.ToLower(limitStr) == "all" {
				limit = 0
			} else {
				var err error
				limit, err = strconv.Atoi(limitStr)
				if err != nil || limit < 1 {
					embed := generateEmbed(EmbedData{
						Title:       "List Failed",
						Description: "Invalid limit provided",
						Color:       0xFF0000,
						Fields: []EmbedField{
							{Name: "Error", Value: "Use a number or 'all'", Inline: true},
						},
					})
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{embed},
						},
					})
					return
				}
			}
		}

		backups, err := backupmgr.GlobalBackupManager.ListBackups(limit)
		if err != nil {
			embed := generateEmbed(EmbedData{
				Title:       "List Failed",
				Description: "Error fetching backups",
				Color:       0xFF0000,
				Fields: []EmbedField{
					{Name: "Error", Value: "Failed to fetch backup list", Inline: true},
				},
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}

		if len(backups) == 0 {
			embed := generateEmbed(EmbedData{
				Title:       "Backup List",
				Description: "No backups found",
				Color:       0xFFD700,
				Fields:      []EmbedField{},
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}

		// Sort backups by ModTime, newest first
		sort.Slice(backups, func(i, j int) bool {
			return backups[i].ModTime.After(backups[j].ModTime)
		})

		// Build embeds, 20 backups per embed
		batchSize := 20
		embeds := []*discordgo.MessageEmbed{}
		for i := 0; i < len(backups); i += batchSize {
			end := i + batchSize
			if end > len(backups) {
				end = len(backups)
			}
			batch := backups[i:end]

			fields := make([]EmbedField, len(batch))
			for j, backup := range batch {
				fields[j] = EmbedField{
					Name:   fmt.Sprintf("📂 Backup #%d", backup.Index),
					Value:  fmt.Sprintf("⏰ %s", backup.ModTime.Format("January 2, 2006, 3:04 PM")),
					Inline: false,
				}
			}

			embed := generateEmbed(EmbedData{
				Title:       "📜 Backup Archives",
				Description: fmt.Sprintf("Showing %d-%d of %d backups", i+1, end, len(backups)),
				Color:       0xFFD700,
				Fields:      fields,
			})
			embeds = append(embeds, embed)
		}

		// Send first embed as interaction response
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embeds[0]},
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /list with first embed: %v\n", err)
			return
		}

		// Send additional embeds
		for _, embed := range embeds[1:] {
			time.Sleep(500 * time.Millisecond)
			_, err = s.ChannelMessageSendEmbed(i.ChannelID, embed)
			if err != nil {
				fmt.Printf("[DISCORD] Error sending additional /list embed: %v\n", err)
			}
		}

	case "bansteamid":
		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			embed := generateEmbed(EmbedData{
				Title:       "Ban Failed",
				Description: "No SteamID provided",
				Color:       0xFF0000,
				Fields: []EmbedField{
					{Name: "Error", Value: "Please provide a SteamID", Inline: true},
				},
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}

		steamID := options[0].StringValue()
		err := banSteamID(steamID)
		if err != nil {
			embed := generateEmbed(EmbedData{
				Title:       "Ban Failed",
				Description: fmt.Sprintf("Could not ban SteamID %s", steamID),
				Color:       0xFF0000,
				Fields: []EmbedField{
					{Name: "Error", Value: err.Error(), Inline: true},
				},
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}

		embed := generateEmbed(EmbedData{
			Title:       "Player Banned",
			Description: fmt.Sprintf("SteamID %s has been banned", steamID),
			Color:       0xFF0000,
			Fields: []EmbedField{
				{Name: "Status", Value: "✅ Completed", Inline: true},
			},
		})
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /ban: %v\n", err)
		}

	case "unbansteamid":
		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			embed := generateEmbed(EmbedData{
				Title:       "Unban Failed",
				Description: "No SteamID provided",
				Color:       0xFF0000,
				Fields: []EmbedField{
					{Name: "Error", Value: "Please provide a SteamID", Inline: true},
				},
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}

		steamID := options[0].StringValue()
		err := unbanSteamID(steamID)
		if err != nil {
			embed := generateEmbed(EmbedData{
				Title:       "Unban Failed",
				Description: fmt.Sprintf("Could not unban SteamID %s", steamID),
				Color:       0xFF0000,
				Fields: []EmbedField{
					{Name: "Error", Value: err.Error(), Inline: true},
				},
			})
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}

		embed := generateEmbed(EmbedData{
			Title:       "Player Unbanned",
			Description: fmt.Sprintf("SteamID %s has been unbanned", steamID),
			Color:       0x00FF00,
			Fields: []EmbedField{
				{Name: "Status", Value: "✅ Completed", Inline: true},
			},
		})
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /unban: %v\n", err)
		}

	}
}
