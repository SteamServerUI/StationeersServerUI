package web

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/gamemgr"
	"StationeersServerUI/src/ssestream"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Version string
	Branch  string
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./UIMod/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Version: config.Version,
		Branch:  config.Branch,
	}
	if data.Version == "" {
		data.Version = "unknown"
	}
	if data.Branch == "" {
		data.Branch = "unknown"
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ServeDetectionManager(w http.ResponseWriter, r *http.Request) {

	htmlFile, err := os.ReadFile("./UIMod/detectionmanager/detectionmanager.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading detectionmanager.html: %v", err), http.StatusInternalServerError)
		return
	}

	htmlContent := string(htmlFile)

	fmt.Fprint(w, htmlContent)
}

func ServeConfigPage(w http.ResponseWriter, r *http.Request) {

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

// StartServer HTTP handler
func StartServer(w http.ResponseWriter, r *http.Request) {
	if err := gamemgr.InternalStartServer(); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, "Server started.")
}

// StopServer HTTP handler
func StopServer(w http.ResponseWriter, r *http.Request) {
	if err := gamemgr.InternalStopServer(); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, "Server stopped.")
}

// handler for the /console endpoint
func GetLogOutput(w http.ResponseWriter, r *http.Request) {
	StartConsoleStream()(w, r)
}

// handler for the /console endpoint
func GetEventOutput(w http.ResponseWriter, r *http.Request) {
	StartDetectionEventStream()(w, r)
}

// StartConsoleStream creates an HTTP handler for console log SSE streaming
func StartConsoleStream() http.HandlerFunc {
	return ssestream.ConsoleStreamManager.CreateStreamHandler("Console")
}

// StartDetectionEventStream creates an HTTP handler for detection event SSE streaming
func StartDetectionEventStream() http.HandlerFunc {
	return ssestream.EventStreamManager.CreateStreamHandler("Event")
}
