package setup

import (
	"StationeersServerUI/src/config"
	"StationeersServerUI/src/loader"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

var downloadBranch string // Holds the branch to download from

// Install performs the entire installation process and ensures the server waits for it to complete
func Install(wg *sync.WaitGroup) {
	defer wg.Done()             // Signal that installation is complete
	time.Sleep(1 * time.Second) // Small pause for effect

	loader.ReloadConfig()

	// Step 0:  Check for updates
	if err := UpdateExecutable(); err != nil {
		fmt.Printf("❌ Update check went sideways: %v\n", err)
	}

	// Step 1: Check and download the UIMod folder contents
	fmt.Println("🔄 Checking UIMod folder contents...")
	CheckAndDownloadUIMod()
	fmt.Println("✅ UIMod folder setup complete.")
	time.Sleep(1 * time.Second)

	// Step 2: Check for Blacklist.txt and create it if it doesn't exist
	fmt.Println("🔄 Checking for Blacklist.txt...")
	checkAndCreateBlacklist()
	fmt.Println("✅ Blacklist.txt verified or created.")
	time.Sleep(1 * time.Second)

	// Step 3: Install and run SteamCMD
	fmt.Println("🔄 Installing and running SteamCMD...")
	InstallAndRunSteamCMD()
	fmt.Println("Thank you for using StationeersServerUI! 🙏")
	fmt.Println(string(colorCyan), "Setup complete!", string(colorReset))
}

func CheckAndDownloadUIMod() {
	workingDir := "./UIMod/"
	loginDir := "./UIMod/login/"
	detectionmanagerDir := "./UIMod/detectionmanager/"

	// Check if the directory exists
	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		fmt.Println("⚠️ Folder ./UIMod does not exist. Creating it...")

		// Create the UIMod folder
		err := os.MkdirAll(workingDir, os.ModePerm)
		if err != nil {
			fmt.Printf("❌ Error creating folder: %v\n", err)
			return
		}

		if _, err := os.Stat(loginDir); os.IsNotExist(err) {
			fmt.Println("⚠️ Folder ./UIMod/login/ does not exist. Creating it...")
			err := os.MkdirAll(loginDir, os.ModePerm)
			if err != nil {
				fmt.Printf("❌ Error creating folder: %v\n", err)
				return
			}
		}

		if _, err := os.Stat(detectionmanagerDir); os.IsNotExist(err) {
			fmt.Println("⚠️ Folder ./UIMod/detectionmanager/ does not exist. Creating it...")
			err := os.MkdirAll(detectionmanagerDir, os.ModePerm)
			if err != nil {
				fmt.Printf("❌ Error creating folder: %v\n", err)
				return
			}
		}

		if config.Branch == "release" || config.Branch == "Release" {
			downloadBranch = "main"
		} else {
			downloadBranch = config.Branch
		}

		// List of files to download with their destination paths
		files := map[string]string{
			workingDir + "apiinfo.html":                   fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/apiinfo.html", downloadBranch),
			workingDir + "config.html":                    fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/config.html", downloadBranch),
			workingDir + "config.json":                    fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/config.json", downloadBranch),
			workingDir + "index.html":                     fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/index.html", downloadBranch),
			workingDir + "script.js":                      fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/script.js", downloadBranch),
			workingDir + "stationeers.png":                fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/stationeers.png", downloadBranch),
			workingDir + "style.css":                      fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/style.css", downloadBranch),
			workingDir + "favicon.ico":                    fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/favicon.ico", downloadBranch),
			loginDir + "login.css":                        fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/login/login.css", downloadBranch),
			loginDir + "login.js":                         fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/login/login.js", downloadBranch),
			loginDir + "login.html":                       fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/login/login.html", downloadBranch),
			detectionmanagerDir + "detectionmanager.js":   fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/detectionmanager/detectionmanager.js", downloadBranch),
			detectionmanagerDir + "detectionmanager.html": fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/detectionmanager/detectionmanager.html", downloadBranch),
			detectionmanagerDir + "detectionmanager.css":  fmt.Sprintf("https://raw.githubusercontent.com/JacksonTheMaster/StationeersServerUI/%s/UIMod/detectionmanager/detectionmanager.css", downloadBranch),
		}

		// Set the first time setup flag to true
		config.IsFirstTimeSetup = true

		// Download each file
		for filepath, url := range files {
			// Extract just the filename for display purposes
			fileName := filepath[strings.LastIndex(filepath, "/")+1:]
			fmt.Printf("Downloading %s... ", fileName)
			err := downloadFileWithProgress(filepath, url)
			if err != nil {
				fmt.Printf("\n❌ Error downloading (setup may has been left incomplete) %s: %v\n", fileName, err)
				return
			}
			fmt.Printf("\n✅ Downloaded %s successfully from branch %s\n", fileName, downloadBranch)
		}

		fmt.Println("✅ All files downloaded successfully.")
	} else {
		fmt.Println("♻️ Folder ./UIMod already exists. Skipping download.")
		config.IsFirstTimeSetup = false
	}
}

// checkAndCreateBlacklist ensures Blacklist.txt exists in the root directory
func checkAndCreateBlacklist() {
	blacklistFile := "./Blacklist.txt"

	// Check if Blacklist.txt exists
	if _, err := os.Stat(blacklistFile); os.IsNotExist(err) {
		// Create Blacklist.txt file with a dummy steamID64 so the gameserver doesn't fail reading this file, as it would not be the expected format if it was empty.
		perm := os.FileMode(0644) // Still works cross-platform
		err := os.WriteFile(blacklistFile, []byte("76561197960265728"), perm)
		if err != nil {
			fmt.Printf("❌ Error creating Blacklist.txt: %v\n", err)
			return
		}

		fmt.Println("✅ Created Blacklist.txt with dummy steamID64.")
	} else {
		fmt.Println("♻️ Blacklist.txt already exists. Skipping creation.")
	}
}

// downloadFileWithProgress downloads a file from the given URL and saves it to the given filepath with progress indication
func downloadFileWithProgress(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Get the total size for progress reporting
	size := resp.ContentLength

	// Create a counter for tracking progress
	counter := &writeCounter{
		Total: size,
	}

	// Write the body to file with progress tracking
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	return nil
}

// writeCounter tracks download progress
type writeCounter struct {
	Total int64
	count int64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.count += int64(n)
	wc.printProgress()
	return n, nil
}

func (wc *writeCounter) printProgress() {
	// If we don't know the total size, just show downloaded bytes
	if wc.Total <= 0 {
		fmt.Printf("\r%s downloaded", bytesToHuman(wc.count))
		return
	}

	// Calculate percentage with bounds checking
	percent := float64(wc.count) / float64(wc.Total) * 100
	if percent > 100 {
		percent = 100
	}

	// Create simple progress bar
	width := 20
	complete := int(percent / 100 * float64(width))

	progressBar := "["
	for i := 0; i < width; i++ {
		if i < complete {
			progressBar += "="
		} else if i == complete && complete < width {
			progressBar += ">"
		} else {
			progressBar += " "
		}
	}
	progressBar += "]"

	// Print progress and erase to end of line
	fmt.Printf("\r%s %.1f%% (%s/%s)",
		progressBar,
		percent,
		bytesToHuman(wc.count),
		bytesToHuman(wc.Total))
}

// bytesToHuman converts bytes to human readable format
func bytesToHuman(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
