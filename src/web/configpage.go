package web

import (
	"fmt"
	"io/fs"
	"net/http"
	"text/template"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/localization"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

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
		JXSNDevApiKey:                          config.JXSNDevApiKey,
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
		UIText_JXSNDevApiKey:                localization.GetString("UIText_JXSNDevApiKey"),
		UIText_JXSNDevApiKeyInfo:            localization.GetString("UIText_JXSNDevApiKeyInfo"),

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
