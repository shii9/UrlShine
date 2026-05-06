// Package reporter handles professional terminal and file-based reporting for URLShine scan results.
// Generates color-coded terminal output, machine-readable JSON reports, and human-friendly Markdown summaries.
package reporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shii9/UrlShine/internal/splitter"
	"github.com/shii9/UrlShine/internal/utils"

	"github.com/fatih/color"
)

// Color definitions for consistent styling across reports.
var (
	cCyan   = color.New(color.FgCyan, color.Bold)
	cGreen  = color.New(color.FgGreen, color.Bold)
	cYellow = color.New(color.FgYellow, color.Bold)
	cWhite  = color.New(color.FgWhite, color.Bold)
	cFaint  = color.New(color.Faint)
)

// Stats holds comprehensive pipeline metrics for final reporting.
// Includes raw and processed URL counts, categorization results, and verification data.
type Stats struct {
	Targets      []string       `json:"targets"`
	TotalRaw     int            `json:"total_raw"`
	AfterNorm    int            `json:"after_norm"`
	Groups       map[string]int `json:"groups"`
	AliveGroups  map[string]int `json:"alive_groups"`
	UniqueParams int            `json:"unique_params"`
	OutputDir    string         `json:"output_dir"`
	DurationSec  float64        `json:"duration_seconds"`
}

// groupOrder defines display order and labels for URL categories.
var groupOrder = []struct{ key, label string }{
	{splitter.GroupAPI, "API Endpoints"},
	{splitter.GroupAuth, "Auth / Admin Pages"},
	{splitter.GroupParams, "URLs with Parameters"},
	{splitter.GroupJS, "JS & Config Files"},
	{splitter.GroupDirs, "Directory Paths"},
}

// Print renders a professional ASCII-boxed terminal summary of scan results.
// Uses color-coding to highlight important metrics and uses box-drawing characters for clarity.
func Print(s Stats, dur time.Duration) {
	w := 76
	line := func(msg string) { fmt.Printf("  %s %-*s %s\n", cCyan.Sprint("║"), w, msg, cCyan.Sprint("║")) }
	sep := func() { fmt.Printf("  %s\n", cCyan.Sprint("╠"+strings.Repeat("═", w+2)+"╣")) }
	top := func() { fmt.Printf("  %s\n", cCyan.Sprint("╔"+strings.Repeat("═", w+2)+"╗")) }
	bot := func() { fmt.Printf("  %s\n", cCyan.Sprint("╚"+strings.Repeat("═", w+2)+"╝")) }

	fmt.Println()
	top()
	line(cWhite.Sprintf("  %-*s", w-2, "URLShine — SCAN SUMMARY"))
	sep()

	// Target(s) information
	tStr := strings.Join(s.Targets, ", ")
	if len(tStr) > 50 {
		tStr = tStr[:47] + "..."
	}
	line(fmt.Sprintf("  %-20s %s", cFaint.Sprint("Target(s)  :"), cGreen.Sprint(tStr)))
	line(fmt.Sprintf("  %-20s %s", cFaint.Sprint("Duration   :"), cWhite.Sprint(dur.Round(time.Second))))
	line(fmt.Sprintf("  %-20s %s", cFaint.Sprint("Output Dir :"), cFaint.Sprint(s.OutputDir)))

	// Collection metrics
	sep()
	line(cWhite.Sprint("  COLLECTION METRICS"))
	line(fmt.Sprintf("  %-32s %s", "Raw URLs collected", cGreen.Sprintf("%-12s", utils.FormatN(s.TotalRaw))))
	line(fmt.Sprintf("  %-32s %s  %s", "After Normalization",
		cGreen.Sprintf("%-12s", utils.FormatN(s.AfterNorm)),
		cFaint.Sprintf("(%.1f%% deduped)", utils.Reduction(s.TotalRaw, s.AfterNorm))))

	// Category breakdown
	sep()
	line(cWhite.Sprint("  CATEGORIZED URLS"))
	for _, g := range groupOrder {
		cnt := s.Groups[g.key]
		alive := s.AliveGroups[g.key]
		aliveStr := ""
		if alive > 0 {
			aliveStr = cGreen.Sprintf("  [%s alive]", utils.FormatN(alive))
		} else if cnt > 0 {
			aliveStr = cYellow.Sprint("  [verifying...]")
		}
		line(fmt.Sprintf("  %-32s %s%s", g.label, cYellow.Sprintf("%-12s", utils.FormatN(cnt)), aliveStr))
	}

	// Additional metrics
	line("")
	line(fmt.Sprintf("  %-32s %s", "Unique Parameter Keys", cYellow.Sprintf("%d", s.UniqueParams)))

	// Calculate total live
	totalAlive := 0
	for _, g := range groupOrder {
		totalAlive += s.AliveGroups[g.key]
	}
	if totalAlive > 0 {
		line(fmt.Sprintf("  %-32s %s", "Total Verified Live URLs", cGreen.Sprintf("%d", totalAlive)))
	}

	bot()
	fmt.Println()
}

// WriteReports writes comprehensive machine-readable (JSON) and human-readable (Markdown) reports to disk.
// Both formats include detailed metrics, categorization results, and scan parameters.
func WriteReports(s Stats) error {
	// JSON Report
	jPath := filepath.Join(s.OutputDir, "urlshine_report.json")
	if err := utils.WriteJSON(jPath, s); err != nil {
		return err
	}

	// Markdown Report
	sb := &strings.Builder{}

	// Header
	fmt.Fprintf(sb, "# URLShine — Scan Report\n\n")
	fmt.Fprintf(sb, "**Generated:** %s  \n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(sb, "**Targets:** %s  \n", strings.Join(s.Targets, ", "))
	fmt.Fprintf(sb, "**Duration:** %.1f seconds  \n", s.DurationSec)
	fmt.Fprintf(sb, "**Output Directory:** `%s`\n\n", s.OutputDir)

	// Collection Summary
	fmt.Fprintf(sb, "## Collection Summary\n\n")
	fmt.Fprintf(sb, "- **Raw URLs Collected:** %s\n", utils.FormatN(s.TotalRaw))
	fmt.Fprintf(sb, "- **After Normalization:** %s (%.1f%% reduction)\n",
		utils.FormatN(s.AfterNorm),
		utils.Reduction(s.TotalRaw, s.AfterNorm))
	fmt.Fprintf(sb, "- **Unique Parameters Found:** %d\n\n", s.UniqueParams)

	// Categories Breakdown
	fmt.Fprintf(sb, "## URL Categories\n\n")
	fmt.Fprintf(sb, "| Category | Total URLs | Verified Live |\n")
	fmt.Fprintf(sb, "|----------|------------|---------------|\n")

	totalAlive := 0
	for _, g := range groupOrder {
		alive := s.AliveGroups[g.key]
		totalAlive += alive
		fmt.Fprintf(sb, "| %s | %d | %d |\n",
			g.label,
			s.Groups[g.key],
			alive)
	}
	fmt.Fprintf(sb, "\n**Total Verified Live:** %d\n\n", totalAlive)

	// Output Files
	fmt.Fprintf(sb, "## Output Files\n\n")
	fmt.Fprintf(sb, "### Raw Collection Results\n")
	fmt.Fprintf(sb, "- `raw/` — Individual tool output files\n")
	fmt.Fprintf(sb, "- `merged_raw.txt` — All URLs combined\n\n")

	fmt.Fprintf(sb, "### Processed Results\n")
	fmt.Fprintf(sb, "- `normalized_urls.txt` — Deduplicated and cleaned URLs\n")
	fmt.Fprintf(sb, "- `api_endpoints.txt` — API endpoint URLs\n")
	fmt.Fprintf(sb, "- `auth_admin_urls.txt` — Authentication and admin pages\n")
	fmt.Fprintf(sb, "- `parameters_urls.txt` — URLs with query parameters\n")
	fmt.Fprintf(sb, "- `js_config_urls.txt` — JavaScript and config files\n")
	fmt.Fprintf(sb, "- `directories_urls.txt` — Directory paths\n\n")

	fmt.Fprintf(sb, "### Verification Results\n")
	fmt.Fprintf(sb, "- `alive_api_endpoints.txt` — Verified live API endpoints\n")
	fmt.Fprintf(sb, "- `alive_auth_admin_urls.txt` — Verified live auth pages\n")
	fmt.Fprintf(sb, "- `alive_parameters_urls.txt` — Verified live parameterized URLs\n")
	fmt.Fprintf(sb, "- `alive_js_config_urls.txt` — Verified live JS and config\n")
	fmt.Fprintf(sb, "- `alive_directories_urls.txt` — Verified live directories\n\n")

	fmt.Fprintf(sb, "## Statistics\n\n")
	fmt.Fprintf(sb, "- **Deduplication Rate:** %.1f%%\n", utils.Reduction(s.TotalRaw, s.AfterNorm))
	fmt.Fprintf(sb, "- **Live Verification Rate:** %.1f%%\n",
		float64(totalAlive)*100.0/float64(s.AfterNorm))
	fmt.Fprintf(sb, "- **Scan Efficiency:** %d URLs/second\n\n",
		int(float64(s.TotalRaw)/s.DurationSec))

	fmt.Fprintf(sb, "---\n\n")
	fmt.Fprintf(sb, "*Generated by URLShine v2.0.0*\n")

	return os.WriteFile(filepath.Join(s.OutputDir, "urlshine_report.md"), []byte(sb.String()), 0644)
}
