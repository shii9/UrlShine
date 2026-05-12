package collector

import "fmt"

// runGAU collects URLs via gau with aggressive parameters.
func runGAU(target, _ string, cfg Config) ([]string, error) {
	// Use maximum threads for GAU for faster collection
	threads := cfg.Threads
	if threads < 100 {
		threads = 100
	}
	args := []string{
		"gau", target,
		"--threads", fmt.Sprintf("%d", threads),
		"--providers", "wayback,commoncrawl,urlscan,otx",
		"--blacklist", "png,jpg,jpeg,gif,bmp,svg,ico,webp,css,eot,ttf,woff,woff2,eot,ttf,pdf,zip,rar,tar,gz,mp4,mp3,avi,webm,mkv,mov,flv,swf",
	}
	if cfg.Subs {
		args = append(args, "--subs")
	}
	return runCmd(args...)
}
