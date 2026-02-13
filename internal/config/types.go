package config

// Config represents the main PHPark configuration
type Config struct {
	// DefaultPHP is the default PHP version to use (e.g., "8.2", "8.3")
	DefaultPHP string `json:"default_php" yaml:"default_php"`

	// Domain is the TLD for local sites (default: "test")
	Domain string `json:"domain" yaml:"domain"`

	// NginxConfigPath is where nginx config files are stored
	NginxConfigPath string `json:"nginx_config_path" yaml:"nginx_config_path"`

	// UseHTTPS indicates if sites should use HTTPS by default
	UseHTTPS bool `json:"use_https" yaml:"use_https"`
}

// Site represents a single parked or linked site
type Site struct {
	// Name is the site name (e.g., "myapp" for myapp.test)
	Name string `json:"name"`

	// Path is the full path to the site directory
	Path string `json:"path"`

	// Type is either "park" or "link"
	Type string `json:"type"`

	// PHPVersion is the PHP version for this site (e.g., "8.2")
	// If empty, uses the default PHP version
	PHPVersion string `json:"php_version,omitempty"`

	// Secured indicates if the site uses HTTPS
	Secured bool `json:"secured"`
}

// SiteRegistry holds all registered sites
type SiteRegistry struct {
	Sites []Site `json:"sites"`
}

// DefaultConfig returns a new Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		DefaultPHP:      "8.2",
		Domain:          "test",
		NginxConfigPath: "/etc/nginx/sites-enabled",
		UseHTTPS:        false,
	}
}

// NewSiteRegistry creates an empty site registry
func NewSiteRegistry() *SiteRegistry {
	return &SiteRegistry{
		Sites: []Site{},
	}
}

// FindSite searches for a site by name
func (sr *SiteRegistry) FindSite(name string) *Site {
	for i := range sr.Sites {
		if sr.Sites[i].Name == name {
			return &sr.Sites[i]
		}
	}
	return nil
}

// AddSite adds or updates a site in the registry
func (sr *SiteRegistry) AddSite(site Site) {
	// Check if site already exists
	for i := range sr.Sites {
		if sr.Sites[i].Name == site.Name {
			// Update existing site
			sr.Sites[i] = site
			return
		}
	}
	// Add new site
	sr.Sites = append(sr.Sites, site)
}

// RemoveSite removes a site from the registry
func (sr *SiteRegistry) RemoveSite(name string) bool {
	for i := range sr.Sites {
		if sr.Sites[i].Name == name {
			// Remove by slicing
			sr.Sites = append(sr.Sites[:i], sr.Sites[i+1:]...)
			return true
		}
	}
	return false
}

// ListSites returns all sites
func (sr *SiteRegistry) ListSites() []Site {
	return sr.Sites
}
