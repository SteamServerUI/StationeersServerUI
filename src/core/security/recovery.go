package security

import (
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/google/uuid"
)

func RecoverOwner(username, password string, now time.Time) (config.IdentityUser, error) {
	if err := ValidateUsername(username); err != nil {
		return config.IdentityUser{}, err
	}
	hash, err := HashIdentityPassword(password)
	if err != nil {
		return config.IdentityUser{}, err
	}
	var recovered config.IdentityUser
	err = config.MutateIdentityConfig(func(identity *config.IdentityConfig) error {
		ownerGroup, ok := identity.Groups[OwnerGroupID]
		if !ok {
			ownerGroup = config.IdentityGroup{
				ID:          OwnerGroupID,
				Name:        "Owner",
				Description: "Full control of this SSUI backend",
				System:      true,
				Permissions: append([]string(nil), AllPermissions...),
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			identity.Groups[OwnerGroupID] = ownerGroup
		}

		normalized := NormalizeUsername(username)
		for id, user := range identity.Users {
			if user.Normalized != normalized {
				continue
			}
			user.Username = strings.TrimSpace(username)
			user.PasswordHash = hash
			user.Enabled = true
			if !stringSliceContains(user.GroupIDs, OwnerGroupID) {
				user.GroupIDs = append(user.GroupIDs, OwnerGroupID)
			}
			user.UpdatedAt = now
			identity.Users[id] = user
			recovered = user
			break
		}
		if recovered.ID == "" {
			recovered = config.IdentityUser{
				ID:           uuid.NewString(),
				Username:     strings.TrimSpace(username),
				Normalized:   normalized,
				PasswordHash: hash,
				Enabled:      true,
				GroupIDs:     []string{OwnerGroupID},
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			identity.Users[recovered.ID] = recovered
		}
		identity.Sessions = make(map[string]config.IdentitySession)
		identity.Tokens = make(map[string]config.IdentityToken)
		identity.SetupRequired = false
		identity.SetupSecretHash = ""
		identity.SetupExpiresAt = time.Time{}
		AppendAudit(identity, "", "local-cli", "owner.recover", "user", recovered.ID, now)
		return nil
	})
	return recovered, err
}

func stringSliceContains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
