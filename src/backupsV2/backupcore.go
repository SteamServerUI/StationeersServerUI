package backupsv2

import (
	"context"
	"os"
	"sync"
	"time"
)

const (
	defaultWaitTime = 30 * time.Second
)

// BackupConfig holds configuration for backup operations
type BackupConfig struct {
	WorldName       string
	BackupDir       string
	SafeBackupDir   string
	RetentionPolicy RetentionPolicy
	WaitTime        time.Duration
}

// RetentionPolicy defines backup retention rules
type RetentionPolicy struct {
	KeepLastN       int           // Keep last N backups regardless of age
	KeepDailyFor    time.Duration // Keep daily backups for this duration
	KeepWeeklyFor   time.Duration // Keep weekly backups for this duration
	KeepMonthlyFor  time.Duration // Keep monthly backups for this duration
	CleanupInterval time.Duration // How often to run cleanup
}

// BackupGroup represents a set of backup files
type BackupGroup struct {
	Index    int
	BinFile  string
	XMLFile  string
	MetaFile string
	ModTime  time.Time
}

// BackupManager manages backup operations
type BackupManager struct {
	config  BackupConfig
	mu      sync.Mutex
	watcher *fsWatcher
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewBackupManager creates a new BackupManager instance
func NewBackupManager(cfg BackupConfig) *BackupManager {
	ctx, cancel := context.WithCancel(context.Background())

	if cfg.WaitTime == 0 {
		cfg.WaitTime = defaultWaitTime
	}

	return &BackupManager{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Initialize sets up required directories
func (m *BackupManager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := os.MkdirAll(m.config.BackupDir, os.ModePerm); err != nil {
		return err
	}
	return os.MkdirAll(m.config.SafeBackupDir, os.ModePerm)
}

// Shutdown stops all backup operations
func (m *BackupManager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cancel != nil {
		m.cancel()
	}
	if m.watcher != nil {
		m.watcher.close()
	}
}
