// Package splitter categorizes URLs into logical attack groups.
package splitter

import (
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/shii9/UrlShine/internal/utils"
)

const (
	GroupAPI    = "api_urls"
	GroupAuth   = "auth_admin_urls"
	GroupParams = "params_urls"
	GroupJS     = "js_config_urls"
	GroupDirs   = "directories_urls"
)

// Groups holds categorised URLs.
type Groups struct {
	API    []string
	Auth   []string
	Params []string
	JS     []string
	Dirs   []string
}

var (
	// API Endpoint Patterns: /api/*, /graphql, /v1-v9, /swagger, etc.
	reAPI = regexp.MustCompile(
		`(?i)(/api[/_v-]|/graphql|/gql|/rest[/_]|/v[0-9]+[./]|/rpc[./]|/grpc|` +
			`/swagger|/openapi|/redoc|/docs|/api-docs|/api\.json|/api\.yaml|` +
			`/endpoint|/service|/actuator|/metrics|/health|/status)(?:[/?#]|$)`)

	// Authentication & Admin Pages: login, auth, admin, dashboard, etc.
	reAuth = regexp.MustCompile(
		`(?i)\b(login|logout|log-in|log-out|signin|sign-in|sign-up|signup|` +
			`register|registration|auth[^o]|oauth|oidc|sso|saml|ldap|` +
			`admin|administrator|wp-admin|wp-login|phpmyadmin|cpanel|` +
			`dashboard|panel|console|portal|backoffice|manage|manager|` +
			`control|moderator|superuser|root|` +
			`forgot.*pass|reset.*pass|change.*pass|recover|` +
			`two.*factor|2fa|mfa|otp|verify|confirm|activate|unlock)(?:[/?#]|$)`)

	// JavaScript & Configuration Files
	reJS = regexp.MustCompile(
		`(?i)\.(js|mjs|jsx|ts|tsx|json|jsonp|php|phtml|env|yaml|yml|cfg|conf|ini|config|` +
			`xml|toml|properties|gradle|pom|lock|secret|key|pem|crt|cert)(\?|#|$)`)

	// Sensitive Configuration Paths
	reConfigPath = regexp.MustCompile(
		`(?i)/(config|settings|env|\.env|\.env\.local|secrets|credentials|backup|dump|` +
			`private|internal|admin|include|\.git|\.svn|\.hg|cvs|web\.config|app\.config|appsettings|` +
			`database|db|db\.php|config\.php|settings\.php|bootstrap)(?:[/?#]|$)`)

	reFileExt = regexp.MustCompile(`(?i)\.[a-z0-9]{1,8}(\?|#|$)`)
)

// Split maps URLs to groups based on regex patterns with priority-based categorization.
// Categories are checked in priority order to avoid overlaps:
// 1. JS/Config files (highest priority, most specific)
// 2. API endpoints
// 3. Auth pages
// 4. Query parameters
// 5. Directories (catchall)
func Split(urls []string) Groups {
	var g Groups
	for _, u := range urls {
		if u == "" {
			continue
		}
		parsed, err := url.Parse(u)
		if err != nil {
			continue
		}

		// Priority 1: JS/Config files (most specific)
		if reJS.MatchString(u) || reConfigPath.MatchString(parsed.Path) {
			g.JS = append(g.JS, u)
			continue
		}

		// Priority 2: API endpoints
		if reAPI.MatchString(u) {
			g.API = append(g.API, u)
			continue
		}

		// Priority 3: Auth/Admin pages
		if reAuth.MatchString(u) {
			g.Auth = append(g.Auth, u)
			continue
		}

		// Priority 4: Query parameters
		if parsed.RawQuery != "" {
			g.Params = append(g.Params, u)
			continue
		}

		// Priority 5: Directories (catchall)
		lastSeg := filepath.Base(parsed.Path)
		if lastSeg != "" && lastSeg != "." && lastSeg != "/" {
			if !reFileExt.MatchString(lastSeg) {
				g.Dirs = append(g.Dirs, u)
			}
		}
	}

	g.API = utils.DeduplicateSort(g.API)
	g.Auth = utils.DeduplicateSort(g.Auth)
	g.Params = utils.DeduplicateSort(g.Params)
	g.JS = utils.DeduplicateSort(g.JS)
	g.Dirs = utils.DeduplicateSort(g.Dirs)

	return g
}

// WriteGroups writes categorised URLs to output directory.
func WriteGroups(g Groups, outDir string) (map[string]string, error) {
	if err := utils.EnsureDir(outDir); err != nil {
		return nil, err
	}
	mapping := map[string][]string{
		GroupAPI:    g.API,
		GroupAuth:   g.Auth,
		GroupParams: g.Params,
		GroupJS:     g.JS,
		GroupDirs:   g.Dirs,
	}
	out := make(map[string]string, len(mapping))
	for name, lines := range mapping {
		path := outDir + "/" + name + ".txt"
		if err := utils.WriteLines(path, lines); err != nil {
			return nil, err
		}
		out[name] = path
	}
	return out, nil
}

// ParamKeys extracts sorted unique query parameter keys.
func ParamKeys(paramURLs []string) []string {
	seen := make(map[string]struct{})
	for _, u := range paramURLs {
		p, err := url.Parse(u)
		if err != nil {
			continue
		}
		for k := range p.Query() {
			seen[strings.ToLower(k)] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Counts returns the count of URLs in each group.
func Counts(g Groups) map[string]int {
	return map[string]int{
		GroupAPI:    len(g.API),
		GroupAuth:   len(g.Auth),
		GroupParams: len(g.Params),
		GroupJS:     len(g.JS),
		GroupDirs:   len(g.Dirs),
	}
}
