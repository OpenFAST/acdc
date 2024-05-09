package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Config struct {
	RecentProjects []string `json:"RecentProjects"`
	Version        string   `json:"Version"`
}

func NewConfig() Config {
	return Config{
		RecentProjects: []string{},
		Version:        version,
	}
}

var ConfigPath = filepath.Join(xdg.ConfigHome, "acdc", "config.json")
var ConfigDir = filepath.Dir(ConfigPath)

func (app *App) LoadConfig() (Config, error) {

	// Create new config structure
	c := NewConfig()

	// Read Config File
	bs, err := os.ReadFile(ConfigPath)

	// Set version
	c.Version = version

	// If file doesn't exist, create it and return; otherwise, return error
	if os.IsNotExist(err) {
		if err := app.SaveConfig(c); err != nil {
			return c, fmt.Errorf("error creating '%s': %w", ConfigPath, err)
		}
		return c, nil
	} else if err != nil {
		return c, fmt.Errorf("error reading '%s': %w", ConfigPath, err)
	}

	// Read file into structure
	if err := json.Unmarshal(bs, &c); err != nil {
		return c, fmt.Errorf("error parsing '%s': %w", ConfigPath, err)
	}

	return c, nil
}

// SaveConfig saves the config file
func (app *App) SaveConfig(c Config) error {

	// Convert config into JSON
	bs, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}

	// Create config file directory
	if err := os.MkdirAll(ConfigDir, 0777); err != nil {
		return fmt.Errorf("error creating config dir '%s': %w", ConfigDir, err)
	}

	// Write config file
	if err := os.WriteFile(ConfigPath, bs, 0777); err != nil {
		return fmt.Errorf("error writing config file '%s': %w", ConfigPath, err)
	}
	return nil
}
