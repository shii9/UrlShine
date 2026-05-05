<div align="center">

  <h1>🔗 URLShine v2.0.0</h1>
  <p><b>Advanced URL Enumeration & Attack Surface Mapper</b></p>
  <p>Professional reconnaissance engine for bug bounty hunters and red teams. Automates aggressive URL extraction from 9 concurrent tools, with deduplication, categorization, and live host verification.</p>

  <p>
    <a href="#-installation">Installation</a> •
    <a href="#-quick-start">Quick Start</a> •
    <a href="#-cli-reference">CLI Reference</a> •
    <a href="#-features">Features</a> •
    <a href="#-license">License</a>
  </p>
</div>

---

## 📥 Installation

**Requirements:**
- Go 1.21+ ([Install Go](https://golang.org/dl/))
- Linux/macOS/Windows

**Fast Installation:**
```bash
git clone https://github.com/shii9/UrlShine.git && cd UrlShine && go build -o urlshine . && sudo mv urlshine /usr/local/bin/ && chmod +x /usr/local/bin/urlshine
```

**Step-by-step Setup (Linux/macOS):**
```bash
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
go build -o urlshine .
sudo mv urlshine /usr/local/bin/
chmod +x /usr/local/bin/urlshine
```

**Windows Setup:**
```bash
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
go build -o urlshine.exe .
# Then add the directory to your PATH or move urlshine.exe to a directory in PATH
```

**Verify Installation:**
```bash
urlshine -h
```

---

### Optional Collection Tools

All tools are optional. URLShine gracefully skips tools that aren't installed.

| Tool | Installation Command |
|------|----------------------|
| [GAU](https://github.com/lc/gau) | `go install github.com/lc/gau/v2/cmd/gau@latest` |
| [Katana](https://github.com/projectdiscovery/katana) | `go install github.com/projectdiscovery/katana/cmd/katana@latest` |
| [GoSpider](https://github.com/jaeles-project/gospider) | `go install github.com/jaeles-project/gospider@latest` |
| [Waybackurls](https://github.com/tomnomnom/waybackurls) | `go install github.com/tomnomnom/waybackurls@latest` |
| [Hakrawler](https://github.com/hakluke/hakrawler) | `go install github.com/hakluke/hakrawler@latest` |
| [Gobuster](https://github.com/OJ/gobuster) | `go install github.com/OJ/gobuster/v3@latest` |
| [Waymore](https://github.com/xnl-h4ck3r/waymore) | `pip3 install waymore` |
| [xnLinkFinder](https://github.com/xnl-h4ck3r/xnLinkFinder) | `pip3 install xnLinkFinder` |
| [Dirb](https://sourceforge.net/projects/dirb/) | `apt-get install dirb` or `brew install dirb` |
| [httpx](https://github.com/projectdiscovery/httpx) | `go install github.com/projectdiscovery/httpx/cmd/httpx@latest` |

---

## 🚀 Quick Start

### ✅ Flexible Flag Format
All of these formats work perfectly (choose whichever you prefer):

```bash
# Format 1: Long flags with double dash
urlshine --all --complete google.com

# Format 2: Short flags with single dash (letters)
urlshine -a -c google.com

# Format 3: Long flags with single dash (NEW!)
urlshine -all -complete google.com

# Format 4: Mix and match
urlshine --gau -katana -c google.com
urlshine -a --complete google.com
```

### Examples

**Collect URLs only (any format works):**
```bash
urlshine --gau --katana google.com
urlshine -gau -katana google.com
urlshine -g -k google.com
urlshine --all google.com
urlshine -all google.com
urlshine -a google.com
```

**Collect + Complete processing (any format works):**
```bash
# All of these are equivalent:
urlshine --all --complete google.com
urlshine -all -complete google.com
urlshine -a -c google.com
urlshine --gau --katana --complete google.com
urlshine -gau -katana -complete google.com
urlshine -g -k -c google.com
urlshine -a -c -t 100 -d 5 google.com
```

**Process multiple domains from file:**
```bash
urlshine -f targets.txt -a -c -t 150 -d 3 -o ./results
```

**Skip live verification:**
```bash
urlshine -a -c --no-alive google.com
```

---

## 💻 CLI Reference

### Main Flags

| Shorthand | Flag | Type | Default | Description |
|-----------|------|------|---------|-------------|
| `-a` | `--all` | boolean | false | Use all 9 collection tools |
| `-c` | `--complete` | boolean | false | Run complete processing pipeline (merge, normalize, categorize, alive check) |
| `-f` | `--file` | string | - | Input file with targets (one per line) |
| `-o` | `--output` | string | urlshine_TIMESTAMP | Output directory |
| `-t` | `--threads` | integer | 50 | Parallel threads (recommended: 50-200) |
| `-d` | `--depth` | integer | 5 | Crawl depth for active tools |
| `-v` | `--verbose` | boolean | false | Debug/verbose logging |
| `-s` | `--subs` | boolean | true | Include subdomains when supported |
| | `--no-alive` | boolean | false | Skip live host verification |
| | `--skip-collect` | boolean | false | Skip collection, reprocess existing data |

### Collection Tools (choose one or more)

| Short | Flag | Tool | Purpose |
|-------|------|------|---------|
| `-g` | `--gau` | GetAllUrls | Archive & passive sources |
| `-k` | `--katana` | Katana | Active JS crawler |
| `-w` | `--gospider` | GoSpider | HTML & JS crawler |
| `-m` | `--waymore` | Waymore | Advanced wayback scraper |
| `-b` | `--waybackurls` | Waybackurls | Wayback machine scraper |
| `-r` | `--hakrawler` | Hakrawler | HTML content crawler |
| `-x` | `--xnlinkfinder` | xnLinkFinder | JS endpoint extractor |
| `-u` | `--gobuster` | Gobuster | Directory brute-force |
| `-i` | `--dirb` | Dirb | Directory enumeration |

### Understanding `-all` vs `-complete`

**What does `-all` do?**
```
-all means: Use all 9 collection tools
```
- Runs: GAU, Katana, GoSpider, Waymore, Waybackurls, Hakrawler, xnLinkFinder, Gobuster, Dirb
- Parallel execution with 10 concurrent tool executors (50 threads per tool)
- Output: Per-tool files (gau.txt, katana.txt, etc.)

**What does `-complete` do?**
```
-complete means: Complete all processing steps
```
1. **Merging** — Deduplicates all results
2. **Normalization** — Cleans URLs
3. **Categorization** — Splits into 5 attack groups (API, Auth, Params, JS, Directories)
4. **Alive Checking** — Verifies live hosts (unless -no-alive used)
5. **Reporting** — Generates JSON & HTML reports

Output: Merged, normalized, categorized files + reports

### Example Commands

**Simple collection (fast):**
```bash
# Specific tools, collection only (long form)
urlshine --gau --katana google.com

# Specific tools, collection only (short form)
urlshine -g -k google.com

# All tools, collection only (long form)
urlshine --all google.com

# All tools, collection only (short form)
urlshine -a google.com
```

**Collection + complete processing:**
```bash
# Long form with specific tools
urlshine --gau --katana --complete google.com

# Short form with specific tools  
urlshine -g -k -c google.com

# Long form with all tools
urlshine --all --complete google.com

# Short form with all tools
urlshine -a -c google.com
```

**Batch processing:**
```bash
# Multiple targets from file (short form)
urlshine -f targets.txt -a -c -t 100 -d 5 -o ./results

# Skip alive check for speed (short form)
urlshine -f targets.txt -a -c --no-alive -t 150 -d 3

# Verbose mode for debugging (short form)
urlshine -f targets.txt -a -c -v
```

**Reprocessing:**
```bash
# Reprocess existing results without recollecting
urlshine --skip-collect --no-alive google.com

# Apply different settings to cached results
urlshine --skip-collect -t 100 google.com
```

---

## ✨ Features

- ✅ **9-Engine URL Collection**: GAU, Katana, GoSpider, Waymore, Waybackurls, Hakrawler, xnLinkFinder, Gobuster, Dirb
- ✅ **Parallel Execution**: 50 default threads with 10 concurrent tool executors (configurable to 200+)
- ✅ **Deep Crawling**: Default depth of 5 layers for maximum URL discovery
- ✅ **Massive Batch Processing**: Handles 10,000+ subdomain lists efficiently
- ✅ **Advanced Categorization**: Splits URLs into 5 attack groups:
  - **API Endpoints** (grep patterns + httpx verification)
  - **Auth/Admin Pages** (grep patterns + gobuster discovery)
  - **Parameters** (grep + arjun discovery)
  - **JavaScript & Config** (grep + linkfinder extraction)
  - **Directories** (path analysis + brute-force tools)
- ✅ **Professional Deduplication**: Hash-based deduplication handles 1M+ URLs
- ✅ **Live Verification**: httpx integration with Go fallback prober
- ✅ **Professional Reports**: JSON and Markdown summaries with statistics

---

## 📊 Output Structure

**Collection only** (without `-complete`):
```
google_url/
├── gau.txt
├── katana.txt
├── gospider.txt
├── waymore.txt
├── waybackurls.txt
├── hakrawler.txt
├── xnlinkfinder.txt
├── gobuster.txt
└── dirb.txt
```

**Complete processing** (with `-complete`):
```
google_url/
├── merged_urls.txt (all tools combined)
├── normalized_urls.txt (cleaned & deduplicated)
├── api_endpoints.txt (API paths)
├── auth_admin_urls.txt (authentication pages)
├── parameters.txt (URLs with parameters)
├── js_config.txt (JavaScript & config files)
├── directories.txt (directory paths)
├── alive_urls.txt (verified live hosts)
├── report.json
└── report.html
```

---

## 🎯 Performance Benchmarks

| Scope | Time | Settings |
|-------|------|----------|
| 1-100 domains | 2-5 min | `-all -complete -t 50 -d 5` |
| 100-1000 domains | 10-30 min | `-all -complete -t 100 -d 3` |
| 1000+ domains | 30-90 min | `-all -complete -t 150 -d 2` |

With `-no-alive` flag: 30-50% faster

---

## 💡 Use Cases

- **Bug Bounty** — Feed a scope file, get categorized URLs ready for sqlmap, nuclei, or manual testing
- **Attack Surface Mapping** — Discover hidden APIs and legacy infrastructure
- **Security Research** — Analyze web applications and their attack surfaces
- **Web Scraping** — Generate specific URL lists for mass requests

---

## ⚠️ Legal Notice

**Authorized use only.** Only scan domains you own or have explicit written permission to test (e.g., Bug Bounty programs). Unauthorized network scanning may violate laws including the Computer Fraud and Abuse Act (CFAA) and similar legislation in other jurisdictions.

---

## 🤝 Contributing

Contributions welcome! Please:
- Report bugs via GitHub Issues
- Suggest features
- Submit pull requests

---

## 📄 License

MIT License - See LICENSE for details.

<div align="center">
  <i>Professional. Focused. Efficient.</i>
</div>

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
