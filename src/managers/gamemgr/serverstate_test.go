package gamemgr

import (
	"testing"
	"time"
)

func TestStartupStateBecomesUncertain(t *testing.T) {
	SetServerState(ServerStateStarting)
	serverStateMu.RLock()
	generation := serverStateGeneration
	serverStateMu.RUnlock()

	markStartupUncertainAfter(generation, time.Millisecond)
	if state := GetServerState(); state != ServerStateUncertain {
		t.Fatalf("expected uncertain after startup timeout, got %s", state)
	}
}

func TestCompletedStartupIgnoresOldTimeout(t *testing.T) {
	SetServerState(ServerStateStarting)
	serverStateMu.RLock()
	generation := serverStateGeneration
	serverStateMu.RUnlock()
	SetServerState(ServerStateRunning)

	markStartupUncertainAfter(generation, time.Millisecond)
	if state := GetServerState(); state != ServerStateRunning {
		t.Fatalf("expected running state to survive old timeout, got %s", state)
	}
}
