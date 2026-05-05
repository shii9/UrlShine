// Package normalizer cleans and sanitizes URLs.
package normalizer

import (
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/shii9/UrlShine/internal/utils"
)

var staticExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
	".svg": true, ".ico": true, ".webp": true, ".tiff": true,
	".css": true, ".eot": true, ".ttf": true, ".woff": true, ".woff2": true,
	".pdf": true, ".zip": true, ".rar": true, ".tar": true, ".gz": true, ".7z": true,
	".mp4": true, ".mp3": true, ".avi": true, ".mkv": true, ".mov": true,
	".swf": true, ".flv": true, ".webm": true,
}

var doubleSlash = regexp.MustCompile(`([^:])/{2,}`)

// NormalizeFile reads raw URLs, normalizes them, and writes to outFile.
func NormalizeFile(inFile, outFile string) (int, int, error) {
	raw, err := utils.ReadLines(inFile)
	if err != nil {
		return 0, 0, err
	}
	clean := Normalize(raw)
	if err := utils.WriteLines(outFile, clean); err != nil {
		return 0, 0, err
	}
	return len(raw), len(clean), nil
}

// Normalize sanitizes URLs, removes static assets, and deduplicates.
func Normalize(raw []string) []string {
	seen := make(map[string]struct{}, len(raw))
	out := make([]string, 0, len(raw))

	for _, line := range raw {
		n := normalizeOne(strings.TrimSpace(line))
		if n == "" {
			continue
		}
		if _, ok := seen[n]; !ok {
			seen[n] = struct{}{}
			out = append(out, n)
		}
	}
	sort.Strings(out)
	return out
}

func normalizeOne(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return ""
	}

	host := strings.ToLower(u.Host)
	if host == "" {
		return ""
	}

	if scheme == "http" && strings.HasSuffix(host, ":80") {
		host = strings.TrimSuffix(host, ":80")
	} else if scheme == "https" && strings.HasSuffix(host, ":443") {
		host = strings.TrimSuffix(host, ":443")
	}

	path := doubleSlash.ReplaceAllString(u.Path, "$1/")
	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}

	ext := strings.ToLower(filepath.Ext(stripQuery(path)))
	if staticExt[ext] {
		return ""
	}

	rebuilt := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: u.RawQuery,
	}
	return rebuilt.String()
}

func stripQuery(s string) string {
	if i := strings.Index(s, "?"); i != -1 {
		return s[:i]
	}
	return s
}
