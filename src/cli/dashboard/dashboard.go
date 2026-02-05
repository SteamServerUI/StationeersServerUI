// Package dashboard provides an interactive terminal UI using Bubble Tea.
// It offers a dashboard-like experience with server status, logs, and controls.
package dashboard

import (
	"fmt"
	"os"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
)

func init() {
	// Register hooks with the logger to enable log capture during dashboard mode.
	// This avoids import cycles (logger doesn't import dashboard).
	logger.RegisterDashboardHooks(IsDashboardActive, CaptureLog)
}

var (
	// dashboardActive tracks whether the dashboard is currently running.
	// Used by the logger to suppress console output during dashboard mode.
	dashboardActive bool
	dashboardMu     sync.RWMutex

	// logBuffer holds SSUI log messages captured while dashboard is active.
	logBuffer   []string
	logBufferMu sync.Mutex
	maxLogLines = 200 // Rolling buffer size
)

// IsDashboardActive returns whether the dashboard is currently running.
// The logger can use this to suppress direct console writes.
func IsDashboardActive() bool {
	dashboardMu.RLock()
	defer dashboardMu.RUnlock()
	return dashboardActive
}

// setDashboardActive sets the dashboard active state.
func setDashboardActive(active bool) {
	dashboardMu.Lock()
	defer dashboardMu.Unlock()
	dashboardActive = active
}

// CaptureLog captures a SSUI log line for display in the dashboard.
// Called by the logger when dashboard is active.
func CaptureLog(line string) {
	logBufferMu.Lock()
	defer logBufferMu.Unlock()
	logBuffer = append(logBuffer, line)
	if len(logBuffer) > maxLogLines {
		logBuffer = logBuffer[1:]
	}
}

// GetLogBuffer returns a copy of the current SSUI log buffer.
func GetLogBuffer() []string {
	logBufferMu.Lock()
	defer logBufferMu.Unlock()
	result := make([]string, len(logBuffer))
	copy(result, logBuffer)
	return result
}

// ClearLogBuffer clears the SSUI log buffer.
func ClearLogBuffer() {
	logBufferMu.Lock()
	defer logBufferMu.Unlock()
	logBuffer = nil
}

// IsInteractiveTerminal checks if stdin/stdout are connected to a terminal.
// Returns false for Docker, systemd, or piped/redirected IO.
// Note: go-isatty works on Windows too (uses GetConsoleMode).
func IsInteractiveTerminal() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())
}

// Run starts the interactive dashboard.
// It takes over the terminal until the user exits.
// Returns an error if the dashboard fails to start.
func Run() error {
	// Signal that dashboard is active - logger should suppress console output
	setDashboardActive(true)
	ClearLogBuffer()

	// Create the initial model
	m := NewModel()

	// Create the Bubble Tea program with options for full terminal control
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer (preserves shell history)
		tea.WithMouseCellMotion(), // Enable mouse support for future interactions
	)

	// Run the program (blocking)
	finalModel, err := p.Run()

	// Signal that dashboard is no longer active
	setDashboardActive(false)

	if err != nil {
		return fmt.Errorf("dashboard error: %w", err)
	}

	// Check if there was an error in the final model
	if m, ok := finalModel.(Model); ok && m.err != nil {
		return m.err
	}

	return nil
}

// RunWithOutput is similar to Run but captures output for testing.
func RunWithOutput(input *os.File, output *os.File) error {
	setDashboardActive(true)
	ClearLogBuffer()

	m := NewModel()

	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
	}
	if input != nil {
		opts = append(opts, tea.WithInput(input))
	}
	if output != nil {
		opts = append(opts, tea.WithOutput(output))
	}

	p := tea.NewProgram(m, opts...)

	_, err := p.Run()
	setDashboardActive(false)

	return err
}
