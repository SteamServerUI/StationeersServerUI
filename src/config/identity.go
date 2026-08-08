package config

import "time"

const IdentitySchemaVersion = 1

type IdentityConfig struct {
	SchemaVersion   int                        `json:"schemaVersion"`
	SetupRequired   bool                       `json:"setupRequired"`
	SetupSecretHash string                     `json:"setupSecretHash,omitempty"`
	SetupExpiresAt  time.Time                  `json:"setupExpiresAt,omitempty"`
	Users           map[string]IdentityUser    `json:"users"`
	Groups          map[string]IdentityGroup   `json:"groups"`
	Sessions        map[string]IdentitySession `json:"sessions"`
	Tokens          map[string]IdentityToken   `json:"tokens"`
	Audit           []IdentityAuditEvent       `json:"audit"`
}

type IdentityUser struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Normalized   string    `json:"normalized"`
	PasswordHash string    `json:"passwordHash"`
	Enabled      bool      `json:"enabled"`
	GroupIDs     []string  `json:"groupIds"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type IdentityGroup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	System      bool      `json:"system"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type IdentitySession struct {
	ID                string    `json:"id"`
	SecretHash        string    `json:"secretHash"`
	CSRFHash          string    `json:"csrfHash"`
	UserID            string    `json:"userId"`
	CreatedAt         time.Time `json:"createdAt"`
	LastUsedAt        time.Time `json:"lastUsedAt"`
	IdleExpiresAt     time.Time `json:"idleExpiresAt"`
	AbsoluteExpiresAt time.Time `json:"absoluteExpiresAt"`
}

type IdentityToken struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	SecretHash string     `json:"secretHash"`
	OwnerID    string     `json:"ownerId"`
	Scopes     []string   `json:"scopes"`
	CreatedAt  time.Time  `json:"createdAt"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
	RevokedAt  *time.Time `json:"revokedAt,omitempty"`
}

type IdentityAuditEvent struct {
	ID         string    `json:"id"`
	ActorID    string    `json:"actorId,omitempty"`
	ActorName  string    `json:"actorName"`
	Action     string    `json:"action"`
	TargetType string    `json:"targetType"`
	TargetID   string    `json:"targetId,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

func NewIdentityConfig() IdentityConfig {
	return IdentityConfig{
		SchemaVersion: IdentitySchemaVersion,
		SetupRequired: true,
		Users:         make(map[string]IdentityUser),
		Groups:        make(map[string]IdentityGroup),
		Sessions:      make(map[string]IdentitySession),
		Tokens:        make(map[string]IdentityToken),
		Audit:         make([]IdentityAuditEvent, 0),
	}
}

func normalizeIdentityConfig(value IdentityConfig) IdentityConfig {
	if value.SchemaVersion == 0 {
		return NewIdentityConfig()
	}
	if value.Users == nil {
		value.Users = make(map[string]IdentityUser)
	}
	if value.Groups == nil {
		value.Groups = make(map[string]IdentityGroup)
	}
	if value.Sessions == nil {
		value.Sessions = make(map[string]IdentitySession)
	}
	if value.Tokens == nil {
		value.Tokens = make(map[string]IdentityToken)
	}
	if value.Audit == nil {
		value.Audit = make([]IdentityAuditEvent, 0)
	}
	return value
}

func GetIdentityConfig() IdentityConfig {
	ConfigMu.RLock()
	defer ConfigMu.RUnlock()
	return cloneIdentityConfig(Identity)
}

func SetIdentityConfig(value IdentityConfig) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()
	Identity = normalizeIdentityConfig(cloneIdentityConfig(value))
	return safeSaveConfigAtomic()
}

func MutateIdentityConfig(change func(*IdentityConfig) error) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	working := cloneIdentityConfig(Identity)
	if err := change(&working); err != nil {
		return err
	}
	Identity = normalizeIdentityConfig(working)
	return safeSaveConfigAtomic()
}

func cloneIdentityConfig(value IdentityConfig) IdentityConfig {
	clone := value
	clone.Users = make(map[string]IdentityUser, len(value.Users))
	for id, user := range value.Users {
		user.GroupIDs = append([]string(nil), user.GroupIDs...)
		clone.Users[id] = user
	}
	clone.Groups = make(map[string]IdentityGroup, len(value.Groups))
	for id, group := range value.Groups {
		group.Permissions = append([]string(nil), group.Permissions...)
		clone.Groups[id] = group
	}
	clone.Sessions = make(map[string]IdentitySession, len(value.Sessions))
	for id, session := range value.Sessions {
		clone.Sessions[id] = session
	}
	clone.Tokens = make(map[string]IdentityToken, len(value.Tokens))
	for id, token := range value.Tokens {
		token.Scopes = append([]string(nil), token.Scopes...)
		clone.Tokens[id] = token
	}
	clone.Audit = append([]IdentityAuditEvent(nil), value.Audit...)
	return clone
}
