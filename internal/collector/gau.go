package collector

import (
	"fmt"
	"sync"
)

// runGAU collects URLs via gau with aggressive parameters using multiple providers and options.
func runGAU(target, _ string, cfg Config) ([]string, error) {
	threads := cfg.Threads
	if threads < 10 {
		threads = 10
	}

	commands := [][]string{
		{"gau", target, "--providers", "wayback,commoncrawl,otx,urlscan", "--threads", fmt.Sprintf("%d", threads)},
		{"gau", target, "--blacklist", "png,jpg,gif,svg,woff,ttf,css,js"},
		{"gau", target, "--from", "201801", "--to", "202312", "--providers", "wayback,commoncrawl,otx,urlscan"},
		{"gau", target, "--mc", "200,301,302,403"},
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
