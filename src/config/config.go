package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type JsonConfig struct {
	DiscordToken            string            `json:"discordToken"`
	ControlChannelID        string            `json:"controlChannelID"`
	StatusChannelID         string            `json:"statusChannelID"`
	ConnectionListChannelID string            `json:"connectionListChannelID"`
	LogChannelID            string            `json:"logChannelID"`
	SaveChannelID           string            `json:"saveChannelID"`
	ControlPanelChannelID   string            `json:"controlPanelChannelID"`
	DiscordCharBufferSize   int               `json:"DiscordCharBufferSize"`
	BlackListFilePath       string            `json:"blackListFilePath"`
	IsDiscordEnabled        *bool             `json:"isDiscordEnabled"`
	ErrorChannelID          string            `json:"errorChannelID"`
	BackupKeepLastN         int               `json:"backupKeepLastN"`
	IsCleanupEnabled        *bool             `json:"isCleanupEnabled"`
	BackupKeepDailyFor      int               `json:"backupKeepDailyFor"`
	BackupKeepWeeklyFor     int               `json:"backupKeepWeeklyFor"`
	BackupKeepMonthlyFor    int               `json:"backupKeepMonthlyFor"`
	BackupCleanupInterval   int               `json:"backupCleanupInterval"`
	BackupWaitTime          int               `json:"backupWaitTime"`
	GameBranch              string            `json:"gameBranch"`
	ServerName              string            `json:"ServerName"`
	SaveInfo                string            `json:"SaveInfo"`
	ServerMaxPlayers        string            `json:"ServerMaxPlayers"`
	ServerPassword          string            `json:"ServerPassword"`
	ServerAuthSecret        string            `json:"ServerAuthSecret"`
	AdminPassword           string            `json:"AdminPassword"`
	GamePort                string            `json:"GamePort"`
	UpdatePort              string            `json:"UpdatePort"`
	UPNPEnabled             *bool             `json:"UPNPEnabled"`
	AutoSave                *bool             `json:"AutoSave"`
	SaveInterval            string            `json:"SaveInterval"`
	AutoPauseServer         *bool             `json:"AutoPauseServer"`
	LocalIpAddress          string            `json:"LocalIpAddress"`
	StartLocalHost          *bool             `json:"StartLocalHost"`
	ServerVisible           *bool             `json:"ServerVisible"`
	UseSteamP2P             *bool             `json:"UseSteamP2P"`
	ExePath                 string            `json:"ExePath"`
	AdditionalParams        string            `json:"AdditionalParams"`
	Users                   map[string]string `json:"users"`       // Map of username to hashed password
	AuthEnabled             *bool             `json:"authEnabled"` // Toggle for enabling/disabling auth
	JwtKey                  string            `json:"JwtKey"`
	AuthTokenLifetime       int               `json:"AuthTokenLifetime"`
	Debug                   *bool             `json:"Debug"`
	CreateSSUILogFile       *bool             `json:"CreateSSUILogFile"`
	LogLevel                int               `json:"LogLevel"`
	IsUpdateEnabled         *bool             `json:"IsUpdateEnabled"`
	AllowPrereleaseUpdates  *bool             `json:"AllowPrereleaseUpdates"`
	AllowMajorUpdates       *bool             `json:"AllowMajorUpdates"`
}

type CustomDetection struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Pattern   string `json:"pattern"`
	EventType string `json:"eventType"`
	Message   string `json:"message"`
}

var (
	Version = "5.0.13"
	Branch                  = "nightly-auth"
	GameBranch              string
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
	IsCleanupEnabled        bool
	BackupKeepLastN         int
	BackupKeepDailyFor      time.Duration
	BackupKeepWeeklyFor     time.Duration
	BackupKeepMonthlyFor    time.Duration
	BackupCleanupInterval   time.Duration
	ConfiguredBackupDir     string
	ConfiguredSafeBackupDir string
	BackupWaitTime          time.Duration
	ServerName              string
	ServerMaxPlayers        string
	ServerPassword          string
	ServerAuthSecret        string
	AdminPassword           string
	GamePort                string
	UpdatePort              string
	LocalIpAddress          string
	ServerVisible           bool
	UseSteamP2P             bool
	BlackListFilePath       string
	SaveInfo                string
	BackupWorldName         string
	WorldName               string
	ExePath                 string
	TLSCertPath             = "./UIMod/cert.pem"
	TLSKeyPath              = "./UIMod/key.pem"
	ConfigPath              = "./UIMod/config.json"
	GameServerAppID         = "600760"
	SaveInterval            string
	AdditionalParams        string
	AutoPauseServer         bool
	UPNPEnabled             bool
	AutoSave                bool
	StartLocalHost          bool
	IsDebugMode             bool
	CreateSSUILogFile       bool
	LogLevel                int
	IsFirstTimeSetup        bool
	LogMessageBuffer        string
	DiscordCharBufferSize   int
	SSEMessageBufferSize    = 2000
	MaxSSEConnections       = 20
	BufferFlushTicker       *time.Ticker
	ControlMessageID        string
	ExceptionMessageID      string
	Users                   map[string]string
	AuthEnabled             bool
	JwtKey                  string
	AuthTokenLifetime       int
	IsUpdateEnabled         bool
	AllowPrereleaseUpdates  bool
	AllowMajorUpdates       bool
)

// LoadConfig loads and initializes the configuration
func LoadConfig() (*JsonConfig, error) {
	var jsonConfig JsonConfig
	file, err := os.Open(ConfigPath)
	if err == nil {
		// File exists, proceed to decode it
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&jsonConfig); err != nil {
			return nil, fmt.Errorf("failed to decode config: %v", err)
		}
	} else if os.IsNotExist(err) {
		// File is missing, log it and proceed with defaults
		fmt.Println("Config file does not exist. Using defaults and environment variables.")
	} else {
		// Other errors (e.g., permissions), fail immediately
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	// Apply configuration with hierarchy
	applyConfig(&jsonConfig)

	return &jsonConfig, nil
}

// applyConfig applies the configuration with JSON -> env -> fallback hierarchy
func applyConfig(cfg *JsonConfig) {
	// Apply values with hierarchy
	DiscordToken = getString(cfg.DiscordToken, "DISCORD_TOKEN", "")
	ControlChannelID = getString(cfg.ControlChannelID, "CONTROL_CHANNEL_ID", "")
	StatusChannelID = getString(cfg.StatusChannelID, "STATUS_CHANNEL_ID", "")
	ConnectionListChannelID = getString(cfg.ConnectionListChannelID, "CONNECTION_LIST_CHANNEL_ID", "")
	LogChannelID = getString(cfg.LogChannelID, "LOG_CHANNEL_ID", "")
	SaveChannelID = getString(cfg.SaveChannelID, "SAVE_CHANNEL_ID", "")
	ControlPanelChannelID = getString(cfg.ControlPanelChannelID, "CONTROL_PANEL_CHANNEL_ID", "")
	DiscordCharBufferSize = getInt(cfg.DiscordCharBufferSize, "DISCORD_CHAR_BUFFER_SIZE", 1000)
	BlackListFilePath = getString(cfg.BlackListFilePath, "BLACKLIST_FILE_PATH", "./Blacklist.txt")

	isDiscordEnabledVal := getBool(cfg.IsDiscordEnabled, "IS_DISCORD_ENABLED", false)
	IsDiscordEnabled = isDiscordEnabledVal
	cfg.IsDiscordEnabled = &isDiscordEnabledVal

	ErrorChannelID = getString(cfg.ErrorChannelID, "ERROR_CHANNEL_ID", "")
	BackupKeepLastN = getInt(cfg.BackupKeepLastN, "BACKUP_KEEP_LAST_N", 2000)

	isCleanupEnabledVal := getBool(cfg.IsCleanupEnabled, "IS_CLEANUP_ENABLED", false)
	IsCleanupEnabled = isCleanupEnabledVal
	cfg.IsCleanupEnabled = &isCleanupEnabledVal

	BackupKeepDailyFor = time.Duration(getInt(cfg.BackupKeepDailyFor, "BACKUP_KEEP_DAILY_FOR", 24)) * time.Hour
	BackupKeepWeeklyFor = time.Duration(getInt(cfg.BackupKeepWeeklyFor, "BACKUP_KEEP_WEEKLY_FOR", 168)) * time.Hour
	BackupKeepMonthlyFor = time.Duration(getInt(cfg.BackupKeepMonthlyFor, "BACKUP_KEEP_MONTHLY_FOR", 730)) * time.Hour
	BackupCleanupInterval = time.Duration(getInt(cfg.BackupCleanupInterval, "BACKUP_CLEANUP_INTERVAL", 730)) * time.Hour
	BackupWaitTime = time.Duration(getInt(cfg.BackupWaitTime, "BACKUP_WAIT_TIME", 30)) * time.Second
	GameBranch = getString(cfg.GameBranch, "GAME_BRANCH", "public")
	ServerName = getString(cfg.ServerName, "SERVER_NAME", "Stationeers Server UI")
	SaveInfo = getString(cfg.SaveInfo, "SAVE_INFO", "Moon Moon")
	ServerMaxPlayers = getString(cfg.ServerMaxPlayers, "SERVER_MAX_PLAYERS", "6")
	ServerPassword = getString(cfg.ServerPassword, "SERVER_PASSWORD", "")
	ServerAuthSecret = getString(cfg.ServerAuthSecret, "SERVER_AUTH_SECRET", "")
	AdminPassword = getString(cfg.AdminPassword, "ADMIN_PASSWORD", "")
	GamePort = getString(cfg.GamePort, "GAME_PORT", "27016")
	UpdatePort = getString(cfg.UpdatePort, "UPDATE_PORT", "27017")

	upnpEnabledVal := getBool(cfg.UPNPEnabled, "UPNP_ENABLED", false)
	UPNPEnabled = upnpEnabledVal
	cfg.UPNPEnabled = &upnpEnabledVal

	autoSaveVal := getBool(cfg.AutoSave, "AUTO_SAVE", true)
	AutoSave = autoSaveVal
	cfg.AutoSave = &autoSaveVal

	SaveInterval = getString(cfg.SaveInterval, "SAVE_INTERVAL", "300")

	autoPauseServerVal := getBool(cfg.AutoPauseServer, "AUTO_PAUSE_SERVER", true)
	AutoPauseServer = autoPauseServerVal
	cfg.AutoPauseServer = &autoPauseServerVal

	LocalIpAddress = getString(cfg.LocalIpAddress, "LOCAL_IP_ADDRESS", "")

	startLocalHostVal := getBool(cfg.StartLocalHost, "START_LOCAL_HOST", true)
	StartLocalHost = startLocalHostVal
	cfg.StartLocalHost = &startLocalHostVal

	serverVisibleVal := getBool(cfg.ServerVisible, "SERVER_VISIBLE", true)
	ServerVisible = serverVisibleVal
	cfg.ServerVisible = &serverVisibleVal

	useSteamP2PVal := getBool(cfg.UseSteamP2P, "USE_STEAM_P2P", true)
	UseSteamP2P = useSteamP2PVal
	cfg.UseSteamP2P = &useSteamP2PVal

	ExePath = getString(cfg.ExePath, "EXE_PATH", getDefaultExePath())
	AdditionalParams = getString(cfg.AdditionalParams, "ADDITIONAL_PARAMS", "")
	Users = getUsers(cfg.Users, "SSUI_USERS", map[string]string{})

	authEnabledVal := getBool(cfg.AuthEnabled, "SSUI_AUTH_ENABLED", false)
	AuthEnabled = authEnabledVal
	cfg.AuthEnabled = &authEnabledVal

	JwtKey = getString(cfg.JwtKey, "SSUI_JWT_KEY", generateJwtKey())
	AuthTokenLifetime = getInt(cfg.AuthTokenLifetime, "SSUI_AUTH_TOKEN_LIFETIME", 1440)

	debugVal := getBool(cfg.Debug, "DEBUG", false)
	IsDebugMode = debugVal
	cfg.Debug = &debugVal

	createSSUILogFileVal := getBool(cfg.CreateSSUILogFile, "CREATE_SSUI_LOGFILE", false)
	CreateSSUILogFile = createSSUILogFileVal
	cfg.CreateSSUILogFile = &createSSUILogFileVal

	LogLevel = getInt(cfg.LogLevel, "LOG_LEVEL", 20)

	isUpdateEnabledVal := getBool(cfg.IsUpdateEnabled, "IS_UPDATE_ENABLED", true)
	IsUpdateEnabled = isUpdateEnabledVal
	cfg.IsUpdateEnabled = &isUpdateEnabledVal

	allowPrereleaseUpdatesVal := getBool(cfg.AllowPrereleaseUpdates, "ALLOW_PRERELEASE_UPDATES", false)
	AllowPrereleaseUpdates = allowPrereleaseUpdatesVal
	cfg.AllowPrereleaseUpdates = &allowPrereleaseUpdatesVal

	allowMajorUpdatesVal := getBool(cfg.AllowMajorUpdates, "ALLOW_MAJOR_UPDATES", false)
	AllowMajorUpdates = allowMajorUpdatesVal
	cfg.AllowMajorUpdates = &allowMajorUpdatesVal

	// Process SaveInfo
	parts := strings.Split(SaveInfo, " ")
	if len(parts) > 0 {
		WorldName = parts[0]
	}
	if len(parts) > 1 {
		BackupWorldName = parts[1]
	}

	// Set backup paths
	ConfiguredBackupDir = filepath.Join("./saves/", WorldName, "Backup")
	ConfiguredSafeBackupDir = filepath.Join("./saves/", WorldName, "Safebackups")
}
