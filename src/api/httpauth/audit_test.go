package httpauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
)

func TestAuditHandlerReturnsNewestFirst(t *testing.T) {
	originalPath := config.ConfigPath
	config.ConfigPath = filepath.Join(t.TempDir(), "config.json")
	t.Cleanup(func() { config.ConfigPath = originalPath })

	identity := config.NewIdentityConfig()
	identity.Audit = []config.IdentityAuditEvent{
		{ID: "old", ActorName: "owner", Action: "user.create", TargetType: "user", CreatedAt: time.Unix(1, 0)},
		{ID: "new", ActorName: "owner", Action: "group.update", TargetType: "group", CreatedAt: time.Unix(2, 0)},
	}
	if err := config.SetIdentityConfig(identity); err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	AuditHandler(response, httptest.NewRequest(http.MethodGet, "/api/v3/auth/audit", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	var body struct {
		Events []config.IdentityAuditEvent `json:"events"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Events) != 2 || body.Events[0].ID != "new" {
		t.Fatal("audit events were not returned newest first")
	}
}
