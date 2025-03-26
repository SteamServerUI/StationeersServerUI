// config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
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
	ServerName              string `json:"ServerName"`
	SaveFileName            string `json:"SaveFileName"`
	ServerMaxPlayers        string `json:"ServerMaxPlayers"`
	ServerPassword          string `json:"ServerPassword"`
	ServerAuthSecret        string `json:"ServerAuthSecret"`
	AdminPassword           string `json:"AdminPassword"`
	GamePort                string `json:"GamePort"`
	UpdatePort              string `json:"UpdatePort"`
	UPNPEnabled             bool   `json:"UPNPEnabled"`
	AutoSave                bool   `json:"AutoSave"`
	SaveInterval            string `json:"SaveInterval"`
	AutoPauseServer         bool   `json:"AutoPauseServer"`
	LocalIpAddress          string `json:"LocalIpAddress"`
	StartLocalHost          bool   `json:"StartLocalHost"`
	ServerVisible           bool   `json:"ServerVisible"`
	UseSteamP2P             bool   `json:"UseSteamP2P"`
	ExePath                 string `json:"ExePath"`
	AdditionalParams        string `json:"AdditionalParams"`
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
	ServerName                string
	SaveFileName              string
	ServerMaxPlayers          string
	ServerPassword            string
	ServerAuthSecret          string
	AdminPassword             string
	GamePort                  string
	UpdatePort                string
	UPNPEnabled               bool
	AutoSave                  bool
	SaveInterval              string
	AutoPauseServer           bool
	LocalIpAddress            string
	StartLocalHost            bool
	ServerVisible             bool
	UseSteamP2P               bool
	ExePath                   string
	AdditionalParams          string
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
	Version                   = "3.0.18"
	Branch                    = "release"
	GameServerAppID           = "600760" // Steam App ID for Stationeers Dedicated Server
)

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
	// Set the default executable path if not specified in the config file
	if config.ExePath == "" {
		if runtime.GOOS == "windows" {
			config.ExePath = "./rocketstation_DedicatedServer.exe"
		} else {
			config.ExePath = "./rocketstation_DedicatedServer.x86_64"
		}
		fmt.Println("Warning: No executable path specified in config file. Defaulting to", config.ExePath)
	}

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
	ServerName = config.ServerName
	SaveFileName = config.SaveFileName
	ServerMaxPlayers = config.ServerMaxPlayers
	ServerPassword = config.ServerPassword
	ServerAuthSecret = config.ServerAuthSecret
	AdminPassword = config.AdminPassword
	GamePort = config.GamePort
	UpdatePort = config.UpdatePort
	UPNPEnabled = config.UPNPEnabled
	AutoSave = config.AutoSave
	SaveInterval = config.SaveInterval
	AutoPauseServer = config.AutoPauseServer
	LocalIpAddress = config.LocalIpAddress
	StartLocalHost = config.StartLocalHost
	ServerVisible = config.ServerVisible
	UseSteamP2P = config.UseSteamP2P
	ExePath = config.ExePath
	AdditionalParams = config.AdditionalParams
	return &config, nil
}
