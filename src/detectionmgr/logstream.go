// logstream.go
package detectionmgr

import (
	"github.com/JacksonTheMaster/StationeersServerUI/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/src/discordbot"
	"github.com/JacksonTheMaster/StationeersServerUI/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/src/ssestream"
)

/*
Real-time Log Processing Pipeline
- Bridges internal SSE stream to detection system
- Performs log enrichment and distribution:
  - Adds logs to the Discord integrations log buffer if enabled
  - Feeds messages to Detector
*/

// StartLogStream starts processing logs directly from the internal SSE manager
func StreamLogs(detector *Detector) {
	logChan := ssestream.ConsoleStreamManager.AddInternalSubscriber()

	go func() {
		logger.Detection.Info("Connected to internal log stream.")
		for logMessage := range logChan {
			if config.IsDiscordEnabled {
				discordbot.PassLogStreamToDiscordLogBuffer(logMessage)
			}
			ProcessLog(detector, logMessage)
		}
	}()
}
