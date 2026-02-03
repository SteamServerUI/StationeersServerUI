// passwordrotation.go
package gamemgr

import (
	"math/rand"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// Predefined list of Stationeers-themed passwords
var stationeersPasswords = []string{
	"Stationeer",
	"Mimas",
	"Europa",
	"Lunar",
	"Vulcan",
	"Venus",
	"Rocket",
	"Oxygen",
	"Nitrogen",
	"Volatiles",
	"Hardsuit",
	"Jetpack",
	"Airlock",
	"Station",
}

// rotatePasswordIfEnabled checks if password rotation is enabled and sets a new random password
func rotatePasswordIfEnabled() {

	if !config.GetIsDiscordEnabled() || config.GetServerInfoPanelChannelID() == "" || !config.GetRotateServerPassword() {
		return
	}

	// Seed the random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Select a random password from the list
	newPassword := stationeersPasswords[rng.Intn(len(stationeersPasswords))]

	// Set the new password in config
	err := config.SetServerPassword(newPassword)
	if err != nil {
		logger.Core.Error("Failed to set rotated server password: " + err.Error())
		return
	}

	logger.Core.Info("Password rotation enabled - new server password set: " + newPassword)
}
