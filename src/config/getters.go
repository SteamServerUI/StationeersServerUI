package config

import "time"

func GetDiscordToken() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return DiscordToken
}

func GetControlChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ControlChannelID
}

func GetStatusChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return StatusChannelID
}

func GetConnectionListChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ConnectionListChannelID
}

func GetLogChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LogChannelID
}

func GetSaveChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SaveChannelID
}

func GetControlPanelChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ControlPanelChannelID
}

func GetDiscordCharBufferSize() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return DiscordCharBufferSize
}

func GetBlackListFilePath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BlackListFilePath
}

func GetIsDiscordEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsDiscordEnabled
}

func GetErrorChannelID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ErrorChannelID
}

func GetBackupKeepLastN() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupKeepLastN
}

func GetIsCleanupEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsCleanupEnabled
}

// GetBackupKeepDailyFor returns the retention period for daily backups in hours.
func GetBackupKeepDailyFor() time.Duration {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupKeepDailyFor
}

func GetBackupKeepWeeklyFor() time.Duration {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupKeepWeeklyFor
}

// GetBackupKeepMonthlyFor returns the retention period for monthly backups in hours.
func GetBackupKeepMonthlyFor() time.Duration {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupKeepMonthlyFor
}

// GetBackupCleanupInterval returns the cleanup interval in hours.
func GetBackupCleanupInterval() time.Duration {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupCleanupInterval
}

func GetIsNewTerrainAndSaveSystem() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsNewTerrainAndSaveSystem
}

func GetGameBranch() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return GameBranch
}

func GetDifficulty() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return Difficulty
}

func GetStartCondition() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return StartCondition
}

func GetStartLocation() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return StartLocation
}

func GetServerName() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ServerName
}

func GetSaveInfo() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SaveInfo
}

func GetWorldName() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return WorldName
}

func GetBackupWorldName() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BackupWorldName
}

func GetServerMaxPlayers() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ServerMaxPlayers
}

func GetServerPassword() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ServerPassword
}

func GetServerAuthSecret() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ServerAuthSecret
}

func GetAdminPassword() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AdminPassword
}

func GetGamePort() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return GamePort
}

func GetUpdatePort() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return UpdatePort
}

func GetUPNPEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return UPNPEnabled
}

func GetAutoSave() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AutoSave
}

func GetSaveInterval() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SaveInterval
}

func GetAutoPauseServer() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AutoPauseServer
}

func GetLocalIpAddress() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LocalIpAddress
}

func GetStartLocalHost() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return StartLocalHost
}

func GetServerVisible() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ServerVisible
}

func GetUseSteamP2P() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return UseSteamP2P
}

func GetExePath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ExePath
}

func GetAdditionalParams() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AdditionalParams
}

func GetUsers() map[string]string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return Users
}

func GetAuthEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AuthEnabled
}

func GetJwtKey() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return JwtKey
}

func GetAuthTokenLifetime() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AuthTokenLifetime
}

func GetIsDebugMode() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsDebugMode
}

func GetCreateSSUILogFile() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return CreateSSUILogFile
}

func GetLogLevel() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LogLevel
}

func GetLogClutterToConsole() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LogClutterToConsole
}

func GetSubsystemFilters() []string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SubsystemFilters
}

func GetIsUpdateEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsUpdateEnabled
}

func GetIsSSCMEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsSSCMEnabled
}

func GetSSCMFilePath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSCMFilePath
}

func GetSSCMPluginDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSCMPluginDir
}

func GetSSCMWebDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSCMWebDir
}

func GetAutoRestartServerTimer() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AutoRestartServerTimer
}

func GetAllowPrereleaseUpdates() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AllowPrereleaseUpdates
}

func GetAllowMajorUpdates() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AllowMajorUpdates
}

func GetIsConsoleEnabled() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsConsoleEnabled
}

func GetLanguageSetting() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LanguageSetting
}

func GetAutoStartServerOnStartup() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return AutoStartServerOnStartup
}

func GetSSUIIdentifier() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSUIIdentifier
}

func GetSSUIWebPort() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSUIWebPort
}

// GetIsFirstTimeSetup returns the IsFirstTimeSetup
func GetIsFirstTimeSetup() bool {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return IsFirstTimeSetup
}

func GetConfigPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ConfigPath
}

func GetVersion() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return Version
}

func GetBranch() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return Branch
}

func GetTLSCertPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return TLSCertPath
}

func GetTLSKeyPath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return TLSKeyPath
}

func GetUIModFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return UIModFolder
}

func GetMaxSSEConnections() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return MaxSSEConnections
}

func GetSSEMessageBufferSize() int {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return SSEMessageBufferSize
}

func GetBufferFlushTicker() *time.Ticker {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return BufferFlushTicker
}

func GetLogFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return LogFolder
}

func GetConfiguredBackupDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ConfiguredBackupDir
}

func GetConfiguredSafeBackupDir() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return ConfiguredSafeBackupDir
}

func GetCustomDetectionsFilePath() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return CustomDetectionsFilePath
}

func GetGameServerAppID() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return GameServerAppID
}
