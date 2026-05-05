package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"urlshine/internal/runner"
	"urlshine/internal/utils"

	"github.com/spf13/cobra"
)

var (
	flagFile      string
	flagOutputDir string
	flagThreads   int
	flagDepth     int
	flagSubs      bool
	flagNoAlive   bool
	flagSkipCol   bool
	flagVerbose   bool

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
	Short: "URLShine v2 — professional URL enumeration & reconnaissance",
	Long:  "URLShine collects URLs from 9 tools with aggressive parallel execution, normalizes, splits into groups, checks alive, and generates professional reports.",
	Example: `  urlshine google.com
  urlshine google.com yahoo.com
  urlshine -f domains.txt
  urlshine -f domains.txt -o ./results -t 50
  urlshine google.com --no-alive -d 2
  urlshine google.com --verbose`,
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

		// Info logged by RunProfessional

		return runner.RunProfessional(runner.Options{
			Targets:     targets,
			OutputDir:   outDir,
			Threads:     flagThreads,
			Depth:       flagDepth,
			Subs:        flagSubs,
			SkipAlive:   flagNoAlive,
			SkipCollect: flagSkipCol,
			Verbose:     flagVerbose,

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
	// Professional custom help menu
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		fmt.Println("  ╔════════════════════════════════════════════════════════════╗")
		fmt.Println("  ║        URLShine v2.0.0 - Professional Help Menu             ║")
		fmt.Println("  ║     Aggressive URL Enumeration & Attack Surface Mapper      ║")
		fmt.Println("  ╚════════════════════════════════════════════════════════════╝")
		fmt.Println("")

		fmt.Println("  📖 OVERVIEW")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Println("  URLShine automates large-scale URL enumeration from 9 advanced tools:")
		fmt.Println("  GAU, Katana, GoSpider, Waymore, Wayback URLs, Hakrawler, xnLinkFinder,")
		fmt.Println("  Gobuster, and Dirb for comprehensive coverage.")
		fmt.Println("")
		fmt.Println("  Features:")
		fmt.Println("    • Creates domain-specific folders for organized results")
		fmt.Println("    • Per-tool URL files for individual tool analysis")
		fmt.Println("    • Automatic URL normalization and deduplication")
		fmt.Println("    • Categorizes URLs: API, Auth, Parameters, JS/Config, Directories")
		fmt.Println("    • Professional progress tracking and reporting")
		fmt.Println("    • Aggressive parallel execution for maximum performance")
		fmt.Println("")

		fmt.Println("  🚀 QUICK START")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Println("    # Single domain - auto-runs all 9 tools")
		fmt.Println("    urlshine google.com")
		fmt.Println("")
		fmt.Println("    # Specific tools only")
		fmt.Println("    urlshine -gau -katana google.com")
		fmt.Println("")
		fmt.Println("    # All tools with aggressive settings (includes -gobuster -dirb)")
		fmt.Println("    urlshine -all -t 100 -d 5 google.com")
		fmt.Println("")
		fmt.Println("    # Process file with all tools")
		fmt.Println("    urlshine -f targets.txt -all -t 150 -d 3 -o ./results")
		fmt.Println("")
		fmt.Println("    # Fast mode (skip alive checking)")
		fmt.Println("    urlshine -all -no-alive -t 100 google.com")
		fmt.Println("")

		fmt.Println("  📋 USAGE")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Println("    urlshine [target ...] [flags]")
		fmt.Println("")

		fmt.Println("  🎯 TARGET OPTIONS")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Printf("    %-25s %s\n", "urlshine google.com", "Single domain")
		fmt.Printf("    %-25s %s\n", "urlshine target1 target2", "Multiple domains")
		fmt.Printf("    %-25s %s\n", "urlshine -f domains.txt", "Domains from file")
		fmt.Println("")

		fmt.Println("  ⚙️  CONFIGURATION FLAGS")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Printf("    %-25s %s\n", "-o, --output", "Output directory (default: urlshine_<timestamp>)")
		fmt.Printf("    %-25s %s\n", "-t, --threads", "Threads for collection (default: 50, max: 200)")
		fmt.Printf("    %-25s %s\n", "-d, --depth", "Crawl depth for active tools (default: 5)")
		fmt.Printf("    %-25s %s\n", "-s, --subs", "Include subdomains (default: true)")
		fmt.Println("")

		fmt.Println("  🔧 URL COLLECTION TOOLS")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Printf("    %-25s %s\n", "-all", "Run all 9 tools + gobuster + dirb automatically")
		fmt.Printf("    %-25s %s\n", "-gau", "GAU (GetAllUrls) - archive sources")
		fmt.Printf("    %-25s %s\n", "-katana", "Katana - JS crawler")
		fmt.Printf("    %-25s %s\n", "-gospider", "GoSpider - HTML & JS crawler")
		fmt.Printf("    %-25s %s\n", "-waymore", "Waymore - advanced wayback scraper")
		fmt.Printf("    %-25s %s\n", "-waybackurls", "Wayback URLs - wayback machine")
		fmt.Printf("    %-25s %s\n", "-hakrawler", "Hakrawler - HTML crawler")
		fmt.Printf("    %-25s %s\n", "-xnlinkfinder", "xnLinkFinder - JS extractor")
		fmt.Printf("    %-25s %s\n", "-gobuster", "Gobuster - directory brute-force")
		fmt.Printf("    %-25s %s\n", "-dirb", "Dirb - directory enumeration")
		fmt.Println("")

		fmt.Println("  🔍 PROCESSING OPTIONS")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Printf("    %-25s %s\n", "-no-alive", "Skip HTTP alive verification (fast mode)")
		fmt.Printf("    %-25s %s\n", "-skip-collect", "Reprocess existing data (skip collection)")
		fmt.Printf("    %-25s %s\n", "-v, --verbose", "Debug/verbose logging output")
		fmt.Println("")

		fmt.Println("  📁 OUTPUT STRUCTURE")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Println("    For domain 'google.com', creates:")
		fmt.Println("")
		fmt.Println("    google_url/")
		fmt.Println("    ├── gau.txt                    (GAU collected URLs)")
		fmt.Println("    ├── katana.txt                 (Katana crawler results)")
		fmt.Println("    ├── gospider.txt               (GoSpider results)")
		fmt.Println("    ├── waymore.txt                (Waymore results)")
		fmt.Println("    ├── waybackurls.txt            (Wayback URLs results)")
		fmt.Println("    ├── hakrawler.txt              (Hakrawler results)")
		fmt.Println("    ├── xnlinkfinder.txt           (xnLinkFinder results)")
		fmt.Println("    ├── gobuster.txt               (Gobuster directory discovery)")
		fmt.Println("    ├── dirb.txt                   (Dirb directory brute-force)")
		fmt.Println("    ├── merged_urls.txt            (All URLs combined)")
		fmt.Println("    ├── normalized_urls.txt        (Cleaned & deduplicated)")
		fmt.Println("    ├── api_urls.txt               (API endpoints)")
		fmt.Println("    ├── auth_admin_urls.txt        (Auth & admin pages)")
		fmt.Println("    ├── params_urls.txt            (URLs with parameters)")
		fmt.Println("    ├── js_config_urls.txt         (JS & config files)")
		fmt.Println("    └── directories_urls.txt       (Directory paths)")
		fmt.Println("")

		fmt.Println("  📊 COMMAND EXAMPLES")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Println("    Small scope (1-100 targets):")
		fmt.Println("      urlshine -all -t 50 -d 5 -f targets.txt")
		fmt.Println("")
		fmt.Println("    Medium scope (100-1000 targets):")
		fmt.Println("      urlshine -all -t 100 -d 3 -f targets.txt")
		fmt.Println("")
		fmt.Println("    Large scope (1000+ targets):")
		fmt.Println("      urlshine -all -t 150 -d 2 -no-alive -f targets.txt")
		fmt.Println("")
		fmt.Println("    Selective tools:")
		fmt.Println("      urlshine -gau -katana -gobuster -t 50 google.com")
		fmt.Println("")

		fmt.Println("  💡 ABOUT -all FLAG")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Println("    When using -all, URLShine automatically:")
		fmt.Println("      ✓ Runs all 9 URL collection tools")
		fmt.Println("      ✓ Includes Gobuster for directory discovery")
		fmt.Println("      ✓ Includes Dirb for directory brute-forcing")
		fmt.Println("      ✓ Performs URL normalization and deduplication")
		fmt.Println("      ✓ Advanced extraction for 5 specialized groups")
		fmt.Println("      ✓ Optional live host verification (disable with -no-alive)")
		fmt.Println("")
		fmt.Println("  📚 RESOURCES")
		fmt.Println("  ──────────────────────────────────────────────────────────────")
		fmt.Println("    GitHub: https://github.com/shii9/UrlShine")
		fmt.Println("    Docs:   See README.md for detailed documentation")
		fmt.Println("")
	})

	f := rootCmd.Flags()
	f.StringVarP(&flagFile, "file", "f", "", "Input file with targets (one per line)")
	f.StringVarP(&flagOutputDir, "output", "o", "", "Output directory (default: urlshine_<timestamp>)")
	f.IntVarP(&flagThreads, "threads", "t", 50, "Parallel threads (default: 50, max: 200+)")
	f.IntVarP(&flagDepth, "depth", "d", 5, "Crawl depth for active tools (default: 5)")
	f.BoolVarP(&flagSubs, "subs", "s", true, "Include subdomains (default: true)")
	f.BoolVar(&flagNoAlive, "no-alive", false, "Skip live host verification (fast mode)")
	f.BoolVar(&flagSkipCol, "skip-collect", false, "Reprocess existing data (skip collection)")
	f.BoolVarP(&flagVerbose, "verbose", "v", false, "Debug/verbose logging output")

	f.BoolVar(&flagAll, "all", false, "Run all 9 tools + gobuster + dirb")
	f.BoolVar(&flagGau, "gau", false, "Run GAU (GetAllUrls)")
	f.BoolVar(&flagGospider, "gospider", false, "Run GoSpider")
	f.BoolVar(&flagKatana, "katana", false, "Run Katana")
	f.BoolVar(&flagWaymore, "waymore", false, "Run Waymore")
	f.BoolVar(&flagWaybackurls, "waybackurls", false, "Run Waybackurls")
	f.BoolVar(&flagHakrawler, "hakrawler", false, "Run Hakrawler")
	f.BoolVar(&flagXnlinkfinder, "xnlinkfinder", false, "Run xnLinkFinder")
	f.BoolVar(&flagGobuster, "gobuster", false, "Run Gobuster (directory discovery)")
	f.BoolVar(&flagDirb, "dirb", false, "Run Dirb (directory brute-forcing)")
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
