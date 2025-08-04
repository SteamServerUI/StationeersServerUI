// processmanagement.go
package gamemgr

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/commandmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/google/uuid"
)

var (
	cmd             *exec.Cmd
	mu              sync.Mutex
	logDone         chan struct{}
	err             error
	autoRestartDone chan struct{}
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
					clearGameServerUUID()
					return false
				}
			}
			cmd = nil
			clearGameServerUUID()
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
			clearGameServerUUID()
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

	logger.Core.Info("=== GAMESERVER STARTING ===")

	if config.IsSSCMEnabled && runtime.GOOS == "linux" {

		var envVars []string
		// Set up SSCM (BepInEx/Doorstop) environment
		envVars, err = SetupBepInExEnvironment()
		if err != nil {
			return fmt.Errorf("failed to set up SSCM environment: %v", err)
		}
		// Create command after environment is set
		cmd = exec.Command(config.ExePath, args...)
		// Set the environment for the command
		if envVars != nil {
			cmd.Env = envVars
			logger.Core.Info("BepInEx/Doorstop environment configured for server process")
		}
		logger.Core.Info("• Executable: " + config.ExePath + " (with SSCM)")
		logger.Core.Info("• Arguments: " + strings.Join(args, " "))
	}

	if !config.IsSSCMEnabled && runtime.GOOS == "linux" {
		// Use ExePath directly as the command
		cmd = exec.Command(config.ExePath, args...)
		logger.Core.Info("• Executable: " + config.ExePath)
		logger.Core.Info("• Arguments: " + strings.Join(args, " "))
	}

	if runtime.GOOS == "windows" {

		// On Windows, set the command to use the executable path and arguments
		cmd = exec.Command(config.ExePath, args...)
		logger.Core.Info("• Executable: " + config.ExePath)
		logger.Core.Debug("Switching to pipes for logs as we are on Windows!")

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
		logger.Core.Info("• Arguments: " + strings.Join(args, " "))
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

		// check if debug.log file exists, if not, create it
		if _, err := os.Stat("./debug.log"); os.IsNotExist(err) {
			file, err := os.OpenFile("./debug.log", os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating debug.log file: %v", err)
			}
			defer file.Close()
		}
		// Start tailing the debug.log file on Linux
		go tailLogFile("./debug.log")
	}
	// create a UUID for this specific run
	createGameServerUUID()
	logger.Core.Debug("Created Game Server with internal UUID: " + config.GameServerUUID.String())

	// Start auto-restart goroutine if AutoRestartServerTimer is set greater than 0
	if config.AutoRestartServerTimer != "0" {
		if autoRestartDone != nil {
			close(autoRestartDone)
		}
		autoRestartDone = make(chan struct{})
		go startAutoRestart(config.AutoRestartServerTimer, autoRestartDone)
		logger.Core.Info("Auto-restart scheduled every " + config.AutoRestartServerTimer + " minutes")
	}

	return nil
}

func InternalStopServer() error {
	mu.Lock()
	defer mu.Unlock()

	if !internalIsServerRunningNoLock() {
		return fmt.Errorf("server not running")
	}

	// Stop auto-restart goroutine
	if autoRestartDone != nil {
		close(autoRestartDone)
		autoRestartDone = nil
		logger.Core.Info("Auto-restart cycle interrupted due to manaual stop")
	}

	// Process is running, stop it
	isWindows := runtime.GOOS == "windows"
	var killErr error

	if isWindows {
		// On Windows, terminate the process (no graceful shutdown)
		killErr = cmd.Process.Kill()
	} else {
		// On Linux/Unix, send SIGTERM for graceful shutdown
		if termErr := cmd.Process.Signal(syscall.SIGTERM); termErr != nil {
			logger.Core.Debug("SIGTERM failed: " + termErr.Error())
			killErr = cmd.Process.Kill() // Fallback to Kill if SIGTERM fails
		} else {
			// Wait for graceful shutdown
			waitErrChan := make(chan error, 1)
			go func() {
				waitErrChan <- cmd.Wait()
			}()

			select {
			case waitErr := <-waitErrChan:
				if waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") {
					logger.Core.Debug("Wait error after SIGTERM: " + waitErr.Error())
				}
			case <-time.After(10 * time.Second): // Increased timeout
				logger.Core.Warn("Timeout waiting for graceful shutdown, sending SIGKILL")
				killErr = cmd.Process.Kill() // Fallback to SIGKILL
				select {
				case waitErr := <-waitErrChan:
					if waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") {
						logger.Core.Debug("Wait error after SIGKILL: " + waitErr.Error())
					}
				case <-time.After(2 * time.Second): // Additional wait for SIGKILL
					return fmt.Errorf("timeout waiting for process to exit after SIGKILL")
				}
			}
		}

		// Stop log tailing (Linux only)
		if logDone != nil {
			close(logDone)
			logDone = nil
		}
	}

	// For Windows, wait briefly after Kill to ensure process is gone
	if isWindows {
		waitErrChan := make(chan error, 1)
		go func() {
			waitErrChan <- cmd.Wait()
		}()

		select {
		case waitErr := <-waitErrChan:
			if waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") &&
				!strings.Contains(waitErr.Error(), "The handle is invalid") {
				return fmt.Errorf("error during server shutdown: %v", waitErr)
			}
		case <-time.After(1 * time.Second):
			return fmt.Errorf("timeout waiting for process to exit")
		}
	}

	if killErr != nil {
		return fmt.Errorf("error stopping server: %v", killErr)
	}

	// Process is confirmed stopped, clear cmd
	cmd = nil
	clearGameServerUUID()
	return nil
}

// startAutoRestart runs a goroutine that restarts the server after the specified timeframe in minutes.
func startAutoRestart(minutes string, done chan struct{}) {
	minutesInt, _ := strconv.Atoi(minutes)
	ticker := time.NewTicker(time.Duration(minutesInt) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mu.Lock()
			if !internalIsServerRunningNoLock() {
				mu.Unlock()
				logger.Core.Info("Auto-restart skipped: server is not running")
				return
			}
			mu.Unlock()

			if config.IsSSCMEnabled {
				commandmgr.WriteCommand("say Attention, server is restarting in 30 seconds!")
				time.Sleep(10 * time.Second)
				commandmgr.WriteCommand("say Attention, server is restarting in 20 seconds!")
				time.Sleep(10 * time.Second)
				commandmgr.WriteCommand("say Attention, server is restarting in 10 seconds!")
				time.Sleep(5 * time.Second)
				commandmgr.WriteCommand("say Attention, server is restarting in 5 seconds!")
				time.Sleep(5 * time.Second)
			}
			logger.Core.Info("Auto-restart triggered: stopping server")
			if err := InternalStopServer(); err != nil {
				logger.Core.Error("Auto-restart failed to stop server: " + err.Error())
				return
			}

			logger.Core.Info("Auto-restart: waiting 5 seconds before restarting")
			time.Sleep(5 * time.Second)

			logger.Core.Info("Auto-restart: starting server")
			if err := InternalStartServer(); err != nil {
				logger.Core.Error("Auto-restart failed to start server: " + err.Error())
				return
			}
		case <-done:
			return
		}
	}
}

func clearGameServerUUID() {
	config.ConfigMu.Lock()
	defer config.ConfigMu.Unlock()
	config.GameServerUUID = uuid.Nil
}

func createGameServerUUID() {
	config.ConfigMu.Lock()
	defer config.ConfigMu.Unlock()
	config.GameServerUUID = uuid.New()
}
