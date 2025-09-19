package loader

import (
	"runtime"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// PrintStartupMessage prints a stylish startup message to the terminal
func printStartupMessage() {
	// Clear some space
	logger.Core.Cleanf("")
	logger.Core.Cleanf("")

	// Main ASCII art logo
	logger.Core.Cleanf("  ███████╗████████╗ █████╗ ████████╗██╗ ██████╗ ███╗   ██╗███████╗███████╗██████╗ ███████╗      ███████╗██╗   ██╗██╗")
	logger.Core.Cleanf("  ██╔════╝╚══██╔══╝██╔══██╗╚══██╔══╝██║██╔═══██╗████╗  ██║██╔════╝██╔════╝██╔══██╗██╔════╝      ██╔════╝██║   ██║██║")
	logger.Core.Cleanf("  ███████╗   ██║   ███████║   ██║   ██║██║   ██║██╔██╗ ██║█████╗  █████╗  ██████╔╝███████╗█████╗███████╗██║   ██║██║")
	logger.Core.Cleanf("  ╚════██║   ██║   ██╔══██║   ██║   ██║██║   ██║██║╚██╗██║██╔══╝  ██╔══╝  ██╔══██╗╚════██║╚════╝╚════██║██║   ██║██║")
	logger.Core.Cleanf("  ███████║   ██║   ██║  ██║   ██║   ██║╚██████╔╝██║ ╚████║███████╗███████╗██║  ██║███████║      ███████║╚██████╔╝██║")
	logger.Core.Cleanf("  ╚══════╝   ╚═╝   ╚═╝  ╚═╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝      ╚══════╝ ╚═════╝ ╚═╝")
	logger.Core.Cleanf("  ╔═══════════════════════════════════════════════════════════════════════════════════════════════════╗")
	logger.Core.Cleanf("  ║                      🎮 YOUR ONE-STOP SHOP FOR RUNNING A STATIONEERS SERVER  🎮                   ║")
	logger.Core.Cleanf("  ║  🚀 Version: %s       📅 %s       💻 Runtime: %.3s/%s                        ║",
		config.GetVersion(),
		time.Now().Format("2006-01-02 15:04"),
		runtime.GOOS,
		runtime.GOARCH)
	logger.Core.Cleanf("  ╚═══════════════════════════════════════════════════════════════════════════════════════════════════╝")

	// Web UI info
	logger.Core.Cleanf("\n  🌐 Web UI available at: https://localhost:8443 (default) or https://<server-ip>:" + config.GetSSUIWebPort())
	logger.Core.Cleanf("\n  🌐 Support available at: https://discord.gg/8n3vN92MyJ")

	// Quote
	logger.Core.Cleanf("\n  JacksonTheMaster: \"Managing game servers shouldn't be rocket science... unless it's a rocket game!\"")
}

func printFirstTimeSetupMessage() {
	// Setup guide
	logger.Core.Cleanf("  📋 GETTING STARTED:")
	logger.Core.Cleanf("  ┌─────────────────────────────────────────────────────────────────────────────────────────────┐")
	logger.Core.Cleanf("  │ • Ready, set, go! Welcome to StationeersServerUI, new User!                                 │")
	logger.Core.Cleanf("  │ • The good news: you made it here, which means you are likely ready to run your server!     │")
	logger.Core.Cleanf("  │ • If this is your first time here, no worries: SSUI is made to be easy to use.              │")
	logger.Core.Cleanf("  │ • Configure your server by visiting the WebUI!                                              │")
	logger.Core.Cleanf("  │ • Support is provided at https://discord.gg/8n3vN92MyJ                                      │")
	logger.Core.Cleanf("  │ • For more details, check the GitHub Wiki:                                                  │")
	logger.Core.Cleanf("  │ • https://github.com/SteamServerUI/StationeersServerUI/v5/wiki                              │")
	logger.Core.Cleanf("  └─────────────────────────────────────────────────────────────────────────────────────────────┘")
}
