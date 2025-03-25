package api

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// LoadConfig loads the configuration from an XML file
func loadConfig() (*Config, error) {
	configPath := "./UIMod/config.xml"
	xmlFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = xml.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %v", err)
	}

	return &config, nil
}

func HandleConfig(w http.ResponseWriter, r *http.Request) {
	config, err := loadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading config: %v", err), http.StatusInternalServerError)
		return
	}

	htmlFile, err := os.ReadFile("./UIMod/config.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading config.html: %v", err), http.StatusInternalServerError)
		return
	}

	htmlContent := string(htmlFile)

	// Process settings with guaranteed order
	settings := strings.Split(config.Server.Settings, " ")
	settingsMap := make(map[string]string)
	var localIPValue string

	// First pass: extract LocalIpAddress and populate known params
	for i := 0; i < len(settings)-1; i += 2 {
		key := settings[i]
		value := settings[i+1]

		if key == "LocalIpAddress" {
			localIPValue = value
			continue
		}
		settingsMap[key] = value
	}

	// Prepare additional params (with LocalIpAddress last if present)
	additionalParamsStr := getAdditionalParams(settings)

	// Build replacements with consistent ordering
	replacements := map[string]string{
		"{{ExePath}}":          config.Server.ExePath,
		"{{StartLocalHost}}":   settingsMap["StartLocalHost"],
		"{{ServerVisible}}":    settingsMap["ServerVisible"],
		"{{GamePort}}":         settingsMap["GamePort"],
		"{{UpdatePort}}":       settingsMap["UpdatePort"],
		"{{AutoSave}}":         settingsMap["AutoSave"],
		"{{SaveInterval}}":     settingsMap["SaveInterval"],
		"{{ServerPassword}}":   settingsMap["ServerPassword"],
		"{{AdminPassword}}":    settingsMap["AdminPassword"],
		"{{ServerMaxPlayers}}": settingsMap["ServerMaxPlayers"],
		"{{ServerName}}":       settingsMap["ServerName"],
		"{{AdditionalParams}}": additionalParamsStr,
		"{{LocalIpAddress}}":   localIPValue, // Handled separately
		"{{SaveFileName}}":     config.SaveFileName,
	}

	// Apply replacements
	for placeholder, value := range replacements {
		if value != "" {
			htmlContent = strings.ReplaceAll(htmlContent, placeholder, value)
		}
	}

	fmt.Fprint(w, htmlContent)
}

func getAdditionalParams(settings []string) string {
	// List of known parameters (excluding LocalIpAddress)
	knownParams := map[string]bool{
		"StartLocalHost":   true,
		"ServerVisible":    true,
		"GamePort":         true,
		"UpdatePort":       true,
		"AutoSave":         true,
		"SaveInterval":     true,
		"ServerPassword":   true,
		"AdminPassword":    true,
		"ServerMaxPlayers": true,
		"ServerName":       true,
	}

	var regularAdditional []string // For unknown parameters
	var localIPParam string        // For LocalIpAddress only
	var otherKnownParams []string  // For known parameters that aren't LocalIpAddress

	// Process settings in pairs (key-value)
	for i := 0; i < len(settings)-1; i += 2 {
		key := settings[i]
		value := settings[i+1]

		switch {
		case key == "LocalIpAddress":
			localIPParam = key + " " + value
		case knownParams[key]:
			otherKnownParams = append(otherKnownParams, key+" "+value)
		default:
			regularAdditional = append(regularAdditional, key+" "+value)
		}
	}

	// Build the final parameter list in correct order:
	// 1. Other known parameters first
	// 2. Regular additional parameters
	// 3. LocalIpAddress (always last)
	var resultParams []string
	resultParams = append(resultParams, otherKnownParams...)
	resultParams = append(resultParams, regularAdditional...)
	if localIPParam != "" {
		resultParams = append(resultParams, localIPParam)
	}

	return strings.Join(resultParams, " ")
}

// SaveConfig saves the updated configuration to the XML file
func SaveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// We'll collect parameters in three separate groups
		var regularParams []string    // Main parameters in order
		var additionalParams []string // Additional parameters
		var localIPParam string       // LocalIpAddress (will be last)

		// Collect all regular parameters except LocalIpAddress
		if val := r.FormValue("StartLocalHost"); val != "" {
			regularParams = append(regularParams, "StartLocalHost", val)
		}
		if val := r.FormValue("ServerVisible"); val != "" {
			regularParams = append(regularParams, "ServerVisible", val)
		}
		if val := r.FormValue("GamePort"); val != "" {
			regularParams = append(regularParams, "GamePort", val)
		}
		if val := r.FormValue("UpdatePort"); val != "" {
			regularParams = append(regularParams, "UpdatePort", val)
		}
		if val := r.FormValue("AutoSave"); val != "" {
			regularParams = append(regularParams, "AutoSave", val)
		}
		if val := r.FormValue("SaveInterval"); val != "" {
			regularParams = append(regularParams, "SaveInterval", val)
		}
		if val := r.FormValue("ServerPassword"); val != "" {
			regularParams = append(regularParams, "ServerPassword", val)
		}
		if val := r.FormValue("AdminPassword"); val != "" {
			regularParams = append(regularParams, "AdminPassword", val)
		}
		if val := r.FormValue("ServerMaxPlayers"); val != "" {
			regularParams = append(regularParams, "ServerMaxPlayers", val)
		}
		if val := r.FormValue("ServerName"); val != "" {
			regularParams = append(regularParams, "ServerName", val)
		}

		// Collect AdditionalParams if they exist
		if extraParams := r.FormValue("AdditionalParams"); extraParams != "" {
			additionalParams = strings.Split(extraParams, " ")
		}

		// Collect LocalIpAddress separately
		if localIP := r.FormValue("LocalIpAddress"); localIP != "" {
			localIPParam = "LocalIpAddress " + localIP
		}

		// Combine all parameters in the correct order
		var finalSettings []string
		finalSettings = append(finalSettings, regularParams...)
		finalSettings = append(finalSettings, additionalParams...)
		if localIPParam != "" {
			finalSettings = append(finalSettings, localIPParam)
		}

		// Create the final settings string
		settingsStr := strings.Join(finalSettings, " ")

		// Determine the executable path based on OS
		var exePath string
		if runtime.GOOS == "windows" {
			exePath = "./rocketstation_DedicatedServer.exe"
		} else {
			exePath = "./rocketstation_DedicatedServer.x86_64"
		}

		// Create config structure
		config := Config{
			Server: struct {
				ExePath  string `xml:"exePath"`
				Settings string `xml:"settings"`
			}{
				ExePath:  exePath,
				Settings: settingsStr,
			},
			SaveFileName: r.FormValue("saveFileName"),
		}

		// Write to config file
		configPath := "./UIMod/config.xml"
		file, err := os.Create(configPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating config file: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		encoder := xml.NewEncoder(file)
		encoder.Indent("", "  ")
		if err := encoder.Encode(config); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding config: %v", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
