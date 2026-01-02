// processmanagement.go
package gamemgr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/runfile"
)

var (
	cmd           *exec.Cmd
	mu            sync.Mutex
	logDone       chan struct{}
	processExited chan struct{}
	// autoRestartDone is defined in autorestart.go
)

func InternalStartServer() error {
	mu.Lock()
	defer mu.Unlock()

	if config.GetIsGameServerRunning() {
		return fmt.Errorf("server is already running")
	}

	args, err := runfile.BuildCommandArgs()
	if err != nil {
		logger.Core.Error("Failed to build command args: " + err.Error())
		return err
	}

	executable, err := runfile.CurrentRunfile.GetExecutable()
	if err != nil {
		logger.Core.Error("Failed to get executable path from runfile: " + err.Error())
		return err
	}
	executablePath := executable

	logger.Core.Info("=== GAMESERVER STARTING ===")
	logger.Core.Info("BepInEx/Doorstop enabled: " + strconv.FormatBool(config.GetIsBepInExEnabled()))

	if config.GetIsBepInExEnabled() && runtime.GOOS == "linux" {

		var envVars []string
		// Set up SSCM (BepInEx/Doorstop) environment
		envVars, err = SetupBepInExEnvironment()
		if err != nil {
			return fmt.Errorf("failed to set up SSCM environment: %v", err)
		}
		// Create command after environment is set
		cmd = exec.Command(executablePath, args...)
		// Set the environment for the command
		if envVars != nil {
			cmd.Env = envVars
			logger.Core.Info("BepInEx/Doorstop environment configured for server process")
		}
	} else {
		// Use ExePath directly as the command for non-BepInEx or Windows
		cmd = exec.Command(executablePath, args...)
	}

	// Log executable and arguments
	logger.Core.Info("• Executable: " + executablePath)
	var formattedArgs []string
	for _, arg := range args {
		if strings.ContainsAny(arg, " \t\n\"'") {
			formattedArgs = append(formattedArgs, `"`+strings.ReplaceAll(arg, `"`, `\"`)+`"`)
		} else {
			formattedArgs = append(formattedArgs, arg)
		}
	}
	logger.Core.Info("• Arguments: " + strings.Join(formattedArgs, " "))

	// get current working directory of SSUI
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %v", err)
	}
	childWD := filepath.Join(cwd, config.GetRunfileIdentifier())

	logger.Core.Debug("Child working directory: " + childWD)

	cmd.Dir = childWD
	logger.Core.Debug("Set gamservers working directory to: " + cmd.Dir)
	// Handle log reading based on GetGameLogFromLogFile
	if config.GetGameLogFromLogFile() {
		logger.Core.Debug("Switching to log file tailing for logs")

		// Check if gameserver.log file exists, if not, create it
		logFilePath := filepath.Join(config.GetRunfileIdentifier(), "gameserver.log")
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating gameserver.log file: %v", err)
			}
			file.Close()
		}

		// Start tailing the gameserver.log file
		if logDone != nil {
			close(logDone) // Close any existing channel
		}
		logDone = make(chan struct{})
		go tailLogFile(logFilePath)

		// Start the command without pipes
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting server: %v", err)
		}
		logger.Core.Debug("Server process started with PID: " + strconv.Itoa(cmd.Process.Pid))
	} else {
		// Use pipes for log reading
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
		logger.Core.Debug("Server process started with PID: " + strconv.Itoa(cmd.Process.Pid))
		logger.Core.Debug("Created pipes")

		// Start reading stdout and stderr pipes
		go readPipe(stdout)
		go readPipe(stderr)
	}

	// Monitor process exit
	processExited = make(chan struct{})
	go func() {
		err := cmd.Wait()
		if err != nil {
			logger.Core.Debug("Process exited with error: " + err.Error())
		} else {
			logger.Core.Debug("Process exited successfully")
		}
		close(processExited)
	}()

	// Create a UUID for this specific run
	createGameServerUUID()
	config.SetIsGameServerRunning(true)

	// Start auto-restart goroutine if AutoRestartServerTimer is set greater than 0
	if config.GetAutoRestartServerTimer() != "0" {
		if autoRestartDone != nil {
			close(autoRestartDone)
		}
		autoRestartDone = make(chan struct{})
		go startAutoRestart(config.GetAutoRestartServerTimer(), autoRestartDone)
		logger.Core.Info("New Auto-restart scheduled: " + config.GetAutoRestartServerTimer())
	}

	return nil
}

func InternalStopServer() error {
	mu.Lock()
	defer mu.Unlock()

	if !config.GetIsGameServerRunning() {
		return fmt.Errorf("server not running")
	}

	// Stop auto-restart goroutine
	if autoRestartDone != nil {
		close(autoRestartDone)
		autoRestartDone = nil
	}

	// Process is running, stop it
	isWindows := runtime.GOOS == "windows"
	var killErr error

	if isWindows {
		// On Windows, terminate the process (no graceful shutdown)
		killErr = cmd.Process.Kill()
		// Wait for the processExited channel to confirm exit
		if processExited != nil {
			select {
			case <-processExited:
				logger.Core.Debug("processExited channel confirmed server shutdown")
			case <-time.After(2 * time.Second):
				logger.Core.Warn("Timeout waiting for processExited confirmation")
			}
		}
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

	if killErr != nil {
		return fmt.Errorf("error stopping server: %v", killErr)
	}

	// Process is confirmed stopped, clear cmd
	cmd = nil
	config.SetIsGameServerRunning(false)
	clearGameServerUUID()
	return nil
}
