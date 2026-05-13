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
  • Verifies all URL enumeration tools are installed
  • Checks Go version and system compatibility
  • Categorizes tools by type (passive/active)
  • Provides intelligent installation recommendations
  • Suggests optimal tool combinations based on available tools`,
	Example: `  urlshine doctor
  urlshine doctor --verbose
  urlshine doctor --recommend`,
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
	fmt.Printf("  OS:        %s\n", runtime.GOOS)
	fmt.Printf("  Arch:      %s\n", runtime.GOARCH)
	fmt.Printf("  Go Version: %s\n", runtime.Version())
	fmt.Printf("  URLShine:  v%s\n", version)
	fmt.Println()

	statuses := utils.CheckDependencies()

	// Categorize tools
	passiveTools := []string{"gau", "waymore", "waybackurls", "xnlinkfinder"}
	activeTools := []string{"katana", "gospider", "gobuster", "dirb"}
	utilityTools := []string{"httpx"}

	var installedPassive, installedActive, installedUtil int
	var missingPassive, missingActive, missingUtil int

	for _, s := range statuses {
		if s.Status == "installed" {
			for _, p := range passiveTools {
				if s.Name == p {
					installedPassive++
				}
			}
			for _, a := range activeTools {
				if s.Name == a {
					installedActive++
				}
			}
			for _, u := range utilityTools {
				if s.Name == u {
					installedUtil++
				}
			}
		} else {
			for _, p := range passiveTools {
				if s.Name == p {
					missingPassive++
				}
			}
			for _, a := range activeTools {
				if s.Name == a {
					missingActive++
				}
			}
			for _, u := range utilityTools {
				if s.Name == u {
					missingUtil++
				}
			}
		}
	}

	totalTools := len(statuses)
	installedCount := installedPassive + installedActive + installedUtil
	missingCount := totalTools - installedCount

	fmt.Printf("Total tools: %d | Installed: %d (passive: %d, active: %d, utility: %d) | Missing: %d\n\n",
		totalTools, installedCount, installedPassive, installedActive, installedUtil, missingCount)

	// Display passive tools (archives)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔍 PASSIVE TOOLS (Archive-based, no target traffic)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, s := range statuses {
		for _, p := range passiveTools {
			if s.Name == p {
				if s.Status == "installed" {
					fmt.Printf("  ✓ %s\n", s.Name)
				}
			}
		}
	}

	// Display active tools (crawlers)
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🕷️  ACTIVE TOOLS (Crawlers, generates target traffic)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, s := range statuses {
		for _, a := range activeTools {
			if s.Name == a {
				if s.Status == "installed" {
					fmt.Printf("  ✓ %s\n", s.Name)
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
					fmt.Printf("  ✓ %s (faster HTTP probing)\n", s.Name)
				}
			}
		}
	}

	// ✅ INSTALLED TOOLS section
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ INSTALLED TOOLS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	installedFound := 0
	for _, s := range statuses {
		if s.Status == "installed" {
			fmt.Printf("  ✓ %s\n", s.Name)
			installedFound++
		}
	}
	if installedFound == 0 {
		fmt.Println("  (none installed)")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("❌ MISSING TOOLS")
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
		fmt.Println("  (all tools installed!)")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 RECOMMENDATIONS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Provide recommendations based on installed tools
	if installedCount == 0 {
		fmt.Println("  ⚠️  No tools installed. Run the installer to get started:")
		if runtime.GOOS == "windows" {
			fmt.Println("    install.bat")
		} else {
			fmt.Println("    bash install.sh")
		}
	} else if installedPassive > 0 && installedActive == 0 {
		fmt.Printf("  ✓ You have %d passive tools. Recommended: add active crawlers (katana, gospider)\n", installedPassive)
		fmt.Println("    Usage: urlshine -a -c target.com")
	} else if installedPassive == 0 && installedActive > 0 {
		fmt.Printf("  ✓ You have %d active tools. Recommended: add passive tools (gau, waymore)\n", installedActive)
		fmt.Println("    Usage: urlshine -a -c target.com")
	} else if installedCount > 0 {
		fmt.Printf("  ✓ Great setup! You have %d tools installed.\n", installedCount)
		fmt.Println("    Usage: urlshine -a -c target.com")
		if installedUtil == 0 {
			fmt.Println("    Tip: Install httpx for faster live host verification")
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🚀 QUICK INSTALLATION")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if runtime.GOOS == "windows" {
		fmt.Println()
		fmt.Println("  Windows: Run the installer batch file:")
		fmt.Println("    install.bat")
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("  Linux/macOS: Run the installer script:")
		fmt.Println("    bash install.sh")
		fmt.Println()
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if missingFound > 0 {
		fmt.Printf("\n⚠️  You have %d missing tools. URLShine will skip unavailable tools gracefully.\n", missingFound)
		fmt.Println("    For optimal results, install all tools using the installer scripts.\n")
	} else {
		fmt.Println("\n✅ All tools are installed! You're ready for comprehensive URL enumeration.\n")
	}
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
