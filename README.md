<div align="center">

# 🛰️ URLShine v3.2.0
### *Professional URL Enumeration & Attack Surface Intelligence Engine*

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)]()
[![Tools](https://img.shields.io/badge/engines-10%20Core%20Tools-orange.svg)]()

**Operationalizing Reconnaissance for High-Fidelity Attack Surface Mapping**

URLShine orchestrates 10 specialized URL enumeration engines across 4 operational tiers into a unified reconnaissance pipeline. It aggregates passive historical archives, queries threat intelligence APIs, crawls active JavaScript single-page applications (SPAs), and performs multi-threaded directory enumeration.

[Features](#-key-capabilities) • [The 10-Tool Arsenal](#-the-10-tool-arsenal) • [Installation](#-installation) • [Usage & Workflows](#-quick-start--workflows) • [CLI Reference](#-cli-reference) • [Output Architecture](#-output-architecture)

</div>

---

## ⚡ Key Capabilities

- 🔥 **10-Engine Concurrent Execution (Default)**: Running `urlshine example.com` automatically executes all 10 collection tools in parallel for maximum URL yield.
- 📂 **Universal Input Support**: Accepts single domains (`example.com`), full URLs (`https://target.com/app`), or target text files (`targets.txt` or `-f targets.txt`).
- ⚡ **Native Common Crawl CDX Engine**: Built-in zero-dependency Go client that queries Common Crawl CDX index directly without requiring external Python or curl tools.
- 🧼 **Surgical Noise Reduction**: Industrial URL normalizer and deduplication engine that strips redundant query strings, static asset noise, and duplicate endpoints.
- 🎯 **Attack Vector Categorization**: Automatically categorizes discovered endpoints into 5 attack groups:
  1. `api_urls.txt` — `/api`, `/v1-v4`, `/graphql`, `/swagger`, `/actuator`
  2. `auth_admin_urls.txt` — Login portals, admin panels, OAuth callbacks, registration
  3. `params_urls.txt` — Parameterized URLs ready for XSS, SQLi, SSRF, and IDOR testing
  4. `js_config_urls.txt` — JavaScript files, `.env`, `.git`, `.json`, sensitive configs
  5. `directories_urls.txt` — Internal paths, backup directories, route trees
- 🌐 **Automated Live Probing**: Verifies responsive endpoints using multi-threaded HTTP/HTTPS probing (`alive_urls.txt`).
- 📊 **Comprehensive Reporting**: Generates formatted JSON (`urlshine_report.json`) and Markdown (`urlshine_report.md`) summaries.

---

## 🛠️ The 10-Tool Arsenal

URLShine categorizes its 10 core enumeration engines into **4 operational tiers**:

```
 ┌──────────────────────────────────────────────────────────────────────────┐
 │                           URLSHINE PIPELINE                              │
 └──────────────────────────────────────────────────────────────────────────┘
       │                │                     │                   │
  Tier 1 Archives   Tier 2 OSINT & APIs   Tier 3 Crawlers    Tier 4 Brute-Force
  • GAU             • CommonCrawl CDX     • Katana (JS)      • Gobuster
  • Waymore         • URLFinder           • Hakrawler
  • ParamSpider     • GitHub Endpoints
                    • xnLinkFinder
```

| Tier | Tool Name | Engine Type | Strategic Purpose |
| :--- | :--- | :--- | :--- |
| **Tier 1: Passive Archives** | `gau` | Archive Aggregator | Queries Wayback Machine, CommonCrawl, URLScan, and AlienVault OTX simultaneously. |
| **Tier 1: Passive Archives** | `waymore` | Archive & Body Parser | Scraping Wayback data and downloading response bodies to extract internal endpoints. |
| **Tier 1: Passive Archives** | `paramspider` | Parameter Miner | Mines web archives specifically for parameter-rich URLs (`?id=`, `?redirect=`). |
| **Tier 2: Passive APIs** | `commoncrawl` | Native Go CDX Client | Built-in zero-dependency Go client querying the Common Crawl CDX API directly. |
| **Tier 2: Passive APIs** | `urlfinder` | Passive Intelligence | Discovers endpoints from crt.sh, VirusTotal, and passive datasets. |
| **Tier 2: Passive OSINT** | `github-endpoints`| GitHub Repo OSINT | Discovers leaked internal endpoints in public GitHub repositories. |
| **Tier 2: Passive Analysis** | `xnlinkfinder` | JS/HTML Link Finder | Extracts relative routes, endpoints, and APIs from JavaScript & HTML source files. |
| **Tier 3: Active Crawlers** | `katana` | Headless Chromium | Next-gen JS crawler supporting React, Vue, Angular, XHR capture, and form filling. |
| **Tier 3: Active Crawlers** | `hakrawler` | Fast Go Crawler | Lightweight, high-speed crawler that parses HTML and JS without browser overhead. |
| **Tier 4: Brute-Force** | `gobuster` | Directory Scanner | High-speed multi-threaded directory and file path brute-forcer. |

---

## 📥 Installation

### Prerequisites
- **Go**: `1.21+` installed and configured in your system `PATH`.

### 1. Build from Source
```bash
# Clone the repository
git clone https://github.com/shii9/UrlShine.git
cd UrlShine

# Build the executable
go build -o urlshine .
```

### 2. Verify Tool Dependencies
URLShine graceful-degrades if optional external binaries are missing. Run the `doctor` command to check installed tools:

```bash
./urlshine doctor
```

#### Install Optional External Tools (Recommended for Max Coverage):
```bash
# Passive & Archive Tools
go install github.com/lc/gau/v2/cmd/gau@latest
pip3 install waymore
pip3 install git+https://github.com/devanshbatham/ParamSpider
pip3 install xnlinkfinder

# Passive Intelligence Tools
go install github.com/projectdiscovery/urlfinder/cmd/urlfinder@latest
go install github.com/gwen001/github-endpoints@latest

# Active Crawlers & Brute-Force
go install github.com/projectdiscovery/katana/cmd/katana@latest
go install github.com/hakluke/hakrawler@latest
go install github.com/OJ/gobuster/v3@latest
go install github.com/projectdiscovery/httpx/cmd/httpx@latest
```

---

## 🚀 Quick Start & Workflows

### 1. Single Target Recon (Runs All 10 Tools)
```bash
urlshine example.com
```

### 2. Scan Targets from a File
```bash
urlshine targets.txt
# OR using the explicit flag:
urlshine -f targets.txt
```

### 3. Full Processing Pipeline (Collection + Normalization + Attack Grouping + Verification)
```bash
urlshine example.com -c
```

### 4. Multi-Target Enterprise Reconnaissance (Aggressive)
```bash
urlshine -f targets.txt -c -t 100 -d 5 -o output_directory
```

---

## 📖 CLI Reference

```text
USAGE
  urlshine [target/file ...] [flags]

TARGET INPUT SUPPORT
  • Single Domains:    example.com, sub.example.com
  • Full URLs:         https://example.com/api/v1
  • Target Lists:      targets.txt (direct argument or via -f targets.txt)

PIPELINE FLAGS
  -a, --all                  Run all 10 collection tools (ENABLED BY DEFAULT)
  -c, --complete             Run full processing pipeline (merge, normalize, categorize, verify)

COLLECTION TOOL SELECTION (Run specific tools only)
  -g, --gau                  GetAllUrls (Wayback, CommonCrawl, URLScan, OTX)
  -m, --waymore              Waymore (enhanced Wayback & HTTP response parser)
      --paramspider          ParamSpider (parameter-focused archive miner)
      --commoncrawl          Common Crawl CDX API (built-in native Go client)
      --urlfinder            URLFinder (crt.sh, VirusTotal, passive intel)
      --github-endpoints     GitHub endpoint discovery (leaked repo code)
  -x, --xnlinkfinder         xnLinkFinder (JS/HTML link & endpoint extraction)
  -k, --katana               Katana (JS-capable headless Chromium crawler)
      --hakrawler            Hakrawler (fast Go web crawler)
  -z, --gobuster             Gobuster (high-speed directory/file brute-force)

INPUT / OUTPUT OPTIONS
  -f, --file FILE            Input file containing targets (one target per line)
  -o, --output DIR           Output directory (default: urlshine_<timestamp>)

ADVANCED TUNING
  -t, --threads INT          Parallel thread pool workers (default: 50, range: 1-500)
  -d, --depth INT            Crawl depth for active tools (default: 5)
  -v, --verbose              Enable debug/verbose logging
  -s, --subs                 Include subdomains in enumeration (default: true)
      --no-alive             Skip live host verification stage (fast mode)
      --skip-collect         Skip collection phase and reprocess existing raw data

DIAGNOSTICS
  urlshine doctor            Perform system audit & dependency check
```

---

## 📂 Output Architecture

When running URLShine with the `--complete` (`-c`) pipeline, results are organized into domain-specific folders:

```text
urlshine_20260801_150000/
├── example_com_url/
│   ├── raw/                       # Raw outputs from individual tools
│   │   ├── gau_example_com.txt
│   │   ├── katana_example_com.txt
│   │   └── ...
│   ├── per-tool txt files         # Cleaned per-tool results
│   │   ├── gau.txt
│   │   ├── waymore.txt
│   │   ├── katana.txt
│   │   └── ...
│   ├── merged_urls.txt            # Consolidated per-domain URL set
│   ├── normalized_urls.txt        # Deduplicated & cleaned URLs
│   ├── api_urls.txt               # API & Web Service Endpoints
│   ├── auth_admin_urls.txt        # Authentication & Admin Entry Points
│   ├── params_urls.txt            # Parameterized URLs (Injection Surface)
│   ├── js_config_urls.txt         # JavaScript Files & Sensitive Configs
│   ├── directories_urls.txt       # Directory Structure & Hidden Paths
│   ├── alive_urls.txt             # Verified Live/Responsive Endpoints
│   └── urlshine_report.json       # Structured JSON Recon Report
└── urlshine_report.md             # Markdown Summary Report
```

---

## 📄 License

Distributed under the **MIT License**. See [`LICENSE`](LICENSE) for complete terms.
