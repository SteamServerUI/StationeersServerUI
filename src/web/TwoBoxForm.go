package web

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/localization"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

func ServeTwoBoxFormTemplate(w http.ResponseWriter, r *http.Request) {
	type Step struct {
		ID                       string
		Title                    string
		HeaderTitle              string
		StepMessage              string
		PrimaryLabel             string
		SecondaryLabel           string
		SecondaryLabelType       string
		SubmitButtonText         string
		SkipButtonText           string
		PrimaryPlaceholderText   string
		SecondaryPlaceholderText string
		ConfigField              string // Field name to save in config
		NextStep                 string // ID of the next step
	}

	type TemplateData struct {
		IsFirstTimeSetup         bool
		Path                     string
		Title                    string
		HeaderTitle              string
		StepMessage              string
		PrimaryLabel             string
		SecondaryLabel           string
		SecondaryLabelType       string
		SubmitButtonText         string
		SkipButtonText           string
		Mode                     string
		ShowExtraButtons         bool
		FooterText               string
		Step                     string
		ConfigField              string
		NextStep                 string
		PrimaryPlaceholderText   string
		SecondaryPlaceholderText string
	}

	twoboxformAssetsFS, err := fs.Sub(config.GetV1UIFS(), "UIMod/onboard_bundled/twoboxform")
	if err != nil {
		logger.Web.Error("Failed to get bundled FS")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(twoboxformAssetsFS, "twoboxform.html")
	if err != nil {
		logger.Web.Error("Failed to parse 2BoxForm template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	path := r.URL.Path
	stepID := r.URL.Query().Get("step")
	if stepID == "" && path == "/setup" {
		stepID = "welcome" // Start with welcome page for first-time setup
	}

	// Define all steps in a map for easy access and modification
	steps := map[string]Step{
		"welcome": {
			ID:                 "welcome",
			Title:              localization.GetString("UIText_Welcome_Title"),
			HeaderTitle:        "",
			StepMessage:        "",
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   localization.GetString("UIText_Welcome_SubmitButton"),
			SkipButtonText:     localization.GetString("UIText_Welcome_SkipButton"),
			NextStep:           "pls_read",
		},
		"pls_read": {
			ID:                 "pls_read",
			Title:              localization.GetString("UIText_PlsRead_Title"),
			HeaderTitle:        localization.GetString("UIText_PlsRead_HeaderTitle"),
			StepMessage:        localization.GetString("UIText_PlsRead_StepMessage"),
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   localization.GetString("UIText_PlsRead_SubmitButton"),
			SkipButtonText:     localization.GetString("UIText_PlsRead_SkipButton"),
			NextStep:           "server_name",
		},
		"server_name": {
			ID:                     "server_name",
			Title:                  localization.GetString("UIText_ServerName_Title"),
			HeaderTitle:            localization.GetString("UIText_ServerName_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_ServerName_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_ServerName_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_ServerName_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_ServerName_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_ServerName_SkipButton"),
			ConfigField:            "ServerName",
			NextStep:               "save_identifier",
		},
		"save_identifier": {
			ID:                     "save_identifier",
			Title:                  localization.GetString("UIText_SaveIdentifier_Title"),
			HeaderTitle:            localization.GetString("UIText_SaveIdentifier_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_SaveIdentifier_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_SaveIdentifier_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_SaveIdentifier_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_SaveIdentifier_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_SaveIdentifier_SkipButton"),
			ConfigField:            "SaveInfo",
			NextStep:               "max_players",
		},
		"max_players": {
			ID:                     "max_players",
			Title:                  localization.GetString("UIText_MaxPlayers_Title"),
			HeaderTitle:            localization.GetString("UIText_MaxPlayers_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_MaxPlayers_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_MaxPlayers_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_MaxPlayers_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_MaxPlayers_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_MaxPlayers_SkipButton"),
			ConfigField:            "ServerMaxPlayers",
			NextStep:               "server_password",
		},
		"server_password": {
			ID:                     "server_password",
			Title:                  localization.GetString("UIText_ServerPassword_Title"),
			HeaderTitle:            localization.GetString("UIText_ServerPassword_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_ServerPassword_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_ServerPassword_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_ServerPassword_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_ServerPassword_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_ServerPassword_SkipButton"),
			ConfigField:            "ServerPassword",
			NextStep:               "game_branch",
		},
		"game_branch": {
			ID:                     "game_branch",
			Title:                  localization.GetString("UIText_GameBranch_Title"),
			HeaderTitle:            localization.GetString("UIText_GameBranch_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_GameBranch_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_GameBranch_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_GameBranch_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_GameBranch_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_GameBranch_SkipButton"),
			ConfigField:            "gameBranch",
			NextStep:               "newterrain_and_savesystem",
		},
		"newterrain_and_savesystem": {
			ID:                     "newterrain_and_savesystem",
			Title:                  localization.GetString("UIText_NewTerrainAndSaveSystem_Title"),
			HeaderTitle:            localization.GetString("UIText_NewTerrainAndSaveSystem_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_NewTerrainAndSaveSystem_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_NewTerrainAndSaveSystem_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_NewTerrainAndSaveSystem_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_NewTerrainAndSaveSystem_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_NewTerrainAndSaveSystem_SkipButton"),
			ConfigField:            "IsNewTerrainAndSaveSystem",
			NextStep:               "network_config_choice",
		},
		"discord_enabled": {
			ID:                     "discord_enabled",
			Title:                  localization.GetString("UIText_DiscordEnabled_Title"),
			HeaderTitle:            localization.GetString("UIText_DiscordEnabled_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_DiscordEnabled_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_DiscordEnabled_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_DiscordEnabled_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_DiscordEnabled_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_DiscordEnabled_SkipButton"),
			ConfigField:            "isDiscordEnabled", // We'll handle the boolean conversion in JS
			NextStep:               "discord_token",    // Default next step if enabled
			// The actual next step will be determined by JS based on the answer
		},
		"discord_token": {
			ID:                     "discord_token",
			Title:                  localization.GetString("UIText_DiscordToken_Title"),
			HeaderTitle:            localization.GetString("UIText_DiscordToken_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_DiscordToken_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_DiscordToken_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_DiscordToken_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_DiscordToken_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_DiscordToken_SkipButton"),
			ConfigField:            "discordToken",
			NextStep:               "control_panel_channel",
		},
		"control_panel_channel": {
			ID:                     "control_panel_channel",
			Title:                  localization.GetString("UIText_ControlPanelChannel_Title"),
			HeaderTitle:            localization.GetString("UIText_ControlPanelChannel_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_ControlPanelChannel_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_ControlPanelChannel_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_ControlPanelChannel_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_ControlPanelChannel_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_ControlPanelChannel_SkipButton"),
			ConfigField:            "controlPanelChannelID",
			NextStep:               "save_channel",
		},
		"save_channel": {
			ID:                     "save_channel",
			Title:                  localization.GetString("UIText_SaveChannel_Title"),
			HeaderTitle:            localization.GetString("UIText_SaveChannel_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_SaveChannel_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_SaveChannel_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_SaveChannel_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_SaveChannel_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_SaveChannel_SkipButton"),
			ConfigField:            "saveChannelID",
			NextStep:               "log_channel",
		},
		"log_channel": {
			ID:                     "log_channel",
			Title:                  localization.GetString("UIText_LogChannel_Title"),
			HeaderTitle:            localization.GetString("UIText_LogChannel_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_LogChannel_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_LogChannel_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_LogChannel_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_LogChannel_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_LogChannel_SkipButton"),
			ConfigField:            "logChannelID",
			NextStep:               "connection_list_channel",
		},
		"connection_list_channel": {
			ID:                     "connection_list_channel",
			Title:                  localization.GetString("UIText_ConnectionListChannel_Title"),
			HeaderTitle:            localization.GetString("UIText_ConnectionListChannel_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_ConnectionListChannel_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_ConnectionListChannel_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_ConnectionListChannel_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_ConnectionListChannel_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_ConnectionListChannel_SkipButton"),
			ConfigField:            "connectionListChannelID",
			NextStep:               "status_channel",
		},
		"status_channel": {
			ID:                     "status_channel",
			Title:                  localization.GetString("UIText_StatusChannel_Title"),
			HeaderTitle:            localization.GetString("UIText_StatusChannel_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_StatusChannel_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_StatusChannel_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_StatusChannel_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_StatusChannel_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_StatusChannel_SkipButton"),
			ConfigField:            "statusChannelID",
			NextStep:               "control_channel",
		},
		"control_channel": {
			ID:                     "control_channel",
			Title:                  localization.GetString("UIText_ControlChannel_Title"),
			HeaderTitle:            localization.GetString("UIText_ControlChannel_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_ControlChannel_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_ControlChannel_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_ControlChannel_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_ControlChannel_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_ControlChannel_SkipButton"),
			ConfigField:            "controlChannelID",
			NextStep:               "network_config_choice",
		},
		"network_config_choice": {
			ID:                     "network_config_choice",
			Title:                  localization.GetString("UIText_NetworkConfigChoice_Title"),
			HeaderTitle:            localization.GetString("UIText_NetworkConfigChoice_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_NetworkConfigChoice_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_NetworkConfigChoice_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_NetworkConfigChoice_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_NetworkConfigChoice_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_NetworkConfigChoice_SkipButton"),
			ConfigField:            "",          // No config field, just for branching
			NextStep:               "game_port", // Default next step if they choose to configure
			// The actual next step will be determined by JS based on the answer
		},

		"game_port": {
			ID:                     "game_port",
			Title:                  localization.GetString("UIText_GamePort_Title"),
			HeaderTitle:            localization.GetString("UIText_GamePort_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_GamePort_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_GamePort_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_GamePort_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_GamePort_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_GamePort_SkipButton"),
			ConfigField:            "GamePort",
			NextStep:               "update_port",
		},
		"update_port": {
			ID:                     "update_port",
			Title:                  localization.GetString("UIText_UpdatePort_Title"),
			HeaderTitle:            localization.GetString("UIText_UpdatePort_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_UpdatePort_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_UpdatePort_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_UpdatePort_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_UpdatePort_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_UpdatePort_SkipButton"),
			ConfigField:            "UpdatePort",
			NextStep:               "upnp_enabled",
		},
		"upnp_enabled": {
			ID:                     "upnp_enabled",
			Title:                  localization.GetString("UIText_UPnPEnabled_Title"),
			HeaderTitle:            localization.GetString("UIText_UPnPEnabled_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_UPnPEnabled_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_UPnPEnabled_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_UPnPEnabled_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_UPnPEnabled_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_UPnPEnabled_SkipButton"),
			ConfigField:            "UPNPEnabled",
			NextStep:               "local_ip_address",
		},
		//skipped for now, not needed and bugged, wontfix for now
		"server_visible": {
			ID:                     "server_visible",
			Title:                  localization.GetString("UIText_ServerVisible_Title"),
			HeaderTitle:            localization.GetString("UIText_ServerVisible_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_ServerVisible_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_ServerVisible_PrimaryPlaceholder"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_ServerVisible_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_ServerVisible_SkipButton"),
			ConfigField:            "ServerVisible", // We'll handle the boolean conversion in JS
			NextStep:               "local_ip_address",
		},
		"local_ip_address": {
			ID:                     "local_ip_address",
			Title:                  localization.GetString("UIText_LocalIPAddress_Title"),
			HeaderTitle:            localization.GetString("UIText_LocalIPAddress_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_LocalIPAddress_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_LocalIPAddress_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_LocalIPAddress_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_LocalIPAddress_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_LocalIPAddress_SkipButton"),
			ConfigField:            "LocalIpAddress",
			NextStep:               "admin_account",
		},
		"admin_account": {
			ID:                       "admin_account",
			Title:                    localization.GetString("UIText_AdminAccount_Title"),
			HeaderTitle:              localization.GetString("UIText_AdminAccount_HeaderTitle"),
			StepMessage:              localization.GetString("UIText_AdminAccount_StepMessage"),
			PrimaryPlaceholderText:   localization.GetString("UIText_AdminAccount_PrimaryPlaceholder"),
			PrimaryLabel:             localization.GetString("UIText_AdminAccount_PrimaryLabel"),
			SecondaryLabel:           localization.GetString("UIText_AdminAccount_SecondaryLabel"),
			SecondaryPlaceholderText: localization.GetString("UIText_AdminAccount_SecondaryPlaceholder"),
			SecondaryLabelType:       "password",
			SubmitButtonText:         localization.GetString("UIText_AdminAccount_SubmitButton"),
			SkipButtonText:           localization.GetString("UIText_AdminAccount_SkipButton"),
			ConfigField:              "",
			NextStep:                 "finalize",
		},
		"sscm": {
			ID:                     "sscm",
			Title:                  localization.GetString("UIText_SSCM_Title"),
			HeaderTitle:            localization.GetString("UIText_SSCM_HeaderTitle"),
			StepMessage:            localization.GetString("UIText_SSCM_StepMessage"),
			PrimaryPlaceholderText: localization.GetString("UIText_SSCM_PrimaryPlaceholder"),
			PrimaryLabel:           localization.GetString("UIText_SSCM_PrimaryLabel"),
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       localization.GetString("UIText_SSCM_SubmitButton"),
			SkipButtonText:         localization.GetString("UIText_SSCM_SkipButton"),
			ConfigField:            "IsSSCMEnabled",
			NextStep:               "finalize",
		},
		"finalize": {
			ID:               "finalize",
			Title:            localization.GetString("UIText_Finalize_Title"),
			HeaderTitle:      "",
			StepMessage:      localization.GetString("UIText_Finalize_StepMessage"),
			PrimaryLabel:     "",
			SecondaryLabel:   "",
			SubmitButtonText: localization.GetString("UIText_Finalize_SubmitButton"),
			SkipButtonText:   localization.GetString("UIText_Finalize_SkipButton"),
			NextStep:         "welcome", // Return to first step if "Return to Setup" is clicked
		},
	}

	data := TemplateData{
		IsFirstTimeSetup: config.IsFirstTimeSetup,
		Path:             path,
		Step:             stepID,
		FooterText:       localization.GetString("UIText_FooterText"),
	}

	switch {

	case path == "/login" && !config.AuthEnabled:
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return

	case path == "/setup":
		data.Mode = "setup"
		data.ShowExtraButtons = true

		// Find the current step in our map
		if step, exists := steps[stepID]; exists {
			data.Title = step.Title
			data.HeaderTitle = step.HeaderTitle
			data.StepMessage = step.StepMessage
			data.PrimaryLabel = step.PrimaryLabel
			data.SecondaryLabel = step.SecondaryLabel
			data.SecondaryLabelType = step.SecondaryLabelType
			data.SubmitButtonText = step.SubmitButtonText
			data.SkipButtonText = step.SkipButtonText
			data.ConfigField = step.ConfigField
			data.NextStep = step.NextStep
			data.PrimaryPlaceholderText = step.PrimaryPlaceholderText
			data.SecondaryPlaceholderText = step.SecondaryPlaceholderText
			if stepID == "sscm" {
				data.FooterText = localization.GetString("UIText_SSCM_FooterText")
			}
		} else {
			// Default to welcome page if step is invalid
			welcomeStep := steps["welcome"]
			data.Title = welcomeStep.Title
			data.HeaderTitle = welcomeStep.HeaderTitle
			data.StepMessage = welcomeStep.StepMessage
			data.PrimaryLabel = welcomeStep.PrimaryLabel
			data.SecondaryLabel = welcomeStep.SecondaryLabel
			data.SecondaryLabelType = welcomeStep.SecondaryLabelType
			data.SubmitButtonText = welcomeStep.SubmitButtonText
			data.SkipButtonText = welcomeStep.SkipButtonText
			data.ConfigField = welcomeStep.ConfigField
			data.NextStep = welcomeStep.NextStep
			data.Step = "welcome"
		}

	case path == "/changeuser":
		data.Title = localization.GetString("UIText_ChangeUser_Title")
		data.HeaderTitle = localization.GetString("UIText_ChangeUser_HeaderTitle")
		data.PrimaryLabel = localization.GetString("UIText_ChangeUser_PrimaryLabel")
		data.SecondaryLabel = localization.GetString("UIText_ChangeUser_SecondaryLabel")
		data.SecondaryPlaceholderText = localization.GetString("UIText_ChangeUser_SecondaryPlaceholder")
		data.SecondaryLabelType = "password"
		data.SubmitButtonText = localization.GetString("UIText_ChangeUser_SubmitButton")
		data.Mode = "changeuser"
		data.ShowExtraButtons = false

	default:
		data.Title = localization.GetString("UIText_Login_Title")
		data.HeaderTitle = localization.GetString("UIText_Login_HeaderTitle")
		data.PrimaryLabel = localization.GetString("UIText_Login_PrimaryLabel")
		data.SecondaryLabel = localization.GetString("UIText_Login_SecondaryLabel")
		data.PrimaryPlaceholderText = localization.GetString("UIText_Login_PrimaryPlaceholder")
		data.SecondaryPlaceholderText = localization.GetString("UIText_Login_SecondaryPlaceholder")
		data.SecondaryLabelType = "password"
		data.SubmitButtonText = localization.GetString("UIText_Login_SubmitButton")
		data.Mode = "login"
		data.ShowExtraButtons = false
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		logger.Web.Error("Failed to execute 2BoxForm template: %v" + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
