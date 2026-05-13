package collector

import "fmt"

// runGAU collects URLs via gau with aggressive parameters using multiple providers and options.
func runGAU(target, _ string, cfg Config) ([]string, error) {
	threads := cfg.Threads
	if threads < 10 {
		threads = 10
	}

	var allUrls []string

	// Professional bug hunter command sequences to maximize coverage
	commands := [][]string{
		// 1. Full Coverage Baseline
		{"gau", target, "--providers", "wayback,commoncrawl,otx,urlscan", "--threads", fmt.Sprintf("%d", threads)},
		// 2. Static Asset Filter
		{"gau", target, "--blacklist", "png,jpg,gif,svg,woff,ttf,css,js"},
		// 3. Deep Historical Sweep
		{"gau", target, "--from", "201801", "--to", "202312", "--providers", "wayback,commoncrawl,otx,urlscan"},
		// 4. Live Status Filter
		{"gau", target, "--mc", "200,301,302,403"},
		// 5. Parameterized URLs (duplicates filtered via --fp)
		{"gau", target, "--fp"},
	}

	for _, args := range commands {
		if cfg.Subs {
			// Add --subs to every command if subdomains are requested
			args = append(args, "--subs")
		}

		lines, err := runCmd(args...)
		if err != nil {
			// Log error but continue to gather as much as possible
			continue
		}
		allUrls = append(allUrls, lines...)
	}

	return allUrls, nil
}
