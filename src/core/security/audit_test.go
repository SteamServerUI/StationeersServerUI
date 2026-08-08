package security

import (
	"testing"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
)

func TestAppendAuditCapsOldEvents(t *testing.T) {
	identity := config.NewIdentityConfig()
	for index := 0; index < maxAuditEvents+5; index++ {
		AppendAudit(&identity, "actor", "owner", "user.update", "user", "target", time.Unix(int64(index), 0))
	}
	if len(identity.Audit) != maxAuditEvents {
		t.Fatalf("audit size = %d, want %d", len(identity.Audit), maxAuditEvents)
	}
	if identity.Audit[0].CreatedAt != time.Unix(5, 0) {
		t.Fatal("oldest audit events were not discarded")
	}
}
