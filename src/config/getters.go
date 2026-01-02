package config

import "time"

func GetDiscordToken() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return DiscordToken
}

func GetControlChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return ControlChannelID
}

func GetStatusChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return StatusChannelID
}

func GetConnectionListChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return ConnectionListChannelID
}

func GetLogChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return LogChannelID
}

func GetSaveChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return SaveChannelID
}

func GetControlPanelChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return ControlPanelChannelID
}

func GetDiscordCharBufferSize() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return DiscordCharBufferSize
}

func GetIsDiscordEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsDiscordEnabled
}

func GetErrorChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return ErrorChannelID
}

func GetGameBranch() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return GameBranch
}

func GetUsers() map[string]string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return Users
}

func GetAuthEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return AuthEnabled
}

func GetJwtKey() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return JwtKey
}

func GetAuthTokenLifetime() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return AuthTokenLifetime
}

func GetIsDebugMode() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsDebugMode
}

func GetCreateSSUILogFile() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return CreateSSUILogFile
}

func GetLogLevel() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return LogLevel
}

func GetLogClutterToConsole() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return LogClutterToConsole
}

func GetSubsystemFilters() []string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return SubsystemFilters
}

func GetIsUpdateEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsUpdateEnabled
}

func GetIsSSCMEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsSSCMEnabled
}

func GetIsBepInExEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsBepInExEnabled
}

func GetSSCMFilePath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return SSCMFilePath
}

func GetSSCMPluginDir() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return SSCMPluginDir
}

func GetSSCMWebDir() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return SSCMWebDir
}

func GetAutoRestartServerTimer() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return AutoRestartServerTimer
}

func GetAllowPrereleaseUpdates() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return AllowPrereleaseUpdates
}

func GetAllowMajorUpdates() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return AllowMajorUpdates
}

func GetIsSSUICLIConsoleEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsSSUICLIConsoleEnabled
}

func GetLanguageSetting() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return LanguageSetting
}

func GetAutoStartServerOnStartup() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return AutoStartServerOnStartup
}

func GetBackendName() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return BackendName
}

func GetBackendEndpointPort() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return BackendEndpointPort
}

// GetIsFirstTimeSetup returns the IsFirstTimeSetup
func GetIsFirstTimeSetup() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsFirstTimeSetup
}

func GetConfigPath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return ConfigPath
}

func GetVersion() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return Version
}

func GetBranch() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return Branch
}

func GetTLSCertPath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return TLSCertPath
}

func GetTLSKeyPath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return TLSKeyPath
}

func GetSSUIFolder() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return SSUIFolder
}

func GetMaxSSEConnections() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return MaxSSEConnections
}

func GetSSEMessageBufferSize() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return SSEMessageBufferSize
}

func GetLogFolder() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return LogFolder
}

func GetCustomDetectionsFilePath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return CustomDetectionsFilePath
}

func GetGameServerAppID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return GameServerAppID
}

func GetCurrentBranchBuildID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return CurrentBranchBuildID
}

func GetAllowAutoGameServerUpdates() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return AllowAutoGameServerUpdates
}

func GetExtractedGameVersion() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return ExtractedGameVersion
}

func GetSkipSteamCMD() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return SkipSteamCMD
}

func GetNoSanityCheck() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return NoSanityCheck
}

func GetIsDockerContainer() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsDockerContainer
}

func GetRunfilesFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunFilesFolder
}

// GetRunFilesFolder returns the RunFilesFolder
func GetRunFilesFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunFilesFolder
}

func GetPluginsFolder() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return PluginsFolder
}

// GetRunfileGame returns the RunfileGame
func GetRunfileIdentifier() string {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	return RunfileIdentifier
}

func GetRegisteredPlugins() map[string]string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return RegisteredPlugins
}

func GetGameLogFromLogFile() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return GameLogFromLogFile
}

func GetBackupsStoreDir() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return BackupsStoreDir
}

func GetBackupLoopInterval() time.Duration {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return BackupLoopInterval
}

func GetBackupMode() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return BackupMode
}

func GetBackupMaxFileSize() int64 {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return BackupMaxFileSize
}

func GetBackupUseCompression() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return BackupUseCompression
}

func GetBackupKeepSnapshot() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return BackupKeepSnapshot
}

func GetBackupLoopActive() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return BackupLoopActive
}

func GetIsTelemetryEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsTelemetryEnabled
}

func GetIsGameServerRunning() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return IsGameServerRunning
}

func GetCreateGameServerLogFile() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return CreateGameServerLogFile
}
