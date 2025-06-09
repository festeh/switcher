package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Command struct represents a command that can be executed
type Command struct {
	Name string `toml:"name"`
	Run  string `toml:"run"`
	Key  string `toml:"key"`
}

// General represents general application settings
type General struct {
	BookScanPath string `toml:"book_scan_path"`
}

// Config represents the application configuration
type Config struct {
	General  General            `toml:"general"`
	Commands map[string]Command `toml:"commands"`
}

// LoadConfig loads the configuration from the TOML file
func LoadConfig() (Config, error) {
	var config Config
	config.Commands = make(map[string]Command)

	// Look for config in home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}

	configPath := filepath.Join(home, ".config", "switcher", "switcher.toml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, err
	}

	// Parse the TOML file
	_, err = toml.DecodeFile(configPath, &config)
	
	// Set default book scan path if not specified
	if config.General.BookScanPath == "" {
		config.General.BookScanPath = filepath.Join(home, "pCloudDrive")
	}
	
	return config, err
}
