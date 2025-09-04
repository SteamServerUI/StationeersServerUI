package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// All configuration variables can be found in vars.go
	Version = "5.6.2"
	Branch  = "release"
)

type JsonConfig struct {
	DiscordToken              string            `json:"discordToken"`
	ControlChannelID          string            `json:"controlChannelID"`
	StatusChannelID           string            `json:"statusChannelID"`
	ConnectionListChannelID   string            `json:"connectionListChannelID"`
	LogChannelID              string            `json:"logChannelID"`
	SaveChannelID             string            `json:"saveChannelID"`
	ControlPanelChannelID     string            `json:"controlPanelChannelID"`
	DiscordCharBufferSize     int               `json:"DiscordCharBufferSize"`
	BlackListFilePath         string            `json:"blackListFilePath"`
	IsDiscordEnabled          *bool             `json:"isDiscordEnabled"`
	ErrorChannelID            string            `json:"errorChannelID"`
	BackupKeepLastN           int               `json:"backupKeepLastN"`           // Number of most recent backups to keep (default: 2000)
	IsCleanupEnabled          *bool             `json:"isCleanupEnabled"`          // Enable automatic cleanup of backups (default: false)
	BackupKeepDailyFor        int               `json:"backupKeepDailyFor"`        // Retention period in hours for daily backups
	BackupKeepWeeklyFor       int               `json:"backupKeepWeeklyFor"`       // Retention period in hours for weekly backups
	BackupKeepMonthlyFor      int               `json:"backupKeepMonthlyFor"`      // Retention period in hours for monthly backups
	BackupCleanupInterval     int               `json:"backupCleanupInterval"`     // Hours between backup cleanup operations
	BackupWaitTime            int               `json:"backupWaitTime"`            // Seconds to wait before copying backups
	IsNewTerrainAndSaveSystem *bool             `json:"IsNewTerrainAndSaveSystem"` // Use new terrain and save system
	GameBranch                string            `json:"gameBranch"`
	Difficulty                string            `json:"Difficulty"`
	StartCondition            string            `json:"StartCondition"`
	StartLocation             string            `json:"StartLocation"`
	ServerName                string            `json:"ServerName"`
	SaveInfo                  string            `json:"SaveInfo"`
	ServerMaxPlayers          string            `json:"ServerMaxPlayers"`
	ServerPassword            string            `json:"ServerPassword"`
	ServerAuthSecret          string            `json:"ServerAuthSecret"`
	AdminPassword             string            `json:"AdminPassword"`
	GamePort                  string            `json:"GamePort"`
	UpdatePort                string            `json:"UpdatePort"`
	UPNPEnabled               *bool             `json:"UPNPEnabled"`
	AutoSave                  *bool             `json:"AutoSave"`
	SaveInterval              string            `json:"SaveInterval"`
	AutoPauseServer           *bool             `json:"AutoPauseServer"`
	LocalIpAddress            string            `json:"LocalIpAddress"`
	StartLocalHost            *bool             `json:"StartLocalHost"`
	ServerVisible             *bool             `json:"ServerVisible"`
	UseSteamP2P               *bool             `json:"UseSteamP2P"`
	ExePath                   string            `json:"ExePath"`
	AdditionalParams          string            `json:"AdditionalParams"`
	Users                     map[string]string `json:"users"`       // Map of username to hashed password
	AuthEnabled               *bool             `json:"authEnabled"` // Toggle for enabling/disabling auth
	JwtKey                    string            `json:"JwtKey"`
	AuthTokenLifetime         int               `json:"AuthTokenLifetime"`
	Debug                     *bool             `json:"Debug"`
	CreateSSUILogFile         *bool             `json:"CreateSSUILogFile"`
	LogLevel                  int               `json:"LogLevel"`
	LogClutterToConsole       *bool             `json:"LogClutterToConsole"`
	SubsystemFilters          []string          `json:"subsystemFilters"`
	IsUpdateEnabled           *bool             `json:"IsUpdateEnabled"`
	IsSSCMEnabled             *bool             `json:"IsSSCMEnabled"`
	AutoRestartServerTimer    string            `json:"AutoRestartServerTimer"`
	AllowPrereleaseUpdates    *bool             `json:"AllowPrereleaseUpdates"`
	AllowMajorUpdates         *bool             `json:"AllowMajorUpdates"`
	IsConsoleEnabled          *bool             `json:"IsConsoleEnabled"`
	LanguageSetting           string            `json:"LanguageSetting"`
	AutoStartServerOnStartup  *bool             `json:"AutoStartServerOnStartup"`
	SSUIIdentifier            string            `json:"SSUIIdentifier"`
	SSUIWebPort               string            `json:"SSUIWebPort"`
}

type CustomDetection struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Pattern   string `json:"pattern"`
	EventType string `json:"eventType"`
	Message   string `json:"message"`
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
	BlackListFilePath = getString(cfg.BlackListFilePath, "BLACKLIST_FILE_PATH", "./Blacklist.txt")

	isDiscordEnabledVal := getBool(cfg.IsDiscordEnabled, "IS_DISCORD_ENABLED", false)
	IsDiscordEnabled = isDiscordEnabledVal
	cfg.IsDiscordEnabled = &isDiscordEnabledVal

	ErrorChannelID = getString(cfg.ErrorChannelID, "ERROR_CHANNEL_ID", "")
	BackupKeepLastN = getInt(cfg.BackupKeepLastN, "BACKUP_KEEP_LAST_N", 2000)

	isCleanupEnabledVal := getBool(cfg.IsCleanupEnabled, "IS_CLEANUP_ENABLED", false)
	IsCleanupEnabled = isCleanupEnabledVal
	cfg.IsCleanupEnabled = &isCleanupEnabledVal

	BackupKeepDailyFor = time.Duration(getInt(cfg.BackupKeepDailyFor, "BACKUP_KEEP_DAILY_FOR", 24)) * time.Hour
	BackupKeepWeeklyFor = time.Duration(getInt(cfg.BackupKeepWeeklyFor, "BACKUP_KEEP_WEEKLY_FOR", 168)) * time.Hour
	BackupKeepMonthlyFor = time.Duration(getInt(cfg.BackupKeepMonthlyFor, "BACKUP_KEEP_MONTHLY_FOR", 730)) * time.Hour
	BackupCleanupInterval = time.Duration(getInt(cfg.BackupCleanupInterval, "BACKUP_CLEANUP_INTERVAL", 730)) * time.Hour
	BackupWaitTime = time.Duration(getInt(cfg.BackupWaitTime, "BACKUP_WAIT_TIME", 30)) * time.Second

	isNewTerrainAndSaveSystemVal := getBool(cfg.IsNewTerrainAndSaveSystem, "ENABLE_DOT_SAVES", false)
	IsNewTerrainAndSaveSystem = isNewTerrainAndSaveSystemVal
	cfg.IsNewTerrainAndSaveSystem = &isNewTerrainAndSaveSystemVal

	GameBranch = getString(cfg.GameBranch, "GAME_BRANCH", "public")
	Difficulty = getString(cfg.Difficulty, "DIFFICULTY", "")
	StartCondition = getString(cfg.StartCondition, "START_CONDITION", "")
	StartLocation = getString(cfg.StartLocation, "START_LOCATION", "")
	ServerName = getString(cfg.ServerName, "SERVER_NAME", "Stationeers Server UI")
	SaveInfo = getString(cfg.SaveInfo, "SAVE_INFO", "Vulcan Vulcan")
	ServerMaxPlayers = getString(cfg.ServerMaxPlayers, "SERVER_MAX_PLAYERS", "6")
	ServerPassword = getString(cfg.ServerPassword, "SERVER_PASSWORD", "")
	ServerAuthSecret = getString(cfg.ServerAuthSecret, "SERVER_AUTH_SECRET", "")
	AdminPassword = getString(cfg.AdminPassword, "ADMIN_PASSWORD", "")
	GamePort = getString(cfg.GamePort, "GAME_PORT", "27016")
	UpdatePort = getString(cfg.UpdatePort, "UPDATE_PORT", "27015")
	LanguageSetting = getString(cfg.LanguageSetting, "LANGUAGE_SETTING", "en-US")
	SSUIIdentifier = getString(cfg.SSUIIdentifier, "SSUI_IDENTIFIER", "")
	SSUIWebPort = getString(cfg.SSUIWebPort, "SSUI_WEB_PORT", "8443")

	upnpEnabledVal := getBool(cfg.UPNPEnabled, "UPNP_ENABLED", false)
	UPNPEnabled = upnpEnabledVal
	cfg.UPNPEnabled = &upnpEnabledVal

	autoSaveVal := getBool(cfg.AutoSave, "AUTO_SAVE", true)
	AutoSave = autoSaveVal
	cfg.AutoSave = &autoSaveVal

	SaveInterval = getString(cfg.SaveInterval, "SAVE_INTERVAL", "300")

	autoPauseServerVal := getBool(cfg.AutoPauseServer, "AUTO_PAUSE_SERVER", true)
	AutoPauseServer = autoPauseServerVal
	cfg.AutoPauseServer = &autoPauseServerVal

	LocalIpAddress = getString(cfg.LocalIpAddress, "LOCAL_IP_ADDRESS", "")

	startLocalHostVal := getBool(cfg.StartLocalHost, "START_LOCAL_HOST", true)
	StartLocalHost = startLocalHostVal
	cfg.StartLocalHost = &startLocalHostVal

	serverVisibleVal := getBool(cfg.ServerVisible, "SERVER_VISIBLE", true)
	ServerVisible = serverVisibleVal
	cfg.ServerVisible = &serverVisibleVal

	useSteamP2PVal := getBool(cfg.UseSteamP2P, "USE_STEAM_P2P", true)
	UseSteamP2P = useSteamP2PVal
	cfg.UseSteamP2P = &useSteamP2PVal

	ExePath = getString(cfg.ExePath, "EXE_PATH", getDefaultExePath())
	AdditionalParams = getString(cfg.AdditionalParams, "ADDITIONAL_PARAMS", "")
	Users = getUsers(cfg.Users, "SSUI_USERS", map[string]string{})

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

	LogLevel = getInt(cfg.LogLevel, "LOG_LEVEL", 20)

	isUpdateEnabledVal := getBool(cfg.IsUpdateEnabled, "IS_UPDATE_ENABLED", true)
	IsUpdateEnabled = isUpdateEnabledVal
	cfg.IsUpdateEnabled = &isUpdateEnabledVal

	allowPrereleaseUpdatesVal := getBool(cfg.AllowPrereleaseUpdates, "ALLOW_PRERELEASE_UPDATES", false)
	AllowPrereleaseUpdates = allowPrereleaseUpdatesVal
	cfg.AllowPrereleaseUpdates = &allowPrereleaseUpdatesVal

	allowMajorUpdatesVal := getBool(cfg.AllowMajorUpdates, "ALLOW_MAJOR_UPDATES", false)
	AllowMajorUpdates = allowMajorUpdatesVal
	cfg.AllowMajorUpdates = &allowMajorUpdatesVal

	SubsystemFilters = getStringSlice(cfg.SubsystemFilters, "SUBSYSTEM_FILTERS", []string{})
	AutoRestartServerTimer = getString(cfg.AutoRestartServerTimer, "AUTO_RESTART_SERVER_TIMER", "0")
	isSSCMEnabledVal := getBool(cfg.IsSSCMEnabled, "IS_SSCM_ENABLED", true)
	IsSSCMEnabled = isSSCMEnabledVal
	cfg.IsSSCMEnabled = &isSSCMEnabledVal

	isConsoleEnabledVal := getBool(cfg.IsConsoleEnabled, "IS_CONSOLE_ENABLED", true)
	IsConsoleEnabled = isConsoleEnabledVal
	cfg.IsConsoleEnabled = &isConsoleEnabledVal

	logClutterToConsoleVal := getBool(cfg.LogClutterToConsole, "LOG_CLUTTER_TO_CONSOLE", false)
	LogClutterToConsole = logClutterToConsoleVal
	cfg.LogClutterToConsole = &logClutterToConsoleVal

	autoStartServerOnStartupVal := getBool(cfg.AutoStartServerOnStartup, "AUTO_START_SERVER_ON_STARTUP", false)
	AutoStartServerOnStartup = autoStartServerOnStartupVal
	cfg.AutoStartServerOnStartup = &autoStartServerOnStartupVal

	// Process SaveInfo
	parts := strings.Split(SaveInfo, " ")
	if len(parts) > 0 {
		WorldName = parts[0]
	}
	if len(parts) > 1 {
		BackupWorldName = parts[1]
	}

	// Set backup paths for old or new style saves
	if IsNewTerrainAndSaveSystem {
		// use new new style autosave folder
		ConfiguredBackupDir = filepath.Join("./saves/", WorldName, "autosave")
	} else {
		// use old style Backups folder
		ConfiguredBackupDir = filepath.Join("./saves/", WorldName, "Backup")
	}
	// use Safebackups folder either way.
	ConfiguredSafeBackupDir = filepath.Join("./saves/", WorldName, "Safebackups")
}

// use safeSaveConfig EXCLUSIVELY though setter functions
// M U S T be called while holding a lock on ConfigMu!
func safeSaveConfig() error {

	cfg := JsonConfig{
		DiscordToken:              DiscordToken,
		ControlChannelID:          ControlChannelID,
		StatusChannelID:           StatusChannelID,
		ConnectionListChannelID:   ConnectionListChannelID,
		LogChannelID:              LogChannelID,
		SaveChannelID:             SaveChannelID,
		ControlPanelChannelID:     ControlPanelChannelID,
		DiscordCharBufferSize:     DiscordCharBufferSize,
		BlackListFilePath:         BlackListFilePath,
		IsDiscordEnabled:          &IsDiscordEnabled,
		ErrorChannelID:            ErrorChannelID,
		BackupKeepLastN:           BackupKeepLastN,
		IsCleanupEnabled:          &IsCleanupEnabled,
		BackupKeepDailyFor:        int(BackupKeepDailyFor / time.Hour),    // Convert to hours
		BackupKeepWeeklyFor:       int(BackupKeepWeeklyFor / time.Hour),   // Convert to hours
		BackupKeepMonthlyFor:      int(BackupKeepMonthlyFor / time.Hour),  // Convert to hours
		BackupCleanupInterval:     int(BackupCleanupInterval / time.Hour), // Convert to hours
		BackupWaitTime:            int(BackupWaitTime / time.Second),      // Convert to seconds
		IsNewTerrainAndSaveSystem: &IsNewTerrainAndSaveSystem,
		GameBranch:                GameBranch,
		Difficulty:                Difficulty,
		StartCondition:            StartCondition,
		StartLocation:             StartLocation,
		ServerName:                ServerName,
		SaveInfo:                  SaveInfo,
		ServerMaxPlayers:          ServerMaxPlayers,
		ServerPassword:            ServerPassword,
		ServerAuthSecret:          ServerAuthSecret,
		AdminPassword:             AdminPassword,
		GamePort:                  GamePort,
		UpdatePort:                UpdatePort,
		UPNPEnabled:               &UPNPEnabled,
		AutoSave:                  &AutoSave,
		SaveInterval:              SaveInterval,
		AutoPauseServer:           &AutoPauseServer,
		LocalIpAddress:            LocalIpAddress,
		StartLocalHost:            &StartLocalHost,
		ServerVisible:             &ServerVisible,
		UseSteamP2P:               &UseSteamP2P,
		ExePath:                   ExePath,
		AdditionalParams:          AdditionalParams,
		Users:                     Users,
		AuthEnabled:               &AuthEnabled,
		JwtKey:                    JwtKey,
		AuthTokenLifetime:         AuthTokenLifetime,
		Debug:                     &IsDebugMode,
		CreateSSUILogFile:         &CreateSSUILogFile,
		LogLevel:                  LogLevel,
		LogClutterToConsole:       &LogClutterToConsole,
		SubsystemFilters:          SubsystemFilters,
		IsUpdateEnabled:           &IsUpdateEnabled,
		IsSSCMEnabled:             &IsSSCMEnabled,
		AutoRestartServerTimer:    AutoRestartServerTimer,
		AllowPrereleaseUpdates:    &AllowPrereleaseUpdates,
		AllowMajorUpdates:         &AllowMajorUpdates,
		IsConsoleEnabled:          &IsConsoleEnabled,
		LanguageSetting:           LanguageSetting,
		AutoStartServerOnStartup:  &AutoStartServerOnStartup,
		SSUIIdentifier:            SSUIIdentifier,
		SSUIWebPort:               SSUIWebPort,
	}

	file, err := os.Create(ConfigPath)
	if err != nil {
		return fmt.Errorf("error creating config.json: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("error encoding config.json: %v", err)
	}

	return nil
}

// use SaveConfig EXCLUSIVELY though loader.SaveConfig to trigger a reload afterwards!
// when the config gets updated, changes do not get reflected at runtime UNLESS a backend reload / config reload is triggered
// This can be done via configchanger.SaveConfig
func SaveConfigToFile(cfg *JsonConfig) error {

	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	file, err := os.Create(ConfigPath)
	if err != nil {
		return fmt.Errorf("error creating config.json: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("error encoding config.json: %v", err)
	}

	return nil
}
