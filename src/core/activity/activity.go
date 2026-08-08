package activity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/google/uuid"
)

const maxEvents = 500

var pathForStore = defaultStorePath

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Tone      string    `json:"tone"`
	Message   string    `json:"message"`
	Detail    string    `json:"detail,omitempty"`
	ActorID   string    `json:"actorId,omitempty"`
	ActorName string    `json:"actorName,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

var store = struct {
	sync.Mutex
	loaded      bool
	events      []Event
	subscribers map[chan Event]struct{}
}{subscribers: make(map[chan Event]struct{})}

func Record(event Event) Event {
	store.Lock()
	defer store.Unlock()
	loadLocked()
	if event.ID == "" {
		event.ID = uuid.NewString()
	}
	if event.Type == "" {
		event.Type = "operation"
	}
	if event.Tone == "" {
		event.Tone = "info"
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	store.events = append([]Event{event}, store.events...)
	if len(store.events) > maxEvents {
		store.events = store.events[:maxEvents]
	}
	persistLocked()
	for subscriber := range store.subscribers {
		select {
		case subscriber <- event:
		default:
		}
	}
	return event
}

func List(limit int) []Event {
	store.Lock()
	defer store.Unlock()
	loadLocked()
	if limit <= 0 || limit > len(store.events) {
		limit = len(store.events)
	}
	result := make([]Event, limit)
	copy(result, store.events[:limit])
	return result
}

func Subscribe() (<-chan Event, func()) {
	store.Lock()
	defer store.Unlock()
	channel := make(chan Event, 16)
	store.subscribers[channel] = struct{}{}
	return channel, func() {
		store.Lock()
		defer store.Unlock()
		delete(store.subscribers, channel)
		close(channel)
	}
}

func loadLocked() {
	if store.loaded {
		return
	}
	store.loaded = true
	data, err := os.ReadFile(pathForStore())
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &store.events)
	if len(store.events) > maxEvents {
		store.events = store.events[:maxEvents]
	}
}

func persistLocked() {
	path := pathForStore()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return
	}
	data, err := json.MarshalIndent(store.events, "", "  ")
	if err != nil {
		return
	}
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, data, 0600); err == nil {
		_ = os.Rename(temporary, path)
	}
}

func defaultStorePath() string {
	return filepath.Join(config.GetSSUIFolder(), "config", "activity-v3.json")
}
