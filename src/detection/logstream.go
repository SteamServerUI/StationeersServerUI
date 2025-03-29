// logstream.go
package detection

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/discord"
	"StationeersServerUI/src/ssestream"
	"fmt"
)

// StartLogStream starts processing logs directly from the internal SSE manager
func StreamLogs(detector *Detector) {
	logChan := ssestream.ConsoleStreamManager.AddInternalSubscriber()

	go func() {
		fmt.Println(string(colorGreen), "Connected to internal log stream.", string(colorReset))
		for logMessage := range logChan {
			if config.IsDiscordEnabled {
				discord.AddToLogBuffer(logMessage)
			}
			ProcessLog(detector, logMessage)
		}
	}()
}
