package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"StationeersServerUI/src/config"
)

// Severity Levels
const (
	DEBUG = 10 // Fine-grained debugging
	INFO  = 20 // Normal operations
	WARN  = 30 // Potential issues
	ERROR = 40 // Critical errors
)

// Subsystems (not levels, just identifiers)
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

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

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

// shouldLog checks severity against config.LogLevel
func (l *Logger) shouldLog(severity int) bool {
	effectiveLevel := config.LogLevel
	if config.IsDebugMode && effectiveLevel < DEBUG {
		effectiveLevel = 10 // Force DEBUG if IsDebugMode is true
	}
	// Add subsystem filtering later if needed via config
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
	// Console version with colors
	consoleLine := fmt.Sprintf("%s%s [%s/%s] %s%s\n", entry.color, timestamp, entry.prefix, l.prefix, entry.message, colorReset)
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
	l.log(logEntry{DEBUG, "DEBUG", colorMagenta, message})
}

func (l *Logger) Info(message string) {
	l.log(logEntry{INFO, "INFO", colorReset, message})
}

func (l *Logger) Warn(message string) {
	l.log(logEntry{WARN, "WARN", colorYellow, message})
}

func (l *Logger) Error(message string) {
	l.log(logEntry{ERROR, "ERROR", colorRed, message})
}

// Subsystem-specific methods
func (l *Logger) Backup(message string) {
	l.log(logEntry{INFO, "BACKUP", colorGreen, message})
}

func (l *Logger) Detection(message string) {
	l.log(logEntry{INFO, "DETECT", colorGreen, message})
}

func (l *Logger) Discord(message string) {
	l.log(logEntry{INFO, "DISCORD", colorGreen, message})
}

func (l *Logger) Core(message string) {
	l.log(logEntry{WARN, "CORE", colorYellow, message})
}

func (l *Logger) Config(message string) {
	l.log(logEntry{WARN, "CONFIG", colorYellow, message})
}

func (l *Logger) Install(message string) {
	l.log(logEntry{INFO, "INSTALL", colorCyan, message})
}

func (l *Logger) SSE(message string) {
	l.log(logEntry{INFO, "SSE", colorMagenta, message})
}

func (l *Logger) Security(message string) {
	l.log(logEntry{ERROR, "SECURITY", colorRed, message})
}

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
