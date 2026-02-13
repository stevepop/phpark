package nginx

const nginxTemplate = `server {
    listen {{.ListenPort}};
    {{if .UseSSL}}listen 443 ssl http2;{{end}}
    server_name {{.ServerName}};
    root {{.Root}};

    {{if .UseSSL}}
    ssl_certificate {{.CertPath}};
    ssl_certificate_key {{.KeyPath}};
    {{end}}

    index index.php index.html index.htm;

    # Logging
    access_log /var/log/nginx/{{.SiteName}}.access.log;
    error_log /var/log/nginx/{{.SiteName}}.error.log;

    # Laravel/PHP framework friendly
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    # PHP-FPM configuration
    location ~ \.php$ {
        fastcgi_pass unix:{{.PHPSocket}};
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
    }

    # Deny access to hidden files
    location ~ /\. {
        deny all;
    }
}
`

// GetTemplate returns the nginx configuration template
func GetTemplate() string {
	return nginxTemplate
}
