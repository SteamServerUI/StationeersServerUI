package backupmgr

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMissingArchiveIsHandledAgain(t *testing.T) {
	for _, mode := range []string{"running", "restart", "reload", "reload without manifest", "cleanup before detection"} {
		t.Run(mode, func(t *testing.T) {
			cfg := BackupConfig{BackupDir: t.TempDir(), SafeBackupDir: t.TempDir()}
			name := "040926_172441_auto.save"
			source := filepath.Join(cfg.BackupDir, name)
			if err := copyFile(analysisFixture(t), source); err != nil {
				t.Fatal(err)
			}
			m := NewBackupManager(cfg)
			if err := m.handleNewBackup(source); err != nil {
				t.Fatal(err)
			}
			archive := filepath.Join(cfg.SafeBackupDir, name)
			if err := os.Remove(archive); err != nil {
				t.Fatal(err)
			}
			if mode == "cleanup before detection" {
				if err := m.Cleanup(); err == nil {
					t.Fatal("cleanup did not notice the missing archive")
				}
				if m.retired[name] {
					t.Fatal("cleanup claimed an already-missing archive as its own deletion")
				}
			}
			if mode == "reload without manifest" {
				if err := os.Remove(filepath.Join(cfg.SafeBackupDir, manifestFilename)); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "restart" || strings.HasPrefix(mode, "reload") {
				m.Shutdown()
				next := NewBackupManager(cfg)
				if strings.HasPrefix(mode, "reload") {
					inheritDetectorState(next, m)
				}
				m = next
				if err := loadInventory(m); err != nil {
					t.Fatal(err)
				}
			}
			defer m.Shutdown()
			now := time.Now()
			if err := pollBackupsAt(m, now); err != nil {
				t.Fatal(err)
			}
			if len(m.pending) != 0 {
				t.Fatal("repair skipped the stability interval")
			}
			if saves, err := m.ListBackups(0); err != nil || len(saves) != 0 {
				t.Fatalf("missing archive still listed: %+v, %v", saves, err)
			}
			if err := pollBackupsAt(m, now.Add(m.config.WaitTime)); err != nil {
				t.Fatal(err)
			}
			ready, identity := nextAutosave(m)
			if ready != name {
				t.Fatalf("missing archive was not queued: %q", ready)
			}
			if err := handleBackup(m, ready, identity); err != nil {
				t.Fatal(err)
			}
			if !analysisReady(m.records[name]) || m.records[name].Analysis.Players != 3 {
				t.Fatal("repair did not run the complete handling pipeline")
			}
			original, _ := os.ReadFile(source)
			repaired, err := os.ReadFile(archive)
			if err != nil || string(original) != string(repaired) {
				t.Fatal("repair changed the archive bytes")
			}
		})
	}
}

func TestRetentionIntentSurvivesInterruptedManifestUpdate(t *testing.T) {
	cfg := BackupConfig{BackupDir: t.TempDir(), SafeBackupDir: t.TempDir()}
	name := "040926_172441_auto.save"
	source := filepath.Join(cfg.BackupDir, name)
	if err := copyFile(analysisFixture(t), source); err != nil {
		t.Fatal(err)
	}
	m := NewBackupManager(cfg)
	if err := m.handleNewBackup(source); err != nil {
		t.Fatal(err)
	}
	// Emulate the durable deletion intent followed by a crash after unlink,
	// before the removed inventory record is written to the manifest.
	m.retired[name] = true
	m.revision++
	if err := saveManifest(m); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(cfg.SafeBackupDir, name)); err != nil {
		t.Fatal(err)
	}
	next := NewBackupManager(cfg)
	if err := loadInventory(next); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	for _, at := range []time.Time{now, now.Add(next.config.WaitTime)} {
		if err := pollBackupsAt(next, at); err != nil {
			t.Fatal(err)
		}
	}
	if len(next.pending) != 0 || !next.retired[name] {
		t.Fatal("intentional deletion was mistaken for a missing archive")
	}
	if err := os.Remove(source); err != nil {
		t.Fatal(err)
	}
	if err := pollBackupsAt(next, now.Add(2*next.config.WaitTime)); err != nil {
		t.Fatal(err)
	}
	if len(next.retired) != 0 || len(next.handled) != 0 {
		t.Fatal("retention tombstones outlived the source")
	}
}

func TestArchiveRepairDoesNotOverwriteChangedFile(t *testing.T) {
	cfg := BackupConfig{BackupDir: t.TempDir(), SafeBackupDir: t.TempDir()}
	source := filepath.Join(cfg.BackupDir, "changed.save")
	if err := copyFile(analysisFixture(t), source); err != nil {
		t.Fatal(err)
	}
	m := NewBackupManager(cfg)
	if err := m.handleNewBackup(source); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(cfg.SafeBackupDir, "changed.save")
	if err := os.WriteFile(archive, []byte("externally replaced"), 0600); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	for _, at := range []time.Time{now, now.Add(m.config.WaitTime)} {
		if err := pollBackupsAt(m, at); err != nil {
			t.Fatal(err)
		}
	}
	if len(m.pending) != 0 {
		t.Fatal("queued an existing archive for replacement")
	}
	data, _ := os.ReadFile(archive)
	if string(data) != "externally replaced" {
		t.Fatal("overwrote the user's archive")
	}
}

func TestUnavailableArchiveFolderPreservesDetectorState(t *testing.T) {
	cfg := BackupConfig{BackupDir: t.TempDir(), SafeBackupDir: filepath.Join(t.TempDir(), "safe")}
	source := filepath.Join(cfg.BackupDir, "kept.save")
	if err := copyFile(analysisFixture(t), source); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(cfg.SafeBackupDir, 0755); err != nil {
		t.Fatal(err)
	}
	m := NewBackupManager(cfg)
	if err := m.handleNewBackup(source); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(cfg.SafeBackupDir, cfg.SafeBackupDir+".offline"); err != nil {
		t.Fatal(err)
	}
	if err := pollBackups(m); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected unavailable destination: %v", err)
	}
	if _, exists := m.records["kept.save"]; !exists || !m.handled["kept.save"] || len(m.pending) != 0 {
		t.Fatal("unavailable storage was mistaken for a missing archive")
	}
}
