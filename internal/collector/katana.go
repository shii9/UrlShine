package collector

import "fmt"

// runKatana collects URLs via Katana with aggressive parameters.
func runKatana(target, _ string, cfg Config) ([]string, error) {
	args := []string{
		"katana",
		"-u", ensureHTTPS(target),
		"-d", fmt.Sprintf("%d", cfg.Depth),
		"-c", fmt.Sprintf("%d", cfg.Threads),
		"-f", "url,qurl,js",
		"-js-crawl",
		"-iqp",
		"-no-color",
		"-silent",
	}
	if cfg.Subs {
		args = append(args, "-subdomains")
	}
	return runCmd(args...)
}
