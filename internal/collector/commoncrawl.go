package collector

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shii9/UrlShine/internal/logger"
)

// ccIndexResponse represents a Common Crawl CDX index listing.
type ccIndexResponse struct {
	ID     string `json:"id"`
	APIURL string `json:"cdx-api"`
}

// runCommoncrawl queries the Common Crawl CDX API directly (native Go, no external tool).
// This is a zero-dependency approach that talks to the CC index servers directly.
func runCommoncrawl(target, _ string, cfg Config) ([]string, error) {
	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	// Get available CC indices
	indices, err := getCCIndices(client)
	if err != nil {
		logger.Debug("CC: failed to get indices: %v", err)
		// Fall back to the most recent known index
		indices = []string{
			"https://index.commoncrawl.org/CC-MAIN-2024-10-index",
			"https://index.commoncrawl.org/CC-MAIN-2023-50-index",
			"https://index.commoncrawl.org/CC-MAIN-2023-40-index",
		}
	}

	seen := make(map[string]struct{})
	var allUrls []string

	// Query the most recent indices (limit to 5 for speed)
	maxIndices := 5
	if len(indices) < maxIndices {
		maxIndices = len(indices)
	}

	for _, indexURL := range indices[:maxIndices] {
		urls, err := queryCCIndex(client, indexURL, target)
		if err != nil {
			logger.Debug("CC: index %s failed: %v", indexURL, err)
			continue
		}

		for _, u := range urls {
			if _, ok := seen[u]; !ok {
				seen[u] = struct{}{}
				allUrls = append(allUrls, u)
			}
		}
	}

	return allUrls, nil
}

// getCCIndices fetches the list of available Common Crawl indices.
func getCCIndices(client *http.Client) ([]string, error) {
	resp, err := client.Get("https://index.commoncrawl.org/collinfo.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var indices []ccIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&indices); err != nil {
		return nil, err
	}

	var urls []string
	for _, idx := range indices {
		if idx.APIURL != "" {
			urls = append(urls, strings.TrimSuffix(idx.APIURL, "/"))
		}
	}
	return urls, nil
}

// queryCCIndex queries a single Common Crawl index for URLs matching the target.
func queryCCIndex(client *http.Client, indexURL, target string) ([]string, error) {
	// Query the CDX API with wildcard domain matching
	queryURL := fmt.Sprintf("%s?url=*.%s&output=text&fl=url&limit=10000",
		indexURL, url.QueryEscape(target))

	resp, err := client.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("CC index returned status %d", resp.StatusCode)
	}

	// Read response line by line
	var urls []string
	reader := bufio.NewReaderSize(resp.Body, 4*1024*1024)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line != "" && (strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://")) {
			urls = append(urls, line)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return urls, err
		}
	}

	return urls, nil
}
