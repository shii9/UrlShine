package collector

import (
	"fmt"
	"sync"
)

// runUrlfinder collects URLs via urlfinder which uses passive intelligence sources
// including crt.sh, VirusTotal, and other data providers.
func runUrlfinder(target, _ string, cfg Config) ([]string, error) {
	targetUrl := ensureHTTPS(target)

	commands := [][]string{
		// Standard passive discovery
		{"urlfinder", "-u", targetUrl, "-all", "-silent"},
		// Domain-only mode
		{"urlfinder", "-d", target, "-all", "-silent"},
		// With subdomain inclusion
		{"urlfinder", "-d", target, "-all", "-s", "-silent"},
	}

	var wg sync.WaitGroup
	results := make(chan []string, len(commands))

	for _, args := range commands {
		args := args
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

// runUrlfinderFallback tries alternative urlfinder binaries (url-finder, URL-finder).
func runUrlfinderFallback(target, outDir string, cfg Config) ([]string, error) {
	// Some systems install as "url-finder" instead of "urlfinder"
	targetUrl := ensureHTTPS(target)

	lines, err := runCmd("url-finder", "-u", targetUrl, "-all", "-silent")
	if err != nil {
		return nil, fmt.Errorf("url-finder not available: %w", err)
	}
	return lines, nil
}
