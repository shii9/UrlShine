// URLShine — Professional URL Enumeration & Attack Surface Mapper
//
// URLShine orchestrates a sophisticated reconnaissance pipeline combining multiple
// URL enumeration tools into a unified workflow. It collects URLs from passive and
// active sources, deduplicates results, categorizes findings by attack vector, and
// verifies live hosts.
//
// Usage:
//
//	urlshine [target ...] [flags]
//
// Examples:
//
//	urlshine -a google.com                          # All tools, collection only
//	urlshine -a -c google.com                       # All tools, full pipeline
//	urlshine -f targets.txt -a -c -t 150 -d 5       # Aggressive multi-target scan
//	urlshine -gau -katana -c google.com             # Specific tools only
//
// For detailed help, run: urlshine -h
package main

import "github.com/shii9/UrlShine/cmd"

func main() {
	cmd.Execute()
}
