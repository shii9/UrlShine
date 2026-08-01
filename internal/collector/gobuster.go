package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shii9/UrlShine/internal/utils"
)

// runGobuster runs directory brute-forcing via gobuster and returns found URLs.
func runGobuster(target, outDir string, cfg Config) ([]string, error) {
	targetUrl := ensureHTTPS(target)
	wordlist, cleanup := getWordlist()
	defer cleanup()

	if wordlist == "" {
		return nil, fmt.Errorf("no wordlist available")
	}

	tmpOut := filepath.Join(outDir, fmt.Sprintf("_gobuster_%s.txt", utils.SanitizeFilename(target)))

	// Run gobuster
	args := []string{
		"gobuster", "dir",
		"-u", targetUrl,
		"-w", wordlist,
		"-t", fmt.Sprintf("%d", cfg.Threads),
		"-o", tmpOut,
		"-q", "--no-error", "-k",
	}

	_, err := runCmd(args...)
	if err != nil {
		// If gobuster fails but wrote some output, try to parse it anyway
		if !utils.FileExists(tmpOut) {
			return nil, err
		}
	}

	defer os.Remove(tmpOut)

	// Parse Gobuster output
	lines, err := utils.ReadLines(tmpOut)
	if err != nil {
		return nil, err
	}

	var foundUrls []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Gobuster output format: /path (Status: 200) [Size: ...]
		parts := strings.Fields(line)
		if len(parts) > 0 {
			path := parts[0]
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			foundUrls = append(foundUrls, strings.TrimRight(targetUrl, "/")+path)
		}
	}

	return foundUrls, nil
}
