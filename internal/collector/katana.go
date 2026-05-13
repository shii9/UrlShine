package collector

import (
	"sync"
)

// runKatana collects URLs via Katana with professional aggressive parameters.
func runKatana(target, _ string, cfg Config) ([]string, error) {
	target = ensureHTTPS(target)

	commands := [][]string{
		{"katana", "-u", target, "-d", "5", "-jc", "-kf", "all", "-silent"},
		{"katana", "-u", target, "-headless", "-d", "4", "-jc", "-xhr", "-jsonl", "-silent"},
		{"katana", "-u", target, "-fs", "rdn", "-d", "4", "-jc", "-silent"},
		{"katana", "-u", target, "-f", "qurl", "-silent"},
		{"katana", "-u", target, "-d", "5", "-jc", "-em", "js,jsp,json,none", "-ndef", "-silent"},
		{"katana", "-u", target, "-d", "3", "-silent"},
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
