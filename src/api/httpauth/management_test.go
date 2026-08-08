package httpauth

import (
	"testing"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
)

func TestGroupAssignmentCannotExceedManagersAccess(t *testing.T) {
	identity := config.NewIdentityConfig()
	identity.Groups["admin"] = config.IdentityGroup{
		Name:        "Admin",
		Permissions: []string{security.PermissionUsersManage, security.PermissionSettingsManage},
	}
	principal := security.Principal{Permissions: map[string]bool{security.PermissionUsersManage: true}}

	if err := canAssignGroups(identity, []string{"admin"}, principal); err == nil {
		t.Fatal("manager assigned a group containing permissions they do not have")
	}
}

func TestLastEnabledOwnerCount(t *testing.T) {
	identity := config.NewIdentityConfig()
	identity.Users["owner"] = config.IdentityUser{Enabled: true, GroupIDs: []string{security.OwnerGroupID}}
	identity.Users["disabled-owner"] = config.IdentityUser{Enabled: false, GroupIDs: []string{security.OwnerGroupID}}
	identity.Users["operator"] = config.IdentityUser{Enabled: true, GroupIDs: []string{"operator"}}

	if count := enabledOwnerCount(identity); count != 1 {
		t.Fatalf("expected one enabled owner, got %d", count)
	}
}

func TestPermissionGrantIsAdditiveAndUnique(t *testing.T) {
	principal := security.Principal{Permissions: map[string]bool{
		security.PermissionServerView: true,
		security.PermissionLogsView:   true,
	}}
	permissions, err := validatePermissionGrant(principal, []string{
		security.PermissionServerView,
		security.PermissionServerView,
		security.PermissionLogsView,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(permissions) != 2 {
		t.Fatalf("expected two unique permissions, got %d", len(permissions))
	}
}
