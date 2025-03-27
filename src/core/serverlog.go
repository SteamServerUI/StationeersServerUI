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
	"time"
)

// handler for the /console endpoint
func GetLogOutput(w http.ResponseWriter, r *http.Request) {
	ssestream.StartConsoleStream()(w, r)
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

// tailLogFile uses tail to read the log file because using the gameserver's output in pipes to read the serverlog doesn't work on Linux with the Stationeers gameserver.
// I didn't manage to implement proper file tailing (tail behavior) here in go, so I opted to just use the actual tail.. This is a workaround for a workaround.

func tailLogFile(logFilePath string) {
	//if we somehow end up running THIS on windows, hard error and shutdown as the whole point of this software is to read the logs and do stuff with them.
	if runtime.GOOS == "windows" {
		fmt.Println("[MAJOR ISSUE DETECTED] Windows detected while trying to read log files the Linux way, skipping. You might wanna check your environment, as this should not happen.")
		fmt.Println("[MAJOR ISSUE DETECTED] Shutting down...")
		ssestream.BroadcastConsoleOutput("[MAJOR ISSUE DETECTED] Windows detected while trying to read log files the Linux way, skipping. You might wanna check your environment, as this should not happen.")
		ssestream.BroadcastConsoleOutput("[MAJOR ISSUE DETECTED] Shutting down...")
		os.Exit(1)
	}

	// Wait and retry until the log file exists
	for i := range 10 { // Retry up to 10 times
		if _, err := os.Stat(logFilePath); err == nil {
			break // File exists, proceed
		}
		if config.IsDebugMode {
			fmt.Printf("Log file %s not found, retrying in 1s (%d/10)\n", logFilePath, i+1)
		}
		time.Sleep(1 * time.Second)
	}

	// If file still doesn't exist, give up and report
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		if config.IsDebugMode {
			fmt.Printf("Log file %s still not found after retries\n", logFilePath)
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Log file %s not found after retries", logFilePath))
		return
	}

	// Start tail -F (robust against rotation)
	cmd := exec.Command("tail", "-F", logFilePath)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		if config.IsDebugMode {
			fmt.Printf("Error creating stdout pipe for tail: %v\n", err)
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error starting tail -F: %v", err))
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		if config.IsDebugMode {
			fmt.Printf("Error starting tail -F: %v\n", err)
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error starting tail -F: %v", err))
		return
	}

	// Clean up when done
	defer func() {
		cmd.Process.Kill() // Kill tail when logDone triggers
		if err := cmd.Wait(); err != nil && config.IsDebugMode {
			fmt.Printf("Tail process exited with: %v\n", err)
		}
	}()

	scanner := bufio.NewScanner(pipe)
	if config.IsDebugMode {
		fmt.Println("Started tailing log file with tail -F")
	}

	// Goroutine to read and broadcast tail output
	go func() {
		defer pipe.Close() // Close pipe when goroutine exits
		for scanner.Scan() {
			output := scanner.Text()
			//if config.IsDebugMode {
			//	fmt.Println("DEBUG: Read from tail -F:", output)
			//}
			ssestream.BroadcastConsoleOutput(output)
		}
		if err := scanner.Err(); err != nil {
			if config.IsDebugMode {
				fmt.Printf("Error reading tail -F output: %v\n", err)
			}
			ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reading tail -F output: %v", err))
		}
	}()

	// Wait for logDone signal to stop
	<-logDone
	if config.IsDebugMode {
		fmt.Println("Received logDone signal, stopping tail -F")
	}
}
