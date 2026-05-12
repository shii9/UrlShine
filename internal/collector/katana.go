package collector

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shii9/UrlShine/internal/utils"
)

// runKatana collects URLs via Katana with aggressive crawling parameters.
func runKatana(target, outDir string, cfg Config) ([]string, error) {
	depth := cfg.Depth
	if depth < 3 {
		depth = 3
	}
	crawlWorkers := cfg.Threads
	if crawlWorkers < 30 {
		crawlWorkers = 30
	}

	tmpOut := filepath.Join(outDir, "_katana_"+utils.SanitizeFilename(target))
	_ = os.MkdirAll(tmpOut, 0755)

	args := []string{
		"katana",
		"-u", ensureHTTPS(target),
		"-d", fmt.Sprintf("%d", depth),
		"-c", fmt.Sprintf("%d", crawlWorkers),
		"-f", "url,qurl,js,fqdn,rdn,status,title,jsl,jsd,cdx",
		"-js-crawl",
		"-iqp",
		"-crawl-scope", "all",
		"-no-color",
		"-silent",
		"-o", tmpOut,
	}
	if cfg.Subs {
		args = append(args, "-subdomains")
	}

	_, _ = runCmd(args...)

	// Parse results from temp directory
	if !utils.FileExists(tmpOut) {
		return nil, nil
	}

	var lines []string
	f, err := os.Open(tmpOut)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}
