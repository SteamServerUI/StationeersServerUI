package logger

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
)

const (
	supportRecordingLimit = 64 << 20
	supportRecordingsKept = 3
)

var (
	ErrSupportRecordingActive = errors.New("a support recording is already running")
	ErrNoSupportRecording     = errors.New("no support recording is running")
)

type SupportRecording struct {
	Active    bool       `json:"active"`
	Full      bool       `json:"full"`
	Bytes     int64      `json:"bytes"`
	StartedAt time.Time  `json:"startedAt"`
	StoppedAt *time.Time `json:"stoppedAt,omitempty"`
	Dir       string     `json:"-"`
}

type supportRecorder struct {
	mu     sync.Mutex
	status SupportRecording
	ssui   *os.File
	game   *os.File
}

var (
	supportRecording     atomic.Pointer[supportRecorder]
	supportMu            sync.Mutex
	lastSupportRecording SupportRecording
)

func StartSupportRecording() (SupportRecording, error) {
	supportMu.Lock()
	defer supportMu.Unlock()
	if recorder := supportRecording.Load(); recorder != nil {
		return recorder.snapshot(), ErrSupportRecordingActive
	}

	root := filepath.Join(config.GetLogFolder(), "support-recordings")
	now := time.Now()
	dir := filepath.Join(root, now.Format("20060102-150405"))
	if err := os.MkdirAll(dir, 0700); err != nil {
		return SupportRecording{}, err
	}
	ssui, err := os.OpenFile(filepath.Join(dir, "ssui.log"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return SupportRecording{}, err
	}
	game, err := os.OpenFile(filepath.Join(dir, "gameserver.log"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		ssui.Close()
		return SupportRecording{}, err
	}
	pruneSupportRecordings(root)

	recorder := &supportRecorder{status: SupportRecording{Active: true, StartedAt: now, Dir: dir}, ssui: ssui, game: game}
	supportRecording.Store(recorder)
	return recorder.snapshot(), nil
}

func StopSupportRecording() (SupportRecording, error) {
	supportMu.Lock()
	defer supportMu.Unlock()
	recorder := supportRecording.Swap(nil)
	if recorder == nil {
		return lastSupportRecording, ErrNoSupportRecording
	}
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	now := time.Now()
	recorder.status.Active = false
	recorder.status.StoppedAt = &now
	lastSupportRecording = recorder.status
	return lastSupportRecording, errors.Join(recorder.ssui.Close(), recorder.game.Close())
}

func SupportRecordingStatus() SupportRecording {
	if recorder := supportRecording.Load(); recorder != nil {
		return recorder.snapshot()
	}
	supportMu.Lock()
	defer supportMu.Unlock()
	return lastSupportRecording
}

func SupportRecordGameLog(line string) {
	recordSupportLine(line, true)
}

func recordSupportLine(line string, game bool) {
	recorder := supportRecording.Load()
	if recorder == nil {
		return
	}
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	if !recorder.status.Active || recorder.status.Full {
		return
	}
	if recorder.status.Bytes+int64(len(line)) > supportRecordingLimit {
		recorder.status.Full = true
		return
	}
	file := recorder.ssui
	if game {
		file = recorder.game
	}
	written, _ := file.WriteString(line)
	recorder.status.Bytes += int64(written)
}

func (recorder *supportRecorder) snapshot() SupportRecording {
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	return recorder.status
}

func pruneSupportRecordings(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	sort.Strings(dirs)
	for len(dirs) > supportRecordingsKept {
		os.RemoveAll(filepath.Join(root, dirs[0]))
		dirs = dirs[1:]
	}
}
