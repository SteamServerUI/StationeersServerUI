package config

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

/*
config.Version and config.Branch can be found in config.go

ConfigMu protects all config variables. Lock it for writes; reads are safe
if writes only happen via applyConfig or with ConfigMu locked.

WARNING: Do NOT set any config vars without locking ConfigMu:
config.ConfigMu.Lock()
config.SomeConfigVar = newValue
config.ConfigMu.Unlock()
*/

var ConfigMu sync.Mutex

// Game Server configuration
var (
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
	AdditionalParams string
	UPNPEnabled      bool
	StartLocalHost   bool
	WorldName        string
	BackupWorldName  string
	SaveInfo         string
	SaveInterval     string
	AutoPauseServer  bool
	AutoSave         bool
)

// Logging, debugging and misc
var (
	IsDebugMode          bool //only used for pprof server, keep it like this and check the log level instead. Debug = 10
	CreateSSUILogFile    bool
	LogLevel             int
	LogMessageBuffer     string
	IsFirstTimeSetup     bool
	BufferFlushTicker    *time.Ticker
	SSEMessageBufferSize = 2000
	MaxSSEConnections    = 20
	GameServerAppID      = "600760"
	ExePath              string
	GameBranch           string
	SubsystemFilters     []string
)

// Discord integration
var (
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
	DiscordCharBufferSize   int
	ControlMessageID        string
	ExceptionMessageID      string
	BlackListFilePath       string
)

// Backup and cleanup settings
var (
	IsCleanupEnabled        bool
	BackupKeepLastN         int
	BackupKeepDailyFor      time.Duration
	BackupKeepWeeklyFor     time.Duration
	BackupKeepMonthlyFor    time.Duration
	BackupCleanupInterval   time.Duration
	ConfiguredBackupDir     string
	ConfiguredSafeBackupDir string
	BackupWaitTime          time.Duration
)

// Authentication and security
var (
	AuthEnabled       bool
	JwtKey            string
	AuthTokenLifetime int
	Users             map[string]string
)

// SSUI Updates
var (
	IsUpdateEnabled        bool
	AllowPrereleaseUpdates bool
	AllowMajorUpdates      bool
)

// SSCM (Stationeers Server Command Manager) settings

var (
	IsSSCMEnabled = true
)

// File paths
var (
	TLSCertPath              = "./UIMod/tls/cert.pem"
	TLSKeyPath               = "./UIMod/tls/key.pem"
	ConfigPath               = "./UIMod/config/config.json"
	CustomDetectionsFilePath = "./UIMod/detectionmanager/customdetections.json"
	LogFolder                = "./UIMod/logs/"
	UIModFolder              = "./UIMod/"
	TwoBoxFormFolder         = "./UIMod/twoboxform/"
	ConfigHtmlPath           = "./UIMod/ui/config.html"
	DetectionManagerHtmlPath = "./UIMod/detectionmanager/detectionmanager.html"
	TwoBoxFormHtmlPath       = "./UIMod/twoboxform/twoboxform.html"
	IndexHtmlPath            = "./UIMod/ui/index.html"
	SSCMFilePath             = "./BepInEx/plugins/SSCM/SSCM.socket"
)
