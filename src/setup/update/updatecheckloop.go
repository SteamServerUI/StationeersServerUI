package update

import (
	"sync"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// UpdateInfo holds the shared state about available updates
var UpdateInfo struct {
	sync.RWMutex

	Available bool   // true if an update is available
	Version   string // the available version string
}

func init() {
	UpdateInfo.Available = false
	UpdateInfo.Version = ""
}

// StartUpdateCheckLoop runs in the background and checks for updates every 6 hours
func StartUpdateCheckLoop() {
	for {
		// Acquire write lock for the duration of the check
		UpdateInfo.Lock()

		err, newVersion := Update(false)
		if err != nil {
			logger.Install.Warn("⚠️ Automatic SSUI Update check failed: " + err.Error())
		} else if newVersion != "" {
			UpdateInfo.Available = true
			UpdateInfo.Version = newVersion
		} else {
			UpdateInfo.Available = false
			UpdateInfo.Version = ""
		}

		UpdateInfo.Unlock()

		time.Sleep(6 * time.Hour)
	}
}
