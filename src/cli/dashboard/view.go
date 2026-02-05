package dashboard

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Color palette - cohesive with SSUI branding
var (
	// Primary colors
	primaryColor   = lipgloss.Color("#7C3AED") // Purple - main brand
	secondaryColor = lipgloss.Color("#10B981") // Green - success/online
	warningColor   = lipgloss.Color("#F59E0B") // Amber - warnings
	errorColor     = lipgloss.Color("#EF4444") // Red - errors/offline
	mutedColor     = lipgloss.Color("#6B7280") // Gray - muted text

	// Background colors
	bgDark   = lipgloss.Color("#1F2937")
	bgDarker = lipgloss.Color("#111827")
)

// Styles
var (
	// Header style
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F9FAFB")).
			Background(bgDark).
			Padding(0, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F9FAFB")).
			Background(primaryColor).
			Padding(0, 2)

	// Panel styles
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Status indicators
	onlineStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	offlineStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Labels and values
	labelStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(18)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9FAFB"))

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Footer
	footerStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(bgDark).
			Padding(0, 2)

	// Tab indicator style
	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	// Player list styles
	playerNameStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	playerIDStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)
)

// View implements tea.Model - renders the entire UI
func (m Model) View() string {
	if m.quitting {
		return "\n  Exiting dashboard...\n\n"
	}

	// Build the UI
	var b strings.Builder

	// Header with tabs
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")

	// Render the active panel (full view) - all panels use same sizing
	switch m.activePanel {
	case PanelSSUILog:
		b.WriteString(m.renderLogPanel())
	case PanelStatus:
		b.WriteString(m.renderStatusPanel())
	case PanelPlayers:
		b.WriteString(m.renderPlayersPanel())
	}
	b.WriteString("\n")

	// Help
	b.WriteString(helpStyle.Render(m.help.View(m.keys)))
	b.WriteString("\n")

	// Footer
	b.WriteString(m.renderFooter())

	return b.String()
}

// renderHeader renders the dashboard header with tab indicators
func (m Model) renderHeader() string {
	title := titleStyle.Render("🚀 SSUI Dashboard")

	// Tab indicators
	var tabs []string
	for i := Panel(0); i < PanelCount; i++ {
		name := panelNames[i]
		if i == m.activePanel {
			tabs = append(tabs, tabActiveStyle.Render("["+name+"]"))
		} else {
			tabs = append(tabs, tabInactiveStyle.Render(" "+name+" "))
		}
	}
	tabBar := strings.Join(tabs, " ")

	// Session uptime
	uptime := time.Since(m.startTime).Round(time.Second)
	uptimeStr := fmt.Sprintf("Session: %s", uptime)

	// Build header
	return headerStyle.Width(m.width).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", mutedStyle(uptimeStr)),
			tabBar,
		),
	)
}

// renderStatusPanel renders the server status panel (full view)
func (m Model) renderStatusPanel() string {
	// Status icon and text
	var statusDisplay string
	if m.serverRunning {
		statusDisplay = onlineStyle.Render("● Online")
	} else {
		statusDisplay = offlineStyle.Render("○ Offline")
	}

	// Format uptime
	uptimeStr := formatUptime(m.serverUptime)

	// Player count
	playerCount := len(m.connectedPlayers)

	// Build status content with sections
	serverSection := lipgloss.JoinVertical(lipgloss.Left,
		boldStyle("═══ Server Status ═══"),
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Status:"), statusDisplay),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Server Name:"), valueStyle.Render(m.serverName)),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Uptime:"), valueStyle.Render(uptimeStr)),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Players:"), valueStyle.Render(fmt.Sprintf("%d / %d", playerCount, m.maxPlayers))),
	)

	worldSection := lipgloss.JoinVertical(lipgloss.Left,
		"",
		boldStyle("═══ World Info ═══"),
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Save Name:"), valueStyle.Render(m.saveName)),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("World ID:"), valueStyle.Render(m.worldID)),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Game Port:"), valueStyle.Render(m.gamePort)),
	)

	ssuiSection := lipgloss.JoinVertical(lipgloss.Left,
		"",
		boldStyle("═══ SSUI Info ═══"),
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Version:"), valueStyle.Render(m.ssuiVersion)),
		lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Runtime:"), valueStyle.Render(m.goRuntime)),
	)

	controlsSection := lipgloss.JoinVertical(lipgloss.Left,
		"",
		boldStyle("═══ Quick Actions ═══"),
		"",
		mutedStyle("  [s] Start Server    [x] Stop Server    [r] Refresh    [tab] Switch View"),
	)

	content := lipgloss.JoinVertical(lipgloss.Left, serverSection, worldSection, ssuiSection, controlsSection)

	// Use consistent panel sizing
	return m.renderPanelWithContent(content)
}

// renderPlayersPanel renders the connected players panel (full view)
func (m Model) renderPlayersPanel() string {
	var content string

	if len(m.connectedPlayers) == 0 {
		content = mutedStyle("No players connected\n\n")
		if !m.serverRunning {
			content += mutedStyle("Server is offline. Press 's' to start.")
		} else {
			content += mutedStyle("Waiting for players to join...")
		}
	} else {
		// Sort player names for consistent display
		var lines []string
		lines = append(lines, boldStyle(fmt.Sprintf("═══ Connected Players (%d/%d) ═══", len(m.connectedPlayers), m.maxPlayers)))
		lines = append(lines, "")

		// Get sorted keys
		var steamIDs []string
		for steamID := range m.connectedPlayers {
			steamIDs = append(steamIDs, steamID)
		}
		sort.Strings(steamIDs)

		for i, steamID := range steamIDs {
			playerName := m.connectedPlayers[steamID]
			line := fmt.Sprintf("  %d. %s  %s",
				i+1,
				playerNameStyle.Render(playerName),
				playerIDStyle.Render("("+steamID+")"),
			)
			lines = append(lines, line)
		}

		content = strings.Join(lines, "\n")
	}

	// Use consistent panel sizing
	return m.renderPanelWithContent(content)
}

// renderLogPanel renders the SSUI log viewer panel (full view)
func (m Model) renderLogPanel() string {
	// Header with scroll indicator
	scrollInfo := fmt.Sprintf("(%d%%)", int(m.logViewport.ScrollPercent()*100))
	header := boldStyle("═══ SSUI Log ═══") + "  " + mutedStyle(scrollInfo)

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		m.logViewport.View(),
	)

	return panelStyle.Width(m.width - 4).Render(content)
}

// renderPanelWithContent renders content in a consistently-sized panel
func (m Model) renderPanelWithContent(content string) string {
	// Calculate consistent content height
	// Account for: header (3 lines), panel borders/padding (4 lines), help (1 line), footer (1 line), newlines (3)
	contentHeight := m.height - 12
	if contentHeight < 5 {
		contentHeight = 5
	}

	return panelStyle.Width(m.width - 4).Height(contentHeight).Render(content)
}

// renderFooter renders the footer with version and status info
func (m Model) renderFooter() string {
	// Left side: SSUI branding with real version
	left := fmt.Sprintf("SSUI %s", m.ssuiVersion)

	// Right side: Quick status indicator
	var serverIndicator string
	if m.serverRunning {
		serverIndicator = onlineStyle.Render("●")
	} else {
		serverIndicator = offlineStyle.Render("○")
	}
	right := fmt.Sprintf("%s Server  |  %d players", serverIndicator, len(m.connectedPlayers))

	// Calculate spacing
	width := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if width < 0 {
		width = 0
	}
	spacer := strings.Repeat(" ", width)

	return footerStyle.Width(m.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, left, spacer, right),
	)
}

// Helper style functions
func boldStyle(s string) string {
	return lipgloss.NewStyle().Bold(true).Render(s)
}

func mutedStyle(s string) string {
	return lipgloss.NewStyle().Foreground(mutedColor).Render(s)
}
