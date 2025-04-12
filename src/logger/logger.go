package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"StationeersServerUI/src/config"
)

// Logger instances
var (
	Main      = &Logger{prefix: SYS_MAIN}
	Web       = &Logger{prefix: SYS_WEB}
	Discord   = &Logger{prefix: SYS_DISCORD}
	Backup    = &Logger{prefix: SYS_BACKUP}
	Detection = &Logger{prefix: SYS_DETECT}
	Core      = &Logger{prefix: SYS_CORE}
	Config    = &Logger{prefix: SYS_CONFIG}
	Install   = &Logger{prefix: SYS_INSTALL}
	SSE       = &Logger{prefix: SYS_SSE}
	Security  = &Logger{prefix: SYS_SECURITY}
)

// Severity Levels
const (
	DEBUG = 10 // Fine-grained debugging
	INFO  = 20 // Normal operations
	WARN  = 30 // Potential issues
	ERROR = 40 // Critical errors
)

// Subsystems
const (
	SYS_MAIN     = "MAIN"
	SYS_WEB      = "WEB"
	SYS_DISCORD  = "DISCORD"
	SYS_BACKUP   = "BACKUP"
	SYS_DETECT   = "DETECT"
	SYS_CORE     = "CORE"
	SYS_CONFIG   = "CONFIG"
	SYS_INSTALL  = "INSTALL"
	SYS_SSE      = "SSE"
	SYS_SECURITY = "SECURITY"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

// Subsystem color map (distinct colors, cohesive vibe)
var subsystemColors = map[string]string{
	SYS_MAIN:     colorBlue,    // Calm, default system
	SYS_WEB:      colorCyan,    // Clean, UI-related
	SYS_DISCORD:  colorMagenta, // Flashy, chatty subsystem
	SYS_BACKUP:   colorGreen,   // Safe, reliable vibe
	SYS_DETECT:   colorYellow,  // Attention-grabbing for detection
	SYS_CORE:     colorMagenta, // Critical, stands out
	SYS_CONFIG:   colorYellow,  // Warning-like, config tweaks
	SYS_INSTALL:  colorBlue,    // Matches MAIN, setup-related
	SYS_SSE:      colorCyan,    // Matches WEB, streaming vibe
	SYS_SECURITY: colorRed,     // Screams "pay attention"
}

type Logger struct {
	mu     sync.Mutex
	prefix string // Subsystem identifier (e.g., "DISCORD")
}

type logEntry struct {
	severity int
	prefix   string // Log type (e.g., "INFO", "CORE")
	color    string
	message  string
}

// shouldLog checks severity and subsystem filters
func (l *Logger) shouldLog(severity int) bool {
	// Subsystem filtering first
	if len(config.SubsystemFilters) > 0 {
		allowed := false
		for _, sub := range config.SubsystemFilters {
			if sub == l.prefix {
				allowed = true
				break
			}
		}
		if !allowed {
			return false // Subsystem not in filter, skip it
		}
	}

	// Existing severity logic
	effectiveLevel := config.LogLevel
	if config.IsDebugMode && effectiveLevel < DEBUG {
		effectiveLevel = 10 // Force DEBUG if IsDebugMode is true
	}
	return severity >= effectiveLevel
}

// log handles the core logging logic
func (l *Logger) log(entry logEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.shouldLog(entry.severity) {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// Use subsystem color by default, override with severity color if set
	entryColor := subsystemColors[l.prefix]
	if entry.color != colorReset {
		entryColor = entry.color
	}
	// Console version with colors
	consoleLine := fmt.Sprintf("%s%s [%s/%s] %s%s\n", entryColor, timestamp, entry.prefix, l.prefix, entry.message, colorReset)
	// File version without colors
	fileLine := fmt.Sprintf("%s [%s/%s] %s\n", timestamp, entry.prefix, l.prefix, entry.message)
	// Console output
	fmt.Print(consoleLine)

	// File output if enabled
	if config.CreateSSUILogFile {
		l.writeToFile(fileLine)
	}
}

func (l *Logger) writeToFile(logLine string) {
	// Retry settings
	const maxRetries = 5
	const retryDelay = 100 * time.Millisecond

	// Retry loop with timeout
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Open file with proper flags
		file, err := os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			// Successfully opened file, proceed with writing
			defer file.Close()
			if _, err := file.WriteString(logLine); err != nil {
				fmt.Printf("%s%s [ERROR/LOGGER] Failed to write to log file: %v%s\n",
					colorRed,
					time.Now().Format("2006-01-02 15:04:05"),
					err,
					colorReset)
			}
			return
		}

		// If the error is due to the directory not existing
		if os.IsNotExist(err) {
			if attempt == maxRetries-1 {
				// Last attempt failed, ditch the log line
				return
			}
			// Wait before retrying
			time.Sleep(retryDelay)
			continue
		}

		// For other errors, log and exit immediately
		fmt.Printf("%s%s [ERROR/LOGGER] Failed to open log file: %v%s\n",
			colorRed,
			time.Now().Format("2006-01-02 15:04:05"),
			err,
			colorReset)
		return
	}
}

// Severity-based methods
func (l *Logger) Debug(message string) {
	l.log(logEntry{DEBUG, "DEBUG", colorReset, message}) // Subsystem color
}

func (l *Logger) Info(message string) {
	l.log(logEntry{INFO, "INFO", colorReset, message}) // Subsystem color
}

func (l *Logger) Warn(message string) {
	l.log(logEntry{WARN, "WARN", colorYellow, message}) // Yellow for warnings
}

func (l *Logger) Error(message string) {
	l.log(logEntry{ERROR, "ERROR", colorRed, message}) // Red for errors
}

// Subsystem-specific methods (update colors for consistency)
func (l *Logger) Backup(message string) {
	l.log(logEntry{INFO, "BACKUP", colorReset, message}) // Green via subsystem
}

func (l *Logger) Detection(message string) {
	l.log(logEntry{INFO, "DETECT", colorReset, message}) // Yellow via subsystem
}

func (l *Logger) Discord(message string) {
	l.log(logEntry{INFO, "DISCORD", colorReset, message}) // Magenta via subsystem
}

func (l *Logger) Core(message string) {
	l.log(logEntry{WARN, "CORE", colorReset, message}) // Magenta via subsystem
}

func (l *Logger) Config(message string) {
	l.log(logEntry{WARN, "CONFIG", colorReset, message}) // Yellow via subsystem
}

func (l *Logger) Install(message string) {
	l.log(logEntry{INFO, "INSTALL", colorReset, message}) // Blue via subsystem
}

func (l *Logger) SSE(message string) {
	l.log(logEntry{INFO, "SSE", colorReset, message}) // Cyan via subsystem
}

func (l *Logger) Security(message string) {
	l.log(logEntry{ERROR, "SECURITY", colorReset, message}) // Red via subsystem
}
