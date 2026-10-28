package diagnostics

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
)

func StartupCheck() error {
	if runtime.GOOS != "windows" && os.Geteuid() == 0 && !config.GetIsDockerContainer() {
		return errors.New("root: SSUI should not be run as root")
	}
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	if err := writeTest(workDir); err != nil {
		return fmt.Errorf("cannot write to working directory %s: %w", workDir, err)
	}
	scan := scanAccess(0)
	if scan.Problems > 0 {
		return fmt.Errorf("SSUI cannot access these paths in its data or save folders (%d): %s", scan.Problems, strings.Join(scan.Issues, "; "))
	}
	return nil
}
