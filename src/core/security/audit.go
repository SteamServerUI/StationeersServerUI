package security

import (
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/google/uuid"
)

const maxAuditEvents = 1000

func AppendAudit(identity *config.IdentityConfig, actorID, actorName, action, targetType, targetID string, now time.Time) {
	identity.Audit = append(identity.Audit, config.IdentityAuditEvent{
		ID:         uuid.NewString(),
		ActorID:    actorID,
		ActorName:  actorName,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		CreatedAt:  now,
	})
	if overflow := len(identity.Audit) - maxAuditEvents; overflow > 0 {
		identity.Audit = append([]config.IdentityAuditEvent(nil), identity.Audit[overflow:]...)
	}
}

func RecordAudit(actorID, actorName, action, targetType, targetID string, now time.Time) error {
	return config.MutateIdentityConfig(func(identity *config.IdentityConfig) error {
		AppendAudit(identity, actorID, actorName, action, targetType, targetID, now)
		return nil
	})
}
