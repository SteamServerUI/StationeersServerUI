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
	serverName       string
	serverMaxPlayers string
	serverPassword   string
	serverAuthSecret string
	adminPassword    string
	gamePort         string
	updatePort       string
	localIpAddress   string
	serverVisible    bool
	useSteamP2P      bool
	additionalParams string
	uPNPEnabled      bool
	startLocalHost   bool
	saveInfo         string
	saveName         string
	worldID          string
	saveInterval     string
	autoPauseServer  bool
	autoSave         bool
	difficulty       string
	startCondition   string
	startLocation    string
)

// Logging, debugging and misc
var (
	isDebugMode              bool //only used for pprof server, keep it like this and check the log level instead. Debug = 10
	createSSUILogFile        bool
	logLevel                 int
	isFirstTimeSetup         bool
	sseMessageBufferSize     = 2000
	maxSSEConnections        = 20
	gameServerAppID          = "600760"
	exePath                  string
	gameBranch               string
	subsystemFilters         []string
	autoRestartServerTimer   string
	isConsoleEnabled         bool
	logClutterToConsole      bool // surpresses clutter mono logs from the gameserver
	languageSetting          string
	autoStartServerOnStartup bool
	ssuiIdentifier           string
	overrideAdvertisedIp     string
)

// Runtime only variables

var (
	currentBranchBuildID string // ONLY RUNTIME
	extractedGameVersion string // ONLY RUNTIME
	skipSteamCMD         bool   // ONLY RUNTIME
	isDockerContainer    bool   // ONLY RUNTIME
	noSanityCheck        bool   // ONLY RUNTIME
)

// Discord integration
var (
	discordToken            string
	discordSession          *discordgo.Session
	isDiscordEnabled        bool
	controlChannelID        string
	statusChannelID         string
	logChannelID            string
	errorChannelID          string
	connectionListChannelID string
	saveChannelID           string
	controlPanelChannelID   string
	discordCharBufferSize   int
	exceptionMessageID      string
	blackListFilePath       string
)

// Backup and cleanup settings
var (
	isCleanupEnabled          bool
	backupKeepLastN           int
	backupKeepDailyFor        time.Duration
	backupKeepWeeklyFor       time.Duration
	backupKeepMonthlyFor      time.Duration
	backupCleanupInterval     time.Duration
	configuredBackupDir       string
	configuredSafeBackupDir   string
	backupWaitTime            time.Duration
	isNewTerrainAndSaveSystem bool
)

// Authentication and security
var (
	authEnabled       bool
	jwtKey            string
	authTokenLifetime int
	users             map[string]string
	ssuiWebPort       string
)

// SSUI Updates and Game Server Updates
var (
	isUpdateEnabled            bool
	allowPrereleaseUpdates     bool
	allowMajorUpdates          bool
	allowAutoGameServerUpdates bool
)

// SSCM (Stationeers Server Command Manager) settings

var (
	isSSCMEnabled bool
)

// File paths
var (
	tlsCertPath              = "./UIMod/tls/cert.pem"
	tlsKeyPath               = "./UIMod/tls/key.pem"
	configPath               = "./UIMod/config/config.json"
	customDetectionsFilePath = "./UIMod/config/customdetections.json"
	logFolder                = "./UIMod/logs/"
	uiModFolder              = "./UIMod/"
	sscmWebDir               = "./UIMod/sscm/"
	sscmFilePath             = "./BepInEx/plugins/SSCM/SSCM.socket"
	sscmPluginDir            = "./BepInEx/plugins/SSCM/"
)

// Bundled Assets

var v1UIFS embed.FS
