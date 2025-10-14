// auth.go
package security

//repurposed from a Jacksonthemaster private repo

import (
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserCredentials for login JSON
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GenerateJWT creates a JWT for a given username
func GenerateJWT(username string, apikeyduration ...int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.GetAuthTokenLifetime()) * time.Minute)
	if strings.HasPrefix(username, "apikey-") {
		durationMonths := 1
		if len(apikeyduration) > 0 {
			durationMonths = apikeyduration[0]
		}
		expirationTime = time.Now().AddDate(0, durationMonths, 0)
	}
	claims := &jwt.MapClaims{
		"exp": expirationTime.Unix(),
		"iss": "StationeersServerUI",
		"id":  username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetJwtKey()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateCredentials checks username and password against stored users
func ValidateCredentials(creds UserCredentials) (bool, error) {
	storedHash, exists := config.GetUsers()[creds.Username]
	if !exists {
		return false, nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password))
	return err == nil, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ValidateJWT checks if a JWT token is valid
func ValidateJWT(tokenString string) (bool, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJwtKey()), nil
	})
	if err != nil || !token.Valid {
		return false, err
	}
	return true, nil
}
