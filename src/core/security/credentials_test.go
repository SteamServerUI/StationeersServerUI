package security

import (
	"testing"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
)

func TestPrincipalUsesCurrentGroupPermissions(t *testing.T) {
	identity := config.NewIdentityConfig()
	identity.Groups["viewer"] = config.IdentityGroup{Permissions: []string{PermissionServerView}}
	user := config.IdentityUser{ID: "user-1", Username: "Human", GroupIDs: []string{"viewer"}}

	principal := principalForUser(user, identity, "session-1", "session")
	if !principal.Permissions[PermissionServerView] {
		t.Fatal("principal did not get current group permission")
	}
	if principal.Permissions[PermissionServerControl] {
		t.Fatal("principal got a permission the group does not have")
	}
}

func TestSessionCSRFValidation(t *testing.T) {
	session := config.IdentitySession{CSRFHash: hashSecret("some csrf secret")}
	if !ValidateSessionCSRF(session, "some csrf secret") {
		t.Fatal("correct CSRF value did not verify")
	}
	if ValidateSessionCSRF(session, "wrong") {
		t.Fatal("wrong CSRF value verified")
	}
}

func TestExpiredTokenShape(t *testing.T) {
	now := time.Now()
	token := config.IdentityToken{ExpiresAt: &now}
	if token.ExpiresAt.After(now) {
		t.Fatal("test token should already be expired")
	}
}
