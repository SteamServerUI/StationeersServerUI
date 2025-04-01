package discord

import (
	"StationeersServerUI/src/config"
	"fmt"
	"strings"
)

// v4 OK
func parseBackupList(rawData string) string {
	lines := strings.Split(rawData, "\n")
	var formattedLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ", ")
		if len(parts) == 2 {
			formattedLines = append(formattedLines, fmt.Sprintf("**%s** - %s", parts[0], parts[1]))
		}
	}
	return strings.Join(formattedLines, "\n")
}

// v4 OK
func formatConnectedPlayers() string {
	if len(config.ConnectedPlayers) == 0 {
		return "No players are currently connected."
	}

	var sb strings.Builder
	sb.WriteString("**Connected Players:**\n")
	sb.WriteString("```\n")
	sb.WriteString("Username              | Steam ID\n")
	sb.WriteString("----------------------|------------------------\n")

	for steamID, username := range config.ConnectedPlayers {
		sb.WriteString(fmt.Sprintf("%-20s | %s\n", username, steamID))
	}

	sb.WriteString("```")
	return sb.String()
}
