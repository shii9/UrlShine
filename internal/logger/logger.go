// Package logger provides structured, leveled, color-coded terminal output for URLShine.
// It handles all console logging with consistent formatting, timestamps, and color schemes.
package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Level defines the log verbosity level for filtering output.
type Level int

const (
	LevelDebug  Level = iota // Lowest: All messages including debug
	LevelInfo                // Standard: Info level and above
	LevelWarn                // Warnings and errors only
	LevelError               // Errors only
	LevelSilent              // Highest: Suppress all output
)

var (
	current = LevelInfo
	mu      sync.Mutex
	verbose = false
)

// SetLevel configures the minimum log level threshold.
func SetLevel(l Level) { mu.Lock(); current = l; mu.Unlock() }

// SetVerbose enables debug-level output for detailed execution tracking.
func SetVerbose(v bool) { mu.Lock(); verbose = v; mu.Unlock() }

// Color-coded log level prefixes for consistent visual hierarchy.
var (
	pInfo  = color.New(color.FgCyan, color.Bold).Sprint("INF")
	pOK    = color.New(color.FgGreen, color.Bold).Sprint(" ✔ ")
	pWarn  = color.New(color.FgYellow, color.Bold).Sprint(" ! ")
	pErr   = color.New(color.FgRed, color.Bold).Sprint(" ✘ ")
	pRun   = color.New(color.FgMagenta, color.Bold).Sprint("RUN")
	pSkip  = color.New(color.Faint).Sprint("---")
	pDebug = color.New(color.FgWhite, color.Faint).Sprint("DBG")

	// Spinner frames for visual activity indicator
	spinnerFrames = []string{"⠋", "⠙", "⠹", "⠴", "⠧", "⠽", "⠏"}

	// Convenience color functions for message formatting.
	cFaint  = color.New(color.Faint).SprintFunc()
	cWhite  = color.New(color.FgWhite, color.Bold).SprintFunc()
	cCyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	cGreen  = color.New(color.FgGreen, color.Bold).SprintFunc()
	cYellow = color.New(color.FgYellow, color.Bold).SprintFunc()
)

// RunWithSpinner logs that a tool is running and returns a function to stop the indicator.
// Since multiple tools run in parallel, we can't easily do a multi-line animated spinner
// without a full TUI library. Instead, we'll use a a "running" sign that looks active.
func RunWithSpinner(tool, target string) {
	mu.Lock()
	defer mu.Unlock()
	// We use a magenta "running" indicator.
	// The "spinner" part is simulated by the user's expectation of progress
	// unless we implement a full TUI. For now, we'll make the RUN prefix more dynamic.
	fmt.Printf("  %s  %s  %-20s  %s %s\n",
		ts(), pRun, cWhite(tool), cFaint(target), cCyan("↻ processing..."))
}

// ts returns formatted current timestamp for log messages.
func ts() string {
	return color.New(color.Faint).Sprintf("%s", time.Now().Format("15:04:05"))
}

// print is the core logging function that outputs formatted messages with timestamp and level prefix.
func print(prefix, format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("  %s  %s  %s\n", ts(), prefix, msg)
}

// Debug logs verbose debug-level information (only shown with -v flag).
func Debug(format string, args ...interface{}) {
	if verbose {
		print(pDebug, format, args...)
	}
}

// Info logs standard informational messages.
func Info(format string, args ...interface{}) { print(pInfo, format, args...) }

// Success logs successful operations with checkmark indicator.
func Success(format string, args ...interface{}) { print(pOK, format, args...) }

// Warn logs warning messages in yellow.
func Warn(format string, args ...interface{}) { print(pWarn, format, args...) }

// Error logs error messages in red.
func Error(format string, args ...interface{}) { print(pErr, format, args...) }

// Run logs active operations (tools running, collection starting, etc.).
func Run(format string, args ...interface{}) { print(pRun, format, args...) }

// Skip logs skipped operations (disabled tools, optional steps, etc.).
func Skip(format string, args ...interface{}) { print(pSkip, cFaint(fmt.Sprintf(format, args...))) }

// Step prints a prominent pipeline step header with ASCII separator for clear visual breaks.
// This is used to separate major phases of the reconnaissance pipeline.
func Step(n, total int, title string) {
	mu.Lock()
	defer mu.Unlock()
	bar := cCyan(strings.Repeat("─", 72))
	fmt.Println()
	fmt.Printf("  %s  %s\n",
		cCyan(fmt.Sprintf("STEP %d/%d", n, total)),
		cWhite(title),
	)
	fmt.Printf("  %s\n", bar)
}

// ToolMatrix displays a formatted table of external tool availability and status.
// Used during verification phase to show which tools are ready for collection.
func ToolMatrix(tools []struct {
	Name  string
	Found bool
}) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println()
	fmt.Printf("  %s  %s\n", cCyan("▸"), cWhite("TOOLS"))
	for _, t := range tools {
		if t.Found {
			fmt.Printf("  %s  %-25s %s\n", cGreen("✔"), t.Name, cFaint("ready"))
		} else {
			fmt.Printf("  %s  %-25s %s\n", cYellow("○"), t.Name, cFaint("not found"))
		}
	}
	fmt.Println()
}

// ToolResult logs the result of a single tool execution against a target.
// Displays URL count collected and handles both success and skip cases.
func ToolResult(tool, target string, count int, skipped bool) {
	mu.Lock()
	defer mu.Unlock()
	if skipped {
		fmt.Printf("  %s  %s  %-20s  %s\n", ts(), pSkip, tool, cFaint("not installed"))
		return
	}
	var result string
	if count == 0 {
		result = cYellow(fmt.Sprintf("→ 0 URLs"))
	} else if count > 5000 {
		result = cGreen(fmt.Sprintf("→ %s URLs (high yield)", FormatN(count)))
	} else {
		result = cGreen(fmt.Sprintf("→ %s URLs", FormatN(count)))
	}
	fmt.Printf("  %s  %s  %-20s  %-30s  %s\n",
		ts(), pOK, cWhite(tool), cFaint(target), result,
	)
}

// BlankLine outputs a blank line for visual spacing in logs.
func BlankLine() {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println()
}

// SectionHeader prints a professional section divider with title.
// Creates clear visual separation between major operations.
func SectionHeader(title string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("  %s  %s\n", cCyan("▸"), cWhite(title))
}

// CollectionProgress displays collection progress with visual progress bar.
func CollectionProgress(completed, total int, target string) {
	mu.Lock()
	defer mu.Unlock()
	percent := (completed * 100) / total
	bar := strings.Repeat("█", percent/5) + strings.Repeat("░", 20-percent/5)
	fmt.Printf("  %s  [%s] %d/%d  %s\n",
		ts(), bar, completed, total, cFaint(target))
}

// FormatN formats large numbers with K/M/B suffixes for readable output.
// Examples: 1500 → "1.5K", 2000000 → "2.0M"
func FormatN(n int) string {
	switch {
	case n >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(n)/1_000_000_000)
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}
