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
	Version = "5.7.1"
	Branch  = "release"
)

/*
If you read this, you are likely a developer. I sincerly apologize for the way the config works.
While I would love to refactor the config to not write to file then read the file every time a config value is changed,
I have not found the time to do so. So, for now, we save to file, then read the file and rely on whatever the file says. Although this is not ideal, it works for now. Deal with it.
JacksonTheMaster
*/

type JsonConfig struct {
	// reordered in 5.6.4 to simplify the order of the config file.

	// Gameserver Settings
	GameBranch       string `json:"gameBranch"`
	GamePort         string `json:"GamePort"`
	ServerName       string `json:"ServerName"`
	SaveInfo         string `json:"SaveInfo,omitempty"` // deprecated, kept for backwards compatibility
	SaveName         string `json:"SaveName"`           // replaces SaveInfo
	WorldID          string `json:"WorldID"`            // replaces SaveInfo
	ServerMaxPlayers string `json:"ServerMaxPlayers"`
	ServerPassword   string `json:"ServerPassword"`
	ServerAuthSecret string `json:"ServerAuthSecret"`
	AdminPassword    string `json:"AdminPassword"`
	UpdatePort       string `json:"UpdatePort"`
	UPNPEnabled      *bool  `json:"UPNPEnabled"`
	AutoSave         *bool  `json:"AutoSave"`
	SaveInterval     string `json:"SaveInterval"`
	AutoPauseServer  *bool  `json:"AutoPauseServer"`
	LocalIpAddress   string `json:"LocalIpAddress"`
	StartLocalHost   *bool  `json:"StartLocalHost"`
	ServerVisible    *bool  `json:"ServerVisible"`
	UseSteamP2P      *bool  `json:"UseSteamP2P"`
	AdditionalParams string `json:"AdditionalParams"`
	Difficulty       string `json:"Difficulty"`
	StartCondition   string `json:"StartCondition"`
	StartLocation    string `json:"StartLocation"`

	// Logging and debug settings
	Debug                *bool    `json:"Debug"`
	CreateSSUILogFile    *bool    `json:"CreateSSUILogFile"`
	LogLevel             int      `json:"LogLevel"`
	SubsystemFilters     []string `json:"subsystemFilters"`
	OverrideAdvertisedIp string   `json:"OverrideAdvertisedIp"`

	// Authentication Settings
	Users             map[string]string `json:"users"`       // Map of username to hashed password
	AuthEnabled       *bool             `json:"authEnabled"` // Toggle for enabling/disabling auth
	JwtKey            string            `json:"JwtKey"`
	AuthTokenLifetime int               `json:"AuthTokenLifetime"`

	// SSUI Settings
	IsNewTerrainAndSaveSystem *bool  `json:"IsNewTerrainAndSaveSystem"` // Use new terrain and save system
	ExePath                   string `json:"ExePath"`
	LogClutterToConsole       *bool  `json:"LogClutterToConsole"`
	IsSSCMEnabled             *bool  `json:"IsSSCMEnabled"`
	AutoRestartServerTimer    string `json:"AutoRestartServerTimer"`
	IsConsoleEnabled          *bool  `json:"IsConsoleEnabled"`
	LanguageSetting           string `json:"LanguageSetting"`
	AutoStartServerOnStartup  *bool  `json:"AutoStartServerOnStartup"`
	SSUIIdentifier            string `json:"SSUIIdentifier"`
	SSUIWebPort               string `json:"SSUIWebPort"`

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
	BlackListFilePath       string `json:"blackListFilePath"`
	IsDiscordEnabled        *bool  `json:"isDiscordEnabled"`
	ErrorChannelID          string `json:"errorChannelID"`

	//Backup Settings
	BackupKeepLastN       int   `json:"backupKeepLastN"`       // Number of most recent backups to keep (default: 2000)
	IsCleanupEnabled      *bool `json:"isCleanupEnabled"`      // Enable automatic cleanup of backups (default: false)
	BackupKeepDailyFor    int   `json:"backupKeepDailyFor"`    // Retention period in hours for daily backups
	BackupKeepWeeklyFor   int   `json:"backupKeepWeeklyFor"`   // Retention period in hours for weekly backups
	BackupKeepMonthlyFor  int   `json:"backupKeepMonthlyFor"`  // Retention period in hours for monthly backups
	BackupCleanupInterval int   `json:"backupCleanupInterval"` // Hours between backup cleanup operations
	BackupWaitTime        int   `json:"backupWaitTime"`        // Seconds to wait before copying backups
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
	discordToken = getString(cfg.DiscordToken, "DISCORD_TOKEN", "")
	controlChannelID = getString(cfg.ControlChannelID, "CONTROL_CHANNEL_ID", "")
	statusChannelID = getString(cfg.StatusChannelID, "STATUS_CHANNEL_ID", "")
	connectionListChannelID = getString(cfg.ConnectionListChannelID, "CONNECTION_LIST_CHANNEL_ID", "")
	logChannelID = getString(cfg.LogChannelID, "LOG_CHANNEL_ID", "")
	saveChannelID = getString(cfg.SaveChannelID, "SAVE_CHANNEL_ID", "")
	controlPanelChannelID = getString(cfg.ControlPanelChannelID, "CONTROL_PANEL_CHANNEL_ID", "")
	discordCharBufferSize = getInt(cfg.DiscordCharBufferSize, "DISCORD_CHAR_BUFFER_SIZE", 1000)
	blackListFilePath = getString(cfg.BlackListFilePath, "BLACKLIST_FILE_PATH", "./Blacklist.txt")

	isDiscordEnabledVal := getBool(cfg.IsDiscordEnabled, "IS_DISCORD_ENABLED", false)
	isDiscordEnabled = isDiscordEnabledVal
	cfg.IsDiscordEnabled = &isDiscordEnabledVal

	errorChannelID = getString(cfg.ErrorChannelID, "ERROR_CHANNEL_ID", "")
	backupKeepLastN = getInt(cfg.BackupKeepLastN, "BACKUP_KEEP_LAST_N", 2000)

	isCleanupEnabledVal := getBool(cfg.IsCleanupEnabled, "IS_CLEANUP_ENABLED", false)
	isCleanupEnabled = isCleanupEnabledVal
	cfg.IsCleanupEnabled = &isCleanupEnabledVal

	backupKeepDailyFor = time.Duration(getInt(cfg.BackupKeepDailyFor, "BACKUP_KEEP_DAILY_FOR", 24)) * time.Hour
	backupKeepWeeklyFor = time.Duration(getInt(cfg.BackupKeepWeeklyFor, "BACKUP_KEEP_WEEKLY_FOR", 168)) * time.Hour
	backupKeepMonthlyFor = time.Duration(getInt(cfg.BackupKeepMonthlyFor, "BACKUP_KEEP_MONTHLY_FOR", 730)) * time.Hour
	backupCleanupInterval = time.Duration(getInt(cfg.BackupCleanupInterval, "BACKUP_CLEANUP_INTERVAL", 730)) * time.Hour
	backupWaitTime = time.Duration(getInt(cfg.BackupWaitTime, "BACKUP_WAIT_TIME", 30)) * time.Second

	isNewTerrainAndSaveSystemVal := getBool(cfg.IsNewTerrainAndSaveSystem, "ENABLE_DOT_SAVES", true)
	isNewTerrainAndSaveSystem = isNewTerrainAndSaveSystemVal
	cfg.IsNewTerrainAndSaveSystem = &isNewTerrainAndSaveSystemVal

	gameBranch = getString(cfg.GameBranch, "GAME_BRANCH", "public")
	difficulty = getString(cfg.Difficulty, "DIFFICULTY", "")
	startCondition = getString(cfg.StartCondition, "START_CONDITION", "")
	startLocation = getString(cfg.StartLocation, "START_LOCATION", "")
	ServerName.value = getString(cfg.ServerName, "SERVER_NAME", "Stationeers Server UI")
	saveInfo = getString(cfg.SaveInfo, "SAVE_INFO", "") // deprecated, kept for backwards compatibility - if set, this gets migrated to SaveName and WorldID and the field is not written back to config.json
	saveName = getString(cfg.SaveName, "SAVE_NAME", "MyMapName")
	worldID = getString(cfg.WorldID, "WORLD_ID", "Lunar")
	ServerMaxPlayers.value = getString(cfg.ServerMaxPlayers, "SERVER_MAX_PLAYERS", "6")
	ServerPassword.value = getString(cfg.ServerPassword, "SERVER_PASSWORD", "")
	ServerAuthSecret.value = getString(cfg.ServerAuthSecret, "SERVER_AUTH_SECRET", "")
	AdminPassword.value = getString(cfg.AdminPassword, "ADMIN_PASSWORD", "")
	GamePort.value = getString(cfg.GamePort, "GAME_PORT", "27016")
	UpdatePort.value = getString(cfg.UpdatePort, "UPDATE_PORT", "27015")
	languageSetting = getString(cfg.LanguageSetting, "LANGUAGE_SETTING", "en-US")
	ssuiIdentifier = getString(cfg.SSUIIdentifier, "SSUI_IDENTIFIER", "")
	ssuiWebPort = getString(cfg.SSUIWebPort, "SSUI_WEB_PORT", "8443")

	upnpEnabledVal := getBool(cfg.UPNPEnabled, "UPNP_ENABLED", false)
	uPNPEnabled = upnpEnabledVal
	cfg.UPNPEnabled = &upnpEnabledVal

	autoSaveVal := getBool(cfg.AutoSave, "AUTO_SAVE", true)
	autoSave = autoSaveVal
	cfg.AutoSave = &autoSaveVal

	saveInterval = getString(cfg.SaveInterval, "SAVE_INTERVAL", "300")

	autoPauseServerVal := getBool(cfg.AutoPauseServer, "AUTO_PAUSE_SERVER", true)
	autoPauseServer = autoPauseServerVal
	cfg.AutoPauseServer = &autoPauseServerVal

	LocalIpAddress.value = getString(cfg.LocalIpAddress, "LOCAL_IP_ADDRESS", "0.0.0.0")

	startLocalHostVal := getBool(cfg.StartLocalHost, "START_LOCAL_HOST", true)
	startLocalHost = startLocalHostVal
	cfg.StartLocalHost = &startLocalHostVal

	serverVisibleVal := getBool(cfg.ServerVisible, "SERVER_VISIBLE", true)
	ServerVisible.value = serverVisibleVal
	cfg.ServerVisible = &serverVisibleVal

	useSteamP2PVal := getBool(cfg.UseSteamP2P, "USE_STEAM_P2P", false)
	useSteamP2P = useSteamP2PVal
	cfg.UseSteamP2P = &useSteamP2PVal

	exePath = getString(cfg.ExePath, "EXE_PATH", getDefaultExePath())
	additionalParams = getString(cfg.AdditionalParams, "ADDITIONAL_PARAMS", "")
	users = getUsers(cfg.Users, "SSUI_USERS", map[string]string{})

	authEnabledVal := getBool(cfg.AuthEnabled, "SSUI_AUTH_ENABLED", false)
	authEnabled = authEnabledVal
	cfg.AuthEnabled = &authEnabledVal

	jwtKey = getString(cfg.JwtKey, "SSUI_JWT_KEY", generateJwtKey())
	authTokenLifetime = getInt(cfg.AuthTokenLifetime, "SSUI_AUTH_TOKEN_LIFETIME", 1440)

	debugVal := getBool(cfg.Debug, "DEBUG", false)
	isDebugMode = debugVal
	cfg.Debug = &debugVal

	createSSUILogFileVal := getBool(cfg.CreateSSUILogFile, "CREATE_SSUI_LOGFILE", false)
	createSSUILogFile = createSSUILogFileVal
	cfg.CreateSSUILogFile = &createSSUILogFileVal

	logLevel = getInt(cfg.LogLevel, "LOG_LEVEL", 20)

	isUpdateEnabledVal := getBool(cfg.IsUpdateEnabled, "IS_UPDATE_ENABLED", true)
	isUpdateEnabled = isUpdateEnabledVal
	cfg.IsUpdateEnabled = &isUpdateEnabledVal

	allowPrereleaseUpdatesVal := getBool(cfg.AllowPrereleaseUpdates, "ALLOW_PRERELEASE_UPDATES", false)
	allowPrereleaseUpdates = allowPrereleaseUpdatesVal
	cfg.AllowPrereleaseUpdates = &allowPrereleaseUpdatesVal

	allowMajorUpdatesVal := getBool(cfg.AllowMajorUpdates, "ALLOW_MAJOR_UPDATES", false)
	allowMajorUpdates = allowMajorUpdatesVal
	cfg.AllowMajorUpdates = &allowMajorUpdatesVal

	allowAutoGameServerUpdatesVal := getBool(cfg.AllowAutoGameServerUpdates, "ALLOW_AUTO_GAME_SERVER_UPDATES", false)
	allowAutoGameServerUpdates = allowAutoGameServerUpdatesVal
	cfg.AllowAutoGameServerUpdates = &allowAutoGameServerUpdatesVal

	subsystemFilters = getStringSlice(cfg.SubsystemFilters, "SUBSYSTEM_FILTERS", []string{})
	autoRestartServerTimer = getString(cfg.AutoRestartServerTimer, "AUTO_RESTART_SERVER_TIMER", "0")
	isSSCMEnabledVal := getBool(cfg.IsSSCMEnabled, "IS_SSCM_ENABLED", true)
	isSSCMEnabled = isSSCMEnabledVal
	cfg.IsSSCMEnabled = &isSSCMEnabledVal

	isConsoleEnabledVal := getBool(cfg.IsConsoleEnabled, "IS_CONSOLE_ENABLED", true)
	isConsoleEnabled = isConsoleEnabledVal
	cfg.IsConsoleEnabled = &isConsoleEnabledVal

	logClutterToConsoleVal := getBool(cfg.LogClutterToConsole, "LOG_CLUTTER_TO_CONSOLE", false)
	logClutterToConsole = logClutterToConsoleVal
	cfg.LogClutterToConsole = &logClutterToConsoleVal

	autoStartServerOnStartupVal := getBool(cfg.AutoStartServerOnStartup, "AUTO_START_SERVER_ON_STARTUP", false)
	autoStartServerOnStartup = autoStartServerOnStartupVal
	cfg.AutoStartServerOnStartup = &autoStartServerOnStartupVal

	// Process SaveInfo to maintain backwards compatibility with pre-5.6.6 SaveInfo field (deprecated)
	if saveInfo != "" {
		parts := strings.Split(saveInfo, " ")
		if len(parts) > 0 {
			saveName = parts[0]
			fmt.Println("SaveName: " + saveName)
		}
		if len(parts) > 1 {
			worldID = parts[1]
			fmt.Println("WorldID: " + worldID)
		}
		cfg.SaveInfo = ""
	}

	if gameBranch != "public" && gameBranch != "beta" {
		isNewTerrainAndSaveSystem = false
	} else {
		isNewTerrainAndSaveSystem = true
	}

	// Set backup paths for old or new style saves
	if isNewTerrainAndSaveSystem {
		// use new new style autosave folder
		configuredBackupDir = filepath.Join("./saves/", saveName, "autosave")
	} else {
		// use old style Backups folder
		configuredBackupDir = filepath.Join("./saves/", saveName, "Backup")
	}
	// use Safebackups folder either way.
	configuredSafeBackupDir = filepath.Join("./saves/", saveName, "Safebackups")

	overrideAdvertisedIp = getString(cfg.OverrideAdvertisedIp, "OVERRIDE_ADVERTISED_IP", "")

	safeSaveConfig()
}

// use safeSaveConfig EXCLUSIVELY though setter functions
// M U S T be called while holding a lock on ConfigMu!
func safeSaveConfig() error {
	cfg := JsonConfig{
		DiscordToken:               discordToken,
		ControlChannelID:           controlChannelID,
		StatusChannelID:            statusChannelID,
		ConnectionListChannelID:    connectionListChannelID,
		LogChannelID:               logChannelID,
		SaveChannelID:              saveChannelID,
		ControlPanelChannelID:      controlPanelChannelID,
		DiscordCharBufferSize:      discordCharBufferSize,
		BlackListFilePath:          blackListFilePath,
		IsDiscordEnabled:           &isDiscordEnabled,
		ErrorChannelID:             errorChannelID,
		BackupKeepLastN:            backupKeepLastN,
		IsCleanupEnabled:           &isCleanupEnabled,
		BackupKeepDailyFor:         int(backupKeepDailyFor / time.Hour),    // Convert to hours
		BackupKeepWeeklyFor:        int(backupKeepWeeklyFor / time.Hour),   // Convert to hours
		BackupKeepMonthlyFor:       int(backupKeepMonthlyFor / time.Hour),  // Convert to hours
		BackupCleanupInterval:      int(backupCleanupInterval / time.Hour), // Convert to hours
		BackupWaitTime:             int(backupWaitTime / time.Second),      // Convert to seconds
		IsNewTerrainAndSaveSystem:  &isNewTerrainAndSaveSystem,
		GameBranch:                 gameBranch,
		Difficulty:                 difficulty,
		StartCondition:             startCondition,
		StartLocation:              startLocation,
		ServerName:                 ServerName.value,
		SaveName:                   saveName,
		WorldID:                    worldID,
		ServerMaxPlayers:           ServerMaxPlayers.value,
		ServerPassword:             ServerPassword.value,
		ServerAuthSecret:           ServerAuthSecret.value,
		AdminPassword:              AdminPassword.value,
		GamePort:                   GamePort.value,
		UpdatePort:                 UpdatePort.value,
		UPNPEnabled:                &uPNPEnabled,
		AutoSave:                   &autoSave,
		SaveInterval:               saveInterval,
		AutoPauseServer:            &autoPauseServer,
		LocalIpAddress:             LocalIpAddress.value,
		StartLocalHost:             &startLocalHost,
		ServerVisible:              &ServerVisible.value,
		UseSteamP2P:                &useSteamP2P,
		ExePath:                    exePath,
		AdditionalParams:           additionalParams,
		Users:                      users,
		AuthEnabled:                &authEnabled,
		JwtKey:                     jwtKey,
		AuthTokenLifetime:          authTokenLifetime,
		Debug:                      &isDebugMode,
		CreateSSUILogFile:          &createSSUILogFile,
		LogLevel:                   logLevel,
		LogClutterToConsole:        &logClutterToConsole,
		SubsystemFilters:           subsystemFilters,
		IsUpdateEnabled:            &isUpdateEnabled,
		IsSSCMEnabled:              &isSSCMEnabled,
		AutoRestartServerTimer:     autoRestartServerTimer,
		AllowPrereleaseUpdates:     &allowPrereleaseUpdates,
		AllowMajorUpdates:          &allowMajorUpdates,
		AllowAutoGameServerUpdates: &allowAutoGameServerUpdates,
		IsConsoleEnabled:           &isConsoleEnabled,
		LanguageSetting:            languageSetting,
		AutoStartServerOnStartup:   &autoStartServerOnStartup,
		SSUIIdentifier:             ssuiIdentifier,
		SSUIWebPort:                ssuiWebPort,
		OverrideAdvertisedIp:       overrideAdvertisedIp,
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
