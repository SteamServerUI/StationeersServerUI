package launchpad

import (
	"encoding/xml"
	"os"
	"path/filepath"
)

// ModMetadata represents a parsed mod from About.xml
type ModMetadata struct {
	Name           string
	Author         string
	Version        string
	Description    string
	WorkshopHandle string
}

// aboutXML is the structure for parsing About.xml files
type aboutXML struct {
	Name           string `xml:"Name"`
	Author         string `xml:"Author"`
	Version        string `xml:"Version"`
	Description    string `xml:"Description"`
	WorkshopHandle string `xml:"WorkshopHandle"`
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
		aboutXMLPath := filepath.Join(modsPath, dirName, "About", "About.xml")

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

		// Use WorkshopHandle from XML if available
		workshopHandle := xmlData.WorkshopHandle

		// Build the ModMetadata struct
		mod := ModMetadata{
			Name:           xmlData.Name,
			Author:         xmlData.Author,
			Version:        xmlData.Version,
			Description:    xmlData.Description,
			WorkshopHandle: workshopHandle,
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
