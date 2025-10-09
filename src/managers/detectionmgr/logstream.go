// logstream.go
package detectionmgr

import (
	"os"
	"path/filepath"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/ssestream"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/discordbot"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

/*
Real-time Log Processing Pipeline
- Bridges internal SSE stream to detection system
- Performs log enrichment and distribution:
  - Adds logs to the Discord integrations log buffer if enabled
  - Feeds messages to Detector
  - Optionally logs messages to daily log files
*/

// StartLogStream starts processing logs directly from the internal SSE manager
func StreamLogs(detector *Detector) {
	logChan := ssestream.ConsoleStreamManager.AddInternalSubscriber()

	go func() {
		logger.Detection.Debug("Connected to internal log stream.")
		for logMessage := range logChan {
			if config.GetIsDiscordEnabled() {
				discordbot.PassLogStreamToDiscordLogBuffer(logMessage)
			}
			ProcessLog(detector, logMessage)
		}
	}()

	if config.LogServerOutputToFile {
		logChan2 := ssestream.ConsoleStreamManager.AddInternalSubscriber()
		f, _ := OpenLogFile(nil)
		go func() {
			logger.Core.Debug("LogOutput: Connected to internal log stream.")
			for logMessage := range logChan2 {
				f, _ = OpenLogFile(f)
				SaveGameLogLine(f, logMessage)
			}
		}()
	}
}

func OpenLogFile(prev *os.File) (*os.File, error) {
	logFile := "server_logs/" + time.Now().Format("2006-01-02") + ".log"

	if prev != nil {
		if prev.Name() != logFile {
			prev.Close()
		} else {
			return prev, nil
		}
	}

	logger.Core.Info("Opening log file: " + logFile)

	if err := os.MkdirAll(filepath.Dir(logFile), os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Core.Error("Can't open server.log: " + err.Error())
		return nil, err
	}
	return f, nil
}

func SaveGameLogLine(output *os.File, logLine string) {
	if output == nil {
		return
	}
	if _, err := output.WriteString(logLine + "\n"); err != nil {
		logger.Core.Error("Can't write log message to server.log: " + err.Error())
	}
}
