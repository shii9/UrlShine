// Package cmd implements URLShine command-line interface using Cobra framework.
// It defines all CLI flags, command execution logic, and help documentation.
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shii9/UrlShine/internal/runner"
	"github.com/shii9/UrlShine/internal/utils"

	"github.com/spf13/cobra"
)

const version = "2.0.0"

// CLI flag variables for collection tools.
var (
	flagFile      string // Input file with targets (one per line)
	flagOutputDir string // Output directory for results
	flagThreads   int    // Number of parallel threads
	flagDepth     int    // Crawl depth for active tools
	flagSubs      bool   // Include subdomains
	flagNoAlive   bool   // Skip live host verification
	flagSkipCol   bool   // Skip collection, reprocess existing data
	flagVerbose   bool   // Enable verbose/debug logging
	flagComplete  bool   // Run full pipeline (collection + processing)

	// Individual tool selection flags
	flagAll          bool // Run all collection tools
	flagGau          bool // Run GetAllUrls
	flagGospider     bool // Run GoSpider
	flagKatana       bool // Run Katana
	flagWaymore      bool // Run Waymore
	flagWaybackurls  bool // Run Wayback URLs
	flagHakrawler    bool // Run Hakrawler
	flagXnlinkfinder bool // Run xnLinkFinder
	flagGobuster     bool // Run Gobuster
	flagDirb         bool // Run Dirb
)

// rootCmd is the primary command entry point for URLShine.
// It orchestrates the full reconnaissance pipeline based on CLI flags and arguments.
var rootCmd = &cobra.Command{
	Use:   "urlshine [domain ...] [-f domains.txt]",
	Short: "URLShine v2.0.0 - Professional URL Enumeration & Attack Surface Mapper",
	Long: `URLShine orchestrates a sophisticated reconnaissance pipeline combining multiple
URL enumeration tools into a unified workflow. It collects URLs from passive and active
sources, deduplicates results, categorizes findings by attack vector, and verifies live hosts.

Key features:
  • Concurrent execution of 9 URL collection tools
  • Intelligent URL deduplication and normalization  
  • Automatic categorization into 5 attack vectors
  • Live host verification with status codes
  • Professional HTML, JSON, and Markdown reports

Use --all (or -a) to run every collector.
Use --complete (or -c) to run the full processing pipeline.`,
	Example: `  # Collect URLs from all tools
  urlshine -a google.com

  # Full processing pipeline (collection + categorization + verification)
  urlshine -a -c google.com

  # Specific tools only
  urlshine -gau -katana -waymore google.com

  # High-performance scan with aggressive settings
  urlshine -a -c -t 150 -d 5 google.com

  # Multiple targets from file
  urlshine -f targets.txt -a -c -o ./results

  # Fast mode (collection only, no live verification)
  urlshine -a -c --no-alive google.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		targets, err := resolveTargets(args, flagFile)
		if err != nil {
			return err
		}
		if len(targets) == 0 {
			return fmt.Errorf("no targets provided")
		}

		outDir := flagOutputDir
		if outDir == "" {
			outDir = "urlshine_" + time.Now().Format("20060102_150405")
		}

		return runner.RunProfessional(runner.Options{
			Targets:     targets,
			OutputDir:   outDir,
			Threads:     flagThreads,
			Depth:       flagDepth,
			Subs:        flagSubs,
			SkipAlive:   flagNoAlive,
			SkipCollect: flagSkipCol,
			Verbose:     flagVerbose,
			RunComplete: flagComplete,

			RunAll:          flagAll,
			RunGau:          flagGau,
			RunGospider:     flagGospider,
			RunKatana:       flagKatana,
			RunWaymore:      flagWaymore,
			RunWaybackurls:  flagWaybackurls,
			RunHakrawler:    flagHakrawler,
			RunXnlinkfinder: flagXnlinkfinder,
			RunGobuster:     flagGobuster,
			RunDirb:         flagDirb,
		})
	},
}

func Execute() {
	// Normalize flags: convert -all to --all, -complete to --complete, etc.
	// This allows users to use -all and --all interchangeably
	normalizeFlags()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// normalizeFlags converts single-dash long flags (-all, -complete, -gau, etc.) to double-dash (--all, --complete, --gau)
// This provides a better UX by accepting both formats
func normalizeFlags() {
	// Map of known long flag names that users might try with single dash
	longFlagNames := map[string]bool{
		"all":          true,
		"complete":     true,
		"gau":          true,
		"katana":       true,
		"gospider":     true,
		"waymore":      true,
		"waybackurls":  true,
		"hakrawler":    true,
		"xnlinkfinder": true,
		"gobuster":     true,
		"dirb":         true,
		"file":         true,
		"output":       true,
		"threads":      true,
		"depth":        true,
		"subs":         true,
		"verbose":      true,
		"no-alive":     true,
		"skip-collect": true,
	}

	for i, arg := range os.Args {
		// Check if it's a single-dash flag (starts with - but not --)
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			// Extract flag name (everything after the dash until = or end)
			flagPart := arg[1:] // Remove leading -
			var flagName string
			if strings.Contains(flagPart, "=") {
				flagName = strings.Split(flagPart, "=")[0]
			} else {
				flagName = flagPart
			}

			// Check if this is a known long flag name
			if longFlagNames[flagName] {
				// Convert -flagname to --flagname
				os.Args[i] = "--" + flagPart
			}
		}
	}
}

func init() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Printf(`
╔════════════════════════════════════════════════════════════════════════════╗
║                         URLShine v%s                                      ║
║       Professional URL Enumeration & Attack Surface Mapper                ║
╚════════════════════════════════════════════════════════════════════════════╝

USAGE
  urlshine [target ...] [flags]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

EXAMPLES

  Collect URLs from all tools:
    $ urlshine -a google.com

  Full processing (collect + categorize + verify):
    $ urlshine -a -c google.com

  Specific tools only:
    $ urlshine -gau -katana google.com

  Multiple targets with aggressive settings:
    $ urlshine -f targets.txt -a -c -t 150 -d 5

  Fast mode (collection only):
    $ urlshine -a -c --no-alive google.com

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

MAIN FLAGS

  -a, --all                  Run all collection tools
  -c, --complete             Run full pipeline (collection + processing)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

COLLECTION TOOLS (select one or more)

  -g, --gau                  GetAllUrls (passive archives)
  -k, --katana               Katana (active crawler)
  -w, --gospider             GoSpider (HTML & JS)
  -m, --waymore              Waymore (Wayback Machine)
  -b, --waybackurls          Wayback URLs
  -r, --hakrawler            Hakrawler (HTML crawler)
  -x, --xnlinkfinder         xnLinkFinder (JS extraction)
  -u, --gobuster             Gobuster (directory discovery)
  -i, --dirb                 Dirb (directory brute-force)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

OPTIONS

  -t, --threads INT          Parallel threads (default: 50)
  -d, --depth INT            Crawl depth (default: 5)
  -f, --file FILE            Input file with targets
  -o, --output DIR           Output directory
  -v, --verbose              Enable debug logging
  --no-alive                 Skip live verification
  --skip-collect             Reprocess existing data

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

For more information, visit: https://github.com/shii9/UrlShine
`, version)
	})

	f := rootCmd.Flags()

	// ─── Input/Output ────────────────────────────────────────────────────────────
	f.StringVarP(&flagFile, "file", "f", "", "Input file with targets (one per line)")
	f.StringVarP(&flagOutputDir, "output", "o", "", "Output directory (default: urlshine_<timestamp>)")

	// ─── Execution Parameters ─────────────────────────────────────────────────────
	f.IntVarP(&flagThreads, "threads", "t", 50, "Parallel threads for concurrent execution (recommended: 50-150)")
	f.IntVarP(&flagDepth, "depth", "d", 5, "Crawl depth for active tools like Katana (higher = more thorough but slower)")
	f.BoolVarP(&flagSubs, "subs", "s", true, "Include subdomains in enumeration")

	// ─── Pipeline Control ─────────────────────────────────────────────────────────
	f.BoolVarP(&flagNoAlive, "no-alive", "n", false, "Skip live host verification (fast mode)")
	f.BoolVar(&flagSkipCol, "skip-collect", false, "Skip collection and process existing files")
	f.BoolVarP(&flagVerbose, "verbose", "v", false, "Enable debug/verbose logging")
	f.BoolVarP(&flagComplete, "complete", "c", false, "Run merge, normalize, categorize, reports, and alive checking")

	f.BoolVarP(&flagAll, "all", "a", false, "Run all supported collection tools")
	f.BoolVarP(&flagGau, "gau", "g", false, "Run GAU (GetAllUrls)")
	f.BoolVarP(&flagGospider, "gospider", "w", false, "Run GoSpider")
	f.BoolVarP(&flagKatana, "katana", "k", false, "Run Katana")
	f.BoolVarP(&flagWaymore, "waymore", "m", false, "Run Waymore")
	f.BoolVarP(&flagWaybackurls, "waybackurls", "b", false, "Run Waybackurls")
	f.BoolVarP(&flagHakrawler, "hakrawler", "r", false, "Run Hakrawler")
	f.BoolVarP(&flagXnlinkfinder, "xnlinkfinder", "x", false, "Run xnLinkFinder")
	f.BoolVarP(&flagGobuster, "gobuster", "u", false, "Run Gobuster directory discovery")
	f.BoolVarP(&flagDirb, "dirb", "i", false, "Run Dirb directory enumeration")
}

func resolveTargets(args []string, file string) ([]string, error) {
	seen := make(map[string]struct{})
	var out []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		s = strings.TrimPrefix(s, "https://")
		s = strings.TrimPrefix(s, "http://")
		s = strings.TrimRight(s, "/")
		if s == "" {
			return
		}
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}

	for _, a := range args {
		add(a)
	}
	if file != "" {
		lines, err := utils.ReadLines(file)
		if err != nil {
			return nil, fmt.Errorf("read file %q: %w", file, err)
		}
		for _, l := range lines {
			add(l)
		}
	}
	return out, nil
}
