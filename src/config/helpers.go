//helpers.go

package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func getString(jsonVal, envKey, defaultVal string) string {
	if jsonVal != "" {
		return jsonVal
	}
	if envVal := os.Getenv(envKey); envVal != "" {
		return envVal
	}
	return defaultVal
}

func getInt(jsonVal int, envKey string, defaultVal int) int {
	if jsonVal != 0 {
		return jsonVal
	}
	if envVal := os.Getenv(envKey); envVal != "" {
		if val, err := strconv.Atoi(envVal); err == nil {
			return val
		}
	}
	return defaultVal
}

func getBool(jsonVal bool, envKey string, defaultVal bool) bool {
	if jsonVal {
		return true
	}
	if envVal := os.Getenv(envKey); envVal != "" {
		if val, err := strconv.ParseBool(envVal); err == nil {
			return val
		}
	}
	return defaultVal
}

func getDefaultExePath() string {
	if runtime.GOOS == "windows" {
		return "./rocketstation_DedicatedServer.exe"
	}
	return "./rocketstation_DedicatedServer.x86_64"
}

func setDefaults(cfg *JsonConfig) {
	if cfg.ExePath == "" {
		cfg.ExePath = getDefaultExePath()
	}
	if cfg.DiscordCharBufferSize <= 0 {
		cfg.DiscordCharBufferSize = 1000
	}
	if cfg.GameBranch == "" {
		cfg.GameBranch = "public"
	}
	if cfg.SaveInfo == "" {
		cfg.SaveInfo = "Moon Moon"
	}
	if cfg.BackupKeepDailyFor <= 0 {
		cfg.BackupKeepDailyFor = 24
	}
	if cfg.BackupKeepWeeklyFor <= 0 {
		cfg.BackupKeepWeeklyFor = 168
	}
	if cfg.BackupKeepMonthlyFor <= 0 {
		cfg.BackupKeepMonthlyFor = 730
	}
	if cfg.BackupCleanupInterval <= 0 {
		cfg.BackupCleanupInterval = 730
	}
	if cfg.BackupWaitTime <= 0 {
		cfg.BackupWaitTime = 30
	}
}

func logConfigDetails() {

	logger.Log(LogDebug, "Gameserver config values loaded")
	logger.Log(LogDebug, "---- GENERAL CONFIG VARS ----")
	logger.Log(LogDebug, fmt.Sprintf("Branch: %s", Branch))
	logger.Log(LogDebug, fmt.Sprintf("GameBranch: %s", GameBranch))
	logger.Log(LogDebug, "IsDiscordEnabled: "+strconv.FormatBool(IsDiscordEnabled))
	logger.Log(LogDebug, "IsCleanupEnabled: "+strconv.FormatBool(IsCleanupEnabled))
	logger.Log(LogDebug, "IsDebugMode: "+strconv.FormatBool(IsDebugMode))
	logger.Log(LogDebug, "IsFirstTimeSetup: "+strconv.FormatBool(IsFirstTimeSetup))

	logger.Log(LogDebug, "---- DISCORD CONFIG VARS ----")
	logger.Log(LogDebug, fmt.Sprintf("BlackListFilePath: %s", BlackListFilePath))
	logger.Log(LogDebug, fmt.Sprintf("ConnectionListChannelID: %s", ConnectionListChannelID))
	logger.Log(LogDebug, fmt.Sprintf("ControlChannelID: %s", ControlChannelID))
	logger.Log(LogDebug, fmt.Sprintf("ControlPanelChannelID: %s", ControlPanelChannelID))
	logger.Log(LogDebug, fmt.Sprintf("DiscordCharBufferSize: %d", DiscordCharBufferSize))
	logger.Log(LogDebug, fmt.Sprintf("DiscordToken: %s", DiscordToken))
	logger.Log(LogDebug, fmt.Sprintf("ErrorChannelID: %s", ErrorChannelID))
	logger.Log(LogDebug, fmt.Sprintf("IsDiscordEnabled: %v", IsDiscordEnabled))
	logger.Log(LogDebug, fmt.Sprintf("LogChannelID: %s", LogChannelID))
	logger.Log(LogDebug, fmt.Sprintf("LogMessageBuffer: %s", LogMessageBuffer))
	logger.Log(LogDebug, fmt.Sprintf("SaveChannelID: %s", SaveChannelID))
	logger.Log(LogDebug, fmt.Sprintf("StatusChannelID: %s", StatusChannelID))

	logger.Log(LogDebug, "---- BACKUP CONFIG VARS ----")
	logger.Log(LogDebug, fmt.Sprintf("BackupKeepLastN: %d", BackupKeepLastN))
	logger.Log(LogDebug, fmt.Sprintf("BackupKeepDailyFor: %s", BackupKeepDailyFor))
	logger.Log(LogDebug, fmt.Sprintf("BackupKeepWeeklyFor: %s", BackupKeepWeeklyFor))
	logger.Log(LogDebug, fmt.Sprintf("BackupKeepMonthlyFor: %s", BackupKeepMonthlyFor))
	logger.Log(LogDebug, fmt.Sprintf("BackupCleanupInterval: %s", BackupCleanupInterval))
	logger.Log(LogDebug, fmt.Sprintf("ConfiguredBackupDir: %s", ConfiguredBackupDir))
	logger.Log(LogDebug, fmt.Sprintf("ConfiguredSafeBackupDir: %s", ConfiguredSafeBackupDir))
	logger.Log(LogDebug, fmt.Sprintf("BackupWaitTime: %s", BackupWaitTime))

	logger.Log(LogDebug, "---- AUTHENTICATION CONFIG VARS ----")
	logger.Log(LogDebug, fmt.Sprintf("AuthTokenLifetime: %d", AuthTokenLifetime))
	logger.Log(LogDebug, fmt.Sprintf("JwtKey: %s", JwtKey))
	logger.Log(LogDebug, fmt.Sprintf("Password: %s", Password))
	logger.Log(LogDebug, fmt.Sprintf("Username: %s", Username))

	logger.Log(LogDebug, "---- MISC CONFIG VARS ----")
	logger.Log(LogDebug, fmt.Sprintf("Branch: %s", Branch))
	logger.Log(LogDebug, fmt.Sprintf("GameServerAppID: %s", GameServerAppID))
	logger.Log(LogDebug, fmt.Sprintf("Version: %s", Version))

}

func generateJwtKey() string {

	// ensure we return JwtKey if it's set
	if len(JwtKey) > 0 {
		return JwtKey
	}
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		logger.Log(LogError, "Failed to generate JWT key, using fallback")
		return "i-am-a-fallback-32-byte-secret-key!!"
	}
	return base64.RawURLEncoding.EncodeToString(key)
}
