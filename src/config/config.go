// config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	DiscordToken            string `json:"discordToken"`
	ControlChannelID        string `json:"controlChannelID"`
	StatusChannelID         string `json:"statusChannelID"`
	ConnectionListChannelID string `json:"connectionListChannelID"`
	LogChannelID            string `json:"logChannelID"`
	SaveChannelID           string `json:"saveChannelID"`
	ControlPanelChannelID   string `json:"controlPanelChannelID"`
	BlackListFilePath       string `json:"blackListFilePath"`
	IsDiscordEnabled        bool   `json:"isDiscordEnabled"`
	ErrorChannelID          string `json:"errorChannelID"`
	GameBranch              string `json:"gameBranch"`
	WorldType               string `json:"worldType"`
}

var (
	DiscordToken              string
	ControlChannelID          string
	StatusChannelID           string
	LogChannelID              string
	ErrorChannelID            string
	ConnectionListChannelID   string
	SaveChannelID             string
	BlackListFilePath         string
	DiscordSession            *discordgo.Session
	LogMessageBuffer          string
	MaxBufferSize             = 1000
	BufferFlushTicker         *time.Ticker
	ConnectedPlayers          = make(map[string]string) // SteamID -> Username
	ConnectedPlayersMessageID string
	ControlMessageID          string
	ExceptionMessageID        string
	BackupRestoreMessageID    string
	ControlPanelChannelID     string
	IsDiscordEnabled          bool
	IsFirstTimeSetup          bool
	GameBranch                string
	Version                   = "3.0.1"
	Branch                    = "release"
	GameServerAppID           = "600760" // Steam App ID for Stationeers Dedicated Server
)

func (c *Config) ValidateWorldType() error {
	allowedWorldTypes := map[string]bool{
		"moon":    true,
		"mars":    true,
		"europa":  true,
		"europa2": true,
		"mimas":   true,
		"vulcan":  true,
		"vulcan2": true,
		"space":   true,
		"loulan":  true,
		"venus":   true,
		"":        true, //Allow empty string for no world type
	}
	if !allowedWorldTypes[c.WorldType] {
		return fmt.Errorf("invalid WorldType: %s. Allowed values are: moon, mars, europa, europa2, mimas, vulcan, vulcan2, space, loulan, venus", c.WorldType)
	}
	return nil
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	if err := config.ValidateWorldType(); err != nil {
		return nil, err
	}
	//print all the values to console
	//fmt.Println("DiscordToken:", config.DiscordToken)
	//fmt.Println("ControlChannelID:", config.ControlChannelID)
	//fmt.Println("StatusChannelID:", config.StatusChannelID)
	//fmt.Println("ConnectionListChannelID:", config.ConnectionListChannelID)
	//fmt.Println("LogChannelID:", config.LogChannelID)
	//fmt.Println("SaveChannelID:", config.SaveChannelID)
	//fmt.Println("BlackListFilePath:", config.BlackListFilePath)
	//fmt.Println("IsDiscordEnabled:", config.IsDiscordEnabled)
	//fmt.Println("ErrorChannelID:", config.ErrorChannelID)
	//fmt.Println("IsFirstTimeSetup:", IsFirstTimeSetup)
	DiscordToken = config.DiscordToken
	ControlChannelID = config.ControlChannelID
	StatusChannelID = config.StatusChannelID
	LogChannelID = config.LogChannelID
	ConnectionListChannelID = config.ConnectionListChannelID
	SaveChannelID = config.SaveChannelID
	BlackListFilePath = config.BlackListFilePath
	ControlPanelChannelID = config.ControlPanelChannelID
	IsDiscordEnabled = config.IsDiscordEnabled
	ErrorChannelID = config.ErrorChannelID
	GameBranch = config.GameBranch
	return &config, nil
}
