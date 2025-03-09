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

// Config represents the application configuration
type Config struct {
	Commands map[string]Command
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
	return config, err
}
