package collector

import "fmt"

// runKatana collects URLs via Katana with aggressive parameters.
func runKatana(target, _ string, cfg Config) ([]string, error) {
	depth := cfg.Depth
	if depth < 3 {
		depth = 3
	}
	crawlWorkers := cfg.Threads
	if crawlWorkers < 20 {
		crawlWorkers = 20
	}
	args := []string{
		"katana",
		"-u", ensureHTTPS(target),
		"-d", fmt.Sprintf("%d", depth),
		"-c", fmt.Sprintf("%d", crawlWorkers),
		"-f", "url,qurl,js,fqdn,rdn,status,title,jsl,jsd",
		"-js-crawl",
		"-iqp",
		"-fx", ".png,.jpg,.jpeg,.gif,.svg,.css,.woff,.woff2,.ttf,.eot,.mp4,.mp3,.avi,.webm,.mkv,.pdf,.zip,.rar",
		"-no-color",
		"-silent",
	}
	if cfg.Subs {
		args = append(args, "-subdomains")
	}
	return runCmd(args...)
}
