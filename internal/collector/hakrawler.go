package collector

import (
	"fmt"
	"sync"
)

// runHakrawler collects URLs via hakrawler with multiple crawling strategies.
// hakrawler is a fast Go web crawler that parses HTML and JavaScript to find endpoints.
func runHakrawler(target, _ string, cfg Config) ([]string, error) {
	targetUrl := ensureHTTPS(target)

	commands := [][]string{
		// Standard crawl — parse HTML, follow links
		{"hakrawler", "-url", targetUrl, "-depth", fmt.Sprintf("%d", cfg.Depth), "-plain"},
		// Subdomain-inclusive crawl
		{"hakrawler", "-url", targetUrl, "-depth", fmt.Sprintf("%d", cfg.Depth), "-subs", "-plain"},
		// JavaScript-focused crawl
		{"hakrawler", "-url", targetUrl, "-depth", "3", "-js", "-plain"},
		// Robots.txt and sitemap parsing
		{"hakrawler", "-url", targetUrl, "-depth", "2", "-robots", "-sitemap", "-plain"},
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
