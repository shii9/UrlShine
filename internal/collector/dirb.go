package collector

import (
	"strings"
)

// runDirb discovers directories and files using dirb.
func runDirb(target, _ string, cfg Config) ([]string, error) {
	// Normalize target to domain only (remove http/https)
	domain := strings.TrimPrefix(target, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimRight(domain, "/")

	// dirb https://domain.com /usr/share/wordlists/dirb/common.txt -o output.txt -N 404
	args := []string{
		"dirb",
		"https://" + domain,
		"/usr/share/wordlists/dirb/common.txt",
		"-N", "404", // skip 404 responses
		"-r", // not recursive (avoid explosion)
	}

	// Try to get output
	lines, err := runCmd(args...)
	if err != nil {
		return nil, err
	}

	// Parse dirb output and convert to full URLs
	var urls []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "=") || strings.HasPrefix(line, "-") {
			continue
		}

		// dirb outputs lines like: http://domain.com/path (CODE)
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			// Extract URL (before the status code in parentheses)
			if idx := strings.Index(line, "("); idx > 0 {
				url := strings.TrimSpace(line[:idx])
				urls = append(urls, url)
			} else {
				urls = append(urls, line)
			}
		}
	}

	return urls, nil
}
