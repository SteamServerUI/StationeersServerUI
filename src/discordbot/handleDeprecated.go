package discordbot

import (
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"

	"github.com/bwmarrin/discordgo"
)

// As of v4.5, the ![command] commands are deprecated and will be removed in a future version.

func handleHelpCommand() {
	helpMessage := `
**Available Commands:**
- ` + "`Attention`" + `: The ![command] commands are deprecated and will be removed in a future version. Please use the slash commands instead.
- ` + "`/help`" + `: Displays the help message for the slash commands.`

	SendMessageToControlChannel(helpMessage)
}

// DEPRECATED
func handleUpdateCommand() {
	SendMessageToControlChannel("🙏 Sorry, this feature has been deprecated. Server Updates are now handled automatically at Software Startup. If you are interested in bringing this feature back, please report it on the GitHub repository. We will be happy to implement it.")
}

// DEPRECATED
func handleValidateCommand() {
	SendMessageToControlChannel(" 🙏Sorry, this feature has been deprecated. Server File Validation is now handled automatically at Software Startup. If you are interested in bringing this feature back, please report it on the GitHub repository. We will be happy to implement it.")
}

// DEPRECATED
func handleListCommand() {
	SendMessageToControlChannel("🙏 Sorry, this feature has been deprecated. The /list command has taken over this functionality. Please use slash (/) commands instead.")
}

// DEPRECATED
func handleRestoreCommand() {
	SendMessageToControlChannel("🙏 Sorry, this feature has been deprecated. The /restore command has taken over this functionality. Please use slash (/) commands instead.")
}

// DEPRECATED
func handleBanCommand() {
	SendMessageToControlChannel("🙏 Sorry, this feature has been deprecated. The /bansteamid command has taken over this functionality. Please use slash (/) commands instead.")
}

// DEPRECATED
func handleUnbanCommand() {
	SendMessageToControlChannel("🙏 Sorry, this feature has been deprecated. The /unbansteamid command has taken over this functionality. Please use slash (/) commands instead.")
}

// DEPRECATED
func listenToDiscordMessages(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || m.ChannelID != config.ControlChannelID {
		logger.Discord.Debug("Ignoring message from " + m.Author.Username)
		logger.Discord.Debug("Ignored message: " + m.Content)
		logger.Discord.Debug("Message channel: " + m.ChannelID)

		return
	}

	// log the message if debug is enabled
	logger.Discord.Debug("Ignoring message from " + m.Author.Username)
	logger.Discord.Debug("Ignored message: " + m.Content)
	logger.Discord.Debug("Message channel: " + m.ChannelID)

	content := strings.TrimSpace(m.Content)

	switch {
	case strings.HasPrefix(content, "!start"):
		gamemgr.InternalStartServer()

	case strings.HasPrefix(content, "!stop"):
		gamemgr.InternalStopServer()

	case strings.HasPrefix(content, "!restore"):
		handleRestoreCommand()

	case strings.HasPrefix(content, "!list"):
		handleListCommand()

	case strings.HasPrefix(content, "!update"):
		handleUpdateCommand()

	case strings.HasPrefix(content, "!help"):
		handleHelpCommand()

	case strings.HasPrefix(content, "!ban"):
		handleBanCommand()

	case strings.HasPrefix(content, "!unban"):
		handleUnbanCommand()

	case strings.HasPrefix(content, "!validate"):
		handleValidateCommand()
	default:
		// Optionally handle unrecognized commands or ignore them
	}
}
