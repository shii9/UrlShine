<div align="center">

# 🔗 URLShine

**Professional URL Enumeration & Attack Surface Mapper**

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-lightgrey.svg)]()

Advanced reconnaissance tool for security professionals. Automates URL extraction from 9 concurrent tools with intelligent deduplication, categorization, and live host verification.

[Installation](#-installation) • [Quick Start](#-quick-start) • [CLI Reference](#-cli-reference) • [Features](#-features) • [Output](#-output) • [FAQ](#-faq)

</div>

---

## 📋 Overview

URLShine orchestrates a sophisticated reconnaissance pipeline combining multiple URL enumeration tools into a unified workflow. It collects URLs from passive and active sources, deduplicates results, categorizes findings by attack vector, and verifies live hosts.

**Perfect for:**
- 🎯 Bug Bounty Hunters
- 🛡️ Red Team Engagements
- 🔍 Security Researchers
- 🏢 Penetration Testers

---

## ⚡ Features

### Collection Engines
| Tool | Type | Purpose |
|------|------|---------|
| **GAU** | Passive | Archive aggregation (100+ threads) |
| **Katana** | Active | JavaScript crawler with parameter extraction |
| **GoSpider** | Active | HTML & JS crawler (sitemaps, robots.txt) |
| **Waymore** | Passive | Advanced Wayback Machine scraper |
| **Waybackurls** | Passive | Wayback Machine URL extraction |
| **Hakrawler** | Active | HTML content crawler with headers |
| **xnLinkFinder** | Passive | JavaScript link & config extraction |
| **Gobuster** | Active | Directory discovery (optimized for threads) |
| **Dirb** | Active | Directory brute-force enumeration |

### Processing Pipeline
- ✅ **Concurrent Collection** — All tools run in parallel with thread pooling
- ✅ **Smart Deduplication** — Hash-based URL deduplication
- ✅ **URL Normalization** — Removes static assets, redundant ports, invalid URLs
- ✅ **Intelligent Categorization** — Splits URLs into 5 attack vectors
- ✅ **Live Verification** — HTTP probing with status codes and metadata
- ✅ **Professional Reports** — JSON, Markdown, and terminal summaries

### Attack Vector Categories
1. **API Endpoints** — `/api`, `/graphql`, `/v1`, `/swagger`, `/openapi`, etc.
2. **Auth & Admin** — Login pages, dashboards, admin panels, 2FA, SSO
3. **Parameters** — URLs with query strings (`?param=value`)
4. **JS & Config Files** — `.js`, `.json`, `.yaml`, `.env`, `.config`, secrets
5. **Directory Paths** — Path-based resources and endpoints

---

## 📥 Installation

### Requirements
- **Go:** 1.21 or higher
- **OS:** Linux, macOS, or Windows
- **External Tools:** Automatically installed via setup scripts

### Automated Installation (Recommended)

**Linux/macOS:**
```bash
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
bash install.sh
```

**Windows (PowerShell):**
```powershell
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
.\install.bat
```

The installer automatically detects your OS and installs all dependencies.

### Manual Installation

**Linux/macOS:**
```bash
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
make build
sudo make install
```

**Windows:**
```powershell
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
go build -o urlshine.exe .
# Add to PATH or move to system directory
```

### Verify Installation
```bash
# Show help
urlshine -h

# Check dependencies
urlshine doctor
```

The `doctor` command displays which tools are installed and provides installation commands for missing ones.

---

## 🚀 Quick Start

### Basic Usage

**Collect URLs from single target:**
```bash
urlshine -a google.com
```

**Full processing pipeline:**
```bash
urlshine -a -c google.com
```

**Multiple targets from file:**
```bash
urlshine -a -c -f targets.txt -o ./results
```

**Selective tools:**
```bash
urlshine -gau -katana -waymore google.com
```

**High-performance mode:**
```bash
urlshine -a -c -t 200 -d 5 google.com
```

### Understanding Output

URLShine creates timestamped output directories:
```
urlshine_20240506_143022/
├── raw/
│   ├── gau.txt              # GAU results
│   ├── katana.txt           # Katana results
│   ├── gospider.txt         # GoSpider results
│   └── ...                  # One file per tool
├── merged_raw.txt           # All URLs combined
├── normalized_urls.txt      # Cleaned & deduplicated
├── api_endpoints.txt        # Categorized: API endpoints
├── auth_admin_urls.txt      # Categorized: Auth pages
├── parameters_urls.txt      # Categorized: URLs with params
├── js_config_urls.txt       # Categorized: JS & config files
├── directories_urls.txt     # Categorized: Directories
├── alive_api_endpoints.txt  # Verified live API endpoints
├── alive_auth_admin_urls.txt # Verified live auth pages
├── alive_parameters_urls.txt # Verified live parameterized URLs
├── alive_js_config_urls.txt  # Verified live JS & config
├── alive_directories_urls.txt # Verified live directories
├── urlshine_report.json     # Machine-readable report
└── urlshine_report.md       # Human-readable summary
```

---

## 🎯 CLI Reference

### Primary Flags

#### Collection Control
```bash
-a, --all                 Run all 9 URL collection tools (default if no tools specified)
-c, --complete            Full pipeline: merge, normalize, categorize, alive-check
```

#### Collection Tools (Choose One or More)
```bash
-g, --gau                 GetAllUrls (passive, multi-source)
-k, --katana              Katana (active, JavaScript crawler)
-w, --gospider            GoSpider (active, HTML & JS crawler)
-m, --waymore             Waymore (passive, advanced Wayback scraper)
-b, --waybackurls         Wayback URLs (passive, Wayback Machine)
-r, --hakrawler           Hakrawler (active, HTML crawler)
-x, --xnlinkfinder        xnLinkFinder (passive, JS link extraction)
-u, --gobuster            Gobuster (active, directory discovery)
-i, --dirb                Dirb (active, directory brute-force)
```

#### Execution Parameters
```bash
-t, --threads INT         Parallel threads (default: 50, max: 500)
-d, --depth INT           Crawl depth for active tools (default: 3)
-f, --file FILE           Input file with targets (one per line)
-o, --output DIR          Output directory (default: urlshine_<timestamp>)
-s, --subs                Include subdomains (default: true)
-v, --verbose             Enable debug logging
```

#### Pipeline Control
```bash
--no-alive                Skip live host verification
--skip-collect            Skip collection, reprocess existing raw data
```

### Command Examples

**Single target, all tools:**
```bash
urlshine -a domain.com
```

**Full pipeline with high performance:**
```bash
urlshine -a -c -t 150 -d 5 domain.com
```

**Multiple targets from file:**
```bash
urlshine -f domains.txt -a -c -o ./scan_results
```

**Specific tools only:**
```bash
urlshine -gau -katana -waymore -c domain.com
```

**Fast mode (collection only, no live check):**
```bash
urlshine -a -c --no-alive domain.com
```

**Reprocess existing data:**
```bash
urlshine --skip-collect -c -o ./urlshine_20240506_143022
```

**Aggressive reconnaissance:**
```bash
urlshine -f massive-list.txt -a -c -t 200 -d 5 -o ./results -v
```

### Help & Diagnostics
```bash
urlshine -h              # Show help menu
urlshine --help          # Show detailed help
urlshine doctor          # Check tool dependencies and status
urlshine -v              # Show version
```

---

## 📊 Output & Reports

### Terminal Summary
URLShine displays a professional ASCII box with scan statistics:
```
╔══════════════════════════════════════════════════════════════════════════════╗
║                                URLShine — SCAN SUMMARY
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  Target(s)  :  google.com                                                  ║
║  Duration   :  2m 34s                                                       ║
║  Output Dir :  ./urlshine_20240506_143022                                   ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  COLLECTION                                                                  ║
║  Total Raw URLs collected          2,847,391                               ║
║  After Normalization               1,243,891 (56% reduction)                ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  GROUPS                                                                      ║
║  API Endpoints                     18,234     → 14,892 alive               ║
║  Auth / Admin Pages                3,492      → 2,891 alive                ║
║  URLs with Parameters              156,234    → 98,234 alive               ║
║  JS & Config Files                 234,567    → 156,234 alive              ║
║  Directory Paths                   834,367    → 672,440 alive              ║
║                                                                              ║
║  Unique Param Keys                 3,847                                    ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### JSON Report
Machine-readable report with full statistics:
```json
{
  "targets": ["google.com"],
  "total_raw": 2847391,
  "after_norm": 1243891,
  "groups": {
    "api_endpoints": 18234,
    "auth_admin": 3492,
    "parameters": 156234,
    "js_config": 234567,
    "directories": 834367
  },
  "alive_groups": {
    "api_endpoints": 14892,
    "auth_admin": 2891,
    "parameters": 98234,
    "js_config": 156234,
    "directories": 672440
  },
  "unique_params": 3847,
  "output_dir": "./urlshine_20240506_143022",
  "duration_seconds": 154.23
}
```

### Markdown Report
Human-readable summary for documentation:
```markdown
# URLShine — Scan Report

**Generated:** 2024-05-06T14:30:22Z
**Targets:** google.com

### Collection

- **Raw URLs:** 2,847,391
- **Normalized:** 1,243,891 (56% reduction)

### Groups

| Group | URLs | Alive |
|---|---|---|
| API Endpoints | 18,234 | 14,892 |
| Auth / Admin Pages | 3,492 | 2,891 |
| URLs with Parameters | 156,234 | 98,234 |
| JS & Config Files | 234,567 | 156,234 |
| Directory Paths | 834,367 | 672,440 |

**Unique Parameters:** 3,847
**Verified Live URLs:** 1,144,891
```

---

## ⚙️ Configuration

### Performance Tuning

**Default Configuration:**
- Threads: 50
- Depth: 3
- Timeout: 10 seconds per tool

**Recommended Settings:**

| Scenario | Command |
|----------|---------|
| Quick scan | `urlshine -a -t 50 domain.com` |
| Standard scan | `urlshine -a -c -t 100 domain.com` |
| Aggressive scan | `urlshine -a -c -t 150 -d 5 domain.com` |
| Massive scanning | `urlshine -f targets.txt -a -c -t 200 -d 5` |

### Environment Variables
```bash
# Optional: Set custom timeouts
export URLSHINE_TIMEOUT=15

# Optional: Custom output directory
export URLSHINE_OUTPUT=/custom/path
```

---

## 🔧 Troubleshooting

### Installation Issues

**Problem:** `command not found: urlshine`
```bash
# Solution: Verify installation
which urlshine
# If not found, add to PATH
export PATH=$PATH:/usr/local/bin
```

**Problem:** Missing external tools
```bash
# Check which tools are missing
urlshine doctor

# Install missing tools manually (examples)
go install -v github.com/projectdiscovery/katana/cmd/katana@latest
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
```

### Execution Issues

**Problem:** "No targets provided"
```bash
# Solution: Provide target via argument or file
urlshine -a domain.com              # Direct argument
urlshine -f targets.txt -a          # From file
```

**Problem:** Permission denied on output
```bash
# Solution: Check output directory permissions
chmod 755 urlshine_*/
# Or specify different output directory
urlshine -a -o /tmp/results domain.com
```

**Problem:** Out of memory with large wordlists
```bash
# Solution: Reduce threads and depth
urlshine -a -t 30 -d 2 domain.com
```

**Problem:** Tools timing out
```bash
# Solution: Increase timeout and reduce threads
urlshine -a -t 50 domain.com
```

### Output Issues

**Problem:** Empty results
```bash
# Verify target connectivity
ping google.com

# Check if tools are working
urlshine -g google.com  # Test single tool
urlshine doctor         # Verify tool installation
```

---

## 📈 Performance Tips

### For Large Scans
1. **Use aggressive threading:** `-t 150-200` for large wordlists
2. **Increase depth:** `-d 5` for thorough crawling
3. **Run overnight:** Large scans can take hours
4. **Use SSDs:** Faster I/O for large result sets
5. **Monitor resources:** Watch CPU and memory usage

### For Multiple Targets
1. **Batch processing:** Put all targets in a file with `-f`
2. **Parallel execution:** Use OS-level parallelization
3. **Filter results:** Use `grep` to post-process output files

### For Production Use
1. **Verify setup:** Run `urlshine doctor` before scanning
2. **Test first:** Start with small target set
3. **Monitor output:** Check logs for errors
4. **Archive results:** Move results to cold storage after review

---

## 🔒 Security & Legal

### Disclaimer
**URLShine is designed for authorized security testing only.** Unauthorized access to computer systems is illegal. Always obtain written permission before scanning.

### Best Practices
- ✅ Only scan systems you own or have explicit permission to scan
- ✅ Respect robots.txt and terms of service
- ✅ Rate-limit your scans appropriately
- ✅ Review output before taking action
- ✅ Keep scanning logs for compliance

### Data Privacy
- URLShine stores all results locally
- No data is transmitted to external servers
- Results are written to disk in plaintext
- Secure your output directories appropriately

---

## 🤝 Contributing

Contributions are welcome! Areas for improvement:
- Additional collection tools
- Better categorization patterns
- Performance optimizations
- Output format improvements

Please open issues for bugs and feature requests.

---

## 📄 License

URLShine is released under the MIT License. See [LICENSE](LICENSE) for details.

---

## ❓ FAQ

**Q: How long does a typical scan take?**
A: Depends on target size and depth. Small targets: 1-5 minutes. Large targets: 10-30 minutes. Massive scans: Hours.

**Q: Can I run URLShine on Windows?**
A: Yes, fully supported. Use `install.bat` for setup or build manually with `go build`.

**Q: Can I use only specific tools?**
A: Yes, specify tools individually: `urlshine -gau -katana domain.com`

**Q: What's the difference between `-a` and `-a -c`?**
A: `-a` collects URLs only. `-a -c` adds full processing: merging, deduplication, categorization, and live verification.

**Q: How do I update URLShine?**
A: `git pull` to get latest code, then `make build && sudo make install`

**Q: Can I reprocess existing data?**
A: Yes: `urlshine --skip-collect -c -o ./urlshine_<existing_dir>`

**Q: What tools must I have installed?**
A: All required tools are installed automatically. Use `urlshine doctor` to verify.

---

<div align="center">

Made with ❤️ for security researchers

[⬆ Back to top](#-urlshine)

</div>
