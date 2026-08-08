package security

import (
	"errors"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/google/uuid"
)

const tokenPrefix = "ssui_pat_"

func CreateNamedToken(ownerID, name string, scopes []string, expiresAt *time.Time, now time.Time) (config.IdentityToken, string, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 80 {
		return config.IdentityToken{}, "", errors.New("token name must be between 1 and 80 characters")
	}
	identity := config.GetIdentityConfig()
	owner, ok := identity.Users[ownerID]
	if !ok || !owner.Enabled {
		return config.IdentityToken{}, "", errors.New("token owner is disabled or missing")
	}
	permissions := ResolvePermissions(owner, identity)
	uniqueScopes := make([]string, 0, len(scopes))
	seen := make(map[string]bool)
	for _, scope := range scopes {
		if !IsPermission(scope) || !permissions[scope] {
			return config.IdentityToken{}, "", errors.New("token scope exceeds owner permissions")
		}
		if !seen[scope] {
			seen[scope] = true
			uniqueScopes = append(uniqueScopes, scope)
		}
	}
	if len(uniqueScopes) == 0 {
		return config.IdentityToken{}, "", errors.New("token requires at least one scope")
	}
	if expiresAt != nil && !expiresAt.After(now) {
		return config.IdentityToken{}, "", errors.New("token expiry must be in the future")
	}
	secret, err := randomSecret(32)
	if err != nil {
		return config.IdentityToken{}, "", err
	}
	token := config.IdentityToken{
		ID:         uuid.NewString(),
		Name:       name,
		SecretHash: hashSecret(secret),
		OwnerID:    ownerID,
		Scopes:     uniqueScopes,
		CreatedAt:  now,
		ExpiresAt:  expiresAt,
	}
	if err := config.MutateIdentityConfig(func(value *config.IdentityConfig) error {
		value.Tokens[token.ID] = token
		AppendAudit(value, owner.ID, owner.Username, "token.create", "token", token.ID, now)
		return nil
	}); err != nil {
		return config.IdentityToken{}, "", err
	}
	return token, tokenPrefix + token.ID + "_" + secret, nil
}

func AuthenticateToken(value string, now time.Time) (Principal, config.IdentityToken, error) {
	if !strings.HasPrefix(value, tokenPrefix) {
		return Principal{}, config.IdentityToken{}, errors.New("invalid token")
	}
	id, secret, ok := strings.Cut(strings.TrimPrefix(value, tokenPrefix), "_")
	if !ok || id == "" || secret == "" {
		return Principal{}, config.IdentityToken{}, errors.New("invalid token")
	}
	identity := config.GetIdentityConfig()
	token, ok := identity.Tokens[id]
	if !ok || token.RevokedAt != nil || !secretMatches(token.SecretHash, secret) {
		return Principal{}, config.IdentityToken{}, errors.New("invalid token")
	}
	if token.ExpiresAt != nil && !token.ExpiresAt.After(now) {
		return Principal{}, config.IdentityToken{}, errors.New("token expired")
	}
	user, ok := identity.Users[token.OwnerID]
	if !ok || !user.Enabled {
		return Principal{}, config.IdentityToken{}, errors.New("token owner is disabled or missing")
	}
	ownerPermissions := ResolvePermissions(user, identity)
	granted := make(map[string]bool)
	for _, scope := range token.Scopes {
		if ownerPermissions[scope] {
			granted[scope] = true
		}
	}
	principal := Principal{
		UserID:       user.ID,
		Username:     user.Username,
		CredentialID: token.ID,
		Credential:   "token",
		Permissions:  granted,
	}
	if token.LastUsedAt == nil || now.Sub(*token.LastUsedAt) >= sessionTouchInterval {
		token.LastUsedAt = &now
		_ = config.MutateIdentityConfig(func(value *config.IdentityConfig) error {
			if _, exists := value.Tokens[id]; exists {
				value.Tokens[id] = token
			}
			return nil
		})
	}
	return principal, token, nil
}

func RevokeToken(id string, now time.Time) error {
	return RevokeTokenAs(id, "", "", now)
}

func RevokeTokenAs(id, actorID, actorName string, now time.Time) error {
	return config.MutateIdentityConfig(func(identity *config.IdentityConfig) error {
		token, ok := identity.Tokens[id]
		if !ok {
			return errors.New("token not found")
		}
		token.RevokedAt = &now
		identity.Tokens[id] = token
		if actorName != "" {
			AppendAudit(identity, actorID, actorName, "token.revoke", "token", id, now)
		}
		return nil
	})
}
