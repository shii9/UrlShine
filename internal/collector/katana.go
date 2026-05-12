package collector

import "fmt"

// runKatana collects URLs via Katana with professional aggressive parameters.
// Uses multiple command variations to maximize URL discovery:
// - High concurrency crawling
// - JavaScript execution and extraction
// - Query parameter discovery
// - Multi-depth traversal
// - Subdomain scope coverage
func runKatana(target, _ string, cfg Config) ([]string, error) {
	depth := cfg.Depth
	if depth < 3 {
		depth = 3
	}
	crawlWorkers := cfg.Threads
	if crawlWorkers < 50 {
		crawlWorkers = 50
	}

	// Professional bug hunter combination: JS crawling + high concurrency + deep traversal
	args := []string{
		"katana",
		"-u", ensureHTTPS(target),
		"-d", fmt.Sprintf("%d", depth),
		"-c", fmt.Sprintf("%d", crawlWorkers),
		"-js-crawl",           // Execute JavaScript to discover dynamic URLs
		"-iqp",                // Include query parameters in output
		"-crawl-scope", "all", // Crawl all scopes (in-scope + out-of-scope)
		"-aff",      // Allow all file formats
		"-no-color", // Clean output
		"-silent",   // No progress messages
	}
	if cfg.Subs {
		args = append(args, "-subdomains")
	}

	// Run with stdout capture
	return runCmd(args...)
}
