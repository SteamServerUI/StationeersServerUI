package discordbot

import (
	"StationeersServerUI/src/config"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// InitializeDiscordBot starts or restarts the Discord bot and connects it to the Discord API.
func InitializeDiscordBot() {
	var err error

	if config.DiscordSession != nil {
		if config.IsDebugMode {
			fmt.Println("[DISCORD] Previous Discord session found, closing it...")
		}
		config.DiscordSession.Close()
	}

	if config.BufferFlushTicker != nil {
		config.BufferFlushTicker.Stop()
	}

	config.DiscordSession, err = discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		fmt.Println("[DISCORD] Error creating Discord session:", err)
		return
	}

	// Set intents explicitly: Guilds for channel info, GuildMessages for message handling, GuildMessageReactions for reactions
	config.DiscordSession.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsMessageContent

	fmt.Println("[DISCORD] Starting Discord integration...")
	if config.IsDebugMode {
		fmt.Println("[DISCORD] Discord token:", config.DiscordToken)
		fmt.Println("[DISCORD] ControlChannelID:", config.ControlChannelID)
		fmt.Println("[DISCORD] StatusChannelID:", config.StatusChannelID)
		fmt.Println("[DISCORD] ConnectionListChannelID:", config.ConnectionListChannelID)
		fmt.Println("[DISCORD] LogChannelID:", config.LogChannelID)
		fmt.Println("[DISCORD] SaveChannelID:", config.SaveChannelID)
	}

	config.DiscordSession.AddHandler(listenToDiscordMessages)
	config.DiscordSession.AddHandler(listenToDiscordReactions)
	config.DiscordSession.AddHandler(listenToSlashCommands)

	err = config.DiscordSession.Open()
	if err != nil {
		fmt.Println("[DISCORD] Error opening Discord connection:", err)
		return
	}

	registerSlashCommands(config.DiscordSession)

	fmt.Println("[DISCORD] Bot is now running.")
	// Start the buffer flush ticker to send the remaining buffer every 5 seconds
	config.BufferFlushTicker = time.NewTicker(5 * time.Second)
	SendMessageToStatusChannel("🤖 Bot Version " + config.Version + " Branch " + config.Branch + " connected to Discord.")
	go func() {
		for range config.BufferFlushTicker.C {
			flushLogBufferToDiscord()
		}
	}()
	sendControlMessage()
	select {} // Keep the program running
}

// Updates the bot status with a string message (unused in 4.3)
func UpdateBotStatusWithMessage(message string) {
	err := config.DiscordSession.UpdateGameStatus(0, message)
	if err != nil {
		fmt.Println("Error updating bot status:", err)
	}
}
