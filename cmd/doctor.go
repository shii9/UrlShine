package cmd

import (
	"fmt"
	"runtime"

	"github.com/shii9/UrlShine/internal/utils"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check and verify tool dependencies",
	Long: `Check if all URLShine dependencies are installed.
Shows which tools are available and provides installation commands for missing ones.`,
	Example: `  urlshine doctor
  urlshine doctor --verbose`,
	Run: func(cmd *cobra.Command, args []string) {
		runDoctor()
	},
}

func runDoctor() {
	fmt.Printf(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                        URLShine Dependency Check                             ║
╚══════════════════════════════════════════════════════════════════════════════╝

`)

	statuses := utils.CheckDependencies()

	// Count installed and missing
	installed := 0
	missing := 0
	for _, s := range statuses {
		if s.Status == "installed" {
			installed++
		} else {
			missing++
		}
	}

	fmt.Printf("Total tools: %d | Installed: %d | Missing: %d\n\n", len(statuses), installed, missing)

	// Display by status
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ INSTALLED TOOLS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	installedCount := 0
	for _, s := range statuses {
		if s.Status == "installed" {
			fmt.Printf("  ✓ %s\n", s.Name)
			installedCount++
		}
	}
	if installedCount == 0 {
		fmt.Println("  (none installed)")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("❌ MISSING TOOLS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	missingCount := 0
	for _, s := range statuses {
		if s.Status == "missing" {
			fmt.Printf("  ✗ %s\n", s.Name)
			fmt.Printf("    Install: %s\n", s.InstallCmd)
			missingCount++
		}
	}
	if missingCount == 0 {
		fmt.Println("  (all tools installed!)")
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

	if missingCount > 0 {
		fmt.Printf("\n⚠️  You have %d missing tools. URLShine will skip unavailable tools gracefully.\n\n", missingCount)
	} else {
		fmt.Println("\n✅ All tools are installed! You're ready to go.\n")
	}
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
