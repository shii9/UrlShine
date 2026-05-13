package collector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shii9/UrlShine/internal/utils"
)

// runWaymore collects URLs via Waymore with aggressive mode 3 and targeted filters.
func runWaymore(target, outDir string, cfg Config) ([]string, error) {
	var allUrls []string

	// Professional bug hunter command sequences to maximize coverage
	commands := [][]string{
		// 1. Main Deep Recon
		{"waymore", "-i", target, "-mode", "B", "-oU"},
		// 2. Historical Discovery (Deep sweep with all Common Crawl collections)
		{"waymore", "-i", target, "-mode", "B", "-from", "2010", "-to", "202512", "-lcc", "0", "-oU"},
		// 3. High-Value Target Filtering (API, Auth, Admin, etc.)
		{"waymore", "-i", target, "-mode", "B", "-ko", "api|auth|admin|user|v1|v2", "-oU"},
		// 4. Response Code Filter
		{"waymore", "-i", target, "-mode", "B", "-fc", "200,301,302,403", "-oU"},
	}

	for i, args := range commands {
		// Waymore requires a file path for -oU
		tmpFile := filepath.Join(outDir, fmt.Sprintf("_waymore_%s_%d.txt", utils.SanitizeFilename(target), i))
		fullArgs := append(args, tmpFile)

		_, err := runCmd(fullArgs...)
		if err != nil {
			continue
		}

		// Read the resulting file
		lines, err := utils.ReadLines(tmpFile)
		if err == nil {
			allUrls = append(allUrls, lines...)
		}

		// Cleanup tmp file
		os.Remove(tmpFile)
	}

	return allUrls, nil
}
