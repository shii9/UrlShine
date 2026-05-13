package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/shii9/UrlShine/internal/utils"
)

// runXnLinkFinder collects URLs via xnLinkFinder with aggressive parameters.
func runXnLinkFinder(target, outDir string, cfg Config) ([]string, error) {
	targetUrl := ensureHTTPS(target)
	targetDomain := target
	if idx := strings.Index(target, "://"); idx != -1 {
		targetDomain = target[idx+3:]
	}
	targetDomain = strings.Split(targetDomain, "/")[0]

	commands := [][]string{
		{"xnLinkFinder", "-i", targetUrl, "-sp", targetUrl, "-sf", targetDomain, "-d", "5", "-o"},
		{"xnLinkFinder", "-i", targetUrl, "-sp", targetUrl, "-sf", targetDomain, "-d", "5", "-all", "-o"},
		{"xnLinkFinder", "-i", filepath.Join(outDir, fmt.Sprintf("waymore_%s.txt", utils.SanitizeFilename(target))), "-sp", targetUrl, "-sf", targetDomain, "-d", "3", "-o"},
	}

	var wg sync.WaitGroup
	results := make(chan []string, len(commands))

	for i, args := range commands {
		i, args := i, args // capture for closure
		wg.Add(1)
		go func() {
			defer wg.Done()
			tmpFile := filepath.Join(outDir, fmt.Sprintf("_xnlf_%s_%d.txt", utils.SanitizeFilename(target), i))
			fullArgs := append(args, tmpFile)

			_, err := runCmd(fullArgs...)
			if err != nil {
				return
			}

			lines, err := utils.ReadLines(tmpFile)
			if err == nil {
				filtered := make([]string, 0)
				for _, l := range lines {
					l = strings.TrimSpace(l)
					if strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://") {
						filtered = append(filtered, l)
					}
				}
				results <- filtered
			}
			os.Remove(tmpFile)
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
