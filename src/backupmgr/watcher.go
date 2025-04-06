package backupmgr

import (
	"StationeersServerUI/src/logger"
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// fsWatcher wraps fsnotify.Watcher with additional safety
type fsWatcher struct {
	watcher *fsnotify.Watcher
	events  chan fsnotify.Event
	errors  chan error
	done    chan struct{}
}

// newFsWatcher creates a new file system watcher
func newFsWatcher(path string) (*fsWatcher, error) {
	// Normalize path
	normalizedPath := filepath.Clean(path)
	fmt.Printf("Creating watcher for path: %s\n", normalizedPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	logger.Backup.Info("Watcher created successfully")

	if err := watcher.Add(normalizedPath); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to add path %s to watcher: %w", normalizedPath, err)
	}
	fmt.Printf("Successfully watching path: %s\n", normalizedPath)

	w := &fsWatcher{
		watcher: watcher,
		events:  make(chan fsnotify.Event),
		errors:  make(chan error),
		done:    make(chan struct{}),
	}

	go w.forwardEvents()
	return w, nil
}

// forwardEvents forwards events and errors from the underlying watcher
func (w *fsWatcher) forwardEvents() {

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				close(w.events)
				return
			}
			select {
			case w.events <- event:
				// Successfully sent event
			case <-w.done:
				return
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				close(w.errors)
				return
			}
			select {
			case w.errors <- err:
				// Successfully sent error
			case <-w.done:
				return
			}
		case <-w.done:
			return
		}
	}
}

// close stops the watcher and closes all channels
func (w *fsWatcher) close() {
	close(w.done)
	if w.watcher != nil {
		w.watcher.Close()
	}
}
