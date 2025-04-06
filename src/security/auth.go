// auth.go
package security

//repurposed from a Jacksonthemaster private repo

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/logger"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserCredentials for login JSON
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginPage serves static login files from ./UIMod/login.html, ./UIMod/login.js, ./UIMod/login.css
func LoginPage(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		http.ServeFile(w, r, "./UIMod/login.html")
	case "/login/login.js":
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "./UIMod/login.js")
	case "/login/login.css":
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./UIMod/login.css")
	default:
		http.NotFound(w, r)
	}
}

// LoginHandler issues a JWT cookie
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Check against config creds
	if creds.Username != config.Username || creds.Password != config.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	expirationTime := time.Now().Add(time.Duration(config.AuthTokenLifetime) * time.Minute)
	claims := &jwt.MapClaims{
		"exp": expirationTime.Unix(),
		"iss": "StationeersServerUI",
		"id":  creds.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JwtKey))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "AuthToken",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// AuthMiddleware protects routes with cookie-based JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("AuthToken")
		if err != nil {
			if err == http.ErrNoCookie {
				// Check if it's a browser (accepts HTML)
				if acceptsHTML(r) {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
				http.Error(w, "Unauthorized - No token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtKey), nil
		})

		if err != nil || !token.Valid {
			if acceptsHTML(r) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			if config.IsDebugMode {
				logger.Security.Warn("Unauthorized Request - Invalid token")
			}
			return
		}

		// Token valid, proceed
		next.ServeHTTP(w, r)
	})
}

// Helper to detect browser requests
func acceptsHTML(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "text/html")
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

	// For browser requests, redirect to login page
	if acceptsHTML(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// For API requests, return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}
