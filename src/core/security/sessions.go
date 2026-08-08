package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/google/uuid"
)

const (
	SessionCookieName       = "SSUISession"
	SessionIdleLifetime     = 24 * time.Hour
	SessionAbsoluteLifetime = 30 * 24 * time.Hour
	sessionTouchInterval    = 5 * time.Minute
)

type SessionCredential struct {
	Value string
	CSRF  string
}

func CreateSession(userID string, now time.Time) (config.IdentitySession, SessionCredential, error) {
	secret, err := randomSecret(32)
	if err != nil {
		return config.IdentitySession{}, SessionCredential{}, err
	}
	csrf := deriveCSRF(secret)
	session := config.IdentitySession{
		ID:                uuid.NewString(),
		SecretHash:        hashSecret(secret),
		CSRFHash:          hashSecret(csrf),
		UserID:            userID,
		CreatedAt:         now,
		LastUsedAt:        now,
		IdleExpiresAt:     now.Add(SessionIdleLifetime),
		AbsoluteExpiresAt: now.Add(SessionAbsoluteLifetime),
	}
	err = config.MutateIdentityConfig(func(identity *config.IdentityConfig) error {
		user, ok := identity.Users[userID]
		if !ok || !user.Enabled {
			return errors.New("user is disabled or missing")
		}
		identity.Sessions[session.ID] = session
		return nil
	})
	if err != nil {
		return config.IdentitySession{}, SessionCredential{}, err
	}
	return session, SessionCredential{Value: session.ID + "." + secret, CSRF: csrf}, nil
}

func AuthenticateSession(value string, now time.Time) (Principal, config.IdentitySession, error) {
	id, secret, ok := strings.Cut(value, ".")
	if !ok || id == "" || secret == "" {
		return Principal{}, config.IdentitySession{}, errors.New("invalid session")
	}
	identity := config.GetIdentityConfig()
	session, ok := identity.Sessions[id]
	if !ok || !secretMatches(session.SecretHash, secret) {
		return Principal{}, config.IdentitySession{}, errors.New("invalid session")
	}
	if !session.AbsoluteExpiresAt.After(now) || !session.IdleExpiresAt.After(now) {
		_ = RevokeSession(id)
		return Principal{}, config.IdentitySession{}, errors.New("session expired")
	}
	user, ok := identity.Users[session.UserID]
	if !ok || !user.Enabled {
		_ = RevokeSession(id)
		return Principal{}, config.IdentitySession{}, errors.New("user is disabled or missing")
	}
	if now.Sub(session.LastUsedAt) >= sessionTouchInterval {
		session.LastUsedAt = now
		session.IdleExpiresAt = now.Add(SessionIdleLifetime)
		if session.IdleExpiresAt.After(session.AbsoluteExpiresAt) {
			session.IdleExpiresAt = session.AbsoluteExpiresAt
		}
		_ = config.MutateIdentityConfig(func(value *config.IdentityConfig) error {
			if _, exists := value.Sessions[id]; exists {
				value.Sessions[id] = session
			}
			return nil
		})
	}
	return principalForUser(user, identity, id, "session"), session, nil
}

func ValidateSessionCSRF(session config.IdentitySession, value string) bool {
	return value != "" && secretMatches(session.CSRFHash, value)
}

func CSRFForSessionCredential(value string) string {
	_, secret, ok := strings.Cut(value, ".")
	if !ok || secret == "" {
		return ""
	}
	return deriveCSRF(secret)
}

func deriveCSRF(sessionSecret string) string {
	mac := hmac.New(sha256.New, []byte(sessionSecret))
	_, _ = mac.Write([]byte("ssui-csrf-v1"))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func RevokeSession(id string) error {
	return config.MutateIdentityConfig(func(identity *config.IdentityConfig) error {
		delete(identity.Sessions, id)
		return nil
	})
}

func RevokeUserCredentials(userID string) error {
	return config.MutateIdentityConfig(func(identity *config.IdentityConfig) error {
		for id, session := range identity.Sessions {
			if session.UserID == userID {
				delete(identity.Sessions, id)
			}
		}
		for id, token := range identity.Tokens {
			if token.OwnerID == userID {
				delete(identity.Tokens, id)
			}
		}
		return nil
	})
}

func principalForUser(user config.IdentityUser, identity config.IdentityConfig, credentialID, credential string) Principal {
	return Principal{
		UserID:       user.ID,
		Username:     user.Username,
		CredentialID: credentialID,
		Credential:   credential,
		Permissions:  ResolvePermissions(user, identity),
	}
}
