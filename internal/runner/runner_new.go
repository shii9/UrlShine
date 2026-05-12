// Package runner coordinates professional URL enumeration with domain-specific organization.
package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shii9/UrlShine/internal/alive"
	"github.com/shii9/UrlShine/internal/banner"
	"github.com/shii9/UrlShine/internal/collector"
	"github.com/shii9/UrlShine/internal/extractor"
	"github.com/shii9/UrlShine/internal/installer"
	"github.com/shii9/UrlShine/internal/logger"
	"github.com/shii9/UrlShine/internal/normalizer"
	"github.com/shii9/UrlShine/internal/reporter"
	"github.com/shii9/UrlShine/internal/splitter"
	"github.com/shii9/UrlShine/internal/utils"
)

var toolOutputNames = map[string]string{
	"gau":          "gau.txt",
	"gospider":     "gospider.txt",
	"katana":       "katana.txt",
	"waymore":      "waymore.txt",
	"waybackurls":  "waybackurls.txt",
	"hakrawler":    "hakrawler.txt",
	"xnLinkFinder": "xnlinkfinder.txt",
	"gobuster":     "gobuster.txt",
	"dirb":         "dirb.txt",
}

// RunProfessional executes collection and, when requested, the complete processing pipeline.
func RunProfessional(opts Options) error {
	start := time.Now()
	logger.SetVerbose(opts.Verbose)

	// Display professional banner
	banner.Print()

	// Verify and check for missing tools
	installer.CheckAndInstall()

	if err := utils.EnsureDir(opts.OutputDir); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	logger.Info("Creating domain-specific folders")
	targetDirs := make(map[string]string, len(opts.Targets))
	for _, target := range opts.Targets {
		folderName := sanitizeDomainName(target) + "_url"
		targetDir := filepath.Join(opts.OutputDir, folderName)
		if err := utils.EnsureDir(targetDir); err != nil {
			return fmt.Errorf("create target dir: %w", err)
		}
		targetDirs[target] = targetDir
		logger.Success("%s/", folderName)
	}

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
		RunHakrawler:    opts.RunHakrawler,
		RunXnlinkfinder: opts.RunXnlinkfinder,
		RunGobuster:     opts.RunGobuster,
		RunDirb:         opts.RunDirb,
	}

	logger.Info("")
	logger.Info("Phase 1: URL collection")
	logger.Info("Configuration: %d threads | depth %d | all tools %v | complete %v",
		opts.Threads, opts.Depth, opts.RunAll, opts.RunComplete)

	logger.SectionHeader("External Tools Status")

	rawDir := filepath.Join(opts.OutputDir, "raw")
	if !opts.SkipCollect {
		logger.SectionHeader("Collecting URLs from All Tools")
		if _, err := collector.RunAll(opts.Targets, rawDir, cfg); err != nil {
			logger.Warn("collection error: %v", err)
		}
		logger.BlankLine()
	} else {
		logger.Skip("collection skipped; using existing per-tool or raw files")
	}

	copied, err := organizeToolResults(rawDir, targetDirs)
	if err != nil {
		logger.Warn("organizing raw tool results failed: %v", err)
	}
	if copied > 0 {
		logger.Success("organized %d per-tool result files", copied)
	}

	if !opts.RunComplete {
		logger.Info("")
		logger.SectionHeader("Collection Complete")
		logger.Info("Use -complete flag to run: merge, normalize, categorize, alive-check, and reports")
		printOutputSummary(opts, start, false)
		return nil
	}

	logger.Info("")
	logger.SectionHeader("Phase 2: Processing Pipeline")

	stats := reporter.Stats{
		Targets:     opts.Targets,
		OutputDir:   opts.OutputDir,
		Groups:      make(map[string]int),
		AliveGroups: make(map[string]int),
	}
	paramKeys := make(map[string]struct{})

	for idx, target := range opts.Targets {
		targetDir := targetDirs[target]
		logger.BlankLine()
		logger.SectionHeader(fmt.Sprintf("Target %d/%d: %s", idx+1, len(opts.Targets), target))

		logger.Step(1, 4, "Merging and deduplicating")
		mergedFile := filepath.Join(targetDir, "merged_urls.txt")
		totalMerged, err := mergeToolsForTarget(targetDir, mergedFile)
		if err != nil {
			return fmt.Errorf("merge %s: %w", target, err)
		}
		stats.TotalRaw += totalMerged
		logger.Success("Merged: %s URLs", utils.FormatN(totalMerged))

		logger.Step(2, 4, "Normalizing URLs")
		normFile := filepath.Join(targetDir, "normalized_urls.txt")
		in, out, err := normalizer.NormalizeFile(mergedFile, normFile)
		if err != nil {
			return fmt.Errorf("normalize %s: %w", target, err)
		}
		stats.AfterNorm += out
		logger.Success("Normalized: %s URLs (%.1f%% reduction)", utils.FormatN(out), utils.Reduction(in, out))

		logger.Step(3, 4, "Categorizing attack groups")
		lines, err := utils.ReadLines(normFile)
		if err != nil {
			return fmt.Errorf("read normalized urls for %s: %w", target, err)
		}

		extr := extractor.NewGroupExtractor(lines, targetDir)
		extractedGroups := extr.ExtractAll()
		if err := extr.SaveResults(extractedGroups); err != nil {
			logger.Warn("saving extracted groups failed: %v", err)
		}

		groups := splitter.Split(lines)
		if _, err := splitter.WriteGroups(groups, targetDir); err != nil {
			logger.Warn("failed to write categorized groups: %v", err)
		}
		for _, key := range splitter.ParamKeys(groups.Params) {
			paramKeys[key] = struct{}{}
		}

		for _, groupName := range []string{
			splitter.GroupAPI,
			splitter.GroupAuth,
			splitter.GroupParams,
			splitter.GroupJS,
			splitter.GroupDirs,
		} {
			count := splitter.Counts(groups)[groupName]
			stats.Groups[groupName] += count
			logger.Success("%-20s %s URLs", strings.ReplaceAll(groupName, "_", " "), utils.FormatN(count))
		}

		if !opts.SkipAlive {
			logger.Step(4, 4, "Alive checking")
			aliveFile := filepath.Join(targetDir, "alive_urls.txt")
			if _, err := alive.ProbeFile(normFile, aliveFile, opts.Threads); err != nil {
				logger.Warn("alive check failed: %v", err)
			} else {
				aliveCount := utils.FileLineCount(aliveFile)
				stats.AliveGroups["verified"] += aliveCount
				logger.Success("Live URLs: %s", utils.FormatN(aliveCount))
			}
		} else {
			logger.Skip("alive verification skipped")
		}
	}

	stats.UniqueParams = len(paramKeys)
	stats.DurationSec = time.Since(start).Seconds()
	if err := reporter.WriteReports(stats); err != nil {
		logger.Warn("failed to write reports: %v", err)
	} else {
		logger.Success("JSON and Markdown reports generated")
	}

	logger.Info("")
	logger.Success("Complete pipeline finished")
	printOutputSummary(opts, start, true)
	return nil
}

func printOutputSummary(opts Options, start time.Time, complete bool) {
	logger.Info("")
	logger.Info("Results saved to:")
	if len(opts.Targets) == 1 {
		logger.Success("  %s/", sanitizeDomainName(opts.Targets[0])+"_url")
	} else {
		logger.Success("  %s/ (multiple domains)", opts.OutputDir)
	}

	logger.Info("")
	logger.Info("Files generated per domain:")
	logger.Info("  Per-tool results: gau.txt, katana.txt, gospider.txt, waymore.txt, waybackurls.txt")
	logger.Info("                    hakrawler.txt, xnlinkfinder.txt, gobuster.txt, dirb.txt")
	if complete {
		logger.Info("  Complete outputs: merged_urls.txt, normalized_urls.txt, api_urls.txt")
		logger.Info("                    auth_admin_urls.txt, params_urls.txt, js_config_urls.txt")
		logger.Info("                    directories_urls.txt")
		if !opts.SkipAlive {
			logger.Info("                    alive_urls.txt")
		}
		logger.Info("  Reports:          urlshine_report.json, urlshine_report.md")
	}
	logger.Success("Total time: %s", formatDuration(time.Since(start)))
	logger.Info("")
}

func organizeToolResults(rawDir string, targetDirs map[string]string) (int, error) {
	if _, err := os.Stat(rawDir); os.IsNotExist(err) {
		return 0, nil
	}

	copied := 0
	for target, targetDir := range targetDirs {
		targetPart := utils.SanitizeFilename(target)
		legacyTargetPart := sanitizeDomainName(target)

		for toolName, dstName := range toolOutputNames {
			candidates := []string{
				filepath.Join(rawDir, toolName+"_"+targetPart+".txt"),
				filepath.Join(rawDir, legacyTargetPart+"-"+dstName),
			}

			for _, srcPath := range candidates {
				if !utils.FileExists(srcPath) {
					continue
				}
				data, err := os.ReadFile(srcPath)
				if err != nil {
					return copied, err
				}
				if err := os.WriteFile(filepath.Join(targetDir, dstName), data, 0644); err != nil {
					return copied, err
				}
				copied++
				break
			}
		}
	}
	return copied, nil
}

func mergeToolsForTarget(targetDir, outFile string) (int, error) {
	tools := []string{
		"gau.txt",
		"katana.txt",
		"gospider.txt",
		"waymore.txt",
		"waybackurls.txt",
		"hakrawler.txt",
		"xnlinkfinder.txt",
		"gobuster.txt",
		"dirb.txt",
	}
	seen := make(map[string]struct{})
	var merged []string

	for _, tool := range tools {
		toolPath := filepath.Join(targetDir, tool)
		if !utils.FileExists(toolPath) {
			continue
		}

		lines, err := utils.ReadLines(toolPath)
		if err != nil {
			return 0, err
		}
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if _, ok := seen[line]; ok {
				continue
			}
			seen[line] = struct{}{}
			merged = append(merged, line)
		}
	}

	if err := utils.WriteLines(outFile, merged); err != nil {
		return 0, err
	}
	return len(merged), nil
}

func sanitizeDomainName(domain string) string {
	replacer := strings.NewReplacer(
		".", "_",
		"/", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	name := strings.ToLower(replacer.Replace(domain))
	return strings.Trim(name, "_")
}

func formatDuration(d time.Duration) string {
	if d.Seconds() < 60 {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm %ds", minutes, seconds)
}
