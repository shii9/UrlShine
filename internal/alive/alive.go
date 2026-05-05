// Package alive verifies URL availability using httpx or native Go HTTP.
package alive

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"urlshine/internal/logger"
	"urlshine/internal/utils"
)

// Result holds live URL metadata.
type Result struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Length     int64  `json:"length"`
	Title      string `json:"title,omitempty"`
	Tech       string `json:"tech,omitempty"`
}

func (r Result) String() string {
	return fmt.Sprintf("%-6d  %-80s  %d", r.StatusCode, r.URL, r.Length)
}

// ProbeFile reads URLs from inFile, checks them, and writes results.
func ProbeFile(inFile, outFile string, threads int) ([]Result, error) {
	if utils.ToolExists("httpx") {
		return probeHTTPX(inFile, outFile, threads)
	}
	logger.Warn("httpx not found in PATH — using built-in Go prober")
	return probeNative(inFile, outFile, threads)
}

func probeHTTPX(inFile, outFile string, threads int) ([]Result, error) {
	args := []string{
		"httpx",
		"-l", inFile,
		"-threads", strconv.Itoa(threads),
		"-silent",
		"-status-code",
		"-title",
		"-content-length",
		"-tech-detect",
		"-o", outFile,
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("httpx: %w", err)
	}

	lines, err := utils.ReadLines(outFile)
	if err != nil {
		return nil, err
	}
	var results []Result
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		r := Result{URL: parts[0]}
		for _, p := range parts[1:] {
			p = strings.Trim(p, "[]")
			if code, err := strconv.Atoi(p); err == nil {
				r.StatusCode = code
			}
		}
		results = append(results, r)
	}
	return results, nil
}

func probeNative(inFile, outFile string, threads int) ([]Result, error) {
	urls, err := utils.ReadLines(inFile)
	if err != nil {
		return nil, err
	}
	if len(urls) == 0 {
		return nil, nil
	}

	bar := progressbar.NewOptions(len(urls),
		progressbar.OptionSetDescription("  probing"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() { fmt.Println() }),
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	type job struct{ url string }
	jobs := make(chan job, len(urls))
	for _, u := range urls {
		jobs <- job{u}
	}
	close(jobs)

	resultsCh := make(chan Result, len(urls))
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				req, err := http.NewRequestWithContext(ctx, "GET", j.url, nil)
				if err == nil {
					req.Header.Set("User-Agent", "Mozilla/5.0 (URLShine/2.0)")
					resp, err := client.Do(req)
					if err == nil {
						resultsCh <- Result{
							URL:        j.url,
							StatusCode: resp.StatusCode,
							Length:     resp.ContentLength,
						}
						resp.Body.Close()
					}
				}
				cancel()
				_ = bar.Add(1)
			}
		}()
	}
	wg.Wait()
	close(resultsCh)

	var results []Result
	for r := range resultsCh {
		results = append(results, r)
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].StatusCode != results[j].StatusCode {
			return results[i].StatusCode < results[j].StatusCode
		}
		return results[i].URL < results[j].URL
	})

	f, err := os.Create(outFile)
	if err != nil {
		return results, err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, r := range results {
		fmt.Fprintln(w, r.String())
	}
	return results, w.Flush()
}

// ProbeGroups iterates over categorised files and probes each.
func ProbeGroups(files map[string]string, aliveDir string, threads int) (map[string]string, error) {
	if err := utils.EnsureDir(aliveDir); err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for group, path := range files {
		count := utils.FileLineCount(path)
		if count == 0 {
			logger.Skip("%s — empty file", group)
			continue
		}
		logger.Run("%-26s  %d URLs ...", group, count)
		stem := strings.TrimSuffix(filepath.Base(path), ".txt")
		outFile := filepath.Join(aliveDir, stem+"_alive.txt")
		res, err := ProbeFile(path, outFile, threads)
		if err != nil {
			logger.Warn("alive %s: %v", group, err)
			continue
		}
		logger.Success("%-26s  %d live", group, len(res))
		out[group] = outFile
	}
	return out, nil
}
