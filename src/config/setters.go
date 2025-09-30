package config

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Although this is a not a real setter, this function can be used to save the config safely
func SetSaveConfig() error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return safeSaveConfig()
}

// Setup and System Settings
func SetIsFirstTimeSetup(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	isFirstTimeSetup = value
	return safeSaveConfig()
}

func SetIsSSCMEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	isSSCMEnabled = value
	return safeSaveConfig()
}

func SetCurrentBranchBuildID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	currentBranchBuildID = value
	return nil
}

func SetExtractedGameVersion(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	extractedGameVersion = value
	return nil
}

func SetSkipSteamCMD(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	skipSteamCMD = value
	return nil
}

func SetIsDockerContainer(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	isDockerContainer = value
	return nil
}

func SetNoSanityCheck(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	noSanityCheck = value
	return nil
}

func SetSaveName(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	saveName = value
	return nil
}

func SetWorldID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	worldID = value
	return nil
}

// ALL SETTERS BELOW THIS LINE ARE UNUSED AT THE MOMENT
// ALL SETTERS BELOW THIS LINE ARE UNUSED AT THE MOMENT
// ALL SETTERS BELOW THIS LINE ARE UNUSED AT THE MOMENT
// ALL SETTERS BELOW THIS LINE ARE UNUSED AT THE MOMENT
// ALL SETTERS BELOW THIS LINE ARE UNUSED AT THE MOMENT

// Debug and Logging Settings
func SetIsDebugMode(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	isDebugMode = value
	return safeSaveConfig()
}

// SetCreateSSUILogFile sets the CreateSSUILogFile with validation
func SetCreateSSUILogFile(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	createSSUILogFile = value
	return safeSaveConfig()
}

// SetLogLevel sets the LogLevel with validation
func SetLogLevel(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("log level cannot be negative")
	}

	logLevel = value
	return safeSaveConfig()
}

func SetLogClutterToConsole(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	logClutterToConsole = value
	return safeSaveConfig()
}

func SetSubsystemFilters(value []string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	for _, v := range value {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("subsystem filter cannot be empty")
		}
	}

	subsystemFilters = value
	return safeSaveConfig()
}

// SetSSEMessageBufferSize sets the SSEMessageBufferSize with validation
func SetSSEMessageBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("SSE message buffer size must be positive")
	}

	sseMessageBufferSize = value
	return safeSaveConfig()
}

// SetMaxSSEConnections sets the MaxSSEConnections with validation
func SetMaxSSEConnections(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("max SSE connections must be positive")
	}

	maxSSEConnections = value
	return safeSaveConfig()
}

func SetLanguageSetting(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	languageSetting = value
	return safeSaveConfig()
}

func SetSSUIWebPort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("SSUI web port cannot be empty")
	}

	ssuiWebPort = value
	return safeSaveConfig()
}

func SetSSUIIdentifier(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ssuiIdentifier = value
	return safeSaveConfig()
}

// Game Settings
func SetGameBranch(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("game branch cannot be empty")
	}

	gameBranch = value
	return safeSaveConfig()
}

func SetDifficulty(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	difficulty = value
	return safeSaveConfig()
}

func SetStartCondition(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	startCondition = value
	return safeSaveConfig()
}

func SetStartLocation(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	startLocation = value
	return safeSaveConfig()
}

func SetIsNewTerrainAndSaveSystem(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	isNewTerrainAndSaveSystem = value
	return safeSaveConfig()
}

// Server Settings
func SetServerName(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	serverName = value
	return safeSaveConfig()
}

func SetSaveInfo(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	saveInfo = value
	return safeSaveConfig()
}

func SetServerMaxPlayers(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	serverMaxPlayers = value
	return safeSaveConfig()
}

func SetServerPassword(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	serverPassword = value
	return safeSaveConfig()
}

func SetServerAuthSecret(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	serverAuthSecret = value
	return safeSaveConfig()
}

func SetAdminPassword(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	adminPassword = value
	return safeSaveConfig()
}

func SetGamePort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	gamePort = value
	return safeSaveConfig()
}

func SetUpdatePort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	updatePort = value
	return safeSaveConfig()
}

func SetUPNPEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	uPNPEnabled = value
	return safeSaveConfig()
}

func SetAutoSave(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	autoSave = value
	return safeSaveConfig()
}

func SetSaveInterval(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	saveInterval = value
	return safeSaveConfig()
}

func SetAutoPauseServer(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	autoPauseServer = value
	return safeSaveConfig()
}

func SetLocalIpAddress(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	localIpAddress = value
	return safeSaveConfig()
}

func SetStartLocalHost(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	startLocalHost = value
	return safeSaveConfig()
}

func SetServerVisible(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	serverVisible = value
	return safeSaveConfig()
}

func SetUseSteamP2P(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	useSteamP2P = value
	return safeSaveConfig()
}

func SetExePath(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	exePath = value
	return safeSaveConfig()
}

func SetAdditionalParams(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	additionalParams = value
	return safeSaveConfig()
}

func SetAutoStartServerOnStartup(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	autoStartServerOnStartup = value
	return safeSaveConfig()
}

func SetAutoRestartServerTimer(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	autoRestartServerTimer = value
	return safeSaveConfig()
}

// Backup Settings
func SetBackupKeepLastN(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep last N cannot be negative")
	}

	backupKeepLastN = value
	return safeSaveConfig()
}

func SetIsCleanupEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	isCleanupEnabled = value
	return safeSaveConfig()
}

func SetBackupKeepDailyFor(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep daily for cannot be negative")
	}

	backupKeepDailyFor = time.Duration(value) * time.Hour
	return safeSaveConfig()
}

func SetBackupKeepWeeklyFor(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep weekly for cannot be negative")
	}

	backupKeepWeeklyFor = time.Duration(value) * time.Hour
	return safeSaveConfig()
}

func SetBackupKeepMonthlyFor(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep monthly for cannot be negative")
	}

	backupKeepMonthlyFor = time.Duration(value) * time.Hour
	return safeSaveConfig()
}

func SetBackupCleanupInterval(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("backup cleanup interval must be positive")
	}

	backupCleanupInterval = time.Duration(value) * time.Hour
	return safeSaveConfig()
}

func SetBackupWaitTime(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup wait time cannot be negative")
	}

	backupWaitTime = time.Duration(value) * time.Second
	return safeSaveConfig()
}

// Discord Settings
func SetIsDiscordEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	isDiscordEnabled = value
	return safeSaveConfig()
}

func SetDiscordToken(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	discordToken = value
	return safeSaveConfig()
}

func SetDiscordSession(value *discordgo.Session) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	discordSession = value
	return safeSaveConfig()
}

func SetControlChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	controlChannelID = value
	return safeSaveConfig()
}

// SetStatusChannelID sets the StatusChannelID
func SetStatusChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	statusChannelID = value
	return safeSaveConfig()
}

// SetLogChannelID sets the LogChannelID
func SetLogChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	logChannelID = value
	return safeSaveConfig()
}

// SetErrorChannelID sets the ErrorChannelID
func SetErrorChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	errorChannelID = value
	return safeSaveConfig()
}

// SetConnectionListChannelID sets the ConnectionListChannelID
func SetConnectionListChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	connectionListChannelID = value
	return safeSaveConfig()
}

// SetSaveChannelID sets the SaveChannelID
func SetSaveChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	saveChannelID = value
	return safeSaveConfig()
}

// SetControlPanelChannelID sets the ControlPanelChannelID
func SetControlPanelChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	controlPanelChannelID = value
	return safeSaveConfig()
}

// SetDiscordCharBufferSize sets the DiscordCharBufferSize with validation
func SetDiscordCharBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("discord char buffer size must be positive")
	}

	discordCharBufferSize = value
	return safeSaveConfig()
}

func SetExceptionMessageID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	exceptionMessageID = value
	return safeSaveConfig()
}

// SetBlackListFilePath sets the BlackListFilePath with validation
func SetBlackListFilePath(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("blacklist file path cannot be empty")
	}

	blackListFilePath = value
	return safeSaveConfig()
}

func SetAuthEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	authEnabled = value
	return safeSaveConfig()
}

// SetJwtKey sets the JwtKey with validation
func SetJwtKey(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if len(value) < 32 {
		return fmt.Errorf("JWT key must be at least 32 bytes")
	}

	jwtKey = value
	return safeSaveConfig()
}

// SetAuthTokenLifetime sets the AuthTokenLifetime with validation
func SetAuthTokenLifetime(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("auth token lifetime must be positive")
	}

	authTokenLifetime = value
	return safeSaveConfig()
}

// SetUsers merges the provided key-value pairs into the existing Users map with validation
func SetUsers(value map[string]string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	// Initialize Users map if it's nil
	if users == nil {
		users = make(map[string]string)
	}

	// Validate and merge each key-value pair
	for k, v := range value {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			return fmt.Errorf("user key or value cannot be empty")
		}
		users[k] = v // Update or add the key-value pair
	}

	return safeSaveConfig()
}

// Update Settings
func SetIsUpdateEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	isUpdateEnabled = value
	return safeSaveConfig()
}

func SetAllowPrereleaseUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	allowPrereleaseUpdates = value
	return safeSaveConfig()
}

func SetAllowMajorUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	allowMajorUpdates = value
	return safeSaveConfig()
}

func SetIsConsoleEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	isConsoleEnabled = value
	return safeSaveConfig()
}

func SetAllowAutoGameServerUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	allowAutoGameServerUpdates = value
	return safeSaveConfig()
}

func SetOverrideAdvertisedIp(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	overrideAdvertisedIp = value
	return safeSaveConfig()
}

func SetV1UIFS(value embed.FS) {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	v1UIFS = value
}
