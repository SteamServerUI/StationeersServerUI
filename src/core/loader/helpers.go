package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

func PrintConfigDetails(logLevel ...string) {
	logger.Config.Debug("=== Game Server Configuration Details ===")

	// Helper function to print sections
	printSection := func(title string, fields map[string]string) {

		if logLevel == nil {
			logger.Config.Debug(fmt.Sprintf("\n%s", title))
			logger.Config.Debug(strings.Repeat("-", len(title)))
			for key, value := range fields {
				logger.Config.Debug(fmt.Sprintf("%-30s: %s", key, value))
			}
		}
		if len(logLevel) > 0 && logLevel[0] == "Info" {
			logger.Config.Info(fmt.Sprintf("\n%s", title))
			logger.Config.Info(strings.Repeat("-", len(title)))
			for key, value := range fields {
				logger.Config.Info(fmt.Sprintf("%-30s: %s", key, value))
			}
		}
	}

	// General Configuration
	general := map[string]string{
		"Branch":                   config.GetBranch(),
		"Version":                  config.GetVersion(),
		"IsFirstTimeSetup":         fmt.Sprintf("%v", config.GetIsFirstTimeSetup()),
		"IsDebugMode":              fmt.Sprintf("%v", config.GetIsDebugMode()),
		"IsSSUICLIConsoleEnabled":  fmt.Sprintf("%v", config.GetIsSSUICLIConsoleEnabled()),
		"AutoStartServerOnStartup": fmt.Sprintf("%v", config.GetAutoStartServerOnStartup()),
		"LanguageSetting":          config.GetLanguageSetting(),
		"ConfigPath":               config.GetConfigPath(),
	}
	printSection("General Configuration", general)

	// Server Configuration
	server := map[string]string{
		"GameBranch":      config.GetGameBranch(),
		"GameServerAppID": config.GetGameServerAppID(),
	}
	printSection("Server Configuration", server)

	// Discord Configuration
	discord := map[string]string{
		"IsDiscordEnabled":        fmt.Sprintf("%v", config.GetIsDiscordEnabled()),
		"ControlChannelID":        config.GetControlChannelID(),
		"StatusChannelID":         config.GetStatusChannelID(),
		"ConnectionListChannelID": config.GetConnectionListChannelID(),
		"LogChannelID":            config.GetLogChannelID(),
		"SaveChannelID":           config.GetSaveChannelID(),
		"ControlPanelChannelID":   config.GetControlPanelChannelID(),
		"ErrorChannelID":          config.GetErrorChannelID(),
		"DiscordCharBufferSize":   fmt.Sprintf("%d", config.GetDiscordCharBufferSize()),
		"BlackListFilePath":       config.GetBlackListFilePath(),
	}
	printSection("Discord Configuration", discord)

	// Authentication Configuration
	auth := map[string]string{
		"AuthEnabled":       fmt.Sprintf("%v", config.GetAuthEnabled()),
		"AuthTokenLifetime": fmt.Sprintf("%d", config.GetAuthTokenLifetime()),
	}
	printSection("Authentication Configuration", auth)

	// Logging Configuration
	logging := map[string]string{
		"CreateSSUILogFile":   fmt.Sprintf("%v", config.GetCreateSSUILogFile()),
		"LogLevel":            fmt.Sprintf("%d", config.GetLogLevel()),
		"LogClutterToConsole": fmt.Sprintf("%v", config.GetLogClutterToConsole()),
		"SubsystemFilters":    fmt.Sprintf("%v", config.GetSubsystemFilters()),
		"LogFolder":           config.GetLogFolder(),
	}
	printSection("Logging Configuration", logging)

	// Updater Configuration
	updater := map[string]string{
		"IsUpdateEnabled":        fmt.Sprintf("%v", config.GetIsUpdateEnabled()),
		"AllowPrereleaseUpdates": fmt.Sprintf("%v", config.GetAllowPrereleaseUpdates()),
		"AllowMajorUpdates":      fmt.Sprintf("%v", config.GetAllowMajorUpdates()),
		"AutoRestartServerTimer": config.GetAutoRestartServerTimer(),
		"AutoGameServerUpdates":  fmt.Sprintf("%v", config.GetAllowAutoGameServerUpdates()),
		"CurrentBranchBuildID":   config.GetCurrentBranchBuildID(),
	}
	printSection("Updater Configuration", updater)

	// SSCM Configuration
	sscm := map[string]string{
		"IsSSCMEnabled": fmt.Sprintf("%v", config.GetIsSSCMEnabled()),
		"SSCMFilePath":  config.GetSSCMFilePath(),
		"SSCMPluginDir": config.GetSSCMPluginDir(),
		"SSCMWebDir":    config.GetSSCMWebDir(),
	}
	printSection("SSCM Configuration", sscm)

	// UI Configuration
	ui := map[string]string{
		"BackendName":          config.GetBackendName(),
		"BackendEndpointPort":  config.GetBackendEndpointPort(),
		"SSUIFolder":           config.GetSSUIFolder(),
		"MaxSSEConnections":    fmt.Sprintf("%d", config.GetMaxSSEConnections()),
		"SSEMessageBufferSize": fmt.Sprintf("%d", config.GetSSEMessageBufferSize()),
	}
	printSection("UI Configuration", ui)

	// TLS Configuration
	tls := map[string]string{
		"TLSCertPath": config.GetTLSCertPath(),
		"TLSKeyPath":  config.GetTLSKeyPath(),
	}
	printSection("TLS Configuration", tls)

	// Custom Detections
	custom := map[string]string{
		"CustomDetectionsFilePath": config.GetCustomDetectionsFilePath(),
	}
	printSection("Custom Detections Configuration", custom)

	logger.Config.Debug("=======================================")
}

func IsInsideContainer(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	// Check .dockerenv file (Docker-specific)
	if _, err := os.Stat("/.dockerenv"); err == nil {
		config.SetIsDockerContainer(true)
		return
	}
	// Check cgroup (works for Docker and other container runtimes)
	file, err := os.Open("/proc/1/cgroup")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check for various container runtime indicators
		if strings.Contains(line, "docker") ||
			strings.Contains(line, "containerd") ||
			strings.Contains(line, "kubepods") ||
			strings.Contains(line, "crio") ||
			strings.Contains(line, "libpod") {
			config.SetIsDockerContainer(true)
			return
		}
	}
}
