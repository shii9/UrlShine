// Package installer handles runtime tool installation and verification.
package installer

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

var tools = map[string]string{
	// Tier 1: Passive Archives
	"gau":         "github.com/lc/gau/v2/cmd/gau@latest",
	"waymore":     "waymore",
	"paramspider": "git+https://github.com/devanshbatham/ParamSpider",

	// Tier 2: Passive APIs & OSINT
	"urlfinder":        "github.com/projectdiscovery/urlfinder/cmd/urlfinder@latest",
	"github-endpoints": "github.com/gwen001/github-endpoints@latest",
	"xnLinkFinder":     "xnlinkfinder",
	"commoncrawl":      "commoncrawl", // Not a true binary but listed for symmetry if checked

	// Tier 3: Active Crawlers
	"katana":    "github.com/projectdiscovery/katana/cmd/katana@latest",
	"hakrawler": "github.com/hakluke/hakrawler@latest",

	// Tier 4: Active Brute-Force
	"gobuster": "github.com/OJ/gobuster/v3@latest",
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

	fmt.Printf("\n%s Some tools are missing (%d tools). Install now? [Y/n] ", yellow("!"), len(missing))

	// Auto-install if running in non-interactive mode or first run
	// For now, show the command to install
	fmt.Printf("\nRun this to install missing tools:\n")
	fmt.Printf("%s bash install.sh\n", cyan("$"))
	fmt.Printf("Or visit: https://github.com/shii9/UrlShine#-installation\n\n")
}

// isToolAvailable checks if a tool is available in PATH
func isToolAvailable(tool string) bool {
	if tool == "commoncrawl" {
		return true // Built-in
	}
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
