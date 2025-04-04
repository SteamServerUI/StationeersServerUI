package core

import (
	"StationeersServerUI/src/backupsv2"
	"StationeersServerUI/src/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func HandleConfigJSON(w http.ResponseWriter, r *http.Request) {

	htmlFile, err := os.ReadFile("./UIMod/config.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading config.html: %v", err), http.StatusInternalServerError)
		return
	}

	htmlContent := string(htmlFile)

	// Determine selected attributes for boolean fields
	upnpTrueSelected := ""
	upnpFalseSelected := ""
	if config.UPNPEnabled {
		upnpTrueSelected = "selected"
	} else {
		upnpFalseSelected = "selected"
	}

	discordTrueSelected := ""
	discordFalseSelected := ""
	if config.IsDiscordEnabled {
		discordTrueSelected = "selected"
	} else {
		discordFalseSelected = "selected"
	}

	autoSaveTrueSelected := ""
	autoSaveFalseSelected := ""
	if config.AutoSave {
		autoSaveTrueSelected = "selected"
	} else {
		autoSaveFalseSelected = "selected"
	}

	autoPauseTrueSelected := ""
	autoPauseFalseSelected := ""
	if config.AutoPauseServer {
		autoPauseTrueSelected = "selected"
	} else {
		autoPauseFalseSelected = "selected"
	}

	startLocalTrueSelected := ""
	startLocalFalseSelected := ""
	if config.StartLocalHost {
		startLocalTrueSelected = "selected"
	} else {
		startLocalFalseSelected = "selected"
	}

	serverVisibleTrueSelected := ""
	serverVisibleFalseSelected := ""
	if config.ServerVisible {
		serverVisibleTrueSelected = "selected"
	} else {
		serverVisibleFalseSelected = "selected"
	}

	steamP2PTrueSelected := ""
	steamP2PFalseSelected := ""
	if config.UseSteamP2P {
		steamP2PTrueSelected = "selected"
	} else {
		steamP2PFalseSelected = "selected"
	}

	// Replace placeholders in the HTML with actual config values
	replacements := map[string]string{
		"{{discordToken}}":                  config.DiscordToken,
		"{{controlChannelID}}":              config.ControlChannelID,
		"{{statusChannelID}}":               config.StatusChannelID,
		"{{connectionListChannelID}}":       config.ConnectionListChannelID,
		"{{logChannelID}}":                  config.LogChannelID,
		"{{saveChannelID}}":                 config.SaveChannelID,
		"{{controlPanelChannelID}}":         config.ControlPanelChannelID,
		"{{blackListFilePath}}":             config.BlackListFilePath,
		"{{errorChannelID}}":                config.ErrorChannelID,
		"{{isDiscordEnabled}}":              fmt.Sprintf("%v", config.IsDiscordEnabled),
		"{{IsDiscordEnabledTrueSelected}}":  discordTrueSelected,
		"{{IsDiscordEnabledFalseSelected}}": discordFalseSelected,
		"{{gameBranch}}":                    config.GameBranch,
		"{{ServerName}}":                    config.ServerName,
		"{{SaveInfo}}":                      config.SaveInfo,
		"{{ServerMaxPlayers}}":              config.ServerMaxPlayers,
		"{{ServerPassword}}":                config.ServerPassword,
		"{{ServerAuthSecret}}":              config.ServerAuthSecret,
		"{{AdminPassword}}":                 config.AdminPassword,
		"{{GamePort}}":                      config.GamePort,
		"{{UpdatePort}}":                    config.UpdatePort,
		"{{UPNPEnabled}}":                   fmt.Sprintf("%v", config.UPNPEnabled), //unused, but maybe useful for future use
		"{{UPNPEnabledTrueSelected}}":       upnpTrueSelected,
		"{{UPNPEnabledFalseSelected}}":      upnpFalseSelected,
		"{{AutoSave}}":                      fmt.Sprintf("%v", config.AutoSave), //all of them
		"{{AutoSaveTrueSelected}}":          autoSaveTrueSelected,
		"{{AutoSaveFalseSelected}}":         autoSaveFalseSelected,
		"{{SaveInterval}}":                  config.SaveInterval,
		"{{AutoPauseServer}}":               fmt.Sprintf("%v", config.AutoPauseServer), //all of them
		"{{AutoPauseServerTrueSelected}}":   autoPauseTrueSelected,
		"{{AutoPauseServerFalseSelected}}":  autoPauseFalseSelected,
		"{{LocalIpAddress}}":                config.LocalIpAddress,
		"{{StartLocalHost}}":                fmt.Sprintf("%v", config.StartLocalHost), //all of them
		"{{StartLocalHostTrueSelected}}":    startLocalTrueSelected,
		"{{StartLocalHostFalseSelected}}":   startLocalFalseSelected,
		"{{ServerVisible}}":                 fmt.Sprintf("%v", config.ServerVisible), //all of them
		"{{ServerVisibleTrueSelected}}":     serverVisibleTrueSelected,
		"{{ServerVisibleFalseSelected}}":    serverVisibleFalseSelected,
		"{{UseSteamP2P}}":                   fmt.Sprintf("%v", config.UseSteamP2P), //all of them
		"{{UseSteamP2PTrueSelected}}":       steamP2PTrueSelected,
		"{{UseSteamP2PFalseSelected}}":      steamP2PFalseSelected,
		"{{ExePath}}":                       config.ExePath,
		"{{AdditionalParams}}":              config.AdditionalParams,
	}

	for placeholder, value := range replacements {
		htmlContent = strings.ReplaceAll(htmlContent, placeholder, value)
	}

	fmt.Fprint(w, htmlContent)
}

func SaveConfigForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusInternalServerError)
		return
	}

	// Load existing configuration
	existingConfig, err := config.LoadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading existing configuration: %v", err), http.StatusInternalServerError)
		return
	}

	// Use reflection to update fields present in the form
	v := reflect.ValueOf(existingConfig).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Tag.Get("json")
		formValues, exists := r.Form[fieldName] // Check if the field exists in the form data

		if !exists {
			continue // Skip fields not submitted in the form
		}

		// If the field exists, use the first value (even if it's empty)
		formValue := ""
		if len(formValues) > 0 {
			formValue = formValues[0]
		}

		field := v.Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(formValue) // Set the value, even if it's empty to allow clearing the field
		case reflect.Bool:
			field.SetBool(formValue == "true")
		}
	}

	// Save the updated config to file
	file, err := os.Create(config.ConfigPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating config.json: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(existingConfig); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding config.json: %v", err), http.StatusInternalServerError)
		return
	}

	// Reload the saved config into globals
	if _, err := config.LoadConfig(); err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}
	// Re-Initialize the backup manager with its global Interface
	if err := backupsv2.ReloadBackupManagerFromConfig(); err != nil {
		log.Printf("Failed to reload backup manager: %v", err)
		return
	}
	if config.IsDebugMode {
		fmt.Println("[BACKUP/DEBUG]Config and backup manager reloaded successfully")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SaveConfigRestful(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Parse the JSON data into a map
	var requestData map[string]interface{}
	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Load existing configuration
	existingConfig, err := config.LoadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading existing configuration: %v", err), http.StatusInternalServerError)
		return
	}

	// Use reflection to update fields present in the request data
	v := reflect.ValueOf(existingConfig).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Tag.Get("json")
		value, exists := requestData[fieldName] // Check if the field exists in the request data

		if !exists {
			continue // Skip fields not submitted in the request
		}

		field := v.Field(i)
		switch field.Kind() {
		case reflect.String:
			// Convert value to string if possible
			if strValue, ok := value.(string); ok {
				field.SetString(strValue)
			}
		case reflect.Bool:
			// Convert value to bool if possible
			if boolValue, ok := value.(bool); ok {
				field.SetBool(boolValue)
			}
		case reflect.Int:
			// Handle integers - this was missing in the original handler
			switch v := value.(type) {
			case float64: // JSON numbers are parsed as float64 by default
				field.SetInt(int64(v))
			case int:
				field.SetInt(int64(v))
			case int64:
				field.SetInt(v)
			case string:
				// Try to convert string to int if provided as string
				if intValue, err := strconv.ParseInt(v, 10, 64); err == nil {
					field.SetInt(intValue)
				}
			}
		}
	}

	// Save the updated config to file
	file, err := os.Create(config.ConfigPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating config.json: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(existingConfig); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding config.json: %v", err), http.StatusInternalServerError)
		return
	}

	// Reload the saved config into globals
	if _, err := config.LoadConfig(); err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}
	// Re-Initialize the backup manager with its global Interface
	if err := backupsv2.ReloadBackupManagerFromConfig(); err != nil {
		log.Printf("Failed to reload backup manager: %v", err)
		return
	}

	if config.IsDebugMode {
		fmt.Println("[BACKUP/DEBUG]Config and backup manager reloaded")
	}

	// Return success response in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Configuration updated successfully"})
}
