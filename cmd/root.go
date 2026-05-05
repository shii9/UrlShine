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

var (
	flagFile      string
	flagOutputDir string
	flagThreads   int
	flagDepth     int
	flagSubs      bool
	flagNoAlive   bool
	flagSkipCol   bool
	flagVerbose   bool
	flagComplete  bool

	flagAll          bool
	flagGau          bool
	flagGospider     bool
	flagKatana       bool
	flagWaymore      bool
	flagWaybackurls  bool
	flagHakrawler    bool
	flagXnlinkfinder bool
	flagGobuster     bool
	flagDirb         bool
)

var rootCmd = &cobra.Command{
	Use:   "urlshine [domain ...] [-f domains.txt]",
	Short: "URLShine v2.0.0 - professional URL enumeration and reconnaissance",
	Long: "URLShine collects URLs from selected reconnaissance tools. " +
		"Use --all (or -a) to run every collector and --complete (or -c) to run the full post-collection pipeline.",
	Example: `  urlshine --gau --katana google.com
  urlshine -g -k google.com
  urlshine --gau --katana --complete google.com
  urlshine -g -k -c google.com
  urlshine --all google.com
  urlshine -a google.com
  urlshine --all --complete google.com
  urlshine -a -c google.com
  urlshine -f domains.txt -a -c -o ./results -t 100
  urlshine -a -c --no-alive google.com
  urlshine -k -c -v google.com`,
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
URLShine v%s
Professional URL Enumeration and Attack Surface Mapper

USAGE
  urlshine [target ...] [flags]

✅ FLEXIBLE FLAG FORMAT
  You can use ANY of these formats (all work the same):
  
  Long flags (double dash):
    urlshine --all --complete google.com
    
  Short flags (single dash with letter):
    urlshine -a -c google.com
    
  Long flags (single dash):
    urlshine -all -complete google.com
    
  Mixed formats:
    urlshine -a --complete google.com
    urlshine --gau -katana -c google.com

FLAGS EXPLAINED
  --all, -a, -all         Use all 9 tools (GAU, Katana, GoSpider, Waymore, Waybackurls, 
                          Hakrawler, xnLinkFinder, Gobuster, Dirb)
  --complete, -c, -complete  Complete all processing steps:
                            • Merging — Deduplicates all results
                            • Normalization — Cleans URLs
                            • Categorization — Splits into 5 attack groups
                            • Alive Checking — Verifies live hosts (unless --no-alive or -no-alive)

COLLECTION TOOLS (use any format)
  --gau, -g, -gau              GetAllUrls - archive & passive sources
  --katana, -k, -katana        Katana - active JS crawler
  --gospider, -w, -gospider    GoSpider - HTML & JS crawler
  --waymore, -m, -waymore      Waymore - advanced wayback scraper
  --waybackurls, -b, -waybackurls  Wayback URLs - wayback machine scraper
  --hakrawler, -r, -hakrawler  Hakrawler - HTML content crawler
  --xnlinkfinder, -x, -xnlinkfinder  xnLinkFinder - JS endpoint extractor
  --gobuster, -u, -gobuster    Gobuster - directory brute-force discovery
  --dirb, -i, -dirb            Dirb - directory enumeration

EXAMPLES
  All of these work:
    urlshine -all -complete google.com
    urlshine --all --complete google.com
    urlshine -a -c google.com
    urlshine -gau -katana google.com
    urlshine --gau --katana google.com
    urlshine -g -k google.com
    urlshine -f targets.txt -a -c -t 100 google.com
    urlshine --file targets.txt --all --complete --threads 100 google.com

OPTIONS
  -t, --threads, -threads      Parallel threads for tools and probing (default: 50)
  -d, --depth, -depth          Crawl depth for active tools (default: 5)
  -o, --output, -output        Output directory (default: urlshine_<timestamp>)
  -f, --file, -file            Input file with targets (one per line)
  -s, --subs, -subs            Include subdomains when supported (default: true)
  -v, --verbose, -verbose      Debug/verbose logging
  --no-alive, -no-alive        Skip live host verification
  --skip-collect, -skip-collect   Skip collection and process existing files

OUTPUT STRUCTURE
  Without --complete or -complete or -c:
    {domain}_url/
    ├── gau.txt
    ├── katana.txt
    ├── gospider.txt
    ├── waymore.txt
    ├── waybackurls.txt
    ├── hakrawler.txt
    ├── xnlinkfinder.txt
    ├── gobuster.txt
    └── dirb.txt

  With --complete or -complete or -c:
    {domain}_url/
    ├── merged_urls.txt (all tools combined)
    ├── normalized_urls.txt (cleaned & deduplicated)
    ├── api_endpoints.txt (API paths)
    ├── auth_admin_urls.txt (authentication pages)
    ├── parameters.txt (URLs with parameters)
    ├── js_config.txt (JavaScript & config files)
    ├── directories.txt (directory paths)
    ├── alive_urls.txt (verified live hosts)
    ├── report.json
    └── report.html

ALL FLAGS
`, version)
		cmd.Flags().PrintDefaults()
		fmt.Println("")
	})

	f := rootCmd.Flags()
	f.StringVarP(&flagFile, "file", "f", "", "Input file with targets (one per line)")
	f.StringVarP(&flagOutputDir, "output", "o", "", "Output directory (default: urlshine_<timestamp>)")
	f.IntVarP(&flagThreads, "threads", "t", 50, "Parallel threads for tools and probing")
	f.IntVarP(&flagDepth, "depth", "d", 5, "Crawl depth for active tools")
	f.BoolVarP(&flagSubs, "subs", "s", true, "Include subdomains when supported")
	f.BoolVar(&flagNoAlive, "no-alive", false, "Skip live host verification during -complete")
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
