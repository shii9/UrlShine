<div align="center">

  <img src="https://raw.githubusercontent.com/shii9/NetRange/main/assets/banner.png" alt="URLShine Banner" width="800" onerror="this.style.display='none'"/>

  <h1>URLShine v2.0.0</h1>
  <p><b>Advanced URL Enumeration & Attack Surface Mapper</b></p>
  <p>Professional reconnaissance engine for bug bounty hunters and red teams. Automates aggressive URL extraction from 9 concurrent tools, normalization, categorization, and live host verification. Optimized to handle massive subdomain lists.</p>

  <p>
    <a href="#-installation">Installation</a> •
    <a href="#-quick-start">Usage</a> •
    <a href="#-features">Features</a> •
    <a href="#-contributing">Contributing</a> •
    <a href="#-license">License</a>
  </p>
</div>

---

Professional tool for expanding your attack surface by chaining the industry's best enumeration tools into a single, high-performance, deduplicated, and categorized pipeline with aggressive parallel execution.
Optimized for large-scale reconnaissance on massive subdomain lists. Essential for web application security testing, bug bounties, and infrastructure analysis.

## 📥 Installation

**One-liner Installation:**
```bash
git clone https://github.com/shii9/UrlShine.git && cd UrlShine && go build -o urlshine . && sudo mv urlshine /usr/local/bin/
```

**Step-by-step Setup:**
```bash
# Clone repository
git clone https://github.com/shii9/UrlShine.git
cd UrlShine

# Build binary
go build -o urlshine .

# Linux/macOS: Install globally (one-time setup)
sudo mv urlshine /usr/local/bin/
chmod +x /usr/local/bin/urlshine

# Verify installation
urlshine --help
```

**Windows Setup:**
```powershell
# Clone repository
git clone https://github.com/shii9/UrlShine.git
cd UrlShine

# Build binary
go build -o urlshine.exe .

# Option 1: Use directly from directory
.\urlshine.exe --help

# Option 2: Add to PATH (Recommended)
# 1. Copy urlshine.exe to a folder (e.g., C:\Tools\)
# 2. Add C:\Tools\ to Windows PATH environment variable
# 3. Then use: urlshine --help
```

**Requirements:**
- Go 1.21+ ([Install Go](https://golang.org/dl/))
- Recommended: Add to PATH for global access (no `./` prefix needed)

**Optional URL Collection Tools** (install for maximum effectiveness):
- [GAU](https://github.com/lc/gau) — GetAllUrls
- [Katana](https://github.com/projectdiscovery/katana) — ProjectDiscovery Katana crawler
- [GoSpider](https://github.com/jaeles-project/gospider) — Spider & crawler
- [Waymore](https://github.com/xnl-h4ck3r/waymore) — Advanced wayback machine scraper
- [Waybackurls](https://github.com/tomnomnom/waybackurls) — Tomnomnom wayback URLs
- [Hakrawler](https://github.com/hakluke/hakrawler) — Hakluke HTML crawler
- [xnLinkFinder](https://github.com/xnl-h4ck3r/xnLinkFinder) — xnl link finder
- [Gobuster](https://github.com/OJ/gobuster) — Directory brute-forcing tool
- [Dirb](https://sourceforge.net/projects/dirb/) — Directory brute-force tool
- [httpx](https://github.com/projectdiscovery/httpx) — HTTP probe tool (optional, for alive verification)

## 🚀 Quick Start

```bash
# Enumerate single domain (auto-runs all 9 tools + gobuster + dirb)
urlshine google.com

# Specific tools only
urlshine -gau -katana google.com

# All tools with aggressive settings
urlshine -all -t 100 -d 5 google.com

# Process file with multiple targets
urlshine -f targets.txt -all -t 150 -d 3

# Combine specific tools with -all for additional discovery
urlshine -gau -katana -all -t 100 -v google.com

# Fast mode (skip alive checking)
urlshine -all -no-alive google.com

# Professional output with custom directory
urlshine -all -t 100 -d 5 -o ./reports google.com

# Verbose mode for debugging
urlshine -all -v google.com
```

## ✨ Features

- ✅ **Aggressive Parallel Execution**: 50 default threads with 10 concurrent tool executors (configurable up to 200+)
- ✅ **Deep Crawling**: Default depth of 5 layers for maximum URL discovery
- ✅ **9-Engine URL Collection**: Harnesses `gau`, `gospider`, `katana`, `waymore`, `waybackurls`, `hakrawler`, `xnLinkFinder`, `gobuster`, and `dirb` for comprehensive URL discovery and directory brute-forcing.
- ✅ **Large Batch Processing**: Efficiently processes massive subdomain files with professional parallelism
- ✅ **Granular Tool Selection**: Explicitly toggle which tools to run via CLI flags.
- ✅ **Regex-Based Categorization**: Intelligently splits URLs into high-value targets: `API`, `Auth`, `Params`, `JS`, and `Directories`.
- ✅ **Advanced Tool-Based Extraction**: Uses specialized tools and commands:
  - **API Detection**: grep patterns + httpx live verification
  - **Auth/Admin**: grep + gobuster/feroxbuster path discovery
  - **Parameters**: grep + arjun/paramspider discovery
  - **JS/Config**: grep + linkfinder/subjs extraction
  - **Directories**: Path analysis + gobuster/dirb/ffuf/feroxbuster
- ✅ **High-Performance Deduplication**: Fast hash-based deduplication handles 1M+ URLs
- ✅ **Live Verification**: Integrated `httpx` support with blazing fast, native Go fallback prober.
- ✅ **Professional Reporting**: Generates JSON and Markdown summaries with comprehensive statistics.

## 🛠️ The Pipeline

When using the `--all` flag, URLShine executes a comprehensive 5-step reconnaissance pipeline with aggressive settings:

1. **URL Collection**: Concurrently executes all 9 URL collection engines with 10 parallel tool executors (50 threads per tool)
2. **Merging**: Aggregates all raw tool outputs into a consolidated list
3. **Normalization**: Strips duplicate ports, sanitizes schemes, and removes static junk extensions (`.jpg`, `.css`, etc.)
4. **Advanced Extraction**: Uses specialized tools to identify:
   - **API Endpoints** (via grep patterns, httpx verification)
   - **Auth/Admin Pages** (via grep patterns, gobuster/feroxbuster)
   - **Parameters** (via grep, arjun, paramspider)
   - **JS/Config Files** (via grep, linkfinder, subjs)
   - **Directories** (via gobuster, dirb, ffuf, feroxbuster, dirsearch)
5. **Alive Checking**: Probes all URLs to filter dead links (httpx or native Go fallback)
6. **Reporting**: Generates professional Markdown and JSON reports

## 💻 CLI Options & Commands

### Flag Reference

**Target Input:**
```
  <domain>               Single domain (example: urlshine google.com)
  -f, --file FILE        File with targets, one per line
```

**URL Collection Tools** (run individually or combine):
```
  -gau                   GAU (GetAllUrls) - archive & passive sources
  -katana                Katana - active JS crawler
  -gospider              GoSpider - HTML & JS crawler
  -waymore               Waymore - advanced wayback scraper
  -waybackurls           Wayback URLs - wayback machine scraper
  -hakrawler             Hakrawler - HTML content crawler
  -xnlinkfinder          xnLinkFinder - JS endpoint extractor
  -gobuster              Gobuster - directory brute-force discovery
  -dirb                  Dirb - directory enumeration
  -all                   Run ALL 9 tools + gobuster + dirb automatically
```

**Configuration Options:**
```
  -t, --threads INT      Parallel threads (default: 50, recommended: 50-200)
  -d, --depth INT        Crawl depth for active tools (default: 5)
  -o, --output DIR       Output directory (default: urlshine_<timestamp>)
  -f, --file FILE        Input file with targets
```

**Processing & Control:**
```
  -no-alive              Skip live host verification (fast mode)
  -skip-collect          Reprocess existing data (skip collection)
  -v, --verbose          Debug/verbose logging output
  -h, --help             Display help menu
```

### Command Examples

**Basic Usage:**
```bash
# Single domain - auto-runs all 9 tools
urlshine google.com

# Multiple domains
urlshine google.com yahoo.com facebook.com

# Show help
urlshine --help
```

**Selective Tool Usage:**
```bash
# Specific tools only
urlshine -gau -katana google.com

# Add directory discovery
urlshine -katana -gobuster -dirb google.com

# Combine multiple sources
urlshine -gau -katana -gospider -t 50 google.com
```

**Professional Reconnaissance:**
```bash
# All tools with aggressive settings (includes gobuster + dirb by default)
urlshine -all -t 100 -d 5 google.com

# Large batch processing
urlshine -f targets.txt -all -t 150 -d 2 -o ./results

# Fast mode without alive checking
urlshine -all -no-alive -t 100 google.com

# Verbose output for debugging
urlshine -all -t 100 -v google.com
```

**Advanced Scenarios:**
```bash
# Combine -all with individual tools (tools won't duplicate)
urlshine -all -gau -t 100 google.com

# All tools + full verbosity + custom output
urlshine -all -t 100 -d 5 -o ./enterprise_scan -v google.com

# Reprocess existing results with different settings
urlshine -skip-collect -no-alive google.com

# Batch with deep crawling and custom threads
urlshine -f massive-targets.txt -all -t 150 -d 3 -o ./batch_results -v
```

### About the `-all` Flag

When you use `-all`, URLShine automatically:
- ✅ Runs all 9 URL collection tools (GAU, Katana, GoSpider, Waymore, Waybackurls, Hakrawler, xnLinkFinder)
- ✅ Includes Gobuster for aggressive directory discovery
- ✅ Includes Dirb for directory brute-forcing
- ✅ Performs full URL normalization and deduplication
- ✅ Advanced extraction for API endpoints, Auth pages, Parameters, JS/Config, Directories
- ✅ Optional: Live host verification (disable with `-no-alive`)

**Example with combinations:**
```bash
# All tools are already in -all, this runs the same as 'urlshine -all google.com'
urlshine -all google.com

# Add extra verbosity
urlshine -all -v google.com

# Increase threads and depth
urlshine -all -t 150 -d 5 google.com

# Fast mode with all tools
urlshine -all -no-alive -t 100 google.com
```

## ⚙️ Requirements

URLShine requires the following tools to be installed and available in your system's `$PATH` for maximum effectiveness:

*   [GAU](https://github.com/lc/gau)
*   [Katana](https://github.com/projectdiscovery/katana)
*   [Waymore](https://github.com/xnl-h4ck3r/waymore)
*   [GoSpider](https://github.com/jaeles-project/gospider)
*   [Waybackurls](https://github.com/tomnomnom/waybackurls)
*   [Hakrawler](https://github.com/hakluke/hakrawler)
*   [xnLinkFinder](https://github.com/xnl-h4ck3r/xnLinkFinder)
*   [Gobuster](https://github.com/OJ/gobuster)
*   [Dirb](https://sourceforge.net/projects/dirb/)
*   [httpx](https://github.com/projectdiscovery/httpx)

*(Note: URLShine gracefully skips tools that are not installed on your system without crashing).*

### Optional Extraction Tools

Install any of these to unlock advanced URL extraction and specialized group detection:

**API Extraction:**
- [httpx](https://github.com/projectdiscovery/httpx) — Live API endpoint verification

**Auth/Admin Detection:**
- [Feroxbuster](https://github.com/epi052/feroxbuster) — Recursive directory discovery

**Parameter Discovery:**
- [Arjun](https://github.com/s0md3v/Arjun) — HTTP parameter discovery
- [ParamSpider](https://github.com/devanshbatham/ParamSpider) — Parameter mining

**JavaScript/Config Extraction:**
- [LinkFinder](https://github.com/GerbenJavado/LinkFinder) — Extract endpoints from JS
- [SubJS](https://github.com/lc/subjs) — Fetch JS from URLs
- [Katana](https://github.com/projectdiscovery/katana) — Already included (JS extraction)

**Directory Discovery:**
- [Ffuf](https://github.com/ffuf/ffuf) — Fast fuzzing tool
- [Feroxbuster](https://github.com/epi052/feroxbuster) — Recursive discovery
- [Dirsearch](https://github.com/maurosoria/dirsearch) — Directory scanner

## 🔍 Advanced Extraction Workflow

When using `--all`, URLShine automatically extracts specialized groups using best-practice tools and techniques:

### 1️⃣ API Endpoints Extraction
```bash
# Detection Method:
# - Grep for patterns: /api/, graphql, /v1/, /v2/, rest, swagger, openapi
# - Verify live with: httpx -mc 200,201,204,301,302

# Output: api_endpoints.txt
# Contains: All API-like URLs, filtered for live responses
```

### 2️⃣ Auth/Admin Pages Extraction
```bash
# Detection Method:
# - Grep for patterns: admin, login, signin, signup, auth, panel, dashboard
# - Brute-force with: gobuster dir -x php,txt,js,html
# - Include: Feroxbuster for recursive discovery (if installed)

# Output: auth_admin.txt
# Contains: Authentication, admin panels, privileged endpoints
```

### 3️⃣ Parameters Extraction
```bash
# Detection Method:
# - Grep for URLs with query strings (?param=value)
# - Mine additional params with: arjun (if installed)
# - Fallback: paramspider support planned

# Output: parameters.txt
# Contains: All URLs with query parameters for fuzzing/testing
```

### 4️⃣ JavaScript & Config Extraction
```bash
# Detection Method:
# - Grep for: .js, .json, .xml, .env, .bak, config files
# - Extract from crawling results (Katana JS crawl)
# - Support: LinkFinder, SubJS (if installed)

# Output: js_config.txt
# Contains: JavaScript files, configuration files, API definitions
```

### 5️⃣ Directory Extraction
```bash
# Detection Method:
# - Parse URLs to identify directory-like paths
# - Brute-force with: gobuster, dirb, ffuf, feroxbuster (if installed)

# Output: directories.txt
# Contains: Directory paths for further enumeration
```

## ⚠️ Aggressive Collection Mode — Professional-Grade Reconnaissance

URLShine is engineered for **massive-scale URL enumeration** from large subdomain lists:

### Aggressive Features:
- **50 Default Threads** — Parallelizes URL collection across all 9 tools simultaneously
- **10 Concurrent Tool Executors** — Each tool runs in parallel against different targets
- **Depth 5 Crawling** — Deep active crawling (GoSpider, Katana, Hakrawler) for maximum coverage
- **Enhanced Tool Parameters**:
  - **GAU**: Uses 100+ threads with 4 source providers (Wayback, CommonCrawl, URLScan, OTX)
  - **Katana**: Includes JS crawling, parameter extraction, and query string capture
  - **GoSpider**: Fetches JS, sitemaps, robots.txt, and uses all crawl methods
  - **Waymore**: Mode 3 (maximum) with 25-50 concurrent requests
  - **Hakrawler**: Includes custom User-Agent headers for better coverage
  - **Gobuster**: Directory discovery with aggressive wordlist brute-forcing
  - **Dirb**: Directory brute-forcing tool for finding hidden directories and files

### Example: Processing Massive Subdomain Lists
```bash
# Professional reconnaissance on 10,000+ subdomains
urlshine -f big-bug-bounty-scope.txt --all -t 100 -d 5 -o ./professional_scan

# This will:
# - Execute 9 tools × N subdomains in aggressive parallel mode
# - Crawl up to 5 layers deep on each subdomain
# - Collect URLs from: Historical archives, JS files, API endpoints, config files
# - Deduplicate 1M+ URLs at high speed
# - Verify which URLs are alive with httpx
# - Generate professional JSON/Markdown reports
```

### Performance Characteristics:
- **100-500 domains**: ~5-15 minutes with default settings
- **1000+ domains**: Recommend `-t 150-200 -d 3` for optimal speed
- **Memory**: Handles 1M+ URL deduplication efficiently
- **Network**: Respectful rate limiting while maintaining aggressive throughput

## �💡 Use Cases

- **Bug Bounty** — Feed URLShine a scope file, go grab a coffee, and return to categorized files ready for `sqlmap`, `nuclei`, or manual testing.
- **Attack Surface Management** — Discover hidden API endpoints and legacy infrastructure left behind by developers.
- **Web Scraping** — Generate highly specific lists for mass web requests and API testing.

## ⚠️ Legal Notice

Authorized use only. Only scan domains you own or have explicit written permission to test (e.g., Bug Bounty programs). Unauthorized network scanning may violate laws including the Computer Fraud and Abuse Act (CFAA) in the United States and similar legislation in other jurisdictions.

## 🤝 Contributing

We welcome contributions! Feel free to:
- Report bugs
- Suggest features
- Submit pull requests

## 📄 License

MIT License - See LICENSE for details.

<div align="center">
  <i>Professional. Focused. Efficient.</i>
</div>
