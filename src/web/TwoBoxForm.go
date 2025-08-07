package web

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
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
			Title:              "Stationeers Server UI",
			HeaderTitle:        "",
			StepMessage:        "",
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   "Start Setup",
			SkipButtonText:     "Skip Setup",
			NextStep:           "pls_read",
		},
		"pls_read": {
			ID:                 "pls_read",
			Title:              "Please read!",
			HeaderTitle:        "We strongly recommend you to read the texts in this setup wizard!",
			StepMessage:        "Most reported issues occur because of a misconfiguration.",
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   "I understand",
			SkipButtonText:     "I understand",
			NextStep:           "server_name",
		},
		"server_name": {
			ID:                     "server_name",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Server Name Setup",
			StepMessage:            "Give your server a name like 'Space Station 13'",
			PrimaryPlaceholderText: "My Stationeers Server with UI",
			PrimaryLabel:           "Server Name",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "ServerName",
			NextStep:               "save_identifier",
		},
		"save_identifier": {
			ID:                     "save_identifier",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Save Identifier Setup",
			StepMessage:            "Set a save identifier like 'SpaceStation13 Vulcan'. Capitalize the first letter of each word. Possible World types can be found in the Stationeers Wiki -> Dedicated Server",
			PrimaryPlaceholderText: "Requires a SaveName and WorldType for first start!",
			PrimaryLabel:           "Save Identifier",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "SaveInfo",
			NextStep:               "max_players",
		},
		"max_players": {
			ID:                     "max_players",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Player Limit Setup",
			StepMessage:            "Choose the maximum number of players that can connect to the server.",
			PrimaryPlaceholderText: "8",
			PrimaryLabel:           "Max Players",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "ServerMaxPlayers",
			NextStep:               "server_password",
		},
		"server_password": {
			ID:                     "server_password",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Server Password Setup",
			StepMessage:            "Set a gameserver password or skip this step.",
			PrimaryPlaceholderText: "Server Password",
			PrimaryLabel:           "Server Password",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "ServerPassword",
			NextStep:               "game_branch",
		},
		"game_branch": {
			ID:                     "game_branch",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Game Branch Setup",
			StepMessage:            "Enter a beta branch or skip this to use the release version. If switching branches, make sure to r e s t a r t SSUI after completing this wizzard.",
			PrimaryPlaceholderText: "beta",
			PrimaryLabel:           "Game Branch",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Use Release Version",
			ConfigField:            "gameBranch",
			NextStep:               "newterrain_and_savesystem",
		},

		"newterrain_and_savesystem": {
			ID:                     "newterrain_and_savesystem",
			Title:                  "CHOOSE TERRAIN SYSTEM",
			HeaderTitle:            "Very important step!",
			StepMessage:            "Just switched to Beta? Flip Terrain and Save System to support that! Enter 'yes' to enable or 'no' to disable.",
			PrimaryPlaceholderText: "yes/no",
			PrimaryLabel:           "Enable new System",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "IsNewTerrainAndSaveSystem",
			NextStep:               "network_config_choice",
		},

		"discord_enabled": {
			ID:                     "discord_enabled",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Discord Integration",
			StepMessage:            "Do you want to enable Discord integration? Enter 'yes' to enable or Skip to disable.",
			PrimaryPlaceholderText: "yes",
			PrimaryLabel:           "Enable Discord",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip (Disable Discord)",
			ConfigField:            "isDiscordEnabled", // We'll handle the boolean conversion in JS
			NextStep:               "discord_token",    // Default next step if enabled
			// The actual next step will be determined by JS based on the answer
		},

		"discord_token": {
			ID:                     "discord_token",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Discord Bot Token",
			StepMessage:            "Enter your Discord bot token for server integration",
			PrimaryPlaceholderText: "Discord Bot Token",
			PrimaryLabel:           "Discord Token",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "discordToken",
			NextStep:               "control_panel_channel",
		},

		"control_panel_channel": {
			ID:                     "control_panel_channel",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Discord Channel Setup (1/6)",
			StepMessage:            "Enter Discord Control Panel Channel ID",
			PrimaryPlaceholderText: "Channel ID",
			PrimaryLabel:           "Control Panel Channel ID",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "controlPanelChannelID",
			NextStep:               "save_channel",
		},

		"save_channel": {
			ID:                     "save_channel",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Discord Channel Setup (2/6)",
			StepMessage:            "Enter Discord Save Channel ID",
			PrimaryPlaceholderText: "Channel ID",
			PrimaryLabel:           "Save Channel ID",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "saveChannelID",
			NextStep:               "log_channel",
		},

		"log_channel": {
			ID:                     "log_channel",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Discord Channel Setup (3/6)",
			StepMessage:            "Enter Discord Log Channel ID",
			PrimaryPlaceholderText: "Channel ID",
			PrimaryLabel:           "Log Channel ID",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "logChannelID",
			NextStep:               "connection_list_channel",
		},

		"connection_list_channel": {
			ID:                     "connection_list_channel",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Discord Channel Setup (4/6)",
			StepMessage:            "Enter Discord Connection List Channel ID",
			PrimaryPlaceholderText: "Channel ID",
			PrimaryLabel:           "Connection List Channel ID",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "connectionListChannelID",
			NextStep:               "status_channel",
		},

		"status_channel": {
			ID:                     "status_channel",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Discord Channel Setup (5/6)",
			StepMessage:            "Enter Discord Status Channel ID",
			PrimaryPlaceholderText: "Channel ID",
			PrimaryLabel:           "Status Channel ID",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "statusChannelID",
			NextStep:               "control_channel",
		},

		"control_channel": {
			ID:                     "control_channel",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Discord Channel Setup (6/6)",
			StepMessage:            "Enter Discord Control Channel ID",
			PrimaryPlaceholderText: "Channel ID",
			PrimaryLabel:           "Control Channel ID",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "controlChannelID",
			NextStep:               "network_config_choice",
		},

		"network_config_choice": {
			ID:                     "network_config_choice",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Network Configuration",
			StepMessage:            "Do you want to configure network settings? Enter 'yes' to configure or Skip to use defaults. Note: Network configuration is especially important on Linux servers.",
			PrimaryPlaceholderText: "yes",
			PrimaryLabel:           "Configure Network",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Continue",
			SkipButtonText:         "Skip (Use Defaults)",
			ConfigField:            "",          // No config field, just for branching
			NextStep:               "game_port", // Default next step if they choose to configure
			// The actual next step will be determined by JS based on the answer
		},

		"game_port": {
			ID:                     "game_port",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Network Setup (1/6)",
			StepMessage:            "Enter the port number for game connections",
			PrimaryPlaceholderText: "27016",
			PrimaryLabel:           "Game Port",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "GamePort",
			NextStep:               "update_port",
		},

		"update_port": {
			ID:                     "update_port",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Network Setup (2/6)",
			StepMessage:            "Enter the port number for update connections",
			PrimaryPlaceholderText: "27015",
			PrimaryLabel:           "Update Port",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "UpdatePort",
			NextStep:               "upnp_enabled",
		},

		"upnp_enabled": {
			ID:                     "upnp_enabled",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Network Setup (3/6)",
			StepMessage:            "Enable UPnP? Enter 'yes' to enable or 'no' to disable.",
			PrimaryPlaceholderText: "yes/no",
			PrimaryLabel:           "Enable UPnP",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "UPNPEnabled", // We'll handle the boolean conversion in JS
			NextStep:               "server_visible",
		},

		"server_visible": {
			ID:                     "server_visible",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Network Setup (4/6)",
			StepMessage:            "Make server visible in the Server list? Enter 'yes' to make visible or 'no' to hide.",
			PrimaryPlaceholderText: "yes/no",
			PrimaryLabel:           "Server Visible",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "ServerVisible", // We'll handle the boolean conversion in JS
			NextStep:               "use_steam_p2p",
		},

		"use_steam_p2p": {
			ID:                     "use_steam_p2p",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Network Setup (5/6)",
			StepMessage:            "Use Steam P2P networking? Enter 'yes' to enable or 'no' to disable.",
			PrimaryPlaceholderText: "yes/no",
			PrimaryLabel:           "Use Steam P2P",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "UseSteamP2P", // We'll handle the boolean conversion in JS
			NextStep:               "local_ip_address",
		},

		"local_ip_address": {
			ID:                     "local_ip_address",
			Title:                  "Stationeers Server UI",
			HeaderTitle:            "Network Setup (6/6)",
			StepMessage:            "Enter server's local IP address in format 0.0.0.0 (no CIDR notation)",
			PrimaryPlaceholderText: "0.0.0.0",
			PrimaryLabel:           "Local IP Address",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Save & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "LocalIpAddress",
			NextStep:               "admin_account", // Continue to admin account setup after network config
		},
		"admin_account": {
			ID:                       "admin_account",
			Title:                    "Stationeers Server UI",
			HeaderTitle:              "Admin Account Setup",
			StepMessage:              "Set up your admin account.",
			PrimaryPlaceholderText:   "Username",
			PrimaryLabel:             "Username",
			SecondaryLabel:           "Password",
			SecondaryPlaceholderText: "Password",
			SecondaryLabelType:       "password",
			SubmitButtonText:         "Save & Continue",
			SkipButtonText:           "Skip Authentication",
			ConfigField:              "", // Special handling for admin account
			NextStep:                 "sscm_opt_in",
		},
		"sscm_opt_in": {
			ID:                     "sscm_opt_in",
			Title:                  "Stationeers Server Command Manager",
			HeaderTitle:            "Unique Feature",
			StepMessage:            "SSCM is a custom server plugin that allows you to execute commands directly from SSUI. It doesn't affect vanilla server functionality while giving you the ability to run commands from the Web console.",
			PrimaryPlaceholderText: "yes",
			PrimaryLabel:           "Enable SSCM",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Enable & Continue",
			SkipButtonText:         "Skip",
			ConfigField:            "IsSSCMEnabled",
			NextStep:               "finalize",
		},
		"finalize": {
			ID:                 "finalize",
			Title:              "Finalize Setup",
			HeaderTitle:        "",
			StepMessage:        "Ready to finalize? Your configuration has already been saved while you completed this setup. If you want to change any of the settings, you may click Return to Start and skip whatever you want to keep. Most options can also be changed on the config Tab in the UI.",
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   "Return to Start",
			SkipButtonText:     "Skip Authentication",
			NextStep:           "welcome", // Return to first step if "Return to Setup" is clicked
		},
	}

	data := TemplateData{
		IsFirstTimeSetup: config.IsFirstTimeSetup,
		Path:             path,
		Step:             stepID,
		FooterText:       "Need help? Check the Stationeers Server UI Github Wiki.",
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
			if stepID == "sscm_opt_in" {
				data.FooterText = "Opt in to SSCM for the most powerful Stationeers server management! This license protects this unique feature, ensuring it stays exclusive to SSUI users. Check the terms in the SSUI GitHub Wiki. Don’t be worried, the license simply protects SSCM’s integrity and its integration with SSUI."
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
		data.Title = "Stationeers Server UI"
		data.HeaderTitle = "Manage Users"
		data.PrimaryLabel = "Username to Add/Update"
		data.SecondaryLabel = "New Password"
		data.SecondaryPlaceholderText = "Password"
		data.SecondaryLabelType = "password"
		data.SubmitButtonText = "Add/Update User"
		data.Mode = "changeuser"
		data.ShowExtraButtons = false

	default:
		data.Title = "Stationeers Server UI"
		data.HeaderTitle = ""
		data.PrimaryLabel = "Username"
		data.SecondaryLabel = "Password"
		data.PrimaryPlaceholderText = "Enter Username"
		data.SecondaryPlaceholderText = "Enter Password"
		data.SecondaryLabelType = "password"
		data.SubmitButtonText = "Login"
		data.Mode = "login"
		data.ShowExtraButtons = false
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		logger.Web.Error("Failed to execute 2BoxForm template: %v" + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
