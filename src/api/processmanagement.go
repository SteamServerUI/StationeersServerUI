// processmanagement.go
package api

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func StartServer(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		fmt.Fprintf(w, "Server is already running.")
		return
	}

	config, err := loadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading config: %v", err), http.StatusInternalServerError)
		return
	}

	// Fix: Properly construct the parameters array
	alwaysNeededParams := []string{"-batchmode", "-nographics"}
	args := alwaysNeededParams
	if runtime.GOOS == "windows" {
		args = append(alwaysNeededParams, "-LOAD", config.SaveFileName, "-settings")
	} else if runtime.GOOS == "linux" {
		args = append(alwaysNeededParams, "-LOAD", "-logfile \"./debug.log\"", config.SaveFileName, "-settings")
	}
	args = append(args, strings.Split(config.Server.Settings, " ")...)

	// Create command based on OS
	if runtime.GOOS == "windows" {
		cmd = exec.Command(config.Server.ExePath, args...)
	} else {
		// Linux-specific setup with output redirection
		cmd = exec.Command(config.Server.ExePath, args...)

		// Create or truncate output file
		outFile, err := os.Create("./proc.out")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating output file: %v", err), http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Redirect both stdout and stderr to the file
		cmd.Stdout = outFile
		cmd.Stderr = outFile
	}

	exePath := colorGreen + colorBold + config.Server.ExePath + colorReset
	fmt.Printf("\n%s%s=== GAMESERVER STARTING ===%s\n", colorCyan, colorBold, colorReset)
	fmt.Printf("• Executable: %s\n", exePath)
	fmt.Printf("• Parameters: ")

	// Fix: Print parameters with proper spacing
	for i, arg := range args {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%s%s%s", colorYellow, arg, colorReset)
	}
	fmt.Printf("\n\n")

	// For Windows only: keep pipe reading functionality
	if runtime.GOOS == "windows" {
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

		// Start reading stdout and stderr
		go readPipe(stdout)
		go readPipe(stderr)
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(w, "Error starting server: %v", err)
		return
	}

	fmt.Fprintf(w, "Server started. PID: %d", cmd.Process.Pid)
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

func GetOutput(w http.ResponseWriter, r *http.Request) {
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
