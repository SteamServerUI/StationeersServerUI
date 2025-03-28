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

type JsonConfig struct {
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
	Debug                   bool   `json:"Debug,omitempty"` //Optional, default false
}

var (
	// Discord-related settings
	DiscordToken            string
	DiscordSession          *discordgo.Session
	IsDiscordEnabled        bool
	ControlChannelID        string
	StatusChannelID         string
	LogChannelID            string
	ErrorChannelID          string
	ConnectionListChannelID string
	SaveChannelID           string
	ControlPanelChannelID   string

	// Server configuration
	ServerName       string
	ServerMaxPlayers string
	ServerPassword   string
	ServerAuthSecret string
	AdminPassword    string
	GamePort         string
	UpdatePort       string
	LocalIpAddress   string
	ServerVisible    bool
	UseSteamP2P      bool

	// File paths and constants
	BlackListFilePath string
	SaveFileName      string
	ExePath           string
	TLSCertPath       = "./UIMod/cert.pem"
	TLSKeyPath        = "./UIMod/key.pem"
	ConfigPath        = "./UIMod/config.json"
	GameServerAppID   = "600760" // Steam App ID for Stationeers Dedicated Server

	// Runtime settings
	SaveInterval     string
	AdditionalParams string
	AutoPauseServer  bool
	UPNPEnabled      bool
	AutoSave         bool
	StartLocalHost   bool
	TLSEnabled       bool
	IsDebugMode      bool
	IsFirstTimeSetup bool

	// Logging and buffers
	LogMessageBuffer  string
	MaxBufferSize     = 1000
	BufferFlushTicker *time.Ticker

	// Player tracking
	ConnectedPlayers          = make(map[string]string) // SteamID -> Username
	ConnectedPlayersMessageID string

	// Message IDs
	ControlMessageID       string
	ExceptionMessageID     string
	BackupRestoreMessageID string

	// Versioning
	Version = "4.0.13"
	Branch     = "nightly"
	GameBranch string
)

func LoadConfig() (*JsonConfig, error) {

	file, err := os.Open(ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jsonconfig JsonConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonconfig)
	if err != nil {
		return nil, err
	}
	// Set the default executable path if not specified in the config file
	if jsonconfig.ExePath == "" {
		if runtime.GOOS == "windows" {
			jsonconfig.ExePath = "./rocketstation_DedicatedServer.exe"
		} else {
			jsonconfig.ExePath = "./rocketstation_DedicatedServer.x86_64"
		}
	}

	DiscordToken = jsonconfig.DiscordToken
	ControlChannelID = jsonconfig.ControlChannelID
	StatusChannelID = jsonconfig.StatusChannelID
	LogChannelID = jsonconfig.LogChannelID
	ConnectionListChannelID = jsonconfig.ConnectionListChannelID
	SaveChannelID = jsonconfig.SaveChannelID
	BlackListFilePath = jsonconfig.BlackListFilePath
	ControlPanelChannelID = jsonconfig.ControlPanelChannelID
	IsDiscordEnabled = jsonconfig.IsDiscordEnabled
	ErrorChannelID = jsonconfig.ErrorChannelID
	GameBranch = jsonconfig.GameBranch
	ServerName = jsonconfig.ServerName
	SaveFileName = jsonconfig.SaveFileName
	ServerMaxPlayers = jsonconfig.ServerMaxPlayers
	ServerPassword = jsonconfig.ServerPassword
	ServerAuthSecret = jsonconfig.ServerAuthSecret
	AdminPassword = jsonconfig.AdminPassword
	GamePort = jsonconfig.GamePort
	UpdatePort = jsonconfig.UpdatePort
	UPNPEnabled = jsonconfig.UPNPEnabled
	AutoSave = jsonconfig.AutoSave
	SaveInterval = jsonconfig.SaveInterval
	AutoPauseServer = jsonconfig.AutoPauseServer
	LocalIpAddress = jsonconfig.LocalIpAddress
	StartLocalHost = jsonconfig.StartLocalHost
	ServerVisible = jsonconfig.ServerVisible
	UseSteamP2P = jsonconfig.UseSteamP2P
	ExePath = jsonconfig.ExePath
	AdditionalParams = jsonconfig.AdditionalParams
	IsDebugMode = jsonconfig.Debug

	if jsonconfig.Debug {
		fmt.Println("DiscordToken:", DiscordToken)
		fmt.Println("ControlChannelID:", ControlChannelID)
		fmt.Println("StatusChannelID:", StatusChannelID)
		fmt.Println("LogChannelID:", LogChannelID)
		fmt.Println("ConnectionListChannelID:", ConnectionListChannelID)
		fmt.Println("SaveChannelID:", SaveChannelID)
		fmt.Println("BlackListFilePath:", BlackListFilePath)
		fmt.Println("ControlPanelChannelID:", ControlPanelChannelID)
		fmt.Println("IsDiscordEnabled:", IsDiscordEnabled)
		fmt.Println("ErrorChannelID:", ErrorChannelID)
		fmt.Println("GameBranch:", GameBranch)
		fmt.Println("ServerName:", ServerName)
		fmt.Println("SaveFileName:", SaveFileName)
		fmt.Println("ServerMaxPlayers:", ServerMaxPlayers)
		fmt.Println("ServerPassword:", ServerPassword)
		fmt.Println("ServerAuthSecret:", ServerAuthSecret)
		fmt.Println("AdminPassword:", AdminPassword)
		fmt.Println("GamePort:", GamePort)
		fmt.Println("UpdatePort:", UpdatePort)
		fmt.Println("UPNPEnabled:", UPNPEnabled)
		fmt.Println("AutoSave:", AutoSave)
		fmt.Println("SaveInterval:", SaveInterval)
		fmt.Println("AutoPauseServer:", AutoPauseServer)
		fmt.Println("LocalIpAddress:", LocalIpAddress)
		fmt.Println("StartLocalHost:", StartLocalHost)
		fmt.Println("ServerVisible:", ServerVisible)
		fmt.Println("UseSteamP2P:", UseSteamP2P)
		fmt.Println("ExePath:", ExePath)
		fmt.Println("AdditionalParams:", AdditionalParams)
		fmt.Println("Version:", Version)
		fmt.Println("Branch:", Branch)
		fmt.Println("GameServerAppID:", GameServerAppID)
		fmt.Println("DiscordSession:", DiscordSession)
		fmt.Println("LogMessageBuffer:", LogMessageBuffer)
		fmt.Println("MaxBufferSize:", MaxBufferSize)
		fmt.Println("BufferFlushTicker:", BufferFlushTicker)
		fmt.Println("ConnectedPlayers:", ConnectedPlayers)
		fmt.Println("ConnectedPlayersMessageID:", ConnectedPlayersMessageID)
	}

	return &jsonconfig, nil
}
