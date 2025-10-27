package advertiser

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/detectionmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
)

var StationeersAdvertisementEndpoint = config.GetStationeersServerPingEndpoint()

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

func getIpFromAdvertiserOverride(address string) (string, error) {
	// If the address is "auto", we need to check our public IPv4 via ipify
	if address == "auto" {
		resp, err := http.Get("https://api4.ipify.org")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return buf.String(), nil
	}
	// If the address is an IP quad, return it as is
	if ip := net.ParseIP(address); ip != nil {
		if ip.To4() != nil {
			return address, nil
		} else if ip.To16() != nil {
			return "", errors.New("IPv6 addresses are not supported for advertiser override")
		}
	}
	// If the address is a DNS name, resolve it
	ips, err := net.LookupIP(address)
	if err != nil {
		return "", err
	}
	// Return the first resolved IPv4 address
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}
	// If the address is invalid, return an error
	return "", errors.New("unable to resolve IP from advertiser override")
}

func StartAdvertiser() {
	if config.GetServerVisible() {
		logger.Advertiser.Warn("Server advertisement is enabled. Disable it in the config and restart SSUI to use manual advertisement. Skipping for now...")
		return
	}
	go func() {
		sessionId := -1
		for {
			// Only advertise if we are running
			if gamemgr.InternalIsServerRunning() {
				// Get max players
				maxplayers, err := strconv.Atoi(config.GetServerMaxPlayers())
				if err != nil {
					logger.Advertiser.Errorf("ServerAdvertiser failed to convert max players number to int: %s", config.GetServerMaxPlayers())
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
				// Get IP address
				ipAddress, err := getIpFromAdvertiserOverride(config.GetAdvertiserOverride())
				if err != nil {
					logger.Advertiser.Errorf("ServerAdvertiser failed to get IP address from config value '%s': %v", config.GetAdvertiserOverride(), err)
					return
				}
				adMessage := ServerAdMessage{
					SessionId:  sessionId,
					Name:       config.GetServerName(),
					Password:   config.GetServerPassword() != "",
					Version:    config.GetExtractedGameVersion(),
					Address:    ipAddress,
					Port:       config.GetGamePort(),
					Players:    players,
					MaxPlayers: maxplayers,
					Type:       platform,
				}
				body, err := json.Marshal(adMessage)
				if err != nil {
					logger.Advertiser.Errorf("ServerAdvertiser failed to Serialize to JSON from native Go struct type: %v", err)
					return
				}
				// Send advertisement
				resp, err := http.Post(StationeersAdvertisementEndpoint, "application/json", bytes.NewBuffer(body))
				// Check for errors
				if err != nil {
					logger.Advertiser.Errorf("ServerAdvertiser failed to send request: %v", err)
					return
				}
				defer resp.Body.Close()
				// Check the status
				if resp.StatusCode != 200 {
					logger.Advertiser.Warnf("ServerAdvertiser received non-200 status: %d", resp.StatusCode)
				}
				// Read the response and update our sessionId if needed
				adResponse := ServerAdResponse{}
				err = json.NewDecoder(resp.Body).Decode(&adResponse)
				if err != nil {
					logger.Advertiser.Errorf("Failed to decode response body: %v", err)
					return
				}
				if adResponse.Status != "Success" {
					logger.Advertiser.Warnf("ServerAdvertiser received unexpected status: %s", adResponse.Status)
				}
				sessionId = adResponse.SessionId
			} else {
				// Reset sessionid for the next run
				sessionId = -1
			}
			// Sleep for 30 seconds to follow the standard advertisement timer
			time.Sleep(30 * time.Second)
		}
	}()
}
