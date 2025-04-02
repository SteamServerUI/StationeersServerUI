package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func handleHelpCommand(s *discordgo.Session, channelID string) {
	helpMessage := `
**Available Commands:**
- ` + "`!start`" + `: Starts the server.
- ` + "`!stop`" + `: Stops the server.
- ` + "`!restore:<index>`" + `: Restores a backup at the specified index. Usage: ` + "`!restore:1`" + `.
- ` + "`!list:<number/all>`" + `: Lists the most recent backups. Use ` + "`!list:all`" + ` to list all backups or ` + "`!list:<number>`" + ` to specify how many to list.
- ` + "`!ban:<SteamID>`" + `: Bans a player by their SteamID. Usage: ` + "`!ban:76561198334231312`" + `.
- ` + "`!unban:<SteamID>`" + `: Unbans a player by their SteamID. Usage: ` + "`!unban:76561198334231312`" + `.
- ` + "`!update`" + `: Updates the server files if there is a game update available. (Currently Stable Branch only)
- ` + "`!validate`" + `: Validates the server files if there is a game update available. (Currently Stable Branch only)
- ` + "`!help`" + `: Displays this help message.

Please stop the server before using update commands.
	`

	_, err := s.ChannelMessageSend(channelID, helpMessage)
	if err != nil {
		fmt.Println("Error sending help message:", err)
		SendMessageToControlChannel("Error sending help message.")
	} else {
		fmt.Println("Help message sent to control channel.")
	}
}
