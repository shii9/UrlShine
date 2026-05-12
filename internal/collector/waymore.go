package collector

import "fmt"

// runWaymore collects URLs via Waymore with aggressive parameters.
func runWaymore(target, _ string, cfg Config) ([]string, error) {
	conc := 50
	if cfg.Threads < 50 && cfg.Threads > 0 {
		conc = cfg.Threads
	}
	if cfg.Threads > 100 {
		conc = 100
	}
	args := []string{
		"waymore",
		"-i", target,
		"-mode", "3",
		"-p",
		"-formatted",
		"-include", "api|auth|admin|user|v1|v2|v3|endpoint",
		"-exclude", "png|jpg|jpeg|gif|bmp|svg|ico|webp|css|woff|woff2|eot|ttf|pdf|zip|rar|tar|gz|mp4|mp3|avi|webm|mkv|mov|flv|swf",
		"-concurrency", fmt.Sprintf("%d", conc),
	}
	return runCmd(args...)
}
