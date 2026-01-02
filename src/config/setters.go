package config

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Although this is a not a real setter, this function can be used to save the config safely
func SetSaveConfig() error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return safeSaveConfigAtomic()
}

// Setup and System Settings
func SetIsFirstTimeSetup(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsFirstTimeSetup = value
	return safeSaveConfigAtomic()
}

func SetIsSSCMEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsSSCMEnabled = value
	return safeSaveConfigAtomic()
}

func SetIsBepInExEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsBepInExEnabled = value
	return safeSaveConfigAtomic()
}

func SetCurrentBranchBuildID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	CurrentBranchBuildID = value
	return nil
}

func SetExtractedGameVersion(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ExtractedGameVersion = value
	return nil
}

func SetSkipSteamCMD(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SkipSteamCMD = value
	return nil
}

func SetIsDockerContainer(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDockerContainer = value
	return nil
}

func SetNoSanityCheck(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	NoSanityCheck = value
	return nil
}

// SetRunfileIdentifier sets the RunfileIdentifier with validation
func SetRunfileIdentifier(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("runfile identifier cannot be empty")
	}
	if !isValidRunfileIdentifier(value) {
		return fmt.Errorf("invalid runfile identifier; must be alphanumeric, dash or underscore only")
	}

	RunfileIdentifier = value
	return safeSaveConfigAtomic()
}

// isValidRunfileIdentifier checks that the identifier is safe to use as a file component
func isValidRunfileIdentifier(s string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, s)
	return err == nil && matched
}

// Debug and Logging Settings
func SetIsDebugMode(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDebugMode = value
	return safeSaveConfigAtomic()
}

// SetCreateSSUILogFile sets the CreateSSUILogFile with validation
func SetCreateSSUILogFile(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	CreateSSUILogFile = value
	return safeSaveConfigAtomic()
}

// SetLogLevel sets the LogLevel with validation
func SetLogLevel(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value < 0 {
		return fmt.Errorf("log level cannot be negative")
	}

	LogLevel = value
	return safeSaveConfigAtomic()
}

func SetLogClutterToConsole(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogClutterToConsole = value
	return safeSaveConfigAtomic()
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
	return safeSaveConfigAtomic()
}

// SetSSEMessageBufferSize sets the SSEMessageBufferSize with validation
func SetSSEMessageBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("SSE message buffer size must be positive")
	}

	SSEMessageBufferSize = value
	return safeSaveConfigAtomic()
}

// SetMaxSSEConnections sets the MaxSSEConnections with validation
func SetMaxSSEConnections(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("max SSE connections must be positive")
	}

	MaxSSEConnections = value
	return safeSaveConfigAtomic()
}

func SetLanguageSetting(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LanguageSetting = value
	return safeSaveConfigAtomic()
}

func SetBackendEndpointPort(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("BackendEndpointPort cannot be empty")
	}

	BackendEndpointPort = value
	return safeSaveConfigAtomic()
}

func SetBackendName(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackendName = value
	return safeSaveConfigAtomic()
}

// Game Settings
func SetGameBranch(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("game branch cannot be empty")
	}

	GameBranch = value
	return safeSaveConfigAtomic()
}

func SetAutoStartServerOnStartup(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoStartServerOnStartup = value
	return safeSaveConfigAtomic()
}

func SetAutoRestartServerTimer(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AutoRestartServerTimer = value
	return safeSaveConfigAtomic()
}

// Discord Settings
func SetIsDiscordEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsDiscordEnabled = value
	return safeSaveConfigAtomic()
}

func SetDiscordToken(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	DiscordToken = value
	return safeSaveConfigAtomic()
}

func SetControlChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlChannelID = value
	return safeSaveConfigAtomic()
}

// SetStatusChannelID sets the StatusChannelID
func SetStatusChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	StatusChannelID = value
	return safeSaveConfigAtomic()
}

// SetLogChannelID sets the LogChannelID
func SetLogChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	LogChannelID = value
	return safeSaveConfigAtomic()
}

// SetErrorChannelID sets the ErrorChannelID
func SetErrorChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ErrorChannelID = value
	return safeSaveConfigAtomic()
}

// SetConnectionListChannelID sets the ConnectionListChannelID
func SetConnectionListChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ConnectionListChannelID = value
	return safeSaveConfigAtomic()
}

// SetSaveChannelID sets the SaveChannelID
func SetSaveChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	SaveChannelID = value
	return safeSaveConfigAtomic()
}

// SetControlPanelChannelID sets the ControlPanelChannelID
func SetControlPanelChannelID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ControlPanelChannelID = value
	return safeSaveConfigAtomic()
}

// SetDiscordCharBufferSize sets the DiscordCharBufferSize with validation
func SetDiscordCharBufferSize(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("discord char buffer size must be positive")
	}

	DiscordCharBufferSize = value
	return safeSaveConfigAtomic()
}

func SetExceptionMessageID(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	ExceptionMessageID = value
	return safeSaveConfigAtomic()
}

func SetAuthEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AuthEnabled = value
	return safeSaveConfigAtomic()
}

// SetJwtKey sets the JwtKey with validation
func SetJwtKey(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if len(value) < 32 {
		return fmt.Errorf("JWT key must be at least 32 bytes")
	}

	JwtKey = value
	return safeSaveConfigAtomic()
}

// SetAuthTokenLifetime sets the AuthTokenLifetime with validation
func SetAuthTokenLifetime(value int) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if value <= 0 {
		return fmt.Errorf("auth token lifetime must be positive")
	}

	AuthTokenLifetime = value
	return safeSaveConfigAtomic()
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

	return safeSaveConfigAtomic()
}

// Update Settings
func SetIsUpdateEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsUpdateEnabled = value
	return safeSaveConfigAtomic()
}

func SetAllowPrereleaseUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowPrereleaseUpdates = value
	return safeSaveConfigAtomic()
}

func SetAllowMajorUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowMajorUpdates = value
	return safeSaveConfigAtomic()
}
func SetIsSSUICLIConsoleEnabled(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsSSUICLIConsoleEnabled = value
	return safeSaveConfigAtomic()
}

func SetAllowAutoGameServerUpdates(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	AllowAutoGameServerUpdates = value
	return safeSaveConfigAtomic()
}

// SetUsers merges the provided key-value pairs into the existing Users map with validation
func SetRegisteredPlugins(value map[string]string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	// Initialize Users map if it's nil
	if RegisteredPlugins == nil {
		RegisteredPlugins = make(map[string]string)
	}

	// Validate and merge each key-value pair
	for k, v := range value {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			return fmt.Errorf("RegisteredPlugin key or value cannot be empty")
		}
		RegisteredPlugins[k] = v // Update or add the key-value pair
	}

	return safeSaveConfigAtomic()
}

func SetGameLogFromLogFile(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	GameLogFromLogFile = value
	return safeSaveConfigAtomic()
}

func SetBackupsStoreDir(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupsStoreDir = value
	return safeSaveConfigAtomic()
}

func SetBackupLoopInterval(value time.Duration) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupLoopInterval = value
	return safeSaveConfigAtomic()
}

func SetBackupMode(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupMode = value
	return safeSaveConfigAtomic()
}

func SetBackupMaxFileSize(value int64) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupMaxFileSize = value
	return safeSaveConfigAtomic()
}

func SetBackupUseCompression(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupUseCompression = value
	return safeSaveConfigAtomic()
}

func SetBackupKeepSnapshot(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupKeepSnapshot = value
	return safeSaveConfigAtomic()
}

func SetBackupLoopActive(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	BackupLoopActive = value
	return safeSaveConfigAtomic()
}

func SetIsGameServerRunning(value bool) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	IsGameServerRunning = value
	return nil
}
