package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Although this is a not a real setter, this function can be used to save the config safely
func SetSaveConfig() error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SafeSaveConfig()
}

// Debug and Logging Settings
func SetIsDebugMode(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDebugMode = value
	return SafeSaveConfig()
}

// SetCreateSSUILogFile sets the CreateSSUILogFile with validation
func SetCreateSSUILogFile(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	CreateSSUILogFile = value
	return SafeSaveConfig()
}

// SetLogLevel sets the LogLevel with validation
func SetLogLevel(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("log level cannot be negative")
	}

	LogLevel = value
	return SafeSaveConfig()
}

// SetLogMessageBuffer sets the LogMessageBuffer with validation
func SetLogMessageBuffer(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogMessageBuffer = value
	return SafeSaveConfig()
}

func SetLogClutterToConsole(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogClutterToConsole = value
	return SafeSaveConfig()
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
	return SafeSaveConfig()
}

// Setup and System Settings
func SetIsFirstTimeSetup(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsFirstTimeSetup = value
	return SafeSaveConfig()
}

// SetBufferFlushTicker sets the BufferFlushTicker
func SetBufferFlushTicker(value *time.Ticker) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BufferFlushTicker = value
	return SafeSaveConfig()
}

// SetSSEMessageBufferSize sets the SSEMessageBufferSize with validation
func SetSSEMessageBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("SSE message buffer size must be positive")
	}

	SSEMessageBufferSize = value
	return SafeSaveConfig()
}

// SetMaxSSEConnections sets the MaxSSEConnections with validation
func SetMaxSSEConnections(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("max SSE connections must be positive")
	}

	MaxSSEConnections = value
	return SafeSaveConfig()
}

func SetLanguageSetting(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LanguageSetting = value
	return SafeSaveConfig()
}

func SetSSUIWebPort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("SSUI web port cannot be empty")
	}

	SSUIWebPort = value
	return SafeSaveConfig()
}

func SetAdditionalLoginHeaderText(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AdditionalLoginHeaderText = value
	return SafeSaveConfig()
}

// Game Settings
func SetGameBranch(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("game branch cannot be empty")
	}

	GameBranch = value
	return SafeSaveConfig()
}

func SetGameServerUUID(value uuid.UUID) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value == uuid.Nil {
		return fmt.Errorf("game server UUID cannot be nil")
	}

	GameServerUUID = value
	return SafeSaveConfig()
}

func SetDifficulty(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	Difficulty = value
	return SafeSaveConfig()
}

func SetStartCondition(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StartCondition = value
	return SafeSaveConfig()
}

func SetStartLocation(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StartLocation = value
	return SafeSaveConfig()
}

func SetIsNewTerrainAndSaveSystem(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsNewTerrainAndSaveSystem = value
	return SafeSaveConfig()
}

// Server Settings
func SetServerName(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerName = value
	return SafeSaveConfig()
}

func SetSaveInfo(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SaveInfo = value
	return SafeSaveConfig()
}

func SetServerMaxPlayers(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerMaxPlayers = value
	return SafeSaveConfig()
}

func SetServerPassword(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerPassword = value
	return SafeSaveConfig()
}

func SetServerAuthSecret(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerAuthSecret = value
	return SafeSaveConfig()
}

func SetAdminPassword(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AdminPassword = value
	return SafeSaveConfig()
}

func SetGamePort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	GamePort = value
	return SafeSaveConfig()
}

func SetUpdatePort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	UpdatePort = value
	return SafeSaveConfig()
}

func SetUPNPEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	UPNPEnabled = value
	return SafeSaveConfig()
}

func SetAutoSave(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoSave = value
	return SafeSaveConfig()
}

func SetSaveInterval(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SaveInterval = value
	return SafeSaveConfig()
}

func SetAutoPauseServer(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoPauseServer = value
	return SafeSaveConfig()
}

func SetLocalIpAddress(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LocalIpAddress = value
	return SafeSaveConfig()
}

func SetStartLocalHost(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StartLocalHost = value
	return SafeSaveConfig()
}

func SetServerVisible(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ServerVisible = value
	return SafeSaveConfig()
}

func SetUseSteamP2P(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	UseSteamP2P = value
	return SafeSaveConfig()
}

func SetExePath(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ExePath = value
	return SafeSaveConfig()
}

func SetAdditionalParams(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AdditionalParams = value
	return SafeSaveConfig()
}

func SetAutoStartServerOnStartup(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoStartServerOnStartup = value
	return SafeSaveConfig()
}

func SetAutoRestartServerTimer(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoRestartServerTimer = value
	return SafeSaveConfig()
}

// Backup Settings
func SetBackupKeepLastN(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep last N cannot be negative")
	}

	BackupKeepLastN = value
	return SafeSaveConfig()
}

func SetIsCleanupEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsCleanupEnabled = value
	return SafeSaveConfig()
}

func SetBackupKeepDailyFor(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep daily for cannot be negative")
	}

	BackupKeepDailyFor = time.Duration(value) * time.Hour
	return SafeSaveConfig()
}

func SetBackupKeepWeeklyFor(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep weekly for cannot be negative")
	}

	BackupKeepWeeklyFor = time.Duration(value) * time.Hour
	return SafeSaveConfig()
}

func SetBackupKeepMonthlyFor(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup keep monthly for cannot be negative")
	}

	BackupKeepMonthlyFor = time.Duration(value) * time.Hour
	return SafeSaveConfig()
}

func SetBackupCleanupInterval(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("backup cleanup interval must be positive")
	}

	BackupCleanupInterval = time.Duration(value) * time.Hour
	return SafeSaveConfig()
}

func SetBackupWaitTime(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("backup wait time cannot be negative")
	}

	BackupWaitTime = time.Duration(value) * time.Second
	return SafeSaveConfig()
}

// Discord Settings
func SetIsDiscordEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDiscordEnabled = value
	return SafeSaveConfig()
}

func SetDiscordToken(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	DiscordToken = value
	return SafeSaveConfig()
}

func SetControlChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlChannelID = value
	return SafeSaveConfig()
}

// SetStatusChannelID sets the StatusChannelID
func SetStatusChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StatusChannelID = value
	return SafeSaveConfig()
}

// SetLogChannelID sets the LogChannelID
func SetLogChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogChannelID = value
	return SafeSaveConfig()
}

// SetErrorChannelID sets the ErrorChannelID
func SetErrorChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ErrorChannelID = value
	return SafeSaveConfig()
}

// SetConnectionListChannelID sets the ConnectionListChannelID
func SetConnectionListChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ConnectionListChannelID = value
	return SafeSaveConfig()
}

// SetSaveChannelID sets the SaveChannelID
func SetSaveChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SaveChannelID = value
	return SafeSaveConfig()
}

// SetControlPanelChannelID sets the ControlPanelChannelID
func SetControlPanelChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlPanelChannelID = value
	return SafeSaveConfig()
}

// SetDiscordCharBufferSize sets the DiscordCharBufferSize with validation
func SetDiscordCharBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("discord char buffer size must be positive")
	}

	DiscordCharBufferSize = value
	return SafeSaveConfig()
}

func SetControlMessageID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlMessageID = value
	return SafeSaveConfig()
}

func SetExceptionMessageID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ExceptionMessageID = value
	return SafeSaveConfig()
}

// SetBlackListFilePath sets the BlackListFilePath with validation
func SetBlackListFilePath(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("blacklist file path cannot be empty")
	}

	BlackListFilePath = value
	return SafeSaveConfig()
}

func SetAuthEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AuthEnabled = value
	return SafeSaveConfig()
}

// SetJwtKey sets the JwtKey with validation
func SetJwtKey(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if len(value) < 32 {
		return fmt.Errorf("JWT key must be at least 32 bytes")
	}

	JwtKey = value
	return SafeSaveConfig()
}

// SetAuthTokenLifetime sets the AuthTokenLifetime with validation
func SetAuthTokenLifetime(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("auth token lifetime must be positive")
	}

	AuthTokenLifetime = value
	return SafeSaveConfig()
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

	return SafeSaveConfig()
}

// Update Settings
func SetIsUpdateEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsUpdateEnabled = value
	return SafeSaveConfig()
}

func SetAllowPrereleaseUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowPrereleaseUpdates = value
	return SafeSaveConfig()
}

func SetAllowMajorUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowMajorUpdates = value
	return SafeSaveConfig()
}

func SetIsSSCMEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsSSCMEnabled = value
	return SafeSaveConfig()
}

func SetIsConsoleEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsConsoleEnabled = value
	return SafeSaveConfig()
}
