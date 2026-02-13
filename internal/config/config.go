package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads the configuration from config.yaml
// If the file doesn't exist, returns default config
func LoadConfig() (*Config, error) {
	paths, err := GetPaths()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return defaults
	if _, err := os.Stat(paths.Config); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read the file
	data, err := os.ReadFile(paths.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves the configuration to config.yaml
func SaveConfig(cfg *Config) error {
	paths, err := GetPaths()
	if err != nil {
		return err
	}

	// Ensure directories exist
	if err := paths.EnsureDirectories(); err != nil {
		return err
	}

	// Convert to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file (0644 = rw-r--r--)
	if err := os.WriteFile(paths.Config, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadSites loads the site registry from sites.json
// If the file doesn't exist, returns empty registry
func LoadSites() (*SiteRegistry, error) {
	paths, err := GetPaths()
	if err != nil {
		return nil, err
	}

	// If sites file doesn't exist, return empty registry
	if _, err := os.Stat(paths.Sites); os.IsNotExist(err) {
		return NewSiteRegistry(), nil
	}

	// Read the file
	data, err := os.ReadFile(paths.Sites)
	if err != nil {
		return nil, fmt.Errorf("failed to read sites file: %w", err)
	}

	// Parse JSON
	var registry SiteRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse sites file: %w", err)
	}

	return &registry, nil
}

// SaveSites saves the site registry to sites.json
func SaveSites(registry *SiteRegistry) error {
	paths, err := GetPaths()
	if err != nil {
		return err
	}

	// Ensure directories exist
	if err := paths.EnsureDirectories(); err != nil {
		return err
	}

	// Convert to pretty JSON
	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sites: %w", err)
	}

	// Write to file (0644 = rw-r--r--)
	if err := os.WriteFile(paths.Sites, data, 0644); err != nil {
		return fmt.Errorf("failed to write sites file: %w", err)
	}

	return nil
}
