package discordbot

import (
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"

	"github.com/bwmarrin/discordgo"
)

// InitializeDiscordBot starts or restarts the Discord bot and connects it to the Discord API.
func InitializeDiscordBot() {
	var err error

	// Clean up previous session
	if config.GetDiscordSession() != nil {
		logger.Discord.Debug("Previous Discord session found, closing it...")
		config.GetDiscordSession().Close()
	}
	if BufferFlushTicker != nil {
		BufferFlushTicker.Stop()
	}

	// Create new session
	session, err := discordgo.New("Bot " + config.GetDiscordToken())
	if err != nil {
		logger.Discord.Error("Error creating Discord session: " + err.Error())
		return
	}
	config.SetDiscordSession(session)

	// Set intents
	config.GetDiscordSession().Identify.Intents = discordgo.IntentsGuildMessageReactions

	logger.Discord.Info("Starting Discord integration...")
	//logger.Discord.Debug("Discord token: " + config.GetDiscordToken())
	logger.Discord.Debug("ControlChannelID: " + config.GetControlChannelID())
	logger.Discord.Debug("StatusChannelID: " + config.GetStatusChannelID())
	logger.Discord.Debug("ConnectionListChannelID: " + config.GetConnectionListChannelID())
	logger.Discord.Debug("LogChannelID: " + config.GetLogChannelID())
	logger.Discord.Debug("SaveChannelID: " + config.GetSaveChannelID())

	// Open session first
	err = config.GetDiscordSession().Open()
	if err != nil {
		logger.Discord.Error("Error opening Discord connection: " + err.Error())
		return
	}

	// Register handlers and commands after session is open
	config.GetDiscordSession().AddHandler(listenToDiscordReactions)
	config.GetDiscordSession().AddHandler(listenToSlashCommands)
	registerSlashCommands(config.GetDiscordSession())

	logger.Discord.Info("Bot is now running.")
	SendMessageToStatusChannel("🤖 SSUI Version " + config.GetVersion() + " connected to Discord.")
	sendControlPanel() // Send control panel message to Discord
	UpdateBotStatusWithMessage("StationeersServerUI v" + config.GetVersion())
	// Start buffer flush ticker
	BufferFlushTicker = time.NewTicker(5 * time.Second)
	go func() {
		for range BufferFlushTicker.C {
			flushLogBufferToDiscord()
		}
	}()

	select {} // Keep it running
}

// Updates the bot status with a string message
func UpdateBotStatusWithMessage(message string) {
	err := config.GetDiscordSession().UpdateGameStatus(0, message)
	if err != nil {
		logger.Discord.Error("Error updating bot status: " + err.Error())
	}
}
