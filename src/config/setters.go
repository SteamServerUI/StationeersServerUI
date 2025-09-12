package config

import (
	"fmt"
	"strings"
	"time"
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

	IsFirstTimeSetup = value
	return safeSaveConfig()
}

func SetIsSSCMEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsSSCMEnabled = value
	return safeSaveConfig()
}

func SetCurrentBranchBuildID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	CurrentBranchBuildID = value
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

	IsDebugMode = value
	return safeSaveConfig()
}

// SetCreateSSUILogFile sets the CreateSSUILogFile with validation
func SetCreateSSUILogFile(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	CreateSSUILogFile = value
	return safeSaveConfig()
}

// SetLogLevel sets the LogLevel with validation
func SetLogLevel(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("log level cannot be negative")
	}

	LogLevel = value
	return safeSaveConfig()
}

func SetLogClutterToConsole(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogClutterToConsole = value
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

	SubsystemFilters = value
	return safeSaveConfig()
}

// SetSSEMessageBufferSize sets the SSEMessageBufferSize with validation
func SetSSEMessageBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("SSE message buffer size must be positive")
	}

	SSEMessageBufferSize = value
	return safeSaveConfig()
}

// SetMaxSSEConnections sets the MaxSSEConnections with validation
func SetMaxSSEConnections(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("max SSE connections must be positive")
	}

	MaxSSEConnections = value
	return safeSaveConfig()
}

func SetLanguageSetting(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LanguageSetting = value
	return safeSaveConfig()
}

func SetSSUIWebPort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("SSUI web port cannot be empty")
	}

	SSUIWebPort = value
	return safeSaveConfig()
}

func SetSSUIIdentifier(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SSUIIdentifier = value
	return safeSaveConfig()
}

// Game Settings
func SetGameBranch(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("game branch cannot be empty")
	}

	GameBranch = value
	return safeSaveConfig()
}

func SetDifficulty(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	Difficulty = value
	return safeSaveConfig()
}

func SetStartCondition(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StartCondition = value
	return safeSaveConfig()
}

func SetStartLocation(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StartLocation = value
	return safeSaveConfig()
}

func SetIsNewTerrainAndSaveSystem(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsNewTerrainAndSaveSystem = value
	return safeSaveConfig()
}

// Server Settings
func SetServerName(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerName = value
	return safeSaveConfig()
}

func SetSaveInfo(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SaveInfo = value
	return safeSaveConfig()
}

func SetServerMaxPlayers(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerMaxPlayers = value
	return safeSaveConfig()
}

func SetServerPassword(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerPassword = value
	return safeSaveConfig()
}

func SetServerAuthSecret(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerAuthSecret = value
	return safeSaveConfig()
}

func SetAdminPassword(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AdminPassword = value
	return safeSaveConfig()
}

func SetGamePort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	GamePort = value
	return safeSaveConfig()
}

func SetUpdatePort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	UpdatePort = value
	return safeSaveConfig()
}

func SetUPNPEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	UPNPEnabled = value
	return safeSaveConfig()
}

func SetAutoSave(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoSave = value
	return safeSaveConfig()
}

func SetSaveInterval(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SaveInterval = value
	return safeSaveConfig()
}

func SetAutoPauseServer(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoPauseServer = value
	return safeSaveConfig()
}

func SetLocalIpAddress(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LocalIpAddress = value
	return safeSaveConfig()
}

func SetStartLocalHost(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StartLocalHost = value
	return safeSaveConfig()
}

func SetServerVisible(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerVisible = value
	return safeSaveConfig()
}

func SetUseSteamP2P(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	UseSteamP2P = value
	return safeSaveConfig()
}

func SetExePath(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ExePath = value
	return safeSaveConfig()
}

func SetAdditionalParams(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AdditionalParams = value
	return safeSaveConfig()
}

func SetAutoStartServerOnStartup(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoStartServerOnStartup = value
	return safeSaveConfig()
}

func SetAutoRestartServerTimer(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoRestartServerTimer = value
	return safeSaveConfig()
}

// Backup Settings
func SetBackupKeepLastN(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep last N cannot be negative")
	}

	BackupKeepLastN = value
	return safeSaveConfig()
}

func SetIsCleanupEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsCleanupEnabled = value
	return safeSaveConfig()
}

func SetBackupKeepDailyFor(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep daily for cannot be negative")
	}

	BackupKeepDailyFor = time.Duration(value) * time.Hour
	return safeSaveConfig()
}

func SetBackupKeepWeeklyFor(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep weekly for cannot be negative")
	}

	BackupKeepWeeklyFor = time.Duration(value) * time.Hour
	return safeSaveConfig()
}

func SetBackupKeepMonthlyFor(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep monthly for cannot be negative")
	}

	BackupKeepMonthlyFor = time.Duration(value) * time.Hour
	return safeSaveConfig()
}

func SetBackupCleanupInterval(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("backup cleanup interval must be positive")
	}

	BackupCleanupInterval = time.Duration(value) * time.Hour
	return safeSaveConfig()
}

func SetBackupWaitTime(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup wait time cannot be negative")
	}

	BackupWaitTime = time.Duration(value) * time.Second
	return safeSaveConfig()
}

// Discord Settings
func SetIsDiscordEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDiscordEnabled = value
	return safeSaveConfig()
}

func SetDiscordToken(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	DiscordToken = value
	return safeSaveConfig()
}

func SetControlChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlChannelID = value
	return safeSaveConfig()
}

// SetStatusChannelID sets the StatusChannelID
func SetStatusChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StatusChannelID = value
	return safeSaveConfig()
}

// SetLogChannelID sets the LogChannelID
func SetLogChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogChannelID = value
	return safeSaveConfig()
}

// SetErrorChannelID sets the ErrorChannelID
func SetErrorChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ErrorChannelID = value
	return safeSaveConfig()
}

// SetConnectionListChannelID sets the ConnectionListChannelID
func SetConnectionListChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ConnectionListChannelID = value
	return safeSaveConfig()
}

// SetSaveChannelID sets the SaveChannelID
func SetSaveChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SaveChannelID = value
	return safeSaveConfig()
}

// SetControlPanelChannelID sets the ControlPanelChannelID
func SetControlPanelChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlPanelChannelID = value
	return safeSaveConfig()
}

// SetDiscordCharBufferSize sets the DiscordCharBufferSize with validation
func SetDiscordCharBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("discord char buffer size must be positive")
	}

	DiscordCharBufferSize = value
	return safeSaveConfig()
}

func SetExceptionMessageID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ExceptionMessageID = value
	return safeSaveConfig()
}

// SetBlackListFilePath sets the BlackListFilePath with validation
func SetBlackListFilePath(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("blacklist file path cannot be empty")
	}

	BlackListFilePath = value
	return safeSaveConfig()
}

func SetAuthEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AuthEnabled = value
	return safeSaveConfig()
}

// SetJwtKey sets the JwtKey with validation
func SetJwtKey(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if len(value) < 32 {
		return fmt.Errorf("JWT key must be at least 32 bytes")
	}

	JwtKey = value
	return safeSaveConfig()
}

// SetAuthTokenLifetime sets the AuthTokenLifetime with validation
func SetAuthTokenLifetime(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("auth token lifetime must be positive")
	}

	AuthTokenLifetime = value
	return safeSaveConfig()
}

// SetUsers merges the provided key-value pairs into the existing Users map with validation
func SetUsers(value map[string]string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	// Initialize Users map if it's nil
	if Users == nil {
		Users = make(map[string]string)
	}

	// Validate and merge each key-value pair
	for k, v := range value {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			return fmt.Errorf("user key or value cannot be empty")
		}
		Users[k] = v // Update or add the key-value pair
	}

	return safeSaveConfig()
}

// Update Settings
func SetIsUpdateEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsUpdateEnabled = value
	return safeSaveConfig()
}

func SetAllowPrereleaseUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowPrereleaseUpdates = value
	return safeSaveConfig()
}

func SetAllowMajorUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowMajorUpdates = value
	return safeSaveConfig()
}

func SetIsConsoleEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsConsoleEnabled = value
	return safeSaveConfig()
}

func SetAllowAutoGameServerUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowAutoGameServerUpdates = value
	return safeSaveConfig()
}
