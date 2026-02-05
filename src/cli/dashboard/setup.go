// setup.go - Bubble Tea based interactive setup wizard
package dashboard

import (
	"fmt"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/security"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SetupStep represents each step in the setup wizard
type SetupStep int

const (
	StepWelcome SetupStep = iota
	StepGameBranch
	StepServerName
	StepSaveName
	StepWorldType
	StepMaxPlayers
	StepServerPassword
	StepNetworkChoice
	StepGamePort
	StepUpdatePort
	StepUPnP
	StepLocalIP
	StepAdminUsername
	StepAdminPassword
	StepDashboardPref
	StepComplete
)

// SetupModel is the Bubble Tea model for the setup wizard
type SetupModel struct {
	step      SetupStep
	width     int
	height    int
	textInput textinput.Model
	err       error
	quitting  bool

	// Selection state for dropdown-style steps
	selectedIndex int
	options       []Option

	// Collected values
	gameBranch       string
	serverName       string
	saveName         string
	worldType        string
	maxPlayers       string
	serverPassword   string
	configureNetwork bool
	gamePort         string
	updatePort       string
	upnpEnabled      bool
	localIP          string
	adminUsername    string
	adminPassword    string
	dashboardEnabled bool
}

// Option represents a selectable option
type Option struct {
	Display string
	Value   string
}

// Setup styles
var (
	setupTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F9FAFB")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 2).
			MarginBottom(1)

	setupBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2).
			Width(70)

	setupLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true)

	setupDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			MarginBottom(1)

	setupSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10B981")).
				Bold(true)

	setupUnselectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6B7280"))

	setupSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10B981")).
				Bold(true)

	setupHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)

	setupProgressStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7C3AED"))
)

// NewSetupModel creates a new setup wizard model
func NewSetupModel() SetupModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 50

	return SetupModel{
		step:        StepWelcome,
		textInput:   ti,
		gameBranch:  "public",
		serverName:  "My Stationeers Server",
		saveName:    "MyWorld",
		worldType:   getDefaultWorld(),
		maxPlayers:  "10",
		gamePort:    "27016",
		updatePort:  "27015",
		upnpEnabled: true,
	}
}

func getDefaultWorld() string {
	if config.GetIsNewTerrainAndSaveSystem() {
		return "Lunar"
	}
	return "Moon"
}

// Init implements tea.Model
func (m SetupModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "up", "k":
			if len(m.options) > 0 && m.selectedIndex > 0 {
				m.selectedIndex--
			}

		case "down", "j":
			if len(m.options) > 0 && m.selectedIndex < len(m.options)-1 {
				m.selectedIndex++
			}

		case "tab":
			// Toggle for boolean options
			if m.step == StepNetworkChoice || m.step == StepUPnP || m.step == StepDashboardPref {
				if len(m.options) > 0 {
					m.selectedIndex = (m.selectedIndex + 1) % len(m.options)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update text input for text entry steps
	if m.isTextEntryStep() {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m SetupModel) isTextEntryStep() bool {
	switch m.step {
	case StepServerName, StepSaveName, StepMaxPlayers, StepServerPassword,
		StepGamePort, StepUpdatePort, StepLocalIP, StepAdminUsername, StepAdminPassword:
		return true
	}
	return false
}

func (m SetupModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepWelcome:
		m.step = StepGameBranch
		m.setupBranchOptions()

	case StepGameBranch:
		m.gameBranch = m.options[m.selectedIndex].Value
		m.step = StepServerName
		m.textInput.SetValue(m.serverName)
		m.textInput.Focus()

	case StepServerName:
		if m.textInput.Value() != "" {
			m.serverName = m.textInput.Value()
		}
		m.step = StepSaveName
		m.textInput.SetValue(m.saveName)

	case StepSaveName:
		if m.textInput.Value() != "" {
			m.saveName = m.textInput.Value()
		}
		m.step = StepWorldType
		m.setupWorldOptions()

	case StepWorldType:
		m.worldType = m.options[m.selectedIndex].Value
		m.step = StepMaxPlayers
		m.textInput.SetValue(m.maxPlayers)

	case StepMaxPlayers:
		if m.textInput.Value() != "" {
			m.maxPlayers = m.textInput.Value()
		}
		m.step = StepServerPassword
		m.textInput.SetValue("")
		m.textInput.EchoMode = textinput.EchoPassword

	case StepServerPassword:
		m.serverPassword = m.textInput.Value()
		m.textInput.EchoMode = textinput.EchoNormal
		m.step = StepNetworkChoice
		m.setupYesNoOptions(false)

	case StepNetworkChoice:
		m.configureNetwork = m.selectedIndex == 0
		if m.configureNetwork {
			m.step = StepGamePort
			m.textInput.SetValue(m.gamePort)
		} else {
			m.step = StepAdminUsername
			m.textInput.SetValue("admin")
		}

	case StepGamePort:
		if m.textInput.Value() != "" {
			m.gamePort = m.textInput.Value()
		}
		m.step = StepUpdatePort
		m.textInput.SetValue(m.updatePort)

	case StepUpdatePort:
		if m.textInput.Value() != "" {
			m.updatePort = m.textInput.Value()
		}
		m.step = StepUPnP
		m.setupYesNoOptions(true)

	case StepUPnP:
		m.upnpEnabled = m.selectedIndex == 0
		m.step = StepLocalIP
		m.textInput.SetValue("")

	case StepLocalIP:
		m.localIP = m.textInput.Value()
		m.step = StepAdminUsername
		m.textInput.SetValue("admin")

	case StepAdminUsername:
		if m.textInput.Value() != "" {
			m.adminUsername = m.textInput.Value()
		}
		m.step = StepAdminPassword
		m.textInput.SetValue("")
		m.textInput.EchoMode = textinput.EchoPassword

	case StepAdminPassword:
		m.adminPassword = m.textInput.Value()
		m.textInput.EchoMode = textinput.EchoNormal
		m.step = StepDashboardPref
		m.setupYesNoOptions(false)

	case StepDashboardPref:
		m.dashboardEnabled = m.selectedIndex == 0
		m.step = StepComplete
		m.saveConfig()

	case StepComplete:
		return m, tea.Quit
	}

	return m, nil
}

func (m *SetupModel) setupBranchOptions() {
	m.selectedIndex = 0
	m.options = []Option{
		{Display: "Stable (public) - Recommended", Value: "public"},
		{Display: "Beta - Latest features, may be unstable", Value: "beta"},
		{Display: "Pre-terrain rework", Value: "preterrain"},
		{Display: "Pre-rocket refactor", Value: "prerocket"},
		{Display: "Previous version", Value: "previous"},
		{Display: "Multiplayer-safe rollback", Value: "multiplayersafe"},
	}
}

func (m *SetupModel) setupWorldOptions() {
	m.selectedIndex = 0
	if config.GetIsNewTerrainAndSaveSystem() {
		m.options = []Option{
			{Display: "Lunar - Earth's Moon", Value: "Lunar"},
			{Display: "Vulcan - Volcanic world", Value: "Vulcan2"},
			{Display: "Venus - Dense atmosphere", Value: "Venus"},
			{Display: "Mars - Red planet", Value: "Mars2"},
			{Display: "Europa - Ice moon", Value: "Europa3"},
			{Display: "Mimas Herschel - Saturn's moon", Value: "MimasHerschel"},
		}
	} else {
		m.options = []Option{
			{Display: "Moon - Earth's Moon", Value: "Moon"},
			{Display: "Vulcan - Volcanic world", Value: "Vulcan"},
			{Display: "Venus - Dense atmosphere", Value: "Venus"},
			{Display: "Mars - Red planet", Value: "Mars"},
			{Display: "Europa - Ice moon", Value: "Europa"},
			{Display: "Mimas - Saturn's moon", Value: "Mimas"},
		}
	}
}

func (m *SetupModel) setupYesNoOptions(defaultYes bool) {
	if defaultYes {
		m.selectedIndex = 0
	} else {
		m.selectedIndex = 1
	}
	m.options = []Option{
		{Display: "Yes", Value: "yes"},
		{Display: "No", Value: "no"},
	}
}

func (m *SetupModel) saveConfig() {
	// Save all collected values
	config.SetGameBranch(m.gameBranch)
	config.SetServerName(m.serverName)
	config.SetSaveName(m.saveName)
	config.SetWorldID(m.worldType)
	config.SetServerMaxPlayers(m.maxPlayers)
	config.SetServerPassword(m.serverPassword)

	if m.configureNetwork {
		config.SetGamePort(m.gamePort)
		config.SetUpdatePort(m.updatePort)
		config.SetUPNPEnabled(m.upnpEnabled)
		if m.localIP != "" {
			config.SetLocalIpAddress(m.localIP)
		}
	}

	// Set up admin user if provided
	if m.adminUsername != "" && m.adminPassword != "" {
		// Hash the password
		hashedPassword, err := security.HashPassword(m.adminPassword)
		if err == nil {
			// Initialize Users map if nil
			if config.GetUsers() == nil {
				config.SetUsers(make(map[string]string))
			}
			config.SetUsers(map[string]string{m.adminUsername: hashedPassword})
			config.SetAuthEnabled(true)
		}
	}

	config.SetIsCLIDashboardEnabled(m.dashboardEnabled)
	config.SetIsFirstTimeSetup(false)
}

// View implements tea.Model
func (m SetupModel) View() string {
	if m.quitting {
		return "\n  Setup cancelled.\n\n"
	}

	var content string

	switch m.step {
	case StepWelcome:
		content = m.renderWelcome()
	case StepGameBranch:
		content = m.renderSelection("Game Branch", "Select the game version branch to use:", m.options)
	case StepServerName:
		content = m.renderTextInput("Server Name", "This name appears in the server browser.", m.serverName)
	case StepSaveName:
		content = m.renderTextInput("Save Name", "The name of your world save folder.", m.saveName)
	case StepWorldType:
		content = m.renderSelection("World Type", "Select the starting planet/moon:", m.options)
	case StepMaxPlayers:
		content = m.renderTextInput("Max Players", "Maximum number of players (1-30).", m.maxPlayers)
	case StepServerPassword:
		content = m.renderTextInput("Server Password", "Leave empty for public server, or set a password.", "")
	case StepNetworkChoice:
		content = m.renderSelection("Network Configuration", "Configure ports and network? (Most users can skip)", m.options)
	case StepGamePort:
		content = m.renderTextInput("Game Port", "UDP port for game connections.", m.gamePort)
	case StepUpdatePort:
		content = m.renderTextInput("Update Port", "UDP port for Steam queries.", m.updatePort)
	case StepUPnP:
		content = m.renderSelection("Enable UPnP", "Automatically configure router port forwarding?", m.options)
	case StepLocalIP:
		content = m.renderTextInput("Local IP", "Leave empty for auto-detect, or specify your server's LAN IP.", "")
	case StepAdminUsername:
		content = m.renderTextInput("Admin Username", "Username for Web UI login.", "admin")
	case StepAdminPassword:
		content = m.renderTextInput("Admin Password", "Password for Web UI login. Leave empty to skip auth.", "")
	case StepDashboardPref:
		content = m.renderSelection("CLI Dashboard", "Auto-launch the fancy dashboard on startup?", m.options)
	case StepComplete:
		content = m.renderComplete()
	}

	// Progress indicator
	progress := m.renderProgress()

	return lipgloss.JoinVertical(lipgloss.Center, progress, "", content)
}

func (m SetupModel) renderWelcome() string {
	title := setupTitleStyle.Render("🚀 SSUI Setup Wizard")

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		setupLabelStyle.Render("Welcome to Stationeers Server UI!"),
		"",
		"This wizard will help you configure your server.",
		"",
		setupDescStyle.Render("• Use ↑/↓ arrows to navigate options"),
		setupDescStyle.Render("• Press Enter to confirm your selection"),
		setupDescStyle.Render("• Press Esc to cancel setup"),
		"",
		setupHintStyle.Render("Press Enter to begin..."),
	)

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		setupBoxStyle.Render(body),
	)
}

func (m SetupModel) renderSelection(title, description string, options []Option) string {
	titleRendered := setupTitleStyle.Render("🚀 " + title)

	var optionLines []string
	for i, opt := range options {
		cursor := "  "
		style := setupUnselectedStyle
		if i == m.selectedIndex {
			cursor = "▸ "
			style = setupSelectedStyle
		}
		optionLines = append(optionLines, cursor+style.Render(opt.Display))
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		setupDescStyle.Render(description),
		"",
		strings.Join(optionLines, "\n"),
		"",
		setupHintStyle.Render("↑/↓ to select, Enter to confirm"),
	)

	return lipgloss.JoinVertical(lipgloss.Center,
		titleRendered,
		setupBoxStyle.Render(body),
	)
}

func (m SetupModel) renderTextInput(title, description, placeholder string) string {
	titleRendered := setupTitleStyle.Render("🚀 " + title)

	hint := "Enter to confirm"
	if placeholder != "" {
		hint = fmt.Sprintf("Default: %s • Enter to confirm", placeholder)
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		setupDescStyle.Render(description),
		"",
		m.textInput.View(),
		"",
		setupHintStyle.Render(hint),
	)

	return lipgloss.JoinVertical(lipgloss.Center,
		titleRendered,
		setupBoxStyle.Render(body),
	)
}

func (m SetupModel) renderComplete() string {
	title := setupTitleStyle.Render("✅ Setup Complete!")

	summary := lipgloss.JoinVertical(lipgloss.Left,
		"",
		setupSuccessStyle.Render("Your server is configured!"),
		"",
		fmt.Sprintf("  Server:    %s", m.serverName),
		fmt.Sprintf("  World:     %s (%s)", m.worldType, m.saveName),
		fmt.Sprintf("  Branch:    %s", m.gameBranch),
		fmt.Sprintf("  Players:   %s max", m.maxPlayers),
		"",
		setupLabelStyle.Render(fmt.Sprintf("🌐 Web UI: https://localhost:%s", config.GetSSUIWebPort())),
		"",
	)

	if m.dashboardEnabled {
		summary = lipgloss.JoinVertical(lipgloss.Left,
			summary,
			setupDescStyle.Render("Dashboard will auto-launch on next startup."),
		)
	}

	summary = lipgloss.JoinVertical(lipgloss.Left,
		summary,
		"",
		setupHintStyle.Render("Press Enter to exit setup..."),
	)

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		setupBoxStyle.Render(summary),
	)
}

func (m SetupModel) renderProgress() string {
	totalSteps := 10 // Approximate, depends on path
	currentStep := int(m.step)
	if currentStep > totalSteps {
		currentStep = totalSteps
	}

	filled := strings.Repeat("█", currentStep)
	empty := strings.Repeat("░", totalSteps-currentStep)

	return setupProgressStyle.Render(fmt.Sprintf("Progress: [%s%s] Step %d", filled, empty, currentStep))
}

// RunSetup starts the setup wizard
func RunSetup() error {
	m := NewSetupModel()

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("setup wizard error: %w", err)
	}

	// Check if user completed or cancelled
	if fm, ok := finalModel.(SetupModel); ok && fm.quitting {
		return nil // User cancelled, don't mark as error
	}

	return nil
}
