<div align="center">

  <h1>🔗 URLShine v2.0.0</h1>
  <p><b>Advanced URL Enumeration & Attack Surface Mapper</b></p>
  <p>Professional reconnaissance engine for bug bounty hunters and red teams. Automates aggressive URL extraction from 9 concurrent tools, with deduplication, categorization, and live host verification.</p>

  <p>
    <a href="#-installation">Installation</a> •
    <a href="#-quick-start">Quick Start</a> •
    <a href="#-cli-reference">CLI Reference</a> •
    <a href="#-features">Features</a> •
    <a href="#-use-cases">Use Cases</a> •
    <a href="#-legal-notice">Legal</a>
  </p>
</div>

---

## 📥 Installation

**Requirements:** Go 1.21+ • Linux/macOS/Windows

### 🚀 Automated Installation (Recommended)

**Linux/macOS:**
```bash
git clone https://github.com/shii9/UrlShine.git && cd UrlShine
bash install.sh
```

**Windows:**
```bash
git clone https://github.com/shii9/UrlShine.git && cd UrlShine
install.bat
```

The installer automatically installs all 10 tools and dependencies.

### Manual Installation

**Linux/macOS:**
```bash
git clone https://github.com/shii9/UrlShine.git && cd UrlShine
go build -o urlshine . && sudo mv urlshine /usr/local/bin/
```

**Windows:**
```bash
git clone https://github.com/shii9/UrlShine.git && cd UrlShine
go build -o urlshine.exe .
# Add directory to PATH or move urlshine.exe to a PATH directory
```

### ✅ Verify Installation

```bash
urlshine -h              # Show help menu
urlshine doctor          # Check tool dependencies
```

The `doctor` command shows which tools are installed and provides installation commands for missing ones.

---

## 📦 Collection Tools

| Tool | Type | Purpose |
|------|------|---------|
| **GAU** | Go | Archive & passive source aggregation (100+ threads) |
| **Katana** | Go | Active JS crawler with parameter extraction |
| **GoSpider** | Go | HTML & JS crawler (sitemaps, robots.txt) |
| **Waymore** | Python | Advanced wayback machine scraper |
| **Waybackurls** | Go | Wayback machine URL extraction |
| **Hakrawler** | Go | HTML content crawler with custom headers |
| **xnLinkFinder** | Python | JavaScript link and config extraction |
| **Gobuster** | Go | Directory discovery (50 threads, quiet mode) |
| **Dirb** | System | Directory brute-force enumeration |
| **httpx** | Go | Live host verification (optional but recommended) |

---

## 🚀 Quick Start

### 📌 Basic Commands

```bash
# Single target - all tools
urlshine -a google.com

# Single target - specific tools
urlshine -g -k -w google.com

# Multiple targets
urlshine -a -f targets.txt

# Full pipeline (collection + processing)
urlshine -a -c google.com

# Custom settings
urlshine -a -c -t 150 -d 3 google.com
```

### 🎛️ Flexible Flag Format

All formats work identically (use what you prefer):
```bash
urlshine --all --complete google.com     # Long flags
urlshine -a -c google.com                # Short flags  
urlshine -all -complete google.com       # Single-dash long
urlshine --gau -katana -c google.com     # Mix and match
```

---

## 💻 CLI Reference

### 🎯 Main Flags

| Short | Flag | Type | Description |
|-------|------|------|-------------|
| `-a` | `--all` | bool | Use all 9 collection tools |
| `-c` | `--complete` | bool | Full pipeline: collect → merge → normalize → categorize → alive-check |
| `-f` | `--file` | string | Input file with targets (one per line) |
| `-o` | `--output` | string | Output directory (default: urlshine_TIMESTAMP) |
| `-t` | `--threads` | int | Parallel threads (default: 50, recommended: 50-200) |
| `-d` | `--depth` | int | Crawl depth for active tools (default: 5) |
| `-v` | `--verbose` | bool | Debug/verbose logging |
| `-s` | `--subs` | bool | Include subdomains (default: true) |
| `--no-alive` | - | bool | Skip live host verification |
| `--skip-collect` | - | bool | Reprocess existing data without recollecting |

### 🔧 Collection Tool Flags (choose one or more)

| Short | Flag | Tool | Purpose |
|-------|------|------|---------|
| `-g` | `--gau` | GetAllUrls | Archive & passive sources |
| `-k` | `--katana` | Katana | Active JS crawler |
| `-w` | `--gospider` | GoSpider | HTML & JS crawler |
| `-m` | `--waymore` | Waymore | Wayback machine scraper |
| `-b` | `--waybackurls` | Waybackurls | Wayback extraction |
| `-r` | `--hakrawler` | Hakrawler | HTML crawler |
| `-x` | `--xnlinkfinder` | xnLinkFinder | JS link finder |
| `-u` | `--gobuster` | Gobuster | Directory discovery |
| `-i` | `--dirb` | Dirb | Directory brute-force |

### 📂 Output Structure

**Collection only** (default behavior):
```
google_com_url/
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

**Complete pipeline** (with `-complete` flag):
```
google_com_url/
├── merged_urls.txt           # All tools combined
├── normalized_urls.txt       # Cleaned & deduplicated
├── api_endpoints.txt         # API paths
├── auth_admin_urls.txt       # Authentication pages
├── parameters.txt            # URLs with parameters
├── js_config.txt             # JavaScript & config files
├── directories.txt           # Directory paths
├── alive_urls.txt            # Verified live hosts
├── report.json               # JSON summary
└── report.html               # HTML report
```

---

## ✨ Features

- ✅ **9-Tool URL Collection** — GAU, Katana, GoSpider, Waymore, Waybackurls, Hakrawler, xnLinkFinder, Gobuster, Dirb
- ✅ **Parallel Execution** — 50 default threads with 10 concurrent tool executors (configurable to 200+)
- ✅ **Deep Crawling** — Default depth of 5 layers for maximum URL discovery
- ✅ **Batch Processing** — Handles 10,000+ subdomain lists efficiently
- ✅ **Smart Categorization** — Splits URLs into 5 attack groups:
  - API Endpoints (grep patterns + httpx verification)
  - Auth/Admin Pages (grep patterns + gobuster discovery)
  - Parameters (grep + arjun discovery)
  - JavaScript & Config (grep + linkfinder extraction)
  - Directories (path analysis + brute-force tools)
- ✅ **Deduplication** — Hash-based deduplication handles 1M+ URLs
- ✅ **Live Verification** — httpx integration with Go fallback prober
- ✅ **Professional Reports** — JSON and Markdown summaries with statistics

---

## 💡 Use Cases

| Use Case | Command |
|----------|---------|
| Bug Bounty — quick enumeration | `urlshine -a target.com` |
| Bug Bounty — full analysis | `urlshine -a -c target.com` |
| Multiple targets | `urlshine -a -f scope.txt -c` |
| Skip alive check | `urlshine -a -c --no-alive target.com` |
| Reprocess existing data | `urlshine --skip-collect -c target.com` |
| Custom threads & depth | `urlshine -a -c -t 200 -d 3 target.com` |
| Attack surface mapping | `urlshine -a -c -t 100 -f domains.txt` |
| Feed to other tools | `urlshine -a target.com` → use with nuclei/sqlmap |

---

## ⚠️ Legal Notice

**Authorized use only.** Only scan domains you own or have explicit written permission to test (e.g., Bug Bounty programs). Unauthorized network scanning may violate laws including the Computer Fraud and Abuse Act (CFAA) and similar legislation in other jurisdictions.

---

## 📄 License

MIT License - See LICENSE for details.

<div align="center">
  <i>Professional. Focused. Efficient.</i>
</div>
