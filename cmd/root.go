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

const version = "3.2.0"

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

	// Tool selection flags (10 core, non-redundant tools)
	flagAll             bool // Run all collection tools
	flagGau             bool // Run GetAllUrls
	flagWaymore         bool // Run Waymore
	flagParamspider     bool // Run ParamSpider
	flagCommoncrawl     bool // Run Common Crawl CDX (native Go)
	flagUrlfinder       bool // Run URLFinder
	flagGithubEndpoints bool // Run github-endpoints
	flagXnlinkfinder    bool // Run xnLinkFinder
	flagKatana          bool // Run Katana
	flagHakrawler       bool // Run Hakrawler
	flagGobuster        bool // Run Gobuster
)

// rootCmd is the primary command entry point for URLShine.
var rootCmd = &cobra.Command{
	Use:   "urlshine [domain/url/file ...] [-f targets.txt]",
	Short: "URLShine v3.2.0 - Professional URL Enumeration & Attack Surface Mapper",
	Long: `URLShine orchestrates a sophisticated reconnaissance pipeline combining 10 distinct,
high-value URL enumeration tools into a unified workflow. It collects URLs from passive and active
sources, deduplicates results, categorizes findings by attack vector, and verifies live hosts.

Accepts target domains (e.g. google.com), target URLs (e.g. https://site.com/app),
or text files containing lists of targets (e.g. targets.txt or -f targets.txt).

Running urlshine runs ALL 10 tools by default to collect maximum URLs.

For detailed usage and scenarios, see: https://github.com/shii9/UrlShine#-usage-guide--practical-examples`,
	Args: cobra.ArbitraryArgs, // Accept domains/files as positional arguments
	RunE: func(cmd *cobra.Command, args []string) error {
		targets, err := resolveTargets(args, flagFile)
		if err != nil {
			return err
		}
		if len(targets) == 0 {
			return fmt.Errorf("no targets provided. Pass domains, URLs, or a text file (e.g. urlshine example.com or urlshine -f targets.txt)")
		}

		outDir := flagOutputDir
		if outDir == "" {
			outDir = "urlshine_" + time.Now().Format("20060102_150405")
		}

		// Check tool selection: if no individual tool flags were explicitly passed, enable ALL tools by default!
		runAllTools := flagAll
		if !flagGau && !flagWaymore && !flagParamspider && !flagCommoncrawl &&
			!flagUrlfinder && !flagGithubEndpoints && !flagXnlinkfinder &&
			!flagKatana && !flagHakrawler && !flagGobuster {
			runAllTools = true
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

			// 10 Core Tools (RunAll defaults to true to maximize URL collection)
			RunAll:             runAllTools,
			RunGau:             flagGau,
			RunWaymore:         flagWaymore,
			RunParamspider:     flagParamspider,
			RunCommoncrawl:     flagCommoncrawl,
			RunUrlfinder:       flagUrlfinder,
			RunGithubEndpoints: flagGithubEndpoints,
			RunXnlinkfinder:    flagXnlinkfinder,
			RunKatana:          flagKatana,
			RunHakrawler:       flagHakrawler,
			RunGobuster:        flagGobuster,
		})
	},
}

func Execute() {
	normalizeFlags()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// normalizeFlags converts single-dash long flags to double-dash.
func normalizeFlags() {
	longFlagNames := map[string]bool{
		"all":         true,
		"complete":    true,

		"gau":              true,
		"waymore":          true,
		"paramspider":      true,
		"commoncrawl":      true,
		"urlfinder":        true,
		"github-endpoints": true,
		"xnlinkfinder":     true,
		"katana":           true,
		"hakrawler":        true,
		"gobuster":         true,

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
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			flagPart := arg[1:]
			var flagName string
			if strings.Contains(flagPart, "=") {
				flagName = strings.Split(flagPart, "=")[0]
			} else {
				flagName = flagPart
			}

			if longFlagNames[flagName] {
				os.Args[i] = "--" + flagPart
			}
		}
	}
}

func init() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Printf(`
╔══════════════════════════════════════════════════════════════════════════════════════════╗
║                             URLShine v%s                                             ║
║               Professional URL Enumeration & Attack Surface Mapper                       ║
╚══════════════════════════════════════════════════════════════════════════════════════════╝

USAGE
  urlshine [target/file ...] [flags]

QUICK START
  $ urlshine example.com             # Run all 10 tools against a single domain
  $ urlshine targets.txt             # Run all 10 tools against targets from a file
  $ urlshine -f targets.txt -c       # Full pipeline: collect + merge + normalize + verify
  $ urlshine doctor                  # Audit installed tools and dependencies

TARGET INPUT SUPPORT
  • Domains:           example.com, sub.example.com
  • Full URLs:         https://example.com/api/v1
  • Target Lists:      targets.txt (direct positional or via -f targets.txt)

PIPELINE FLAGS
  -a, --all                  Run all 10 collection tools (ENABLED BY DEFAULT)
  -c, --complete             Run full processing pipeline (merge, normalize, categorize, verify)

COLLECTION TOOLS (10 Core Tools in 4 Tiers)
  Tier 1 — Passive Archives (Zero Target Traffic):
    -g, --gau                GetAllUrls (Wayback, CommonCrawl, URLScan, OTX)
    -m, --waymore            Waymore (enhanced Wayback & HTTP response parser)
        --paramspider        ParamSpider (parameter-focused archive miner)

  Tier 2 — Passive APIs & OSINT (Zero Target Traffic):
        --commoncrawl        Common Crawl CDX API (built-in native Go client)
        --urlfinder          URLFinder (crt.sh, VirusTotal, passive intel)
        --github-endpoints   GitHub endpoint discovery (leaked repo code)
    -x, --xnlinkfinder       xnLinkFinder (JS/HTML link & endpoint extraction)

  Tier 3 — Active Crawlers (Generates Target Traffic):
    -k, --katana             Katana (JS-capable headless Chromium crawler)
        --hakrawler          Hakrawler (fast Go web crawler)

  Tier 4 — Active Brute-Force (Directory Enumeration):
    -z, --gobuster           Gobuster (high-speed directory/file brute-force)

INPUT / OUTPUT OPTIONS
  -f, --file FILE            Input file containing targets (one per line)
  -o, --output DIR           Output directory (default: urlshine_<timestamp>)

ADVANCED TUNING
  -t, --threads INT          Parallel thread pool (default: 50, range: 1-500)
  -d, --depth INT            Crawl depth for active tools (default: 5)
  -v, --verbose              Enable debug/verbose logging
  -s, --subs                 Include subdomains in scans (default: true)
      --no-alive             Skip live host verification stage (fast mode)
      --skip-collect         Skip collection phase and reprocess existing raw data

DOCUMENTATION & DIAGNOSTICS
  📖 README:  https://github.com/shii9/UrlShine#-usage-guide--practical-examples
  🛟 Doctor:  urlshine doctor

`, version)
	})

	// Define flags
	f := rootCmd.Flags()
	f.StringVarP(&flagFile, "file", "f", "", "Input file with targets (one per line)")
	f.StringVarP(&flagOutputDir, "output", "o", "", "Output directory (default: urlshine_<timestamp>)")
	f.IntVarP(&flagThreads, "threads", "t", 50, "Parallel threads for concurrent execution (recommended: 50-150)")
	f.IntVarP(&flagDepth, "depth", "d", 5, "Crawl depth for active tools like Katana (higher = more thorough but slower)")
	f.BoolVarP(&flagSubs, "subs", "s", true, "Include subdomains in enumeration")
	f.BoolVarP(&flagNoAlive, "no-alive", "n", false, "Skip live host verification (fast mode)")
	f.BoolVar(&flagSkipCol, "skip-collect", false, "Skip collection and process existing files")
	f.BoolVarP(&flagVerbose, "verbose", "v", false, "Enable debug/verbose logging")
	f.BoolVarP(&flagComplete, "complete", "c", false, "Run merge, normalize, categorize, reports, and alive checking")

	// 10 Core tool flags
	f.BoolVarP(&flagAll, "all", "a", false, "Run all 10 collection tools (default: true)")
	f.BoolVarP(&flagGau, "gau", "g", false, "Run GAU (GetAllUrls)")
	f.BoolVarP(&flagWaymore, "waymore", "m", false, "Run Waymore")
	f.BoolVar(&flagParamspider, "paramspider", false, "Run ParamSpider")
	f.BoolVar(&flagCommoncrawl, "commoncrawl", false, "Run Common Crawl CDX (native Go)")
	f.BoolVar(&flagUrlfinder, "urlfinder", false, "Run URLFinder")
	f.BoolVar(&flagGithubEndpoints, "github-endpoints", false, "Run github-endpoints")
	f.BoolVarP(&flagXnlinkfinder, "xnlinkfinder", "x", false, "Run xnLinkFinder")
	f.BoolVarP(&flagKatana, "katana", "k", false, "Run Katana")
	f.BoolVar(&flagHakrawler, "hakrawler", false, "Run Hakrawler")
	f.BoolVarP(&flagGobuster, "gobuster", "z", false, "Run Gobuster")

	rootCmd.AddCommand(doctorCmd)
}

// resolveTargets parses positional arguments and -f file inputs.
// Automatically detects if a positional argument is a text file and reads targets from it!
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

	// Add positional arguments
	for _, a := range args {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		// If positional argument is a text file on disk, read targets line by line!
		if utils.FileExists(a) {
			lines, err := utils.ReadLines(a)
			if err == nil && len(lines) > 0 {
				for _, l := range lines {
					add(l)
				}
				continue
			}
		}
		// Otherwise, treat as domain or URL
		add(a)
	}

	// Add file passed via -f / --file flag
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
