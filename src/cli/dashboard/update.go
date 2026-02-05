package dashboard

import (
	"runtime"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/detectionmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model - handles all messages and key presses
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key presses
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Tab):
			// Cycle through panels (full view switch)
			m.activePanel = (m.activePanel + 1) % PanelCount

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Up):
			if m.activePanel == PanelSSUILog {
				m.logViewport.LineUp(3)
			}

		case key.Matches(msg, m.keys.Down):
			if m.activePanel == PanelSSUILog {
				m.logViewport.LineDown(3)
			}

		case key.Matches(msg, m.keys.Start):
			if !m.serverRunning {
				return m, startServerCmd()
			}

		case key.Matches(msg, m.keys.Stop):
			if m.serverRunning {
				return m, stopServerCmd()
			}

		case key.Matches(msg, m.keys.Refresh):
			// Force refresh of all data
			return m, tea.Batch(fetchStatusCmd(), fetchLogsCmd())
		}

	case tea.WindowSizeMsg:
		// Window was resized
		m.width = msg.Width
		m.height = msg.Height

		// Calculate content area height (same for all panels)
		headerHeight := 5 // Title + tabs
		helpHeight := 2   // Help line
		footerHeight := 2 // Footer
		panelPadding := 6 // Border + padding

		contentHeight := m.height - headerHeight - helpHeight - footerHeight - panelPadding
		if contentHeight < 5 {
			contentHeight = 5
		}
		contentWidth := m.width - 8 // Account for panel borders and padding
		if contentWidth < 20 {
			contentWidth = 20
		}

		// Update viewport
		m.logViewport.Width = contentWidth
		m.logViewport.Height = contentHeight
		m.help.Width = m.width

	case tickMsg:
		// Periodic update - refresh data every second
		cmds = append(cmds, tickCmd()) // Schedule next tick
		cmds = append(cmds, fetchStatusCmd())
		cmds = append(cmds, fetchLogsCmd())

	case statusUpdateMsg:
		// Update status from fetched data
		m.serverRunning = msg.running
		m.serverUptime = msg.uptime
		m.serverName = msg.serverName
		m.saveName = msg.saveName
		m.worldID = msg.worldID
		m.gamePort = msg.gamePort
		m.maxPlayers = msg.maxPlayers
		m.connectedPlayers = msg.players
		m.ssuiVersion = msg.version
		m.goRuntime = msg.goRuntime

	case logUpdateMsg:
		// Update log viewport content
		if len(msg.logs) > 0 {
			content := strings.Join(msg.logs, "")
			m.logViewport.SetContent(content)
			m.logViewport.GotoBottom()
		}

	case serverActionMsg:
		// Handle server action results
		if msg.err != nil {
			m.err = msg.err
		}
		// Refresh status after action
		cmds = append(cmds, fetchStatusCmd())
	}

	return m, tea.Batch(cmds...)
}

// fetchStatusCmd returns a command that fetches current server status from real sources
func fetchStatusCmd() tea.Cmd {
	return func() tea.Msg {
		// Get real data from config and managers
		running := config.GetIsGameServerRunning()
		uptime := gamemgr.GetServerUptime()

		// Get connected players
		detector := detectionmgr.GetDetector()
		players := detectionmgr.GetPlayers(detector)

		return statusUpdateMsg{
			running:    running,
			uptime:     uptime,
			serverName: config.GetServerName(),
			saveName:   config.GetSaveName(),
			worldID:    config.GetWorldID(),
			gamePort:   config.GetGamePort(),
			maxPlayers: parseMaxPlayers(config.GetServerMaxPlayers()),
			players:    players,
			version:    config.GetVersion(),
			goRuntime:  runtime.GOOS + "/" + runtime.GOARCH,
		}
	}
}

// fetchLogsCmd returns a command that fetches SSUI log entries
func fetchLogsCmd() tea.Cmd {
	return func() tea.Msg {
		logs := GetLogBuffer()

		if len(logs) == 0 {
			logs = []string{
				"Waiting for SSUI log entries...\n",
				"\n",
				"Logs will appear here when backend activity occurs.\n",
			}
		}

		return logUpdateMsg{logs: logs}
	}
}

// startServerCmd starts the game server
func startServerCmd() tea.Cmd {
	return func() tea.Msg {
		err := gamemgr.InternalStartServer()
		return serverActionMsg{action: "start", err: err}
	}
}

// stopServerCmd stops the game server
func stopServerCmd() tea.Cmd {
	return func() tea.Msg {
		err := gamemgr.InternalStopServer()
		return serverActionMsg{action: "stop", err: err}
	}
}
