package discord

// v4 NOT OK

import (
	"StationeersServerUI/src/config"
	"fmt"
	"net/http"
	"strings"
)

// v4 NOT OK
func SendCommandToAPI(endpoint string) {
	url := "http://localhost:8080" + endpoint
	if _, err := http.Get(url); err != nil {
		fmt.Printf("Failed to send %s command: %v\n", endpoint, err)
	}
}

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
