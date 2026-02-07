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
	// (logger doesn't import dashboard, dashboard imports logger).
	// See note in logger/logger.go (bottom)
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
	maxLogLines = 200
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

// CaptureLog captures a Backend log line for display in the dashboard.
// Called by the logger when dashboard is active.
func CaptureLog(line string) {
	logBufferMu.Lock()
	defer logBufferMu.Unlock()
	logBuffer = append(logBuffer, line)
	if len(logBuffer) > maxLogLines {
		logBuffer = logBuffer[1:]
	}
}

// GetLogBuffer returns a copy of the current Backend log buffer. (Dashboard-Local)
func GetLogBuffer() []string {
	logBufferMu.Lock()
	defer logBufferMu.Unlock()
	result := make([]string, len(logBuffer))
	copy(result, logBuffer)
	return result
}

// ClearLogBuffer clears the Backend log buffer. (Dashboard-Local)
func ClearLogBuffer() {
	logBufferMu.Lock()
	defer logBufferMu.Unlock()
	logBuffer = nil
}

// IsInteractiveTerminal checks if stdin/stdout are connected to a terminal.
// Returns false for Docker, systemd, or piped/redirected IO.
func IsInteractiveTerminal() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())
}

// Run takes over the terminal until the user exits.
func Run() error {

	setDashboardActive(true)
	ClearLogBuffer()

	m := NewModel()

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(), // Use alternate screen buffer (preserves shell history)
	)

	finalModel, err := p.Run()

	setDashboardActive(false)

	if err != nil {
		return fmt.Errorf("dashboard error: %w", err)
	}

	if m, ok := finalModel.(Model); ok && m.err != nil {
		return m.err
	}

	return nil
}
