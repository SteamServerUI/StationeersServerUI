package cli

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/loader"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/managers/gamemgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/setup/update"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/steamcmd"
)

// init registers default cli commands and their aliases.
func init() {
	RegisterCommand("help", helpCommand, "h")
	RegisterCommand("reloadbackend", WrapNoReturn(loader.ReloadBackend), "rlb", "rb", "r")
	RegisterCommand("reloadconfig", WrapNoReturn(loader.ReloadConfig), "rlc", "rc")
	RegisterCommand("restartbackend", WrapNoReturn(loader.RestartBackend), "rsb")
	RegisterCommand("exit", WrapNoReturn(exitfromcli), "e")
	RegisterCommand("deleteconfig", WrapNoReturn(deleteConfig), "delc", "dc")
	RegisterCommand("startserver", WrapNoReturn(startServer), "start")
	RegisterCommand("stopserver", WrapNoReturn(stopServer), "stop")
	RegisterCommand("runsteamcmd", WrapNoReturn(runSteamCMD), "steamcmd", "stcmd")
	RegisterCommand("testlocalization", WrapNoReturn(testLocalization), "tl")
	RegisterCommand("supportmode", WrapNoReturn(supportMode), "sm")
	RegisterCommand("supportpackage", WrapNoReturn(supportPackage), "sp")
	RegisterCommand("getbuildid", WrapNoReturn(getBuildID), "gbid")
	RegisterCommand("printconfig", WrapNoReturn(printConfig), "pc")
	RegisterCommand("update", WrapNoReturn(triggerUpdateCheck), "u")
	RegisterCommand("applyupdate", WrapNoReturn(applyUpdate), "au")
	RegisterCommand("listmods", WrapNoReturn(listmods), "lm")
	RegisterCommand("listworkshophandles", WrapNoReturn(listworkshophandles), "lwh")
	RegisterCommand("downloadworkshopupdates", WrapNoReturn(downloadWorkshopUpdates), "dwu")
	RegisterCommand("downloadworkshopitemtest", WrapNoReturn(downloadWorkshopItemTest), "dwmodcon")
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

	if isSupportMode {
		config.SetIsDebugMode(false)
		config.SetLogLevel(20)
		config.SetCreateSSUILogFile(false)
		isSupportMode = false
		logger.Core.Info("Support mode disabled.")
		return
	}
	config.SetIsDebugMode(true)
	config.SetLogLevel(10)
	config.SetCreateSSUILogFile(true)
	isSupportMode = true
	loader.ReloadBackend()
	time.Sleep(1000 * time.Millisecond)
	logger.Core.Info("Support mode enabled. To generate a support package, type 'supportpackage' or 'sp'.")
}

func supportPackage() {
	if !isSupportMode {
		logger.Core.Error("Support mode is not enabled.")
		return
	}
	zipFileName := fmt.Sprintf("support_package_%s.zip", time.Now().Format("20060102_150405"))
	zipFile, _ := os.Create(zipFileName)
	defer zipFile.Close()
	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	filepath.Walk("./UIMod/logs", func(p string, i os.FileInfo, err error) error {
		if err != nil || i.IsDir() {
			return nil
		}
		f, _ := os.Open(p)
		defer f.Close()
		w, _ := zw.Create(strings.TrimPrefix(p, "./"))
		io.Copy(w, f)
		return nil
	})

	configData, _ := os.ReadFile("./UIMod/config/config.json")

	var configMap map[string]interface{}
	if err := json.Unmarshal(configData, &configMap); err != nil {
		logger.Core.Error("Failed to unmarshal config.json for support package")
		return
	}
	delete(configMap, "discordToken")
	delete(configMap, "users")
	delete(configMap, "JwtKey")
	delete(configMap, "AdminPassword")
	delete(configMap, "ServerAuthSecret")
	delete(configMap, "ServerPassword")
	sanitizedConfig, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		logger.Core.Error("Failed to marshal sanitized config into support package")
		return
	}

	// Write sanitized config to zip
	w, _ := zw.Create("UIMod/config/config.json")

	if _, err := w.Write(sanitizedConfig); err != nil {
		logger.Core.Error("Failed to write sanitized config to support package")
	}

	// Gather system information
	var osVersion string
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "ver")
		output, _ := cmd.Output()
		osVersion = strings.TrimSpace(string(output))
	} else if runtime.GOOS == "linux" {
		d, _ := os.ReadFile("/etc/os-release")
		for _, l := range strings.Split(string(d), "\n") {
			if strings.HasPrefix(l, "PRETTY_NAME=") {
				osVersion = strings.TrimPrefix(l, "PRETTY_NAME=")
				break
			}
		}
	} else {
		osVersion = "unknown"
	}

	info := fmt.Sprintf("OS: %s\nVersion: %s\nArch: %s\nBranch: %s\nVersion: %s\nTime: %s",
		runtime.GOOS, osVersion, runtime.GOARCH, config.GetBranch(), config.GetVersion(), time.Now().Format(time.RFC3339))
	w, _ = zw.Create("system_info.txt")
	w.Write([]byte(info))
}
