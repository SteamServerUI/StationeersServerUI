package runfile

import (
	"fmt"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
)

// unused
func Examples() {
	// Load gameTemplate into the global CurrentRunfile
	err := LoadRunfile("Stationeers", config.GetRunFilesFolder())
	if err != nil {
		panic(err)
	}

	// Discover all available arguments
	allArgs := GetAllArgs()
	fmt.Println("Available arguments:")
	for _, arg := range allArgs {
		fmt.Printf("%s (%s): %s\n", arg.Flag, arg.Type, arg.Description)
	}

	// Get arguments for UI display
	groups := GetUIGroups()
	for _, group := range groups {
		fmt.Printf("\n%s Settings:\n", group)
		for _, arg := range GetArgsByGroup(group) {
			fmt.Printf("- %s: %s (Current: %s)\n", arg.UILabel, arg.Description, arg.RuntimeValue)
		}
	}

	// Update a parameter
	if err := SetArgValue("GamePort", "28015"); err != nil {
		panic(err)
	}

	// Build final command line
	args, err := BuildCommandArgs()
	if err != nil {
		panic(err)
	}
	fmt.Println("\nCommand line:", args)

	arg, err := GetSingleArg("someflag")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(arg.RuntimeValue) // Access the argument's properties
}

func TestArgBuilder() {
	args, err := BuildCommandArgs()
	if err != nil {
		panic(err)
	}
	fmt.Println(args)
}
