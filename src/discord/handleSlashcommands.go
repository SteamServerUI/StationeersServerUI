package discord

import (
	"StationeersServerUI/src/backupsv2"
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/core"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// registerSlashCommands defines and registers slash commands
func registerSlashCommands(s *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "start",
			Description: "Start the server",
		},
		{
			Name:        "stop",
			Description: "Stop the server",
		},
		{
			Name:        "help",
			Description: "Show command help",
		},
		{
			Name:        "restore",
			Description: "Restore a backup at the specified index",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "index",
					Description: "Backup index to restore (e.g., 1)",
					Required:    true,
				},
			},
		},
		{
			Name:        "list",
			Description: "List the most recent backups",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "limit",
					Description: "Number of backups to list or 'all' (default: 5)",
					Required:    false,
				},
			},
		},
	}

	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			fmt.Printf("[DISCORD] Error registering command %s: %v\n", cmd.Name, err)
		} else {
			fmt.Printf("[DISCORD] Successfully registered command: %s\n", cmd.Name)
		}
	}
}

// listenToSlashCommands handles slash command interactions
func listenToSlashCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if i.ChannelID != config.ControlChannelID {
		return
	}

	cmdName := i.ApplicationCommandData().Name

	switch cmdName {
	case "start":
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🕛Server is starting...",
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /start: %v\n", err)
			return
		}
		core.InternalStartServer()
		SendMessageToStatusChannel("🕛Start command received from Server Controller, Server is Starting...")

	case "stop":
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🕛Server is stopping...",
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /stop: %v\n", err)
			return
		}
		core.InternalStopServer()
		SendMessageToStatusChannel("🕛Stop command received from Server Controller, flatlining Server in 5 Seconds...")

	case "help":
		helpMessage := `
**Available Commands:**
- ` + "`/start`" + `: Starts the server.
- ` + "`/stop`" + `: Stops the server.
- ` + "`/restore <index>`" + `: Restores a backup at the specified index. Usage: ` + "`/restore 1`" + `.
- ` + "`/list [limit]`" + `: Lists the most recent backups. Use ` + "`/list all`" + ` or ` + "`/list 3`" + `.
- ` + "`/help`" + `: Displays this help message.
- ` + "`!ban:<SteamID>`" + `: Bans a player by their SteamID. Usage: ` + "`!ban:76561198334231312`" + `.
- ` + "`!unban:<SteamID>`" + `: Unbans a player by their SteamID. Usage: ` + "`!unban:76561198334231312`" + `.
`
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: helpMessage,
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /help: %v\n", err)
		}

	case "restore":
		options := i.ApplicationCommandData().Options
		indexStr := options[0].StringValue() // Required option
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌Invalid index provided for restore.",
				},
			})
			SendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
			return
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("🕛Restoring backup %d...", index),
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /restore: %v\n", err)
			return
		}

		core.InternalStopServer()
		err = backupsv2.GlobalBackupManager.RestoreBackup(index)
		if err != nil {
			SendMessageToControlChannel(fmt.Sprintf("❌Failed to restore backup at index %d: %v", index, err))
			SendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
			return
		}
		SendMessageToControlChannel(fmt.Sprintf("✅Backup %d restored successfully, Starting Server...", index))
		time.Sleep(5 * time.Second)
		core.InternalStartServer()

	case "list":
		var limit int = 5 // Default
		options := i.ApplicationCommandData().Options
		if len(options) > 0 {
			limitStr := options[0].StringValue()
			if strings.ToLower(limitStr) == "all" {
				limit = 0 // 0 means all in ListBackups
			} else {
				var err error
				limit, err = strconv.Atoi(limitStr)
				if err != nil || limit < 1 {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "❌Invalid limit provided. Use a number or 'all'.",
						},
					})
					return
				}
			}
		}

		backups, err := backupsv2.GlobalBackupManager.ListBackups(limit)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌Failed to fetch backup list.",
				},
			})
			return
		}

		if len(backups) == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "No backups found.",
				},
			})
			return
		}

		// Build a table-like response
		table := "```\nIndex | Created\n------|--------------------\n"
		for _, backup := range backups {
			table += fmt.Sprintf("%5d | %s\n", backup.Index, backup.ModTime.Format("2006-01-02 15:04:05"))
		}
		table += "```"

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: table,
			},
		})
		if err != nil {
			fmt.Printf("[DISCORD] Error responding to /list: %v\n", err)
		}
	}
}
