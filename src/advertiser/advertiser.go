package advertiser

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/detectionmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
)

type ServerAdMessage struct {
	SessionId  int
	Name       string
	Password   bool
	Version    string
	Address    string
	Port       string
	Players    int
	MaxPlayers int
	Type       int
}

type ServerAdResponse struct {
	SessionId int
	Status    string
}

func StartAdvertiser(wg *sync.WaitGroup) {
	if config.GetServerVisible() {
		logger.Core.Warn("Server advertisement is enabled. Disable it in the config and restart SSUI to use manual advertisement. Skipping for now...")
		return
	}
	wg.Go(func() {
		sessionId := -1
		for {
			// Only advertise if we are running
			if gamemgr.InternalIsServerRunning() {
				// Get max players
				maxplayers, err := strconv.Atoi(config.GetServerMaxPlayers())
				if err != nil {
					logger.Core.Errorf("ServerAdvertiser failed to convert port number to int: %s", config.GetServerMaxPlayers())
					return
				}
				// Get connected players
				detector := detectionmgr.GetDetector()
				players := len(detectionmgr.GetPlayers(detector))
				// Get platform
				platform := 0
				switch runtime.GOOS {
				case "windows":
					platform = 1
				case "linux":
					platform = 2
				}
				adMessage := ServerAdMessage{
					SessionId:  sessionId,
					Name:       config.GetServerName(),
					Password:   config.GetServerPassword() != "",
					Version:    config.GetExtractedGameVersion(),
					Address:    "127.0.0.1", //TODO - we need to pass this in from the command line
					Port:       config.GetGamePort(),
					Players:    players,
					MaxPlayers: maxplayers,
					Type:       platform,
				}
				body, err := json.Marshal(adMessage)
				if err != nil {
					logger.Core.Errorf("ServerAdvertiser failed to Serialize to JSON from native Go struct type: %v", err)
					return
				}
				// Send advertisement
				resp, err := http.Post("http://40.82.200.175:8081/Ping", "application/json", bytes.NewBuffer(body))
				// Check for errors
				if err != nil {
					logger.Core.Errorf("ServerAdvertiser failed to send request: %v", err)
					return
				}
				defer resp.Body.Close()
				// Check the status
				if resp.StatusCode != 200 {
					logger.Core.Warnf("ServerAdvertiser received non-200 status: %d", resp.StatusCode)
				}
				// Read the response and update our sessionId if needed
				adResponse := ServerAdResponse{}
				err = json.NewDecoder(resp.Body).Decode(&adResponse)
				if err != nil {
					logger.Core.Errorf("Failed to decode response body: %v", err)
					return
				}
				if adResponse.Status != "Success" {
					logger.Core.Warnf("ServerAdvertiser received unexpeted status: %s", adResponse.Status)
				}
				sessionId = adResponse.SessionId
			} else {
				// Reset sessionid for the next run
				sessionId = -1
			}
			// Sleep for 30 seconds to follow the standard advertisement timer
			time.Sleep(30 * time.Second)
		}
	})
}
