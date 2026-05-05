// Package collector provides URL aggregation by tool.
package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"urlshine/internal/logger"
	"urlshine/internal/utils"
)

// ToolAggregator tracks and aggregates URLs collected by each tool.
type ToolAggregator struct {
	mu    sync.Mutex
	Calls map[string][]string // tool name -> collected URLs
}

// NewToolAggregator creates a new aggregator.
func NewToolAggregator() *ToolAggregator {
	return &ToolAggregator{
		Calls: make(map[string][]string),
	}
}

// AddURLs adds URLs from a specific tool.
func (ta *ToolAggregator) AddURLs(toolName string, urls []string) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.Calls[toolName] = append(ta.Calls[toolName], urls...)
}

// WriteToolFiles writes aggregated URLs to per-tool files in output directory.
func (ta *ToolAggregator) WriteToolFiles(outputDir string) map[string]int {
	ta.mu.Lock()
	defer ta.mu.Unlock()

	counts := make(map[string]int)

	for toolName, urls := range ta.Calls {
		if len(urls) == 0 {
			continue
		}

		// Deduplicate
		seen := make(map[string]struct{})
		var unique []string
		for _, url := range urls {
			if _, ok := seen[url]; !ok {
				seen[url] = struct{}{}
				unique = append(unique, url)
			}
		}

		// Write file
		toolFile := filepath.Join(outputDir, fmt.Sprintf("%s.txt", toolName))
		if err := utils.WriteLines(toolFile, unique); err != nil {
			logger.Warn("Failed to write %s.txt: %v", toolName, err)
			continue
		}

		counts[toolName] = len(unique)
		logger.Success("%-20s: %s URLs saved", toolName, utils.FormatN(len(unique)))
	}

	return counts
}

// MergeToolFiles merges all per-tool files into a single merged file.
func MergeToolFiles(outputDir, mergedFile string) (int, error) {
	if _, err := os.Stat(outputDir); err != nil {
		return 0, fmt.Errorf("output dir: %w", err)
	}

	seen := make(map[string]struct{})
	var merged []string

	// Only read .txt files that are tool results (not raw, not merged, not groups)
	toolFiles := []string{"gau.txt", "katana.txt", "gospider.txt", "waymore.txt", "waybackurls.txt", "hakrawler.txt", "xnlinkfinder.txt"}

	for _, toolFile := range toolFiles {
		path := filepath.Join(outputDir, toolFile)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue // Skip if file doesn't exist
		}

		lines, err := utils.ReadLines(path)
		if err != nil {
			logger.Warn("Failed to read %s: %v", toolFile, err)
			continue
		}

		for _, line := range lines {
			if _, ok := seen[line]; !ok {
				seen[line] = struct{}{}
				merged = append(merged, line)
			}
		}
	}

	if err := utils.WriteLines(mergedFile, merged); err != nil {
		return 0, fmt.Errorf("write merged: %w", err)
	}

	return len(merged), nil
}
