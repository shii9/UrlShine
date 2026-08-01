// Package collector orchestrates external URL collection tools concurrently.
package collector

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/shii9/UrlShine/internal/logger"
	"github.com/shii9/UrlShine/internal/utils"
)

// Config controls collector behaviour.
type Config struct {
	Threads int
	Depth   int
	Subs    bool
	Timeout int // seconds per command

	// Core tool selection (10 high-value, non-redundant tools)
	RunAll          bool
	RunGau          bool
	RunWaymore      bool
	RunParamspider  bool
	RunCommoncrawl  bool
	RunUrlfinder    bool
	RunGithubEndpoints bool
	RunXnlinkfinder bool
	RunKatana       bool
	RunHakrawler    bool
	RunGobuster     bool
}

// DefaultConfig returns production-ready defaults optimized for aggressive collection.
func DefaultConfig() Config {
	return Config{
		Threads: 50,
		Depth:   5,
		Subs:    true,
		Timeout: 300,

		RunAll:             false,
		RunGau:             true,
		RunWaymore:         true,
		RunParamspider:     true,
		RunCommoncrawl:     true,
		RunUrlfinder:       true,
		RunGithubEndpoints: true,
		RunXnlinkfinder:    true,
		RunKatana:          true,
		RunHakrawler:       true,
		RunGobuster:        true,
	}
}

// tool pairs a name with its execution function.
type tool struct {
	name string
	fn   func(target, outDir string, cfg Config) ([]string, error)
}

// allTools defines the 10 distinct, non-redundant URL collection engines in 4 operational tiers:
//
//   Tier 1 — Passive Archives (Wayback, CommonCrawl, OTX, URLScan, parameter mining)
//   Tier 2 — Passive APIs & OSINT (Native CDX, crt.sh/VirusTotal, GitHub leaks, JS link finder)
//   Tier 3 — Active Crawlers (Headless JS Chromium crawler, fast Go crawler)
//   Tier 4 — Active Brute-Force (High-speed directory/file brute-force)
var allTools = []tool{
	// ─── Tier 1: Passive Archives ─────────────────────────────────────────
	{"gau", runGAU},                 // Wayback, CommonCrawl, URLScan, OTX
	{"waymore", runWaymore},         // Enhanced Wayback Machine queries & response parsing
	{"paramspider", runParamspider}, // Parameter-focused archive mining

	// ─── Tier 2: Passive APIs & OSINT ─────────────────────────────────────
	{"commoncrawl", runCommoncrawl},         // Common Crawl CDX API (native Go, zero dependency)
	{"urlfinder", runUrlfinder},             // crt.sh, VirusTotal, passive intel
	{"github-endpoints", runGithubEndpoints}, // GitHub repository endpoint leaks
	{"xnLinkFinder", runXnLinkFinder},       // JS/HTML link & endpoint extraction

	// ─── Tier 3: Active Crawlers ──────────────────────────────────────────
	{"katana", runKatana},       // Advanced JS-capable crawler (headless Chromium)
	{"hakrawler", runHakrawler}, // Fast Go web crawler, HTML/JS parsing

	// ─── Tier 4: Active Brute-Force ───────────────────────────────────────
	{"gobuster", runGobuster}, // Fast multi-threaded directory/file brute-force
}

// RunAll executes every installed tool against every target concurrently.
// Writes per-tool, per-target files to rawDir. Returns list of written file paths.
func RunAll(targets []string, rawDir string, cfg Config) ([]string, error) {
	if err := utils.EnsureDir(rawDir); err != nil {
		return nil, fmt.Errorf("create raw dir: %w", err)
	}

	// Tool availability and selection matrix
	var activeTools []tool
	for _, t := range allTools {
		selected := cfg.RunAll
		if !selected {
			switch t.name {
			case "gau":
				selected = cfg.RunGau
			case "waymore":
				selected = cfg.RunWaymore
			case "paramspider":
				selected = cfg.RunParamspider
			case "commoncrawl":
				selected = cfg.RunCommoncrawl
			case "urlfinder":
				selected = cfg.RunUrlfinder
			case "github-endpoints":
				selected = cfg.RunGithubEndpoints
			case "xnLinkFinder":
				selected = cfg.RunXnlinkfinder
			case "katana":
				selected = cfg.RunKatana
			case "hakrawler":
				selected = cfg.RunHakrawler
			case "gobuster":
				selected = cfg.RunGobuster
			}
		}

		if selected {
			activeTools = append(activeTools, t)
		}
	}

	// If no tools selected and --all not set, default to all tools
	if len(activeTools) == 0 {
		activeTools = allTools
	}

	toolStatus := make([]struct {
		Name  string
		Found bool
	}, len(activeTools))
	for i, t := range activeTools {
		toolStatus[i] = struct {
			Name  string
			Found bool
		}{Name: t.name, Found: utils.ToolExists(t.name)}
	}
	logger.ToolMatrix(toolStatus)

	type job struct {
		tool   tool
		target string
	}

	jobs := make(chan job)
	go func() {
		for _, t := range activeTools {
			for _, tgt := range targets {
				jobs <- job{t, tgt}
			}
		}
		close(jobs)
	}()

	var (
		mu             sync.Mutex
		outFiles       []string
		wg             sync.WaitGroup
		sem            = make(chan struct{}, cfg.Threads)
		completedJobs  = 0
		targetProgress = make(map[string]int)
	)

	// Initialize per-target tracking
	for _, t := range targets {
		targetProgress[t] = 0
	}

	for j := range jobs {
		wg.Add(1)
		go func(j job) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if !utils.ToolExists(j.tool.name) {
				logger.ToolResult(j.tool.name, j.target, 0, true)
				mu.Lock()
				completedJobs++
				targetProgress[j.target]++
				mu.Unlock()
				return
			}

			logger.RunWithSpinner(j.tool.name, j.target)
			lines, err := j.tool.fn(j.target, rawDir, cfg)
			if err != nil {
				logger.Warn("%-20s [%s] failed", j.tool.name, j.target)
				mu.Lock()
				completedJobs++
				targetProgress[j.target]++
				mu.Unlock()
				return
			}

			// Deduplicate and keep only HTTP/HTTPS
			unique := make(map[string]struct{})
			clean := make([]string, 0)
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if (strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://")) && l != "" {
					if _, ok := unique[l]; !ok {
						unique[l] = struct{}{}
						clean = append(clean, l)
					}
				}
			}

			outFile := fmt.Sprintf("%s/%s_%s.txt", rawDir, j.tool.name, utils.SanitizeFilename(j.target))
			if err := utils.WriteLines(outFile, clean); err != nil {
				logger.Error("write %s: %v", outFile, err)
				mu.Lock()
				completedJobs++
				targetProgress[j.target]++
				mu.Unlock()
				return
			}

			logger.ToolResult(j.tool.name, j.target, len(clean), false)

			mu.Lock()
			outFiles = append(outFiles, outFile)
			completedJobs++
			targetProgress[j.target]++
			mu.Unlock()
		}(j)
	}

	wg.Wait()
	return outFiles, nil
}

// ─── Shared Helpers ──────────────────────────────────────────────────────────

func runCmd(args ...string) ([]string, error) {
	logger.Debug("Exec: %s", strings.Join(args, " "))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	var lines []string
	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)
	for sc.Scan() {
		if l := strings.TrimSpace(sc.Text()); l != "" {
			lines = append(lines, l)
		}
	}
	if err := sc.Err(); err != nil {
		_ = cmd.Wait()
		return lines, err
	}
	if err := cmd.Wait(); err != nil {
		return lines, err
	}
	return lines, nil
}

func runCmdStdin(input string, args ...string) ([]string, error) {
	logger.Debug("Exec(stdin): %s", strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = strings.NewReader(input)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	var lines []string
	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)
	for sc.Scan() {
		if l := strings.TrimSpace(sc.Text()); l != "" {
			lines = append(lines, l)
		}
	}
	if err := sc.Err(); err != nil {
		_ = cmd.Wait()
		return lines, err
	}
	if err := cmd.Wait(); err != nil {
		return lines, err
	}
	return lines, nil
}

func ensureHTTPS(target string) string {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		return target
	}
	return "https://" + target
}
