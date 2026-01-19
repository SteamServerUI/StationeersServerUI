package launchpad

import (
	"encoding/base64"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

// ModMetadata represents a parsed mod from About.xml
type ModMetadata struct {
	Name           string
	Author         string
	Version        string
	Description    string
	WorkshopHandle string
	Images         map[string]string // filename -> base64 encoded image data
}

// aboutXML is the structure for parsing About.xml files
type aboutXML struct {
	Name           string `xml:"Name"`
	Author         string `xml:"Author"`
	Version        string `xml:"Version"`
	Description    string `xml:"Description"`
	WorkshopHandle string `xml:"WorkshopHandle"`
}

// loadModImages loads all images from the About folder and converts them to base64
func loadModImages(aboutPath string) map[string]string {
	images := make(map[string]string)

	// List all files in the About folder
	entries, err := os.ReadDir(aboutPath)
	if err != nil {
		return images // Return empty map if we can't read the directory
	}

	// Common image extensions (case-insensitive)
	imageExtensions := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".webp": true,
	}

	// Iterate over files and look for images
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories
		}

		filename := entry.Name()
		ext := strings.ToLower(filepath.Ext(filename))

		// Check if file has an image extension
		if !imageExtensions[ext] {
			continue
		}

		// Read the image file
		imagePath := filepath.Join(aboutPath, filename)
		data, err := os.ReadFile(imagePath)
		if err != nil {
			continue // Skip if we can't read the file
		}

		// Convert to base64
		encoded := base64.StdEncoding.EncodeToString(data)
		images[filename] = encoded
	}

	return images
}

// GetModList returns an array of installed mods and their details
func GetModList() []ModMetadata {
	var mods []ModMetadata

	// Check if ./mods folder exists
	modsPath := "./mods"
	info, err := os.Stat(modsPath)
	if err != nil || !info.IsDir() {
		return mods // Return empty slice if folder doesn't exist or isn't a directory
	}

	// Read all entries in the mods folder
	entries, err := os.ReadDir(modsPath)
	if err != nil {
		return mods // Return empty slice if we can't read the directory
	}

	// Iterate over each entry in the mods folder
	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip if not a directory
		}

		dirName := entry.Name()
		aboutPath := filepath.Join(modsPath, dirName, "About")
		aboutXMLPath := filepath.Join(aboutPath, "About.xml")

		// Check if About.xml exists
		if _, err := os.Stat(aboutXMLPath); os.IsNotExist(err) {
			continue // Skip if About.xml doesn't exist
		}

		// Try to parse the XML file
		data, err := os.ReadFile(aboutXMLPath)
		if err != nil {
			continue // Skip if we can't read the file
		}

		var xmlData aboutXML
		err = xml.Unmarshal(data, &xmlData)
		if err != nil {
			continue // Skip if we can't parse the XML
		}

		// Load images from the About folder
		images := loadModImages(aboutPath)

		// Use WorkshopHandle from XML if available
		workshopHandle := xmlData.WorkshopHandle

		// Build the ModMetadata struct
		mod := ModMetadata{
			Name:           xmlData.Name,
			Author:         xmlData.Author,
			Version:        xmlData.Version,
			Description:    xmlData.Description,
			WorkshopHandle: workshopHandle,
			Images:         images,
		}

		mods = append(mods, mod)
	}

	return mods
}

// GetModWorkshopHandles returns an array of workshop handles for installed mods that have one
func GetModWorkshopHandles() []string {
	var handles []string
	mods := GetModList()

	for _, mod := range mods {
		if mod.WorkshopHandle != "" {
			handles = append(handles, mod.WorkshopHandle)
		}
	}

	return handles
}
