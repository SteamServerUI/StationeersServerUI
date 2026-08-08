package middleware

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
)

func TestIdentityMiddlewareRequiresAuthentication(t *testing.T) {
	withTestIdentity(t, func() {
		handler := IdentityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		request := httptest.NewRequest(http.MethodGet, "https://ssui.test/api", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", response.Code)
		}
	})
}

func TestSessionRequestNeedsCSRFWhenChangingState(t *testing.T) {
	withTestIdentity(t, func() {
		credential := createTestSession(t)
		handler := IdentityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		request := httptest.NewRequest(http.MethodPost, "https://ssui.test/api", nil)
		request.AddCookie(&http.Cookie{Name: security.SessionCookieName, Value: credential.Value})
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusForbidden {
			t.Fatalf("expected missing CSRF to return 403, got %d", response.Code)
		}

		request = httptest.NewRequest(http.MethodPost, "https://ssui.test/api", nil)
		request.AddCookie(&http.Cookie{Name: security.SessionCookieName, Value: credential.Value})
		request.Header.Set("X-SSUI-CSRF", credential.CSRF)
		response = httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusNoContent {
			t.Fatalf("expected valid CSRF to pass, got %d", response.Code)
		}
	})
}

func TestPermissionMiddlewareDeniesMissingGrant(t *testing.T) {
	withTestIdentity(t, func() {
		credential := createTestSession(t)
		handler := IdentityMiddleware(RequirePermission(security.PermissionUsersManage, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		request := httptest.NewRequest(http.MethodGet, "https://ssui.test/api", nil)
		request.AddCookie(&http.Cookie{Name: security.SessionCookieName, Value: credential.Value})
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)
		if response.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", response.Code)
		}
	})
}

func withTestIdentity(t *testing.T, test func()) {
	t.Helper()
	originalPath := config.ConfigPath
	config.ConfigPath = filepath.Join(t.TempDir(), "config.json")
	t.Cleanup(func() { config.ConfigPath = originalPath })

	now := time.Now()
	identity := config.NewIdentityConfig()
	identity.SetupRequired = false
	identity.Groups["viewer"] = config.IdentityGroup{
		ID:          "viewer",
		Name:        "Viewer",
		Permissions: []string{security.PermissionServerView},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	identity.Users["user-1"] = config.IdentityUser{
		ID:         "user-1",
		Username:   "viewer",
		Normalized: "viewer",
		Enabled:    true,
		GroupIDs:   []string{"viewer"},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := config.SetIdentityConfig(identity); err != nil {
		t.Fatal(err)
	}
	test()
}

func createTestSession(t *testing.T) security.SessionCredential {
	t.Helper()
	_, credential, err := security.CreateSession("user-1", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return credential
}
