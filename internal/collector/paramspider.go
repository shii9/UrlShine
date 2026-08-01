package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/shii9/UrlShine/internal/utils"
)

// runParamspider collects parameter-rich URLs via ParamSpider.
// ParamSpider mines web archives specifically for URLs containing query parameters,
// which are high-value targets for injection testing.
func runParamspider(target, outDir string, cfg Config) ([]string, error) {
	commands := [][]string{
		// Standard parameter mining
		{"paramspider", "-d", target, "--level", "high", "-o"},
		// Exclude common static extensions for cleaner output
		{"paramspider", "-d", target, "--exclude", "png,jpg,gif,svg,css,js,woff,ttf,ico", "-o"},
	}

	var wg sync.WaitGroup
	results := make(chan []string, len(commands))

	for i, args := range commands {
		i, args := i, args
		wg.Add(1)
		go func() {
			defer wg.Done()
			tmpFile := filepath.Join(outDir, fmt.Sprintf("_paramspider_%s_%d.txt", utils.SanitizeFilename(target), i))
			fullArgs := append(args, tmpFile)

			_, err := runCmd(fullArgs...)
			if err != nil {
				return
			}

			lines, err := utils.ReadLines(tmpFile)
			if err == nil {
				filtered := make([]string, 0, len(lines))
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
