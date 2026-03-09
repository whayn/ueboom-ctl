package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config stores speaker and host information.
type Config struct {
	TargetMAC string `json:"target_mac"`
	HostMAC   string `json:"host_mac"`
}

const (
	SystemConfigPath = "/etc/ueboom/config.json"
)

// UserConfigPath returns the path to the user configuration file.
func UserConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ueboom", "config.json"), nil
}

// Load reads the configuration from system or user paths.
func Load() (*Config, string, error) {
	if cfg, err := readFromFile(SystemConfigPath); err == nil {
		return cfg, SystemConfigPath, nil
	}

	userPath, err := UserConfigPath()
	if err != nil {
		return nil, "", err
	}

	cfg, err := readFromFile(userPath)
	return cfg, userPath, err
}

func readFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the configuration to a file.
func (c *Config) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
