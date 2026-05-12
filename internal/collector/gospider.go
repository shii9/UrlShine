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
	tmpOut := filepath.Join(outDir, "_gospider_"+utils.SanitizeFilename(target))
	_ = os.MkdirAll(tmpOut, 0755)

	// Set aggressive concurrency
	concurrency := cfg.Threads
	if concurrency < 30 {
		concurrency = 30
	}
	depth := cfg.Depth
	if depth < 3 {
		depth = 3
	}

	args := []string{
		"gospider",
		"-s", ensureHTTPS(target),
		"-c", fmt.Sprintf("%d", concurrency),
		"-d", fmt.Sprintf("%d", depth),
		"-o", tmpOut,
		"--js", "--sitemap", "--robots", "--other-source",
		"-a", "-w", "-r",
		"--blacklist", `\.(png|jpg|jpeg|gif|bmp|svg|ico|webp|css|woff|woff2|eot|ttf|pdf|zip|rar|tar|gz|mp4|mp3|avi|webm|mkv|mov|flv|swf)$`,
	}
	if cfg.Subs {
		args = append(args, "--subs")
	}

	_, _ = runCmd(args...)
	return gospiderParseDir(tmpOut)
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
