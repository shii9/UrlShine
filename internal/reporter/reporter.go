// Package reporter handles final terminal and markdown/JSON reporting.
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

var (
	cCyan   = color.New(color.FgCyan, color.Bold)
	cGreen  = color.New(color.FgGreen, color.Bold)
	cYellow = color.New(color.FgYellow, color.Bold)
	cWhite  = color.New(color.FgWhite, color.Bold)
	cFaint  = color.New(color.Faint)
)

// Stats holds all pipeline metrics.
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

var groupOrder = []struct{ key, label string }{
	{splitter.GroupAPI, "API Endpoints"},
	{splitter.GroupAuth, "Auth / Admin Pages"},
	{splitter.GroupParams, "URLs with Parameters"},
	{splitter.GroupJS, "JS & Config Files"},
	{splitter.GroupDirs, "Directory Paths"},
}

// Print renders the box-drawing summary.
func Print(s Stats, dur time.Duration) {
	w := 68
	line := func(msg string) { fmt.Printf("  %s %-*s %s\n", cCyan.Sprint("║"), w, msg, cCyan.Sprint("║")) }
	sep := func() { fmt.Printf("  %s\n", cCyan.Sprint("╠"+strings.Repeat("═", w+2)+"╣")) }
	top := func() { fmt.Printf("  %s\n", cCyan.Sprint("╔"+strings.Repeat("═", w+2)+"╗")) }
	bot := func() { fmt.Printf("  %s\n", cCyan.Sprint("╚"+strings.Repeat("═", w+2)+"╝")) }

	fmt.Println()
	top()
	line(cWhite.Sprintf("  %-*s", w-2, "URLShine — SCAN SUMMARY"))
	sep()

	tStr := strings.Join(s.Targets, ", ")
	if len(tStr) > 40 {
		tStr = tStr[:37] + "..."
	}
	line(fmt.Sprintf("  %-20s %s", cFaint.Sprint("Target(s)  :"), cGreen.Sprint(tStr)))
	line(fmt.Sprintf("  %-20s %s", cFaint.Sprint("Duration   :"), cWhite.Sprint(dur.Round(time.Second))))
	line(fmt.Sprintf("  %-20s %s", cFaint.Sprint("Output Dir :"), cFaint.Sprint(s.OutputDir)))

	sep()
	line(cWhite.Sprint("  COLLECTION"))
	line(fmt.Sprintf("  %-32s %s", "Total Raw URLs collected", cGreen.Sprintf("%-10s", utils.FormatN(s.TotalRaw))))
	line(fmt.Sprintf("  %-32s %s  %s", "After Normalization",
		cGreen.Sprintf("%-10s", utils.FormatN(s.AfterNorm)),
		cFaint.Sprintf("(%.0f%% reduction)", utils.Reduction(s.TotalRaw, s.AfterNorm))))

	sep()
	line(cWhite.Sprint("  GROUPS"))
	for _, g := range groupOrder {
		cnt := s.Groups[g.key]
		alive := s.AliveGroups[g.key]
		aliveStr := ""
		if alive > 0 {
			aliveStr = cGreen.Sprintf("  → %s alive", utils.FormatN(alive))
		}
		line(fmt.Sprintf("  %-30s %s%s", g.label, cYellow.Sprintf("%-10s", utils.FormatN(cnt)), aliveStr))
	}
	line("")
	line(fmt.Sprintf("  %-30s %s", "Unique Param Keys", cYellow.Sprintf("%d", s.UniqueParams)))
	bot()
	fmt.Println()
}

// WriteReports writes markdown and JSON reports.
func WriteReports(s Stats) error {
	// JSON
	jPath := filepath.Join(s.OutputDir, "urlshine_report.json")
	if err := utils.WriteJSON(jPath, s); err != nil {
		return err
	}

	// Markdown
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "# URLShine — Scan Report\n\n")
	fmt.Fprintf(sb, "**Generated:** %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(sb, "**Targets:** %s\n\n", strings.Join(s.Targets, ", "))
	fmt.Fprintf(sb, "### Collection\n\n")
	fmt.Fprintf(sb, "- **Raw URLs:** %d\n", s.TotalRaw)
	fmt.Fprintf(sb, "- **Normalized:** %d (%.0f%% reduction)\n\n", s.AfterNorm, utils.Reduction(s.TotalRaw, s.AfterNorm))
	fmt.Fprintf(sb, "### Groups\n\n| Group | URLs | Alive |\n|---|---|---|\n")
	for _, g := range groupOrder {
		fmt.Fprintf(sb, "| %s | %d | %d |\n", g.label, s.Groups[g.key], s.AliveGroups[g.key])
	}
	fmt.Fprintf(sb, "\n**Unique Parameters:** %d\n", s.UniqueParams)
	if verified := s.AliveGroups["verified"]; verified > 0 {
		fmt.Fprintf(sb, "\n**Verified Live URLs:** %d\n", verified)
	}

	return os.WriteFile(filepath.Join(s.OutputDir, "urlshine_report.md"), []byte(sb.String()), 0644)
}
