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
	if stepID == "" && config.IsFirstTimeSetup {
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
			NextStep:            "admin_account",
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
			StepMessage:      "Ready to finalize?",
			UsernameLabel:    "",
			PasswordLabel:    "",
			PasswordType:     "hidden",
			SubmitButtonText: "Return to Setup",
			SkipButtonText:   "Skip Authentication",
			NextStep:         "server_name", // Return to first step if "Return to Setup" is clicked
		},
	}

	data := TemplateData{
		IsFirstTimeSetup: config.IsFirstTimeSetup,
		Path:             path,
		Step:             stepID,
		FooterText:       "Need help? Check the 'Users System' page on the Github Wiki.",
	}

	switch {
	case config.IsFirstTimeSetup && (path == "/setup" || path == "/"):
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
