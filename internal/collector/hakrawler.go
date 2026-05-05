package collector

import "fmt"

// runHakrawler collects URLs via hakrawler with aggressive parameters.
func runHakrawler(target, _ string, cfg Config) ([]string, error) {
	args := []string{
		"hakrawler",
		"-u",
		"-d", fmt.Sprintf("%d", cfg.Depth),
		"-h", "User-Agent: Mozilla/5.0",
	}
	if cfg.Subs {
		args = append(args, "-subs")
	}
	return runCmdStdin(ensureHTTPS(target)+"\n", args...)
}
