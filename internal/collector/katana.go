package collector

// runKatana collects URLs via Katana with professional aggressive parameters.
// Uses multiple command variations to maximize URL discovery:
// - High concurrency crawling
// - JavaScript execution and extraction
// - Query parameter discovery
// - Multi-depth traversal
// - Subdomain scope coverage
func runKatana(target, _ string, cfg Config) ([]string, error) {
	target = ensureHTTPS(target)
	var allUrls []string

	// Professional bug hunter command sequences to maximize coverage
	commands := [][]string{
		// 1. Deep Crawl + JS Parsing
		{"katana", "-u", target, "-d", "5", "-jc", "-kf", "all", "-silent"},
		// 2. Headless JS Rendering (Extracts XHR and browser-rendered content)
		{"katana", "-u", target, "-headless", "-d", "4", "-jc", "-xhr", "-jsonl", "-silent"},
		// 3. Root Domain + Subdomains scope
		{"katana", "-u", target, "-fs", "rdn", "-d", "4", "-jc", "-silent"},
		// 4. Query Parameter Focus
		{"katana", "-u", target, "-f", "qurl", "-silent"},
		// 5. Asset/Endpoint Filter (JS, JSON, JSP, no-extension)
		{"katana", "-u", target, "-d", "5", "-jc", "-em", "js,jsp,json,none", "-ndef", "-silent"},
		// 6. Basic Fast Pass
		{"katana", "-u", target, "-d", "3", "-silent"},
	}

	for _, args := range commands {
		lines, err := runCmd(args...)
		if err != nil {
			// Log error but continue with other commands to gather as much as possible
			continue
		}
		allUrls = append(allUrls, lines...)
	}

	return allUrls, nil
}
