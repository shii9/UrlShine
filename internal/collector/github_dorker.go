package collector

import (
	"os"
	"strings"
	"sync"
)

// runGithubEndpoints discovers URLs/endpoints leaked in GitHub repositories.
// Uses github-endpoints tool to search public repos for target domain references.
// Requires GITHUB_TOKEN environment variable for higher rate limits.
func runGithubEndpoints(target, _ string, cfg Config) ([]string, error) {
	// Check for GitHub token (optional but recommended)
	token := os.Getenv("GITHUB_TOKEN")

	commands := [][]string{
		// Standard endpoint search
		{"github-endpoints", "-d", target, "-raw"},
	}

	if token != "" {
		// With authentication for higher rate limits
		commands = append(commands,
			[]string{"github-endpoints", "-d", target, "-t", token, "-raw"},
		)
	}

	// Also try github-subdomains for additional coverage
	commands = append(commands,
		[]string{"github-subdomains", "-d", target, "-raw"},
	)

	var wg sync.WaitGroup
	results := make(chan []string, len(commands))

	for _, args := range commands {
		args := args
		wg.Add(1)
		go func() {
			defer wg.Done()
			lines, err := runCmd(args...)
			if err == nil {
				// Filter to only include URLs or construct URLs from endpoints
				var filtered []string
				for _, l := range lines {
					l = strings.TrimSpace(l)
					if strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://") {
						filtered = append(filtered, l)
					} else if l != "" && !strings.HasPrefix(l, "#") {
						// Construct URL from path-like output
						filtered = append(filtered, ensureHTTPS(target)+"/"+strings.TrimPrefix(l, "/"))
					}
				}
				results <- filtered
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
