package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"strconv"
)

// getConfigValue retrieves a string value from JSON, env var, or default
func getConfigValue(jsonValue, envKey, defaultValue string) string {
	if jsonValue != "" {
		return jsonValue
	}
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	return defaultValue
}

// getConfigValueInt retrieves an int value from JSON, env var, or default
func getConfigValueInt(jsonValue int, envKey string, defaultValue int) int {
	if jsonValue != 0 { // 0 means unset in JSON due to omitempty
		return jsonValue
	}
	if envValue := os.Getenv(envKey); envValue != "" {
		if val, err := strconv.Atoi(envValue); err == nil {
			return val
		}
	}
	return defaultValue
}

// generateJwtKey creates a random 32-byte key for JWT if not provided
func generateJwtKey() string {
	if JwtKey != "" {
		return JwtKey
	}
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		// Fallback in case of error (rare), but this should log in a real app
		return "i-am-a-fallback-32-byte-secret-key!!"
	}
	return base64.RawURLEncoding.EncodeToString(key)
}

func GetSecretsFromEnv(jsonconfig JsonConfig) {
	// Try to locate User config from JSON, if not found, use env vars, and fallback to defaults, but generate JwtKey (no default)
	Username = getConfigValue(jsonconfig.Username, "SSUI_USERNAME", "admin")
	Password = getConfigValue(jsonconfig.Password, "SSUI_PASSWORD", "password")
	JwtKey = getConfigValue(jsonconfig.JwtKey, "SSUI_JWT_KEY", generateJwtKey())
	AuthTokenLifetime = getConfigValueInt(jsonconfig.AuthTokenLifetime, "AUTH_TOKEN_LIFETIME", 1440)
}
