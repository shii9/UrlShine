# URLShine Development Guide

This document provides technical documentation for URLShine developers and maintainers.

## Architecture Overview

### Pipeline Stages

URLShine executes a 5-stage reconnaissance pipeline:

```
┌─────────────────────────────────────────────────────────────┐
│ STAGE 1: COLLECTION                                         │
│ - Concurrent execution of 9 URL enumeration tools          │
│ - Results saved to raw/ directory                          │
│ - Per-tool output files (gau.txt, katana.txt, etc.)       │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│ STAGE 2: MERGING                                            │
│ - Combine all tool outputs into merged_raw.txt            │
│ - Simple concatenation, maintains duplicates              │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│ STAGE 3: NORMALIZATION                                      │
│ - Remove static assets (.jpg, .png, .css, etc.)           │
│ - Strip redundant ports (:80, :443)                       │
│ - Deduplicate URLs via hash map                           │
│ - Validate URL format                                      │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│ STAGE 4: CATEGORIZATION                                     │
│ - Regex pattern matching for attack vectors               │
│ - 5 output categories:                                     │
│   * API endpoints                                          │
│   * Authentication & Admin pages                          │
│   * URLs with parameters                                   │
│   * JavaScript & config files                             │
│   * Directory paths                                        │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│ STAGE 5: VERIFICATION                                       │
│ - HTTP probe with httpx (or Go fallback)                  │
│ - Extract status codes, titles, tech stack               │
│ - Per-category alive verification                        │
└─────────────────────────────────────────────────────────────┘
```

## Core Packages

### `cmd/` — Command-Line Interface
- **root.go**: Main command definition using Cobra
- **doctor.go**: Dependency verification tool
- Handles flag parsing and pipeline orchestration

### `internal/banner/` — ASCII Art Header
- Displays professional banner on startup
- Shows version, system info, tool list
- All output is color-coded

### `internal/logger/` — Logging System
- Structured, leveled, color-coded logging
- Log levels: Debug, Info, Warn, Error, Silent
- Thread-safe with mutex protection
- Functions: `Info()`, `Success()`, `Warn()`, `Error()`, `Step()`, etc.

### `internal/collector/` — URL Collection
- Executes 9 external tools: GAU, Katana, GoSpider, Waymore, etc.
- Concurrent execution with semaphore-based thread pooling
- Per-tool configuration and timeout handling
- Returns count of URLs collected per tool

**Key functions:**
- `RunAll(targets, outputDir, config)` — Run all configured tools
- `Config` struct — Tool selection and performance settings

### `internal/merger/` — Output Aggregation
- Combines all tool outputs into single file
- Simple concatenation (deduplication happens in normalizer)
- Reports total URL count before normalization

**Key functions:**
- `MergeDir(inputDir, outputFile)` — Merge directory of files

### `internal/normalizer/` — URL Cleaning
- Removes static asset extensions
- Strips redundant ports
- Deduplicates using hash map
- Validates URL format

**Key functions:**
- `NormalizeFile(inputFile, outputFile)` — Process and deduplicate URLs
- Tracks input/output counts for statistics

### `internal/splitter/` — URL Categorization
- Regex pattern matching for attack vectors
- 5 categories: API, Auth, Params, JS/Config, Directories
- Per-category output files

**Key functions:**
- `Split(inputFile, outputDir)` — Categorize URLs
- Category constants: `GroupAPI`, `GroupAuth`, `GroupParams`, etc.

### `internal/alive/` — Live Host Verification
- HTTP probing with status codes
- Uses `httpx` tool if available, falls back to Go HTTP client
- Extracts: status code, content length, title, tech stack
- Per-category verification

**Key functions:**
- `ProbeFile(inputFile, outputFile)` — Verify URLs are live

### `internal/reporter/` — Report Generation
- Terminal output with ASCII boxes
- JSON report for machine reading
- Markdown report for documentation
- Statistics: dedup rate, live rate, URLs/second

**Key functions:**
- `Print(stats, duration)` — Terminal summary
- `WriteReports(stats)` — JSON and Markdown files

### `internal/runner/` — Pipeline Orchestration
- Coordinates all 5 pipeline stages
- Error handling and progress tracking
- Statistics collection and reporting

**Key functions:**
- `Run(options)` — Execute full pipeline
- `Options` struct — Configuration from CLI flags

### `internal/utils/` — Utility Functions
- File I/O: `ReadLines()`, `WriteLines()`, `WriteJSON()`
- Tool detection: `ToolExists()`, `ExecutableExists()`
- Number formatting: `FormatN()`, `Reduction()`
- Directory management: `EnsureDir()`

## Data Structures

### `collector.Config`
```go
type Config struct {
    Threads         int
    Depth           int
    Subs            bool
    RunAll          bool
    // Individual tool flags...
}
```

### `runner.Options`
```go
type Options struct {
    Targets        []string
    OutputDir      string
    Threads        int
    Depth          int
    SkipAlive      bool
    SkipCollect    bool
    // Tool selection flags...
}
```

### `reporter.Stats`
```go
type Stats struct {
    Targets      []string
    TotalRaw     int
    AfterNorm    int
    Groups       map[string]int      // URLs per category
    AliveGroups  map[string]int      // Live URLs per category
    UniqueParams int
    OutputDir    string
    DurationSec  float64
}
```

## File Organization

### Output Directory Structure
```
urlshine_20240506_143022/
├── raw/                       # Original tool outputs
│   ├── gau.txt
│   ├── katana.txt
│   ├── gospider.txt
│   └── ...
├── merged_raw.txt             # All URLs combined (stage 2)
├── normalized_urls.txt        # Deduplicated URLs (stage 3)
├── api_endpoints.txt          # Categorized: API (stage 4)
├── auth_admin_urls.txt        # Categorized: Auth
├── parameters_urls.txt        # Categorized: Parameters
├── js_config_urls.txt         # Categorized: JS & Config
├── directories_urls.txt       # Categorized: Directories
├── alive_api_endpoints.txt    # Verified live: API (stage 5)
├── alive_auth_admin_urls.txt  # Verified live: Auth
├── alive_parameters_urls.txt  # Verified live: Parameters
├── alive_js_config_urls.txt   # Verified live: JS & Config
├── alive_directories_urls.txt # Verified live: Directories
├── urlshine_report.json       # Machine-readable report
└── urlshine_report.md         # Human-readable report
```

## Key Algorithms

### URL Deduplication
```go
// Hash-based deduplication
seen := make(map[string]bool)
for _, url := range urls {
    if !seen[url] {
        seen[url] = true
        dedupedURLs = append(dedupedURLs, url)
    }
}
```

### Categorization (Regex Patterns)
- **API**: Contains `/api`, `/graphql`, `/v1`, `/swagger`, etc.
- **Auth**: Contains `login`, `admin`, `dashboard`, `2fa`, `sso`, etc.
- **Parameters**: Contains `?` with query parameters
- **JS/Config**: File extension matching (`.js`, `.json`, `.yaml`, `.env`)
- **Directories**: Path-based patterns

### Thread Pooling
Uses semaphore-based worker pool for concurrent tool execution:
```go
sem := make(chan struct{}, threads)
for _, tool := range tools {
    go func(t Tool) {
        sem <- struct{}{}        // Acquire
        defer func() { <-sem }() // Release
        t.Execute()
    }(tool)
}
```

## Performance Considerations

### Defaults
- **Threads**: 50 (configurable 1-500)
- **Depth**: 3 layers (for active crawlers)
- **Timeout**: 10 seconds per tool
- **Memory**: ~500MB for 1M URLs

### Optimization Tips
1. **Increase threads** for large URL sets (150-200)
2. **Reduce depth** to save time on large scans (2-3)
3. **Use SSD** for faster I/O
4. **Skip alive check** for speed: `--no-alive`
5. **Filter tools** — don't run unnecessary tools

## Testing Guidelines

### Unit Testing
```bash
go test ./...                    # Run all tests
go test -v ./internal/merger/    # Verbose output
go test -cover ./...             # Coverage report
```

### Integration Testing
```bash
./urlshine -a example.com
./urlshine -a -c example.com
./urlshine -f targets.txt -a -c
./urlshine doctor
```

### Performance Testing
```bash
time ./urlshine -a -c massive-list.txt
```

## Common Development Tasks

### Adding a New Collection Tool
1. Create `internal/collector/newtool.go`
2. Implement tool execution function
3. Add to `collector.Config` struct
4. Add CLI flag in `cmd/root.go`
5. Test: `./urlshine --newtool domain.com`

### Adding a New URL Category
1. Define pattern in `internal/splitter/splitter.go`
2. Add to `groupOrder` in `internal/reporter/reporter.go`
3. Test categorization output
4. Update README with new category

### Improving URL Normalization
1. Edit patterns in `internal/normalizer/normalizer.go`
2. Test with various URLs
3. Benchmark performance impact
4. Document changes

### Enhancing Reports
1. Edit `internal/reporter/reporter.go`
2. Add new fields to `Stats` struct
3. Update output formatting
4. Test with sample output

## Debugging

### Enable Verbose Logging
```bash
./urlshine -a -v google.com
```

### Check Tool Installation
```bash
./urlshine doctor
```

### Test Individual Tools
```bash
gau google.com
katana -u google.com
gospider -s google.com
```

### Profile Performance
```bash
go tool pprof ./urlshine
```

## Dependencies

### Go Modules
- **cobra**: CLI framework
- **fatih/color**: Terminal colors
- **schollz/progressbar**: Progress bars
- **golang.org/x/net**: Network utilities

### External Tools
- GAU, Katana, GoSpider, Waymore, Waybackurls, Hakrawler, xnLinkFinder, Gobuster, Dirb
- httpx (optional but recommended)

## Release Checklist

- [ ] Update version in `banner.go`
- [ ] Update CHANGELOG
- [ ] Run tests: `go test ./...`
- [ ] Build binaries: `make build-all`
- [ ] Update README if needed
- [ ] Tag release: `git tag v2.x.x`
- [ ] Push to repository
- [ ] Create GitHub release
- [ ] Update website/docs

## Contributing Code

1. **Fork** the repository
2. **Branch** for your feature: `git checkout -b feature/my-feature`
3. **Code** with proper documentation
4. **Test** thoroughly
5. **Commit** with clear messages
6. **Push** to your fork
7. **Pull Request** with description

See CONTRIBUTING.md for more details.

---

For questions, open an issue or start a discussion on GitHub.
