package discord

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/core"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// registerSlashCommands defines and registers the /start and /stop slash commands
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
	}

	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd) // GuildID is empty for global commands
		if err != nil {
			fmt.Printf("[DISCORD] Error registering command %s: %v\n", cmd.Name, err)
		} else {
			fmt.Printf("[DISCORD] Successfully registered command: %s\n", cmd.Name)
		}
	}
}

// listenToSlashCommands handles slash command interactions
func listenToSlashCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Only process slash commands (ignore other interaction types)
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Restrict to control channel
	if i.ChannelID != config.ControlChannelID {
		return
	}

	// Get the command name
	cmdName := i.ApplicationCommandData().Name

	switch cmdName {
	case "start":
		// Acknowledge the interaction and send a response
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
		// Execute the command and update status
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
		// Execute the command and update status
		core.InternalStopServer()
		SendMessageToStatusChannel("🕛Stop command received from Server Controller, flatlining Server in 5 Seconds...")
	}
}
