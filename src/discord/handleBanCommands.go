package discord

import (
	"StationeersServerUI/src/config"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var blacklistMutex sync.Mutex

func handleBanCommand(s *discordgo.Session, channelID string, content string) {
	// Extract the SteamID from the command
	parts := strings.Split(content, ":")
	if len(parts) != 2 {
		s.ChannelMessageSend(channelID, "❌Invalid ban command. Use `!ban:<SteamID>`.")
		return
	}
	steamID := strings.TrimSpace(parts[1])

	// Read the current blacklist
	blacklist, err := readBlacklist(config.BlackListFilePath)
	if err != nil {
		s.ChannelMessageSend(channelID, "❌Error reading blacklist file.")
		return
	}

	// Check if the SteamID is already in the blacklist
	entries := strings.Split(blacklist, ",")
	exists := false
	for _, entry := range entries {
		if strings.TrimSpace(entry) == steamID {
			exists = true
			break
		}
	}

	if exists {
		s.ChannelMessageSend(channelID, fmt.Sprintf("❌SteamID %s is already banned.", steamID))
		return
	}

	// Add the SteamID to the blacklist
	blacklist = appendToBlacklist(blacklist, steamID)

	// Write the updated blacklist back to the file
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()
	err = os.WriteFile(config.BlackListFilePath, []byte(blacklist), 0644)
	if err != nil {
		s.ChannelMessageSend(channelID, "❌Error writing to blacklist file.")
		return
	}

	s.ChannelMessageSend(channelID, fmt.Sprintf("✅SteamID %s has been banned.", steamID))
}

func handleUnbanCommand(s *discordgo.Session, channelID string, content string) {
	// Extract the SteamID from the command
	parts := strings.Split(content, ":")
	if len(parts) != 2 {
		s.ChannelMessageSend(channelID, "❌Invalid unban command. Use `!unban:<SteamID>`.")
		return
	}
	steamID := strings.TrimSpace(parts[1])

	// Read the current blacklist
	blacklist, err := readBlacklist(config.BlackListFilePath)
	if err != nil {
		s.ChannelMessageSend(channelID, "❌Error reading blacklist file.")
		return
	}

	// Check if the SteamID is in the blacklist
	entries := strings.Split(blacklist, ",")
	exists := false
	for _, entry := range entries {
		if strings.TrimSpace(entry) == steamID {
			exists = true
			break
		}
	}

	if !exists {
		s.ChannelMessageSend(channelID, fmt.Sprintf("✅SteamID %s is not banned.", steamID))
		return
	}

	// Remove the SteamID from the blacklist
	updatedBlacklist := removeFromBlacklist(blacklist, steamID)

	// Write the updated blacklist back to the file
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()
	err = os.WriteFile(config.BlackListFilePath, []byte(updatedBlacklist), 0644)
	if err != nil {
		s.ChannelMessageSend(channelID, "❌Error writing to blacklist file.")
		return
	}

	s.ChannelMessageSend(channelID, fmt.Sprintf("✅SteamID %s has been unbanned.", steamID))
}

func readBlacklist(blackListFilePath string) (string, error) {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	file, err := os.Open(blackListFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close() // Ensure the file is closed after reading

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func appendToBlacklist(blacklist, steamID string) string {
	if len(blacklist) > 0 && !strings.HasSuffix(blacklist, ",") {
		blacklist += ","
	}
	return blacklist + steamID
}

func removeFromBlacklist(blacklist, steamID string) string {
	entries := strings.Split(blacklist, ",")
	var updatedEntries []string
	for _, entry := range entries {
		if strings.TrimSpace(entry) != steamID {
			updatedEntries = append(updatedEntries, entry)
		}
	}
	return strings.Join(updatedEntries, ",")
}
