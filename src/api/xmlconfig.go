package api

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// LoadConfig loads the configuration from an XML file
func uloadConfig() (*Config, error) {
	configPath := "./UIMod/config.xml"
	xmlFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = xml.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %v", err)
	}

	return &config, nil
}
