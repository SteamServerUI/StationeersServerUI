// Package misc provides a non-blocking command-line interface for entering commands
// while allowing the application to continue its operations normally.
package cli

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
	"github.com/SteamServerUI/SteamServerUI/v7/src/localization"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamcmd"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/runfile"
)

// ANSI escape codes for green text and reset
const (
	cliPrompt = "\033[32m" + "SSUICLI" + " » " + "\033[0m"
)

var isSupportMode bool

// CommandFunc defines the signature for command handler functions.
type CommandFunc func(args []string) error

// commandRegistry holds the map of command names to their handler functions.
var commandRegistry = make(map[string]CommandFunc)
var mu sync.Mutex

var commandAliases = make(map[string][]string)

// RegisterCommand adds a new command and its handler to the registry.
func RegisterCommand(name string, handler CommandFunc, aliases ...string) {
	mu.Lock()
	defer mu.Unlock()
	commandRegistry[name] = handler
	if len(aliases) > 0 {
		commandAliases[name] = append(commandAliases[name], aliases...)
		for _, alias := range aliases {
			commandRegistry[alias] = handler
		}
	}
}

// StartConsole starts a non-blocking console input loop in a separate goroutine.
func StartConsole(wg *sync.WaitGroup) {
	if !config.GetIsConsoleEnabled() {
		logger.Core.Info("SSUICLI runtime console is disabled in config, skipping...")
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		logger.Core.Info("SSUICLI runtime console started. Type 'help' for commands.")
		time.Sleep(10 * time.Millisecond)

		for {
			fmt.Print(cliPrompt)
			os.Stdout.Sync() // Force flush the output buffer
			if !scanner.Scan() {
				break
			}
			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}
			ProcessCommand(input)
		}

		if err := scanner.Err(); err != nil {
			logger.Core.Error("SSUICLI input error:" + err.Error())
		}
		logger.Core.Info("SSUICLI runtime console stopped.")
	}()
}

// ProcessCommand parses and executes a command from the input string.
func ProcessCommand(input string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	commandName := strings.ToLower(args[0])
	args = args[1:] // Remove command name from args

	mu.Lock()
	handler, exists := commandRegistry[commandName]
	mu.Unlock()

	if !exists {
		logger.Core.Error("Unknown command:" + commandName + ". Type 'help' for available commands.")
		return
	}

	if err := handler(args); err != nil {
		logger.Core.Error("Command " + commandName + " failed:" + err.Error())
	}
}

// WrapNoReturn wraps a function with no return value to match CommandFunc.
func WrapNoReturn(fn func()) CommandFunc {
	return func(args []string) error {
		if len(args) > 0 {
			return errors.New("command does not accept arguments")
		}
		fn()
		logger.Core.Info("Runtime CLI Command executed successfully")
		return nil
	}
}

// helpCommand displays available commands along with their aliases.
func helpCommand(args []string) error {
	mu.Lock()
	defer mu.Unlock()
	logger.Core.Info("Available commands:")
	// Collect primary commands (those in commandAliases keys)
	primaryCommands := make([]string, 0, len(commandAliases))
	for cmd := range commandAliases {
		primaryCommands = append(primaryCommands, cmd)
	}
	sort.Strings(primaryCommands)
	for _, cmd := range primaryCommands {
		aliases := commandAliases[cmd]
		if len(aliases) > 0 {
			logger.Core.Info("- " + cmd + " (aliases: " + strings.Join(aliases, ", ") + ")")
		} else {
			logger.Core.Info("- %s" + cmd)
		}
	}
	return nil
}

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
	RegisterCommand("setdummybuildid", WrapNoReturn(setDummyBuildID), "sdbid")
	RegisterCommand("printconfig", WrapNoReturn(printConfig), "pc")
	RegisterCommand("testargbuilder", WrapNoReturn(TestArgBuilder), "targb")
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

func setDummyBuildID() {
	config.SetCurrentBranchBuildID("dummy")
	logger.Core.Info("Dummy build ID set")
}

func testLocalization() {
	currentLanguageSetting := config.GetLanguageSetting()
	s := localization.GetString("UIText_StartButton")
	logger.Core.Info("Start Server Button text (current language: " + currentLanguageSetting + "): " + s)
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

func TestArgBuilder() {
	runfile.TestArgBuilder()
}
