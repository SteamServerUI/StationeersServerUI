package logger

import (
	"fmt"
	"sync"
	"time"

	"StationeersServerUI/src/config"
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
	prefix string
}

func (l *Logger) log(prefix, color string, message string, isDebug bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if isDebug && !config.IsDebugMode {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s%s [%s/%s] %s%s\n", color, timestamp, prefix, l.prefix, message, colorReset)
}

func (l *Logger) Info(message string) {
	l.log("INFO", colorGreen, message, false)
}

func (l *Logger) Warn(message string) {
	l.log("WARN", colorYellow, message, false)
}

func (l *Logger) Error(message string) {
	l.log("ERROR", colorRed, message, false)
}

func (l *Logger) Debug(message string) {
	l.log("DEBUG", colorMagenta, message, true)
}

func (l *Logger) Backup(message string) {
	l.log("BACKUP", colorGreen, message, false)
}

func (l *Logger) Detection(message string) {
	l.log("DETECT", colorGreen, message, false)
}

func (l *Logger) Discord(message string) {
	l.log("DISCORD", colorGreen, message, false)
}

func (l *Logger) Core(message string) {
	l.log("CORE", colorYellow, message, false)
}

func (l *Logger) Config(message string) {
	l.log("CONFIG", colorYellow, message, false)
}

func (l *Logger) Install(message string) {
	l.log("INSTALL", colorCyan, message, false)
}

func (l *Logger) SSE(message string) {
	l.log("SSE", colorMagenta, message, false)
}

var (
	Main      = &Logger{prefix: "MAIN"}
	Web       = &Logger{prefix: "WEB"}
	Discord   = &Logger{prefix: "DISCORD"}
	Backup    = &Logger{prefix: "BACKUP"}
	Detection = &Logger{prefix: "DETECT"}
	Core      = &Logger{prefix: "CORE"}
	Config    = &Logger{prefix: "CONFIG"}
	Install   = &Logger{prefix: "INSTALL"}
	SSE       = &Logger{prefix: "SSE"}
	Auth      = &Logger{prefix: "AUTH"}
)
