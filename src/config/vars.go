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
	SaveInfo         string
	SaveName         string
	WorldID          string
	SaveInterval     string
	AutoPauseServer  bool
	AutoSave         bool
	Difficulty       string
	StartCondition   string
	StartLocation    string
)

// Logging, debugging and misc
var (
	IsDebugMode              bool //only used for pprof server, keep it like this and check the log level instead. Debug = 10
	CreateSSUILogFile        bool
	LogLevel                 int
	IsFirstTimeSetup         bool
	SSEMessageBufferSize     = 2000
	MaxSSEConnections        = 20
	GameServerAppID          = "600760"
	ExePath                  string
	GameBranch               string
	SubsystemFilters         []string
	AutoRestartServerTimer   string
	IsConsoleEnabled         bool
	LogClutterToConsole      bool // surpresses clutter mono logs from the gameserver
	LanguageSetting          string
	AutoStartServerOnStartup bool
	SSUIIdentifier           string
)

// SteamServerUI Settings
var (
	UseRunfiles       bool
	RunfileIdentifier string
	IsStationeersMode bool
)

// Runtime only variables

var (
	CurrentBranchBuildID string // ONLY RUNTIME
	ExtractedGameVersion string // ONLY RUNTIME
	SkipSteamCMD         bool   // ONLY RUNTIME
	IsDockerContainer    bool   // ONLY RUNTIME
	NoSanityCheck        bool   // ONLY RUNTIME
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
	BlackListFilePath       string
)

// Backup and cleanup settings
var (
	IsCleanupEnabled          bool
	BackupKeepLastN           int
	BackupKeepDailyFor        time.Duration
	BackupKeepWeeklyFor       time.Duration
	BackupKeepMonthlyFor      time.Duration
	BackupCleanupInterval     time.Duration
	ConfiguredBackupDir       string
	ConfiguredSafeBackupDir   string
	BackupWaitTime            time.Duration
	IsNewTerrainAndSaveSystem bool
)

// Authentication and security
var (
	AuthEnabled       bool
	JwtKey            string
	AuthTokenLifetime int
	Users             map[string]string
	SSUIWebPort       string
)

// SSUI Updates and Game Server Updates
var (
	IsUpdateEnabled            bool
	AllowPrereleaseUpdates     bool
	AllowMajorUpdates          bool
	AllowAutoGameServerUpdates bool
)

// SSCM (Stationeers Server Command Manager) settings

var (
	IsSSCMEnabled bool
)

// File paths
var (
	TLSCertPath              = "./UIMod/tls/cert.pem"
	TLSKeyPath               = "./UIMod/tls/key.pem"
	ConfigPath               = "./UIMod/config/config.json"
	CustomDetectionsFilePath = "./UIMod/config/customdetections.json"
	LogFolder                = "./UIMod/logs/"
	UIModFolder              = "./UIMod/"
	TwoBoxFormFolder         = "./UIMod/twoboxform/"
	ConfigHtmlPath           = "./UIMod/ui/config.html"
	DetectionManagerHtmlPath = "./UIMod/ui/detectionmanager.html"
	TwoBoxFormHtmlPath       = "./UIMod/twoboxform/twoboxform.html"
	IndexHtmlPath            = "./UIMod/ui/index.html"
	SSCMWebDir               = "./UIMod/sscm/"
	SSCMFilePath             = "./BepInEx/plugins/SSCM/SSCM.socket"
	SSCMPluginDir            = "./BepInEx/plugins/SSCM/"
	RunFilesFolder           = "./UIMod/runfiles/"
)

// Bundled Assets

var V1UIFS embed.FS
