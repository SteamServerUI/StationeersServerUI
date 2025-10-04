// handlers.go
package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config/configchanger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/loader"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/security"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/google/uuid"
)

var setupReminderCount = 0 // to limit the number of setup reminders shown to the user

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
		Expires:  time.Now().Add(time.Duration(config.GetAuthTokenLifetime()) * time.Minute),
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
		// Log request details for debugging
		//logger.Web.Debug("Request Path:" + r.URL.Path) //very spammy

		// Check for first-time setup redirect
		if config.GetIsFirstTimeSetup() {
			totalSetupReminderCount := 3 // Defines how often we redirect the users reqests to the setup page
			if setupReminderCount < totalSetupReminderCount {
				if r.URL.Path == "/" && (r.Referer() == "" || r.Referer() != "/setup") {
					remainingReminderCount := totalSetupReminderCount - setupReminderCount
					logger.Web.Warn("🔍Redirecting to setup page, you should really enable authentication...")
					logger.Web.Warn(fmt.Sprintf("You will be remined %s more times.", strconv.Itoa(remainingReminderCount)))
					http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
					setupReminderCount++
					return
				}
			}
		}

		if !config.GetAuthEnabled() {
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

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var creds security.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Invalid JSON"})
		return
	}

	if strings.HasPrefix(creds.Username, "apikey-") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Invalid Username"})
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

	// Initialize Users map if nil
	if config.GetUsers() == nil {
		config.SetUsers(make(map[string]string))
	}

	// Add or update the user
	config.SetUsers(map[string]string{creds.Username: hashedPassword})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User registered successfully",
		"username": creds.Username,
	})
}

// SetupFinalizeHandler marks setup as complete
func SetupFinalizeHandler(w http.ResponseWriter, r *http.Request) {

	//check if users map is nil or empty
	if len(config.GetUsers()) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No users registered - cannot finalize setup at this time. You should really enable authentication - or click 'Skip authentication'"})
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
	config.SetIsFirstTimeSetup(false)
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

	logger.Web.Info("User Setup finalized successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Setup finalized successfully",
		"restart_hint": "You will be redirected to the login page...",
	})
	loader.ReloadBackend()
}

func RegisterAPIKeyHandler(w http.ResponseWriter, r *http.Request) {

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Allow only GET or POST methods
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method Not Allowed"})
		return
	}

	// Set default duration for GET requests, require duration for POST
	durationMonths := 1
	if r.Method == http.MethodPost {
		var reqBody struct {
			DurationMonths *int `json:"durationMonths"` // Use pointer to distinguish between 0 and unspecified
		}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Invalid JSON"})
			return
		}
		if reqBody.DurationMonths == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - durationMonths is required for POST"})
			return
		}
		if *reqBody.DurationMonths <= 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Duration must be positive"})
			return
		}
		durationMonths = *reqBody.DurationMonths
	}

	var creds security.UserCredentials

	// Generate a random UUID as the username
	creds.Username = "apikey-" + uuid.NewString()

	// Hash a random UUID as the password
	hashedPassword, err := security.HashPassword(uuid.NewString())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Initialize Users map if nil
	if config.GetUsers() == nil {
		config.SetUsers(make(map[string]string))
	}

	// Add or update the user
	config.SetUsers(map[string]string{creds.Username: hashedPassword})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	apikey, err := security.GenerateJWT(creds.Username, durationMonths)
	expires := time.Now().AddDate(0, durationMonths, 0)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "APIKey registered successfully",
		"apikey":  apikey,
		"expires": expires.Format(time.RFC3339),
	})
	logger.Security.Infof("APIKey %s registered successfully. Expires: %s ", creds.Username, expires.Format(time.RFC3339))
}
