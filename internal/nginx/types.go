package nginx

// SiteConfig represents nginx configuration for a site
type SiteConfig struct {
	// Site information
	SiteName   string // e.g., "myapp"
	Domain     string // e.g., "test"
	ServerName string // e.g., "myapp.test"

	// Paths
	Root     string // Document root (e.g., /Users/steve/sites/myapp/public)
	SitePath string // Full site path

	// PHP configuration
	PHPVersion string // e.g., "8.2"
	PHPSocket  string // e.g., "/var/run/php/php8.2-fpm.sock"

	// SSL
	UseSSL   bool
	CertPath string
	KeyPath  string

	// Additional
	ListenPort int // 80 or 443
}

// NginxConfig holds all nginx-related paths
type NginxConfig struct {
	SitesAvailable string // /etc/nginx/sites-available
	SitesEnabled   string // /etc/nginx/sites-enabled
	ConfigPath     string // Where PHPark puts configs
}
