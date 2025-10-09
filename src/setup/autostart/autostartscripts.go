package autostart

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

func SetupAutostartScripts() error {
	scriptFS, err := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/scripts")
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		script, err := scriptFS.Open("autostart.ps1")
		if err != nil {
			return err
		}
		defer script.Close()
		data, err := io.ReadAll(script)
		if err != nil {
			return err
		}

		_, err = os.Stat("autostart.ps1")
		if err == nil {
			err = os.Remove("autostart.ps1")
			if err != nil {
				return fmt.Errorf("failed to remove autostart.ps1: %w", err)
			}
		}

		err = os.WriteFile("autostart.ps1", data, 0755)
		if err != nil {
			return err
		}
	}
	if runtime.GOOS == "linux" {
		script, err := scriptFS.Open("autostart.sh")
		if err != nil {
			return err
		}
		defer script.Close()
		data, err := io.ReadAll(script)
		if err != nil {
			return err
		}
		_, err = os.Stat("autostart.sh")
		if err == nil {
			err = os.Remove("autostart.sh")
			if err != nil {
				return fmt.Errorf("failed to remove autostart.sh: %w", err)
			}
		}

		err = os.WriteFile("autostart.sh", data, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
