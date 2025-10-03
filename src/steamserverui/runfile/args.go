package runfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// Package-level mutex for file operations
var runfileMutex sync.Mutex

var CurrentRunfile *RunFile

// Custom error types
type ErrRunfileNotLoaded struct{ Msg string }

func (e ErrRunfileNotLoaded) Error() string { return e.Msg }

type ErrArgNotFound struct{ Flag string }

func (e ErrArgNotFound) Error() string { return fmt.Sprintf("argument %s not found", e.Flag) }

type ErrInvalidGameName struct{ Name string }

func (e ErrInvalidGameName) Error() string {
	return fmt.Sprintf("invalid game name %q: must start with uppercase letter, no spaces, alphanumeric", e.Name)
}

type ErrUnsetIdentifier struct{ Name string }

func (e ErrUnsetIdentifier) Error() string {
	return fmt.Sprintf("undefined runfile Identifier %q: If this is a first time setup, you can safely ignore this warning and proceed to select a runfile from the Runfile Gallery on your UI", e.Name)
}

type ErrValidation struct {
	Issues []string
}

func (e ErrValidation) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Issues, "; "))
}

type GameArg struct {
	Flag          string `json:"flag"`
	Value         string `json:"value"`
	RuntimeValue  string `json:"-"`
	Required      bool   `json:"required"`
	RequiresValue bool   `json:"requires_value"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	Special       string `json:"special,omitempty"`
	UILabel       string `json:"ui_label"`
	UIGroup       string `json:"ui_group"`
	Weight        int    `json:"weight"`
	Min           int    `json:"min,omitempty"`
	Max           int    `json:"max,omitempty"`
	Disabled      bool   `json:"disabled,omitempty"`
}

type Meta struct {
	Name    string `json:"name"`    // SSUI Specific Game Identifier, must match the one in the filename.
	Version string `json:"version"` // Runfile version
}

type RunFile struct {
	Meta               Meta                 `json:"meta"`
	Architecture       string               `json:"architecture,omitempty"`
	SteamAppID         string               `json:"steam_app_id"`
	SteamLoginRequired bool                 `json:"steam_login_required,omitempty"` //unused & unsupported, will later be used in combination with some way to provide a steam login
	WindowsExecutable  string               `json:"windows_executable"`
	LinuxExecutable    string               `json:"linux_executable"`
	Args               map[string][]GameArg `json:"args"`
}

// Validate checks the RunFile state
func (rf *RunFile) Validate() error {
	var issues []string

	// Validate SteamAppID: non-empty, numeric
	if rf.SteamAppID == "" {
		issues = append(issues, "SteamAppID is required")
	} else if _, err := strconv.Atoi(rf.SteamAppID); err != nil {
		issues = append(issues, fmt.Sprintf("SteamAppID must be numeric, got %s", rf.SteamAppID))
	}

	// Validate WindowsExecutable: if non-empty, must end with .exe
	if rf.WindowsExecutable != "" && !strings.HasSuffix(strings.ToLower(rf.WindowsExecutable), ".exe") {
		issues = append(issues, fmt.Sprintf("WindowsExecutable must end with .exe, got %s", rf.WindowsExecutable))
	}

	// Validate LinuxExecutable: if non-empty, must not end with .exe
	if rf.LinuxExecutable != "" && strings.HasSuffix(strings.ToLower(rf.LinuxExecutable), ".exe") {
		issues = append(issues, fmt.Sprintf("LinuxExecutable must not end with .exe, got %s", rf.LinuxExecutable))
	}

	// Validate Meta: ensure Name is non-empty
	if rf.Meta.Name == "" {
		issues = append(issues, "Meta.Name is required")
	}

	// Validate args
	for _, arg := range rf.getAllArgs() {
		if arg.Disabled {
			continue
		}
		if arg.Required && arg.RequiresValue && arg.RuntimeValue == "" {
			issues = append(issues, fmt.Sprintf("required argument %s has no value", arg.Flag))
		}
		switch arg.Type {
		case "int":
			if arg.RuntimeValue != "" {
				if _, err := strconv.Atoi(arg.RuntimeValue); err != nil {
					issues = append(issues, fmt.Sprintf("invalid integer value for %s: %s", arg.Flag, arg.RuntimeValue))
				}
			}
		case "bool":
			if arg.RuntimeValue != "" && arg.RuntimeValue != "true" && arg.RuntimeValue != "false" {
				issues = append(issues, fmt.Sprintf("invalid boolean value for %s: %s", arg.Flag, arg.RuntimeValue))
			}
		}
	}

	if len(issues) > 0 {
		return ErrValidation{Issues: issues}
	}
	return nil
}

// getAllArgs returns all GameArgs (internal method for validation)
func (rf *RunFile) getAllArgs() []GameArg {
	var allArgs []GameArg
	for _, category := range []string{"basic", "network", "advanced"} {
		if args, exists := rf.Args[category]; exists {
			allArgs = append(allArgs, args...)
		}
	}
	return allArgs
}

// LoadRunfile loads the runfile and stores it in CurrentRunfile
func LoadRunfile(gameName, runFilesFolder string) error {
	runfileMutex.Lock()
	defer runfileMutex.Unlock()

	// check if gameName is set to empty string
	if gameName == "" {
		err := ErrUnsetIdentifier{Name: gameName}
		return err
	}

	// Edge case: empty runFilesFolder Setting
	if runFilesFolder == "" {
		err := fmt.Errorf("runFilesFolder cannot be empty")
		logger.Runfile.Error(err.Error())
		return err
	}

	// Edge case: validate gameName (uppercase first letter, no spaces, alphanumeric)
	if gameName == "" || !regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`).MatchString(gameName) {
		err := ErrInvalidGameName{Name: gameName}
		logger.Runfile.Error(err.Error())
		return err
	}

	filePath := filepath.Join(runFilesFolder, fmt.Sprintf("run%s.ssui", gameName))
	logger.Runfile.Debug(fmt.Sprintf("loading runfile: path=%s", filePath))

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to read runfile: path=%s, error=%v", filePath, err))
		return fmt.Errorf("failed to read runfile: %w", err)
	}

	var runfile RunFile
	if err := json.Unmarshal(fileData, &runfile); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to parse runfile: path=%s, error=%v", filePath, err))
		return fmt.Errorf("failed to parse runfile: %w", err)
	}

	// Check executable availability
	if _, err := runfile.GetExecutable(); err != nil {
		logger.Runfile.Debug(fmt.Sprintf("executable validation failed: error=%v", err))
		return err
	}

	// Initialize runtime values *before* validation
	for category := range runfile.Args {
		for i := range runfile.Args[category] {
			runfile.Args[category][i].RuntimeValue = runfile.Args[category][i].Value
			logger.Runfile.Debug(fmt.Sprintf("initialized arg: flag=%s, value=%s, runtime=%s",
				runfile.Args[category][i].Flag,
				runfile.Args[category][i].Value,
				runfile.Args[category][i].RuntimeValue))
		}
	}

	// Validate runfile
	if err := runfile.Validate(); err != nil {
		logger.Runfile.Error(fmt.Sprintf("runfile validation failed: path=%s, error=%v", filePath, err))
		CurrentRunfile = nil // Ensure no partial state
		return err
	}

	CurrentRunfile = &runfile
	logger.Runfile.Info(fmt.Sprintf("runfile loaded: path=%s", filePath))
	return nil
}

// SaveRunfile persists the current RunFile to disk
func SaveRunfile() error {
	runfileMutex.Lock()
	defer runfileMutex.Unlock()

	if CurrentRunfile == nil {
		err := ErrRunfileNotLoaded{Msg: "runfile not loaded"}
		logger.Runfile.Error(err.Error())
		return err
	}

	// Build filepath
	filePath := filepath.Join(config.GetRunFilesFolder(), fmt.Sprintf("run%s.ssui", config.GetRunfileIdentifier()))
	logger.Runfile.Debug(fmt.Sprintf("saving runfile: path=%s", filePath))

	// Update Value from RuntimeValue
	for category := range CurrentRunfile.Args {
		for i := range CurrentRunfile.Args[category] {
			CurrentRunfile.Args[category][i].Value = CurrentRunfile.Args[category][i].RuntimeValue
		}
	}

	// Validate state
	if err := CurrentRunfile.Validate(); err != nil {
		logger.Runfile.Error(fmt.Sprintf("runfile validation failed: path=%s, error=%v", filePath, err))
		return err
	}

	// Serialize to JSON
	data, err := json.MarshalIndent(CurrentRunfile, "", "  ")
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to serialize runfile: path=%s, error=%v", filePath, err))
		return fmt.Errorf("failed to serialize runfile: %w", err)
	}

	// Write to file with retries
	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			logger.Runfile.Warn(fmt.Sprintf("failed to write runfile: path=%s, attempt=%d, error=%v", filePath, attempt, err))
			if attempt == maxRetries {
				logger.Runfile.Error(fmt.Sprintf("failed to write runfile after %d attempts: path=%s, error=%v", maxRetries, filePath, err))
				return fmt.Errorf("failed to write runfile after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}

	logger.Runfile.Info(fmt.Sprintf("runfile saved: path=%s", filePath))
	return nil
}

// SetArgValue updates an argument's runtime value and saves the runfile
func SetArgValue(flag string, value string) error {
	if CurrentRunfile == nil {
		err := ErrRunfileNotLoaded{Msg: "runfile not loaded"}
		logger.Runfile.Error(err.Error())
		return err
	}

	for category := range CurrentRunfile.Args {
		for i := range CurrentRunfile.Args[category] {
			if CurrentRunfile.Args[category][i].Flag != flag {
				continue
			}

			// Validate value
			arg := CurrentRunfile.Args[category][i]
			switch arg.Type {
			case "int":
				if _, err := strconv.Atoi(value); err != nil {
					err := ErrValidation{Issues: []string{fmt.Sprintf("invalid integer value for %s: %s", flag, value)}}
					logger.Runfile.Error(fmt.Sprintf("validation failed: flag=%s, value=%s, error=%v", flag, value, err))
					return err
				}
			case "bool":
				if value != "true" && value != "false" {
					err := ErrValidation{Issues: []string{fmt.Sprintf("invalid boolean value for %s: %s", flag, value)}}
					logger.Runfile.Error(fmt.Sprintf("validation failed: flag=%s, value=%s, error=%v", flag, value, err))
					return err
				}
			}

			// Transactional update
			originalValue := arg.RuntimeValue // Clone state
			CurrentRunfile.Args[category][i].RuntimeValue = value
			if err := SaveRunfile(); err != nil {
				// Rollback on failure
				CurrentRunfile.Args[category][i].RuntimeValue = originalValue
				logger.Runfile.Error(fmt.Sprintf("failed to save runfile: flag=%s, value=%s, error=%v", flag, value, err))
				return fmt.Errorf("failed to save runfile: %w", err)
			}

			logger.Runfile.Debug(fmt.Sprintf("set arg: flag=%s, value=%s", flag, value))
			return nil
		}
	}

	err := ErrArgNotFound{Flag: flag}
	logger.Runfile.Error(fmt.Sprintf("arg not found: flag=%s", flag))
	return err
}

// BuildCommandArgs builds the command-line arguments
func BuildCommandArgs() ([]string, error) {
	if CurrentRunfile == nil {
		err := ErrRunfileNotLoaded{Msg: "no runfile is currently loaded"}
		logger.Runfile.Error(err.Error())
		return nil, err
	}

	// Validate before building
	if err := CurrentRunfile.Validate(); err != nil {
		logger.Runfile.Error(fmt.Sprintf("runfile validation failed: error=%v", err))
		return nil, err
	}

	var args []string
	allArgs := CurrentRunfile.getAllArgs()

	// Sort by weight (primary) and UIGroup (secondary)
	sort.Slice(allArgs, func(i, j int) bool {
		if allArgs[i].Weight != allArgs[j].Weight {
			return allArgs[i].Weight < allArgs[j].Weight
		}
		return switchCategoryWeight(allArgs[i].UIGroup) < switchCategoryWeight(allArgs[j].UIGroup)
	})

	for _, arg := range allArgs {
		if arg.Disabled || (!arg.Required && arg.RequiresValue && arg.RuntimeValue == "") { //Clear text for clarity: skip if disabled OR if it's an optional argument that needs a value but doesn't have one set
			continue
		}

		// Handle space_delimited: split and append non-empty parts
		if arg.Special == "space_delimited" {
			parts := strings.Split(arg.RuntimeValue, " ")
			for _, part := range parts {
				if part != "" {
					args = append(args, part)
				}
			}
			continue
		}

		// Append non-empty flags
		if arg.Flag != "" && arg.Special != "dont_append_flag_just_value" {
			args = append(args, arg.Flag)
		}

		// Append value if required and non-empty
		if arg.RequiresValue && arg.RuntimeValue != "" {
			args = append(args, arg.RuntimeValue)
		}
	}

	return args, nil
}

// switchCategoryWeight maps UIGroup to a weight for sorting
func switchCategoryWeight(group string) int {
	switch group {
	case "Basic":
		return 0
	case "Network":
		return 1
	case "Advanced":
		return 2
	default:
		return 3
	}
}
