package collector

import "fmt"

// runWaymore collects URLs via Waymore with aggressive parameters.
func runWaymore(target, _ string, cfg Config) ([]string, error) {
	conc := 25
	if cfg.Threads > 50 {
		conc = 50
	}
	args := []string{
		"waymore",
		"-i", target,
		"-mode", "3",
		"-p",
		"-exclude", "png|jpg|jpeg|gif|bmp|svg|ico|css|woff|woff2|eot|ttf|pdf|zip|rar|tar|gz|mp4|mp3",
		"-concurrency", fmt.Sprintf("%d", conc),
	}
	return runCmd(args...)
}
