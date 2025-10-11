package gamemgr

import (
	"runtime"
	"syscall"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

// InternalIsServerRunning checks if the server process is running.
// Safe to call standalone as it manages its own locking.
func InternalIsServerRunning() bool {
	mu.Lock()
	defer mu.Unlock()
	return internalIsServerRunningNoLock()
}

// internalIsServerRunningNoLock checks if the server process is running.
// Caller M U S T hold mu.Lock().
func internalIsServerRunningNoLock() bool {
	if cmd == nil || cmd.Process == nil {
		return false
	}

	if runtime.GOOS == "windows" {
		select {
		case <-processExited:
			cmd = nil
			clearGameServerUUID()
			return false
		default:
			// Process is still running
			return true
		}
	}

	if runtime.GOOS == "linux" {
		// On Unix-like systems, use Signal(0)
		if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
			logger.Core.Debug("Signal(0) failed, assuming process is dead: " + err.Error())
			cmd = nil
			clearGameServerUUID()
			return false
		}
		return true
	}

	logger.Core.Warn("Failed to check if server is running, assuming it's dead")
	return false
}
