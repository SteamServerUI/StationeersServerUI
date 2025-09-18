//go:build linux

package autostart

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

func SetupService() error {
	// run ./autostart.sh, which will ask for elevation when needed and setup the service

	//stat /etc/systemd/system/ssui.service and it it exists tell the user that the service is already installed
	_, err := os.Stat("/etc/systemd/system/stationeersserverui.service")
	if err == nil {
		logger.Main.Error("Autostart service is already installed. Exiting in 5 seconds...")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	exec.Command("chmod", "+x", "autostart.sh").Run()
	exec.Command("./autostart.sh", "install").Run()

	err = exec.Command("chmod", "+x", "autostart.sh").Run()
	if err != nil {
		return fmt.Errorf("failed to chmod autostart.sh: %w", err)
	}
	cmd := exec.Command("./autostart.sh")
	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// handle exit codes in a switch statement
			switch exitErr.ExitCode() {
			case 10:
				return fmt.Errorf("autostart.sh failed: For security reasons, it is not recommended to run this service as a root user")
			case 2:
				return fmt.Errorf("autostart.sh Error: systemd is not available on this system")
			case 3:
				return fmt.Errorf("autostart.sh Error: Could not determine base directory")
			case 4:
				return fmt.Errorf("autostart.sh Failed to stop ssui.service")
			case 5:
				return fmt.Errorf("autostart.sh Error: Failed to create service file")
			case 6:
				return fmt.Errorf("autostart.sh Error: Failed to reload systemd daemon")
			case 7:
				return fmt.Errorf("autostart.sh Error: Failed to enable ssui.service")
			case 8:
				return fmt.Errorf("autostart.sh Error: Failed to start ssui.service")
			default:
				return fmt.Errorf("autostart.sh failed with exit code %d: %w", exitErr.ExitCode(), err)
			}
		}

		return fmt.Errorf("failed to run autostart.sh: %w", err)
	}
	// stat for the file to check if it exists, if it does, remove it
	_, err = os.Stat("autostart.sh")
	if err == nil {
		err = os.Remove("autostart.sh")
		if err != nil {
			return fmt.Errorf("failed to remove autostart.sh: %w", err)
		}
	}
	logger.Main.Info("Autostart setup complete. Restart SSUI without the --setupautostart flag to start normally. Exiting in 5 seconds...")
	return nil

}
