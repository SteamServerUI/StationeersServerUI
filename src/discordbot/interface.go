package discordbot

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/logger"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// InitializeDiscordBot starts or restarts the Discord bot and connects it to the Discord API.
func InitializeDiscordBot() {
	var err error

	// Clean up previous session
	if config.DiscordSession != nil {
		if config.IsDebugMode {
			fmt.Println("[DISCORD] Previous Discord session found, closing it...")
		}
		config.DiscordSession.Close()
	}
	if config.BufferFlushTicker != nil {
		config.BufferFlushTicker.Stop()
	}

	// Create new session
	config.DiscordSession, err = discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		fmt.Println("[DISCORD] Error creating Discord session:", err)
		return
	}

	// Set intents
	config.DiscordSession.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsMessageContent

	logger.Discord.Info("Starting Discord integration...")
	logger.Discord.Debug("Discord token: " + config.DiscordToken)
	logger.Discord.Debug("ControlChannelID: " + config.ControlChannelID)
	logger.Discord.Debug("StatusChannelID: " + config.StatusChannelID)
	logger.Discord.Debug("ConnectionListChannelID: " + config.ConnectionListChannelID)
	logger.Discord.Debug("LogChannelID: " + config.LogChannelID)
	logger.Discord.Debug("SaveChannelID: " + config.SaveChannelID)

	// Open session first
	err = config.DiscordSession.Open()
	if err != nil {
		fmt.Println("[DISCORD] Error opening Discord connection:", err)
		return
	}

	// Register handlers and commands after session is open
	config.DiscordSession.AddHandler(listenToDiscordMessages)
	config.DiscordSession.AddHandler(listenToDiscordReactions)
	config.DiscordSession.AddHandler(listenToSlashCommands)
	registerSlashCommands(config.DiscordSession)

	fmt.Println("[DISCORD] Bot is now running.")
	logger.Discord.Info("Bot is now running.")
	SendMessageToStatusChannel("🤖 Bot Version " + config.Version + " Branch " + config.Branch + " connected to Discord.")
	sendControlPanel() // Send control panel message to Discord
	UpdateBotStatusWithMessage("StationeersServerUI v" + config.Version)
	// Start buffer flush ticker
	config.BufferFlushTicker = time.NewTicker(5 * time.Second)
	go func() {
		for range config.BufferFlushTicker.C {
			flushLogBufferToDiscord()
		}
	}()

	select {} // Keep it running
}

// Updates the bot status with a string message (unused in 4.3)
func UpdateBotStatusWithMessage(message string) {
	err := config.DiscordSession.UpdateGameStatus(0, message)
	if err != nil {
		fmt.Println("Error updating bot status:", err)
	}
}
