package web

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/commandmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/gamemgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/ssestream"
)

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Version string
	Branch  string
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	htmlFS, err := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/ui")
	if err != nil {
		http.Error(w, "Error accessing Virt FS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(htmlFS, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Core.Error("failed to serve v1 Index.html")
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
	detectionmanagerFS, err := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/detectionmanager")
	if err != nil {
		http.Error(w, "Error accessing Virt FS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlFile, err := detectionmanagerFS.Open("detectionmanager.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading detectionmanager.html: %v", err), http.StatusInternalServerError)
		return
	}
	defer htmlFile.Close()

	htmlContent, err := io.ReadAll(htmlFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading detectionmanager.html content: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(htmlContent)
}

func ServeConfigPage(w http.ResponseWriter, r *http.Request) {

	htmlFS, err := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/ui")
	if err != nil {
		http.Error(w, "Error accessing Virt FS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlFile, err := htmlFS.Open("config.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading config.html: %v", err), http.StatusInternalServerError)
		return
	}
	defer htmlFile.Close()

	htmlContentBytes, err := io.ReadAll(htmlFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading config.html content: %v", err), http.StatusInternalServerError)
		return
	}

	htmlContent := string(htmlContentBytes)

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

	isNewTerrainAndSaveSystemTrueSelected := ""
	isNewTerrainAndSaveSystemFalseSelected := ""

	if config.IsNewTerrainAndSaveSystem {
		isNewTerrainAndSaveSystemTrueSelected = "selected"
	} else {
		isNewTerrainAndSaveSystemFalseSelected = "selected"
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
		"{{discordToken}}":                           config.DiscordToken,
		"{{controlChannelID}}":                       config.ControlChannelID,
		"{{statusChannelID}}":                        config.StatusChannelID,
		"{{connectionListChannelID}}":                config.ConnectionListChannelID,
		"{{logChannelID}}":                           config.LogChannelID,
		"{{saveChannelID}}":                          config.SaveChannelID,
		"{{controlPanelChannelID}}":                  config.ControlPanelChannelID,
		"{{blackListFilePath}}":                      config.BlackListFilePath,
		"{{errorChannelID}}":                         config.ErrorChannelID,
		"{{isDiscordEnabled}}":                       fmt.Sprintf("%v", config.IsDiscordEnabled),
		"{{IsDiscordEnabledTrueSelected}}":           discordTrueSelected,
		"{{IsDiscordEnabledFalseSelected}}":          discordFalseSelected,
		"{{gameBranch}}":                             config.GameBranch,
		"{{Difficulty}}":                             config.Difficulty,
		"{{StartCondition}}":                         config.StartCondition,
		"{{StartLocation}}":                          config.StartLocation,
		"{{ServerName}}":                             config.ServerName,
		"{{SaveInfo}}":                               config.SaveInfo,
		"{{ServerMaxPlayers}}":                       config.ServerMaxPlayers,
		"{{ServerPassword}}":                         config.ServerPassword,
		"{{ServerAuthSecret}}":                       config.ServerAuthSecret,
		"{{AdminPassword}}":                          config.AdminPassword,
		"{{GamePort}}":                               config.GamePort,
		"{{UpdatePort}}":                             config.UpdatePort,
		"{{UPNPEnabled}}":                            fmt.Sprintf("%v", config.UPNPEnabled),
		"{{UPNPEnabledTrueSelected}}":                upnpTrueSelected,
		"{{UPNPEnabledFalseSelected}}":               upnpFalseSelected,
		"{{AutoSave}}":                               fmt.Sprintf("%v", config.AutoSave),
		"{{AutoSaveTrueSelected}}":                   autoSaveTrueSelected,
		"{{AutoSaveFalseSelected}}":                  autoSaveFalseSelected,
		"{{SaveInterval}}":                           config.SaveInterval,
		"{{AutoPauseServer}}":                        fmt.Sprintf("%v", config.AutoPauseServer),
		"{{AutoPauseServerTrueSelected}}":            autoPauseTrueSelected,
		"{{AutoPauseServerFalseSelected}}":           autoPauseFalseSelected,
		"{{LocalIpAddress}}":                         config.LocalIpAddress,
		"{{StartLocalHost}}":                         fmt.Sprintf("%v", config.StartLocalHost),
		"{{StartLocalHostTrueSelected}}":             startLocalTrueSelected,
		"{{StartLocalHostFalseSelected}}":            startLocalFalseSelected,
		"{{ServerVisible}}":                          fmt.Sprintf("%v", config.ServerVisible),
		"{{ServerVisibleTrueSelected}}":              serverVisibleTrueSelected,
		"{{ServerVisibleFalseSelected}}":             serverVisibleFalseSelected,
		"{{UseSteamP2P}}":                            fmt.Sprintf("%v", config.UseSteamP2P),
		"{{UseSteamP2PTrueSelected}}":                steamP2PTrueSelected,
		"{{UseSteamP2PFalseSelected}}":               steamP2PFalseSelected,
		"{{ExePath}}":                                config.ExePath,
		"{{AdditionalParams}}":                       config.AdditionalParams,
		"{{AutoRestartServerTimer}}":                 config.AutoRestartServerTimer,
		"{{IsNewTerrainAndSaveSystem}}":              fmt.Sprintf("%v", config.IsNewTerrainAndSaveSystem),
		"{{IsNewTerrainAndSaveSystemTrueSelected}}":  isNewTerrainAndSaveSystemTrueSelected,
		"{{IsNewTerrainAndSaveSystemFalseSelected}}": isNewTerrainAndSaveSystemFalseSelected,
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

var lastSteamCMDExecution time.Time // last time SteamCMD was executed via API.

// run SteamCMD from API, but only allow once every 5 minutes to "kinda" prevent concurrent executions although that woluldnt hurn.
// If the user has a 5mbit connection, I cannot help them anyways.
func HandleRunSteamCMD(w http.ResponseWriter, r *http.Request) {
	const rateLimitDuration = 30 * time.Second

	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check rate limit
	if time.Since(lastSteamCMDExecution) < rateLimitDuration {
		json.NewEncoder(w).Encode(map[string]string{"statuscode": "200", "status": "Rejected", "message": "Slow down, you just called SteamCMD.", "advanced": "Use SSUICLI or restart SSUI to run SteamCMD repeatedly without limit."})
		return
	}

	if gamemgr.InternalIsServerRunning() {
		logger.Core.Warn("Server is running, stopping server first...")
		gamemgr.InternalStopServer()
		time.Sleep(10000 * time.Millisecond)
	}
	logger.Core.Info("Running SteamCMD")
	setup.InstallAndRunSteamCMD()

	// Update last execution time
	lastSteamCMDExecution = time.Now()

	// Success: return 202 Accepted and JSON
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"statuscode": "202", "status": "Accepted", "message": "SteamCMD ran successfully."})
}
