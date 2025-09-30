package config

import (
	"fmt"
	"strings"
)

// SetRunfileGame sets the RunfileGame with validation
func SetRunfileGame(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("runfile game cannot be empty")
	}

	RunfileGame = value
	return nil
	//return saveConfig()
}
