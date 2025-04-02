package backupsv2

import (
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
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watcher.Add(path); err != nil {
		watcher.Close()
		return nil, err
	}

	w := &fsWatcher{
		watcher: watcher,
		events:  make(chan fsnotify.Event),
		errors:  make(chan error),
		done:    make(chan struct{}),
	}

	// Start forwarding events and errors
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
