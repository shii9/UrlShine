package collector

import (
	"fmt"
	"sync"
	"time"
)

// runGAU collects URLs via gau with aggressive parameters using multiple providers and options.
// GAU (GetAllUrls) fetches URLs from Wayback Machine, Common Crawl, URLScan, and AlienVault OTX.
func runGAU(target, _ string, cfg Config) ([]string, error) {
	threads := cfg.Threads
	if threads < 10 {
		threads = 10
	}

	// Dynamic date range — always extends to current month
	currentDate := time.Now().Format("200601")

	commands := [][]string{
		// All providers with max threads
		{"gau", target, "--providers", "wayback,commoncrawl,otx,urlscan", "--threads", fmt.Sprintf("%d", threads), "--retries", "3"},
		// Filtered output — exclude static assets to focus on interesting endpoints
		{"gau", target, "--blacklist", "png,jpg,gif,svg,woff,woff2,ttf,eot,css,ico,mp4,mp3,avi,webp,webm"},
		// Full date range for comprehensive historical coverage
		{"gau", target, "--from", "200801", "--to", currentDate, "--providers", "wayback,commoncrawl,otx,urlscan", "--retries", "3"},
		// Filter by interesting HTTP status codes
		{"gau", target, "--mc", "200,301,302,307,308,403,405,500"},
		// Fetch parameters mode — prioritizes URLs with parameters
		{"gau", target, "--fp"},
	}

	var wg sync.WaitGroup
	results := make(chan []string, len(commands))

	for _, args := range commands {
		args := args // capture for closure
		wg.Add(1)
		go func() {
			defer wg.Done()
			if cfg.Subs {
				args = append(args, "--subs")
			}
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
