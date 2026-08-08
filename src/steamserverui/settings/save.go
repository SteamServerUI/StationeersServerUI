package settings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
)

// package settings handles API communication with the config values in package config via getter /setter functions.

// setterFunc defines the signature for setter functions
type setterFunc func(interface{}) error

// setterMap maps JSON keys (global variable names) to setter functions with type checking
var setterMap = map[string]setterFunc{
	"BackendEndpointPort": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetBackendEndpointPort(str)
		}
		return fmt.Errorf("invalid type for BackendEndpointPort: expected string")
	},
	"DiscordToken": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetDiscordToken(str)
		}
		return fmt.Errorf("invalid type for DiscordToken: expected string")
	},
	"ControlChannelID": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetControlChannelID(str)
		}
		return fmt.Errorf("invalid type for ControlChannelID: expected string")
	},
	"StatusChannelID": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetStatusChannelID(str)
		}
		return fmt.Errorf("invalid type for StatusChannelID: expected string")
	},
	"ConnectionListChannelID": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetConnectionListChannelID(str)
		}
		return fmt.Errorf("invalid type for ConnectionListChannelID: expected string")
	},
	"LogChannelID": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetLogChannelID(str)
		}
		return fmt.Errorf("invalid type for LogChannelID: expected string")
	},
	"SaveChannelID": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetSaveChannelID(str)
		}
		return fmt.Errorf("invalid type for SaveChannelID: expected string")
	},
	"ControlPanelChannelID": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetControlPanelChannelID(str)
		}
		return fmt.Errorf("invalid type for ControlPanelChannelID: expected string")
	},
	"DiscordCharBufferSize": func(v interface{}) error {
		if f, ok := v.(float64); ok {
			return config.SetDiscordCharBufferSize(int(f))
		}
		return fmt.Errorf("invalid type for DiscordCharBufferSize: expected number")
	},
	"IsDiscordEnabled": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetIsDiscordEnabled(b)
		}
		return fmt.Errorf("invalid type for IsDiscordEnabled: expected bool")
	},
	"ErrorChannelID": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetErrorChannelID(str)
		}
		return fmt.Errorf("invalid type for ErrorChannelID: expected string")
	},
	"GameBranch": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetGameBranch(str)
		}
		return fmt.Errorf("invalid type for GameBranch: expected string")
	},
	"Users": func(v interface{}) error {
		if m, ok := v.(map[string]interface{}); ok {
			users := make(map[string]string)
			for k, val := range m {
				if strVal, ok := val.(string); ok {
					users[k] = strVal
				} else {
					return fmt.Errorf("invalid value type for Users: expected string")
				}
			}
			return config.SetUsers(users)
		}
		return fmt.Errorf("invalid type for Users: expected map[string]string")
	},
	"AuthEnabled": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetAuthEnabled(b)
		}
		return fmt.Errorf("invalid type for AuthEnabled: expected bool")
	},
	//"JwtKey": func(v interface{}) error {
	//	if str, ok := v.(string); ok {
	//		return config.SetJwtKey(str)
	//	}
	//	return fmt.Errorf("invalid type for JwtKey: expected string")
	//},
	"AuthTokenLifetime": func(v interface{}) error {
		if f, ok := v.(float64); ok {
			return config.SetAuthTokenLifetime(int(f))
		}
		return fmt.Errorf("invalid type for AuthTokenLifetime: expected number")
	},
	"IsDebugMode": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetIsDebugMode(b)
		}
		return fmt.Errorf("invalid type for IsDebugMode: expected bool")
	},
	"CreateSSUILogFile": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetCreateSSUILogFile(b)
		}
		return fmt.Errorf("invalid type for CreateSSUILogFile: expected bool")
	},
	"LogLevel": func(v interface{}) error {
		if f, ok := v.(float64); ok {
			return config.SetLogLevel(int(f))
		}
		return fmt.Errorf("invalid type for LogLevel: expected number")
	},
	"GameLogFromLogFile": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetGameLogFromLogFile(b)
		}
		return fmt.Errorf("invalid type for GameLogFromLogFile: expected bool")
	},
	"SubsystemFilters": func(v interface{}) error {
		if arr, ok := v.([]interface{}); ok {
			filters := make([]string, 0, len(arr))
			for _, val := range arr {
				if strVal, ok := val.(string); ok {
					filters = append(filters, strVal)
				} else {
					return fmt.Errorf("invalid value type for SubsystemFilters: expected string")
				}
			}
			return config.SetSubsystemFilters(filters)
		}
		return fmt.Errorf("invalid type for SubsystemFilters: expected array of strings")
	},
	"IsUpdateEnabled": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetIsUpdateEnabled(b)
		}
		return fmt.Errorf("invalid type for IsUpdateEnabled: expected bool")
	},
	"IsSSCMEnabled": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetIsSSCMEnabled(b)
		}
		return fmt.Errorf("invalid type for IsSSCMEnabled: expected bool")
	},
	"IsBepInExEnabled": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetIsBepInExEnabled(b)
		}
		return fmt.Errorf("invalid type for IsBepInExEnabled: expected bool")
	},
	"AllowPrereleaseUpdates": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetAllowPrereleaseUpdates(b)
		}
		return fmt.Errorf("invalid type for AllowPrereleaseUpdates: expected bool")
	},
	"AllowMajorUpdates": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetAllowMajorUpdates(b)
		}
		return fmt.Errorf("invalid type for AllowMajorUpdates: expected bool")
	},
	"RunfileIdentifier": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetRunfileIdentifier(str)
		}
		return fmt.Errorf("invalid type for RunfileIdentifier: expected string")
	},
	"BackupsStoreDir": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetBackupsStoreDir(str)
		}
		return fmt.Errorf("invalid type for BackupsStoreDir: expected string")
	},
	"BackupLoopActive": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetBackupLoopActive(b)
		}
		return fmt.Errorf("invalid type for BackupLoopActive: expected bool")
	},
	"BackupLoopInterval": func(v interface{}) error {
		if str, ok := v.(string); ok {
			interval, err := time.ParseDuration(str)
			if err != nil {
				return err
			}
			return config.SetBackupLoopInterval(interval)
		}
		return fmt.Errorf("invalid type for BackupLoopInterval: expected string")
	},
	"BackupMode": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetBackupMode(str)
		}
		return fmt.Errorf("invalid type for BackupMode: expected string")
	},
	"BackupMaxFileSize": func(v interface{}) error {
		if f, ok := v.(float64); ok {
			return config.SetBackupMaxFileSize(int64(f))
		}
		return fmt.Errorf("invalid type for BackupMaxFileSize: expected number")
	},
	"BackupUseCompression": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetBackupUseCompression(b)
		}
		return fmt.Errorf("invalid type for BackupUseCompression: expected bool")
	},
	"BackupKeepSnapshot": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetBackupKeepSnapshot(b)
		}
		return fmt.Errorf("invalid type for BackupKeepSnapshot: expected bool")
	},
	//"IsTelemetryEnabled": func(v interface{}) error {
	//	//Accepts both bool and an empty string as a valid value for true.
	//	var expectedtype string
	//	if _, ok := v.(string); ok {
	//		expectedtype = "string"
	//		return config.SetIsTelemetryEnabled(true)
	//	}
	//	if _, ok := v.(bool); ok {
	//		expectedtype = "bool"
	//		return config.SetIsTelemetryEnabled(v.(bool))
	//	}
	//	return fmt.Errorf("invalid type for IsTelemetryEnabled: expected %s", expectedtype)
	//},
	"IsSSUICLIConsoleEnabled": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetIsSSUICLIConsoleEnabled(b)
		}
		return fmt.Errorf("invalid type for IsSSUICLIConsoleEnabled: expected bool")
	},
	"AllowAutoGameServerUpdates": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetAllowAutoGameServerUpdates(b)
		}
		return fmt.Errorf("invalid type for AllowAutoGameServerUpdates: expected bool")
	},
	"AutoStartServerOnStartup": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetAutoStartServerOnStartup(b)
		}
		return fmt.Errorf("invalid type for AutoStartServerOnStartup: expected bool")
	},
	"AutoRestartServerTimer": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return config.SetAutoRestartServerTimer(str)
		}
		return fmt.Errorf("invalid type for AutoRestartServerTimer: expected string")
	},
	"LogClutterToConsole": func(v interface{}) error {
		if b, ok := v.(bool); ok {
			return config.SetLogClutterToConsole(b)
		}
		return fmt.Errorf("invalid type for LogClutterToConsole: expected bool")
	},
	"LanguageSetting": func(v interface{}) error {
		if str, ok := v.(string); ok {
			return setLanguageSetting(str)
		}
		return fmt.Errorf("invalid type for LanguageSetting: expected string")
	},
}

type updateSettingRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func HandleSaveSetting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeSettingError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only PATCH requests are allowed")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	decoder := json.NewDecoder(r.Body)
	var request updateSettingRequest
	if err := decoder.Decode(&request); err != nil {
		writeSettingError(w, http.StatusBadRequest, "invalid_json", "A setting key and value are required")
		return
	}
	request.Key = strings.TrimSpace(request.Key)
	setter, exists := setterMap[request.Key]
	if !exists {
		writeSettingError(w, http.StatusBadRequest, "unknown_setting", "The requested setting does not exist")
		return
	}
	if err := setter(request.Value); err != nil {
		writeSettingError(w, http.StatusBadRequest, "invalid_setting_value", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{
		"key":     request.Key,
		"value":   request.Value,
		"effects": settingEffects(request.Key),
	}})
}

func settingEffects(key string) []string {
	if strings.Contains(key, "Port") || strings.Contains(key, "CLI") ||
		strings.Contains(key, "BackupsStore") || strings.Contains(key, "BackupLoop") ||
		strings.Contains(key, "GameLog") {
		return []string{"backend.reload"}
	}
	return []string{"immediate"}
}

func writeSettingError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func setLanguageSetting(value string) error {
	defer loader.ReloadBackend()
	return config.SetLanguageSetting(value)
}
