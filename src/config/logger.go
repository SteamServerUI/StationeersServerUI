// logger.go
package config

import (
	"fmt"
	"sync"
	"time"
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

var logger *Logger

// Logger configuration
type Logger struct {
	mu     sync.Mutex
	prefix string
	debug  bool
}

type LogLevel int

const (
	LogInfo LogLevel = iota
	LogWarn
	LogError
	LogDebug
	// Package-specific levels
	LogBackup
	LogDetection
	LogDiscord
	LogCore
	LogInstall
	LogSSE
)

// NewLogger creates a new logger instance
func NewLogger(prefix string, debug bool) *Logger {
	return &Logger{
		prefix: prefix,
		debug:  debug,
	}
}

// Log method for the logger
func (l *Logger) Log(level LogLevel, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level == LogDebug && !l.debug {
		return
	}

	var color, prefix string
	switch level {
	case LogInfo:
		color = colorGreen
		prefix = "INFO"
	case LogWarn:
		color = colorYellow
		prefix = "WARN"
	case LogError:
		color = colorRed
		prefix = "ERROR"
	case LogDebug:
		color = colorMagenta
		prefix = "DEBUG"
	case LogBackup:
		color = colorGreen
		prefix = "BACKUP"
	case LogDetection:
		color = colorGreen
		prefix = "DETECT"
	case LogDiscord:
		color = colorGreen
		prefix = "DISCORD"
	case LogCore:
		color = colorYellow
		prefix = "CORE"
	case LogInstall:
		color = colorCyan
		prefix = "INSTALL"
	case LogSSE:
		color = colorMagenta
		prefix = "SSE"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s%s [%s/%s] %s%s\n", color, timestamp, prefix, l.prefix, message, colorReset)
}
