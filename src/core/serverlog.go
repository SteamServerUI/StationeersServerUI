package core

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/ssestream"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
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

// tailLogFile implements a tail -f-like behavior using fsnotify for linux because using the pipes to read the serverlog doesn't work on Linux with the stationeers gameserver
func tailLogFile(logFilePath string) {
	// Wait briefly to ensure the file exists after server start
	time.Sleep(1 * time.Second)

	// Open the log file
	file, err := os.Open(logFilePath)
	if err != nil {
		if config.IsDebugMode {
			fmt.Printf("Error opening log file: %v\n", err)
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error opening log file: %v", err))
		return
	}
	defer file.Close()

	// Create an fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		if config.IsDebugMode {
			fmt.Printf("Error creating watcher: %v\n", err)
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error spawning log file watcher: %v", err))
		return
	}
	defer watcher.Close()

	// Add the log file to the watcher
	err = watcher.Add(logFilePath)
	if err != nil {
		if config.IsDebugMode {
			fmt.Printf("Error adding log file to watcher: %v\n", err)
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error adding log file to watcher: %v", err))
		return
	}

	// Seek to the end of the file initially to start tailing from the current point
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
		fmt.Println("Started tailing log file with fsnotify")
	}

	// Goroutine to handle file events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					// New content has been written to the file
					for scanner.Scan() {
						output := scanner.Text()
						if config.IsDebugMode {
							fmt.Println("DEBUG: Read from log file:", output)
						}
						ssestream.BroadcastConsoleOutput(output)
					}
					if err := scanner.Err(); err != nil {
						if config.IsDebugMode {
							fmt.Printf("Error reading log file: %v\n", err)
						}
						ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reading log file: %v", err))
						return
					}
				}
				// Handle file truncation or rotation (optional)
				if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					if config.IsDebugMode {
						fmt.Println("Log file removed or renamed, attempting to reopen")
					}
					// Reopen the file and re-add to watcher
					file.Close()
					file, err = os.Open(logFilePath)
					if err != nil {
						ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reopening log file: %v", err))
						return
					}
					scanner = bufio.NewScanner(file)
					watcher.Add(logFilePath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				if config.IsDebugMode {
					fmt.Printf("Watcher error: %v\n", err)
				}
				ssestream.BroadcastConsoleOutput(fmt.Sprintf("Watcher error: %v", err))
			case <-logDone:
				return
			}
		}
	}()

	// Keep the function running until we either get a signal to stop or the server stops
	<-logDone
}
