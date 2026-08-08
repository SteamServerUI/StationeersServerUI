package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/google/uuid"
)

const (
	OwnerGroupID        = "system-owner"
	setupSecretLifetime = 30 * time.Minute
)

type Principal struct {
	UserID       string
	Username     string
	CredentialID string
	Credential   string
	Permissions  map[string]bool
}

func NormalizeUsername(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func ValidateUsername(value string) error {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < 3 || len(trimmed) > 64 {
		return errors.New("username must be between 3 and 64 characters")
	}
	for _, character := range trimmed {
		if unicode.IsLetter(character) || unicode.IsDigit(character) || character == '-' || character == '_' || character == '.' {
			continue
		}
		return errors.New("username contains unsupported characters")
	}
	return nil
}

func PrepareIdentitySetup(now time.Time) (string, error) {
	identity := config.GetIdentityConfig()
	if identity.SchemaVersion == config.IdentitySchemaVersion && len(identity.Users) > 0 && !identity.SetupRequired {
		return "", nil
	}

	secret, err := randomSecret(32)
	if err != nil {
		return "", err
	}
	identity = config.NewIdentityConfig()
	identity.SetupSecretHash = hashSecret(secret)
	identity.SetupExpiresAt = now.Add(setupSecretLifetime)
	if err := config.SetIdentityConfig(identity); err != nil {
		return "", err
	}
	return secret, nil
}

func BootstrapOwner(setupSecret, username, password string, now time.Time) (config.IdentityUser, error) {
	if err := ValidateUsername(username); err != nil {
		return config.IdentityUser{}, err
	}
	passwordHash, err := HashIdentityPassword(password)
	if err != nil {
		return config.IdentityUser{}, err
	}

	var owner config.IdentityUser
	err = config.MutateIdentityConfig(func(identity *config.IdentityConfig) error {
		if !identity.SetupRequired || len(identity.Users) != 0 {
			return errors.New("owner setup is already complete")
		}
		if now.After(identity.SetupExpiresAt) || !secretMatches(identity.SetupSecretHash, setupSecret) {
			return errors.New("invalid or expired setup secret")
		}

		group := config.IdentityGroup{
			ID:          OwnerGroupID,
			Name:        "Owner",
			Description: "Full control of this SSUI backend",
			System:      true,
			Permissions: append([]string(nil), AllPermissions...),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		owner = config.IdentityUser{
			ID:           uuid.NewString(),
			Username:     strings.TrimSpace(username),
			Normalized:   NormalizeUsername(username),
			PasswordHash: passwordHash,
			Enabled:      true,
			GroupIDs:     []string{OwnerGroupID},
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		identity.Groups[group.ID] = group
		identity.Users[owner.ID] = owner
		identity.SetupRequired = false
		identity.SetupSecretHash = ""
		identity.SetupExpiresAt = time.Time{}
		return nil
	})
	return owner, err
}

func AuthenticateUser(username, password string) (config.IdentityUser, error) {
	identity := config.GetIdentityConfig()
	normalized := NormalizeUsername(username)
	for _, user := range identity.Users {
		if user.Normalized != normalized || !user.Enabled {
			continue
		}
		if !VerifyIdentityPassword(user.PasswordHash, password) {
			break
		}
		return user, nil
	}
	return config.IdentityUser{}, errors.New("invalid credentials")
}

func ResolvePermissions(user config.IdentityUser, identity config.IdentityConfig) map[string]bool {
	permissions := make(map[string]bool)
	for _, groupID := range user.GroupIDs {
		group, ok := identity.Groups[groupID]
		if !ok {
			continue
		}
		for _, permission := range group.Permissions {
			if IsPermission(permission) {
				permissions[permission] = true
			}
		}
	}
	return permissions
}

func randomSecret(size int) (string, error) {
	value := make([]byte, size)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("generate secret: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func hashSecret(value string) string {
	digest := sha256.Sum256([]byte(value))
	return hex.EncodeToString(digest[:])
}

func secretMatches(expected, value string) bool {
	actual := hashSecret(value)
	left, leftErr := hex.DecodeString(expected)
	right, rightErr := hex.DecodeString(actual)
	if leftErr != nil || rightErr != nil || len(left) != len(right) {
		return false
	}
	return subtle.ConstantTimeCompare(left, right) == 1
}
