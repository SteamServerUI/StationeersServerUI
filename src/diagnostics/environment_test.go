package diagnostics

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func TestAccessScanFindsLockedPaths(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("permissions are not enforced here")
	}
	dir := t.TempDir()
	secret := filepath.Join(dir, "secret.save")
	if err := os.WriteFile(secret, []byte("x"), 0o000); err != nil {
		t.Fatal(err)
	}
	locked := filepath.Join(dir, "locked")
	if err := os.Mkdir(locked, 0o555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(locked, 0o755)

	readOnly := AccessScan{Issues: []string{}}
	readOnly.tree(dir, false, 0)
	if readOnly.Problems != 1 {
		t.Fatalf("read-only scan: %+v", readOnly)
	}
	writable := AccessScan{Issues: []string{}}
	writable.tree(dir, true, 0)
	if writable.Problems != 2 {
		t.Fatalf("writable scan: %+v", writable)
	}
	limited := AccessScan{Issues: []string{}}
	limited.tree(dir, true, 1)
	if !limited.Truncated {
		t.Fatalf("limited scan: %+v", limited)
	}
	missing := AccessScan{Issues: []string{}}
	missing.tree(filepath.Join(dir, "nope"), true, 0)
	if missing.Problems != 0 || missing.Checked != 0 {
		t.Fatalf("missing root: %+v", missing)
	}
}

func TestWriteTextRedactsAndSummarizes(t *testing.T) {
	report := Report{Overall: "failed", Findings: []Finding{{Code: codeDisk, Status: StatusWarn, Summary: "Only 5 MB are free.", Fix: "Free up disk space."}}, Environment: Environment{UID: 1000, GID: 1000, Volumes: []Volume{{Name: "saves", Path: "./saves", Status: "not writable"}}}, Access: &AccessScan{Checked: 3, Problems: 1, Issues: []string{"saves/world.xml cannot be read"}}}
	var text strings.Builder
	if err := WriteText(&text, report); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Overall: failed", "FullBowlCat: Only 5 MB are free.", "Fix: Free up disk space.", "saves    not writable, ./saves", "3 paths checked, 1 problems", "saves/world.xml cannot be read", "Full data:"} {
		if !strings.Contains(text.String(), want) {
			t.Fatalf("missing %q in:\n%s", want, text.String())
		}
	}
}

func TestDNSFindings(t *testing.T) {
	partial := Report{DNS: []DNSResult{{Host: "a", OK: true}, {Host: "b", Purpose: "updates"}}}
	partial.checkDNS()
	if len(partial.Findings) != 1 || partial.Findings[0].Status != StatusWarn || !strings.Contains(partial.Findings[0].Summary, "b (updates)") {
		t.Fatalf("partial: %+v", partial.Findings)
	}
	none := Report{DNS: []DNSResult{{Host: "a"}, {Host: "b"}}}
	none.checkDNS()
	if len(none.Findings) != 1 || none.Findings[0].Status != StatusFail {
		t.Fatalf("none: %+v", none.Findings)
	}
}
