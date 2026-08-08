package backupapi

import "testing"

func TestBackupFormat(t *testing.T) {
	tests := map[string]string{
		"backup_2026-08-09.tar.gz": "tar.gz",
		"backup_2026-08-09.tar":    "tar",
		"snapshot_2026-08-09":      "snapshot",
		"backup_2026-08-09":        "copy",
	}
	for name, want := range tests {
		if got := backupFormat(name); got != want {
			t.Fatalf("backupFormat(%q) = %q, want %q", name, got, want)
		}
	}
}
