package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Command struct represents a command that can be executed
type Command struct {
	Name string
}

// Config represents the application configuration
type Config struct {
	Commands []Command
}

// LoadConfig loads the configuration from the TOML file
func LoadConfig() (Config, error) {
	var config Config
	
	// Look for config in home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}
	
	configPath := filepath.Join(home, ".config", "switcher.toml")
	
	// Check if config file exists, if not, create default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config = Config{
			Commands: []Command{
				{Name: "firefox"},
				{Name: "vscode"},
			},
		}
		return config, nil
	}
	
	// Parse the TOML file
	_, err = toml.DecodeFile(configPath, &config)
	return config, err
}
