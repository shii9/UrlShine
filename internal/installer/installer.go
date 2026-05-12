// Package installer handles runtime tool installation and verification.
package installer

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

var tools = map[string]string{
	"gau":          "github.com/lc/gau/v2/cmd/gau@latest",
	"katana":       "github.com/projectdiscovery/katana/cmd/katana@latest",
	"gospider":     "github.com/jaeles-project/gospider@latest",
	"waymore":      "waymore",
	"waybackurls":  "github.com/tomnomnom/waybackurls@latest",
	"hakrawler":    "github.com/hakluke/hakrawler@latest",
	"xnLinkFinder": "xnlinkfinder",
	"gobuster":     "github.com/OJ/gobuster/v3@latest",
	"httpx":        "github.com/projectdiscovery/httpx/cmd/httpx@latest",
}

// CheckAndInstall verifies if tools are installed and attempts installation if missing
func CheckAndInstall() {
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println()
	fmt.Printf("%s Verifying required tools...\n\n", cyan("▸"))

	missing := []string{}
	for tool := range tools {
		if !isToolAvailable(tool) {
			missing = append(missing, tool)
			fmt.Printf("%s  %s not found\n", yellow("○"), tool)
		}
	}

	if len(missing) == 0 {
		fmt.Printf("%s All tools ready\n\n", green("✔"))
		return
	}

	fmt.Printf("\n%s Some tools are missing. Install now? [Y/n] ", yellow("!"))

	// Auto-install if running in non-interactive mode or first run
	// For now, show the command to install
	fmt.Printf("\nRun this to install missing tools:\n")
	fmt.Printf("%s bash install.sh\n", cyan("$"))
	fmt.Printf("Or visit: https://github.com/shii9/UrlShine#-installation\n\n")
}

// isToolAvailable checks if a tool is available in PATH
func isToolAvailable(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// ShowProgress displays animated progress while tools are running
func ShowProgress(tool string, target string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	frames := []string{"▸", "▹", "▿"}

	for i := 0; i < 3; i++ {
		fmt.Printf("\r%s %s scanning %s...", cyan(frames[i%len(frames)]), tool, target)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Print("\r")
}
