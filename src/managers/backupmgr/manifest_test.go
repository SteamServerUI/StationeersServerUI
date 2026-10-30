package backupmgr

import (
	"context"
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func primeBackupInventory(t *testing.T, m *BackupManager) {
	t.Helper()
	if err := loadInventory(m); err != nil {
		t.Fatal(err)
	}
	for name := range m.records {
		if _, err := analyzeBackup(context.Background(), m, name); err != nil {
			t.Fatal(err)
		}
	}
}

func TestManifestRestoresAnalysisAndRuntimeReadsUseRAM(t *testing.T) {
	path := analysisFixture(t)
	cfg := BackupConfig{SafeBackupDir: filepath.Dir(path)}
	m := NewBackupManager(cfg)
	primeBackupInventory(t, m)
	if err := saveManifest(m); err != nil {
		t.Fatal(err)
	}
	next := NewBackupManager(cfg)
	if err := loadInventory(next); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(path, path+".hidden"); err != nil {
		t.Fatal(err)
	}
	// Runtime reads must not reopen either the ZIP or the manifest.
	if err := os.Remove(filepath.Join(cfg.SafeBackupDir, manifestFilename)); err != nil {
		t.Fatal(err)
	}
	saves, err := next.ListBackups(0)
	if err != nil || len(saves) != 1 || saves[0].Summary.DaysPlayed != 67 {
		t.Fatalf("list: %+v, %v", saves, err)
	}
	analysis, err := next.AnalyzeBackup(context.Background(), filepath.Base(path))
	if err != nil || analysis.Players != 3 || analysis.Furnaces != 2 {
		t.Fatalf("analysis: %+v, %v", analysis, err)
	}
}

func TestReloadKeepsPendingObservation(t *testing.T) {
	cfg := BackupConfig{BackupDir: t.TempDir(), SafeBackupDir: t.TempDir()}
	name := "130226_173454_auto.save"
	if err := copyFile(analysisFixture(t), filepath.Join(cfg.BackupDir, name)); err != nil {
		t.Fatal(err)
	}
	old := NewBackupManager(cfg)
	if err := loadInventory(old); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	if err := pollBackupsAt(old, now); err != nil {
		t.Fatal(err)
	}
	old.Shutdown()
	next := NewBackupManager(cfg)
	inheritDetectorState(next, old)
	if err := loadInventory(next); err != nil {
		t.Fatal(err)
	}
	if err := pollBackupsAt(next, now.Add(45*time.Second)); err != nil {
		t.Fatal(err)
	}
	ready, identity := nextAutosave(next)
	if ready != name {
		t.Fatal("reload lost the first observation")
	}
	if err := handleBackup(next, ready, identity); err != nil {
		t.Fatal(err)
	}
	if next.records[name].ScanVersion != analysisVersion {
		t.Fatal("handling did not analyze the copied save")
	}
	other := NewBackupManager(BackupConfig{BackupDir: t.TempDir(), SafeBackupDir: cfg.SafeBackupDir})
	inheritDetectorState(other, old)
	if len(other.observed) != 0 {
		t.Fatal("state crossed source folders")
	}
}

func TestRetentionDoesNotResurrectAutosaveAfterRestart(t *testing.T) {
	cfg := BackupConfig{BackupDir: t.TempDir(), SafeBackupDir: t.TempDir()}
	name := "130226_173454_auto.save"
	source := filepath.Join(cfg.BackupDir, name)
	if err := copyFile(analysisFixture(t), source); err != nil {
		t.Fatal(err)
	}
	m := NewBackupManager(cfg)
	if err := m.handleNewBackup(source); err != nil {
		t.Fatal(err)
	}
	if err := m.Cleanup(); err != nil {
		t.Fatal(err)
	}
	assertFileExists(t, filepath.Join(cfg.SafeBackupDir, name), false)
	assertFileExists(t, source, true)
	m.Shutdown()
	next := NewBackupManager(cfg)
	if err := loadInventory(next); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	if err := pollBackupsAt(next, now); err != nil {
		t.Fatal(err)
	}
	if err := pollBackupsAt(next, now.Add(45*time.Second)); err != nil {
		t.Fatal(err)
	}
	if len(next.pending) != 0 || !next.handled[name] {
		t.Fatal("retention-deleted archive was queued again")
	}
	if err := os.Remove(source); err != nil {
		t.Fatal(err)
	}
	if err := pollBackupsAt(next, now.Add(90*time.Second)); err != nil {
		t.Fatal(err)
	}
	if len(next.handled) != 0 {
		t.Fatal("source tombstones grew without bound")
	}
}

func TestManifestFailurePreventsRetentionDeletion(t *testing.T) {
	path := analysisFixture(t)
	m := NewBackupManager(BackupConfig{SafeBackupDir: filepath.Dir(path)})
	primeBackupInventory(t, m)
	// A directory at the manifest path makes publishing fail on both supported OSes.
	if err := os.Mkdir(filepath.Join(m.config.SafeBackupDir, manifestFilename), 0755); err != nil {
		t.Fatal(err)
	}
	if err := m.cleanSafeBackupDir(); err == nil {
		t.Fatal("retention ignored persistence failure")
	}
	assertFileExists(t, path, true)
}

func TestManifestRecoveryAndFutureVersionProtection(t *testing.T) {
	for _, test := range []struct {
		name, data string
		future     bool
	}{
		{"broken", "{truncated", false},
		{"future", `{"version":999,"new-format":{}}`, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := analysisFixture(t)
			manifestPath := filepath.Join(filepath.Dir(path), manifestFilename)
			if err := os.WriteFile(manifestPath, []byte(test.data), 0600); err != nil {
				t.Fatal(err)
			}
			m := NewBackupManager(BackupConfig{SafeBackupDir: filepath.Dir(path)})
			err := loadInventory(m)
			if test.future {
				if err == nil {
					t.Fatal("accepted a newer manifest version")
				}
				m.Shutdown()
				data, err := os.ReadFile(manifestPath)
				if err != nil || string(data) != test.data {
					t.Fatal("future manifest was overwritten")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			quarantined, _ := filepath.Glob(manifestPath + ".invalid-*")
			if len(quarantined) != 1 {
				t.Fatal("broken manifest not preserved")
			}
			primeBackupInventory(t, m)
			if err := saveManifest(m); err != nil {
				t.Fatal(err)
			}
			manifest, err := readManifest(manifestPath)
			if err != nil || len(manifest.Backups) != 1 {
				t.Fatalf("rebuilt manifest: %+v, %v", manifest, err)
			}
		})
	}
}

func TestNestedBackupDownloadAndRestore(t *testing.T) {
	t.Chdir(t.TempDir())
	safe := t.TempDir()
	if err := os.Mkdir(filepath.Join(safe, "nested"), 0755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(safe, "nested", "130226_173454_auto.save")
	if err := copyFile(analysisFixture(t), path); err != nil {
		t.Fatal(err)
	}
	m := NewBackupManager(BackupConfig{WorldName: "restored", SafeBackupDir: safe})
	primeBackupInventory(t, m)
	saves, err := m.ListBackups(0)
	if err != nil || len(saves) != 1 || saves[0].Name != "nested/130226_173454_auto.save" {
		t.Fatalf("nested list: %+v, %v", saves, err)
	}
	download, err := m.GetBackupFileData("nested/130226_173454_auto.save")
	if err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if sha256.Sum256(original) != sha256.Sum256(download.Data) {
		t.Fatal("download bytes changed")
	}
	if err := m.RestoreBackup("nested/130226_173454_auto.save"); err != nil {
		t.Fatal(err)
	}
	restored := filepath.Join("saves", "restored", "restored.save")
	identity, err := identifySave(restored)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateBackupSave(context.Background(), restored, identity); err != nil {
		t.Fatal(err)
	}
	analysis, err := AnalyzeSave(context.Background(), restored)
	if err != nil || analysis.Players != 3 {
		t.Fatalf("restored analysis: %+v, %v", analysis, err)
	}
}

func TestCompletedCopySurvivesSourceRotation(t *testing.T) {
	path := analysisFixture(t)
	expected, err := identifySave(path)
	if err != nil {
		t.Fatal(err)
	}
	temp, err := copyBackupToTemp(path, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := checkCopiedSource(path, expected); err != nil {
		t.Fatal(err)
	}
	identity, err := identifySave(temp)
	if err != nil {
		t.Fatal(err)
	}
	record, err := scanCopiedBackup(context.Background(), temp, identity)
	if err != nil || !analysisReady(record) {
		t.Fatalf("lost completed copy during rotation: %+v, %v", record, err)
	}
	if err := os.WriteFile(path, []byte("different content"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := checkCopiedSource(path, expected); err == nil {
		t.Fatal("accepted a source changed during copy")
	}
}
