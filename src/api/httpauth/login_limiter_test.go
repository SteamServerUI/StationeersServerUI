package httpauth

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoginFailuresBlockAndClear(t *testing.T) {
	request := httptest.NewRequest("POST", "/auth/login", nil)
	request.RemoteAddr = "192.0.2.10:1234"
	key := loginAttemptKey(request, " Owner ")
	clearLoginFailures(key)
	now := time.Now()
	for index := 0; index < loginFailureLimit; index++ {
		recordLoginFailure(key, now.Add(time.Duration(index)*time.Second))
	}
	if retry := loginRetryAfter(key, now.Add(5*time.Second)); retry <= 0 {
		t.Fatal("failed logins did not block the key")
	}
	clearLoginFailures(key)
	if retry := loginRetryAfter(key, now.Add(5*time.Second)); retry != 0 {
		t.Fatal("successful login did not clear failures")
	}
}

func TestLoginAttemptKeyNormalizesUsername(t *testing.T) {
	request := httptest.NewRequest("POST", "/auth/login", nil)
	request.RemoteAddr = "192.0.2.11:1234"
	if loginAttemptKey(request, "Owner") != loginAttemptKey(request, " owner ") {
		t.Fatal("username variants produced different login keys")
	}
}

func TestBlockedLoginReturnsRetryAfter(t *testing.T) {
	request := httptest.NewRequest("POST", "/auth/login", nil)
	request.RemoteAddr = "192.0.2.12:1234"
	credentials := identityCredentials{Username: "owner", Password: "wrong password"}
	key := loginAttemptKey(request, credentials.Username)
	clearLoginFailures(key)
	t.Cleanup(func() { clearLoginFailures(key) })
	now := time.Now()
	for index := 0; index < loginFailureLimit; index++ {
		recordLoginFailure(key, now)
	}
	response := httptest.NewRecorder()
	if _, ok := authenticateIdentityRequest(response, request, credentials); ok {
		t.Fatal("blocked login was authenticated")
	}
	if response.Code != 429 {
		t.Fatalf("blocked login status = %d", response.Code)
	}
	if response.Header().Get("Retry-After") == "" {
		t.Fatal("blocked login omitted Retry-After")
	}
}
