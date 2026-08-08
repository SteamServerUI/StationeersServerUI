package activity

import (
	"path/filepath"
	"testing"
)

func TestRecordPersistsAndReloadsActivity(t *testing.T) {
	originalPath := pathForStore
	directory := t.TempDir()
	pathForStore = func() string { return filepath.Join(directory, "activity.json") }
	resetActivityStore()
	t.Cleanup(func() {
		pathForStore = originalPath
		resetActivityStore()
	})

	created := Record(Event{Message: "Server restarted", ActorName: "owner"})
	if created.ID == "" || created.CreatedAt.IsZero() {
		t.Fatalf("record did not receive identity and timestamp: %#v", created)
	}

	store.Lock()
	store.loaded = false
	store.events = nil
	store.Unlock()
	reloaded := List(10)
	if len(reloaded) != 1 || reloaded[0].ID != created.ID {
		t.Fatalf("persisted event was not reloaded: %#v", reloaded)
	}
}

func TestSubscribersReceiveNewActivity(t *testing.T) {
	originalPath := pathForStore
	directory := t.TempDir()
	pathForStore = func() string { return filepath.Join(directory, "activity.json") }
	resetActivityStore()
	t.Cleanup(func() {
		pathForStore = originalPath
		resetActivityStore()
	})

	events, cancel := Subscribe()
	defer cancel()
	written := Record(Event{Message: "Backup created"})
	received := <-events
	if received.ID != written.ID {
		t.Fatalf("subscriber received %q, want %q", received.ID, written.ID)
	}
}

func resetActivityStore() {
	store.Lock()
	defer store.Unlock()
	store.loaded = false
	store.events = nil
	store.subscribers = make(map[chan Event]struct{})
}
