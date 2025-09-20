package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
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
		"IsConsoleEnabled":         fmt.Sprintf("%v", config.GetIsConsoleEnabled()),
		"AutoStartServerOnStartup": fmt.Sprintf("%v", config.GetAutoStartServerOnStartup()),
		"LanguageSetting":          config.GetLanguageSetting(),
		"ConfigPath":               config.GetConfigPath(),
	}
	printSection("General Configuration", general)

	// Server Configuration
	server := map[string]string{
		"GameBranch":                config.GetGameBranch(),
		"ServerName":                config.GetServerName(),
		"WorldName":                 config.GetWorldName(),
		"BackupWorldName":           config.GetBackupWorldName(),
		"ServerMaxPlayers":          config.GetServerMaxPlayers(),
		"GamePort":                  config.GetGamePort(),
		"UpdatePort":                config.GetUpdatePort(),
		"UPNPEnabled":               fmt.Sprintf("%v", config.GetUPNPEnabled()),
		"AutoSave":                  fmt.Sprintf("%v", config.GetAutoSave()),
		"SaveInterval":              config.GetSaveInterval(),
		"AutoPauseServer":           fmt.Sprintf("%v", config.GetAutoPauseServer()),
		"LocalIpAddress":            config.GetLocalIpAddress(),
		"StartLocalHost":            fmt.Sprintf("%v", config.GetStartLocalHost()),
		"ServerVisible":             fmt.Sprintf("%v", config.GetServerVisible()),
		"UseSteamP2P":               fmt.Sprintf("%v", config.GetUseSteamP2P()),
		"ExePath":                   config.GetExePath(),
		"AdditionalParams":          config.GetAdditionalParams(),
		"GameServerAppID":           config.GetGameServerAppID(),
		"Difficulty":                config.GetDifficulty(),
		"StartCondition":            config.GetStartCondition(),
		"StartLocation":             config.GetStartLocation(),
		"SaveInfo":                  config.GetSaveInfo(),
		"IsNewTerrainAndSaveSystem": fmt.Sprintf("%v", config.GetIsNewTerrainAndSaveSystem()),
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

	// Backup Configuration
	backup := map[string]string{
		"BackupKeepLastN":         fmt.Sprintf("%d", config.GetBackupKeepLastN()),
		"IsCleanupEnabled":        fmt.Sprintf("%v", config.GetIsCleanupEnabled()),
		"BackupKeepDailyFor":      fmt.Sprintf("%v", config.GetBackupKeepDailyFor()),
		"BackupKeepWeeklyFor":     fmt.Sprintf("%v", config.GetBackupKeepWeeklyFor()),
		"BackupKeepMonthlyFor":    fmt.Sprintf("%v", config.GetBackupKeepMonthlyFor()),
		"BackupCleanupInterval":   fmt.Sprintf("%v", config.GetBackupCleanupInterval()),
		"ConfiguredBackupDir":     config.GetConfiguredBackupDir(),
		"ConfiguredSafeBackupDir": config.GetConfiguredSafeBackupDir(),
	}
	printSection("Backup Configuration", backup)

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
		"SSUIIdentifier":       config.GetSSUIIdentifier(),
		"SSUIWebPort":          config.GetSSUIWebPort(),
		"UIModFolder":          config.GetUIModFolder(),
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

// SetupWorkingDir sets the working directory to the directory of the executable to prevent user errors
func SetupWorkingDir() error {
	if runtime.GOOS == "windows" {
		// For now Windows doesn't have symlinking issues so we'll just let is use the current working directory
		return nil
	}
	if runtime.GOOS == "linux" {
		// Get the current executable path from /proc/self/exe
		exePath, err := os.Readlink("/proc/self/exe")
		if err != nil {
			return err
		}
		// Get the directory path of the executable
		dirPath := filepath.Dir(exePath)
		// Change the working directory to the executable's directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if cwd != dirPath {
			logger.Core.Debug("Changing working directory to " + dirPath)
			err = os.Chdir(dirPath)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func SanityCheck() error {

	if runtime.GOOS == "windows" {
		return nil
	}
	// Check if running as root (UID 0)
	if os.Geteuid() == 0 {
		return fmt.Errorf("root: SSUI should not be run as root")
	}

	// Check if current working directory is writable
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Try to create a temporary file to test write permissions
	testFile := filepath.Join(workDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("cannot write to working directory %s: %w", workDir, err)
	}
	// Clean up test file
	if err := os.Remove(testFile); err != nil {
		return fmt.Errorf("failed to clean up sanity check writetest file: %w", err)
	}

	// Check if steamcmd package is installed  (requires further testing, disabled for now)
	//cmd := exec.Command("dpkg-query", "-W", "-f='${Status}'", "steamcmd")
	//output, err := cmd.CombinedOutput()
	//if err == nil && strings.Contains(string(output), "install ok installed") {
	//	return fmt.Errorf("steamcmd package is installed")
	//}

	return nil
}
