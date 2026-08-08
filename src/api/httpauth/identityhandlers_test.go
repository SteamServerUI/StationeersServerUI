package httpauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
)

func TestBootstrapOwnerCreatesBrowserSession(t *testing.T) {
	originalPath := config.ConfigPath
	config.ConfigPath = filepath.Join(t.TempDir(), "config.json")
	t.Cleanup(func() { config.ConfigPath = originalPath })
	if err := config.SetIdentityConfig(config.NewIdentityConfig()); err != nil {
		t.Fatal(err)
	}

	secret, err := security.PrepareIdentitySetup(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	body, _ := json.Marshal(map[string]string{
		"setupSecret": secret,
		"username":    "owner",
		"password":    "a long owner password",
	})
	request := httptest.NewRequest(http.MethodPost, "/api/v3/auth/setup/bootstrap", strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	BootstrapOwnerHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != security.SessionCookieName || !cookies[0].HttpOnly || !cookies[0].Secure || cookies[0].SameSite != http.SameSiteStrictMode {
		t.Fatal("bootstrap did not return a hardened session cookie")
	}
	identity := config.GetIdentityConfig()
	if identity.SetupRequired || identity.SetupSecretHash != "" || len(identity.Users) != 1 || len(identity.Audit) != 1 {
		t.Fatal("bootstrap did not complete and consume setup state")
	}
}

func TestSessionLoginRejectsWrongPassword(t *testing.T) {
	originalPath := config.ConfigPath
	config.ConfigPath = filepath.Join(t.TempDir(), "config.json")
	t.Cleanup(func() { config.ConfigPath = originalPath })
	if err := config.SetIdentityConfig(config.NewIdentityConfig()); err != nil {
		t.Fatal(err)
	}

	secret, err := security.PrepareIdentitySetup(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := security.BootstrapOwner(secret, "owner", "a long owner password", time.Now()); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v3/auth/login", strings.NewReader(`{"username":"owner","password":"wrong password"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	SessionLoginHandler(response, request)
	if response.Code != http.StatusUnauthorized || len(response.Result().Cookies()) != 0 {
		t.Fatal("wrong password created a browser session")
	}
}
