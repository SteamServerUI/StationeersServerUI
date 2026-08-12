package gamemgr

import (
	"sync"
	"time"
)

type ServerState string

const (
	ServerStateUncertain      ServerState = "uncertain"
	ServerStateStopped        ServerState = "stopped"
	ServerStateStarting       ServerState = "starting"
	ServerStateLoadingMap     ServerState = "loading-map"
	ServerStateHostingSession ServerState = "hosting-session"
	ServerStateRunning        ServerState = "running"
	ServerStateStopping       ServerState = "stopping"
)

const startupStateTimeout = 5 * time.Minute

var (
	serverStateMu         sync.RWMutex
	serverState           = ServerStateUncertain
	serverStateGeneration uint64
)

func GetServerState() ServerState {
	serverStateMu.RLock()
	defer serverStateMu.RUnlock()
	return serverState
}

func SetServerState(state ServerState) {
	serverStateMu.Lock()
	serverState = state
	if state == ServerStateStarting {
		serverStateGeneration++
		generation := serverStateGeneration
		serverStateMu.Unlock()
		go markStartupUncertainAfter(generation, startupStateTimeout)
		return
	}
	if state == ServerStateRunning || state == ServerStateStopping || state == ServerStateStopped {
		serverStateGeneration++
	}
	serverStateMu.Unlock()
}

func markStartupUncertainAfter(generation uint64, delay time.Duration) {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	<-timer.C

	serverStateMu.Lock()
	defer serverStateMu.Unlock()
	if generation != serverStateGeneration {
		return
	}
	switch serverState {
	case ServerStateStarting, ServerStateLoadingMap, ServerStateHostingSession:
		serverState = ServerStateUncertain
	}
}
