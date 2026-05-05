// Package utils provides shared file I/O, deduplication, tool detection, and helpers.
package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const scanBuf = 4 * 1024 * 1024 // 4 MB scanner buffer

// ─── Tool Detection ───────────────────────────────────────────────────────────

// ToolExists returns true if the named binary is found in PATH.
func ToolExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// ToolStatus represents the status of a single tool.
type ToolStatus struct {
	Name       string
	Status     string // "installed", "missing", "optional"
	InstallCmd string
	Required   bool
}

// CheckDependencies verifies all required tools and returns their status.
func CheckDependencies() []ToolStatus {
	tools := []ToolStatus{
		// Go-based tools
		{Name: "gau", InstallCmd: "go install github.com/lc/gau/v2/cmd/gau@latest", Required: false},
		{Name: "katana", InstallCmd: "go install github.com/projectdiscovery/katana/cmd/katana@latest", Required: false},
		{Name: "gospider", InstallCmd: "go install github.com/jaeles-project/gospider@latest", Required: false},
		{Name: "waymore", InstallCmd: "pip3 install waymore", Required: false},
		{Name: "waybackurls", InstallCmd: "go install github.com/tomnomnom/waybackurls@latest", Required: false},
		{Name: "hakrawler", InstallCmd: "go install github.com/hakluke/hakrawler@latest", Required: false},
		{Name: "xnlinkfinder", InstallCmd: "pip3 install xnlinkfinder", Required: false},
		{Name: "gobuster", InstallCmd: "go install github.com/OJ/gobuster/v3@latest", Required: false},
		{Name: "dirb", InstallCmd: "apt-get install dirb (Linux) | brew install dirb (macOS)", Required: false},
		// Optional but useful
		{Name: "httpx", InstallCmd: "go install github.com/projectdiscovery/httpx/cmd/httpx@latest", Required: false},
	}

	for i := range tools {
		if ToolExists(tools[i].Name) {
			tools[i].Status = "installed"
		} else {
			tools[i].Status = "missing"
		}
	}
	return tools
}

// ─── File Helpers ─────────────────────────────────────────────────────────────

// ReadLines reads all non-empty, trimmed lines from path.
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, scanBuf), scanBuf)
	for sc.Scan() {
		if l := strings.TrimSpace(sc.Text()); l != "" {
			out = append(out, l)
		}
	}
	return out, sc.Err()
}

// WriteLines writes lines to path one per line (creates/truncates).
func WriteLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriterSize(f, scanBuf)
	for _, l := range lines {
		fmt.Fprintln(w, l)
	}
	return w.Flush()
}

// AppendLines appends lines to path (creates if missing).
func AppendLines(path string, lines []string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, l := range lines {
		fmt.Fprintln(w, l)
	}
	return w.Flush()
}

// WriteJSON marshals v as indented JSON to path.
func WriteJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// FileExists returns true if path exists and is a regular file.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// FileLineCount returns count of non-empty lines in path (0 on error).
func FileLineCount(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	n := 0
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, scanBuf), scanBuf)
	for sc.Scan() {
		if strings.TrimSpace(sc.Text()) != "" {
			n++
		}
	}
	return n
}

// EnsureDir creates dir and all parents if they don't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// ─── Deduplication ────────────────────────────────────────────────────────────

// Deduplicate removes duplicates preserving first-seen order.
func Deduplicate(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// DeduplicateSort removes duplicates and returns sorted result.
func DeduplicateSort(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

// ─── String Helpers ───────────────────────────────────────────────────────────

// SanitizeFilename converts a URL/domain to a safe filename component.
func SanitizeFilename(s string) string {
	r := strings.NewReplacer(
		"https://", "", "http://", "",
		"/", "_", ":", "_", ".", "_",
		"*", "_", "?", "_", "&", "_",
	)
	return strings.Trim(r.Replace(s), "_")
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

// Reduction calculates percentage reduction from a to b.
func Reduction(a, b int) float64 {
	if a == 0 {
		return 0
	}
	return float64(a-b) / float64(a) * 100
}
