package supportlink

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

type fakeService struct {
	mu         sync.Mutex
	secretHash string
	claimHash  string
	served     bool
	responses  []response
	deleted    bool
	puts       int
	snapshot   map[string]any
}

func (service *fakeService) authorized(r *http.Request) bool {
	return hashHex(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")) == service.secretHash
}

func (service *fakeService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service.mu.Lock()
	defer service.mu.Unlock()
	switch {
	case r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/snapshot"):
		if !service.authorized(r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		json.NewDecoder(r.Body).Decode(&service.snapshot)
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodPut:
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		service.puts++
		service.secretHash, service.claimHash = body["secretHash"], body["claimHash"]
		json.NewEncoder(w).Encode(map[string]bool{"claimed": false})
	case !service.authorized(r):
		w.WriteHeader(http.StatusForbidden)
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/requests"):
		poll := pollResponse{Claimed: true, Requests: []request{}}
		if !service.served {
			service.served = true
			poll.Requests = []request{{ID: "r1", Command: "status"}, {ID: "r2", Command: "rm -rf /"}, {ID: "r3", Command: "logs"}}
		} else {
			service.mu.Unlock()
			time.Sleep(20 * time.Millisecond)
			service.mu.Lock()
		}
		json.NewEncoder(w).Encode(poll)
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/responses"):
		var reply response
		json.NewDecoder(r.Body).Decode(&reply)
		service.responses = append(service.responses, reply)
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodDelete:
		service.deleted = true
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func TestSupportLinkRoundTrip(t *testing.T) {
	service := &fakeService{}
	server := httptest.NewServer(service)
	defer server.Close()
	endpoint = server.URL

	opened, err := Start(10 * time.Minute)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if !strings.HasPrefix(opened.ClaimCode, "SSUI-") || len(normalizeClaimCode(opened.ClaimCode)) != 26 {
		t.Fatalf("unexpected claim code %q", opened.ClaimCode)
	}
	if _, err := Start(10 * time.Minute); err != ErrActive {
		t.Fatalf("second start: got %v, want ErrActive", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for {
		service.mu.Lock()
		count := len(service.responses)
		pushed := service.snapshot != nil
		service.mu.Unlock()
		if count == 3 && pushed {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("got %d responses and snapshot %v, want 3 and a snapshot", count, pushed)
		}
		time.Sleep(10 * time.Millisecond)
	}

	snapshot := Snapshot()
	if !snapshot.Claimed || snapshot.ClaimCode != "" {
		t.Fatalf("claimed session should hide its code: %+v", snapshot)
	}

	closed := Stop()
	if closed.Active {
		t.Fatal("session still active after Stop")
	}

	service.mu.Lock()
	defer service.mu.Unlock()
	if service.claimHash != hashHex(normalizeClaimCode(opened.ClaimCode)) {
		t.Fatal("claim hash does not match the shown code")
	}
	if !service.responses[0].OK || service.responses[0].RequestID != "r1" || !strings.Contains(service.responses[0].Text, "SSUI") {
		t.Fatalf("status response: %+v", service.responses[0])
	}
	if service.responses[1].OK || service.responses[1].Error == "" {
		t.Fatalf("unknown command must be refused: %+v", service.responses[1])
	}
	if service.responses[2].File == nil || service.responses[2].File.Name != "recent-logs.txt" || service.responses[2].Data == nil {
		t.Fatalf("logs response: %+v", service.responses[2])
	}
	for _, key := range []string{"diagnostics", "players", "backups", "mods", "routines", "integrations", "update", "logs"} {
		if _, ok := service.snapshot[key]; !ok {
			t.Fatalf("snapshot is missing %q", key)
		}
	}
	if !service.deleted {
		t.Fatal("session was not deleted on Stop")
	}
	if len(closed.Activity) != 3 {
		t.Fatalf("activity entries: %d", len(closed.Activity))
	}
}

func TestStartRejectsInvalidDuration(t *testing.T) {
	if _, err := Start(time.Minute); err != ErrInvalidDuration {
		t.Fatalf("got %v", err)
	}
	if _, err := Start(3 * time.Hour); err != ErrInvalidDuration {
		t.Fatalf("got %v", err)
	}
}

func TestNormalizeClaimCode(t *testing.T) {
	if got := normalizeClaimCode(" ssui-abcde-fghij "); got != "ABCDEFGHIJ" {
		t.Fatalf("got %q", got)
	}
}

type releasingService struct {
	mu         sync.Mutex
	secretHash string
	polls      int
	deleted    bool
}

func (service *releasingService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service.mu.Lock()
	defer service.mu.Unlock()
	switch {
	case r.Method == http.MethodPut && !strings.HasSuffix(r.URL.Path, "/snapshot"):
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		service.secretHash = body["secretHash"]
		json.NewEncoder(w).Encode(map[string]bool{"claimed": false})
	case r.Method == http.MethodGet:
		service.polls++
		json.NewEncoder(w).Encode(pollResponse{Claimed: service.polls == 1, Requests: []request{}})
	case r.Method == http.MethodDelete:
		service.deleted = true
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func TestSupportLinkClosesWhenReleased(t *testing.T) {
	service := &releasingService{}
	server := httptest.NewServer(service)
	defer server.Close()
	endpoint = server.URL

	if _, err := Start(10 * time.Minute); err != nil {
		t.Fatalf("start: %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for Snapshot().Active {
		if time.Now().After(deadline) {
			Stop()
			t.Fatal("connection stayed open after the support team released it")
		}
		time.Sleep(10 * time.Millisecond)
	}
	service.mu.Lock()
	defer service.mu.Unlock()
	if !service.deleted {
		t.Fatal("released session was not deleted")
	}
}

type tombstoneService struct {
	mu    sync.Mutex
	polls int
	puts  int
}

func (service *tombstoneService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service.mu.Lock()
	defer service.mu.Unlock()
	switch {
	case r.Method == http.MethodPut && !strings.HasSuffix(r.URL.Path, "/snapshot"):
		service.puts++
		if service.puts > 1 {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"claimed": false})
	case r.Method == http.MethodGet:
		service.polls++
		if service.polls > 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(pollResponse{Claimed: true, Requests: []request{}})
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func TestSupportLinkEndsQuietlyAfterTombstone(t *testing.T) {
	server := httptest.NewServer(&tombstoneService{})
	defer server.Close()
	endpoint = server.URL

	if _, err := Start(10 * time.Minute); err != nil {
		t.Fatalf("start: %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for Snapshot().Active {
		if time.Now().After(deadline) {
			Stop()
			t.Fatal("connection stayed open after the session was ended")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if status := Snapshot(); status.Error != "" {
		t.Fatalf("a normal end must not show an error: %q", status.Error)
	}
}
