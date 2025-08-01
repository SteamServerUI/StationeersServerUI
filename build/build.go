// build/build.go
//go:build ignore
// +build ignore

// run from root with `go run build/build.go`
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

const (
	// ANSI color codes for styling terminal output
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

func main() {
	fmt.Printf("%s=== Starting Build Pipeline ===%s\n", colorCyan, colorReset)

	// Load the config
	config.LoadConfig()
	fmt.Printf("%s✓ Configuration loaded%s\n", colorGreen, colorReset)

	// Increment the version
	newVersion := incrementVersion("./src/config/config.go")

	// Platforms to build for
	platforms := []struct {
		os   string
		arch string
	}{
		{"windows", "amd64"},
		{"linux", "amd64"},
	}

	// Clean up old executables
	cleanupOldExecutables(newVersion)

	// Build for each platform
	for _, platform := range platforms {
		fmt.Printf("%s\nBuilding for %s/%s...%s\n", colorBlue, platform.os, platform.arch, colorReset)

		// Set OS and architecture for cross-compilation
		os.Setenv("GOOS", platform.os)
		os.Setenv("GOARCH", platform.arch)

		// Prepare the output file name with the new version, branch, and platform
		var outputName string
		if config.Branch == "release" {
			outputName = fmt.Sprintf("StationeersServerControlv%s", newVersion)
		} else {
			outputName = fmt.Sprintf("StationeersServerControlv%s_%s", newVersion, config.Branch)
		}

		// Append appropriate extension based on platform
		if platform.os == "windows" {
			outputName += ".exe"
		}
		if platform.os == "linux" {
			outputName += ".x86_64"
		}

		// Output to /build
		outputPath := filepath.Join("build", outputName)

		// Run the go build command targeting server.go at root
		cmd := exec.Command("go", "build", "-ldflags=-s -w", "-gcflags=-l=4", "-o", outputPath, "server.go")

		// Capture any output or errors
		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s✗ Build failed for %s/%s:%s %s\nOutput: %s\n",
				colorRed, platform.os, platform.arch, colorReset, err, string(cmdOutput))
			log.Fatalf("Build process terminated")
		}

		fmt.Printf("%s✓ Build successful!%s Created: %s%s%s\n",
			colorGreen, colorReset, colorYellow, outputPath, colorReset)
	}
	fmt.Printf("%s\n=== Build Pipeline Completed ===%s\n", colorCyan, colorReset)
}

// incrementVersion function to increment the version in config.go
func incrementVersion(configFile string) string {
	fmt.Printf("%sUpdating version...%s\n", colorBlue, colorReset)

	// Read the content of the config.go file
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config.go: %s", err)
	}

	// Use regex to find and increment the patch version (assuming version format is x.y.z)
	versionRegex := regexp.MustCompile(`Version\s*=\s*"(\d+)\.(\d+)\.(\d+)"`)
	matches := versionRegex.FindStringSubmatch(string(content))
	if len(matches) != 4 {
		log.Fatalf("Failed to find version in config.go")
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	// Increment the patch version
	patch++
	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	// Replace the old version with the new version
	newContent := versionRegex.ReplaceAllString(string(content), fmt.Sprintf(`Version = "%s"`, newVersion))

	// Write the updated content back to config.go
	err = os.WriteFile(configFile, []byte(newContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write updated version to config.go: %s", err)
	}

	fmt.Printf("%s✓ Version updated from %s.%s.%s to %s%s\n",
		colorGreen, matches[1], matches[2], matches[3], newVersion, colorReset)
	return newVersion
}

// Modified cleanupOldExecutables to handle both Windows and Linux executables in /build
func cleanupOldExecutables(buildVersion string) {
	fmt.Printf("%s\nCleaning up old executables...%s\n", colorBlue, colorReset)

	currentVersion := buildVersion
	dir := "build"

	// Ensure build directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("%sNo build directory found, skipping cleanup%s\n", colorYellow, colorReset)
		return
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read build directory: %s", err)
	}

	deletedCount := 0
	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) == ".exe" || filepath.Ext(filename) == ".x86_64" {
			match, _ := filepath.Match("StationeersServerControl*", filename)
			if match && !strings.Contains(filename, currentVersion) {
				exePath := filepath.Join(dir, filename)
				fmt.Printf("%s- Removing: %s%s%s\n", colorMagenta, colorYellow, exePath, colorReset)

				err := os.Remove(exePath)
				if err != nil {
					fmt.Printf("%s✗ Failed to delete %s: %s%s\n", colorRed, exePath, err, colorReset)
				} else {
					fmt.Printf("%s✓ Deleted successfully%s\n", colorGreen, colorReset)
					deletedCount++
				}
			}
		}
	}

	if deletedCount == 0 {
		fmt.Printf("%sNo old executables found to clean up%s\n", colorYellow, colorReset)
	} else {
		fmt.Printf("%s✓ Cleaned up %d old executable(s)%s\n", colorGreen, deletedCount, colorReset)
	}
}
