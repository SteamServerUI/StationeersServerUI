package loader

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

// Define flags matching the config variable names
var backendEndpointPortFlag string
var gameBranchFlag string
var logLevelFlag int
var isDebugModeFlag bool
var createSSUILogFileFlag bool
var recoveryPasswordFlag string
var devModeFlag bool
var skipSteamCMDFlag bool
var sanityCheckFlag bool
var SetCustomWorkDirFlag string

// ParseFlags parses command-line arguments ONCE at startup (called from func main)
func ParseFlags() {
	flag.StringVar(&backendEndpointPortFlag, "BackendEndpointPort", "", "Override the backend endpoint port (e.g., 8080)")
	flag.StringVar(&backendEndpointPortFlag, "p", "", "(Alias) Override the backend endpoint port (e.g., 8080)")
	flag.StringVar(&gameBranchFlag, "GameBranch", "", "Override the game branch (e.g., beta)")
	flag.StringVar(&gameBranchFlag, "b", "", "(Alias) Override the game branch (e.g., beta)")
	flag.StringVar(&recoveryPasswordFlag, "RecoveryPassword", "", "Reset or create the recovery owner (expects password as argument)")
	flag.StringVar(&recoveryPasswordFlag, "r", "", "(Alias) Reset or create the recovery owner")
	flag.BoolVar(&devModeFlag, "dev", false, "Enable development owner and CLI console. For development only.")
	flag.IntVar(&logLevelFlag, "LogLevel", 0, "Override the log level (e.g., 10)")
	flag.IntVar(&logLevelFlag, "ll", 0, "(Alias) Override the log level (e.g., 10)")
	flag.BoolVar(&isDebugModeFlag, "IsDebugMode", false, "Enable debug mode")
	flag.BoolVar(&isDebugModeFlag, "debug", false, "(Alias) Enable debug mode")
	flag.BoolVar(&createSSUILogFileFlag, "LogToFiles", false, "Create log files for SSUI")
	flag.BoolVar(&createSSUILogFileFlag, "lf", false, "(Alias) Create log files for SSUI")
	flag.BoolVar(&skipSteamCMDFlag, "NoSteamCMD", false, "Skips SteamCMD installation")
	flag.BoolVar(&sanityCheckFlag, "NoSanityCheck", false, "Skips the sanity check. Not recommended.")
	flag.StringVar(&SetCustomWorkDirFlag, "SetCustomWorkDir", "", "Sets a custom workdir. Not recommended for production use, but possible. (e.g., /home/steam/SSUI/)")

	// Parse command-line flags
	flag.Parse()
}

// HandleCmdArgs handles command-line arguments ONCE at startup (called from func main) and applies them using the config setters.
// Because this is using the config rather than adding features to it, it is a part of the loader package.
func HandleFlags(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	if devModeFlag {
		if _, err := security.RecoverOwner("admin", "adminadmin", time.Now()); err != nil {
			logger.Main.Error("Failed to prepare development owner: " + err.Error())
			return
		}
		config.SetIsSSUICLIConsoleEnabled(true)
		logger.Main.Warn("Dev mode enabled: owner admin/adminadmin and CLI console are active")
	}

	if skipSteamCMDFlag {
		config.SetSkipSteamCMD(true)
	}

	if backendEndpointPortFlag != "" && backendEndpointPortFlag != "8443" {
		oldPort := config.GetBackendEndpointPort()
		config.SetBackendEndpointPort(backendEndpointPortFlag)
		logger.Main.Info(fmt.Sprintf("Overriding BackendEndpointPort from command line: Before=%s, Now=%s", oldPort, backendEndpointPortFlag))
	}

	if gameBranchFlag != "" {
		oldBranch := config.GetGameBranch()
		config.SetGameBranch(gameBranchFlag)
		logger.Main.Info(fmt.Sprintf("Overriding GameBranch from command line: Before=%s, Now=%s", oldBranch, gameBranchFlag))
	}

	if recoveryPasswordFlag != "" {
		recoveryPasswordFlag = strings.TrimSpace(recoveryPasswordFlag)
		if recoveryPasswordFlag == "" {
			logger.Main.Error("Recovery flag provided but password is empty")
		} else {
			if _, err := security.RecoverOwner("recovery", recoveryPasswordFlag, time.Now()); err != nil {
				logger.Main.Error("Failed to create recovery owner: " + err.Error())
				return
			}
			logger.Main.Warn("Recovery owner created; all existing sessions and API tokens were revoked")
		}
	}

	if logLevelFlag != 0 {
		oldLevel := config.GetLogLevel()
		config.SetLogLevel(logLevelFlag)
		logger.Main.Info(fmt.Sprintf("Overriding LogLevel from command line: Before=%d, Now=%d", oldLevel, logLevelFlag))
	}

	if isDebugModeFlag {
		oldDebug := config.GetIsDebugMode()
		config.SetIsDebugMode(true)
		config.SetLogLevel(10)
		logger.Main.Info(fmt.Sprintf("Overriding IsDebugMode from command line: Before=%t, Now=true", oldDebug))
	}

	if createSSUILogFileFlag {
		oldCreateSSUILogFile := config.GetCreateSSUILogFile()
		config.SetCreateSSUILogFile(true)
		logger.Main.Info(fmt.Sprintf("Overriding CreateSSUILogFile from command line: Before=%t, Now=true", oldCreateSSUILogFile))
	}
}

// HandleSanityCheckFlag has special handling to allow usage directly at startup before other systems are initialized.
func HandleSanityCheckFlag() {
	if sanityCheckFlag {
		config.NoSanityCheck = true
		logger.Main.Warn("Sanity check flag enabled, skipping sanity check. Not recommended.")
		logger.Main.Info("Sleeping for 5 seconds to remind you again to not use this flag in production.")
		time.Sleep(5 * time.Second)
	}
}
