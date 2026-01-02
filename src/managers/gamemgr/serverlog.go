package gamemgr

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/ssestream"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/google/uuid"
)

const defaultLogFolderMode = 0755
const defaultLogFileMode = 0644

// readPipe for Windows
func readPipe(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	logger.Core.Debug("Started reading pipe")
	for scanner.Scan() {
		output := scanner.Text()
		ssestream.BroadcastConsoleOutput(output)
	}
	if err := scanner.Err(); err != nil {
		logger.Core.Debug("Pipe error: " + err.Error())
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reading pipe: %v", err))
	}
	logger.Core.Debug("Pipe closed")
}

// tailLogFile uses tail to read the log file because using the gameserver's output in pipes to read the serverlog doesn't work on Linux with the Stationeers gameserver.
// I didn't manage to implement proper file tailing (tail behavior) here in go, so I opted to just use the actual tail.. This is a workaround for a workaround.

func tailLogFile(logFilePath string) {
	var cmd *exec.Cmd

	// Wait and retry until the log file exists
	for i := range 10 { // Retry up to 10 times
		if _, err := os.Stat(logFilePath); err == nil {
			break // File exists, proceed
		}
		logger.Core.Debug("Log file " + logFilePath + " not found, retrying in 1s (" + strconv.Itoa(i+1) + "/10)")
		time.Sleep(1 * time.Second)
	}

	// If file still doesn't exist, give up and report
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		logger.Core.Debug("Log file " + logFilePath + " still not found after retries")
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Log file %s not found after retries", logFilePath))
		return
	}

	// Choose command based on operating system
	if runtime.GOOS == "windows" {
		// Use PowerShell's Get-Content -Wait for Windows
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Get-Content -Path %q -Wait", logFilePath))
	} else {
		// Use tail -F for Linux/Unix
		cmd = exec.Command("tail", "-F", logFilePath)
	}

	// Create a pipe for command output
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Core.Debug("Error creating stdout pipe for tail command: " + err.Error())
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error starting tail command: %v", err))
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		logger.Core.Debug("Error starting tail: " + err.Error())
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error starting tail: %v", err))
		return
	}

	// Clean up when done
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill() // Kill the command process when logDone triggers
			if err := cmd.Wait(); err != nil {
				logger.Core.Debug("Tail process exited with: " + err.Error())
			}
		}
	}()

	scanner := bufio.NewScanner(pipe)
	logger.Core.Debug("Started tailing a log file")

	// Goroutine to read and broadcast tail output
	go func() {
		defer pipe.Close() // Close pipe when goroutine exits
		for scanner.Scan() {
			output := scanner.Text()
			ssestream.BroadcastConsoleOutput(output)
		}
		if err := scanner.Err(); err != nil {

			logger.Core.Debug("Error reading tail output: " + err.Error())

			ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reading tail output: %v", err))
		}
	}()

	// Wait for logDone signal to stop
	<-logDone

	logger.Core.Debug("Received logDone signal, stopping tail -F")

}

func logToFile(message string) {
	if config.GetCreateGameServerLogFile() {
		logFileFolder := config.GetLogFolder()
		// Check if the log folder exists, if not create it
		_, err := os.Stat(logFileFolder)
		if os.IsNotExist(err) {
			err := os.Mkdir(logFileFolder, defaultLogFolderMode)
			if err != nil {
				logger.Core.Error("Error creating log folder: " + err.Error())
				return
			}
		} else if err != nil {
			logger.Core.Error("Error checking log folder: " + err.Error())
			return
		}

		if GameServerUUID != uuid.Nil {
			logFileName := fmt.Sprintf("serverlog_%s_%s.log", time.Now().Format("200601021504"), GameServerUUID.String())
			logFilePath := path.Join(logFileFolder, logFileName)

			// append the log file to the log folder
			file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, defaultLogFileMode)
			if err != nil {
				logger.Core.Error("Error opening log file: " + err.Error())
				return
			}
			defer file.Close()

			_, err = file.WriteString(message + "\n")
			if err != nil {
				logger.Core.Error("Error writing to log file: " + err.Error())
				return
			}
		} else {
			logger.Core.Error("Game Server UUID not set, cannot log to file")
			return
		}
	}
}
