package runfile

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// GetAllArgs returns all GameArgs from the runfile
func GetAllArgs() []GameArg {
	if CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		return nil
	}
	return CurrentRunfile.getAllArgs()
}

func GetUIGroups() []string {
	if CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		return nil
	}

	groups := make(map[string]bool)
	for _, arg := range CurrentRunfile.getAllArgs() {
		groups[arg.UIGroup] = true
	}

	var result []string
	for group := range groups {
		result = append(result, group)
	}
	return result
}

// GetArgsByGroup returns GameArgs from the rf for use in UI context - and skips args that have arg.Special is not "hide_in_ui"
func GetArgsByGroup(group string) []GameArg {
	if CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		return nil
	}

	var result []GameArg
	for _, arg := range CurrentRunfile.getAllArgs() {
		if arg.UIGroup == group {
			if arg.Special != "hide_in_ui" {
				result = append(result, arg)
			}
		}
	}
	return result
}

// GetSingleArg retrieves a specific GameArg by its flag
func GetSingleArg(flag string) (*GameArg, error) {
	if CurrentRunfile == nil {
		err := ErrRunfileNotLoaded{Msg: "runfile not loaded"}
		logger.Runfile.Error(err.Error())
		return nil, err
	}

	for _, arg := range CurrentRunfile.getAllArgs() {
		if arg.Flag == flag {
			return &arg, nil
		}
	}

	err := ErrArgNotFound{Flag: flag}
	logger.Runfile.Error(fmt.Sprintf("arg not found: flag=%s", flag))
	return nil, err
}

// GetExecutable returns the appropriate executable based on GOOS or Architecture
func (rf *RunFile) GetExecutable() (string, error) {
	goos := strings.ToLower(runtime.GOOS)

	// If Architecture is not all, use it exclusively
	if rf.Architecture != "all" && rf.Architecture != "" {
		arch := strings.ToLower(rf.Architecture)
		if arch != "windows" && arch != "linux" {
			return "", fmt.Errorf("invalid architecture in runfile: %s", rf.Architecture)
		}
		if arch != goos {
			return "", fmt.Errorf("runfile architecture %s does not match current OS %s", arch, goos)
		}
		if arch == "windows" {
			if rf.WindowsExecutable == "" {
				return "", fmt.Errorf("WindowsExecutable is required when architecture is set to windows")
			}
			return rf.WindowsExecutable, nil
		}
		if rf.LinuxExecutable == "" {
			return "", fmt.Errorf("LinuxExecutable is required when architecture is set to linux")
		}
		return rf.LinuxExecutable, nil
	}

	// If Architecture is all or empty, select based on GOOS

	if goos == "windows" {
		if rf.WindowsExecutable == "" {
			return "", fmt.Errorf("WindowsExecutable is required for windows OS")
		}
		return rf.WindowsExecutable, nil
	}
	if goos == "linux" {
		if rf.LinuxExecutable == "" {
			return "", fmt.Errorf("LinuxExecutable is required for linux OS")
		}
		return rf.LinuxExecutable, nil
	}
	return "", fmt.Errorf("unsupported OS: %s", goos)
}
