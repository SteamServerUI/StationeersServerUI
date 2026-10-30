package cli

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/core/loader"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/diagnostics"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/logger"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/managers/gamemgr"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/setup/update"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/steamcmd"
)

// init registers default cli commands and their aliases.
func init() {

	//Long command, Description, IsDevCommand, Aliases

	// User commands
	RegisterCommand("help", helpCommand, "Show available commands", false, "h")
	RegisterCommand("startserver", WrapNoReturn(startServer), "Start the game server", false, "start")
	RegisterCommand("stopserver", WrapNoReturn(stopServer), "Stop the game server", false, "stop")
	RegisterCommand("update", WrapNoReturn(triggerUpdateCheck), "Trigger an SSUI update check", false, "u")
	RegisterCommand("applyupdate", WrapNoReturn(applyUpdate), "Apply available SSUI updates", false, "au")
	RegisterCommand("reloadbackend", WrapNoReturn(loader.ReloadBackend), "Reload the SSUI backend", false, "rlb", "rb", "r")
	RegisterCommand("reloadconfig", WrapNoReturn(loader.ReloadConfig), "Reload the SSUI configuration", false, "rlc", "rc")
	RegisterCommand("restartbackend", WrapNoReturn(loader.RestartBackend), "Restart the SSUI backend", false, "rsb")
	RegisterCommand("supportmode", WrapNoReturn(supportMode), "Start or stop a support recording", false, "sm")
	RegisterCommand("supportpackage", WrapNoReturn(supportPackage), "Export a support package (with logs if recorded)", false, "sp")
	RegisterCommand("exit", WrapNoReturn(exitfromcli), "Exit / Shutdown SSUI", false, "e")
	RegisterCommand("runsteamcmd", WrapNoReturn(runSteamCMD), "Run SteamCMD to update the Gameserver", false, "steamcmd", "stcmd")
	RegisterCommand("downloadworkshopupdates", WrapNoReturn(downloadWorkshopUpdates), "Download Steam workshop Mod updates", true, "dwu")

	// Dev commands
	RegisterCommand("deleteconfig", WrapNoReturn(deleteConfig), "DELETE the SSUI configuration file", true, "delc", "dc")
	RegisterCommand("testlocalization", WrapNoReturn(testLocalization), "Test localization strings", true, "tl")
	RegisterCommand("getbuildid", WrapNoReturn(getBuildID), "Get the current game build ID (server must have started once, else empty)", true, "gbid")
	RegisterCommand("printconfig", WrapNoReturn(printConfig), "Print the current SSUI configuration", true, "pc")
	RegisterCommand("listmods", WrapNoReturn(listmods), "List installed SLP mods", true, "lm")
	RegisterCommand("listworkshophandles", WrapNoReturn(listworkshophandles), "List workshop Mod handles", true, "lwh")
	RegisterCommand("downloadworkshopitem", downloadWorkshopItem, "Download a workshop item. When no arguments are provided, 3672138641/BlueprintMod is downloaded. Can be called like downloadworkshopitem 3505169479 or downloadworkshopitem 3505169479 3505115682 3505169479 ", true, "dwi")
	RegisterCommand("dumpheapprofile", WrapNoReturn(dumpHeapProfile), "Dump a pprof heap profile for debugging", true, "dhp")
	RegisterCommand("testserverstatuspaneldiscord", WrapNoReturn(testServerStatusPanelDiscord), "Send a fake player list to the Discord package to test the server status panel", true, "tsspd")
}

// COMMAND HANDLERS WITH COMMANDS USEFUL FOR USERS

func downloadWorkshopUpdates() {
	_, err := steamcmd.UpdateWorkshopItems()
	if err != nil {
		logger.Core.Error("Error downloading workshop updates: " + err.Error())
	}
}

func startServer() {
	err := gamemgr.InternalStartServer()
	if err != nil {
		logger.Core.Error("Error starting server:" + err.Error())
	}
}
func stopServer() {
	err := gamemgr.InternalStopServer()
	if err != nil {
		logger.Core.Error("Error stopping server:" + err.Error())
	}
}

func exitfromcli() {
	// send signal to the main process to exit
	logger.Core.Info("I have to go...")
	err := gamemgr.InternalStopServer()
	if err != nil {
		logger.Core.Error("Error stopping server:" + err.Error())
	}
	os.Exit(0)
}

func deleteConfig() {
	//remove file at config.ConfigPath
	if err := os.Remove(config.GetConfigPath()); err != nil {
		logger.Core.Error("Error deleting config file: " + err.Error())
		return
	}
	logger.Core.Info("Config file deleted successfully")
}

func runSteamCMD() {
	steamcmd.InstallAndRunSteamCMD()
}

func printConfig() {
	loader.PrintConfigDetails("Info")
}

func getBuildID() {
	buildID := config.GetCurrentBranchBuildID()
	if buildID == "" {
		logger.Core.Error("Build ID not found, empty string returned")
		return
	}
	logger.Core.Info("Build ID: " + buildID)
}

func triggerUpdateCheck() {
	err, newVersion := update.Update(false)
	if err != nil {
		logger.Install.Warn("⚠️ Update check failed: " + err.Error())
		return
	}
	if newVersion != "" {
		logger.Install.Infof("✅ Update to %s available, Trigger update from WebUI or with applyupdate command", newVersion)
	}
}

func applyUpdate() {
	err, _ := update.Update(true)
	if err != nil {
		logger.Install.Warn("⚠️ Update failed: " + err.Error())
		return
	}
}

func supportMode() {
	if logger.SupportRecordingStatus().Active {
		if _, err := diagnostics.StopRecording(); err != nil {
			logger.Core.Warn("Support recording stopped with an error: " + err.Error())
		}
		logger.Core.Info("Support recording stopped. Type 'supportpackage' or 'sp' to export it.")
		return
	}
	if _, err := logger.StartSupportRecording(); err != nil {
		logger.Core.Error("Could not start support recording: " + err.Error())
		return
	}
	logger.Core.Info("Support recording started. Reproduce the issue, then type 'supportpackage' or 'sp'.")
}

func supportPackage() {
	if logger.SupportRecordingStatus().Active {
		supportMode()
	}
	write, name := diagnostics.WriteRecordingPackage, "support_recording"
	if logger.SupportRecordingStatus().StoppedAt == nil {
		write, name = diagnostics.WritePackage, "support_package"
	}
	path := fmt.Sprintf("%s_%s.zip", name, time.Now().Format("20060102_150405"))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		logger.Core.Error("Could not create support package: " + err.Error())
		return
	}
	err = errors.Join(write(file), file.Close())
	if err != nil {
		os.Remove(path)
		logger.Core.Error("Could not write support package: " + err.Error())
		return
	}
	logger.Core.Info("Support package written to " + path + ". Review it before sharing.")
}
