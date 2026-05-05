// Package merger combines per-tool outputs into a single unique list.
package merger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shii9/UrlShine/internal/logger"
	"github.com/shii9/UrlShine/internal/utils"
)

// MergeDir reads all .txt files from dir, skipping internal temp files,
// and writes all unique HTTP(S) URLs to outFile.
func MergeDir(dir, outFile string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("read dir: %w", err)
	}

	seen := make(map[string]struct{}, 1_000_000)
	var merged []string
	skipPrefixes := []string{"_gospider_", "_xnlf_", "targets_"}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		if e.Name() == filepath.Base(outFile) {
			continue
		}

		skip := false
		for _, p := range skipPrefixes {
			if strings.HasPrefix(e.Name(), p) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		path := filepath.Join(dir, e.Name())
		lines, err := readURLFile(path)
		if err != nil {
			logger.Warn("merge: skip %s: %v", e.Name(), err)
			continue
		}

		added := 0
		for _, l := range lines {
			if _, ok := seen[l]; !ok {
				seen[l] = struct{}{}
				merged = append(merged, l)
				added++
			}
		}
		logger.Info("  %-42s  +%s URLs", e.Name(), utils.FormatN(added))
	}

	if err := utils.WriteLines(outFile, merged); err != nil {
		return 0, fmt.Errorf("write merged: %w", err)
	}

	return len(merged), nil
}

// readURLFile extracts only http(s):// lines using a large buffer.
func readURLFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 8*1024*1024), 8*1024*1024)
	for sc.Scan() {
		l := strings.TrimSpace(sc.Text())
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://") {
			lines = append(lines, l)
		}
	}
	return lines, sc.Err()
}
