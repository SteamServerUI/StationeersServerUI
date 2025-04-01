package discord

// v4 FIXED

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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

/*
As Backups.go will be rewritten, these actions will be rebuilt later:
*/

// v4 IGNORED
func handleListCommand(s *discordgo.Session, channelID string, content string) {
	s.ChannelMessageSend(channelID, "❌This feature has been soft-deprecated due to backend changes. It will come back soon, but for now we recommend using the WebUI.")
	return
	fmt.Println("!list command received, fetching backup list...")

	// Extract the "top" number or "all" option from the command
	parts := strings.Split(content, ":")
	top := 5 // Default to 5
	var err error
	if len(parts) == 2 {
		if parts[1] == "all" {
			top = -1 // No limit
		} else {
			top, err = strconv.Atoi(parts[1])
			if err != nil || top < 1 {
				s.ChannelMessageSend(channelID, "❌Invalid number provided. Use `!list:<number>` or `!list:all`.")
				return
			}
		}
	}

	// Step 1: Fetch the backup list from the server
	resp, err := http.Get("http://localhost:8080/backups")
	if err != nil {
		fmt.Println("Failed to fetch backup list:", err)
		s.ChannelMessageSend(channelID, "❌Failed to fetch backup list.")
		return
	}
	defer resp.Body.Close()

	// Step 2: Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read backup list response:", err)
		s.ChannelMessageSend(channelID, "❌Failed to read backup list.")
		return
	}

	// Step 3: Output the raw backup list data for debugging
	//fmt.Println("Raw backup list data:", string(body))

	// Step 4: Parse the backup list data into a formatted string
	backupList := parseBackupList(string(body))
	//fmt.Println("Formatted backup list:\n", backupList)

	// Step 5: Split the backup list into individual lines
	lines := strings.Split(backupList, "\n")

	// Step 6: Send each line as a separate message, respecting the "top" limit
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue // Skip empty lines
		}
		if top > 0 && count >= top {
			break // Stop if we've reached the "top" limit
		}
		fmt.Println("Sending line to Discord:", line)
		message, err := s.ChannelMessageSend(channelID, line)
		if err != nil {
			fmt.Println("Error sending line to Discord:", err)
		} else {
			fmt.Println("Successfully sent line to Discord. Message ID:", message.ID)
		}
		count++

		// Optional: Add a small delay to avoid hitting rate limits
		time.Sleep(500 * time.Millisecond)
	}
}

// v4 IGNORED
func handleRestoreCommand(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	s.ChannelMessageSend(m.ChannelID, "❌This feature has been soft-deprecated due to backend changes. It will come back soon, but for now we recommend using the WebUI.")
	return
	parts := strings.Split(content, ":")
	if len(parts) != 2 {
		s.ChannelMessageSend(m.ChannelID, "❌Invalid restore command. Use `!restore:<index>`.")
		sendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
		return
	}
	// SendCommandToAPI("/stop")
	indexStr := parts[1]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌Invalid index provided for restore.")
		sendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/restore?index=%d", index)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌Failed to restore backup at index %d.", index))
		sendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅Backup %d restored successfully, Starting Server...", index))
	//sleep 5 sec to give the server time to start
	time.Sleep(5 * time.Second)
	//SendCommandToAPI("/start")
}

// v4 IGNORED
func handleUpdateCommand(s *discordgo.Session, channelID string) {
	// Notify that the update process is starting
	s.ChannelMessageSend(channelID, "🙏Sorry, this feature has been deprecated. Server Updates are now handled automatically at Software Startup. If you are interested in bringing this feature back, please report it on the GitHub repository. We will be happy to implement it.")
}

// v4 IGNORED
func handleValidateCommand(s *discordgo.Session, channelID string) {
	// Notify that the update process is starting
	s.ChannelMessageSend(channelID, "🙏Sorry, this feature has been deprecated. Server File Validation is now handled automatically at Software Startup. If you are interested in bringing this feature back, please report it on the GitHub repository. We will be happy to implement it.")
}
