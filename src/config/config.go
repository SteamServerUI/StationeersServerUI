// config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
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
	DiscordCharBufferSize   int    `json:"DiscordCharBufferSize"`
	BlackListFilePath       string `json:"blackListFilePath"`
	IsDiscordEnabled        bool   `json:"isDiscordEnabled"`
	ErrorChannelID          string `json:"errorChannelID"`
	GameBranch              string `json:"gameBranch"`
	ServerName              string `json:"ServerName"`
	SaveInfo                string `json:"SaveInfo"`
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
	Username                string `json:"Username,omitempty"`
	Password                string `json:"Password,omitempty"`
	JwtKey                  string `json:"JwtKey,omitempty"`
	AuthTokenLifetime       int    `json:"AuthTokenLifetime,omitempty"`
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
	SaveInfo          string // Save folder name, optionally with backup name ("WorldName BackupName")
	BackupWorldName   string // Only Backup world name
	WorldName         string // Only World name
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
	LogMessageBuffer      string
	DiscordCharBufferSize int
	SSEMessageBufferSize  = 2000
	MaxSSEConnections     = 20
	BufferFlushTicker     *time.Ticker

	// Player tracking
	ConnectedPlayers          = make(map[string]string) // SteamID -> Username
	ConnectedPlayersMessageID string

	// Message IDs
	ControlMessageID       string
	ExceptionMessageID     string
	BackupRestoreMessageID string

	// Authentication
	Username          string
	Password          string
	JwtKey            string
	AuthTokenLifetime int // In minutes, e.g., 1440 (24h)

	// Versioning
	Version = "4.1.4"
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

	if jsonconfig.DiscordCharBufferSize <= 0 {
		jsonconfig.DiscordCharBufferSize = 1000 // Default to 1000 characters
	}

	if jsonconfig.GameBranch == "" {
		jsonconfig.GameBranch = "public" //default to public release of Stationeers if no value is set from the config file
	}

	// Get secrets from env vars, json or defaults
	GetSecretsFromEnv(jsonconfig)

	if jsonconfig.SaveInfo == "" {
		jsonconfig.SaveInfo = "Moon Moon"
	}

	//set BackupWorldName and WorldName from SaveInfo, wich is in the format "WorldName BackupName" for argument handling.
	parts := strings.Split(jsonconfig.SaveInfo, " ")
	if len(parts) > 0 {
		WorldName = parts[0]
	}
	if len(parts) > 1 {
		BackupWorldName = parts[1]
	}

	//set the rest of the config variables from the json config if they are available
	SaveInfo = jsonconfig.SaveInfo
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
	DiscordCharBufferSize = jsonconfig.DiscordCharBufferSize

	if jsonconfig.Debug {
		fmt.Println("----DISCORD CONFIG VARS----")
		fmt.Println("BlackListFilePath:", BlackListFilePath)
		fmt.Println("ConnectedPlayersMessageID:", ConnectedPlayersMessageID)
		fmt.Println("ConnectionListChannelID:", ConnectionListChannelID)
		fmt.Println("ControlChannelID:", ControlChannelID)
		fmt.Println("ControlPanelChannelID:", ControlPanelChannelID)
		fmt.Println("DiscordCharBufferSize:", DiscordCharBufferSize)
		fmt.Println("DiscordToken:", DiscordToken)
		fmt.Println("ErrorChannelID:", ErrorChannelID)
		fmt.Println("IsDiscordEnabled:", IsDiscordEnabled)
		fmt.Println("LogChannelID:", LogChannelID)
		fmt.Println("LogMessageBuffer:", LogMessageBuffer)
		fmt.Println("SaveChannelID:", SaveChannelID)
		fmt.Println("StatusChannelID:", StatusChannelID)
		fmt.Println("----GAMESERVER CONFIG VARS----")
		fmt.Println("AdditionalParams:", AdditionalParams)
		fmt.Println("AdminPassword:", AdminPassword)
		fmt.Println("AutoPauseServer:", AutoPauseServer)
		fmt.Println("AutoSave:", AutoSave)
		fmt.Println("BackupWorldName:", BackupWorldName)
		fmt.Println("ExePath:", ExePath)
		fmt.Println("GameBranch:", GameBranch)
		fmt.Println("GamePort:", GamePort)
		fmt.Println("LocalIpAddress:", LocalIpAddress)
		fmt.Println("SaveInfo:", SaveInfo)
		fmt.Println("SaveInterval:", SaveInterval)
		fmt.Println("ServerAuthSecret:", ServerAuthSecret)
		fmt.Println("ServerMaxPlayers:", ServerMaxPlayers)
		fmt.Println("ServerName:", ServerName)
		fmt.Println("ServerPassword:", ServerPassword)
		fmt.Println("ServerVisible:", ServerVisible)
		fmt.Println("StartLocalHost:", StartLocalHost)
		fmt.Println("UpdatePort:", UpdatePort)
		fmt.Println("UPNPEnabled:", UPNPEnabled)
		fmt.Println("UseSteamP2P:", UseSteamP2P)
		fmt.Println("WorldName:", WorldName)
		fmt.Println("----AUTHENTICATION CONFIG VARS----")
		fmt.Println("AuthTokenLifetime:", AuthTokenLifetime)
		fmt.Println("JwtKey:", JwtKey)
		fmt.Println("Password:", Password)
		fmt.Println("Username:", Username)
		fmt.Println("----MISC CONFIG VARS----")
		fmt.Println("Branch:", Branch)
		fmt.Println("GameServerAppID:", GameServerAppID)
		fmt.Println("Version:", Version)
	}

	return &jsonconfig, nil
}
