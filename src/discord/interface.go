package discord

import (
	"StationeersServerUI/src/config"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// connectedPlayersMessageID tracks the message ID for editing the connected players message
var connectedPlayersMessageID string

// StartDiscordBot starts the Discord bot and connects it to the Discord API.
func StartDiscordBot() {
	var err error
	config.DiscordSession, err = discordgo.New("Bot " + config.DiscordToken)
	fmt.Println("[DISCORD] Starting Discord integration...")
	if config.IsDebugMode {
		fmt.Println("[DISCORD] Discord token:", config.DiscordToken)
		fmt.Println("[DISCORD] ControlChannelID:", config.ControlChannelID)
		fmt.Println("[DISCORD] StatusChannelID:", config.StatusChannelID)
		fmt.Println("[DISCORD] ConnectionListChannelID:", config.ConnectionListChannelID)
		fmt.Println("[DISCORD] LogChannelID:", config.LogChannelID)
		fmt.Println("[DISCORD] SaveChannelID:", config.SaveChannelID)
	}
	if err != nil {
		fmt.Println("[DISCORD] Error creating Discord session:", err)
		return
	}
	fmt.Println("[DISCORD] Bot is now running and connected")

	config.DiscordSession.AddHandler(listenToDiscordMessages)
	config.DiscordSession.AddHandler(listenToDiscordReactions)

	err = config.DiscordSession.Open()
	if err != nil {
		fmt.Println("[DISCORD] Error opening Discord connection:", err)
		return
	}

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
