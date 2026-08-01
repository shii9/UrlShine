package cmd

import (
	"fmt"
	"runtime"

	"github.com/shii9/UrlShine/internal/utils"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Dependency auditing and system health check",
	Long: `URLShine Doctor performs a comprehensive system audit:
  • Verifies all 10 core URL enumeration tools are installed
  • Checks Go version and system compatibility
  • Categorizes tools by type (passive archives, passive APIs, active crawlers, brute-force)
  • Provides intelligent installation recommendations`,
	Example: `  urlshine doctor
  urlshine doctor --verbose`,
	Run: func(cmd *cobra.Command, args []string) {
		runDoctor()
	},
}

func runDoctor() {
	fmt.Printf(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                    URLShine Dependency Audit & Health Check                  ║
╚══════════════════════════════════════════════════════════════════════════════╝

`)

	// System Information
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 SYSTEM INFORMATION")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  OS:         %s\n", runtime.GOOS)
	fmt.Printf("  Arch:       %s\n", runtime.GOARCH)
	fmt.Printf("  Go Version: %s\n", runtime.Version())
	fmt.Printf("  URLShine:   v%s\n", version)
	fmt.Println()

	statuses := utils.CheckDependencies()

	// Categorize tools into 4 tiers + utility
	passiveArchiveTools := []string{"gau", "waymore", "paramspider"}
	passiveAPITools := []string{"commoncrawl", "urlfinder", "github-endpoints", "xnlinkfinder"}
	activeCrawlerTools := []string{"katana", "hakrawler"}
	activeBruteTools := []string{"gobuster"}
	utilityTools := []string{"httpx"}

	var installedArchives, installedAPIs, installedCrawlers, installedBrute, installedUtil int

	for _, s := range statuses {
		if s.Status == "installed" {
			for _, p := range passiveArchiveTools {
				if s.Name == p {
					installedArchives++
				}
			}
			for _, a := range passiveAPITools {
				if s.Name == a {
					installedAPIs++
				}
			}
			for _, c := range activeCrawlerTools {
				if s.Name == c {
					installedCrawlers++
				}
			}
			for _, b := range activeBruteTools {
				if s.Name == b {
					installedBrute++
				}
			}
			for _, u := range utilityTools {
				if s.Name == u {
					installedUtil++
				}
			}
		}
	}

	totalTools := len(statuses)
	installedCount := installedArchives + installedAPIs + installedCrawlers + installedBrute + installedUtil
	missingCount := totalTools - installedCount

	fmt.Printf("Total tools: %d | Installed: %d (Archives: %d, APIs: %d, Crawlers: %d, Brute-Force: %d, Utility: %d) | Missing: %d\n\n",
		totalTools, installedCount, installedArchives, installedAPIs, installedCrawlers, installedBrute, installedUtil, missingCount)

	// Display passive archive tools
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔍 PASSIVE ARCHIVES (Archive-based, zero target traffic)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, s := range statuses {
		for _, p := range passiveArchiveTools {
			if s.Name == p {
				if s.Status == "installed" {
					fmt.Printf("  ✓ %-18s [Installed]\n", s.Name)
				} else {
					fmt.Printf("  ✗ %-18s [Missing]\n", s.Name)
				}
			}
		}
	}

	// Display passive API tools
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📡 PASSIVE APIS & OSINT (External APIs/Feeds, zero target traffic)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, s := range statuses {
		for _, a := range passiveAPITools {
			if s.Name == a {
				if s.Status == "installed" {
					fmt.Printf("  ✓ %-18s [Installed]\n", s.Name)
				} else {
					fmt.Printf("  ✗ %-18s [Missing]\n", s.Name)
				}
			}
		}
	}

	// Display active crawlers
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🕷️  ACTIVE CRAWLERS (Crawlers, generates target traffic)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, s := range statuses {
		for _, c := range activeCrawlerTools {
			if s.Name == c {
				if s.Status == "installed" {
					fmt.Printf("  ✓ %-18s [Installed]\n", s.Name)
				} else {
					fmt.Printf("  ✗ %-18s [Missing]\n", s.Name)
				}
			}
		}
	}

	// Display active brute force tools
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💥 ACTIVE BRUTE-FORCE (Directory/File enumeration)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, s := range statuses {
		for _, b := range activeBruteTools {
			if s.Name == b {
				if s.Status == "installed" {
					fmt.Printf("  ✓ %-18s [Installed]\n", s.Name)
				} else {
					fmt.Printf("  ✗ %-18s [Missing]\n", s.Name)
				}
			}
		}
	}

	// Display utility tools
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔧 UTILITY TOOLS (Optional, improves performance)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, s := range statuses {
		for _, u := range utilityTools {
			if s.Name == u {
				if s.Status == "installed" {
					fmt.Printf("  ✓ %-18s (faster HTTP probing)\n", s.Name)
				} else {
					fmt.Printf("  ✗ %-18s [Missing]\n", s.Name)
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("❌ MISSING TOOLS INSTALLATION COMMANDS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	missingFound := 0
	for _, s := range statuses {
		if s.Status == "missing" {
			fmt.Printf("  ✗ %s\n", s.Name)
			fmt.Printf("    Install: %s\n", s.InstallCmd)
			missingFound++
		}
	}
	if missingFound == 0 {
		fmt.Println("  (all 10 core tools installed!)")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 RECOMMENDATIONS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if installedCount == 0 {
		fmt.Println("  ⚠️  No tools installed. Run the installer to get started:")
		if runtime.GOOS == "windows" {
			fmt.Println("    install.bat")
		} else {
			fmt.Println("    bash install.sh")
		}
	} else if installedArchives+installedAPIs > 0 && installedCrawlers == 0 {
		fmt.Printf("  ✓ You have %d passive tools. Recommended: add active crawlers (katana, hakrawler)\n", installedArchives+installedAPIs)
		fmt.Println("    Usage: urlshine -a -c target.com")
	} else if installedCount > 0 {
		fmt.Printf("  ✓ Great setup! You have %d tools installed out of %d core tools.\n", installedCount, totalTools)
		fmt.Println("    Usage: urlshine -a -c target.com")
		if installedUtil == 0 {
			fmt.Println("    Tip: Install httpx for faster live host verification")
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if missingFound > 0 {
		fmt.Printf("\n⚠️  You have %d missing tools. URLShine will skip unavailable tools gracefully.\n", missingFound)
		fmt.Println("    For optimal results, install missing tools for maximum URL coverage.\n")
	} else {
		fmt.Println("\n✅ All 10 core tools are installed! You're ready for clean, maximum URL enumeration.\n")
	}
}
