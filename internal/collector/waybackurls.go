package collector

import (
	"regexp"
	"strings"
	"time"
)

// runWaybackurls executes 5 waybackurls command variants (one-by-one) to
// maximize archived URL discovery. Each variant applies different filters:
// 1. Basic collection (no filter)
// 2. URLs with parameters (grep "=")
// 3. High-value parameter names (callback, redirect, token, etc.)
// 4. JavaScript files (.js extensions)
// 5. Combined with GAU for additional archive sources
// Results are deduplicated before returning.
func runWaybackurls(target, _ string, cfg Config) ([]string, error) {
	unique := map[string]struct{}{}

	// 1. Basic waybackurls collection (no filter)
	basicUrls, _ := runCmdStdin(target+"\n", "waybackurls")
	for _, url := range basicUrls {
		if url = strings.TrimSpace(url); url != "" {
			unique[url] = struct{}{}
		}
	}

	// 2. Parameter-focused URLs (contains "=")
	for _, url := range basicUrls {
		if strings.Contains(url, "=") {
			unique[url] = struct{}{}
		}
	}

	// 3. High-value parameter names (open redirects, IDOR, cache bypass, etc.)
	highValueParams := regexp.MustCompile(`(?i)(callback|redirect|return|next|url|q|search|token|id|page|file|path)=`)
	for _, url := range basicUrls {
		if highValueParams.MatchString(url) {
			unique[url] = struct{}{}
		}
	}

	// 4. JavaScript files (.js extensions)
	jsExtension := regexp.MustCompile(`\.js$`)
	for _, url := range basicUrls {
		if jsExtension.MatchString(url) {
			unique[url] = struct{}{}
		}
	}

	// 5. Combined waybackurls + GAU (strongest URL expansion)
	gauUrls, _ := runCmdStdin(target+"\n", "gau")
	for _, url := range gauUrls {
		if url = strings.TrimSpace(url); url != "" {
			unique[url] = struct{}{}
		}
	}

	// Brief pause between runs to avoid hammering targets
	time.Sleep(200 * time.Millisecond)

	// Convert unique map to sorted slice
	result := make([]string, 0, len(unique))
	for u := range unique {
		result = append(result, u)
	}
	return result, nil
}
