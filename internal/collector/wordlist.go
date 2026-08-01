package collector

import (
	"os"
	"path/filepath"

	"github.com/shii9/UrlShine/internal/utils"
)

var commonWordlists = []string{
	"/usr/share/wordlists/dirb/common.txt",
	"/usr/share/wordlists/dirbuster/directory-list-2.3-medium.txt",
	"/usr/share/wordlists/dirbuster/directory-list-2.3-small.txt",
	"/usr/share/wordlists/dirb/small.txt",
	"/usr/share/seclists/Discovery/Web-Content/common.txt",
	"/usr/share/seclists/Discovery/Web-Content/raft-medium-words.txt",
}

// Expanded fallback wordlist — covers common directories, files, API paths,
// admin panels, backup files, config files, and typical web application paths.
var defaultWords = []string{
	// Admin & Auth
	"admin", "administrator", "login", "signin", "signup", "register", "auth", "oauth",
	"sso", "logout", "session", "account", "profile", "user", "users", "dashboard",
	"panel", "console", "portal", "manage", "manager", "control", "root", "staff",
	"moderator", "superadmin", "wp-admin", "wp-login.php", "phpmyadmin", "cpanel",
	"webmail", "admin.php", "admin.html",

	// API & Services
	"api", "v1", "v2", "v3", "v4", "rest", "graphql", "gql", "grpc", "rpc",
	"api-docs", "swagger", "swagger-ui", "openapi", "redoc", "docs", "endpoints",
	"actuator", "health", "metrics", "status", "info", "version", "ping",
	"api.php", "api.json", "api.yaml", "api.yml", "service", "services",

	// Configuration & Secrets
	"config", "configuration", ".env", ".env.local", ".env.production", ".env.staging",
	"secrets", "credentials", "database", "db", "db.php", "config.php", "settings",
	"settings.php", "bootstrap", "web.config", "app.config", "appsettings.json",
	".htaccess", ".htpasswd", "wp-config.php", "wp-config.php.bak",

	// Common Files & Resources
	"robots.txt", "sitemap.xml", "sitemap_index.xml", "crossdomain.xml",
	"favicon.ico", "manifest.json", "humans.txt", "security.txt",
	".well-known", ".well-known/security.txt", ".well-known/openid-configuration",
	"index.php", "index.html", "index.asp", "index.jsp", "default.aspx",

	// Backup & Development
	"backup", "backups", "bak", "old", "copy", "temp", "tmp", "test", "testing",
	"demo", "dev", "development", "staging", "prod", "production", "debug",
	"trace", "log", "logs", "error_log", "access_log",

	// Static Assets
	"assets", "static", "public", "media", "js", "css", "images", "img",
	"upload", "uploads", "files", "download", "downloads", "content", "data",
	"resources", "dist", "build", "node_modules",

	// Common App Paths
	"app", "application", "home", "main", "about", "contact", "help", "faq",
	"search", "blog", "news", "forum", "shop", "store", "cart", "checkout",
	"payment", "order", "orders", "invoice", "ticket", "tickets", "support",

	// Version Control & CI
	".git", ".git/config", ".git/HEAD", ".svn", ".svn/entries", ".hg",
	".env.example", ".gitignore", ".dockerignore", "Dockerfile", "docker-compose.yml",
	"Makefile", "Jenkinsfile", ".travis.yml", ".gitlab-ci.yml",

	// Server-specific
	"server-status", "server-info", "phpinfo.php", "info.php", "test.php",
	"cgi-bin", "cgi-bin/test-cgi", "fcgi-bin",
	"xmlrpc.php", "wp-cron.php", "wp-content", "wp-includes",

	// Monitoring & Debug
	"monitor", "monitoring", "grafana", "prometheus", "kibana", "elastic",
	"solr", "jenkins", "travis", "circleci",
}

// getWordlist returns the path to a wordlist and a cleanup function.
// It prioritizes standard pentesting wordlists if they exist.
// Otherwise, it creates a temporary fallback wordlist with comprehensive entries.
func getWordlist() (string, func()) {
	for _, path := range commonWordlists {
		if utils.FileExists(path) {
			return path, func() {} // No cleanup needed
		}
	}

	// Create fallback temporary wordlist
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "urlshine_fallback_wordlist.txt")
	err := utils.WriteLines(tmpFile, defaultWords)
	if err != nil {
		// If temporary creation fails, return empty path
		return "", func() {}
	}

	return tmpFile, func() {
		os.Remove(tmpFile)
	}
}
