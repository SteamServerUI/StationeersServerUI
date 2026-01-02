package config

import (
	"embed"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

/*
config.Version and config.Branch can be found in config.go

ConfigMu protects all config variables. Lock it for writes; reads are safe
if writes only happen via applyConfig or with ConfigMu locked. Uses getters where possible.
*/

var ConfigMu sync.RWMutex

// Logging, debugging and misc
var (
	IsDebugMode              bool //only used for pprof server, keep it like this and check the log level instead. Debug = 10
	CreateSSUILogFile        bool
	CreateGameServerLogFile  bool
	LogLevel                 int
	IsFirstTimeSetup         bool
	SSEMessageBufferSize     = 2000
	MaxSSEConnections        = 20
	GameServerAppID          = "600760"
	GameBranch               string
	SubsystemFilters         []string
	AutoRestartServerTimer   string
	IsSSUICLIConsoleEnabled  bool
	LogClutterToConsole      bool // surpresses clutter mono logs from the gameserver
	LanguageSetting          string
	AutoStartServerOnStartup bool
	BackendName              string
	GameLogFromLogFile       bool
)

// SteamServerUI Settings
var (
	RunfileIdentifier string
)

// Runtime only variables

var (
	CurrentBranchBuildID string  // ONLY RUNTIME
	ExtractedGameVersion string  // ONLY RUNTIME
	SkipSteamCMD         bool    // ONLY RUNTIME
	IsDockerContainer    bool    // ONLY RUNTIME
	NoSanityCheck        bool    // ONLY RUNTIME
	IsTelemetryEnabled   = false // ONLY RUNTIME (for now)
	IsGameServerRunning  bool    // ONLY RUNTIME
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
	ExceptionMessageID      string
)

// Backup settings
var (
	BackupsStoreDir      string
	BackupLoopActive     bool
	BackupLoopInterval   time.Duration
	BackupMode           string
	BackupMaxFileSize    int64 = 20 * 1024 * 1024 * 1024
	BackupUseCompression bool
	BackupKeepSnapshot   bool
)

// Authentication and security
var (
	AuthEnabled         bool
	JwtKey              string
	AuthTokenLifetime   int
	Users               map[string]string
	BackendEndpointPort string
)

// SSUI Updates and Game Server Updates
var (
	IsUpdateEnabled            bool
	AllowPrereleaseUpdates     bool
	AllowMajorUpdates          bool
	AllowAutoGameServerUpdates bool
)

// BepInEx settings

var (
	IsSSCMEnabled    bool
	IsBepInExEnabled bool
)

// Plugins

var (
	RegisteredPlugins map[string]string
)

// File paths
var (
	TLSCertPath              = "./SSUI/tls/cert.pem"
	TLSKeyPath               = "./SSUI/tls/key.pem"
	ConfigPath               = "./SSUI/config/config.json"
	CustomDetectionsFilePath = "./SSUI/config/customdetections.json"
	LogFolder                = "./SSUI/logs/"
	SSUIFolder               = "./SSUI/"
	TwoBoxFormFolder         = "./SSUI/twoboxform/"
	ConfigHtmlPath           = "./SSUI/ui/config.html"
	DetectionManagerHtmlPath = "./SSUI/ui/detectionmanager.html"
	TwoBoxFormHtmlPath       = "./SSUI/twoboxform/twoboxform.html"
	IndexHtmlPath            = "./SSUI/ui/index.html"
	SSCMWebDir               = "./SSUI/sscm/"
	SSCMFilePath             = "/BepInEx/plugins/SSCM/SSCM.socket"
	SSCMPluginDir            = "/BepInEx/plugins/SSCM/"
	RunFilesFolder           = "./SSUI/runfiles/"
	PluginsFolder            = "./SSUI/plugins/"
)

// Bundled Assets

var V1UIFS embed.FS
