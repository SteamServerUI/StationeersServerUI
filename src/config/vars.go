package config

import (
	"embed"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// File paths
const (
	TLSCertPath              = "./UIMod/tls/cert.pem"
	TLSKeyPath               = "./UIMod/tls/key.pem"
	ConfigPath               = "./UIMod/config/config.json"
	CustomDetectionsFilePath = "./UIMod/config/customdetections.json"
	LogFolder                = "./UIMod/logs/"
	UIModFolder              = "./UIMod/"
	SSCMWebDir               = "./UIMod/sscm/"
	SSCMFilePath             = "./BepInEx/plugins/SSCM/SSCM.socket"
	SSCMPluginDir            = "./BepInEx/plugins/SSCM/"
)

/*
config.Version and config.Branch can be found in config.go

ConfigMu protects all config variables. Lock it for writes; reads are safe
if writes only happen via applyConfig or with ConfigMu locked. Uses getters where possible.
*/
var ConfigMu sync.RWMutex

type ConfigValue[T any] struct {
	value     T
	validator func(T) error
}

func (c *ConfigValue[T]) Get() T {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return c.value
}

func (c *ConfigValue[T]) Set(newval T) error {
	if c.validator != nil {
		err := c.validator(newval)
		if err != nil {
			return err
		}
	}

	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	c.value = newval
	return safeSaveConfig()
}

// Updated Configurations
var (
	ServerName       ConfigValue[string]
	ServerMaxPlayers ConfigValue[string] // TODO: why is this a string??
	ServerPassword   ConfigValue[string]
	ServerAuthSecret ConfigValue[string]
	AdminPassword    ConfigValue[string]
	GamePort         ConfigValue[string]
	UpdatePort       ConfigValue[string]
	LocalIpAddress   ConfigValue[string]
	ServerVisible    ConfigValue[bool]
)

// Game Server configuration
var (
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

// Bundled Assets

var v1UIFS embed.FS
