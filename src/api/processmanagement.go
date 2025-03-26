// processmanagement.go
package api

import (
	"StationeersServerUI/src/config"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var cmd *exec.Cmd
var mu sync.Mutex
var clients []chan string
var clientsMu sync.Mutex

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

type Arg struct {
	Flag      string
	Value     string
	Condition func() bool
}

// Define argument order with clearer initialization
var argOrder = []Arg{
	{Flag: "-batchmode"},
	{Flag: "-nographics"},
	{Flag: "-LOAD", Value: config.SaveFileName},
	{Flag: "-logFile", Value: "./debug.log", Condition: func() bool { return runtime.GOOS == "linux" }}, // Attach a logfile on Linux, since piped output is not available
	{Flag: "-settings"},
	{Flag: "StartLocalHost", Value: strconv.FormatBool(config.StartLocalHost)},
	{Flag: "ServerVisible", Value: strconv.FormatBool(config.ServerVisible)},
	{Flag: "GamePort", Value: config.GamePort},
	{Flag: "UpdatePort", Value: config.UpdatePort},
	{Flag: "AutoSave", Value: strconv.FormatBool(config.AutoSave)},
	{Flag: "SaveInterval", Value: config.SaveInterval},
	{Flag: "ServerMaxPlayers", Value: config.ServerMaxPlayers},
	{Flag: "ServerName", Value: config.ServerName},
	{Flag: "ServerPassword", Value: config.ServerPassword, Condition: func() bool { return config.ServerPassword != "" }},
	{Flag: "ServerAuthSecret", Value: config.ServerAuthSecret, Condition: func() bool { return config.ServerAuthSecret != "" }},
	{Flag: "AdminPassword", Value: config.AdminPassword, Condition: func() bool { return config.AdminPassword != "" }},
	{Flag: "UPNPEnabled", Value: strconv.FormatBool(config.UPNPEnabled)},
	{Flag: "AutoPauseServer", Value: strconv.FormatBool(config.AutoPauseServer)},
	{Flag: "UseSteamP2P", Value: strconv.FormatBool(config.UseSteamP2P)},
	{Flag: "LocalIpAddress", Value: config.LocalIpAddress},
}

func buildCommandArgs() []string {
	var args []string

	for _, arg := range argOrder {
		if arg.Condition != nil && !arg.Condition() {
			continue
		}

		args = append(args, arg.Flag)
		if arg.Value != "" {
			args = append(args, arg.Value)
		}
	}

	if config.AdditionalParams != "" {
		extraArgs := strings.Fields(config.AdditionalParams)
		args = append(args, extraArgs...)
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

	args := buildCommandArgs()
	cmd = exec.Command(config.ExePath, args...)

	fmt.Printf("\n%s%s=== GAMESERVER STARTING ===%s\n", colorCyan, colorBold, colorReset)
	fmt.Printf("• Executable: %s\n", (colorGreen + colorBold + config.ExePath + colorReset))
	fmt.Printf("• Parameters: ")

	for i, arg := range args {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%s%s%s", colorYellow, arg, colorReset)
	}

	fmt.Printf("\n\n")
	// Capture stdout and stderr
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

	// Start the command
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(w, "Error starting server: %v", err)
		return
	}

	// Start reading stdout and stderr
	go readPipe(stdout)
	go readPipe(stderr)

	fmt.Fprintf(w, "Server started.")
}

func readPipe(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		output := scanner.Text()
		clientsMu.Lock()
		for _, clientChan := range clients {
			clientChan <- output
		}
		clientsMu.Unlock()
	}
	if err := scanner.Err(); err != nil {
		output := fmt.Sprintf("Error reading pipe: %v", err)
		clientsMu.Lock()
		for _, clientChan := range clients {
			clientChan <- output
		}
		clientsMu.Unlock()
	}
}

func GetLogOutput(w http.ResponseWriter, r *http.Request) {
	// Create a new channel for this client
	clientChan := make(chan string)

	// Register the client
	clientsMu.Lock()
	clients = append(clients, clientChan)
	clientsMu.Unlock()

	// Ensure the channel is removed when the client disconnects
	defer func() {
		clientsMu.Lock()
		for i, ch := range clients {
			if ch == clientChan {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		clientsMu.Unlock()
		close(clientChan)
	}()

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Write data to the client as it comes in
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for msg := range clientChan {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}
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
