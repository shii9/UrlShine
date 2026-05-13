package collector

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shii9/UrlShine/internal/utils"
)

// runGospider collects URLs via gospider with aggressive parameters.
func runGospider(target, outDir string, cfg Config) ([]string, error) {
	target = ensureHTTPS(target)
	var allUrls []string

	// Professional bug hunter command sequences to maximize coverage
	commands := [][]string{
		// 1. Broad Discovery (Baseline)
		{"gospider", "-s", target, "-c", "20", "-d", "3", "--subs", "--other-source", "--include-subs", "--include-other-source", "--sitemap", "--robots", "--js", "-q"},
		// 2. Deep Crawl (Infinite recursion)
		{"gospider", "-s", target, "-c", "30", "-d", "0", "--subs", "--other-source", "--include-subs", "--include-other-source", "--sitemap", "--robots", "--js", "-q"},
		// 3. Fast Third-Party Expansion
		{"gospider", "-s", target, "-c", "20", "-d", "1", "--other-source", "--include-subs", "--include-other-source", "-q"},
		// 4. Focused Crawl (Static Asset Filtered)
		{"gospider", "-s", target, "-c", "20", "-d", "3", "--subs", "--other-source", "--include-subs", "--include-other-source", "--blacklist", `\.(png|jpg|jpeg|gif|css|svg|woff|woff2|ttf|ico)$`, "-q"},
	}

	for i, args := range commands {
		// Create unique tmp dir for each run to avoid file conflicts
		tmpOut := filepath.Join(outDir, fmt.Sprintf("_gospider_%s_%d", utils.SanitizeFilename(target), i))
		_ = os.MkdirAll(tmpOut, 0755)

		_, _ = runCmd(args...)

		// Parse results from this run
		lines, err := gospiderParseDir(tmpOut)
		if err == nil {
			allUrls = append(allUrls, lines...)
		}

		// Cleanup tmp dir after parsing
		os.RemoveAll(tmpOut)
	}

	return allUrls, nil
}

func gospiderParseDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fl, _ := gospiderParseFile(filepath.Join(dir, e.Name()))
		lines = append(lines, fl...)
	}
	return lines, nil
}

func gospiderParseFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if idx := strings.Index(line, " - "); idx != -1 {
			line = strings.TrimSpace(line[idx+3:])
		}
		if idx := strings.Index(line, " ["); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			lines = append(lines, line)
		}
	}
	return lines, sc.Err()
}
