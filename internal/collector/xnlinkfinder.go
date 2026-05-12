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
	linksFile := filepath.Join(outDir, fmt.Sprintf("_xnlf_%s.txt", utils.SanitizeFilename(target)))
	_ = os.Remove(linksFile)

	depth := cfg.Depth
	if depth < 2 {
		depth = 2
	}
	if depth > 3 {
		depth = 3
	}

	args := []string{
		"xnLinkFinder",
		"-i", target,
		"-sf", target,
		"-o", linksFile,
		"-d", fmt.Sprintf("%d", depth),
		"-t", "15",
		"-p", "10",
		"-timeout", "30",
		"-rl", "2",
	}
	if cfg.Threads < 5 {
		args = append(args, "-mtl", "5")
	} else {
		args = append(args, "-mtl", fmt.Sprintf("%d", cfg.Threads/10))
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
