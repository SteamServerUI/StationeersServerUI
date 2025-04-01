// logstream.go
package detection

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/discord"
	"StationeersServerUI/src/ssestream"
	"fmt"
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
		fmt.Println(string(colorGreen), "Connected to internal log stream.", string(colorReset))
		for logMessage := range logChan {
			if config.IsDiscordEnabled {
				discord.PassLogStreamToDiscordLogBuffer(logMessage)
			}
			ProcessLog(detector, logMessage)
		}
	}()
}
