package settings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

// package settings handles API communication with the config values in package config via getter /setter functions.

// ConfigSetting represents metadata for a configuration setting
type ConfigSetting struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Group       string      `json:"group"`
	Description string      `json:"description"`
	Value       interface{} `json:"value"`
	Min         *int        `json:"min,omitempty"`
	Max         *int        `json:"max,omitempty"`
	Required    bool        `json:"required"`
}

// ConfigSettingsResponse represents the API response
type ConfigSettingsResponse struct {
	Data  []ConfigSetting `json:"data"`
	Error string          `json:"error,omitempty"`
}

func RetrieveSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	settings := []ConfigSetting{
		{
			Name:        "IsDebugMode",
			Type:        "bool",
			Group:       "System Settings",
			Description: "Enable pprof server",
			Value:       config.GetIsDebugMode(),
		},
		{
			Name:        "CreateSSUILogFile",
			Type:        "bool",
			Group:       "System Settings",
			Description: "Create SSUI log files",
			Value:       config.GetCreateSSUILogFile(),
		},
		{
			Name:        "LogLevel",
			Type:        "int",
			Group:       "System Settings",
			Description: "Logging verbosity level",
			Value:       config.GetLogLevel(),
			Min:         intPtr(0),
		},
		//{
		//	Name:        "BackendEndpointIP",
		//	Type:        "string",
		//	Group:       "System Settings",
		//	Description: "IP address for backend endpoint",
		//	Value:       config.GetBackendEndpointIP(),
		//	Required:    true,
		//},
		{
			Name:        "BackendEndpointPort",
			Type:        "string",
			Group:       "System Settings",
			Description: "Port for backend endpoint",
			Value:       config.GetSSUIWebPort(), // CHANGE VAR NAME TO MATCH (Get)BackendEndpointPort
			Required:    true,
		},
		//{
		//	Name:        "BackupKeepLastN",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Number of recent backups to keep",
		//	Value:       config.GetBackupKeepLastN(),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "IsCleanupEnabled",
		//	Type:        "bool",
		//	Group:       "Advanced Settings",
		//	Description: "Enable backup cleanup",
		//	Value:       config.GetIsCleanupEnabled(),
		//},
		//{
		//	Name:        "BackupKeepDailyFor",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Hours to keep daily backups",
		//	Value:       int(config.GetBackupKeepDailyFor() / time.Hour),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "BackupKeepWeeklyFor",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Hours to keep weekly backups",
		//	Value:       int(config.GetBackupKeepWeeklyFor() / time.Hour),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "BackupKeepMonthlyFor",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Hours to keep monthly backups",
		//	Value:       int(config.GetBackupKeepMonthlyFor() / time.Hour),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "BackupCleanupInterval",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Hours between backup cleanup runs",
		//	Value:       int(config.GetBackupCleanupInterval() / time.Hour),
		//	Min:         intPtr(0),
		//},
		//{
		//	Name:        "BackupWaitTime",
		//	Type:        "int",
		//	Group:       "Advanced Settings",
		//	Description: "Seconds to wait before backup",
		//	Value:       int(config.GetBackupWaitTime() / time.Second),
		//	Min:         intPtr(0),
		//},
		{
			Name:        "GameBranch",
			Type:        "string",
			Group:       "Update Settings",
			Description: "Game branch for updates (Restart Required)",
			Value:       config.GetGameBranch(),
		},
		//{
		//	Name:        "Users",
		//	Type:        "map",
		//	Group:       "Advanced Settings",
		//	Description: "User authentication mappings",
		//	Value:       config.GetUsers(),
		//},
		{
			Name:        "AuthEnabled",
			Type:        "bool",
			Group:       "System Settings",
			Description: "Enable authentication",
			Value:       config.GetAuthEnabled(),
		},
		{
			Name:        "JwtKey",
			Type:        "string",
			Group:       "Security Settings",
			Description: "Encryption key for Authentication",
			Value:       config.GetJwtKey(),
		},
		{
			Name:        "AuthTokenLifetime",
			Type:        "int",
			Group:       "Security Settings",
			Description: "Token lifetime in seconds",
			Value:       config.GetAuthTokenLifetime(),
			Min:         intPtr(0),
		},
		{
			Name:        "IsUpdateEnabled",
			Type:        "bool",
			Group:       "System Settings",
			Description: "Enable automatic updates",
			Value:       config.GetIsUpdateEnabled(),
		},
		//{
		//	Name:        "IsSSCMEnabled",
		//	Type:        "bool",
		//	Group:       "Advanced Settings",
		//	Description: "Enable SSCM integration",
		//	Value:       config.GetIsSSCMEnabled(),
		//},
		//{
		//	Name:        "IsCodeServerEnabled",
		//	Type:        "bool",
		//	Group:       "System Settings",
		//	Description: "Enables the Code Server integration. Linux only. Attention: Does NOT have an updater at the moment. Be careful with this feature.",
		//	Value:       config.GetIsCodeServerEnabled(),
		//},
		{
			Name:        "AllowPrereleaseUpdates",
			Type:        "bool",
			Group:       "Update Settings",
			Description: "Allow prerelease updates",
			Value:       config.GetAllowPrereleaseUpdates(),
		},
		{
			Name:        "AllowMajorUpdates",
			Type:        "bool",
			Group:       "Update Settings",
			Description: "Allow major version updates",
			Value:       config.GetAllowMajorUpdates(),
		},
		// Discord Settings
		{
			Name:        "IsDiscordEnabled",
			Type:        "bool",
			Group:       "Discord Settings",
			Description: "Enable Discord integration",
			Value:       config.GetIsDiscordEnabled(),
		},
		{
			Name:        "DiscordToken",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Discord bot token",
			Value:       config.GetDiscordToken(),
		},
		{
			Name:        "ControlChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Control channel ID",
			Value:       config.GetControlChannelID(),
		},
		{
			Name:        "StatusChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Status channel ID",
			Value:       config.GetStatusChannelID(),
		},
		{
			Name:        "ConnectionListChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Connection list channel ID",
			Value:       config.GetConnectionListChannelID(),
		},
		{
			Name:        "LogChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Log channel ID",
			Value:       config.GetLogChannelID(),
		},
		{
			Name:        "SaveChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Save channel ID",
			Value:       config.GetSaveChannelID(),
		},
		{
			Name:        "ControlPanelChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Control panel channel ID",
			Value:       config.GetControlPanelChannelID(),
		},
		{
			Name:        "DiscordCharBufferSize",
			Type:        "int",
			Group:       "Discord Settings",
			Description: "Discord character buffer size",
			Value:       config.GetDiscordCharBufferSize(),
		},
		{
			Name:        "ErrorChannelID",
			Type:        "string",
			Group:       "Discord Settings",
			Description: "Error channel ID",
			Value:       config.GetErrorChannelID(),
		},
		//{
		//	Name:        "BackupContentDir",
		//	Type:        "string",
		//	Group:       "Backup Settings",
		//	Description: "Backup content directory",
		//	Value:       config.GetBackupContentDir(),
		//},
		//{
		//	Name:        "BackupsStoreDir",
		//	Type:        "string",
		//	Group:       "Backup Settings",
		//	Description: "Backup stored backups directory",
		//	Value:       config.GetBackupsStoreDir(),
		//},
		//{
		//	Name:        "BackupLoopInterval",
		//	Type:        "string",
		//	Group:       "Backup Settings",
		//	Description: "Backup loop interval",
		//	Value:       config.GetBackupLoopInterval().String(),
		//},
		//{
		//	Name:        "BackupMode",
		//	Type:        "string",
		//	Group:       "Backup Settings",
		//	Description: "Backup mode",
		//	Value:       config.GetBackupMode(),
		//},
		//{
		//	Name:        "BackupMaxFileSize",
		//	Type:        "int",
		//	Group:       "Backup Settings",
		//	Description: "Max file size",
		//	Value:       config.GetBackupMaxFileSize(),
		//},
		//{
		//	Name:        "BackupUseCompression",
		//	Type:        "bool",
		//	Group:       "Backup Settings",
		//	Description: "Use compression",
		//	Value:       config.GetBackupUseCompression(),
		//},
		//{
		//	Name:        "BackupKeepSnapshot",
		//	Type:        "bool",
		//	Group:       "Backup Settings",
		//	Description: "Keep snapshot",
		//	Value:       config.GetBackupKeepSnapshot(),
		//},
		{
			Name:        "IsConsoleEnabled",
			Type:        "bool",
			Group:       "System Settings",
			Description: "Expose various actions directly in the command line (Restart Required)",
			Value:       config.GetIsConsoleEnabled(),
		},
	}

	response := ConfigSettingsResponse{
		Data: settings,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

// intPtr creates a pointer to an int
func intPtr(i int) *int {
	return &i
}
