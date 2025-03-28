package detection

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/discord"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/r3labs/sse"
)

// This setup is a relict of Version 1.0, when the SSE stream was not coming from the same server as the UI.
// This should ideally be replaced, but there currently isn't a good way to tie into log reading for Linux and Windows simultaneously, hence requiring seperade logic.

func StartLogStream(detector *Detector) {
	client := sse.NewClient("https://localhost:443/console")
	client.Headers["Content-Type"] = "text/event-stream"
	client.Headers["Connection"] = "keep-alive"
	client.Headers["Cache-Control"] = "no-cache"
	client.Connection.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	retryDelay := 5 * time.Second

	go func() {
		for {
			fmt.Println(string(colorYellow), "Attempting to connect to SSE stream...", string(colorReset))

			connected := false
			err := client.SubscribeRaw(func(msg *sse.Event) {
				if len(msg.Data) > 0 {
					// Log connection only on the first message
					if !connected {
						fmt.Println(string(colorGreen), "Connected to SSE stream.", string(colorReset))
						connected = true
					}

					logMessage := string(msg.Data)
					if config.IsDiscordEnabled {
						discord.AddToLogBuffer(logMessage)
					}
					ProcessLog(detector, logMessage)
				}
			})

			if err != nil {
				fmt.Println(string(colorYellow), "SSE stream not available yet, retrying in 5 seconds...", string(colorReset))
				fmt.Printf(string(colorRed)+"SubscribeRaw error: %v\n"+string(colorReset), err)
				time.Sleep(retryDelay)
				continue
			}
			// If SubscribeRaw returns without error, break the loop
			break
		}
	}()
}
