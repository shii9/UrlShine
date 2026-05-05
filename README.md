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

### 🚀 Automated Installation (Recommended)

**Linux/macOS:**
```bash
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
bash install.sh
```

**Windows:**
```bash
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
install.bat
```

> The installer scripts automatically install all required tools and dependencies!

### Manual Installation

**Step-by-step Setup (Linux/macOS):**
```bash
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
go build -o urlshine .
sudo mv urlshine /usr/local/bin/
chmod +x /usr/local/bin/urlshine
```

**Windows:**
```bash
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
go build -o urlshine.exe .
# Add directory to PATH or move urlshine.exe to a directory in PATH
```

### 📋 Verify Installation

**Check your installation:**
```bash
urlshine -h
```

**Check tool dependencies:**
```bash
urlshine doctor
```

This shows which tools are installed and provides installation commands for missing ones.

---

### 📦 All Tools Are Installed Automatically

The installer scripts (`install.sh` and `install.bat`) automatically install all dependencies:

| Tool | Type | Purpose |
|------|------|---------|
| GAU | Go | Archive & passive source aggregation (100+ threads) |
| Katana | Go | Active JS crawler with parameter extraction |
| GoSpider | Go | HTML & JS crawler (sitemaps, robots.txt) |
| Waymore | Python | Advanced wayback machine scraper |
| Waybackurls | Go | Wayback machine URL extraction |
| Hakrawler | Go | HTML content crawler with custom headers |
| xnLinkFinder | Python | JavaScript link and config extraction |
| Gobuster | Go | Directory discovery (50 threads, quiet mode) |
| Dirb | System | Directory brute-force enumeration |
| httpx | Go | Live host verification (optional but recommended) |

**If any tool fails to install:**
```bash
# Check what's installed
urlshine doctor

# Install missing tools manually
# See the output of 'urlshine doctor' for installation commands
```

---

## 🚀 Quick Start

### 📌 Flexible Flag Format

All flag formats work perfectly (use what you prefer):
```bash
# Format 1: Double-dash long flags
urlshine --all --complete google.com

# Format 2: Single-dash short letters
urlshine -a -c google.com

# Format 3: Single-dash long flags
urlshine -all -complete google.com

# Format 4: Mix and match
urlshine --gau -katana -c google.com
urlshine -a --complete google.com
```

### 📋 Common Use Cases

**Collection only:**
```bash
urlshine -a google.com           # All tools
urlshine -g -k google.com        # GAU + Katana
urlshine --gau --katana google.com
```

**Full processing pipeline:**
```bash
urlshine -a -c google.com                 # All tools + complete
urlshine -g -k -c google.com              # Specific tools + complete
urlshine -f targets.txt -a -c -t 150 -d 3 # File input with options
```

**Advanced options:**
```bash
urlshine -a -c --no-alive google.com             # Skip alive check
urlshine -a -c -t 200 -d 5 -o ./results google.com # Custom threads & depth
```

---

## ⚙️ Commands

### urlshine doctor
Check and verify all tool dependencies:
```bash
urlshine doctor
```
Shows which tools are installed and provides installation commands for missing ones.

---

## 💻 CLI Reference

### 🎯 Main Flags

| Short | Flag | Type | Default | Description |
|-------|------|------|---------|-------------|
| `-a` | `--all` | bool | false | Use all 9 collection tools |
| `-c` | `--complete` | bool | false | Complete pipeline: merge, normalize, categorize, alive-check |
| `-f` | `--file` | string | - | Input file with targets (one per line) |
| `-o` | `--output` | string | urlshine_TIMESTAMP | Output directory |
| `-t` | `--threads` | int | 50 | Parallel threads (recommended: 50-200) |
| `-d` | `--depth` | int | 5 | Crawl depth for active tools |
| `-v` | `--verbose` | bool | false | Debug/verbose logging |
| `-s` | `--subs` | bool | true | Include subdomains when supported |
| - | `--no-alive` | bool | false | Skip live host verification |
| - | `--skip-collect` | bool | false | Skip collection, reprocess existing data |

### 🔧 Collection Tools (choose one or more)

| Short | Flag | Tool | Purpose |
|-------|------|------|---------|
| `-g` | `--gau` | GetAllUrls | Archive & passive sources (100+ threads) |
| `-k` | `--katana` | Katana | Active JS crawler with parameter extraction |
| `-w` | `--gospider` | GoSpider | HTML & JS crawler (sitemaps, robots.txt) |
| `-m` | `--waymore` | Waymore | Advanced wayback machine scraper |
| `-b` | `--waybackurls` | Waybackurls | Wayback machine URL extraction |
| `-r` | `--hakrawler` | Hakrawler | HTML content crawler with custom headers |
| `-x` | `--xnlinkfinder` | xnLinkFinder | JavaScript link and config extraction |
| `-u` | `--gobuster` | Gobuster | Directory discovery (50 threads, quiet mode) |
| `-i` | `--dirb` | Dirb | Directory brute-force (non-recursive) |

### 📚 Mode Overview

**Collection Only (default):**
- Flag: None (or just `-g -k`, etc.)
- Output: Per-tool files in `{domain}_url/`
- Speed: Fast collection without post-processing
- Use: Quick enumeration, saving results

**Complete Pipeline:**
- Flag: `-complete` or `-c`
- Steps: Collection → Merge → Normalize → Categorize → Alive-check → Report
- Output: Merged, categorized, and report files
- Use: Full reconnaissance with analysis

### 🎯 Common Commands

**Quick enumeration (collection only):**
```bash
urlshine -a google.com                    # All 9 tools
urlshine -g -k -w google.com              # Specific tools
```

**Complete analysis (full pipeline):**
```bash
urlshine -a -c google.com                 # All tools + processing
urlshine -f targets.txt -a -c -t 150      # File input + options
```

**Advanced usage:**
```bash
urlshine -a -c --no-alive google.com      # Skip alive verification
urlshine -a -c -t 200 -d 5 google.com     # Custom threads & depth
urlshine --skip-collect -c google.com     # Reprocess existing data
```

### 📂 Output Examples

**Collection only (9 files):**
```bash
urlshine -a google.com
# Output: google_com_url/
#   ├── gau.txt
#   ├── katana.txt
#   ├── gospider.txt
#   ├── waymore.txt
#   ├── waybackurls.txt
#   ├── hakrawler.txt
#   ├── xnlinkfinder.txt
#   ├── gobuster.txt
#   └── dirb.txt
```

**Complete pipeline (categorized + reports):**
```bash
urlshine -a -c google.com
# Output: google_com_url/
#   ├── merged_urls.txt
#   ├── normalized_urls.txt
#   ├── api_endpoints.txt
#   ├── auth_admin_urls.txt
#   ├── parameters.txt
#   ├── js_config.txt
#   ├── directories.txt
#   ├── alive_urls.txt
#   ├── report.json
#   └── report.html
```

---

## ✨ Features
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
