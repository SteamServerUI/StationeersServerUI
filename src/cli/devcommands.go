package cli

import (
	"fmt"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/localization"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/modding"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/steamcmd"
)

// COMMAND HANDLERS WITH COMMANDS USEFUL FOR DEVELOPMENT AND DEBUGGING

func downloadWorkshopItemTest() {
	workshopHandles := []string{"3505169479"}
	_, err := steamcmd.DownloadWorkshopItems(workshopHandles)
	if err != nil {
		logger.Core.Error("Error downloading workshop items: " + err.Error())
	}
}

func listworkshophandles() {
	handles := modding.GetModWorkshopHandles()
	if len(handles) == 0 {
		logger.Core.Info("No mods with Workshop handles found.")
		return
	}
	logger.Core.Info(fmt.Sprintf("Installed Mod Workshop Handles: (%d):", len(handles)))
	logger.Modding.Info(fmt.Sprintf("%v", handles))

}

func listmods() {
	mods := modding.GetModList()
	if len(mods) == 0 {
		logger.Core.Info("No mods installed.")
		return
	}
	logger.Core.Info(fmt.Sprintf("Installed Mods (%d):", len(mods)))
	for _, mod := range mods {

		// print mod details in one logger call but with /n for new lines
		logger.Modding.Info("Mod Details:\n" +
			fmt.Sprintf("Modname: %s\n", mod.Name) +
			fmt.Sprintf("  Version: %s\n", mod.Version) +
			fmt.Sprintf("  Author: %s\n", mod.Author) +
			fmt.Sprintf("  Workshop Handle: %s\n", mod.WorkshopHandle))
	}
}

func testLocalization() {
	currentLanguageSetting := config.GetLanguageSetting()
	s := localization.GetString("UIText_StartButton")
	logger.Core.Info("Start Server Button text (current language: " + currentLanguageSetting + "): " + s)
}
