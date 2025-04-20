package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/commandmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/gamemgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/ssestream"
)

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Version string
	Branch  string
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(config.IndexHtmlPath)
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

	htmlFile, err := os.ReadFile(config.DetectionManagerHtmlPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading detectionmanager.html: %v", err), http.StatusInternalServerError)
		return
	}

	htmlContent := string(htmlFile)

	fmt.Fprint(w, htmlContent)
}

func ServeConfigPage(w http.ResponseWriter, r *http.Request) {

	htmlFile, err := os.ReadFile(config.ConfigHtmlPath)
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
		"{{UPNPEnabled}}":                   fmt.Sprintf("%v", config.UPNPEnabled),
		"{{UPNPEnabledTrueSelected}}":       upnpTrueSelected,
		"{{UPNPEnabledFalseSelected}}":      upnpFalseSelected,
		"{{AutoSave}}":                      fmt.Sprintf("%v", config.AutoSave),
		"{{AutoSaveTrueSelected}}":          autoSaveTrueSelected,
		"{{AutoSaveFalseSelected}}":         autoSaveFalseSelected,
		"{{SaveInterval}}":                  config.SaveInterval,
		"{{AutoPauseServer}}":               fmt.Sprintf("%v", config.AutoPauseServer),
		"{{AutoPauseServerTrueSelected}}":   autoPauseTrueSelected,
		"{{AutoPauseServerFalseSelected}}":  autoPauseFalseSelected,
		"{{LocalIpAddress}}":                config.LocalIpAddress,
		"{{StartLocalHost}}":                fmt.Sprintf("%v", config.StartLocalHost),
		"{{StartLocalHostTrueSelected}}":    startLocalTrueSelected,
		"{{StartLocalHostFalseSelected}}":   startLocalFalseSelected,
		"{{ServerVisible}}":                 fmt.Sprintf("%v", config.ServerVisible),
		"{{ServerVisibleTrueSelected}}":     serverVisibleTrueSelected,
		"{{ServerVisibleFalseSelected}}":    serverVisibleFalseSelected,
		"{{UseSteamP2P}}":                   fmt.Sprintf("%v", config.UseSteamP2P),
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
	logger.Web.Debug("Received start request from API")
	if err := gamemgr.InternalStartServer(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Web.Core("Error starting server: " + err.Error())
		return
	}
	fmt.Fprint(w, "Server started.")
	logger.Web.Core("Server started.")
}

// StopServer HTTP handler
func StopServer(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("Received stop request from API")
	if err := gamemgr.InternalStopServer(); err != nil {
		if err.Error() == "server not running" {
			fmt.Fprint(w, "Server was not running or was already stopped")
			logger.Web.Core("Server not running or was already stopped")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Web.Core("Error stopping server: " + err.Error())
		return
	}
	fmt.Fprint(w, "Server stopped.")
	logger.Web.Core("Server stopped.")
}

func GetGameServerRunState(w http.ResponseWriter, r *http.Request) {
	runState := gamemgr.InternalIsServerRunning()
	response := map[string]interface{}{
		"isRunning": runState,
		"uuid":      config.GameServerUUID.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to respond with Game Server status", http.StatusInternalServerError)
		return
	}
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

func ServeTwoBoxCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, config.UIModFolder+"twoboxform/twoboxform.css")
}

func ServeTwoBoxJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, config.UIModFolder+"twoboxform/twoboxform.js")
}

func ServeSSCMJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, config.SSCMWebDir+"sscm.js")
}

// CommandHandler handles POST requests to execute commands via commandmgr.
// Expects a command in the request body. Returns 204 on success or error details.
func CommandHandler(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read command from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	command := strings.TrimSpace(string(body))

	// Validate command
	if command == "" {
		http.Error(w, "Command cannot be empty", http.StatusBadRequest)
		return
	}

	// Execute command via commandmgr
	if err := commandmgr.WriteCommand(command); err != nil {
		switch err {
		case os.ErrNotExist:
			http.Error(w, "Command file path not configured", http.StatusInternalServerError)
		case os.ErrInvalid:
			http.Error(w, "Invalid command", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to write command: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Success: return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func HandleIsSSCMEnabled(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if SSCM is enabled
	if !config.IsSSCMEnabled {
		http.Error(w, "SSCM is disabled", http.StatusForbidden)
		return
	}

	// Success: return 200 OK
	w.WriteHeader(http.StatusOK)
}
