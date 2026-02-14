// Package misc provides a non-blocking command-line interface for entering commands
// while allowing the application to continue its operations normally.
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/cli/dashboard"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
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
var commandDescriptions = make(map[string]string)
var commandIsDevOnly = make(map[string]bool)

// RegisterCommand adds a new command and its handler to the registry.
func RegisterCommand(name string, handler CommandFunc, desc string, isDevCommand bool, aliases ...string) {
	mu.Lock()
	defer mu.Unlock()
	commandRegistry[name] = handler
	commandDescriptions[name] = desc
	commandIsDevOnly[name] = isDevCommand
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

		// Auto-launch dashboard on interactive terminals if enabled in config
		if config.GetIsCLIDashboardEnabled() && dashboard.IsInteractiveTerminal() {
			time.Sleep(3 * time.Second) // Give other subsystems time to initialize
			logger.Core.Info("CLI Dashboard is enabled, launching...")
			time.Sleep(500 * time.Millisecond) // Small delay for log to be visible
			if err := dashboard.Run(); err != nil {
				logger.Core.Error("Dashboard exited with error: " + err.Error())
			}
			logger.Core.Info("Dashboard closed, returning to SSUICLI prompt...")
		}

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

// helpCommand displays available commands along with their aliases and descriptions.
func helpCommand(args []string) error {
	mu.Lock()
	defer mu.Unlock()
	showDev := len(args) > 0 && args[0] == "dev"

	// Collect primary commands
	primaryCommands := make([]string, 0, len(commandAliases))
	for cmd := range commandAliases {
		primaryCommands = append(primaryCommands, cmd)
	}
	sort.Strings(primaryCommands)

	// Find the longest command name for alignment
	maxLen := 0
	for _, cmd := range primaryCommands {
		if !showDev && commandIsDevOnly[cmd] {
			continue
		}
		if len(cmd) > maxLen {
			maxLen = len(cmd)
		}
	}

	logger.Core.Cleanf("")
	if showDev {
		logger.Core.Cleanf("  \033[33m═══ SSUICLI Commands (including dev) ═══\033[0m")
	} else {
		logger.Core.Cleanf("  \033[32m═══ SSUICLI Commands ═══\033[0m")
	}
	logger.Core.Cleanf("")

	for _, cmd := range primaryCommands {
		isDev := commandIsDevOnly[cmd]
		if !showDev && isDev {
			continue
		}

		desc := commandDescriptions[cmd]
		aliases := commandAliases[cmd]

		aliasStr := ""
		if len(aliases) > 0 {
			aliasStr = " \033[90m(" + strings.Join(aliases, ", ") + ")\033[0m"
		}

		if isDev {
			logger.Core.Cleanf("  \033[33m%-*s\033[0m  %s%s", maxLen, cmd, desc, aliasStr)
		} else {
			logger.Core.Cleanf("  \033[36m%-*s\033[0m  %s%s", maxLen, cmd, desc, aliasStr)
		}
	}

	logger.Core.Cleanf("")
	return nil
}
