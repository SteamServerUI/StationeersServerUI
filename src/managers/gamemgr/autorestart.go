package gamemgr

import (
	"strconv"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/commandmgr"
)

var (
	autoRestartDone chan struct{}
	// other local vars are defined in processmanagement.go
)

// startAutoRestart runs a goroutine that restarts the server either after a specified duration in minutes
// or at a specific time of day (HH:MM) every day.
func startAutoRestart(schedule string, done chan struct{}) {
	// Try parsing as a time in HH:MM format
	if t, err := time.Parse("15:04", schedule); err == nil {
		// Valid HH:MM format, schedule daily restart
		go scheduleDailyRestart(t, done)
		return
	}

	// Fallback to parsing as minutes duration
	minutesInt, err := strconv.Atoi(schedule)
	if err != nil {
		logger.Core.Error("Invalid AutoRestartServerTimer format: " + schedule)
		return
	}
	if minutesInt <= 0 {
		logger.Core.Error("AutoRestartServerTimer must be a positive number of minutes or valid HH:MM time")
		return
	}

	ticker := time.NewTicker(time.Duration(minutesInt) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mu.Lock()
			if !internalIsServerRunningNoLock() {
				mu.Unlock()
				logger.Core.Info("Auto-restart skipped: server is not running")
				return
			}
			mu.Unlock()

			if config.IsSSCMEnabled {
				commandmgr.WriteCommand("say Attention, server is restarting in 30 seconds!")
				time.Sleep(10 * time.Second)
				commandmgr.WriteCommand("say Attention, server is restarting in 20 seconds!")
				time.Sleep(10 * time.Second)
				commandmgr.WriteCommand("say Attention, server is restarting in 10 seconds!")
				time.Sleep(5 * time.Second)
				commandmgr.WriteCommand("say Attention, server is restarting in 5 seconds!")
				time.Sleep(5 * time.Second)
			}
			logger.Core.Info("Auto-restart triggered: stopping server")
			if err := InternalStopServer(); err != nil {
				logger.Core.Error("Auto-restart failed to stop server: " + err.Error())
				return
			}

			logger.Core.Info("Auto-restart: waiting 5 seconds before restarting")
			time.Sleep(5 * time.Second)

			logger.Core.Info("Auto-restart: starting server")
			if err := InternalStartServer(); err != nil {
				logger.Core.Error("Auto-restart failed to start server: " + err.Error())
				return
			}
		case <-done:
			return
		}
	}
}

// scheduleDailyRestart schedules a server restart at the specified time of day (HH:MM) every day.
func scheduleDailyRestart(t time.Time, done chan struct{}) {
	// Extract hour and minute from the parsed time
	hour, min := t.Hour(), t.Minute()

	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
		if now.After(next) || now.Equal(next) {
			// If the time is in the past or now, schedule for tomorrow
			next = next.Add(24 * time.Hour)
		}
		duration := next.Sub(now)

		// Wait until the next restart time or until interrupted
		timer := time.NewTimer(duration)
		select {
		case <-timer.C:
			mu.Lock()
			if !internalIsServerRunningNoLock() {
				mu.Unlock()
				logger.Core.Info("Auto-restart skipped: server is not running")
				continue
			}
			mu.Unlock()

			if config.IsSSCMEnabled {
				commandmgr.WriteCommand("say Attention, server is restarting in 30 seconds!")
				time.Sleep(10 * time.Second)
				commandmgr.WriteCommand("say Attention, server is restarting in 20 seconds!")
				time.Sleep(10 * time.Second)
				commandmgr.WriteCommand("say Attention, server is restarting in 10 seconds!")
				time.Sleep(5 * time.Second)
				commandmgr.WriteCommand("say Attention, server is restarting in 5 seconds!")
				time.Sleep(5 * time.Second)
			}
			logger.Core.Info("Daily auto-restart triggered: stopping server")
			if err := InternalStopServer(); err != nil {
				logger.Core.Error("Daily auto-restart failed to stop server: " + err.Error())
				continue
			}

			logger.Core.Info("Daily auto-restart: waiting 5 seconds before restarting")
			time.Sleep(5 * time.Second)

			logger.Core.Info("Daily auto-restart: starting server")
			if err := InternalStartServer(); err != nil {
				logger.Core.Error("Daily auto-restart failed to start server: " + err.Error())
				continue
			}
		case <-done:
			timer.Stop()
			return
		}
	}
}
