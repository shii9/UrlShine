<div align="center">

# 🔗 URLShine
**The Definitive URL Enumeration & Attack Surface Mapping Framework**

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-lightgrey.svg)]()

**Weaponizing reconnaissance for elite security professionals.** 
URLShine obliterates the noise and automates the grind, orchestrating 9 concurrent elite enumeration engines with surgical precision. Intelligent deduplication, vector-based categorization, and high-speed host verification—engineered for maximum impact.

[Installation](#-installation) • [Quick Start](#-quick-start) • [CLI Reference](#-cli-reference) • [Features](#-features) • [Output](#-output) • [FAQ](#-faq)

</div>

---

## 📋 Overview

URLShine is not just a wrapper; it is a high-performance reconnaissance pipeline. By fusing passive archival data with aggressive active crawling, it transforms raw target data into a structured map of an organization's attack surface. 

**Engineered for:**
- 🎯 **Elite Bug Bounty Hunters** — Find the hidden endpoints before anyone else.
- 🛡️ **Red Team Operators** — Map infrastructure with surgical accuracy.
- 🔍 **Advanced Security Researchers** — Analyze attack vectors at scale.
- 🏢 **Professional Penetration Testers** — Automate the reconnaissance phase of any engagement.

---

## ⚡ Features

### 🛠️ The Collection Arsenal
| Engine | Type | Strategic Purpose |
|:---|:---|:---|
| **GAU** | Passive | Massive archive aggregation (100+ threads) |
| **Katana** | Active | Next-gen JS crawler with deep parameter extraction |
| **GoSpider** | Active | High-speed HTML/JS crawler (sitemaps, robots.txt) |
| **Waymore** | Passive | Advanced Wayback Machine deep-scraping |
| **Waybackurls** | Passive | Rapid Wayback Machine URL extraction |
| **Hakrawler** | Active | Fast HTML content crawling with header analysis |
| **xnLinkFinder** | Passive | Precise JS link & config extraction |
| **Gobuster** | Active | High-concurrency directory discovery |
| **Dirb** | Active | Comprehensive directory brute-force enumeration |

### ⚙️ The Processing Pipeline
- 🚀 **Extreme Concurrency** — All engines execute in parallel via optimized thread pooling.
- 💎 **Surgical Deduplication** — Hash-based filtering to ensure zero redundancy.
- 🧼 **Industrial Normalization** — Aggressively strips static noise, redundant ports, and invalid payloads.
- 🎯 **Vector-Based Categorization** — Intelligently splits findings into 5 high-value attack vectors.
- ⚡ **Live Host Verification** — High-speed HTTP probing to separate the signal from the noise.
- 📊 **Executive Reporting** — Professional JSON, Markdown, and terminal summaries.

### 🛡️ Targeted Attack Vectors
1. **API Endpoints** — High-value targets: `/api`, `/graphql`, `/v1`, `/swagger`, `/openapi`.
2. **Auth & Admin** — Critical entry points: Login portals, admin panels, 2FA, SSO.
3. **Parameterized URLs** — Prime targets for XSS, SQLi, and SSRF (`?param=value`).
4. **JS & Config Leakage** — Secrets and architecture leaks: `.js`, `.json`, `.yaml`, `.env`, `.config`.
5. **Infrastructure Paths** — Mapping the directory structure and hidden endpoints.

---

## 📥 Installation

### ⚠️ Requirements
- **Go:** 1.21+ (Required for core engine)
- **OS:** Linux, macOS, or Windows
- **External Tools:** Handled automatically by the deployment scripts.

### 🚀 Deployment (Recommended)

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
*The installer automatically detects your environment and deploys all required dependencies.*

### 🛠️ Manual Build
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
# Add to PATH for global execution
```

### 🩺 System Validation
```bash
# Verify installation
urlshine -h

# Run the system health check
urlshine doctor
```
The `doctor` command performs a full dependency audit and provides immediate remediation steps for missing tools.

---

## 🚀 Quick Start

### Operational Commands

**Standard Recon (Single Target):**
```bash
urlshine -a google.com
```

**Full Tactical Pipeline (Collect $\rightarrow$ Process $\rightarrow$ Verify):**
```bash
urlshine -a -c google.com
```

**Mass Scale Operation (File-based):**
```bash
urlshine -a -c -f targets.txt -o ./results
```

**Custom Toolset Engagement:**
```bash
urlshine -gau -katana -waymore google.com
```

**Aggressive High-Performance Mode:**
```bash
urlshine -a -c -t 200 -d 5 google.com
```

### Data Architecture
URLShine generates an organized evidence directory for every operation:
```
urlshine_20240506_143022/
├── raw/                        # Unfiltered tool output
│   ├── gau.txt                 # GAU raw data
│   ├── katana.txt              # Katana raw data
│   └── ...
├── merged_raw.txt              # Total combined dataset
├── normalized_urls.txt         # Deduplicated & cleaned signal
├── api_endpoints.txt           # Vector: API Endpoints
├── auth_admin_urls.txt         # Vector: Auth/Admin
├── parameters_urls.txt         # Vector: Parameterized URLs
├── js_config_urls.txt          # Vector: Config & JS Leaks
├── directories_urls.txt       # Vector: Directory Structure
├── alive_api_endpoints.txt     # Verified Live: API
├── alive_auth_admin_urls.txt   # Verified Live: Auth/Admin
├── alive_parameters_urls.txt   # Verified Live: Params
├── alive_js_config_urls.txt    # Verified Live: Config/JS
├── alive_directories_urls.txt  # Verified Live: Directories
├── urlshine_report.json        # Machine-readable evidence
└── urlshine_report.md          # Executive summary
```

---

## 🎯 CLI Reference

### 🕹️ Operational Flags

#### Execution Logic
```bash
-a, --all                 Execute all 9 enumeration engines (Default)
-c, --complete            Full pipeline: Merge $\rightarrow$ Normalize $\rightarrow$ Categorize $\rightarrow$ Verify
```

#### The Arsenal (Specific Tool Selection)
```bash
-g, --gau                 GetAllUrls (Passive, Multi-source)
-k, --katana              Katana (Active, JS Crawler)
-w, --gospider            GoSpider (Active, HTML/JS Crawler)
-m, --waymore             Waymore (Passive, Deep Archive)
-b, --waybackurls         Wayback URLs (Passive, Wayback Machine)
-r, --hakrawler           Hakrawler (Active, HTML Crawler)
-x, --xnlinkfinder        xnLinkFinder (Passive, JS Link Extractor)
-u, --gobuster            Gobuster (Active, Dir Discovery)
-i, --dirb                Dirb (Active, Dir Brute-force)
```

#### Performance Tuning
```bash
-t, --threads INT         Concurrent workers (Default: 50, Max: 500)
-d, --depth INT           Crawling depth for active engines (Default: 3)
-f, --file FILE           Target list (One per line)
-o, --output DIR          Custom output directory
-s, --subs                Include subdomains (Default: true)
-v, --verbose             Enable debug-level telemetry
```

#### Pipeline Overrides
```bash
--no-alive                Bypass live host verification
--skip-collect            Bypass collection; process existing raw data
```

### ⚡ Combat Examples

**Full-spectrum reconnaissance on a single domain:**
```bash
urlshine -a -c -t 150 -d 5 domain.com
```

**Massive target list processing with high concurrency:**
```bash
urlshine -f domains.txt -a -c -t 200 -d 5 -o ./ops_results
```

**Surgical tool selection for rapid discovery:**
```bash
urlshine -gau -katana -waymore -c domain.com
```

**Fast-track discovery (No live check):**
```bash
urlshine -a -c --no-alive domain.com
```

**Data reprocessing from previous raids:**
```bash
urlshine --skip-collect -c -o ./urlshine_20240506_143022
```

---

## 📊 Output & Reports

### Tactical Summary
URLShine delivers a high-density ASCII summary upon completion:
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

---

## ⚙️ Configuration

### Performance Profiles

| Scenario | Strategic Goal | Command |
|:---|:---|:---|
| **Quick Recon** | Rapid surface check | `urlshine -a -t 50 domain.com` |
| **Standard Op** | Balanced depth & speed | `urlshine -a -c -t 100 domain.com` |
| **Aggressive** | Thorough surface mapping | `urlshine -a -c -t 150 -d 5 domain.com` |
| **Massive** | Infrastructure-scale scan | `urlshine -f targets.txt -a -c -t 200 -d 5` |

### Environment Control
```bash
# Override global timeouts
export URLSHINE_TIMEOUT=15

# Define custom evidence storage
export URLSHINE_OUTPUT=/mnt/evidence/urlshine
```

---

## 🔧 Troubleshooting

### Deployment Failures
**Issue:** `command not found: urlshine`
- **Fix:** Verify installation path: `which urlshine`. Ensure `/usr/local/bin` is in your `$PATH`.

**Issue:** Dependency gaps
- **Fix:** Run `urlshine doctor` to identify missing engines and execute the provided installation commands.

### Operational Issues
**Issue:** "No targets provided"
- **Fix:** Supply a direct target (`urlshine -a domain.com`) or a target list (`urlshine -f targets.txt -a`).

**Issue:** I/O Permission Denied
- **Fix:** Execute with appropriate privileges or specify a writable output directory: `urlshine -a -o /tmp/results domain.com`.

**Issue:** Resource Exhaustion (OOM)
- **Fix:** Scale back concurrency and depth: `urlshine -a -t 30 -d 2 domain.com`.

**Issue:** Tool Timeouts
- **Fix:** Reduce thread pressure to allow tools to complete: `urlshine -a -t 50 domain.com`.

---

## 📈 Performance Optimization

### For Enterprise-Scale Scans
1. **Force Concurrency:** Push threads to `150-200` for massive datasets.
2. **Deep Crawling:** Set `-d 5` to uncover deeply nested endpoints.
3. **I/O Throughput:** Always deploy on NVMe/SSD storage for large result sets.
4. **Noise Reduction:** Use `--no-alive` to skip verification and speed up the pipeline.

### For Multi-Target Operations
1. **Batch Processing:** Centralize targets in a file and use the `-f` flag.
2. **Parallelization:** Deploy multiple instances of URLShine across different target subsets.
3. **Post-Processing:** Use `grep` or `awk` to slice through the output directories.

---

## 🔒 Legal & Ethical Boundaries

### ⚠️ Mandatory Disclaimer
**URLShine is a high-power security tool designed exclusively for authorized testing.** 
The use of this tool for unauthorized access to computer systems is strictly illegal. You must possess explicit, written permission from the target organization before initiating any scan.

### Operator Best Practices
- ✅ **Authorization First:** Only scan assets you own or have legal permission to test.
- ✅ **Compliance:** Respect `robots.txt` and the target's Terms of Service.
- ✅ **Stealth & Stability:** Rate-limit scans to avoid triggering WAFs or crashing services.
- ✅ **Audit Trails:** Maintain detailed logs of all scanning activities for compliance.

---

## 🤝 Contributing

We welcome contributions from the security community. We are looking for:
- **New Engines:** Integration of high-value collection tools.
- **Smarter Patterns:** Improved regex for attack vector categorization.
- **Performance Hacks:** Optimizing the Go pipeline for even higher throughput.

Open an issue to discuss proposed changes.

---

## 📄 License
URLShine is released under the **MIT License**. See [LICENSE](LICENSE) for full details.

---

## ❓ FAQ

**Q: How long does a typical operation take?**
A: Scalable. Small targets: 1-5 mins. Complex surfaces: 10-30 mins. Infrastructure-scale: Hours.

**Q: Is Windows fully supported?**
A: Yes. Use `install.bat` for a zero-touch setup or `go build` for manual binaries.

**Q: What is the difference between `-a` and `-a -c`?**
A: `-a` is **Collection Only**. `-a -c` is the **Full Tactical Pipeline** (Merge $\rightarrow$ Normalize $\rightarrow$ Categorize $\rightarrow$ Verify).

**Q: How do I keep the tool updated?**
A: `git pull` $\rightarrow$ `make build` $\rightarrow$ `sudo make install`.

---

<div align="center">

**Engineered for the Elite. Built for the Hunt.**

[⬆ Back to top](#-urlshine)

</div>
