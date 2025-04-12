// processmanagement.go
package gamemgr

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

var (
	cmd     *exec.Cmd
	mu      sync.Mutex
	logDone chan struct{}
)

// InternalIsServerRunning checks if the server process is running.
// Safe to call standalone as it manages its own locking.
func InternalIsServerRunning() bool {
	mu.Lock()
	defer mu.Unlock()
	return internalIsServerRunningNoLock()
}

// internalIsServerRunningNoLock checks if the server process is running.
// Caller must hold mu.Lock().
func internalIsServerRunningNoLock() bool {
	if cmd == nil || cmd.Process == nil {
		return false
	}

	if runtime.GOOS == "windows" {
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		select {
		case err := <-done:
			// process has likely exited
			if err != nil {
				logger.Core.Debug("Wait failed: " + err.Error())
				if strings.Contains(err.Error(), "The handle is invalid") {
					cmd = nil
					return false
				}
			}
			cmd = nil
			return false
		case <-time.After(50 * time.Millisecond):
			// Process is still running
			return true
		}
	} else {
		// On Unix-like systems, use Signal(0)
		if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
			logger.Core.Debug("Signal(0) failed, assuming process is dead: " + err.Error())
			cmd = nil
			return false
		}
		return true
	}
}

func InternalStartServer() error {
	mu.Lock()
	defer mu.Unlock()

	if internalIsServerRunningNoLock() {
		return fmt.Errorf("server is already running")
	}

	args := buildCommandArgs()
	cmd = exec.Command(config.ExePath, args...)

	logger.Core.Info("=== GAMESERVER STARTING ===")
	logger.Core.Info("• Executable: " + config.ExePath)
	logger.Core.Info("• Parameters: " + strings.Join(args, " "))

	if runtime.GOOS == "windows" {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("error creating StdoutPipe: %v", err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("error creating StderrPipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting server: %v", err)
		}
		logger.Core.Debug("Server process started with PID:" + strconv.Itoa(cmd.Process.Pid))
		logger.Core.Debug("Created pipes")

		// Start reading stdout and stderr pipes on Windows
		go readPipe(stdout)
		go readPipe(stderr)
	} else {

		logger.Core.Debug("Switching to log file for logs as we are on Linux! Hail the Penguin!")

		if logDone != nil {
			close(logDone) // Close any existing channel
		}
		logDone = make(chan struct{})
		// On Linux, start the command without pipes since we're using the log file
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting server: %v", err)
		}
		logger.Core.Debug("Server process started with PID:" + strconv.Itoa(cmd.Process.Pid))

		// Start tailing the debug.log file on Linux
		go tailLogFile("./debug.log")
	}

	return nil
}

func InternalStopServer() error {
	mu.Lock()
	defer mu.Unlock()

	if !internalIsServerRunningNoLock() {
		return fmt.Errorf("server not running")
	}

	// Process is running, stop it
	isWindows := runtime.GOOS == "windows"
	var killErr error

	if isWindows {
		// On Windows, attempt to terminate the process
		killErr = cmd.Process.Kill()
	} else {
		// On Linux/Unix, try SIGTERM first for graceful shutdown
		if termErr := cmd.Process.Signal(syscall.SIGTERM); termErr != nil {
			// If SIGTERM fails, fall back to Kill
			killErr = cmd.Process.Kill()
		}
		// Close the logDone channel to stop the tailing goroutine (Linux only)
		if logDone != nil {
			close(logDone)
			logDone = nil // Reset to nil to avoid double-closing
		}
	}

	// Wait for the process to exit, with a timeout to avoid hanging
	waitErrChan := make(chan error, 1)
	go func() {
		waitErrChan <- cmd.Wait()
	}()

	var waitErr error
	select {
	case waitErr = <-waitErrChan:
	case <-time.After(1 * time.Second):
		waitErr = fmt.Errorf("timeout waiting for process to exit")
	}

	// If kill failed, return that error
	if killErr != nil {
		return fmt.Errorf("error stopping server: %v", killErr)
	}

	// Handle wait errors
	if waitErr != nil {
		// If the error is "handle is invalid" or similar, assume the process is dead
		if isWindows && strings.Contains(waitErr.Error(), "The handle is invalid") {
			logger.Core.Debug("Handle is invalid during shutdown, assuming process is terminated")
		} else if !strings.Contains(waitErr.Error(), "exit status") {
			// Report other unexpected errors
			return fmt.Errorf("error during server shutdown: %v", waitErr)
		}
	}

	// Process is confirmed stopped, clear cmd
	cmd = nil
	return nil
}
