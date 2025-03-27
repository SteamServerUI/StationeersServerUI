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
	IsDebugMode               bool
	GameBranch                string
	Version                   = "4.0.5"
	Branch                    = "nightly"
	GameServerAppID           = "600760" // Steam App ID for Stationeers Dedicated Server
	ConfigPath                = "./UIMod/config.json"
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
