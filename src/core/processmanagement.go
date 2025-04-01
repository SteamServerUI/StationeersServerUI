// processmanagement.go
package core

import (
	"StationeersServerUI/src/config"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var (
	cmd     *exec.Cmd
	mu      sync.Mutex
	logDone chan struct{}
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func InternalStartServer() error {
	mu.Lock()
	defer mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		return fmt.Errorf("server is already running")
	}

	if config.IsDebugMode {
		fmt.Println("Config values:", config.UPNPEnabled, config.StartLocalHost, config.ServerVisible, config.AutoSave, config.AutoPauseServer, config.UseSteamP2P)
	}

	args := buildCommandArgs()
	cmd = exec.Command(config.ExePath, args...)

	fmt.Printf("\n%s%s=== GAMESERVER STARTING ===%s\n", colorCyan, colorBold, colorReset)
	fmt.Printf("• Executable: %s\n", colorGreen+colorBold+config.ExePath+colorReset)
	fmt.Printf("• Parameters: %s\n", colorYellow+strings.Join(args, " ")+colorReset)

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
		if config.IsDebugMode {
			fmt.Println("Created pipes")
		}
		// Start reading stdout and stderr pipes on Windows
		go readPipe(stdout)
		go readPipe(stderr)
	} else {
		if config.IsDebugMode {
			fmt.Println("Switching to log file for logs as we are on Linux! Hail the Penguin!")
		}
		if logDone != nil {
			close(logDone) // Close any existing channel
		}
		logDone = make(chan struct{})
		// On Linux, start the command without pipes since we're using the log file
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting server: %v", err)
		}

		// Start tailing the debug.log file on Linux
		go tailLogFile("./debug.log")
	}

	return nil
}

func InternalStopServer() error {
	mu.Lock()
	defer mu.Unlock()

	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("server is not running")
	}

	isWindows := runtime.GOOS == "windows"

	if isWindows {
		// On Windows, just kill the process directly
		if killErr := cmd.Process.Kill(); killErr != nil {
			return fmt.Errorf("error stopping server: %v", killErr)
		}
	} else {
		// On Linux/Unix, try SIGTERM first for graceful shutdown
		if termErr := cmd.Process.Signal(syscall.SIGTERM); termErr != nil {
			// If SIGTERM fails, fall back to Kill
			if killErr := cmd.Process.Kill(); killErr != nil {
				return fmt.Errorf("error stopping server: %v", killErr)
			}
		}
		// Close the logDone channel to stop the tailing goroutine (Linux only)
		if logDone != nil {
			close(logDone)
			logDone = nil // Reset to nil to avoid double-closing
		}
	}

	// Wait for the process to exit
	if waitErr := cmd.Wait(); waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") {
		// Only report actual errors, not just non-zero exit codes
		return fmt.Errorf("error during server shutdown: %v", waitErr)
	}

	cmd = nil
	return nil
}

// StartServer HTTP handler
func StartServer(w http.ResponseWriter, r *http.Request) {
	if err := InternalStartServer(); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, "Server started.")
}

// StopServer HTTP handler
func StopServer(w http.ResponseWriter, r *http.Request) {
	if err := InternalStopServer(); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, "Server stopped.")
}
