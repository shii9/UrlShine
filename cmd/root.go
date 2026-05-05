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
		"Use -all to run every collector and -complete to run the full post-collection pipeline.",
	Example: `  urlshine -gau -katana google.com
  urlshine -gau -katana -complete google.com
  urlshine -all -complete google.com
  urlshine -f domains.txt -all -complete -o ./results -t 100
  urlshine -all -complete -no-alive google.com
  urlshine -katana -complete -v google.com`,
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
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Printf(`
URLShine v%s
Professional URL Enumeration and Attack Surface Mapper

USAGE
  urlshine [target ...] [flags]

FLAGS EXPLAINED
  -all        Use all tools (GAU, Katana, GoSpider, Waymore, Waybackurls, 
              Hakrawler, xnLinkFinder, Gobuster, Dirb)
  -complete   Complete all processing steps:
              • Merging — Deduplicates all results
              • Normalization — Cleans URLs
              • Categorization — Splits into 5 attack groups
              • Alive Checking — Verifies live hosts (unless -no-alive used)

COLLECTION TOOLS
  -gau            GetAllUrls - archive & passive sources
  -katana         Katana - active JS crawler
  -gospider       GoSpider - HTML & JS crawler
  -waymore        Waymore - advanced wayback scraper
  -waybackurls    Wayback URLs - wayback machine scraper
  -hakrawler      Hakrawler - HTML content crawler
  -xnlinkfinder   xnLinkFinder - JS endpoint extractor
  -gobuster       Gobuster - directory brute-force discovery
  -dirb           Dirb - directory enumeration

EXAMPLES
  urlshine -gau -katana google.com
      Collect URLs with GAU and Katana only (no post-processing)

  urlshine -gau -katana -complete google.com
      Collect with GAU and Katana, then run merging, normalization,
      categorization, and alive checking

  urlshine -all google.com
      Collect with all 9 tools (no post-processing)

  urlshine -all -complete google.com
      Collect with all 9 tools, then run complete processing pipeline

OPTIONS
  -t, --threads INT       Parallel threads (default: 50)
  -d, --depth INT         Crawl depth for active tools (default: 5)
  -o, --output DIR        Output directory
  -f, --file FILE         Input file with targets (one per line)
  -no-alive              Skip live host verification
  -skip-collect          Reprocess existing data
  -v, --verbose          Debug/verbose logging

OUTPUT STRUCTURE
  Without -complete:
    {domain}_url/
    ├── gau.txt
    ├── katana.txt
    ├── gospider.txt
    └── ... (per-tool files)

  With -complete:
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

FLAGS
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
	f.BoolVar(&flagComplete, "complete", false, "Run merge, normalize, categorize, reports, and alive checking")

	f.BoolVar(&flagAll, "all", false, "Run all supported collection tools")
	f.BoolVar(&flagGau, "gau", false, "Run GAU (GetAllUrls)")
	f.BoolVar(&flagGospider, "gospider", false, "Run GoSpider")
	f.BoolVar(&flagKatana, "katana", false, "Run Katana")
	f.BoolVar(&flagWaymore, "waymore", false, "Run Waymore")
	f.BoolVar(&flagWaybackurls, "waybackurls", false, "Run Waybackurls")
	f.BoolVar(&flagHakrawler, "hakrawler", false, "Run Hakrawler")
	f.BoolVar(&flagXnlinkfinder, "xnlinkfinder", false, "Run xnLinkFinder")
	f.BoolVar(&flagGobuster, "gobuster", false, "Run Gobuster directory discovery")
	f.BoolVar(&flagDirb, "dirb", false, "Run Dirb directory enumeration")
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
