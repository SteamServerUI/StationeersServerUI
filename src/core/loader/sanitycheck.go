package loader

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/diagnostics"
)

func runSanityCheck() error {
	if runtime.GOOS != "windows" {
		config.SetIsDockerContainer(len(diagnostics.ContainerSignals()) > 0)
	}
	if config.GetNoSanityCheck() {
		return nil
	}

	if runtime.GOOS != "windows" {
		exePath, err := os.Readlink("/proc/self/exe")
		if err != nil {
			return err
		}
		dirPath := filepath.Dir(exePath)
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		if cwd != dirPath && !isGoRunBuildDir(dirPath) {
			if err := os.Chdir(dirPath); err != nil {
				return err
			}
		}
	}

	return diagnostics.StartupCheck()
}

// isGoRunBuildDir reports whether executableDir belongs to a temporary or
// cached go run build. Go 1.24 began caching go run executables in the Go
// build cache, so they are no longer guaranteed to live below /tmp.
func isGoRunBuildDir(executableDir string) bool {
	cleanDir := filepath.ToSlash(filepath.Clean(executableDir))
	return strings.Contains(cleanDir, "/tmp") ||
		strings.Contains(cleanDir, "/go-build/") ||
		strings.HasSuffix(cleanDir, "/go-build")
}
