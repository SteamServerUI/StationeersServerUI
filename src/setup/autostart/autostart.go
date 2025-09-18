package autostart

import (
	"fmt"
)

func Initialize() error {

	err := SetupAutostartScripts()
	if err != nil {
		return fmt.Errorf("failed to setup autostart script, cannot proceed with autostart setup: %w", err)
	}

	err = SetupBinarySymlink()
	if err != nil {
		return fmt.Errorf("failed to create symlink for autostart: %w", err)
	}

	err = SetupService()
	if err != nil {
		return fmt.Errorf("failed to setup service: %w", err)
	}
	return nil
}
