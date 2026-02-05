package dashboard

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// SSUI ASCII Art Logo

const ssuiLogo = `███████╗███████╗██╗   ██╗██╗
██╔════╝██╔════╝██║   ██║██║
███████╗███████╗██║   ██║██║
╚════██║╚════██║██║   ██║██║
███████║███████║╚██████╔╝██║
╚══════╝╚══════╝ ╚═════╝ ╚═╝`

// Main View

// View implements tea.Model - renders the entire UI
func (m Model) View() string {
	if m.quitting {
		farewell := lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true).
			Render("\n  👋 Dashboard closed. Returning to CLI...\n\n")
		return farewell
	}

	// Build the UI
	var b strings.Builder

	// Header with logo and status bar
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	// Tab bar
	b.WriteString(m.renderTabBar())
	b.WriteString("\n\n")

	// Main content panel
	switch m.activePanel {
	case PanelStatus:
		b.WriteString(m.renderStatusPanel())
	case PanelSSUILog:
		b.WriteString(m.renderLogPanel())
	case PanelPlayers:
		b.WriteString(m.renderPlayersPanel())
	case PanelConfig:
		b.WriteString(m.renderConfigPanel())
	}

	b.WriteString("\n")

	// Footer with all controls
	b.WriteString(m.renderFooter())

	return b.String()
}

// Header Components

// renderHeader renders the header with status indicators
func (m Model) renderHeader() string {
	// Status indicators
	var statusItems []string

	// Server status pill
	if m.serverRunning {
		statusItems = append(statusItems, StatusPillOnline.Render(" "+BulletFilled+" ONLINE "))
	} else {
		statusItems = append(statusItems, StatusPillOffline.Render(" "+BulletEmpty+" OFFLINE "))
	}

	// Player count with mini progress bar
	playerCount := len(m.connectedPlayers)
	playerBar := RenderMiniBar(playerCount, m.maxPlayers, 6)
	playerText := fmt.Sprintf("%d/%d", playerCount, m.maxPlayers)
	if playerCount > 0 {
		playerStyle := lipgloss.NewStyle().Foreground(Cyan).Bold(true)
		statusItems = append(statusItems, playerStyle.Render("👥 "+playerBar+" "+playerText))
	} else {
		statusItems = append(statusItems, MutedStyle.Render("👥 "+playerBar+" "+playerText))
	}

	// Uptime (if running)
	if m.serverRunning && m.serverUptime > 0 {
		uptimeStyle := lipgloss.NewStyle().Foreground(Green)
		statusItems = append(statusItems, uptimeStyle.Render("⏱ "+formatUptimeShort(m.serverUptime)))
	}

	// Version badge
	statusItems = append(statusItems, VersionBadgeStyle.Render("v"+m.ssuiVersion))

	statusBar := lipgloss.JoinHorizontal(lipgloss.Center, strings.Join(statusItems, "  "))

	// Simple centered status bar
	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(statusBar)
}

// renderTabBar renders the navigation tab bar
func (m Model) renderTabBar() string {
	var tabs []string

	for i := Panel(0); i < PanelCount; i++ {
		icon := panelIcons[i]
		name := panelNames[i]
		label := icon + " " + name

		if i == m.activePanel {
			tabs = append(tabs, TabActiveStyle.Render(label))
		} else {
			tabs = append(tabs, TabInactiveStyle.Render(label))
		}
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	// Add subtle divider line under tabs
	dividerWidth := lipgloss.Width(tabBar)
	if m.width > dividerWidth {
		dividerWidth = m.width - 4
	}
	divider := DividerStyle.Render(strings.Repeat("─", dividerWidth))

	return lipgloss.JoinVertical(lipgloss.Left, tabBar, divider)
}

// Status Panel

func (m Model) renderStatusPanel() string {
	// Left column: Server info
	leftCol := m.renderServerInfoBox()

	// Right column: World info + Quick stats
	rightCol := lipgloss.JoinVertical(lipgloss.Left,
		m.renderWorldInfoBox(),
		"",
		m.renderQuickStatsBox(),
	)

	// Combine columns
	leftWidth := m.width/2 - 4
	rightWidth := m.width/2 - 4
	if leftWidth < 30 {
		leftWidth = 30
	}
	if rightWidth < 30 {
		rightWidth = 30
	}

	leftColStyled := lipgloss.NewStyle().Width(leftWidth).Render(leftCol)
	rightColStyled := lipgloss.NewStyle().Width(rightWidth).Render(rightCol)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftColStyled, "  ", rightColStyled)

	return m.renderPanelContainer(content)
}

// renderLogoBlock renders the SSUI ASCII logo with gradient
func (m Model) renderLogoBlock() string {
	logoLines := strings.Split(ssuiLogo, "\n")
	var styledLogo strings.Builder

	// Apply gradient colors to logo lines
	gradientColors := []lipgloss.Color{PurpleLight, Purple, Purple, PurpleDark, PurpleDark, Gray500}
	for i, line := range logoLines {
		color := gradientColors[i%len(gradientColors)]
		styledLogo.WriteString(lipgloss.NewStyle().Foreground(color).Bold(true).Render(line))
		if i < len(logoLines)-1 {
			styledLogo.WriteString("\n")
		}
	}
	return styledLogo.String()
}

func (m Model) renderServerInfoBox() string {
	var lines []string

	lines = append(lines, RenderSectionTitle("Server Status"))
	lines = append(lines, "")

	// Status line (simple, no animation)
	var statusLine string
	if m.serverRunning {
		statusLine = lipgloss.NewStyle().Foreground(Green).Bold(true).Render(BulletFilled+" Online") +
			MutedStyle.Render(" • ") +
			lipgloss.NewStyle().Foreground(Cyan).Render("⏱ "+formatUptime(m.serverUptime))
	} else {
		statusLine = OfflineStyle.Render(BulletEmpty+" Offline") +
			MutedStyle.Render(" • press ") +
			KeyStyle.Render("s") +
			MutedStyle.Render(" to start")
	}
	lines = append(lines, statusLine)
	lines = append(lines, "")

	// Server details with icons
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(PurpleLight).Render("🏷")+"  "+RenderKeyValue("Server Name", m.serverName))
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Cyan).Render("🌐")+"  "+RenderKeyValue("Game Port", m.gamePort))
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Blue).Render("📡")+"  "+RenderKeyValue("Update Port", m.updatePort))
	lines = append(lines, "")

	// Player capacity with fancy progress bar
	playerCount := len(m.connectedPlayers)
	barWidth := 20
	playerBar := RenderProgressBar(playerCount, m.maxPlayers, barWidth)
	playerText := fmt.Sprintf("%d/%d", playerCount, m.maxPlayers)

	// Add player icon and styling
	playerIcon := "👥"
	if playerCount == 0 {
		playerIcon = MutedStyle.Render("👥")
	}

	lines = append(lines, "  "+playerIcon+"  "+LabelStyle.Render("Players:"))
	lines = append(lines, "      "+playerBar+" "+NumberStyle.Render(playerText))
	lines = append(lines, "")

	// Add SSUI Logo underneath player count
	lines = append(lines, m.renderLogoBlock())

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m Model) renderWorldInfoBox() string {
	var lines []string

	lines = append(lines, RenderSectionTitle("World Info"))
	lines = append(lines, "")

	// World details with icons
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Green).Render("💾")+"  "+RenderKeyValue("Save Name", m.saveName))
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Cyan).Render("🔑")+"  "+RenderKeyValue("World ID", m.worldID))
	lines = append(lines, "")

	// Game version with special formatting
	versionDisplay := m.gameVersion
	if versionDisplay == "" {
		versionDisplay = MutedStyle.Render("Unknown")
	} else {
		versionDisplay = ValueHighlightStyle.Render(versionDisplay)
	}
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Yellow).Render("📦")+"  "+LabelStyle.Render("Game Version:")+" "+versionDisplay)

	// Branch
	branchDisplay := m.gameBranch
	if branchDisplay == "" || branchDisplay == "public" {
		branchDisplay = lipgloss.NewStyle().Foreground(Green).Render("public")
	} else {
		branchDisplay = lipgloss.NewStyle().Foreground(Yellow).Render(branchDisplay)
	}
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(PurpleLight).Render("🌿")+"  "+LabelStyle.Render("Branch:")+" "+branchDisplay)

	// Build ID if available
	if m.buildID != "" {
		lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Gray500).Render("🔢")+"  "+RenderKeyValue("Build ID", m.buildID))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m Model) renderQuickStatsBox() string {
	var lines []string

	lines = append(lines, RenderSectionTitle("SSUI Info"))
	lines = append(lines, "")

	// Version with branding
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Purple).Bold(true).Render(Rocket)+"  "+LabelStyle.Render("Version:")+" "+ValueHighlightStyle.Render(m.ssuiVersion))

	// Runtime info
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Blue).Render("⚡")+"  "+RenderKeyValue("Runtime", m.goRuntime))

	// Environment with icon
	if m.isDocker {
		lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Cyan).Render("🐳")+"  "+LabelStyle.Render("Environment:")+" "+ValueHighlightStyle.Render("Docker"))
	} else {
		lines = append(lines, "  "+lipgloss.NewStyle().Foreground(Green).Render("💻")+"  "+LabelStyle.Render("Environment:")+" "+ValueStyle.Render("Native"))
	}

	// Session info with animated indicator
	sessionUptime := time.Since(m.startTime).Round(time.Second)
	lines = append(lines, "")
	sessionIcon := GetSpinnerDot(m.tickCount)
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(PurpleLight).Render(sessionIcon)+"  "+LabelStyle.Render("Session:")+" "+ValueStyle.Render(formatUptime(sessionUptime)))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// Players Panel

func (m Model) renderPlayersPanel() string {
	var content string

	if len(m.connectedPlayers) == 0 {
		// Empty state with nice styling
		emptyIcon := lipgloss.NewStyle().
			Foreground(Gray500).
			Render("👥")

		emptyTitle := lipgloss.NewStyle().
			Foreground(Gray400).
			Bold(true).
			Render("No Players Connected")

		var emptySubtext string
		if !m.serverRunning {
			emptySubtext = MutedStyle.Render("Server is offline. Press ") +
				KeyStyle.Render("s") +
				MutedStyle.Render(" to start.")
		} else {
			// Show animated waiting indicator
			spinner := GetSpinnerFrame(m.tickCount)
			emptySubtext = lipgloss.NewStyle().Foreground(Cyan).Render(spinner) +
				MutedStyle.Render(" Waiting for players to join...")
		}

		content = lipgloss.JoinVertical(lipgloss.Center,
			"",
			"",
			emptyIcon,
			"",
			emptyTitle,
			"",
			emptySubtext,
			"",
			"",
		)

		// Center the content
		content = lipgloss.NewStyle().
			Width(m.width - 8).
			Align(lipgloss.Center).
			Render(content)
	} else {
		// Player list header with player bar
		playerBar := RenderProgressBar(len(m.connectedPlayers), m.maxPlayers, 15)
		headerLine := lipgloss.JoinHorizontal(lipgloss.Top,
			RenderSectionTitle("Connected Players"),
			"  ",
			playerBar,
			" ",
			NumberStyle.Render(fmt.Sprintf("%d/%d", len(m.connectedPlayers), m.maxPlayers)),
		)

		// Column headers with better styling
		colHeader := lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(4).Foreground(Gray600).Render("#"),
			lipgloss.NewStyle().Width(26).Foreground(Gray500).Bold(true).Render("Player Name"),
			lipgloss.NewStyle().Foreground(Gray500).Bold(true).Render("Steam ID"),
		)

		// Fancy divider
		divider := DividerStyle.Render(BoxTeeRight + strings.Repeat(BoxHorizontal, 60) + BoxTeeLeft)

		// Player rows
		var rows []string
		rows = append(rows, headerLine)
		rows = append(rows, "")
		rows = append(rows, colHeader)
		rows = append(rows, divider)

		// Sort players for consistent display
		var steamIDs []string
		for steamID := range m.connectedPlayers {
			steamIDs = append(steamIDs, steamID)
		}
		sort.Strings(steamIDs)

		for i, steamID := range steamIDs {
			playerName := m.connectedPlayers[steamID]

			// Alternate row styling for readability
			indexStyle := lipgloss.NewStyle().Width(4).Foreground(Gray500)
			nameStyle := lipgloss.NewStyle().Width(26).Foreground(GreenLight).Bold(true)
			idStyle := lipgloss.NewStyle().Foreground(Gray500).Italic(true)

			// Add online indicator
			onlineIndicator := lipgloss.NewStyle().Foreground(Green).Render(BulletFilled)

			row := lipgloss.JoinHorizontal(lipgloss.Top,
				indexStyle.Render(fmt.Sprintf("%d.", i+1)),
				onlineIndicator+" "+nameStyle.Render(playerName),
				idStyle.Render("  "+steamID),
			)
			rows = append(rows, row)
		}

		content = lipgloss.JoinVertical(lipgloss.Left, rows...)
	}

	return m.renderPanelContainer(content)
}

// Log Panel

func (m Model) renderLogPanel() string {
	// Header with scroll indicator
	scrollPercent := int(m.logViewport.ScrollPercent() * 100)
	scrollIndicator := LogScrollIndicatorStyle.Render(fmt.Sprintf("(%d%%)", scrollPercent))

	// Scroll position indicator
	var scrollHint string
	if scrollPercent < 100 {
		scrollHint = MutedStyle.Render(" • Use ↑↓ to scroll, 'G' for bottom")
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		RenderSectionTitle("SSUI Logs"),
		"  ",
		scrollIndicator,
		scrollHint,
	)

	// Log content
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		m.logViewport.View(),
	)

	// Use a slightly different style for log panel
	panelWidth := m.width - 4
	if panelWidth < 40 {
		panelWidth = 40
	}

	return LogContainerStyle.
		Width(panelWidth).
		BorderForeground(Purple).
		Render(content)
}

// Footer

func (m Model) renderFooter() string {
	// Left: Branding
	brand := FooterBrandStyle.Render("SSUI")
	version := MutedStyle.Render(" v" + m.ssuiVersion)
	left := brand + version

	// Center: All key hints in one comprehensive line
	var keyHints []string

	keyHints = append(keyHints,
		KeyStyle.Render("tab")+KeyDescStyle.Render(" panels"),
		KeyStyle.Render("↑↓")+KeyDescStyle.Render(" nav"),
		KeyStyle.Render("s")+KeyDescStyle.Render(" start"),
		KeyStyle.Render("x")+KeyDescStyle.Render(" stop"),
	)

	// Add config-specific hints when on config panel
	if m.activePanel == PanelConfig {
		keyHints = append(keyHints,
			KeyStyle.Render("enter")+KeyDescStyle.Render(" edit"),
			KeyStyle.Render("ctrl+s")+KeyDescStyle.Render(" save"),
		)
	}

	keyHints = append(keyHints,
		KeyStyle.Render("r")+KeyDescStyle.Render(" refresh"),
		KeyStyle.Render("q")+KeyDescStyle.Render(" quit"),
	)

	center := lipgloss.JoinHorizontal(lipgloss.Top,
		strings.Join(keyHints, FooterSeparatorStyle.Render(" │ ")),
	)

	// Calculate spacing - simple centered layout
	leftWidth := lipgloss.Width(left)
	centerWidth := lipgloss.Width(center)
	totalContentWidth := leftWidth + centerWidth

	var footerContent string
	if m.width > totalContentWidth+10 {
		// Full layout with spacing
		spacing := max((m.width-totalContentWidth)/2, 2)
		footerContent = left + strings.Repeat(" ", spacing) + center
	} else {
		// Compact layout
		footerContent = left + "  " + center
	}

	return FooterStyle.Width(m.width).Render(footerContent)
}

// Helper Functions

// renderPanelContainer wraps content in the standard panel style
func (m Model) renderPanelContainer(content string) string {
	// Calculate consistent content height
	headerHeight := 6 // Status bar + tabs
	footerHeight := 2 // Footer
	panelPadding := 4 // Border + padding

	contentHeight := max(m.height-headerHeight-footerHeight-panelPadding, 10)

	panelWidth := m.width - 4
	if panelWidth < 40 {
		panelWidth = 40
	}

	return PanelStyle.
		Width(panelWidth).
		Height(contentHeight).
		Render(content)
}
