package backupmgr

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"

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
	logger.Backup.Debug("Creating watcher for path: " + normalizedPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	logger.Backup.Debug("Watcher created successfully")

	// Watch the root save path and all subdirectories
	err = filepath.WalkDir(normalizedPath, func(subPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if err := watcher.Add(subPath); err != nil {
				logger.Backup.Error("Failed to add subdir to watcher: " + subPath + ": " + err.Error())
			} else {
				logger.Backup.Debug("Successfully watching subdir: " + subPath)
			}
		}
		return nil
	})
	if err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to add paths to watcher: %w", err)
	}

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
