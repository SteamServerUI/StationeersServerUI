package discordbot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func handleHelpCommand(s *discordgo.Session, channelID string) {
	helpMessage := `
**Available Commands:**
- ` + "`Attention`" + `: The ![command] commands are deprecated and will be removed in a future version. Please use the slash commands instead.
- ` + "`/help`" + `: Displays the help message for the slash commands.
- ` + "`!start`" + `: Starts the server.
- ` + "`!stop`" + `: Stops the server.
- ` + "`!list:<number/all>`" + `: Lists the most recent backups. Use ` + "`!list:all`" + ` to list all backups or ` + "`!list:<number>`" + ` to specify how many to list.
- ` + "`!ban:<SteamID>`" + `: Bans a player by their SteamID. Usage: ` + "`!ban:76561198334231312`" + `.
- ` + "`!unban:<SteamID>`" + `: Unbans a player by their SteamID. Usage: ` + "`!unban:76561198334231312`" + `.
- ` + "`!help`" + `: Displays this help message.`

	_, err := s.ChannelMessageSend(channelID, helpMessage)
	if err != nil {
		fmt.Println("Error sending help message:", err)
		SendMessageToControlChannel("Error sending help message.")
	}
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
