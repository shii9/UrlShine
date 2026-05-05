// Package collector orchestrates external URL collection tools concurrently.
package collector

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"urlshine/internal/logger"
	"urlshine/internal/utils"
)

// Config controls collector behaviour.
type Config struct {
	Threads int
	Depth   int
	Subs    bool
	Timeout int // seconds

	RunAll          bool
	RunGau          bool
	RunGospider     bool
	RunKatana       bool
	RunWaymore      bool
	RunWaybackurls  bool
	RunHakrawler    bool
	RunXnlinkfinder bool
	RunGobuster     bool
	RunDirb         bool
}

// DefaultConfig returns production-ready defaults optimized for aggressive collection.
func DefaultConfig() Config {
	return Config{Threads: 50, Depth: 5, Subs: true, Timeout: 60, RunAll: true}
}

// tool pairs a name with its execution function.
type tool struct {
	name string
	fn   func(target, outDir string, cfg Config) ([]string, error)
}

var allTools = []tool{
	{"gau", runGAU},
	{"gospider", runGospider},
	{"katana", runKatana},
	{"waymore", runWaymore},
	{"waybackurls", runWaybackurls},
	{"hakrawler", runHakrawler},
	{"xnLinkFinder", runXnLinkFinder},
	{"gobuster", runGobuster},
	{"dirb", runDirb},
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
			case "gospider":
				selected = cfg.RunGospider
			case "katana":
				selected = cfg.RunKatana
			case "waymore":
				selected = cfg.RunWaymore
			case "waybackurls":
				selected = cfg.RunWaybackurls
			case "hakrawler":
				selected = cfg.RunHakrawler
			case "xnLinkFinder":
				selected = cfg.RunXnlinkfinder
			case "gobuster":
				selected = cfg.RunGobuster
			case "dirb":
				selected = cfg.RunDirb
			}
		}

		// If no tools were explicitly selected, but we are running (RunAll = false, and everything else false)
		// default to all tools if we want, but the user explicitly wants only what they toggle.
		// Wait, if NO tool flags are set, and --all is false, what happens?
		// We'll assume the CLI defaults it to everything if no URL tools are specified. But let's just respect the flags.

		if selected {
			activeTools = append(activeTools, t)
		}
	}

	// If the user provided NO url tools and didn't provide --all, they probably still want all tools if they passed domain?
	// The prompt implies if they do "urlshine google.com", they want all default tools.
	// We'll enforce this in RunAll: if activeTools is empty, use allTools.
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
		mu       sync.Mutex
		outFiles []string
		wg       sync.WaitGroup
		sem      = make(chan struct{}, 10) // max 10 concurrent tool executions for aggressive parallelism
	)

	totalJobs := len(activeTools) * len(targets)
	completedJobs := 0

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
				if completedJobs%10 == 0 {
					logger.Debug("[%d/%d] Collection jobs completed", completedJobs, totalJobs)
				}
				mu.Unlock()
				return
			}

			logger.Run("%-20s → %s", j.tool.name, j.target)
			lines, err := j.tool.fn(j.target, rawDir, cfg)
			if err != nil {
				logger.Warn("%-20s [%s] error: %v", j.tool.name, j.target, err)
				mu.Lock()
				completedJobs++
				mu.Unlock()
				return
			}

			// Keep only HTTP/HTTPS
			clean := make([]string, 0, len(lines))
			for _, l := range lines {
				if strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://") {
					clean = append(clean, l)
				}
			}

			outFile := fmt.Sprintf("%s/%s_%s.txt", rawDir, j.tool.name, utils.SanitizeFilename(j.target))
			if err := utils.WriteLines(outFile, clean); err != nil {
				logger.Error("write %s: %v", outFile, err)
				mu.Lock()
				completedJobs++
				mu.Unlock()
				return
			}

			logger.ToolResult(j.tool.name, j.target, len(clean), false)

			mu.Lock()
			outFiles = append(outFiles, outFile)
			completedJobs++
			if completedJobs%10 == 0 {
				logger.Debug("[%d/%d] Collection jobs completed", completedJobs, totalJobs)
			}
			mu.Unlock()
		}(j)
	}

	wg.Wait()
	return outFiles, nil
}

// ─── Shared Helpers ──────────────────────────────────────────────────────────

func runCmd(args ...string) ([]string, error) {
	logger.Debug("Exec: %s", strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	// cmd.Stderr = os.Stderr // Optional: uncomment for verbose stderr logging
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
	// cmd.Stderr = os.Stderr
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
