// handlers.go
package web

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/configchanger"
	"StationeersServerUI/src/loader"
	"StationeersServerUI/src/logger"
	"StationeersServerUI/src/security"
	"encoding/json"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var setupReminderCount = 0 // to limit the number of setup reminders shown to the user

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

// LoginHandler issues a JWT cookie
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds security.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Invalid JSON"})
		return
	}

	// Check credentials using security package
	valid, err := security.ValidateCredentials(creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}
	if !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized - Invalid credentials"})
		return
	}

	// Generate JWT
	tokenString, err := security.GenerateJWT(creds.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "AuthToken",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Duration(config.AuthTokenLifetime) * time.Minute),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// AuthMiddleware protects routes with cookie-based JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !config.AuthEnabled {
			if config.IsFirstTimeSetup {
				if setupReminderCount < 1 {
					http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
					setupReminderCount++
				}
			}
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("AuthToken")
		if err != nil {
			// Browser redirect check
			accept := r.Header.Get("Accept")
			if accept != "" && strings.Contains(accept, "text/html") {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			// API response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized - No token"})
			return
		}

		valid, err := security.ValidateJWT(cookie.Value)
		if err != nil || !valid {
			// Browser redirect check
			accept := r.Header.Get("Accept")
			if accept != "" && strings.Contains(accept, "text/html") {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			// API response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized - Invalid token"})
			logger.Security.Warn("Unauthorized Request - Invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the cookie by setting it with an expired time
	http.SetCookie(w, &http.Cookie{
		Name:     "AuthToken",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set to past time to expire immediately
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	accept := r.Header.Get("Accept")
	if accept != "" && strings.Contains(accept, "text/html") {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	// For API requests, return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}

// RegisterUserHandler registers new users
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var creds security.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Invalid JSON"})
		return
	}

	// Hash the password
	hashedPassword, err := security.HashPassword(creds.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Load existing config to update it
	existingConfig, err := config.LoadConfig()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error - Failed to load config"})
		return
	}

	// Initialize Users map if nil
	if existingConfig.Users == nil {
		existingConfig.Users = make(map[string]string)
	}

	// Add or update the user
	existingConfig.Users[creds.Username] = hashedPassword

	// Persist the updated config
	if err := configchanger.SaveConfig(existingConfig); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error - Failed to save config"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User registered successfully",
		"username": creds.Username,
	})
}

// SetupFinalizeHandler marks setup as complete
func SetupFinalizeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if setup is already done
	if !config.IsFirstTimeSetup {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Setup already finalized"})
		return
	}

	//check if users map is nil or empty
	if config.Users == nil || len(config.Users) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - No users registered - cannot finalize setup at this time"})
		return
	}

	// Load existing config to update it
	newConfig, err := config.LoadConfig()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error - Failed to load config"})
		return
	}

	// Mark setup as complete and enable auth
	config.IsFirstTimeSetup = false
	isTrue := true
	newConfig.AuthEnabled = &isTrue // Set the pointer to true

	// Save the updated config
	err = configchanger.SaveConfig(newConfig)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error - Failed to save config"})
		return
	}

	logger.Web.Core("User Setup finalized successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Setup finalized successfully",
		"restart_hint": "You will be redirected to the login page...",
	})
	loader.ReloadConfig()
}
