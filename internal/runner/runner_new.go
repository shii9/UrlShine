// Package runner coordinates professional URL enumeration with domain-specific organization.
package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"urlshine/internal/alive"
	"urlshine/internal/collector"
	"urlshine/internal/extractor"
	"urlshine/internal/logger"
	"urlshine/internal/normalizer"
	"urlshine/internal/reporter"
	"urlshine/internal/splitter"
	"urlshine/internal/utils"
)

// RunProfessional executes the professional pipeline with domain-specific folders.
func RunProfessional(opts Options) error {
	start := time.Now()
	logger.SetVerbose(opts.Verbose)

	if err := utils.EnsureDir(opts.OutputDir); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	logger.Info("📁 CREATING DOMAIN-SPECIFIC FOLDERS")
	logger.Info("────────────────────────────────────")

	// Create folders for each target
	targetDirs := make(map[string]string)
	for _, target := range opts.Targets {
		folderName := sanitizeDomainName(target) + "_url"
		targetDir := filepath.Join(opts.OutputDir, folderName)
		if err := utils.EnsureDir(targetDir); err != nil {
			return fmt.Errorf("create target dir: %w", err)
		}
		targetDirs[target] = targetDir
		logger.Success("✓ %s/", folderName)
	}

	logger.Info("")
	logger.Info("🚀 PHASE 1: COLLECTING URLs FROM ALL TOOLS")
	logger.Info("────────────────────────────────────────────")
	logger.Info("Configuration: %d threads | Depth: %d | All Tools: %v",
		opts.Threads, opts.Depth, opts.RunAll)
	logger.Info("")

	// Configuration
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

	// Use collector.RunAll() to collect from all targets
	if !opts.SkipCollect {
		rawDir := filepath.Join(opts.OutputDir, "raw")
		_, err := collector.RunAll(opts.Targets, rawDir, cfg)
		if err != nil {
			logger.Error("collection error: %v", err)
		}

		// Organize collected URLs by tool into domain-specific folders
		organizeToolResults(opts.OutputDir, targetDirs)
	}

	logger.Info("")
	logger.Info("📊 PHASE 2: PER-DOMAIN PROCESSING")
	logger.Info("────────────────────────────────────────────")

	stats := reporter.Stats{
		Targets:     opts.Targets,
		OutputDir:   opts.OutputDir,
		Groups:      make(map[string]int),
		AliveGroups: make(map[string]int),
	}

	// Process each domain folder
	for idx, target := range opts.Targets {
		logger.Info("")
		logger.Info("[%d/%d] Processing %s", idx+1, len(opts.Targets), sanitizeDomainName(target)+"_url")
		logger.Info("────────────────────────────────────────────")

		targetDir := targetDirs[target]

		// Step 1: Merge tool files
		logger.Step(1, 4, "Merging URLs from all tools")
		mergedFile := filepath.Join(targetDir, "merged_urls.txt")
		totalMerged := mergeToolsForTarget(targetDir, mergedFile)
		logger.Success("✓ Merged: %s URLs", utils.FormatN(totalMerged))

		// Step 2: Normalize
		logger.Step(2, 4, "Normalizing URLs")
		normFile := filepath.Join(targetDir, "normalized_urls.txt")
		in, out, _ := normalizer.NormalizeFile(mergedFile, normFile)
		if in > 0 {
			reduction := float64((in - out)) / float64(in) * 100
			logger.Success("✓ Normalized: %s URLs (%.1f%% reduction)", utils.FormatN(out), reduction)
		} else {
			logger.Success("✓ Normalized: %s URLs", utils.FormatN(out))
		}

		// Step 3: Advanced Extraction
		logger.Step(3, 4, "Advanced URL Extraction with Specialized Tools")
		logger.Info("")

		lines, _ := utils.ReadLines(normFile)

		// Run advanced extractor for specialized groups
		extr := extractor.NewGroupExtractor(lines, targetDir)
		extractedGroups := extr.ExtractAll()

		// Save extracted groups to files
		logger.Info("")
		logger.Info("  💾 Saving specialized groups:")
		extr.SaveResults(extractedGroups)

		// Also run standard splitter for additional categorization
		logger.Info("")
		logger.Info("  📊 Running standard categorization:")
		groups := splitter.Split(lines)
		_, err := splitter.WriteGroups(groups, targetDir)
		if err != nil {
			logger.Warn("failed to write groups: %v", err)
		}

		counts := splitter.Counts(groups)
		for _, groupName := range []string{"api_urls", "auth_admin_urls", "params_urls", "js_config_urls", "directories_urls"} {
			if count, ok := counts[groupName]; ok && count > 0 {
				displayName := strings.ReplaceAll(groupName, "_", " ")
				logger.Success("  ✓ %-20s: %s URLs", displayName, utils.FormatN(count))
				stats.Groups[groupName] = count
			}
		}

		// Step 4: Optional alive check
		if !opts.SkipAlive {
			logger.Step(4, 4, "Verifying live URLs")
			aliveFile := filepath.Join(targetDir, "alive_urls.txt")
			if _, err := alive.ProbeFile(normFile, aliveFile, opts.Threads); err != nil {
				logger.Warn("alive check failed: %v", err)
			} else {
				aliveCount := utils.FileLineCount(aliveFile)
				logger.Success("✓ Live URLs: %s", utils.FormatN(aliveCount))
				stats.AliveGroups["verified"] = aliveCount
			}
		} else {
			logger.Skip("alive verification skipped")
		}

		logger.Info("")
	}

	logger.Info("")
	logger.Info("✅ COMPLETE")
	logger.Info("────────────────────────────────────────────")
	logger.Info("Results saved to:")
	if len(opts.Targets) == 1 {
		logger.Success("  %s/", sanitizeDomainName(opts.Targets[0])+"_url")
	} else {
		logger.Success("  %s/ (multiple domains)", opts.OutputDir)
	}

	logger.Info("")
	logger.Info("📁 FILES GENERATED (per domain):")
	logger.Info("────────────────────────────────────────────")
	logger.Info("  Per-Tool Results:")
	logger.Info("    • gau.txt                 - GAU tool URLs")
	logger.Info("    • katana.txt              - Katana crawler URLs")
	logger.Info("    • gospider.txt            - GoSpider URLs")
	logger.Info("    • waymore.txt             - Waymore URLs")
	logger.Info("    • waybackurls.txt         - Wayback URLs")
	logger.Info("    • hakrawler.txt           - Hakrawler URLs")
	logger.Info("    • xnlinkfinder.txt        - xnLinkFinder URLs")
	logger.Info("    • gobuster.txt            - Gobuster directory discovery")
	logger.Info("    • dirb.txt                - Dirb directory brute-force")
	logger.Info("")
	logger.Info("  Processed Results:")
	logger.Info("    • merged_urls.txt         - All URLs combined")
	logger.Info("    • normalized_urls.txt     - Cleaned & deduplicated")
	logger.Info("    • api_urls.txt            - API endpoints")
	logger.Info("    • auth_admin_urls.txt     - Auth & admin pages")
	logger.Info("    • params_urls.txt         - URLs with parameters")
	logger.Info("    • js_config_urls.txt      - JS & config files")
	logger.Info("    • directories_urls.txt    - Directory paths")
	if !opts.SkipAlive {
		logger.Info("    • alive_urls.txt          - Verified live URLs")
	}

	logger.Info("")
	dur := time.Since(start)
	logger.Success("Total time: %s", formatDuration(dur))
	logger.Info("")

	return nil
}

// organizeToolResults moves tool results from raw/ to domain-specific folders
func organizeToolResults(outputDir string, targetDirs map[string]string) {
	rawDir := filepath.Join(outputDir, "raw")
	if _, err := os.Stat(rawDir); os.IsNotExist(err) {
		return // raw dir doesn't exist
	}

	// Read raw directory
	entries, err := os.ReadDir(rawDir)
	if err != nil {
		return
	}

	// For each target domain, organize its tool files
	for target, targetDir := range targetDirs {
		sanitized := sanitizeDomainName(target)

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			// Look for files named like: google.com-gau.txt, google.com-katana.txt, etc.
			if strings.HasPrefix(name, sanitized+"-") {
				toolName := strings.TrimPrefix(name, sanitized+"-")
				srcPath := filepath.Join(rawDir, name)
				dstPath := filepath.Join(targetDir, toolName)

				// Copy file
				data, _ := os.ReadFile(srcPath)
				os.WriteFile(dstPath, data, 0644)
			}
		}
	}
}

// mergeToolsForTarget merges all tool files for a single target.
func mergeToolsForTarget(targetDir, outFile string) int {
	tools := []string{"gau.txt", "katana.txt", "gospider.txt", "waymore.txt", "waybackurls.txt", "hakrawler.txt", "xnlinkfinder.txt", "gobuster.txt", "dirb.txt"}
	seen := make(map[string]struct{})
	var merged []string

	for _, tool := range tools {
		toolPath := filepath.Join(targetDir, tool)
		if _, err := os.Stat(toolPath); os.IsNotExist(err) {
			continue // Skip if file doesn't exist
		}

		lines, err := utils.ReadLines(toolPath)
		if err != nil {
			continue
		}

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if _, ok := seen[line]; !ok {
				seen[line] = struct{}{}
				merged = append(merged, line)
			}
		}
	}

	utils.WriteLines(outFile, merged)
	return len(merged)
}

// sanitizeDomainName converts a domain to a safe folder name
func sanitizeDomainName(domain string) string {
	// Replace dots and special chars with underscores, lowercase
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
	// Remove leading/trailing underscores
	name = strings.Trim(name, "_")
	return name
}

// formatDuration formats a duration in human-readable format
func formatDuration(d time.Duration) string {
	if d.Seconds() < 60 {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm %ds", minutes, seconds)
}
