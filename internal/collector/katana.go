package collector

import (
	"fmt"
	"sync"
)

// runKatana collects URLs via Katana with professional aggressive parameters.
// Katana is an advanced JavaScript-capable crawler with headless Chromium support.
func runKatana(target, _ string, cfg Config) ([]string, error) {
	target = ensureHTTPS(target)

	depth := cfg.Depth
	if depth < 3 {
		depth = 3
	}

	commands := [][]string{
		// JS-aware deep crawl with all known field types
		{"katana", "-u", target, "-d", fmt.Sprintf("%d", depth), "-jc", "-kf", "all", "-silent"},
		// Headless browser mode for JS-heavy SPAs with XHR capture
		{"katana", "-u", target, "-headless", "-d", fmt.Sprintf("%d", depth-1), "-jc", "-xhr", "-silent"},
		// Field scope: root domain only (avoids third-party crawl noise)
		{"katana", "-u", target, "-fs", "rdn", "-d", fmt.Sprintf("%d", depth-1), "-jc", "-silent"},
		// Query URL extraction mode — finds URLs with parameters
		{"katana", "-u", target, "-f", "qurl", "-silent"},
		// Extension match: JS, JSP, JSON, and no-extension endpoints
		{"katana", "-u", target, "-d", fmt.Sprintf("%d", depth), "-jc", "-em", "js,jsp,json,php,aspx,none", "-ndef", "-silent"},
		// Shallow but broad crawl for quick wins
		{"katana", "-u", target, "-d", "3", "-silent"},
		// Automatic form filling for discovering hidden endpoints
		{"katana", "-u", target, "-d", fmt.Sprintf("%d", depth-1), "-aff", "-jc", "-silent"},
		// Crawl with custom header for API discovery
		{"katana", "-u", target, "-d", fmt.Sprintf("%d", depth), "-H", "Accept: application/json", "-jc", "-silent"},
	}

	var wg sync.WaitGroup
	results := make(chan []string, len(commands))

	for _, args := range commands {
		args := args // capture for closure
		wg.Add(1)
		go func() {
			defer wg.Done()
			lines, err := runCmd(args...)
			if err == nil {
				results <- lines
			}
		}()
	}

	wg.Wait()
	close(results)

	var allUrls []string
	for lines := range results {
		allUrls = append(allUrls, lines...)
	}

	return allUrls, nil
}
