// processmanagement.go
package api

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/ssestream"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var (
	cmd *exec.Cmd
	mu  sync.Mutex
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

type Arg struct {
	Flag          string
	Value         string
	Condition     func() bool
	RequiresValue bool
}

// print args for debugging
func printArgs(args []string) {
	for _, arg := range args {
		fmt.Println(arg)
	}
}

func buildCommandArgs() []string {

	// Define argument order here
	var argOrder = []Arg{
		{Flag: "-nographics", RequiresValue: false},
		{Flag: "-batchmode", RequiresValue: false},
		{Flag: "-LOAD", Value: config.SaveFileName, RequiresValue: true},
		{Flag: "-logFile", Value: `"./debug.log"`, Condition: func() bool { return runtime.GOOS == "linux" }, RequiresValue: true}, // Attach a logfile on Linux, since piped output is not available
		{Flag: "-settings", RequiresValue: false},
		{Flag: "StartLocalHost", Value: strconv.FormatBool(config.StartLocalHost), RequiresValue: true},
		{Flag: "ServerVisible", Value: strconv.FormatBool(config.ServerVisible), RequiresValue: true},
		{Flag: "GamePort", Value: config.GamePort, RequiresValue: true},
		{Flag: "UPNPEnabled", Value: strconv.FormatBool(config.UPNPEnabled), RequiresValue: true},
		{Flag: "ServerName", Value: config.ServerName, RequiresValue: true},
		{Flag: "ServerPassword", Value: config.ServerPassword, Condition: func() bool { return config.ServerPassword != "" }, RequiresValue: true},
		{Flag: "ServerMaxPlayers", Value: config.ServerMaxPlayers, RequiresValue: true},
		{Flag: "AutoSave", Value: strconv.FormatBool(config.AutoSave), RequiresValue: true},
		{Flag: "SaveInterval", Value: config.SaveInterval, RequiresValue: true},
		{Flag: "ServerAuthSecret", Value: config.ServerAuthSecret, Condition: func() bool { return config.ServerAuthSecret != "" }, RequiresValue: true},
		{Flag: "UpdatePort", Value: config.UpdatePort, RequiresValue: true},
		{Flag: "AutoPauseServer", Value: strconv.FormatBool(config.AutoPauseServer), RequiresValue: true},
		{Flag: "UseSteamP2P", Value: strconv.FormatBool(config.UseSteamP2P), RequiresValue: true},
		{Flag: "AdminPassword", Value: config.AdminPassword, Condition: func() bool { return config.AdminPassword != "" }, RequiresValue: true},
		{Flag: "LocalIpAddress", Value: config.LocalIpAddress, RequiresValue: true},
	}
	var args []string
	for _, arg := range argOrder {
		// Skip if condition exists and fails
		if arg.Condition != nil && !arg.Condition() {
			continue
		}

		// If the flag requires a value and the value is empty, skip it entirely
		if arg.RequiresValue && arg.Value == "" {
			continue
		}

		// Add the flag
		args = append(args, arg.Flag)

		// Add the value if it exists
		if arg.Value != "" {
			// If the value contains a space and isn’t already quoted, wrap it in quotes
			if strings.Contains(arg.Value, " ") && !strings.HasPrefix(arg.Value, `"`) && !strings.HasSuffix(arg.Value, `"`) {
				args = append(args, `"`+arg.Value+`"`)
			} else {
				args = append(args, arg.Value)
			}
		}
	}

	if config.AdditionalParams != "" {
		extraArgs := strings.Fields(config.AdditionalParams)
		for _, extraArg := range extraArgs {
			if strings.Contains(extraArg, " ") {
				args = append(args, `"`+extraArg+`"`)
			} else {
				args = append(args, extraArg)
			}
		}
	}
	if config.IsDebugMode {
		printArgs(args)
	}
	return args
}

func StartServer(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if cmd != nil && cmd.Process != nil {
		fmt.Fprintf(w, "Server is already running.")
		return
	}
	if config.IsDebugMode {
		fmt.Println("Config values:", config.UPNPEnabled, config.StartLocalHost, config.ServerVisible, config.AutoSave, config.AutoPauseServer, config.UseSteamP2P)
	}
	args := buildCommandArgs()
	cmd = exec.Command(config.ExePath, args...)

	fmt.Printf("\n%s%s=== GAMESERVER STARTING ===%s\n", colorCyan, colorBold, colorReset)
	fmt.Printf("• Executable: %s\n", colorGreen+colorBold+config.ExePath+colorReset)
	fmt.Printf("• Parameters: %s\n", colorYellow+strings.Join(args, " ")+colorReset)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(w, "Error creating StdoutPipe: %v", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(w, "Error creating StderrPipe: %v", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(w, "Error starting server: %v", err)
		return
	}

	// Start reading stdout and stderr
	go readPipe(stdout)
	go readPipe(stderr)

	fmt.Fprintf(w, "Server started.")
}

func readPipe(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	if config.IsDebugMode {
		fmt.Println("Started reading pipe") // Debug
	}
	for scanner.Scan() {
		output := scanner.Text()
		if config.IsDebugMode {
			fmt.Println("Pipe output:", output) // Debug
		}
		ssestream.BroadcastConsoleOutput(output)
	}
	if err := scanner.Err(); err != nil {
		if config.IsDebugMode {
			fmt.Println("Pipe error:", err) // Debug
		}
		ssestream.BroadcastConsoleOutput(fmt.Sprintf("Error reading pipe: %v", err))
	}
	if config.IsDebugMode {
		fmt.Println("Pipe closed") // Debug
	}
}

func GetLogOutput(w http.ResponseWriter, r *http.Request) {
	ssestream.StartConsoleStream()(w, r)
}

func StopServer(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if cmd == nil || cmd.Process == nil {
		fmt.Fprintf(w, "Server is not running.")
		return
	}

	isWindows := runtime.GOOS == "windows"

	if isWindows {
		// On Windows, just kill the process directly
		if killErr := cmd.Process.Kill(); killErr != nil {
			fmt.Fprintf(w, "Error stopping server: %v", killErr)
			return
		}
	} else {
		// On Linux/Unix, try SIGTERM first for graceful shutdown
		if termErr := cmd.Process.Signal(syscall.SIGTERM); termErr != nil {
			// If SIGTERM fails, fall back to Kill
			if killErr := cmd.Process.Kill(); killErr != nil {
				fmt.Fprintf(w, "Error stopping server: %v", killErr)
				return
			}
		}
	}

	// Wait for the process to exit
	if waitErr := cmd.Wait(); waitErr != nil && !strings.Contains(waitErr.Error(), "exit status") {
		// Only report actual errors, not just non-zero exit codes
		fmt.Fprintf(w, "Error during server shutdown: %v", waitErr)
		return
	}

	cmd = nil
	fmt.Fprintf(w, "Server stopped.")
}
