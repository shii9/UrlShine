package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shii9/UrlShine/internal/utils"
)

// runXnLinkFinder collects URLs via xnLinkFinder with aggressive parameters.
func runXnLinkFinder(target, outDir string, cfg Config) ([]string, error) {
	targetUrl := ensureHTTPS(target)
	targetDomain := target
	// Strip protocol for -sf
	if idx := strings.Index(target, "://"); idx != -1 {
		targetDomain = target[idx+3:]
	}
	targetDomain = strings.Split(targetDomain, "/")[0]

	var allUrls []string

	// Professional bug hunter command sequences to maximize coverage
	commands := [][]string{
		// 1. Main Live Crawl
		{"xnLinkFinder", "-i", targetUrl, "-sp", targetUrl, "-sf", targetDomain, "-d", "5", "-o"},
		// 2. Max Recall Version (includes less common TLDs)
		{"xnLinkFinder", "-i", targetUrl, "-sp", targetUrl, "-sf", targetDomain, "-d", "5", "-all", "-o"},
		// 3. Waymore Integrated Crawl (feeds waymore results back into xnLinkFinder)
		{"xnLinkFinder", "-i", filepath.Join(outDir, fmt.Sprintf("waymore_%s.txt", utils.SanitizeFilename(target))), "-sp", targetUrl, "-sf", targetDomain, "-d", "3", "-o"},
	}

	for i, args := range commands {
		// Create unique tmp file for each run
		tmpFile := filepath.Join(outDir, fmt.Sprintf("_xnlf_%s_%d.txt", utils.SanitizeFilename(target), i))
		fullArgs := append(args, tmpFile)

		_, err := runCmd(fullArgs...)
		if err != nil {
			continue
		}

		// Read and filter results
		lines, err := utils.ReadLines(tmpFile)
		if err == nil {
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://") {
					allUrls = append(allUrls, l)
				}
			}
		}
		os.Remove(tmpFile)
	}

	return allUrls, nil
}
