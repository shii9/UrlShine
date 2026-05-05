package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"urlshine/internal/utils"
)

// runXnLinkFinder collects URLs via xnLinkFinder.
func runXnLinkFinder(target, outDir string, cfg Config) ([]string, error) {
	linksFile := filepath.Join(outDir, fmt.Sprintf("_xnlf_%s.txt", utils.SanitizeFilename(target)))
	_ = os.Remove(linksFile)

	args := []string{
		"xnLinkFinder",
		"-i", target,
		"-sf", target,
		"-o", linksFile,
		"-nwll", // no wordlist
		"-timeout", "30",
	}
	if cfg.Depth > 0 {
		depth := cfg.Depth
		if depth > 3 {
			depth = 3
		}
		args = append(args, "-d", fmt.Sprintf("%d", depth))
	}

	_, _ = runCmd(args...)

	if !utils.FileExists(linksFile) {
		return nil, nil
	}
	lines, err := utils.ReadLines(linksFile)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://") {
			out = append(out, l)
		}
	}
	return out, nil
}
