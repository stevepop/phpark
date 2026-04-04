package config

import (
	"os"
	"os/user"
	"path/filepath"
)

const (
	// AppName is the application name
	AppName = "phpark"

	// ConfigFileName is the main config file
	ConfigFileName = "config.yaml"

	// SitesFileName stores the site registry
	SitesFileName = "sites.json"
)

// Paths holds all PHPark directory and file paths
type Paths struct {
	Home         string // ~/.phpark
	Config       string // ~/.phpark/config.yaml
	Sites        string // ~/.phpark/sites.json
	Nginx        string // ~/.phpark/nginx (generated configs)
	Certificates string // ~/.phpark/certificates (SSL certs)
	Logs         string // ~/.phpark/logs
}

// getEffectiveHomeDir returns the real user's home directory.
// When running via sudo, it returns SUDO_USER's home instead of /root,
// so PHPark config ends up in the real user's home.
func getEffectiveHomeDir() (string, error) {
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		if u, err := user.Lookup(sudoUser); err == nil {
			return u.HomeDir, nil
		}
	}
	return os.UserHomeDir()
}

// GetPaths returns all PHPark paths
func GetPaths() (*Paths, error) {
	homeDir, err := getEffectiveHomeDir()
	if err != nil {
		return nil, err
	}

	phparkHome := filepath.Join(homeDir, "."+AppName)

	return &Paths{
		Home:         phparkHome,
		Config:       filepath.Join(phparkHome, ConfigFileName),
		Sites:        filepath.Join(phparkHome, SitesFileName),
		Nginx:        filepath.Join(phparkHome, "nginx"),
		Certificates: filepath.Join(phparkHome, "certificates"),
		Logs:         filepath.Join(phparkHome, "logs"),
	}, nil
}

// EnsureDirectories creates all required directories if they don't exist
func (p *Paths) EnsureDirectories() error {
	directories := []string{
		p.Home,
		p.Nginx,
		p.Certificates,
		p.Logs,
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// Exists checks if PHPark is already installed (home directory exists)
func (p *Paths) Exists() bool {
	info, err := os.Stat(p.Home)
	if err != nil {
		return false
	}
	return info.IsDir()
}
