package backupmgr

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRetentionDoesNotDeleteExternallyChangedArchive(t *testing.T) {
	safe := t.TempDir()
	path := writeBackupSave(t, safe, "old.save", time.Now().Add(-48*time.Hour))
	m := NewBackupManager(BackupConfig{SafeBackupDir: safe})
	primeBackupInventory(t, m)
	if err := os.WriteFile(path, []byte("a different archive imported by the user"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := m.cleanSafeBackupDir(); err == nil {
		t.Fatal("retention ignored an archive identity change")
	}
	assertFileExists(t, path, true)
}

func TestSourceCleanupKeepsLastCopyAfterArchiveDisappears(t *testing.T) {
	for _, restart := range []bool{false, true} {
		t.Run(map[bool]string{false: "running", true: "restarted"}[restart], func(t *testing.T) {
			source, safe := t.TempDir(), t.TempDir()
			stamp := time.Now().Add(-48 * time.Hour)
			original := writeBackupSave(t, source, "old.save", stamp)
			if err := os.Chtimes(original, stamp, stamp); err != nil {
				t.Fatal(err)
			}
			cfg := BackupConfig{BackupDir: source, SafeBackupDir: safe, RetentionPolicy: RetentionPolicy{KeepNewestCount: 1}}
			m := NewBackupManager(cfg)
			if err := m.handleNewBackup(original); err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(filepath.Join(safe, "old.save")); err != nil {
				t.Fatal(err)
			}
			if restart {
				m.Shutdown()
				m = NewBackupManager(cfg)
				if err := loadInventory(m); err != nil {
					t.Fatal(err)
				}
			}
			_ = m.Cleanup() // It may report the missing archive, but cannot remove the source.
			assertFileExists(t, original, true)
		})
	}
}

func TestBackupFoldersCannotOverlap(t *testing.T) {
	root := t.TempDir()
	for _, test := range []struct {
		name, source, safe string
		valid              bool
	}{
		{"siblings", filepath.Join(root, "autosave"), filepath.Join(root, "Safebackups"), true},
		{"same", root, root, false},
		{"destination in source", root, filepath.Join(root, "Safebackups"), false},
		{"source in destination", filepath.Join(root, "autosave"), root, false},
		{"similar prefix", filepath.Join(root, "save"), filepath.Join(root, "save-backups"), true},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := checkBackupFolders(BackupConfig{BackupDir: test.source, SafeBackupDir: test.safe})
			if (err == nil) != test.valid {
				t.Fatalf("folder validation: %v", err)
			}
		})
	}
	t.Run("symlink parent", func(t *testing.T) {
		link := filepath.Join(t.TempDir(), "linked")
		if err := os.Symlink(root, link); err != nil {
			t.Skipf("cannot create symlinks: %v", err)
		}
		if err := checkBackupFolders(BackupConfig{BackupDir: root, SafeBackupDir: filepath.Join(link, "not-created-yet")}); err == nil {
			t.Fatal("symlink concealed overlapping directories")
		}
	})
}
