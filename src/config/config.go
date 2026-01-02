package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	// All configuration variables can be found in vars.go
	Version = "7.0.0"
	Branch  = "release"
)

/*
If you read this, you are likely a developer. I sincerely apologize for the way the config works.
While I would love to refactor the config to not write to file then read the file every time a config value is changed,
I have not found the time to do so. So, for now, we save to file, then read the file and rely on whatever the file says. Although this is not ideal, it works for now. Deal with it.
JacksonTheMaster
*/

type JsonConfig struct {
	// reordered in 5.6.4 to simplify the order of the config file.

	// Gameserver Settings
	RunfileIdentifier string `json:"RunfileIdentifier"`
	GameBranch        string `json:"gameBranch"`

	// Logging and debug settings
	Debug                   *bool    `json:"Debug"`
	CreateSSUILogFile       *bool    `json:"CreateSSUILogFile"`
	CreateGameServerLogFile *bool    `json:"CreateGameServerLogFile"`
	LogLevel                int      `json:"LogLevel"`
	GameLogFromLogFile      *bool    `json:"GameLogFromLogFile"`
	SubsystemFilters        []string `json:"subsystemFilters"`

	// Authentication Settings
	Users             map[string]string `json:"users"`       // Map of username to hashed password
	AuthEnabled       *bool             `json:"authEnabled"` // Toggle for enabling/disabling auth
	JwtKey            string            `json:"JwtKey"`
	AuthTokenLifetime int               `json:"AuthTokenLifetime"`

	// SSUI Settings
	LogClutterToConsole      *bool             `json:"LogClutterToConsole"`
	IsSSCMEnabled            *bool             `json:"IsSSCMEnabled"`
	IsBepInExEnabled         *bool             `json:"IsBepInExEnabled"`
	AutoRestartServerTimer   string            `json:"AutoRestartServerTimer"`
	IsSSUICLIConsoleEnabled  *bool             `json:"IsSSUICLIConsoleEnabled"`
	LanguageSetting          string            `json:"LanguageSetting"`
	AutoStartServerOnStartup *bool             `json:"AutoStartServerOnStartup"`
	BackendName              string            `json:"BackendName"`
	BackendEndpointPort      string            `json:"BackendEndpointPort"`
	RegisteredPlugins        map[string]string `json:"RegisteredPlugins"`

	// Update Settings
	IsUpdateEnabled            *bool `json:"IsUpdateEnabled"`
	AllowPrereleaseUpdates     *bool `json:"AllowPrereleaseUpdates"`
	AllowMajorUpdates          *bool `json:"AllowMajorUpdates"`
	AllowAutoGameServerUpdates *bool `json:"AllowAutoGameServerUpdates"`

	// Discord Settings
	DiscordToken            string `json:"discordToken"`
	ControlChannelID        string `json:"controlChannelID"`
	StatusChannelID         string `json:"statusChannelID"`
	ConnectionListChannelID string `json:"connectionListChannelID"`
	LogChannelID            string `json:"logChannelID"`
	SaveChannelID           string `json:"saveChannelID"`
	ControlPanelChannelID   string `json:"controlPanelChannelID"`
	DiscordCharBufferSize   int    `json:"DiscordCharBufferSize"`
	IsDiscordEnabled        *bool  `json:"isDiscordEnabled"`
	ErrorChannelID          string `json:"errorChannelID"`

	// Backup Settings
	BackupsStoreDir      string        `json:"BackupsStoreDir"`
	BackupLoopActive     *bool         `json:"BackupLoopActive"`
	BackupLoopInterval   time.Duration `json:"BackupLoopInterval"`
	BackupMode           string        `json:"BackupMode"`
	BackupMaxFileSize    int64         `json:"BackupMaxFileSize"`
	BackupUseCompression *bool         `json:"BackupUseCompression"`
	BackupKeepSnapshot   *bool         `json:"BackupKeepSnapshot"`
}

// LoadConfig loads and initializes the configuration
func LoadConfig() (*JsonConfig, error) {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

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
		// File is missing, log it and proceed with defaults (probably first time setup)
		fmt.Println("config file was not found, proceeding with defaults.")
	} else {
		// Other errors (e.g., permissions), fail immediately
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	// Apply configuration
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
	RunfileIdentifier = getString(cfg.RunfileIdentifier, "RUNFILE_IDENTIFIER", "")

	isDiscordEnabledVal := getBool(cfg.IsDiscordEnabled, "IS_DISCORD_ENABLED", false)
	IsDiscordEnabled = isDiscordEnabledVal
	cfg.IsDiscordEnabled = &isDiscordEnabledVal

	ErrorChannelID = getString(cfg.ErrorChannelID, "ERROR_CHANNEL_ID", "")
	GameBranch = getString(cfg.GameBranch, "GAME_BRANCH", "public")

	LanguageSetting = getString(cfg.LanguageSetting, "LANGUAGE_SETTING", "en-US")
	BackendName = getString(cfg.BackendName, "SSUI_IDENTIFIER", "")
	BackendEndpointPort = getString(cfg.BackendEndpointPort, "BACKNED_ENDPOINT_PORT", "8443")

	Users = getUsers(cfg.Users, "SSUI_USERS", map[string]string{})
	RegisteredPlugins = getPlugins(cfg.RegisteredPlugins, "SSUI_REGISTERED_PLUGINS", map[string]string{})

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

	createGameServerLogFileVal := getBool(cfg.CreateGameServerLogFile, "CREATE_GAMESERVER_LOGFILE", false)
	CreateGameServerLogFile = createGameServerLogFileVal
	cfg.CreateGameServerLogFile = &createGameServerLogFileVal

	LogLevel = getInt(cfg.LogLevel, "LOG_LEVEL", 20)

	gameLogFromLogFileVal := getBool(cfg.GameLogFromLogFile, "GAME_LOG_FROM_LOG_FILE", false)
	GameLogFromLogFile = gameLogFromLogFileVal
	cfg.GameLogFromLogFile = &gameLogFromLogFileVal

	isUpdateEnabledVal := getBool(cfg.IsUpdateEnabled, "IS_UPDATE_ENABLED", true)
	IsUpdateEnabled = isUpdateEnabledVal
	cfg.IsUpdateEnabled = &isUpdateEnabledVal

	allowPrereleaseUpdatesVal := getBool(cfg.AllowPrereleaseUpdates, "ALLOW_PRERELEASE_UPDATES", false)
	AllowPrereleaseUpdates = allowPrereleaseUpdatesVal
	cfg.AllowPrereleaseUpdates = &allowPrereleaseUpdatesVal

	allowMajorUpdatesVal := getBool(cfg.AllowMajorUpdates, "ALLOW_MAJOR_UPDATES", false)
	AllowMajorUpdates = allowMajorUpdatesVal
	cfg.AllowMajorUpdates = &allowMajorUpdatesVal

	allowAutoGameServerUpdatesVal := getBool(cfg.AllowAutoGameServerUpdates, "ALLOW_AUTO_GAME_SERVER_UPDATES", false)
	AllowAutoGameServerUpdates = allowAutoGameServerUpdatesVal
	cfg.AllowAutoGameServerUpdates = &allowAutoGameServerUpdatesVal

	SubsystemFilters = getStringSlice(cfg.SubsystemFilters, "SUBSYSTEM_FILTERS", []string{})
	AutoRestartServerTimer = getString(cfg.AutoRestartServerTimer, "AUTO_RESTART_SERVER_TIMER", "0")

	isBepInExEnabledVal := getBool(cfg.IsBepInExEnabled, "IS_BEPINEX_ENABLED", false)
	IsBepInExEnabled = isBepInExEnabledVal
	cfg.IsBepInExEnabled = &isBepInExEnabledVal

	isSSCMEnabledVal := getBool(cfg.IsSSCMEnabled, "IS_SSCM_ENABLED", false)
	IsSSCMEnabled = isSSCMEnabledVal
	cfg.IsSSCMEnabled = &isSSCMEnabledVal

	isSSUICLIConsoleEnabledVal := getBool(cfg.IsSSUICLIConsoleEnabled, "IS_CONSOLE_ENABLED", true)
	IsSSUICLIConsoleEnabled = isSSUICLIConsoleEnabledVal
	cfg.IsSSUICLIConsoleEnabled = &isSSUICLIConsoleEnabledVal

	logClutterToConsoleVal := getBool(cfg.LogClutterToConsole, "LOG_CLUTTER_TO_CONSOLE", false)
	LogClutterToConsole = logClutterToConsoleVal
	cfg.LogClutterToConsole = &logClutterToConsoleVal

	autoStartServerOnStartupVal := getBool(cfg.AutoStartServerOnStartup, "AUTO_START_SERVER_ON_STARTUP", false)
	AutoStartServerOnStartup = autoStartServerOnStartupVal
	cfg.AutoStartServerOnStartup = &autoStartServerOnStartupVal

	// Backup Manager v3 Settings
	BackupsStoreDir = getString(cfg.BackupsStoreDir, "STORED_BACKUPS_DIR", SSUIFolder+"backups/storedBackups")
	BackupLoopInterval = getDuration(cfg.BackupLoopInterval, "BACKUP_LOOP_INTERVAL", 0*time.Hour)
	BackupMode = getString(cfg.BackupMode, "BACKUP_MODE", "tar")
	BackupMaxFileSize = getInt64(cfg.BackupMaxFileSize, "MAX_FILE_SIZE", 20*1024*1024*1024)
	backupUseCompressionVal := getBool(cfg.BackupUseCompression, "USE_COMPRESSION", true)
	BackupUseCompression = backupUseCompressionVal
	cfg.BackupUseCompression = &backupUseCompressionVal
	backupKeepSnapshotVal := getBool(cfg.BackupKeepSnapshot, "KEEP_SNAPSHOT", false)
	BackupKeepSnapshot = backupKeepSnapshotVal
	cfg.BackupKeepSnapshot = &backupKeepSnapshotVal
	enableBackupLoopVal := getBool(cfg.BackupLoopActive, "ENABLE_BACKUP_LOOP", false)
	BackupLoopActive = enableBackupLoopVal
	cfg.BackupLoopActive = &enableBackupLoopVal
}

// buildCurrentJsonConfig constructs JsonConfig from current runtime state
func buildCurrentJsonConfig() JsonConfig {
	return JsonConfig{
		DiscordToken:               DiscordToken,
		ControlChannelID:           ControlChannelID,
		StatusChannelID:            StatusChannelID,
		ConnectionListChannelID:    ConnectionListChannelID,
		LogChannelID:               LogChannelID,
		SaveChannelID:              SaveChannelID,
		ControlPanelChannelID:      ControlPanelChannelID,
		DiscordCharBufferSize:      DiscordCharBufferSize,
		IsDiscordEnabled:           &IsDiscordEnabled,
		ErrorChannelID:             ErrorChannelID,
		GameBranch:                 GameBranch,
		Users:                      Users,
		AuthEnabled:                &AuthEnabled,
		JwtKey:                     JwtKey,
		AuthTokenLifetime:          AuthTokenLifetime,
		Debug:                      &IsDebugMode,
		CreateSSUILogFile:          &CreateSSUILogFile,
		CreateGameServerLogFile:    &CreateGameServerLogFile,
		LogLevel:                   LogLevel,
		LogClutterToConsole:        &LogClutterToConsole,
		GameLogFromLogFile:         &GameLogFromLogFile,
		SubsystemFilters:           SubsystemFilters,
		IsUpdateEnabled:            &IsUpdateEnabled,
		IsSSCMEnabled:              &IsSSCMEnabled,
		IsBepInExEnabled:           &IsBepInExEnabled,
		AutoRestartServerTimer:     AutoRestartServerTimer,
		AllowPrereleaseUpdates:     &AllowPrereleaseUpdates,
		AllowMajorUpdates:          &AllowMajorUpdates,
		AllowAutoGameServerUpdates: &AllowAutoGameServerUpdates,
		IsSSUICLIConsoleEnabled:    &IsSSUICLIConsoleEnabled,
		LanguageSetting:            LanguageSetting,
		AutoStartServerOnStartup:   &AutoStartServerOnStartup,
		BackendName:                BackendName,
		BackendEndpointPort:        BackendEndpointPort,
		RunfileIdentifier:          RunfileIdentifier,
		RegisteredPlugins:          RegisteredPlugins,
		BackupsStoreDir:            BackupsStoreDir,
		BackupLoopInterval:         BackupLoopInterval,
		BackupMode:                 BackupMode,
		BackupMaxFileSize:          BackupMaxFileSize,
		BackupUseCompression:       &BackupUseCompression,
		BackupKeepSnapshot:         &BackupKeepSnapshot,
		BackupLoopActive:           &BackupLoopActive,
	}
}

// safeSaveConfigAtomic writes config atomically using temp file + rename
// MUST be called with ConfigMu locked
func safeSaveConfigAtomic() error {
	cfg := buildCurrentJsonConfig()

	// Create temp file in same directory to ensure atomic rename
	dir := filepath.Dir(ConfigPath)
	tmpFile, err := os.CreateTemp(dir, "config-*.json.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp config file: %v", err)
	}
	tmpPath := tmpFile.Name()

	// Encode with pretty print
	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to encode config: %v", err)
	}

	// Ensure data is flushed
	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to sync temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to close temp file: %v", err)
	}

	// Atomic replace
	if err := os.Rename(tmpPath, ConfigPath); err != nil {
		os.Remove(tmpPath) // clean up
		return fmt.Errorf("failed to rename temp file to config: %v", err)
	}

	return nil
}
