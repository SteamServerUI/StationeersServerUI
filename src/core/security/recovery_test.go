package security

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
)

func TestRecoverOwnerRevokesCredentials(t *testing.T) {
	originalPath := config.ConfigPath
	config.ConfigPath = filepath.Join(t.TempDir(), "config.json")
	t.Cleanup(func() { config.ConfigPath = originalPath })

	identity := config.NewIdentityConfig()
	identity.SetupRequired = false
	identity.Sessions["old-session"] = config.IdentitySession{ID: "old-session"}
	identity.Tokens["old-token"] = config.IdentityToken{ID: "old-token"}
	if err := config.SetIdentityConfig(identity); err != nil {
		t.Fatal(err)
	}

	owner, err := RecoverOwner("recovery", "a safe recovery password", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if !owner.Enabled || !stringSliceContains(owner.GroupIDs, OwnerGroupID) {
		t.Fatal("recovered user is not an enabled owner")
	}
	result := config.GetIdentityConfig()
	if len(result.Sessions) != 0 || len(result.Tokens) != 0 {
		t.Fatal("owner recovery did not revoke existing credentials")
	}
}
