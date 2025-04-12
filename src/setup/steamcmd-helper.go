package setup

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// createSteamCMDDirectory creates the SteamCMD directory.
func createSteamCMDDirectory(steamCMDDir string) error {
	if err := os.MkdirAll(steamCMDDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create SteamCMD directory: %w", err)
	}
	logger.Install.Debug("✅ Created SteamCMD directory: " + steamCMDDir + "\n")
	return nil
}

// downloadAndExtractSteamCMD downloads and extracts SteamCMD.
func downloadAndExtractSteamCMD(downloadURL string, steamCMDDir string, extractFunc ExtractorFunc) error {
	// Validate download URL
	if err := validateURL(downloadURL); err != nil {
		return fmt.Errorf("invalid download URL: %w", err)
	}
	logger.Install.Debug("✅ Validated download URL: " + downloadURL + "\n")

	// Download SteamCMD with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	logger.Install.Debug("✅ Created HTTP request for download.\n")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading SteamCMD: %w", err)
	}
	defer resp.Body.Close()
	logger.Install.Debug("✅ Successfully downloaded SteamCMD.\n")

	// Check for successful HTTP response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download SteamCMD: HTTP status %v", resp.Status)
	}

	logger.Install.Debug("✅ Received HTTP status: " + resp.Status + "\n")

	// Read the downloaded content into memory
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading SteamCMD content: %w", err)
	}
	logger.Install.Debug("✅ Read SteamCMD content into memory.\n")

	// Create a reader for the content
	contentReader := bytes.NewReader(content)

	// Extract the content using the provided extractor function
	if err := extractFunc(contentReader, int64(len(content)), steamCMDDir); err != nil {
		return fmt.Errorf("error extracting SteamCMD: %w", err)
	}
	logger.Install.Debug("✅ Successfully extracted SteamCMD.\n")

	return nil
}

// setExecutablePermissions sets executable permissions for SteamCMD files.
func setExecutablePermissions(steamCMDDir string) error {
	if runtime.GOOS == "windows" {
		logger.Install.Debug("✅ Skipping executable permissions on Windows.\n")
		return nil
	}
	// List of files that need executable permissions
	files := []string{
		filepath.Join(steamCMDDir, "steamcmd.sh"),
		filepath.Join(steamCMDDir, "linux32", "steamcmd"),
		filepath.Join(steamCMDDir, "linux32", "steamerrorreporter"),
	}

	for _, file := range files {
		if err := os.Chmod(file, 0755); err != nil {
			return fmt.Errorf("failed to set executable permissions for %s: %w", file, err)
		}
		logger.Install.Debug("✅ Set executable permissions for: " + file + "\n")
	}

	return nil
}

// verifySteamCMDBinary verifies that the steamcmd binary exists.
func verifySteamCMDBinary(steamCMDDir string) error {
	// Different binary paths based on OS
	var binaryPath string
	if runtime.GOOS == "windows" {
		binaryPath = filepath.Join(steamCMDDir, "steamcmd.exe")
	} else {
		binaryPath = filepath.Join(steamCMDDir, "linux32", "steamcmd")
	}

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("steamcmd binary not found: %s", binaryPath)
	}
	logger.Install.Debug("✅ Verified steamcmd binary: " + binaryPath + "\n")
	return nil
}

// validateURL checks if a URL is valid.
func validateURL(rawURL string) error {
	_, err := url.ParseRequestURI(rawURL)
	return err
}

// untar extracts a tar.gz archive.
func untar(dest string, r io.Reader) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		// Ensure the parent directory exists
		parentDir := filepath.Dir(target)
		if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create parent directory %s: %v", parentDir, err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", target, err)
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %v", target, err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return fmt.Errorf("failed to write file %s: %v", target, err)
			}
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink %s: %v", target, err)
			}
		default:
			return fmt.Errorf("unknown type: %v in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

// unzip extracts a zip archive.
func unzip(zipReader io.ReaderAt, size int64, dest string) error {
	reader, err := zip.NewReader(zipReader, size)
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Ensure the destination directory exists
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, f := range reader.File {
		// Sanitize the file path to prevent path traversal
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			// Create directory with the same permissions as in the zip file
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Create the file with the same permissions as in the zip file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer outFile.Close()

		// Open the file in the zip archive
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}
		defer rc.Close()

		// Copy the file contents using a buffer for better performance
		buffer := make([]byte, 32*1024) // 32KB buffer
		if _, err := io.CopyBuffer(outFile, rc, buffer); err != nil {
			return fmt.Errorf("failed to copy file contents: %w", err)
		}
	}

	return nil
}

// untarWrapper adapts the untar function to match the ExtractorFunc signature.
func untarWrapper(r io.ReaderAt, _ int64, dest string) error {
	return untar(dest, io.NewSectionReader(r, 0, 1<<63-1)) // Use a large size for the section reader
}

// installRequiredLibraries installs the required libraries for SteamCMD if they are not already installed.
func installRequiredLibraries() error {
	// Check if the system is Debian-based
	if runtime.GOOS != "linux" {
		return nil // Only Linux systems need this
	}

	// Check if running inside a Docker container
	if _, err := os.Stat("/.dockerenv"); err == nil {
		logger.Install.Debug("✅ Running inside a Docker container, skipping library installation.\n")
		return nil
	}

	// According to https://developer.valvesoftware.com/wiki/SteamCMD#Manually only lib32gcc-s1 is needed
	// List of required libraries
	requiredLibs := []string{
		"lib32gcc-s1",
		//"lib32stdc++6",
	}

	// Check and install each library
	for _, lib := range requiredLibs {
		// Check if the library is already installed
		cmd := exec.Command("dpkg", "-s", lib)
		if err := cmd.Run(); err == nil {
			logger.Install.Debug("✅ Library already installed: " + lib + "\n")
			continue // Library is already installed, skip to the next one
		}

		// Library is not installed, attempt to install it
		logger.Install.Debug("🔄 Installing library: " + lib + "\n")
		installCmd := exec.Command("sudo", "apt-get", "install", "-y", lib)
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install library %s: %w", lib, err)
		}
		logger.Install.Debug("✅ Installed library: " + lib + "\n")
	}

	return nil
}
