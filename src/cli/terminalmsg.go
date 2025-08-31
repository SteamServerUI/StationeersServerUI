package cli

import (
	"fmt"
	"runtime"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

// PrintStartupMessage prints a stylish startup message to the terminal
func PrintStartupMessage() {
	// Clear some space
	fmt.Println()
	fmt.Println()

	// Main ASCII art logo
	fmt.Println("  ███████╗████████╗ █████╗ ████████╗██╗ ██████╗ ███╗   ██╗███████╗███████╗██████╗ ███████╗      ███████╗██╗   ██╗██╗")
	fmt.Println("  ██╔════╝╚══██╔══╝██╔══██╗╚══██╔══╝██║██╔═══██╗████╗  ██║██╔════╝██╔════╝██╔══██╗██╔════╝      ██╔════╝██║   ██║██║")
	fmt.Println("  ███████╗   ██║   ███████║   ██║   ██║██║   ██║██╔██╗ ██║█████╗  █████╗  ██████╔╝███████╗█████╗███████╗██║   ██║██║")
	fmt.Println("  ╚════██║   ██║   ██╔══██║   ██║   ██║██║   ██║██║╚██╗██║██╔══╝  ██╔══╝  ██╔══██╗╚════██║╚════╝╚════██║██║   ██║██║")
	fmt.Println("  ███████║   ██║   ██║  ██║   ██║   ██║╚██████╔╝██║ ╚████║███████╗███████╗██║  ██║███████║      ███████║╚██████╔╝██║")
	fmt.Println("  ╚══════╝   ╚═╝   ╚═╝  ╚═╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝      ╚══════╝ ╚═════╝ ╚═╝")

	// Decorative line
	fmt.Println("  ╔═══════════════════════════════════════════════════════════════════════════════════════════════════╗")
	// Tagline
	fmt.Println("  ║                      🎮 YOUR ONE-STOP SHOP FOR RUNNING A STATIONEERS SERVER  🎮                   ║")
	// System info
	fmt.Printf("  ║  🚀 Version: %s       📅 %s       💻 Runtime: %.3s/%s                       ║\n",
		config.GetVersion(),
		time.Now().Format("2006-01-02 15:04:05"),
		runtime.GOOS,
		runtime.GOARCH)
	// Decorative line
	fmt.Println("  ╚═══════════════════════════════════════════════════════════════════════════════════════════════════╝")

	// Web UI info
	fmt.Println("\n  🌐 Web UI available at: https://localhost:8443 (default) or https://<server-ip>:" + config.GetSSUIWebPort())
	fmt.Println("\n  🌐 Support available at: https://discord.gg/8n3vN92MyJ")

	// Quote
	fmt.Println("\n  JacksonTheMaster: \"Managing game servers shouldn't be rocket science... unless it's a rocket game!\"")
}

func PrintFirstTimeSetupMessage() {
	// Setup guide
	fmt.Println("  📋 GETTING STARTED:")
	fmt.Println("  ┌─────────────────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Println("  │ • Ready, set, go! Welcome to StationeersServerUI, new User!                                 │")
	fmt.Println("  │ • The good news: you made it here, which means you are likely ready to run your server!     │")
	fmt.Println("  │ • If this is your first time here, no worries: SSUI is made to be easy to use.              │")
	fmt.Println("  │ • Configure your server by visiting the WebUI!                                              │")
	fmt.Println("  │ • Support is provided at https://discord.gg/8n3vN92MyJ                                      │")
	fmt.Println("  │ • For more details, check the GitHub Wiki:                                                  │")
	fmt.Println("  │ • https://github.com/JacksonTheMaster/StationeersServerUI/v5/wiki                           │")
	fmt.Println("  └─────────────────────────────────────────────────────────────────────────────────────────────┘")
}
