package collector

import "fmt"

// runHakrawler collects URLs via hakrawler with aggressive parameters.
func runHakrawler(target, _ string, cfg Config) ([]string, error) {
	depth := cfg.Depth
	if depth < 3 {
		depth = 3
	}
	args := []string{
		"hakrawler",
		"-u",
		"-d", fmt.Sprintf("%d", depth),
		"-h", "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}
	if cfg.Subs {
		args = append(args, "-subs")
	}
	return runCmdStdin(ensureHTTPS(target)+"\n", args...)
}
