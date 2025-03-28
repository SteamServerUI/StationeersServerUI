// auth.go
package tlsconfig

//repurposed from a Jacksonthemaster private repo

import (
	"StationeersServerUI/src/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserCredentials for login JSON
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	claims := &jwt.MapClaims{ // Using MapClaims for simplicity
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
		Secure:   true, // Works with your TLS
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	// Return token (optional, UI can ignore this if it uses the cookie)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// AuthMiddleware protects routes with cookie-based JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("AuthToken")
		if err != nil {
			if err == http.ErrNoCookie {
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
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		// Token valid, proceed
		next.ServeHTTP(w, r)
	})
}

// BasicLoginPage serves a simple login form (temporary)
func BasicLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle form submission by redirecting to LoginHandler
		username := r.FormValue("username")
		password := r.FormValue("password")
		creds := UserCredentials{Username: username, Password: password}
		data, _ := json.Marshal(creds)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		LoginHandler(w, req)
		return
	}

	// Serve basic HTML form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html>
        <head><title>Login</title></head>
        <body>
            <h2>Stationeers Server UI - Login</h2>
            <form method="POST" action="/login">
                <label>Username:</label><input type="text" name="username"><br>
                <label>Password:</label><input type="password" name="password"><br>
                <input type="submit" value="Login">
            </form>
        </body>
        </html>
    `)
}
