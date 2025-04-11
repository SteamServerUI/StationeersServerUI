package web

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/logger"
	"net/http"
	"text/template"
)

func ServeLoginTemplate(w http.ResponseWriter, r *http.Request) {
	type Step struct {
		ID                  string
		Title               string
		HeaderTitle         string
		StepMessage         string
		UsernameLabel       string
		PasswordLabel       string
		PasswordType        string
		SubmitButtonText    string
		SkipButtonText      string
		UsernamePlaceholder string
		PasswordPlaceholder string
		ConfigField         string // Field name to save in config
		NextStep            string // ID of the next step
	}

	type TemplateData struct {
		IsFirstTimeSetup    bool
		Path                string
		Title               string
		HeaderTitle         string
		StepMessage         string
		UsernameLabel       string
		PasswordLabel       string
		PasswordType        string
		SubmitButtonText    string
		SkipButtonText      string
		Mode                string
		ShowExtraButtons    bool
		FooterText          string
		Step                string
		ConfigField         string
		NextStep            string
		UsernamePlaceholder string
		PasswordPlaceholder string
	}

	tmpl, err := template.ParseFiles("./UIMod/login/login.html")
	if err != nil {
		logger.Web.Error("Failed to parse login template: %v" + err.Error())
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
			ID:               "welcome",
			Title:            "Stationeers Server UI",
			HeaderTitle:      "",
			StepMessage:      "",
			UsernameLabel:    "",
			PasswordLabel:    "",
			PasswordType:     "hidden",
			SubmitButtonText: "Start Setup",
			SkipButtonText:   "Skip Setup",
			NextStep:         "server_name",
		},
		"server_name": {
			ID:                  "server_name",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Server Name Setup",
			StepMessage:         "Give your server a name like 'Space Station 13'",
			UsernamePlaceholder: "My Stationeers Server with UI",
			UsernameLabel:       "Server Name",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "ServerName",
			NextStep:            "save_identifier",
		},
		"save_identifier": {
			ID:                  "save_identifier",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Save Identifier Setup",
			StepMessage:         "Set a save identifier like 'SpaceStation13 Moon'. Capitalize the first letter of each word. Possible World types can be found in the Stationeers Wiki or the Stationeers Server UI GitHub Wiki.",
			UsernamePlaceholder: "Requires a SaveName, accepts optional WorldType",
			UsernameLabel:       "Save Identifier",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "SaveInfo",
			NextStep:            "max_players",
		},
		"max_players": {
			ID:                  "max_players",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Player Limit Setup",
			StepMessage:         "Choose the maximum number of players that can connect to the server.",
			UsernamePlaceholder: "8",
			UsernameLabel:       "Max Players",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "ServerMaxPlayers",
			NextStep:            "server_password",
		},
		"server_password": {
			ID:                  "server_password",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Server Password Setup",
			StepMessage:         "Set a gameserver password or skip this step.",
			UsernamePlaceholder: "Server Password",
			UsernameLabel:       "Server Password",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "ServerPassword",
			NextStep:            "game_branch",
		},
		"game_branch": {
			ID:                  "game_branch",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Game Branch Setup",
			StepMessage:         "Enter a beta branch or skip this to use the release version. If switching branches, make sure to restart SSUI after completing this wizzard.",
			UsernamePlaceholder: "beta",
			UsernameLabel:       "Game Branch",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Use Release Version",
			ConfigField:         "GameBranch",
			NextStep:            "discord_enabled",
		},

		"discord_enabled": {
			ID:                  "discord_enabled",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Discord Integration",
			StepMessage:         "Do you want to enable Discord integration? Enter 'yes' to enable or Skip to disable.",
			UsernamePlaceholder: "yes",
			UsernameLabel:       "Enable Discord",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip (Disable Discord)",
			ConfigField:         "IsDiscordEnabled", // We'll handle the boolean conversion in JS
			NextStep:            "discord_token",    // Default next step if enabled
			// The actual next step will be determined by JS based on the answer
		},

		"discord_token": {
			ID:                  "discord_token",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Discord Bot Token",
			StepMessage:         "Enter your Discord bot token for server integration",
			UsernamePlaceholder: "Discord Bot Token",
			UsernameLabel:       "Discord Token",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "DiscordToken",
			NextStep:            "control_panel_channel",
		},

		"control_panel_channel": {
			ID:                  "control_panel_channel",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Discord Channel Setup (1/6)",
			StepMessage:         "Enter Discord Control Panel Channel ID",
			UsernamePlaceholder: "Channel ID",
			UsernameLabel:       "Control Panel Channel ID",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "ControlPanelChannelID",
			NextStep:            "save_channel",
		},

		"save_channel": {
			ID:                  "save_channel",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Discord Channel Setup (2/6)",
			StepMessage:         "Enter Discord Save Channel ID",
			UsernamePlaceholder: "Channel ID",
			UsernameLabel:       "Save Channel ID",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "SaveChannelID",
			NextStep:            "log_channel",
		},

		"log_channel": {
			ID:                  "log_channel",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Discord Channel Setup (3/6)",
			StepMessage:         "Enter Discord Log Channel ID",
			UsernamePlaceholder: "Channel ID",
			UsernameLabel:       "Log Channel ID",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "LogChannelID",
			NextStep:            "connection_list_channel",
		},

		"connection_list_channel": {
			ID:                  "connection_list_channel",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Discord Channel Setup (4/6)",
			StepMessage:         "Enter Discord Connection List Channel ID",
			UsernamePlaceholder: "Channel ID",
			UsernameLabel:       "Connection List Channel ID",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "ConnectionListChannelID",
			NextStep:            "status_channel",
		},

		"status_channel": {
			ID:                  "status_channel",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Discord Channel Setup (5/6)",
			StepMessage:         "Enter Discord Status Channel ID",
			UsernamePlaceholder: "Channel ID",
			UsernameLabel:       "Status Channel ID",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "StatusChannelID",
			NextStep:            "control_channel",
		},

		"control_channel": {
			ID:                  "control_channel",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Discord Channel Setup (6/6)",
			StepMessage:         "Enter Discord Control Channel ID",
			UsernamePlaceholder: "Channel ID",
			UsernameLabel:       "Control Channel ID",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "ControlChannelID",
			NextStep:            "network_config_choice",
		},

		"network_config_choice": {
			ID:                  "network_config_choice",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Network Configuration",
			StepMessage:         "Do you want to configure network settings? Enter 'yes' to configure or Skip to use defaults. Note: Network configuration is especially important on Linux servers.",
			UsernamePlaceholder: "yes",
			UsernameLabel:       "Configure Network",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Continue",
			SkipButtonText:      "Skip (Use Defaults)",
			ConfigField:         "",          // No config field, just for branching
			NextStep:            "game_port", // Default next step if they choose to configure
			// The actual next step will be determined by JS based on the answer
		},

		"game_port": {
			ID:                  "game_port",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Network Setup (1/6)",
			StepMessage:         "Enter the port number for game connections",
			UsernamePlaceholder: "27016",
			UsernameLabel:       "Game Port",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "GamePort",
			NextStep:            "update_port",
		},

		"update_port": {
			ID:                  "update_port",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Network Setup (2/6)",
			StepMessage:         "Enter the port number for update connections",
			UsernamePlaceholder: "27017",
			UsernameLabel:       "Update Port",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "UpdatePort",
			NextStep:            "upnp_enabled",
		},

		"upnp_enabled": {
			ID:                  "upnp_enabled",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Network Setup (3/6)",
			StepMessage:         "Enable UPnP? Enter 'yes' to enable or 'no' to disable.",
			UsernamePlaceholder: "yes/no",
			UsernameLabel:       "Enable UPnP",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "UPNPEnabled", // We'll handle the boolean conversion in JS
			NextStep:            "server_visible",
		},

		"server_visible": {
			ID:                  "server_visible",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Network Setup (4/6)",
			StepMessage:         "Make server visible in the Server list? Enter 'yes' to make visible or 'no' to hide.",
			UsernamePlaceholder: "yes/no",
			UsernameLabel:       "Server Visible",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "ServerVisible", // We'll handle the boolean conversion in JS
			NextStep:            "use_steam_p2p",
		},

		"use_steam_p2p": {
			ID:                  "use_steam_p2p",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Network Setup (5/6)",
			StepMessage:         "Use Steam P2P networking? Enter 'yes' to enable or 'no' to disable.",
			UsernamePlaceholder: "yes/no",
			UsernameLabel:       "Use Steam P2P",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "UseSteamP2P", // We'll handle the boolean conversion in JS
			NextStep:            "local_ip_address",
		},

		"local_ip_address": {
			ID:                  "local_ip_address",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Network Setup (6/6)",
			StepMessage:         "Enter server's local IP address in format 0.0.0.0 (no CIDR notation)",
			UsernamePlaceholder: "0.0.0.0",
			UsernameLabel:       "Local IP Address",
			PasswordLabel:       "",
			PasswordType:        "hidden",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip",
			ConfigField:         "LocalIpAddress",
			NextStep:            "admin_account", // Continue to admin account setup after network config
		},
		"admin_account": {
			ID:                  "admin_account",
			Title:               "Stationeers Server UI",
			HeaderTitle:         "Admin Account Setup",
			StepMessage:         "Set up your admin account.",
			UsernamePlaceholder: "Username",
			UsernameLabel:       "Username",
			PasswordLabel:       "Password",
			PasswordPlaceholder: "Password",
			PasswordType:        "password",
			SubmitButtonText:    "Save & Continue",
			SkipButtonText:      "Skip Authentication",
			ConfigField:         "", // Special handling for admin account
			NextStep:            "finalize",
		},
		"finalize": {
			ID:               "finalize",
			Title:            "Stationeers Server UI - Finalize Setup",
			HeaderTitle:      "",
			StepMessage:      "Ready to finalize? Your configuration has already been saved while you completed this setup. If you want to change any of the settings, you may click Return to Start and skip whatever you want to keep. All options can also be changed on the config Tab in the UI.",
			UsernameLabel:    "",
			PasswordLabel:    "",
			PasswordType:     "hidden",
			SubmitButtonText: "Return to Start",
			SkipButtonText:   "Skip Authentication",
			NextStep:         "welcome", // Return to first step if "Return to Setup" is clicked
		},
	}

	data := TemplateData{
		IsFirstTimeSetup: config.IsFirstTimeSetup,
		Path:             path,
		Step:             stepID,
		FooterText:       "Need help? Check the Stationeers Server UI Github Wiki.",
	}

	switch {
	case path == "/setup":
		data.Mode = "setup"
		data.ShowExtraButtons = true

		// Find the current step in our map
		if step, exists := steps[stepID]; exists {
			data.Title = step.Title
			data.HeaderTitle = step.HeaderTitle
			data.StepMessage = step.StepMessage
			data.UsernameLabel = step.UsernameLabel
			data.PasswordLabel = step.PasswordLabel
			data.PasswordType = step.PasswordType
			data.SubmitButtonText = step.SubmitButtonText
			data.SkipButtonText = step.SkipButtonText
			data.ConfigField = step.ConfigField
			data.NextStep = step.NextStep
			data.UsernamePlaceholder = step.UsernamePlaceholder
			data.PasswordPlaceholder = step.PasswordPlaceholder
		} else {
			// Default to welcome page if step is invalid
			welcomeStep := steps["welcome"]
			data.Title = welcomeStep.Title
			data.HeaderTitle = welcomeStep.HeaderTitle
			data.StepMessage = welcomeStep.StepMessage
			data.UsernameLabel = welcomeStep.UsernameLabel
			data.PasswordLabel = welcomeStep.PasswordLabel
			data.PasswordType = welcomeStep.PasswordType
			data.SubmitButtonText = welcomeStep.SubmitButtonText
			data.SkipButtonText = welcomeStep.SkipButtonText
			data.ConfigField = welcomeStep.ConfigField
			data.NextStep = welcomeStep.NextStep
			data.Step = "welcome"
		}

	case path == "/changeuser":
		data.Title = "Stationeers Server UI"
		data.HeaderTitle = "Manage Users"
		data.UsernameLabel = "Username to Add/Update"
		data.PasswordLabel = "New Password"
		data.PasswordPlaceholder = "Password"
		data.PasswordType = "password"
		data.SubmitButtonText = "Add/Update User"
		data.Mode = "changeuser"
		data.ShowExtraButtons = false

	default:
		data.Title = "Stationeers Server UI"
		data.HeaderTitle = ""
		data.UsernameLabel = "Username"
		data.PasswordLabel = "Password"
		data.UsernamePlaceholder = "Enter Username"
		data.PasswordPlaceholder = "Enter Password"
		data.PasswordType = "password"
		data.SubmitButtonText = "Login"
		data.Mode = "login"
		data.ShowExtraButtons = false
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		logger.Web.Error("Failed to execute login template: %v" + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
