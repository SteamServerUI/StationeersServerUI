package security

import (
	"testing"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
)

func TestResolvePermissionsIsAdditive(t *testing.T) {
	identity := config.NewIdentityConfig()
	identity.Groups["operator"] = config.IdentityGroup{Permissions: []string{PermissionServerView, PermissionServerControl}}
	identity.Groups["logs"] = config.IdentityGroup{Permissions: []string{PermissionLogsView, "made.up"}}
	user := config.IdentityUser{GroupIDs: []string{"operator", "logs", "missing"}}

	permissions := ResolvePermissions(user, identity)
	for _, expected := range []string{PermissionServerView, PermissionServerControl, PermissionLogsView} {
		if !permissions[expected] {
			t.Fatalf("permission %q was not resolved", expected)
		}
	}
	if permissions["made.up"] {
		t.Fatal("unknown permission was accepted")
	}
}
