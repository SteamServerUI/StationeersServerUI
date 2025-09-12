package loader

import (
	"flag"
	"fmt"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/security"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// LoadCmdArgs parses command-line arguments ONCE at startup (called from func main) and applies them using the config setters.
// Because this is using the config rather than adding features to it, it is a part of the loader package.
func LoadCmdArgs() {
	// Define flags matching the config variable names
	var backendEndpointPort string
	var gameBranch string
	var logLevel int
	var isDebugMode bool
	var createSSUILogFile bool
	var recoveryPassword string
	var devMode bool

	flag.StringVar(&backendEndpointPort, "BackendEndpointPort", "", "Override the backend endpoint port (e.g., 8080)")
	flag.StringVar(&backendEndpointPort, "p", "", "(Alias) Override the backend endpoint port (e.g., 8080)")
	flag.StringVar(&gameBranch, "GameBranch", "", "Override the game branch (e.g., beta)")
	flag.StringVar(&gameBranch, "b", "", "(Alias) Override the game branch (e.g., beta)")
	flag.StringVar(&recoveryPassword, "RecoveryPassword", "", "Adds a 'recovery' user (expects password as argument)")
	flag.StringVar(&recoveryPassword, "r", "", "(Alias) Adds a 'recovery' user (expects password as argument)")
	flag.BoolVar(&devMode, "dev", false, "Enable dev mode: Auth, and enables cli-console. For development only.")
	flag.IntVar(&logLevel, "LogLevel", 0, "Override the log level (e.g., 10)")
	flag.IntVar(&logLevel, "ll", 0, "(Alias) Override the log level (e.g., 10)")
	flag.BoolVar(&isDebugMode, "IsDebugMode", false, "Enable debug mode")
	flag.BoolVar(&isDebugMode, "debug", false, "(Alias) Enable debug mode")
	flag.BoolVar(&createSSUILogFile, "CreateSSUILogFile", false, "Create a log file for SSUI")
	flag.BoolVar(&createSSUILogFile, "lf", false, "(Alias) Create a log file for SSUI")

	// Parse command-line flags
	flag.Parse()

	if devMode {
		config.SetAuthEnabled(true)
		config.SetIsFirstTimeSetup(false)
		config.SetUsers(map[string]string{"admin": "$2a$10$7QQhPkNAfT.MXhJhnnodXOyn3KKE/1eu7nYb0y2O1UBoAWc0Y/fda"}) // admin:admin
		config.SetIsConsoleEnabled(true)
		logger.Main.Info("Dev mode enabled: Auth enabled, admin user set to admin:admin:superadmin, console enabled")
	}

	if backendEndpointPort != "" && backendEndpointPort != "8443" {
		oldPort := config.GetSSUIWebPort()
		config.SetSSUIWebPort(backendEndpointPort)
		logger.Main.Info(fmt.Sprintf("Overriding SetSSUIWebPort from command line: Before=%s, Now=%s", oldPort, backendEndpointPort))
	}

	if gameBranch != "" {
		oldBranch := config.GetGameBranch()
		config.SetGameBranch(gameBranch)
		logger.Main.Info(fmt.Sprintf("Overriding GameBranch from command line: Before=%s, Now=%s", oldBranch, gameBranch))
	}

	if recoveryPassword != "" {
		recoveryPassword = strings.TrimSpace(recoveryPassword)
		if recoveryPassword == "" {
			logger.Main.Error("Recovery flag provided but password is empty. Skipping recovery user creation.")
		} else {
			hashedPassword, err := security.HashPassword(recoveryPassword)
			if err != nil {
				logger.Main.Error(fmt.Sprintf("Failed to hash recovery password: %v", err))
				return
			}
			config.SetUsers(map[string]string{"recovery": hashedPassword})
			logger.Main.Warn(fmt.Sprintf("Recovery user added with access level superadmin. Login with username 'recovery' and password '%s'", recoveryPassword))
		}
	}

	if logLevel != 0 {
		oldLevel := config.GetLogLevel()
		config.SetLogLevel(logLevel)
		logger.Main.Info(fmt.Sprintf("Overriding LogLevel from command line: Before=%d, Now=%d", oldLevel, logLevel))
	}

	if isDebugMode {
		oldDebug := config.GetIsDebugMode()
		config.SetIsDebugMode(true)
		config.SetLogLevel(10)
		logger.Main.Info(fmt.Sprintf("Overriding IsDebugMode from command line: Before=%t, Now=true", oldDebug))
	}

	if createSSUILogFile {
		oldCreateSSUILogFile := config.GetCreateSSUILogFile()
		config.SetCreateSSUILogFile(true)
		logger.Main.Info(fmt.Sprintf("Overriding CreateSSUILogFile from command line: Before=%t, Now=true", oldCreateSSUILogFile))
	}
}
