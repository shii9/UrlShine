package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shii9/UrlShine/internal/utils"
)

// runWaymore collects URLs via Waymore with aggressive mode and targeted filters.
// Waymore provides enhanced Wayback Machine access with content downloading and filtering.
func runWaymore(target, outDir string, cfg Config) ([]string, error) {
	// Dynamic date range — always extends to current year+month
	currentDate := time.Now().Format("200601")

	commands := [][]string{
		// Standard mode B (both URLs and responses)
		{"waymore", "-i", target, "-mode", "B", "-oU"},
		// Full historical date range
		{"waymore", "-i", target, "-mode", "B", "-from", "2008", "-to", currentDate, "-lcc", "0", "-oU"},
		// Filter for high-value keywords (API, auth, admin, etc.)
		{"waymore", "-i", target, "-mode", "B", "-ko", "api|auth|admin|user|v1|v2|v3|login|signup|register|token|key|secret|config|backup|debug", "-oU"},
		// Filter by interesting status codes
		{"waymore", "-i", target, "-mode", "B", "-fc", "200,301,302,307,403,405,500", "-oU"},
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
