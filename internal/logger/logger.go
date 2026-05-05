// Package logger provides structured, leveled, color-coded terminal output.
package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Level controls log verbosity.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelSilent
)

var (
	current = LevelInfo
	mu      sync.Mutex
	verbose = false
)

// SetLevel configures the minimum log level.
func SetLevel(l Level) { mu.Lock(); current = l; mu.Unlock() }

// SetVerbose enables debug-level output.
func SetVerbose(v bool) { mu.Lock(); verbose = v; mu.Unlock() }

var (
	pInfo  = color.New(color.FgCyan, color.Bold).Sprint("INF")
	pOK    = color.New(color.FgGreen, color.Bold).Sprint(" ✔ ")
	pWarn  = color.New(color.FgYellow, color.Bold).Sprint(" ! ")
	pErr   = color.New(color.FgRed, color.Bold).Sprint(" ✘ ")
	pRun   = color.New(color.FgMagenta, color.Bold).Sprint("RUN")
	pSkip  = color.New(color.Faint).Sprint("---")
	pDebug = color.New(color.FgWhite, color.Faint).Sprint("DBG")

	cFaint  = color.New(color.Faint).SprintFunc()
	cWhite  = color.New(color.FgWhite, color.Bold).SprintFunc()
	cCyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	cGreen  = color.New(color.FgGreen, color.Bold).SprintFunc()
	cYellow = color.New(color.FgYellow, color.Bold).SprintFunc()
)

func ts() string {
	return color.New(color.Faint).Sprintf("%s", time.Now().Format("15:04:05"))
}

func print(prefix, format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("  %s  %s  %s\n", ts(), prefix, msg)
}

func Debug(format string, args ...interface{}) {
	if verbose {
		print(pDebug, format, args...)
	}
}
func Info(format string, args ...interface{})    { print(pInfo, format, args...) }
func Success(format string, args ...interface{}) { print(pOK, format, args...) }
func Warn(format string, args ...interface{})    { print(pWarn, format, args...) }
func Error(format string, args ...interface{})   { print(pErr, format, args...) }
func Run(format string, args ...interface{})     { print(pRun, format, args...) }
func Skip(format string, args ...interface{})    { print(pSkip, cFaint(fmt.Sprintf(format, args...))) }

// Step prints a prominent pipeline step header with a separator bar.
func Step(n, total int, title string) {
	mu.Lock()
	defer mu.Unlock()
	bar := cCyan(strings.Repeat("━", 64))
	fmt.Println()
	fmt.Printf("  %s\n", bar)
	fmt.Printf("  %s  %s  %s\n",
		cCyan(fmt.Sprintf("STEP %d/%d", n, total)),
		cFaint("›"),
		cWhite(strings.ToUpper(title)),
	)
	fmt.Printf("  %s\n", bar)
	fmt.Println()
}

// ToolMatrix prints the availability table for all external tools.
func ToolMatrix(tools []struct {
	Name  string
	Found bool
}) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println()
	fmt.Printf("  %s\n", cWhite("External Tools"))
	fmt.Printf("  %s\n", cFaint(strings.Repeat("─", 42)))
	for _, t := range tools {
		if t.Found {
			fmt.Printf("  %s  %-20s %s\n", cGreen("✔"), t.Name, cFaint("ready"))
		} else {
			fmt.Printf("  %s  %-20s %s\n", cYellow("✘"), t.Name, cFaint("not found — will skip"))
		}
	}
	fmt.Println()
}

// ToolResult prints the result of one tool run against one target.
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
		result = cGreen(fmt.Sprintf("→ %s URLs", FormatN(count)))
	} else {
		result = cGreen(fmt.Sprintf("→ %s URLs", FormatN(count)))
	}
	fmt.Printf("  %s  %s  %-20s  %-30s  %s\n",
		ts(), pOK, cWhite(tool), cFaint(target), result,
	)
}

// Banner is for standalone use outside of banner package.
func BlankLine() {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println()
}

// SectionHeader prints a professional section divider
func SectionHeader(title string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println()
	fmt.Printf("  %s  %s\n", cCyan("▸"), cWhite(strings.ToUpper(title)))
	fmt.Printf("  %s\n", cFaint(strings.Repeat("─", 80)))
	fmt.Println()
}

// CollectionProgress shows overall collection progress
func CollectionProgress(completed, total int, target string) {
	mu.Lock()
	defer mu.Unlock()
	percent := (completed * 100) / total
	bar := strings.Repeat("█", percent/5) + strings.Repeat("░", 20-percent/5)
	fmt.Printf("  %s  [%s] %d/%d  %s\n",
		ts(), bar, completed, total, cFaint(target))
}

// FormatN formats large numbers with K/M suffixes.
func FormatN(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}
