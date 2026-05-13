// Package runner coordinates the full reconnaissance pipeline with per-tool tracking.
// It orchestrates the 5-stage processing pipeline: collection, merging, normalization,
// categorization, and live verification, with comprehensive logging at each stage.
package runner

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/shii9/UrlShine/internal/alive"
	"github.com/shii9/UrlShine/internal/collector"
	"github.com/shii9/UrlShine/internal/logger"
	"github.com/shii9/UrlShine/internal/merger"
	"github.com/shii9/UrlShine/internal/normalizer"
	"github.com/shii9/UrlShine/internal/reporter"
	"github.com/shii9/UrlShine/internal/splitter"
	"github.com/shii9/UrlShine/internal/utils"
)

// Options defines all configurable parameters for the reconnaissance pipeline.
// This structure is populated from CLI flags and controls pipeline behavior.
type Options struct {
	// Target configuration
	Targets   []string // List of domains/targets to scan
	OutputDir string   // Output directory for results

	// Performance tuning
	Threads int  // Number of parallel threads (default: 50)
	Depth   int  // Crawl depth for active tools (default: 3)
	Subs    bool // Include subdomains in scans

	// Pipeline control
	SkipAlive   bool // Skip live host verification stage
	SkipCollect bool // Skip collection, reprocess existing data
	Verbose     bool // Enable verbose debug logging
	RunComplete bool // Run full pipeline including post-collection processing

	// Collection tool selection
	RunAll          bool // Run all available tools
	RunGau          bool // Enable GetAllUrls
	RunGospider     bool // Enable GoSpider
	RunKatana       bool // Enable Katana
	RunWaymore      bool // Enable Waymore
	RunWaybackurls  bool // Enable Wayback URLs
	RunXnlinkfinder bool // Enable xnLinkFinder
	RunGobuster     bool // Enable Gobuster
	RunDirb         bool // Enable Dirb
}

// Run executes the URLShine reconnaissance pipeline.
// It orchestrates a 5-stage process for comprehensive URL enumeration.
func Run(opts Options) error {
	start := time.Now()
	logger.SetVerbose(opts.Verbose)

	if err := utils.EnsureDir(opts.OutputDir); err != nil {
		return fmt.Errorf("create out dir: %w", err)
	}

	rawDir := filepath.Join(opts.OutputDir, "raw")
	mergedFile := filepath.Join(opts.OutputDir, "merged_raw.txt")
	normFile := filepath.Join(opts.OutputDir, "normalized_urls.txt")

	stats := reporter.Stats{
		Targets:     opts.Targets,
		OutputDir:   opts.OutputDir,
		Groups:      make(map[string]int),
		AliveGroups: make(map[string]int),
	}

	// 1. COLLECT
	logger.Step(1, 5, "URL Collection")
	if !opts.SkipCollect {
		cfg := collector.Config{
			Threads:         opts.Threads,
			Depth:           opts.Depth,
			Subs:            opts.Subs,
			RunAll:          opts.RunAll,
			RunGau:          opts.RunGau,
			RunGospider:     opts.RunGospider,
			RunKatana:       opts.RunKatana,
			RunWaymore:      opts.RunWaymore,
			RunWaybackurls:  opts.RunWaybackurls,
			RunXnlinkfinder: opts.RunXnlinkfinder,
		}
		_, err := collector.RunAll(opts.Targets, rawDir, cfg)
		if err != nil {
			logger.Error("collection: %v", err)
		}
	} else {
		logger.Skip("skipping collection phase")
	}

	// 2. MERGE
	logger.Step(2, 5, "Merging Results")
	totalRaw, err := merger.MergeDir(rawDir, mergedFile)
	if err != nil {
		return err
	}
	stats.TotalRaw = totalRaw
	logger.Success("merged_raw.txt → %s URLs", utils.FormatN(totalRaw))

	// 3. NORMALIZE
	logger.Step(3, 5, "URL Normalization")
	in, out, err := normalizer.NormalizeFile(mergedFile, normFile)
	if err != nil {
		return err
	}
	stats.AfterNorm = out
	logger.Success("normalized_urls.txt → %s URLs (%.0f%% reduced)", utils.FormatN(out), utils.Reduction(in, out))

	// 4. SPLIT
	logger.Step(4, 5, "Logical Splitting")
	lines, err := utils.ReadLines(normFile)
	if err != nil {
		return err
	}
	groups := splitter.Split(lines)
	groupFiles, err := splitter.WriteGroups(groups, opts.OutputDir)
	if err != nil {
		return err
	}
	paramKeys := splitter.ParamKeys(groups.Params)
	_ = utils.WriteLines(filepath.Join(opts.OutputDir, "unique_params.txt"), paramKeys)
	stats.UniqueParams = len(paramKeys)

	for k, v := range splitter.Counts(groups) {
		stats.Groups[k] = v
		logger.Success("%-26s %s URLs", k, utils.FormatN(v))
	}

	// 5. ALIVE
	logger.Step(5, 5, "Alive Checking")
	if !opts.SkipAlive {
		probes := make(map[string]string)
		for k, p := range groupFiles {
			probes[k] = p
		}
		probes["normalized_urls"] = normFile
		aliveFiles, err := alive.ProbeGroups(probes, filepath.Join(opts.OutputDir, "alive"), opts.Threads)
		if err != nil {
			logger.Warn("alive check failed: %v", err)
		}
		for label, p := range aliveFiles {
			stats.AliveGroups[label] = utils.FileLineCount(p)
		}
	} else {
		logger.Skip("alive checks skipped")
	}

	// 6. REPORT
	logger.Success("Generating Final Report")
	dur := time.Since(start)
	stats.DurationSec = dur.Seconds()

	if err := reporter.WriteReports(stats); err != nil {
		logger.Warn("failed to write reports: %v", err)
	} else {
		logger.Success("JSON & Markdown reports generated")
	}

	reporter.Print(stats, dur)
	return nil
}
