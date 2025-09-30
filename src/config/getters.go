package config

import (
	"embed"
	"time"

	"github.com/bwmarrin/discordgo"
)

func GetDiscordToken() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return discordToken
}

func GetControlChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return controlChannelID
}

func GetStatusChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return statusChannelID
}

func GetConnectionListChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return connectionListChannelID
}

func GetLogChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return logChannelID
}

func GetSaveChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return saveChannelID
}

func GetControlPanelChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return controlPanelChannelID
}

func GetDiscordCharBufferSize() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return discordCharBufferSize
}

func GetExceptionMessageID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return exceptionMessageID
}

func GetDiscordSession() *discordgo.Session {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return discordSession
}

func GetBlackListFilePath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return blackListFilePath
}

func GetIsDiscordEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return isDiscordEnabled
}

func GetErrorChannelID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return errorChannelID
}

func GetBackupKeepLastN() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return backupKeepLastN
}

func GetIsCleanupEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return isCleanupEnabled
}

// GetBackupKeepDailyFor returns the retention period for daily backups in hours.
func GetBackupKeepDailyFor() time.Duration {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return backupKeepDailyFor
}

func GetBackupKeepWeeklyFor() time.Duration {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return backupKeepWeeklyFor
}

// GetBackupKeepMonthlyFor returns the retention period for monthly backups in hours.
func GetBackupKeepMonthlyFor() time.Duration {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return backupKeepMonthlyFor
}

// GetBackupCleanupInterval returns the cleanup interval in hours.
func GetBackupCleanupInterval() time.Duration {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return backupCleanupInterval
}

func GetIsNewTerrainAndSaveSystem() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return isNewTerrainAndSaveSystem
}

func GetGameBranch() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return gameBranch
}

func GetDifficulty() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return difficulty
}

func GetStartCondition() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return startCondition
}

func GetStartLocation() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return startLocation
}

func GetServerName() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return serverName
}

// special getter for backwards compatibility with SaveInfo
func GetLegacySaveInfo() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	saveinfo := saveName + ";" + worldID
	return saveinfo
}

func GetSaveName() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return saveName
}

func GetWorldID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return worldID
}

func GetServerMaxPlayers() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return serverMaxPlayers
}

func GetServerPassword() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return serverPassword
}

func GetServerAuthSecret() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return serverAuthSecret
}

func GetAdminPassword() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return adminPassword
}

func GetGamePort() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return gamePort
}

func GetUpdatePort() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return updatePort
}

func GetUPNPEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return uPNPEnabled
}

func GetAutoSave() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return autoSave
}

func GetSaveInterval() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return saveInterval
}

func GetAutoPauseServer() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return autoPauseServer
}

func GetLocalIpAddress() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return localIpAddress
}

func GetStartLocalHost() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return startLocalHost
}

func GetServerVisible() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return serverVisible
}

func GetUseSteamP2P() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return useSteamP2P
}

func GetExePath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return exePath
}

func GetAdditionalParams() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return additionalParams
}

func GetUsers() map[string]string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return users
}

func GetAuthEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return authEnabled
}

func GetJwtKey() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return jwtKey
}

func GetAuthTokenLifetime() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return authTokenLifetime
}

func GetIsDebugMode() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return isDebugMode
}

func GetCreateSSUILogFile() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return createSSUILogFile
}

func GetLogLevel() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return logLevel
}

func GetLogClutterToConsole() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return logClutterToConsole
}

func GetSubsystemFilters() []string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return subsystemFilters
}

func GetIsUpdateEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return isUpdateEnabled
}

func GetIsSSCMEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return isSSCMEnabled
}

func GetSSCMFilePath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return sscmFilePath
}

func GetSSCMPluginDir() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return sscmPluginDir
}

func GetSSCMWebDir() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return sscmWebDir
}

func GetAutoRestartServerTimer() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return autoRestartServerTimer
}

func GetAllowPrereleaseUpdates() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return allowPrereleaseUpdates
}

func GetAllowMajorUpdates() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return allowMajorUpdates
}

func GetIsConsoleEnabled() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return isConsoleEnabled
}

func GetLanguageSetting() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return languageSetting
}

func GetAutoStartServerOnStartup() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return autoStartServerOnStartup
}

func GetSSUIIdentifier() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return ssuiIdentifier
}

func GetSSUIWebPort() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return ssuiWebPort
}

// GetIsFirstTimeSetup returns the IsFirstTimeSetup
func GetIsFirstTimeSetup() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return isFirstTimeSetup
}

func GetConfigPath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return configPath
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
	return tlsCertPath
}

func GetTLSKeyPath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return tlsKeyPath
}

func GetUIModFolder() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return uiModFolder
}

func GetMaxSSEConnections() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return maxSSEConnections
}

func GetSSEMessageBufferSize() int {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return sseMessageBufferSize
}

func GetLogFolder() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return logFolder
}

func GetConfiguredBackupDir() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return configuredBackupDir
}

func GetConfiguredSafeBackupDir() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return configuredSafeBackupDir
}

func GetCustomDetectionsFilePath() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return customDetectionsFilePath
}

func GetGameServerAppID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return gameServerAppID
}

func GetCurrentBranchBuildID() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return currentBranchBuildID
}

func GetAllowAutoGameServerUpdates() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return allowAutoGameServerUpdates
}

func GetExtractedGameVersion() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return extractedGameVersion
}

func GetSkipSteamCMD() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return skipSteamCMD
}

func GetNoSanityCheck() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return noSanityCheck
}

func GetIsDockerContainer() bool {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return isDockerContainer
}

func GetOverrideAdvertisedIp() string {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return overrideAdvertisedIp
}

func GetV1UIFS() embed.FS {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return v1UIFS
}
