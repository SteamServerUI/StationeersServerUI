package backupsv2

import (
	"github.com/fsnotify/fsnotify"
)

// fsWatcher wraps fsnotify.Watcher with additional safety
type fsWatcher struct {
	watcher *fsnotify.Watcher
	events  chan fsnotify.Event
	errors  chan error
}

// newFsWatcher creates a new file system watcher
func newFsWatcher(path string) (*fsWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watcher.Add(path); err != nil {
		watcher.Close()
		return nil, err
	}

	return &fsWatcher{
		watcher: watcher,
		events:  make(chan fsnotify.Event),
		errors:  make(chan error),
	}, nil
}

// close stops the watcher
func (w *fsWatcher) close() {
	if w.watcher != nil {
		w.watcher.Close()
	}
}
