// processmanagement.go
package core

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/ssestream"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	cmd *exec.Cmd
	mu  sync.Mutex
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

type Arg struct {
	Flag          string
	Value         string
	RequiresValue bool
	Condition     func() bool
	NoQuote       bool
}

func buildCommandArgs() []string {
	var argOrder = []Arg{
		{Flag: "-nographics", RequiresValue: false},
		{Flag: "-batchmode", RequiresValue: false},
		{Flag: "-LOAD", Value: config.SaveFileName, RequiresValue: true, NoQuote: true}, // LOAD has special handling because the gameserver expects 2 parameters
		{Flag: "-logFile", Value: "./debug.log", Condition: func() bool { return runtime.GOOS == "linux" }, RequiresValue: true},
		{Flag: "-settings", RequiresValue: false},
		{Flag: "StartLocalHost", Value: strconv.FormatBool(config.StartLocalHost), RequiresValue: true},
		{Flag: "ServerVisible", Value: strconv.FormatBool(config.ServerVisible), RequiresValue: true},
		{Flag: "GamePort", Value: config.GamePort, RequiresValue: true},
		{Flag: "UPNPEnabled", Value: strconv.FormatBool(config.UPNPEnabled), RequiresValue: true},
		{Flag: "ServerName", Value: config.ServerName, RequiresValue: true},
		{Flag: "ServerPassword", Value: config.ServerPassword, Condition: func() bool { return config.ServerPassword != "" }, RequiresValue: true},
		{Flag: "ServerMaxPlayers", Value: config.ServerMaxPlayers, RequiresValue: true},
		{Flag: "AutoSave", Value: strconv.FormatBool(config.AutoSave), RequiresValue: true},
		{Flag: "SaveInterval", Value: config.SaveInterval, RequiresValue: true},
		{Flag: "ServerAuthSecret", Value: config.ServerAuthSecret, Condition: func() bool { return config.ServerAuthSecret != "" }, RequiresValue: true},
		{Flag: "UpdatePort", Value: config.UpdatePort, RequiresValue: true},
		{Flag: "AutoPauseServer", Value: strconv.FormatBool(config.AutoPauseServer), RequiresValue: true},
		{Flag: "UseSteamP2P", Value: strconv.FormatBool(config.UseSteamP2P), RequiresValue: true},
		{Flag: "AdminPassword", Value: config.AdminPassword, Condition: func() bool { return config.AdminPassword != "" }, RequiresValue: true},
	}

	var args []string
	for _, arg := range argOrder {
		if arg.Condition != nil && !arg.Condition() {
			continue
		}
		if arg.RequiresValue && arg.Value == "" {
			continue
		}

		args = append(args, arg.Flag)

		if arg.Flag == "-LOAD" && arg.Value != "" {
			parts := strings.SplitN(arg.Value, " ", 2)
			for _, part := range parts {
				if part != "" {
					args = append(args, part)
				}
			}
			continue
		}

		if arg.Value != "" {
			args = append(args, arg.Value)
		}
	}

	if config.AdditionalParams != "" {
		args = append(args, strings.Fields(config.AdditionalParams)...)
	}

	if config.LocalIpAddress != "" {
		args = append(args, "LocalIpAddress")
		args = append(args, config.LocalIpAddress)
	}

	if config.IsDebugMode {
		fmt.Println("=== DEBUG: Raw arguments passed to exec.Command ===")
		for i, arg := range args {
			fmt.Printf("Arg[%d]: %q\n", i, arg)
		}
	}
	return args
}

func StartServer(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		fmt.Fprintf(w, "Server is already running.")
		return
	}
	if config.IsDebugMode {
		fmt.Println("Config values:", config.UPNPEnabled, config.StartLocalHost, config.ServerVisible, config.AutoSave, config.AutoPauseServer, config.UseSteamP2P)
	}
	args := buildCommandArgs()
	cmd = exec.Command(config.ExePath, args...)

	fmt.Printf("\n%s%s=== GAMESERVER STARTING ===%s\n", colorCyan, colorBold, colorReset)
	fmt.Printf("• Executable: %s\n", colorGreen+colorBold+config.ExePath+colorReset)
	fmt.Printf("• Parameters: %s\n", colorYellow+strings.Join(args, " ")+colorReset)

	// Only set up pipes for Windows
	if runtime.GOOS == "windows" {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintf(w, "Error creating StdoutPipe: %v", err)
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Fprintf(w, "Error creating StderrPipe: %v", err)
			return
		}

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(w, "Error starting server: %v", err)
			return
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
		// On Linux, start the command without pipes since we're using the log file
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(w, "Error starting server: %v", err)
			return
		}

		// Start tailing the debug.log file on Linux
		go tailLogFile("./debug.log")
	}

	fmt.Fprintf(w, "Server started.")
}

// readPipe for Windows
func readPipe(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	if config.IsDebugMode {
		fmt.Println("Started reading pipe") // Debug
	}
	for scanner.Scan() {
		output := scanner.Text()
		ssestream.BroadcastConsoleOutput(output)
	}
	if err := scanner.Err(); err != nil {
		if config.IsDebugMode {
			fmt.Println("Pipe error:", err) // Debug
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reading pipe: %v", err))
	}
	if config.IsDebugMode {
		fmt.Println("Pipe closed") // Debug
	}
}

// tailLogFile implements a tail -f-like behavior for Linux
func tailLogFile(logFilePath string) {
	// Wait briefly to ensure the file exists after server start
	time.Sleep(1 * time.Second)

	file, err := os.Open(logFilePath)
	if err != nil {
		if config.IsDebugMode {
			fmt.Printf("Error opening log file: %v\n", err)
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error opening log file: %v", err))
		return
	}
	defer file.Close()

	// Seek to the end of the file initially (like tail -f)
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		if config.IsDebugMode {
			fmt.Printf("Error seeking to end of log file: %v\n", err)
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error seeking log file: %v", err))
		return
	}

	scanner := bufio.NewScanner(file)
	if config.IsDebugMode {
		fmt.Println("Started tailing log file") // Debug
	}

	// Continuously read new lines
	for {
		for scanner.Scan() {
			output := scanner.Text()
			ssestream.BroadcastConsoleOutput(output)
		}

		// If we reach EOF, wait and check for new content
		if err := scanner.Err(); err != nil {
			if config.IsDebugMode {
				fmt.Printf("Error reading log file: %v\n", err)
			}
			ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reading log file: %v", err))
			return
		}

		// Sleep briefly before checking for new content
		time.Sleep(100 * time.Millisecond)

		// Check if the file has been truncated or rotated (optional handling)
		currentPos, _ := file.Seek(0, io.SeekCurrent)
		fileInfo, err := file.Stat()
		if err != nil {
			continue
		}
		if currentPos > fileInfo.Size() {
			// File was truncated or rotated, reset to start
			file.Seek(0, io.SeekStart)
		}
	}
}

func GetLogOutput(w http.ResponseWriter, r *http.Request) {
	ssestream.StartConsoleStream()(w, r)
}

func StopServer(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if cmd == nil || cmd.Process == nil {
		fmt.Fprintf(w, "Server is not running.")
		return
	}

	isWindows := runtime.GOOS == "windows"

	if isWindows {
		// On Windows, just kill the process directly
		if killErr := cmd.Process.Kill(); killErr != nil {
			fmt.Fprintf(w, "Error stopping server: %v", killErr)
			return
		}
	} else {
		// On Linux/Unix, try SIGTERM first for graceful shutdown
		if termErr := cmd.Process.Signal(syscall.SIGTERM); termErr != nil {
			// If SIGTERM fails, fall back to Kill
			if killErr := cmd.Process.Kill(); killErr != nil {
				fmt.Fprintf(w, "Error stopping server: %v", killErr)
				return
			}
		}
	}

	// Wait for the process to exit
	if waitErr := cmd.Wait(); waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") {
		// Only report actual errors, not just non-zero exit codes
		fmt.Fprintf(w, "Error during server shutdown: %v", waitErr)
		return
	}

	cmd = nil
	fmt.Fprintf(w, "Server stopped.")
}
