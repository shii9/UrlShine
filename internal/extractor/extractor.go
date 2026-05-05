// Package extractor provides advanced URL extraction using specialized tools.
package extractor

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"urlshine/internal/logger"
	"urlshine/internal/utils"
)

// GroupExtractor handles advanced extraction for specialized URL groups.
type GroupExtractor struct {
	URLs      []string
	OutputDir string
}

// NewGroupExtractor creates a new extractor instance.
func NewGroupExtractor(urls []string, outputDir string) *GroupExtractor {
	return &GroupExtractor{
		URLs:      urls,
		OutputDir: outputDir,
	}
}

// ExtractAll runs all specialized extractors.
func (ge *GroupExtractor) ExtractAll() map[string][]string {
	results := make(map[string][]string)

	// Run extractors concurrently
	var wg sync.WaitGroup
	mu := sync.Mutex{}

	extractors := []struct {
		name string
		fn   func() ([]string, error)
	}{
		{"API Endpoints", ge.ExtractAPIs},
		{"Auth/Admin", ge.ExtractAuthAdmin},
		{"Parameters", ge.ExtractParameters},
		{"JS/Config", ge.ExtractJSConfig},
		{"Directories", ge.ExtractDirectories},
	}

	for _, ex := range extractors {
		wg.Add(1)
		go func(name string, fn func() ([]string, error)) {
			defer wg.Done()
			logger.Run("  ⟳ Extracting %s...", name)
			urls, err := fn()
			if err != nil {
				logger.Warn("  ✗ %s error: %v", name, err)
				return
			}
			mu.Lock()
			results[name] = urls
			logger.Success("  ✓ %s: %d URLs", name, len(urls))
			mu.Unlock()
		}(ex.name, ex.fn)
	}

	wg.Wait()
	return results
}

// ExtractAPIs finds API endpoints using grep patterns + httpx live checking.
func (ge *GroupExtractor) ExtractAPIs() ([]string, error) {
	// API patterns: /api/, graphql, /v1/, /v2/, rest, swagger, openapi
	patterns := []string{
		"/api/",
		"graphql",
		"/v1/",
		"/v2/",
		"rest",
		"swagger",
		"openapi",
		"/json",
		"/xml",
	}

	var apis []string
	seen := make(map[string]struct{})

	for _, url := range ge.URLs {
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(url), strings.ToLower(pattern)) {
				if _, ok := seen[url]; !ok {
					seen[url] = struct{}{}
					apis = append(apis, url)
					break
				}
			}
		}
	}

	// Live verification with httpx if available
	if commandExists("httpx") {
		liveAPIs := ge.verifyLiveWithHttpx(apis)
		return liveAPIs, nil
	}

	return apis, nil
}

// ExtractAuthAdmin finds authentication and admin endpoints.
func (ge *GroupExtractor) ExtractAuthAdmin() ([]string, error) {
	patterns := []string{
		"admin", "login", "signin", "signup",
		"auth", "oauth", "callback",
		"panel", "dashboard", "manage",
		"private", "internal", "staff",
		"control", "root", "user",
		"account", "register", "forgot",
		"reset", "2fa", "mfa", "security",
		"profile", "settings", "config",
		"secret", "priv", "restricted",
	}

	var authAdmin []string
	seen := make(map[string]struct{})

	for _, url := range ge.URLs {
		urlLower := strings.ToLower(url)
		for _, pattern := range patterns {
			if strings.Contains(urlLower, pattern) {
				if _, ok := seen[url]; !ok {
					seen[url] = struct{}{}
					authAdmin = append(authAdmin, url)
					break
				}
			}
		}
	}

	// Use gobuster/feroxbuster for hidden paths if available
	if commandExists("gobuster") {
		authAdmin = append(authAdmin, ge.extractWithGobuster(authAdmin)...)
	}

	return utils.DeduplicateSort(authAdmin), nil
}

// ExtractParameters finds URLs with parameters.
func (ge *GroupExtractor) ExtractParameters() ([]string, error) {
	// Look for URLs with query parameters
	var withParams []string
	seen := make(map[string]struct{})

	for _, url := range ge.URLs {
		if strings.Contains(url, "?") {
			if _, ok := seen[url]; !ok {
				seen[url] = struct{}{}
				withParams = append(withParams, url)
			}
		}
	}

	// If available, use paramspider or arjun for additional discovery
	if commandExists("arjun") {
		discovered := ge.extractWithArjun(withParams)
		withParams = append(withParams, discovered...)
	}

	return utils.DeduplicateSort(withParams), nil
}

// ExtractJSConfig finds JavaScript and configuration files.
func (ge *GroupExtractor) ExtractJSConfig() ([]string, error) {
	patterns := []string{
		".js", ".json", ".xml",
		".env", ".bak", ".old",
		"config", "settings", "app.js",
		"main.js", "bundle.js", "vendor.js",
		".min.js", "static", "assets",
	}

	var jsConfig []string
	seen := make(map[string]struct{})

	for _, url := range ge.URLs {
		urlLower := strings.ToLower(url)
		for _, pattern := range patterns {
			if strings.Contains(urlLower, strings.ToLower(pattern)) {
				if _, ok := seen[url]; !ok {
					seen[url] = struct{}{}
					jsConfig = append(jsConfig, url)
					break
				}
			}
		}
	}

	return jsConfig, nil
}

// ExtractDirectories finds directory paths.
func (ge *GroupExtractor) ExtractDirectories() ([]string, error) {
	var dirs []string
	seen := make(map[string]struct{})

	for _, url := range ge.URLs {
		// Parse URL to extract path
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			continue
		}

		// Remove scheme and domain
		path := ge.extractPath(url)
		if path != "/" && path != "" {
			// Count slashes to find directory-like paths
			slashCount := strings.Count(path, "/")
			// If path has more than 2 slashes or ends with /, likely a directory
			if slashCount >= 2 || strings.HasSuffix(path, "/") {
				if _, ok := seen[url]; !ok {
					seen[url] = struct{}{}
					dirs = append(dirs, url)
				}
			}
		}
	}

	return dirs, nil
}

// Helper methods

func (ge *GroupExtractor) extractPath(url string) string {
	// Remove scheme
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	// Find first slash
	if idx := strings.Index(url, "/"); idx > -1 {
		return url[idx:]
	}
	return "/"
}

func (ge *GroupExtractor) verifyLiveWithHttpx(urls []string) []string {
	if len(urls) == 0 {
		return urls
	}

	tmpFile := filepath.Join(ge.OutputDir, "temp_urls.txt")
	utils.WriteLines(tmpFile, urls)
	defer os.Remove(tmpFile)

	cmd := exec.Command("httpx", "-l", tmpFile, "-silent", "-mc", "200,201,204,301,302")
	output, err := cmd.Output()
	if err != nil {
		return urls
	}

	var liveURLs []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			liveURLs = append(liveURLs, line)
		}
	}

	return liveURLs
}

func (ge *GroupExtractor) extractWithGobuster(seedURLs []string) []string {
	// Extract unique domains from URLs
	domains := make(map[string]struct{})
	for _, url := range seedURLs {
		domain := ge.extractDomain(url)
		if domain != "" {
			domains[domain] = struct{}{}
		}
	}

	var discovered []string
	for domain := range domains {
		cmd := exec.Command("gobuster", "dir",
			"-u", "https://"+domain,
			"-w", "/usr/share/wordlists/dirb/common.txt",
			"-t", "30",
			"-q", "--no-error", "-k")

		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					url := "https://" + domain + strings.Fields(line)[0]
					discovered = append(discovered, url)
				}
			}
		}
	}

	return discovered
}

func (ge *GroupExtractor) extractWithArjun(urls []string) []string {
	if len(urls) == 0 {
		return []string{}
	}

	var discovered []string

	// For each unique domain, run arjun
	domains := make(map[string]struct{})
	for _, url := range urls {
		if domain := ge.extractDomain(url); domain != "" {
			domains[domain] = struct{}{}
		}
	}

	for domain := range domains {
		cmd := exec.Command("arjun", "-u", "https://"+domain, "--quiet")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "http") {
					discovered = append(discovered, line)
				}
			}
		}
	}

	return discovered
}

func (ge *GroupExtractor) extractDomain(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	if idx := strings.Index(url, "/"); idx > -1 {
		return url[:idx]
	}
	return url
}

// commandExists checks if a command exists in PATH.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// SaveResults saves extracted groups to files.
func (ge *GroupExtractor) SaveResults(groups map[string][]string) error {
	for name, urls := range groups {
		// Convert group name to filename
		filename := safeGroupFilename(name)
		filename += ".txt"

		filepath := filepath.Join(ge.OutputDir, filename)
		if err := utils.WriteLines(filepath, urls); err != nil {
			logger.Warn("Failed to save %s: %v", filename, err)
			continue
		}

		logger.Success("  ✓ Saved: %s (%d URLs)", filename, len(urls))
	}

	return nil
}

func safeGroupFilename(name string) string {
	name = strings.ToLower(name)
	var b strings.Builder
	lastUnderscore := false
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			b.WriteByte('_')
			lastUnderscore = true
		}
	}
	return strings.Trim(b.String(), "_")
}
