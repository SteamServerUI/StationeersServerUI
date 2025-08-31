package loader

import (
	"fmt"
	"strconv"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

func PrintConfigDetails() {
	logger.Config.Debug("Gameserver config values loaded")
	logger.Config.Debug("---- GENERAL CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("Branch: %s", config.Branch))
	logger.Config.Debug(fmt.Sprintf("GameBranch: %s", config.GameBranch))
	logger.Config.Debug("IsDiscordEnabled: " + strconv.FormatBool(config.IsDiscordEnabled))
	logger.Config.Debug("IsCleanupEnabled: " + strconv.FormatBool(config.IsCleanupEnabled))
	logger.Config.Debug("IsDebugMode (pprof Server): " + strconv.FormatBool(config.IsDebugMode))
	logger.Config.Debug("IsFirstTimeSetup: " + strconv.FormatBool(config.IsFirstTimeSetup))

	logger.Config.Debug("---- DISCORD CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("BlackListFilePath: %s", config.BlackListFilePath))
	logger.Config.Debug(fmt.Sprintf("ConnectionListChannelID: %s", config.ConnectionListChannelID))
	logger.Config.Debug(fmt.Sprintf("ControlChannelID: %s", config.ControlChannelID))
	logger.Config.Debug(fmt.Sprintf("ControlPanelChannelID: %s", config.ControlPanelChannelID))
	logger.Config.Debug(fmt.Sprintf("DiscordCharBufferSize: %d", config.DiscordCharBufferSize))
	logger.Config.Debug(fmt.Sprintf("DiscordToken: %s", config.DiscordToken))
	logger.Config.Debug(fmt.Sprintf("ErrorChannelID: %s", config.ErrorChannelID))
	logger.Config.Debug(fmt.Sprintf("IsDiscordEnabled: %v", config.IsDiscordEnabled))
	logger.Config.Debug(fmt.Sprintf("LogChannelID: %s", config.LogChannelID))
	logger.Config.Debug(fmt.Sprintf("LogMessageBuffer: %s", config.LogMessageBuffer))
	logger.Config.Debug(fmt.Sprintf("SaveChannelID: %s", config.SaveChannelID))
	logger.Config.Debug(fmt.Sprintf("StatusChannelID: %s", config.StatusChannelID))

	logger.Config.Debug("---- BACKUP CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("BackupKeepLastN: %d", config.BackupKeepLastN))
	logger.Config.Debug(fmt.Sprintf("BackupKeepDailyFor: %s", config.BackupKeepDailyFor))
	logger.Config.Debug(fmt.Sprintf("BackupKeepWeeklyFor: %s", config.BackupKeepWeeklyFor))
	logger.Config.Debug(fmt.Sprintf("BackupKeepMonthlyFor: %s", config.BackupKeepMonthlyFor))
	logger.Config.Debug(fmt.Sprintf("BackupCleanupInterval: %s", config.BackupCleanupInterval))
	logger.Config.Debug(fmt.Sprintf("ConfiguredBackupDir: %s", config.ConfiguredBackupDir))
	logger.Config.Debug(fmt.Sprintf("ConfiguredSafeBackupDir: %s", config.ConfiguredSafeBackupDir))
	logger.Config.Debug(fmt.Sprintf("BackupWaitTime: %s", config.BackupWaitTime))

	logger.Config.Debug("---- AUTHENTICATION CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("AuthTokenLifetime: %d", config.AuthTokenLifetime))
	logger.Config.Debug(fmt.Sprintf("JwtKey: %s", config.JwtKey))

	logger.Config.Debug("---- MISC CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("Branch: %s", config.Branch))
	logger.Config.Debug(fmt.Sprintf("GameServerAppID: %s", config.GameServerAppID))
	logger.Config.Debug(fmt.Sprintf("Version: %s", config.Version))
	logger.Config.Debug(fmt.Sprintf("IsNewTerrainAndSaveSystem: %v", config.IsNewTerrainAndSaveSystem))

	logger.Config.Debug("----  UPDATER CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("AllowPrereleaseUpdates: %v", config.AllowPrereleaseUpdates))
	logger.Config.Debug(fmt.Sprintf("AllowMajorUpdates: %v", config.AllowMajorUpdates))
	logger.Config.Debug(fmt.Sprintf("IsUpdateEnabled: %v", config.IsUpdateEnabled))

	logger.Config.Debug("----  SSCM CONFIG VARS ----")
	logger.Config.Debug(fmt.Sprintf("SSCMFilePath: %s", config.SSCMFilePath))
	logger.Config.Debug(fmt.Sprintf("IsSSCMEnabled: %v", config.IsSSCMEnabled))
}
