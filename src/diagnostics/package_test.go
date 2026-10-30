package diagnostics

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/logger"
)

func TestPackageWithoutRecordingHasNoLogs(t *testing.T) {
	var buffer bytes.Buffer
	if err := WritePackage(&buffer); err != nil {
		t.Fatal(err)
	}
	files := readZip(t, buffer.Bytes())
	if len(files) != 2 || files["diagnosis.json"] == "" || files["report.html"] == "" {
		t.Fatalf("unexpected package contents: %v", keys(files))
	}
}

func TestRecordingPackageRedactsKnownSecrets(t *testing.T) {
	config.ConfigMu.Lock()
	oldFolder, oldPassword := config.LogFolder, config.ServerPassword
	config.LogFolder, config.ServerPassword = t.TempDir(), "hunter2-secret"
	config.ConfigMu.Unlock()
	t.Cleanup(func() {
		config.ConfigMu.Lock()
		config.LogFolder, config.ServerPassword = oldFolder, oldPassword
		config.ConfigMu.Unlock()
	})

	if _, err := logger.StartSupportRecording(); err != nil {
		t.Fatal(err)
	}
	logger.Core.Debug("joining with password hunter2-secret")
	logger.SupportRecordGameLog("game says hunter2-secret")
	if err := WriteRecordingPackage(io.Discard); !errors.Is(err, ErrNoRecording) {
		t.Fatalf("active recording was exported: %v", err)
	}
	if _, err := StopRecording(); err != nil {
		t.Fatal(err)
	}

	var buffer bytes.Buffer
	if err := WriteRecordingPackage(&buffer); err != nil {
		t.Fatal(err)
	}
	files := readZip(t, buffer.Bytes())
	for _, name := range []string{"ssui.log", "gameserver.log"} {
		if !strings.Contains(files[name], "[REDACTED]") || strings.Contains(files[name], "hunter2-secret") {
			t.Fatalf("%s was not redacted: %q", name, files[name])
		}
	}
}

func readZip(t *testing.T, data []byte) map[string]string {
	t.Helper()
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]string{}
	for _, file := range reader.File {
		content, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(content)
		content.Close()
		if err != nil {
			t.Fatal(err)
		}
		files[file.Name] = string(data)
	}
	return files
}

func keys(files map[string]string) []string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	return names
}
