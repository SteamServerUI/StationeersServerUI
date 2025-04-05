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

func generateJwtKey() string {

	// ensure we return JwtKey if it's set
	if len(JwtKey) > 0 {
		return JwtKey
	}
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Failed to generate JWT key, using fallback")
		return "i-am-a-fallback-32-byte-secret-key!!"
	}
	return base64.RawURLEncoding.EncodeToString(key)
}
