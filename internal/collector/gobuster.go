package collector

import (
	"fmt"
	"strings"
)

// runGobuster discovers directories and files using gobuster.
func runGobuster(target, outDir string, cfg Config) ([]string, error) {
	// Normalize target to domain only (remove http/https)
	domain := strings.TrimPrefix(target, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimRight(domain, "/")

	// Use medium-sized wordlist for better coverage
	// gobuster -u https://domain.com -w /usr/share/wordlists/dirb/common.txt -t 50 -o results.txt
	args := []string{
		"gobuster", "dir",
		"-u", "https://" + domain,
		"-w", "/usr/share/wordlists/dirb/common.txt", // fallback to common wordlist
		"-t", fmt.Sprintf("%d", cfg.Threads),
		"-q",            // quiet mode (no banner)
		"--no-error",    // suppress errors
		"-k",            // skip SSL verification
		"--no-progress", // no progress bar
	}

	// Try to get output, convert to URLs
	lines, err := runCmd(args...)
	if err != nil {
		return nil, err
	}

	// Parse gobuster output and convert to full URLs
	var urls []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "[") {
			continue
		}

		// gobuster outputs paths like: /path (200)
		parts := strings.Fields(line)
		if len(parts) > 0 {
			path := parts[0]
			// Build full URL
			url := "https://" + domain + path
			urls = append(urls, url)
		}
	}

	return urls, nil
}
