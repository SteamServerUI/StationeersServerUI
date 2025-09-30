package config

import (
	"fmt"
	"strings"
)

// SetRunfileGame sets the RunfileGame with validation
func SetRunfileIdentifier(value string) error {
	ConfigMu.Lock()
	defer ConfigMu.Unlock()

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("runfile game cannot be empty")
	}

	RunfileIdentifier = value
	return nil
	//return saveConfig()
}
