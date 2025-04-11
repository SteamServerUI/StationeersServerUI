package web

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/logger"
	"net/http"
	"text/template"
)

func ServeLoginTemplate(w http.ResponseWriter, r *http.Request) {
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
	step := r.URL.Query().Get("step")
	if step == "" && config.IsFirstTimeSetup {
		step = "welcome" // Start with welcome page for first-time setup
	}

	data := TemplateData{
		IsFirstTimeSetup: config.IsFirstTimeSetup,
		Path:             path,
		Step:             step,
		FooterText:       "Need help? Check the 'Users System' page on the Github Wiki.",
	}

	switch {
	case config.IsFirstTimeSetup && (path == "/setup" || path == "/"):
		data.Mode = "setup"
		data.ShowExtraButtons = true
		switch step {
		case "welcome":
			data.Title = "Stationeers Server UI"
			data.HeaderTitle = ""
			data.StepMessage = ""
			data.UsernameLabel = ""
			data.PasswordLabel = ""
			data.PasswordType = "hidden"
			data.SubmitButtonText = "Start Setup"
			data.SkipButtonText = "Skip Setup"
		case "1":
			data.Title = "Stationeers Server UI"
			data.HeaderTitle = "Setup (1/5)"
			data.StepMessage = "Give your server a name like 'Space Station 13'"
			data.UsernamePlaceholder = "My Stationeers Server with UI"
			data.UsernameLabel = "Server Name"
			data.PasswordLabel = ""
			data.PasswordType = "hidden"
			data.SubmitButtonText = "Save & Continue"
			data.SkipButtonText = "Skip"
		case "2":
			data.Title = "Stationeers Server UI"
			data.HeaderTitle = "Setup (2/5)"
			data.StepMessage = "Set a save identifier like 'SpaceStation13 Moon'. Capitalitze the first letter of each word. Possible World types can be found in the Stationeers Wiki or the Stationeers Server UI GitHub Wiki. "
			data.UsernamePlaceholder = "Requires a SaveName, accepts optional WorldType"
			data.UsernameLabel = "Save Identifier"
			data.PasswordLabel = ""
			data.PasswordType = "hidden"
			data.SubmitButtonText = "Save & Continue"
			data.SkipButtonText = "Skip"
		case "3":
			data.Title = "Stationeers Server UI"
			data.HeaderTitle = "Setup (3/5)"
			data.StepMessage = "Choose the maximum number of players that can connect to the server."
			data.UsernamePlaceholder = "8"
			data.UsernameLabel = "Max Players"
			data.PasswordLabel = ""
			data.PasswordType = "hidden"
			data.SubmitButtonText = "Save & Continue"
			data.SkipButtonText = "Skip"
		case "4":
			data.Title = "Stationeers Server UI"
			data.HeaderTitle = "Setup (4/5)"
			data.StepMessage = "Set a gameserver password or skip this step."
			data.UsernamePlaceholder = "Server Password"
			data.UsernameLabel = "Server Password"
			data.PasswordLabel = ""
			data.PasswordType = "hidden"
			data.SubmitButtonText = "Save & Continue"
			data.SkipButtonText = "Skip"
		case "5":
			data.Title = "Stationeers Server UI"
			data.HeaderTitle = "Setup (5/5)"
			data.StepMessage = "Set up your admin account."
			data.UsernamePlaceholder = "Username"
			data.UsernameLabel = "Username"
			data.PasswordLabel = "Password"
			data.PasswordPlaceholder = "Password"
			data.PasswordType = "password"
			data.SubmitButtonText = "Save & Continue"
			data.SkipButtonText = "Skip Authentication"
		case "finalize":
			data.Title = "Stationeers Server UI - Finalize Setup"
			data.HeaderTitle = ""
			data.StepMessage = "Ready to finalize?"
			data.UsernameLabel = ""
			data.PasswordLabel = ""
			data.PasswordType = "hidden"
			data.SubmitButtonText = "Return to Setup"
			data.SkipButtonText = "Skip Authentication"
		default:
			// Redirect to welcome page if step is invalid
			data.Title = "Stationeers Server UI"
			data.HeaderTitle = ""
			data.StepMessage = ""
			data.UsernameLabel = ""
			data.PasswordLabel = ""
			data.PasswordType = "hidden"
			data.SubmitButtonText = "Start Setup"
			data.Step = "welcome"
		}
	case path == "/changeuser":
		data.Title = "Stationeers Server UI - Manage Users"
		data.HeaderTitle = "Manage Users"
		data.UsernameLabel = "Username to Add/Update"
		data.PasswordLabel = "New Password"
		data.PasswordPlaceholder = "Password"
		data.PasswordType = "password"
		data.SubmitButtonText = "Add/Update User"
		data.Mode = "changeuser"
		data.ShowExtraButtons = false
	default:
		data.Title = "Stationeers Server UI - Login"
		data.HeaderTitle = "Login"
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
