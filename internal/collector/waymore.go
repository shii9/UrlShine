package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/shii9/UrlShine/internal/utils"
)

// runWaymore collects URLs via Waymore with aggressive mode 3 and targeted filters.
func runWaymore(target, outDir string, cfg Config) ([]string, error) {
	commands := [][]string{
		{"waymore", "-i", target, "-mode", "B", "-oU"},
		{"waymore", "-i", target, "-mode", "B", "-from", "2010", "-to", "202512", "-lcc", "0", "-oU"},
		{"waymore", "-i", target, "-mode", "B", "-ko", "api|auth|admin|user|v1|v2", "-oU"},
		{"waymore", "-i", target, "-mode", "B", "-fc", "200,301,302,403", "-oU"},
	}

	var wg sync.WaitGroup
	results := make(chan []string, len(commands))

	for i, args := range commands {
		i, args := i, args // capture for closure
		wg.Add(1)
		go func() {
			defer wg.Done()
			tmpFile := filepath.Join(outDir, fmt.Sprintf("_waymore_%s_%d.txt", utils.SanitizeFilename(target), i))
			fullArgs := append(args, tmpFile)

			_, err := runCmd(fullArgs...)
			if err != nil {
				return
			}

			lines, err := utils.ReadLines(tmpFile)
			if err == nil {
				results <- lines
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
