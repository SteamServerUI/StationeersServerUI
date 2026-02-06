package gamemgr

import (
	"fmt"
	"sync"
	"time"
)

var (
	serverStartTime time.Time
	uptimeMu        sync.RWMutex
)

func setServerStartTime() {
	uptimeMu.Lock()
	defer uptimeMu.Unlock()
	serverStartTime = time.Now()
}

func clearServerStartTime() {
	uptimeMu.Lock()
	defer uptimeMu.Unlock()
	serverStartTime = time.Time{}
}

// GetServerUptime returns the server uptime as a DURATION.
// Returns 0 if the server is not running.
func GetServerUptime() time.Duration {
	uptimeMu.RLock()
	defer uptimeMu.RUnlock()
	if serverStartTime.IsZero() {
		return 0
	}
	return time.Since(serverStartTime)
}

// GetServerStartTime returns the time WHEN the server was started.
// Returns zero time if the server is not running.
func GetServerStartTime() time.Time {
	uptimeMu.RLock()
	defer uptimeMu.RUnlock()
	return serverStartTime
}

// FormatUptime formats the uptime duration into a human-readable string
func FormatUptime(d time.Duration) string {
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
