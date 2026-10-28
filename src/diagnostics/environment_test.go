package diagnostics

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestReadOnlyFolderIsReported(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("permissions are not enforced here")
	}
	dir := t.TempDir()
	locked := filepath.Join(dir, "locked")
	if err := os.Mkdir(locked, 0o555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(locked, 0o755)

	if status := inspectDirectory(locked); status != "not writable" {
		t.Fatalf("got %q", status)
	}
	report := Report{Environment: Environment{UID: os.Getuid(), Volumes: []Volume{{Name: "saves", Path: locked, Status: inspectDirectory(locked)}}}}
	report.checkEnvironment()
	if len(report.Findings) != 1 || report.Findings[0].Code != codeReadOnly {
		t.Fatalf("findings: %+v", report.Findings)
	}
}

func TestLowDiskSpaceWarns(t *testing.T) {
	report := Report{Environment: Environment{Volumes: []Volume{{Name: "saves", Status: "ok", FreeMB: 500, TotalMB: 100000}}}}
	report.checkEnvironment()
	if len(report.Findings) != 1 || report.Findings[0].Code != codeDisk {
		t.Fatalf("findings: %+v", report.Findings)
	}
}
