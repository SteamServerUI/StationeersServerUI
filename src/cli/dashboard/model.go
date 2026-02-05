package dashboard

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Panel represents which panel is currently displayed (full screen)
type Panel int

const (
	PanelSSUILog Panel = iota // SSUI backend logs (default)
	PanelStatus               // Server status overview
	PanelPlayers              // Connected players list
	PanelCount                // Used for cycling
)

// PanelNames for display
var panelNames = map[Panel]string{
	PanelSSUILog: "SSUI Log",
	PanelStatus:  "Status",
	PanelPlayers: "Players",
}

// Model represents the dashboard state
type Model struct {
	// Window dimensions
	width  int
	height int

	// Current panel (shown full screen)
	activePanel Panel

	// Components
	logViewport viewport.Model // For SSUI backend logs
	help        help.Model
	keys        keyMap

	// Server state
	serverRunning    bool
	serverUptime     time.Duration
	serverName       string
	saveName         string
	worldID          string
	gamePort         string
	maxPlayers       int
	connectedPlayers map[string]string // SteamID -> PlayerName

	// SSUI info
	ssuiVersion string
	goRuntime   string // e.g., "linux/amd64"

	// Internal
	err       error
	quitting  bool
	startTime time.Time
}

// keyMap defines the key bindings for the dashboard
type keyMap struct {
	Tab     key.Binding
	Up      key.Binding
	Down    key.Binding
	Start   key.Binding
	Stop    key.Binding
	Refresh key.Binding
	Help    key.Binding
	Quit    key.Binding
}

// ShortHelp returns the short help text (displayed in footer)
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Start, k.Stop, k.Help, k.Quit}
}

// FullHelp returns the full help text (displayed when '?' is pressed)
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.Up, k.Down},
		{k.Start, k.Stop, k.Refresh},
		{k.Help, k.Quit},
	}
}

// defaultKeyMap returns the default key bindings
func defaultKeyMap() keyMap {
	return keyMap{
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch view"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "scroll up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "scroll down"),
		),
		Start: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "start server"),
		),
		Stop: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "stop server"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q/esc", "quit"),
		),
	}
}

// NewModel creates a new dashboard model with initial values
func NewModel() Model {
	// Initialize viewport for logs
	vp := viewport.New(80, 20)
	vp.SetContent("Waiting for SSUI logs...")

	return Model{
		activePanel:      PanelSSUILog, // Default to SSUI log view
		logViewport:      vp,
		help:             help.New(),
		keys:             defaultKeyMap(),
		startTime:        time.Now(),
		connectedPlayers: make(map[string]string),
		ssuiVersion:      "Loading...",
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	// Start with a tick to update data periodically
	return tea.Batch(
		tickCmd(),
		tea.SetWindowTitle("SSUI Dashboard"),
	)
}

// tickMsg is sent periodically to update the dashboard
type tickMsg time.Time

// tickCmd returns a command that sends a tick message every second
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// statusUpdateMsg is sent when server status is updated
type statusUpdateMsg struct {
	running    bool
	uptime     time.Duration
	serverName string
	saveName   string
	worldID    string
	gamePort   string
	maxPlayers int
	players    map[string]string
	version    string
	goRuntime  string
}

// logUpdateMsg is sent when new SSUI logs are available
type logUpdateMsg struct {
	logs []string
}

// serverActionMsg is sent after a server action (start/stop)
type serverActionMsg struct {
	action string
	err    error
}

// formatUptime formats duration to a human-readable string
func formatUptime(d time.Duration) string {
	if d == 0 {
		return "N/A"
	}
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// parseMaxPlayers converts the string max players to int
func parseMaxPlayers(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
