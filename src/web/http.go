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

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/localization"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/commandmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/ssestream"
)

// TemplateData holds data to be passed to templates
type IndexTemplateData struct {
	Version                 string
	Branch                  string
	UIText_StartButton      string
	UIText_StopButton       string
	UIText_Settings         string
	UIText_Update_SteamCMD  string
	UIText_Console          string
	UIText_Detection_Events string
	UIText_Backup_Manager   string
	UIText_Discord_Info     string
	UIText_API_Info         string
	UIText_Copyright1       string
	UIText_Copyright2       string
}

// ConfigTemplateData holds data for the config page template
type ConfigTemplateData struct {
	// Config values
	DiscordToken                           string
	ControlChannelID                       string
	StatusChannelID                        string
	ConnectionListChannelID                string
	LogChannelID                           string
	SaveChannelID                          string
	ControlPanelChannelID                  string
	BlackListFilePath                      string
	ErrorChannelID                         string
	IsDiscordEnabled                       string
	IsDiscordEnabledTrueSelected           string
	IsDiscordEnabledFalseSelected          string
	GameBranch                             string
	Difficulty                             string
	StartCondition                         string
	StartLocation                          string
	ServerName                             string
	SaveInfo                               string
	ServerMaxPlayers                       string
	ServerPassword                         string
	ServerAuthSecret                       string
	AdminPassword                          string
	GamePort                               string
	UpdatePort                             string
	UPNPEnabled                            string
	UPNPEnabledTrueSelected                string
	UPNPEnabledFalseSelected               string
	AutoSave                               string
	AutoSaveTrueSelected                   string
	AutoSaveFalseSelected                  string
	SaveInterval                           string
	AutoPauseServer                        string
	AutoPauseServerTrueSelected            string
	AutoPauseServerFalseSelected           string
	LocalIpAddress                         string
	StartLocalHost                         string
	StartLocalHostTrueSelected             string
	StartLocalHostFalseSelected            string
	ServerVisible                          string
	ServerVisibleTrueSelected              string
	ServerVisibleFalseSelected             string
	UseSteamP2P                            string
	UseSteamP2PTrueSelected                string
	UseSteamP2PFalseSelected               string
	ExePath                                string
	AdditionalParams                       string
	AutoRestartServerTimer                 string
	IsNewTerrainAndSaveSystem              string
	IsNewTerrainAndSaveSystemTrueSelected  string
	IsNewTerrainAndSaveSystemFalseSelected string
	AutoStartServerOnStartup               string
	AutoStartServerOnStartupTrueSelected   string
	AutoStartServerOnStartupFalseSelected  string

	UIText_ServerConfig         string
	UIText_DiscordIntegration   string
	UIText_DetectionManager     string
	UIText_ConfigurationWizard  string
	UIText_PleaseSelectSection  string
	UIText_UseWizardAlternative string
	UIText_BasicSettings        string
	UIText_NetworkSettings      string
	UIText_AdvancedSettings     string
	UIText_BetaSettings         string
	UIText_BasicServerSettings  string

	UIText_ServerName                   string
	UIText_ServerNameInfo               string
	UIText_SaveFileName                 string
	UIText_SaveFileNameInfo             string
	UIText_MaxPlayers                   string
	UIText_MaxPlayersInfo               string
	UIText_ServerPassword               string
	UIText_ServerPasswordInfo           string
	UIText_AdminPassword                string
	UIText_AdminPasswordInfo            string
	UIText_AutoSave                     string
	UIText_AutoSaveInfo                 string
	UIText_SaveInterval                 string
	UIText_SaveIntervalInfo             string
	UIText_AutoPauseServer              string
	UIText_AutoPauseServerInfo          string
	UIText_NetworkConfiguration         string
	UIText_GamePort                     string
	UIText_GamePortInfo                 string
	UIText_UpdatePort                   string
	UIText_UpdatePortInfo               string
	UIText_UPNPEnabled                  string
	UIText_UPNPEnabledInfo              string
	UIText_LocalIpAddress               string
	UIText_LocalIpAddressInfo           string
	UIText_StartLocalHost               string
	UIText_StartLocalHostInfo           string
	UIText_ServerVisible                string
	UIText_ServerVisibleInfo            string
	UIText_UseSteamP2P                  string
	UIText_UseSteamP2PInfo              string
	UIText_AdvancedConfiguration        string
	UIText_ServerAuthSecret             string
	UIText_ServerAuthSecretInfo         string
	UIText_ServerExePath                string
	UIText_ServerExePathInfo            string
	UIText_ServerExePathInfo2           string
	UIText_AdditionalParams             string
	UIText_AdditionalParamsInfo         string
	UIText_AutoRestartServerTimer       string
	UIText_AutoRestartServerTimerInfo   string
	UIText_GameBranch                   string
	UIText_GameBranchInfo               string
	UIText_BetaOnlySettings             string
	UIText_BetaWarning                  string
	UIText_UseNewTerrainAndSave         string
	UIText_UseNewTerrainAndSaveInfo     string
	UIText_Difficulty                   string
	UIText_DifficultyInfo               string
	UIText_StartCondition               string
	UIText_StartConditionInfo           string
	UIText_StartLocation                string
	UIText_StartLocationInfo            string
	UIText_AutoStartServerOnStartup     string
	UIText_AutoStartServerOnStartupInfo string

	UIText_DiscordIntegrationTitle    string
	UIText_DiscordBotToken            string
	UIText_DiscordBotTokenInfo        string
	UIText_ChannelConfiguration       string
	UIText_AdminCommandChannel        string
	UIText_AdminCommandChannelInfo    string
	UIText_ControlPanelChannel        string
	UIText_ControlPanelChannelInfo    string
	UIText_StatusChannel              string
	UIText_StatusChannelInfo          string
	UIText_ConnectionListChannel      string
	UIText_ConnectionListChannelInfo  string
	UIText_LogChannel                 string
	UIText_LogChannelInfo             string
	UIText_SaveInfoChannel            string
	UIText_SaveInfoChannelInfo        string
	UIText_ErrorChannel               string
	UIText_ErrorChannelInfo           string
	UIText_BannedPlayersListPath      string
	UIText_BannedPlayersListPathInfo  string
	UIText_DiscordIntegrationBenefits string
	UIText_DiscordBenefit1            string
	UIText_DiscordBenefit2            string
	UIText_DiscordBenefit3            string
	UIText_DiscordBenefit4            string
	UIText_DiscordBenefit5            string
	UIText_DiscordSetupInstructions   string

	UIText_Copyright        string
	UIText_CopyrightConfig1 string
	UIText_CopyrightConfig2 string
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

	data := IndexTemplateData{
		Version:                 config.Version,
		Branch:                  config.Branch,
		UIText_StartButton:      localization.GetString("UIText_StartButton"),
		UIText_StopButton:       localization.GetString("UIText_StopButton"),
		UIText_Settings:         localization.GetString("UIText_Settings"),
		UIText_Update_SteamCMD:  localization.GetString("UIText_Update_SteamCMD"),
		UIText_Console:          localization.GetString("UIText_Console"),
		UIText_Detection_Events: localization.GetString("UIText_Detection_Events"),
		UIText_Backup_Manager:   localization.GetString("UIText_Backup_Manager"),
		UIText_Discord_Info:     localization.GetString("UIText_Discord_Info"),
		UIText_API_Info:         localization.GetString("UIText_API_Info"),
		UIText_Copyright1:       localization.GetString("UIText_Copyright1"),
		UIText_Copyright2:       localization.GetString("UIText_Copyright2"),
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

	tmpl, err := template.ParseFS(htmlFS, "config.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Core.Error("failed to serve config.html")
		return
	}

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

	autoStartServerTrueSelected := ""
	autoStartServerFalseSelected := ""
	if config.AutoStartServerOnStartup {
		autoStartServerTrueSelected = "selected"
	} else {
		autoStartServerFalseSelected = "selected"
	}

	steamP2PTrueSelected := ""
	steamP2PFalseSelected := ""
	if config.UseSteamP2P {
		steamP2PTrueSelected = "selected"
	} else {
		steamP2PFalseSelected = "selected"
	}

	data := ConfigTemplateData{
		// Config values
		DiscordToken:                           config.DiscordToken,
		ControlChannelID:                       config.ControlChannelID,
		StatusChannelID:                        config.StatusChannelID,
		ConnectionListChannelID:                config.ConnectionListChannelID,
		LogChannelID:                           config.LogChannelID,
		SaveChannelID:                          config.SaveChannelID,
		ControlPanelChannelID:                  config.ControlPanelChannelID,
		BlackListFilePath:                      config.BlackListFilePath,
		ErrorChannelID:                         config.ErrorChannelID,
		IsDiscordEnabled:                       fmt.Sprintf("%v", config.IsDiscordEnabled),
		IsDiscordEnabledTrueSelected:           discordTrueSelected,
		IsDiscordEnabledFalseSelected:          discordFalseSelected,
		GameBranch:                             config.GameBranch,
		Difficulty:                             config.Difficulty,
		StartCondition:                         config.StartCondition,
		StartLocation:                          config.StartLocation,
		ServerName:                             config.ServerName,
		SaveInfo:                               config.SaveInfo,
		ServerMaxPlayers:                       config.ServerMaxPlayers,
		ServerPassword:                         config.ServerPassword,
		ServerAuthSecret:                       config.ServerAuthSecret,
		AdminPassword:                          config.AdminPassword,
		GamePort:                               config.GamePort,
		UpdatePort:                             config.UpdatePort,
		UPNPEnabled:                            fmt.Sprintf("%v", config.UPNPEnabled),
		UPNPEnabledTrueSelected:                upnpTrueSelected,
		UPNPEnabledFalseSelected:               upnpFalseSelected,
		AutoSave:                               fmt.Sprintf("%v", config.AutoSave),
		AutoSaveTrueSelected:                   autoSaveTrueSelected,
		AutoSaveFalseSelected:                  autoSaveFalseSelected,
		SaveInterval:                           config.SaveInterval,
		AutoPauseServer:                        fmt.Sprintf("%v", config.AutoPauseServer),
		AutoPauseServerTrueSelected:            autoPauseTrueSelected,
		AutoPauseServerFalseSelected:           autoPauseFalseSelected,
		LocalIpAddress:                         config.LocalIpAddress,
		StartLocalHost:                         fmt.Sprintf("%v", config.StartLocalHost),
		StartLocalHostTrueSelected:             startLocalTrueSelected,
		StartLocalHostFalseSelected:            startLocalFalseSelected,
		ServerVisible:                          fmt.Sprintf("%v", config.ServerVisible),
		ServerVisibleTrueSelected:              serverVisibleTrueSelected,
		ServerVisibleFalseSelected:             serverVisibleFalseSelected,
		UseSteamP2P:                            fmt.Sprintf("%v", config.UseSteamP2P),
		UseSteamP2PTrueSelected:                steamP2PTrueSelected,
		UseSteamP2PFalseSelected:               steamP2PFalseSelected,
		ExePath:                                config.ExePath,
		AdditionalParams:                       config.AdditionalParams,
		AutoRestartServerTimer:                 config.AutoRestartServerTimer,
		IsNewTerrainAndSaveSystem:              fmt.Sprintf("%v", config.IsNewTerrainAndSaveSystem),
		IsNewTerrainAndSaveSystemTrueSelected:  isNewTerrainAndSaveSystemTrueSelected,
		IsNewTerrainAndSaveSystemFalseSelected: isNewTerrainAndSaveSystemFalseSelected,
		AutoStartServerOnStartup:               fmt.Sprintf("%v", config.AutoStartServerOnStartup),
		AutoStartServerOnStartupTrueSelected:   autoStartServerTrueSelected,
		AutoStartServerOnStartupFalseSelected:  autoStartServerFalseSelected,

		// Localized UI text
		UIText_ServerConfig:         localization.GetString("UIText_ServerConfig"),
		UIText_DiscordIntegration:   localization.GetString("UIText_DiscordIntegration"),
		UIText_DetectionManager:     localization.GetString("UIText_DetectionManager"),
		UIText_ConfigurationWizard:  localization.GetString("UIText_ConfigurationWizard"),
		UIText_PleaseSelectSection:  localization.GetString("UIText_PleaseSelectSection"),
		UIText_UseWizardAlternative: localization.GetString("UIText_UseWizardAlternative"),
		UIText_BasicSettings:        localization.GetString("UIText_BasicSettings"),
		UIText_NetworkSettings:      localization.GetString("UIText_NetworkSettings"),
		UIText_AdvancedSettings:     localization.GetString("UIText_AdvancedSettings"),
		UIText_BetaSettings:         localization.GetString("UIText_BetaSettings"),
		UIText_BasicServerSettings:  localization.GetString("UIText_BasicServerSettings"),

		UIText_ServerName:                   localization.GetString("UIText_ServerName"),
		UIText_ServerNameInfo:               localization.GetString("UIText_ServerNameInfo"),
		UIText_SaveFileName:                 localization.GetString("UIText_SaveFileName"),
		UIText_SaveFileNameInfo:             localization.GetString("UIText_SaveFileNameInfo"),
		UIText_MaxPlayers:                   localization.GetString("UIText_MaxPlayers"),
		UIText_MaxPlayersInfo:               localization.GetString("UIText_MaxPlayersInfo"),
		UIText_ServerPassword:               localization.GetString("UIText_ServerPassword"),
		UIText_ServerPasswordInfo:           localization.GetString("UIText_ServerPasswordInfo"),
		UIText_AdminPassword:                localization.GetString("UIText_AdminPassword"),
		UIText_AdminPasswordInfo:            localization.GetString("UIText_AdminPasswordInfo"),
		UIText_AutoSave:                     localization.GetString("UIText_AutoSave"),
		UIText_AutoSaveInfo:                 localization.GetString("UIText_AutoSaveInfo"),
		UIText_SaveInterval:                 localization.GetString("UIText_SaveInterval"),
		UIText_SaveIntervalInfo:             localization.GetString("UIText_SaveIntervalInfo"),
		UIText_AutoPauseServer:              localization.GetString("UIText_AutoPauseServer"),
		UIText_AutoPauseServerInfo:          localization.GetString("UIText_AutoPauseServerInfo"),
		UIText_NetworkConfiguration:         localization.GetString("UIText_NetworkConfiguration"),
		UIText_GamePort:                     localization.GetString("UIText_GamePort"),
		UIText_GamePortInfo:                 localization.GetString("UIText_GamePortInfo"),
		UIText_UpdatePort:                   localization.GetString("UIText_UpdatePort"),
		UIText_UpdatePortInfo:               localization.GetString("UIText_UpdatePortInfo"),
		UIText_UPNPEnabled:                  localization.GetString("UIText_UPNPEnabled"),
		UIText_UPNPEnabledInfo:              localization.GetString("UIText_UPNPEnabledInfo"),
		UIText_LocalIpAddress:               localization.GetString("UIText_LocalIpAddress"),
		UIText_LocalIpAddressInfo:           localization.GetString("UIText_LocalIpAddressInfo"),
		UIText_StartLocalHost:               localization.GetString("UIText_StartLocalHost"),
		UIText_StartLocalHostInfo:           localization.GetString("UIText_StartLocalHostInfo"),
		UIText_ServerVisible:                localization.GetString("UIText_ServerVisible"),
		UIText_ServerVisibleInfo:            localization.GetString("UIText_ServerVisibleInfo"),
		UIText_UseSteamP2P:                  localization.GetString("UIText_UseSteamP2P"),
		UIText_UseSteamP2PInfo:              localization.GetString("UIText_UseSteamP2PInfo"),
		UIText_AdvancedConfiguration:        localization.GetString("UIText_AdvancedConfiguration"),
		UIText_ServerAuthSecret:             localization.GetString("UIText_ServerAuthSecret"),
		UIText_ServerAuthSecretInfo:         localization.GetString("UIText_ServerAuthSecretInfo"),
		UIText_ServerExePath:                localization.GetString("UIText_ServerExePath"),
		UIText_ServerExePathInfo:            localization.GetString("UIText_ServerExePathInfo"),
		UIText_ServerExePathInfo2:           localization.GetString("UIText_ServerExePathInfo2"),
		UIText_AdditionalParams:             localization.GetString("UIText_AdditionalParams"),
		UIText_AdditionalParamsInfo:         localization.GetString("UIText_AdditionalParamsInfo"),
		UIText_AutoRestartServerTimer:       localization.GetString("UIText_AutoRestartServerTimer"),
		UIText_AutoRestartServerTimerInfo:   localization.GetString("UIText_AutoRestartServerTimerInfo"),
		UIText_GameBranch:                   localization.GetString("UIText_GameBranch"),
		UIText_GameBranchInfo:               localization.GetString("UIText_GameBranchInfo"),
		UIText_BetaOnlySettings:             localization.GetString("UIText_BetaOnlySettings"),
		UIText_BetaWarning:                  localization.GetString("UIText_BetaWarning"),
		UIText_UseNewTerrainAndSave:         localization.GetString("UIText_UseNewTerrainAndSave"),
		UIText_UseNewTerrainAndSaveInfo:     localization.GetString("UIText_UseNewTerrainAndSaveInfo"),
		UIText_Difficulty:                   localization.GetString("UIText_Difficulty"),
		UIText_DifficultyInfo:               localization.GetString("UIText_DifficultyInfo"),
		UIText_StartCondition:               localization.GetString("UIText_StartCondition"),
		UIText_StartConditionInfo:           localization.GetString("UIText_StartConditionInfo"),
		UIText_StartLocation:                localization.GetString("UIText_StartLocation"),
		UIText_StartLocationInfo:            localization.GetString("UIText_StartLocationInfo"),
		UIText_AutoStartServerOnStartup:     localization.GetString("UIText_AutoStartServerOnStartup"),
		UIText_AutoStartServerOnStartupInfo: localization.GetString("UIText_AutoStartServerOnStartupInfo"),

		UIText_DiscordIntegrationTitle:    localization.GetString("UIText_DiscordIntegrationTitle"),
		UIText_DiscordBotToken:            localization.GetString("UIText_DiscordBotToken"),
		UIText_DiscordBotTokenInfo:        localization.GetString("UIText_DiscordBotTokenInfo"),
		UIText_ChannelConfiguration:       localization.GetString("UIText_ChannelConfiguration"),
		UIText_AdminCommandChannel:        localization.GetString("UIText_AdminCommandChannel"),
		UIText_AdminCommandChannelInfo:    localization.GetString("UIText_AdminCommandChannelInfo"),
		UIText_ControlPanelChannel:        localization.GetString("UIText_ControlPanelChannel"),
		UIText_ControlPanelChannelInfo:    localization.GetString("UIText_ControlPanelChannelInfo"),
		UIText_StatusChannel:              localization.GetString("UIText_StatusChannel"),
		UIText_StatusChannelInfo:          localization.GetString("UIText_StatusChannelInfo"),
		UIText_ConnectionListChannel:      localization.GetString("UIText_ConnectionListChannel"),
		UIText_ConnectionListChannelInfo:  localization.GetString("UIText_ConnectionListChannelInfo"),
		UIText_LogChannel:                 localization.GetString("UIText_LogChannel"),
		UIText_LogChannelInfo:             localization.GetString("UIText_LogChannelInfo"),
		UIText_SaveInfoChannel:            localization.GetString("UIText_SaveInfoChannel"),
		UIText_SaveInfoChannelInfo:        localization.GetString("UIText_SaveInfoChannelInfo"),
		UIText_ErrorChannel:               localization.GetString("UIText_ErrorChannel"),
		UIText_ErrorChannelInfo:           localization.GetString("UIText_ErrorChannelInfo"),
		UIText_BannedPlayersListPath:      localization.GetString("UIText_BannedPlayersListPath"),
		UIText_BannedPlayersListPathInfo:  localization.GetString("UIText_BannedPlayersListPathInfo"),
		UIText_DiscordIntegrationBenefits: localization.GetString("UIText_DiscordIntegrationBenefits"),
		UIText_DiscordBenefit1:            localization.GetString("UIText_DiscordBenefit1"),
		UIText_DiscordBenefit2:            localization.GetString("UIText_DiscordBenefit2"),
		UIText_DiscordBenefit3:            localization.GetString("UIText_DiscordBenefit3"),
		UIText_DiscordBenefit4:            localization.GetString("UIText_DiscordBenefit4"),
		UIText_DiscordBenefit5:            localization.GetString("UIText_DiscordBenefit5"),
		UIText_DiscordSetupInstructions:   localization.GetString("UIText_DiscordSetupInstructions"),

		UIText_CopyrightConfig1: localization.GetString("UIText_Copyright1"),
		UIText_CopyrightConfig2: localization.GetString("UIText_Copyright2"),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
