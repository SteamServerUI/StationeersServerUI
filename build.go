// build.go
package main

import (
	"StationeersServerUI/src/config"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ServerConfig struct {
	ExePath  string `xml:"exePath"`
	Settings string `xml:"settings"`
}

type Config struct {
	Server       ServerConfig `xml:"server"`
	SaveFileName string       `xml:"saveFileName"`
}

func main() {

	// Load the config
	config.LoadConfig()

	// Increment the version
	newVersion := incrementVersion("src/config/config.go")

	// Platforms to build for
	platforms := []struct {
		os   string
		arch string
	}{
		{"windows", "amd64"},
		{"linux", "amd64"},
	}

	// Build for each platform
	for _, platform := range platforms {
		// Set OS and architecture for cross-compilation
		os.Setenv("GOOS", platform.os)
		os.Setenv("GOARCH", platform.arch)

		// Prepare the output file name with the new version, branch, and platform
		var outputName string
		if config.Branch == "release" {
			outputName = fmt.Sprintf("StationeersServerControl%s", newVersion)
		} else {
			outputName = fmt.Sprintf("StationeersServerControl%s_%s", newVersion, config.Branch)
		}

		// Append .exe only on Windows
		if platform.os == "windows" {
			outputName += ".exe"
		}
		if platform.os == "linux" {
			outputName += ".x86_64"
		}

		// Run the go build command with the custom output name
		cmd := exec.Command("go", "build", "-ldflags=-s -w", "-gcflags=-l=4", "-o", outputName, "./src")

		// Capture any output or errors
		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Build failed for %s-%s: %s\nOutput: %s", platform.os, platform.arch, err, string(cmdOutput))
		}

		fmt.Printf("Build successful for %s-%s! Output: %s\n", platform.os, platform.arch, outputName)
	}

	// Clean up old .exe files that follow the pattern "StationeersServerControl*"
	cleanupOldExecutables(newVersion)
}

// incrementVersion function to increment the version in config.go
func incrementVersion(configFile string) string {
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

	// Construct the new version
	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	// Replace the old version with the new version
	newContent := versionRegex.ReplaceAllString(string(content), fmt.Sprintf(`Version = "%s"`, newVersion))

	// Write the updated content back to config.go
	err = os.WriteFile(configFile, []byte(newContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write updated version to config.go: %s", err)
	}

	fmt.Printf("Version updated to %s\n", newVersion)
	return newVersion
}

// Modified cleanupOldExecutables to handle both Windows and Linux executables
func cleanupOldExecutables(buildVersion string) {
	currentVersion := buildVersion
	// Get the current directory
	dir := "."

	// Get a list of all files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %s", err)
	}

	// Loop through the files and delete old executables
	for _, file := range files {
		filename := file.Name()
		// Check for both .exe and Linux executables
		if filepath.Ext(filename) == ".exe" || filepath.Ext(filename) == ".x86_64" {
			match, _ := filepath.Match("StationeersServerControl*", filename)
			if match {
				// Skip deletion if the filename contains the current version
				if strings.Contains(filename, currentVersion) {
					continue
				}

				exePath := filepath.Join(dir, filename)
				fmt.Printf("Deleting old executable: %s\n", exePath)
				err := os.Remove(exePath)
				if err != nil {
					log.Printf("Failed to delete %s: %s", exePath, err)
				} else {
					fmt.Printf("Successfully deleted: %s\n", exePath)
				}
			}
		}
	}
}
