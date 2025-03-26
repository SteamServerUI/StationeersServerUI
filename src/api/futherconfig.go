package api

import (
	"StationeersServerUI/src/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

// LoadConfigJSON loads the configuration from the JSON file
func loadConfigJSON() (*config.Config, error) {
	configPath := "./UIMod/config.json"
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config.json: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config.json: %v", err)
	}

	var config config.Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config.json: %v", err)
	}

	return &config, nil
}

func HandleConfigJSON(w http.ResponseWriter, r *http.Request) {
	config, err := loadConfigJSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading config.json: %v", err), http.StatusInternalServerError)
		return
	}

	htmlFile, err := os.ReadFile("./UIMod/furtherconfig.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading discord.html: %v", err), http.StatusInternalServerError)
		return
	}

	htmlContent := string(htmlFile)

	// Replace placeholders in the HTML with actual config values, including the new errorChannelID
	replacements := map[string]string{
		"{{discordToken}}":            config.DiscordToken,
		"{{controlChannelID}}":        config.ControlChannelID,
		"{{statusChannelID}}":         config.StatusChannelID,
		"{{connectionListChannelID}}": config.ConnectionListChannelID,
		"{{logChannelID}}":            config.LogChannelID,
		"{{saveChannelID}}":           config.SaveChannelID,
		"{{controlPanelChannelID}}":   config.ControlPanelChannelID,
		"{{blackListFilePath}}":       config.BlackListFilePath,
		"{{errorChannelID}}":          config.ErrorChannelID,
		"{{isDiscordEnabled}}":        fmt.Sprintf("%v", config.IsDiscordEnabled),
		"{{gameBranch}}":              config.GameBranch,
		"{{ServerName}}":              config.ServerName,
		"{{SaveFileName}}":            config.SaveFileName,
		"{{ServerMaxPlayers}}":        config.ServerMaxPlayers,
		"{{ServerPassword}}":          config.ServerPassword,
		"{{ServerAuthSecret}}":        config.ServerAuthSecret,
		"{{AdminPassword}}":           config.AdminPassword,
		"{{GamePort}}":                config.GamePort,
		"{{UpdatePort}}":              config.UpdatePort,
		"{{UPNPEnabled}}":             fmt.Sprintf("%v", config.UPNPEnabled),
		"{{AutoSave}}":                fmt.Sprintf("%v", config.AutoSave),
		"{{SaveInterval}}":            config.SaveInterval,
		"{{AutoPauseServer}}":         fmt.Sprintf("%v", config.AutoPauseServer),
		"{{LocalIpAddress}}":          config.LocalIpAddress,
		"{{StartLocalHost}}":          fmt.Sprintf("%v", config.StartLocalHost),
		"{{ServerVisible}}":           fmt.Sprintf("%v", config.ServerVisible),
		"{{UseSteamP2P}}":             fmt.Sprintf("%v", config.UseSteamP2P),
		"{{ExePath}}":                 config.ExePath,
		"{{AdditionalParams}}":        config.AdditionalParams,
	}

	fmt.Printf("Debug - UseSteamP2P value: %v, string representation: %s\n", config.UseSteamP2P, fmt.Sprintf("%v", config.UseSteamP2P))

	for placeholder, value := range replacements {
		htmlContent = strings.ReplaceAll(htmlContent, placeholder, value)
	}

	fmt.Fprint(w, htmlContent)
}

func SaveConfigJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Load existing configuration
	existingConfig, err := loadConfigJSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading existing configuration: %v", err), http.StatusInternalServerError)
		return
	}

	// Use reflection to update only submitted fields
	v := reflect.ValueOf(existingConfig).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Tag.Get("json")
		formValue := r.FormValue(fieldName)

		if formValue == "" {
			continue
		}

		field := v.Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(formValue)
		case reflect.Bool:
			field.SetBool(formValue == "true")
		}
	}

	configPath := "./UIMod/config.json"
	file, err := os.Create(configPath)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
