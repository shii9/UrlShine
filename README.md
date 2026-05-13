<div align="center">

# 🛰️ URLShine
**The Premier Attack Surface Intelligence Platform for Infrastructure Mapping**

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-lightgrey.svg)]()

**Operationalizing reconnaissance for high-fidelity asset discovery.**
URLShine transforms raw target data into actionable intelligence, orchestrating a multi-layered pipeline of passive archive mining, active surface enumeration, and deterministic verification. It is engineered to uncover shadow assets and hidden endpoints with surgical precision.

[Installation](#-installation) • [Strategic Workflows](#-strategic-workflows) • [The Intelligence Arsenal](#-the-intelligence-arsenal) • [CLI Reference](#-cli-reference) • [Output Architecture](#-output-architecture)
</div>

---

## 🎯 Intelligence Overview

URLShine is a high-performance reconnaissance framework designed to map an organization's total digital footprint. By fusing passive historical data with aggressive active discovery, it eliminates the blind spots in an attack surface, providing a comprehensive map of all reachable endpoints.

**Engineered for:**
- 🛡️ **Red Team Operators** — High-fidelity mapping of infrastructure for surgical engagement.
- 🎯 **Elite Bug Bounty Hunters** — Discovering "unknown unknowns" and shadow assets before the competition.
- 🔍 **Security Researchers** — Analyzing large-scale exposure surfaces and leak patterns.
- 🏢 **Enterprise Security Teams** — Continuous monitoring of the external attack surface.

---

## ⚡ Strategic Capabilities

### 🔭 Surface Discovery
URLShine doesn't just crawl; it discovers. By leveraging a hybrid of passive archival data (Wayback, GAU) and active probing (Katana, GoSpider), it uncovers endpoints that are no longer linked but remain active.

### 🧼 Surgical Noise Reduction
Raw data is noise. URLShine employs a **Deterministic Deduplication** engine and **Industrial Normalization** to strip redundant parameters, static assets, and invalid payloads, ensuring you only analyze high-value signal.

### 📉 Vector Analysis
Findings are intelligently categorized into high-value attack vectors to streamline the exploitation phase:
1. **API Intelligence** — Identifying `/api`, `/graphql`, `/v1`, and undocumented endpoints.
2. **Privileged Access** — Mapping admin panels, login portals, and SSO entry points.
3. **Input Surface** — Isolating parameterized URLs prime for XSS, SQLi, and SSRF.
4. **Sensitive Leaks** — Detecting `.env`, `.git`, `.config`, and JS-based secret leakage.
5. **Infrastructure Mapping** — Deconstructing the directory structure and hidden paths.

---

## 🚀 Strategic Workflows

Depending on the objective, URLShine can be deployed in different operational modes:

| Scenario | Strategic Goal | Operational Command |
|:---|:---|:---|
| **Quick Raid** | Rapid surface check & low noise | `urlshine -a -t 50 target.com` |
| **Deep Dive** | Maximum aggregation & recursive discovery | `urlshine -a -c -t 150 -d 5 target.com` |
| **Enterprise Map** | Infrastructure-scale mapping (multi-target) | `urlshine -f targets.txt -a -c -t 200 -d 5` |
| **Surgical Strike** | Target specific intelligence (e.g., APIs only) | `urlshine -gau -katana -c target.com` |

---

## 🛠️ The Intelligence Arsenal

URLShine orchestrates a diverse set of engines categorized by their intelligence type:

| Engine | Intel Type | Strategic Purpose |
|:---|:---|:---|
| **GAU** | Passive Archive | Massive historical URL aggregation |
| **Waymore** | Passive Archive | Advanced Wayback Machine deep-scraping |
| **Waybackurls** | Passive Archive | Rapid archival endpoint extraction |
| **Katana** | Active Enumeration | Next-gen JS crawler with deep parameter extraction |
| **GoSpider** | Active Enumeration | High-speed HTML/JS surface mapping |
| **xnLinkFinder** | Passive Analysis | Precise JS link & config extraction |
| **Gobuster** | Active Discovery | High-concurrency directory brute-forcing |
| **Dirb** | Active Discovery | Comprehensive path enumeration |

---

## 📥 Installation

### ⚠️ Requirements
- **Go:** 1.21+ (Core engine)
- **OS:** Linux (Optimized)
- **External Tools:** Managed automatically by the deployment script.

### 🚀 Rapid Deployment
```bash
git clone https://github.com/shii9/UrlShine.git
cd UrlShine
bash install.sh
cd ..
mv UrlShine /usr/local/bin
```

After installation, URLShine will be available as a system command:
```bash
urlshine --help          # Full usage documentation
urlshine doctor          # Dependency audit & health check
urlshine -a -c google.com # Run full enumeration pipeline
```

### 🩺 System Validation
```bash
urlshine doctor    # Verify all dependencies are installed
```

---

## � Usage Guide & Practical Examples

### Getting Started (Beginners)

#### 1️⃣ Simple Single-Target Scan
Perfect for quick reconnaissance on a single domain:

```bash
urlshine -a google.com
```
**What it does:**
- Runs all passive collection engines (no traffic to target)
- Gathers historical URLs from archives
- Outputs results to `urlshine_<timestamp>/raw/` directory
- ⏱️ Takes 1-2 minutes

---

#### 2️⃣ Complete Analysis (Recommended for Most Cases)
Performs collection + categorization + live verification:

```bash
urlshine -a -c example.com
```
**What it does:**
- Collects URLs from all 9 tools (passive + active)
- Merges and deduplicates results
- Categorizes into 5 attack vectors (APIs, Auth pages, Parameters, JS configs, Directories)
- Verifies which URLs are actually live/responding
- Generates summary report
- ⏱️ Takes 3-5 minutes

---

### Targeted Reconnaissance (Intermediate)

#### 3️⃣ Passive-Only Scan (Safe, No Noise)
Uses only archive-based tools for stealthy reconnaissance:

```bash
urlshine -gau -waymore -waybackurls -xnlinkfinder -c target.com
```
**What it does:**
- Only queries historical archives (Wayback Machine, CommonCrawl, etc.)
- No active traffic to the target
- Safe for sensitive engagements
- Fast results with lower false positives
- ⏱️ Takes 1-2 minutes

---

#### 4️⃣ Active Crawling Only (Deep Discovery)
Uses only live crawlers for maximum endpoint discovery:

```bash
urlshine -katana -gospider -c target.com -d 5
```
**What it does:**
- Actively crawls the website and follows links
- Extracts parameters from forms and JavaScript
- Depth 5 = crawl 5 levels deep
- Generates significant traffic to target
- Best for authorized penetration testing
- ⏱️ Takes 5-10 minutes

---

#### 5️⃣ API Endpoint Discovery Only
Focus only on API paths and GraphQL endpoints:

```bash
urlshine -a -c target.com | grep api_endpoints.txt
```
**What it does:**
- Runs full pipeline
- Output saved to `api_endpoints.txt` (tagged separately)
- Ideal for API security testing
- Contains `/api`, `/graphql`, `/v1`, `/swagger` patterns

---

### Advanced Usage (Professionals)

#### 6️⃣ Multi-Target Campaign with Custom Threading
Scan multiple targets with aggressive parallelization:

```bash
urlshine -f targets.txt -a -c -t 200 -d 5 -o ./results/
```
**What it does:**
- Reads targets from `targets.txt` (one per line)
- 200 parallel threads (very aggressive)
- Depth 5 (thorough crawling)
- Saves results to `./results/` directory
- Perfect for enterprise-scale audits
- ⏱️ 10-30 minutes depending on target count

**Example targets.txt:**
```
example.com
target.co.uk
vulnerable-app.io
```

---

#### 7️⃣ Fast Enumeration (Collection Only, No Verification)
Speed-optimized for rapid reconnaissance:

```bash
urlshine -a -c --no-alive target.com -t 100
```
**What it does:**
- Runs full collection pipeline
- Skips HTTP verification (saves 50% of time)
- 100 parallel threads
- Returns all URLs without checking if they're live
- ⏱️ Takes 1-2 minutes

---

#### 8️⃣ Reprocess Existing Data (Skip Collection)
Use previously collected data with new analysis:

```bash
urlshine --skip-collect -c urlshine_20260512_143022/
```
**What it does:**
- Skips URL collection phase
- Recategorizes existing URLs with updated patterns
- Useful for tweaking categorization rules
- Much faster for iterative analysis
- ⏱️ Takes 10-30 seconds

---

### Real-World Scenarios

#### 🎯 Scenario 1: Bug Bounty Hunting
You want to find hidden APIs and admin panels quickly:

```bash
# Step 1: Passive reconnaissance (stealthy)
urlshine -gau -waymore -xnlinkfinder -c target.com

# Step 2: Check what you found
cat urlshine_*/api_endpoints.txt        # APIs found
cat urlshine_*/auth_admin_urls.txt      # Admin pages found
```

---

#### 🎯 Scenario 2: Penetration Testing (Full Authorization)
Comprehensive infrastructure mapping with active crawling:

```bash
# Full aggressive scan
urlshine -a -c target.com -t 150 -d 5

# Results breakdown:
ls urlshine_*/
# Shows: raw/, merged.txt, normalized.txt, api_endpoints.txt, etc.

# Analyze the findings:
cat urlshine_*/*.txt | sort | uniq > all_endpoints.txt
```

---

#### 🎯 Scenario 3: Security Audit of Multiple Targets
Enterprise-level assessment of multiple domains:

```bash
# Create target file
echo "company1.com" > targets.txt
echo "company2.com" >> targets.txt
echo "api.company3.io" >> targets.txt

# Run audit with detailed logging
urlshine -f targets.txt -a -c -t 200 -d 3 -v -o /tmp/audit_results/

# Generate report
ls /tmp/audit_results/
```

---

### Workflow: From Discovery to Exploitation

```bash
# 1. Discover all endpoints
urlshine -a -c target.com

# 2. Extract only live API endpoints
cat urlshine_*/alive_api_endpoints.txt > apis.txt

# 3. Test each API with custom tools
while read api; do
  echo "Testing: $api"
  curl -s "$api" | head -20
done < apis.txt

# 4. Find endpoints with parameters
cat urlshine_*/parameters_urls.txt > vulnerable_params.txt

# 5. Fuzz the parameters
# Use with your favorite fuzzer (wfuzz, ffuf, etc.)
```

---

## 🔍 Understanding the Output Categories

When you run `urlshine`, results are automatically split into 5 categories. Here's what each means:

| Category | What It Contains | Example | Use Case |
|:---|:---|:---|:---|
| **api_endpoints.txt** | REST/GraphQL APIs | `/api/users`, `/graphql` | API security testing |
| **auth_admin_urls.txt** | Login & admin pages | `/admin`, `/login`, `/dashboard` | Access control testing |
| **parameters_urls.txt** | URLs with query strings | `/search?q=test&id=123` | Parameter tampering, XSS, SQLi |
| **js_config_urls.txt** | JavaScript & config files | `.js`, `.json`, `.env`, `.git` | Source code analysis, secret leaks |
| **directories_urls.txt** | Path-based resources | `/upload`, `/admin/settings` | Directory traversal, access testing |

---

## 🚀 Performance Tuning Guide

| Use Case | Recommended Settings | Command |
|:---|:---|:---|
| **Quick Check** | Low threads, passive only | `urlshine -gau -waymore -t 30 target.com` |
| **Standard Scan** | Default settings | `urlshine -a -c target.com` |
| **Aggressive** | High threads, deep crawl | `urlshine -a -c -t 150 -d 5 target.com` |
| **Enterprise** | Maximum parallelization | `urlshine -a -c -t 200 -d 5 -f targets.txt` |
| **Stealthy** | Low threads, passive | `urlshine -gau -t 20 target.com` |

---

## 🛠️ CLI Reference

### Operational Logic
```bash
-a, --all                 Execute all enumeration engines (Default)
-c, --complete            Full Intelligence Pipeline: Collect $\rightarrow$ Normalize $\rightarrow$ Categorize $\rightarrow$ Verify
```

### Performance Tuning
```bash
-t, --threads INT         Concurrent workers (Default: 50, Max: 500)
-d, --depth INT           Crawling depth for active engines (Default: 3)
-f, --file FILE           Target list (One per line)
-o, --output DIR          Custom output directory
-s, --subs                Include subdomains (Default: true)
-v, --verbose             Enable debug-level telemetry
```

### Pipeline Overrides
```bash
--no-alive                Bypass live host verification
--skip-collect            Bypass collection; process existing raw data
```

---

## 📊 Output Architecture

URLShine generates a structured evidence directory for every operation, designed for immediate integration into other security tools.

```
urlshine_TIMESTAMP/
├── raw/                        # Unfiltered engine output
│   ├── gau.txt                 # GAU raw data
│   └── ...
├── merged_raw.txt              # Combined dataset
├── normalized_urls.txt         # Deterministic deduplicated signal
├── api_endpoints.txt           # Vector: API Intelligence
├── auth_admin_urls.txt         # Vector: Privileged Access
├── parameters_urls.txt         # Vector: Input Surface
├── js_config_urls.txt          # Vector: Sensitive Leaks
├── directories_urls.txt       # Vector: Infrastructure Mapping
├── alive_api_endpoints.txt     # Verified Live: API
└── urlshine_report.json        # Machine-readable evidence
```

---

## 🔒 Legal & Ethical Boundaries

**URLShine is a high-power security tool designed exclusively for authorized testing.**
The use of this tool for unauthorized access to computer systems is strictly illegal. You must possess explicit, written permission from the target organization before initiating any scan.

- ✅ **Authorization First:** Only scan assets you own or have legal permission to test.
- ✅ **Compliance:** Respect `robots.txt` and target Terms of Service.
- ✅ **Stability:** Rate-limit scans to avoid triggering WAFs or causing DoS.

---

## 🤝 Contributing
We welcome contributions from the security community to enhance the discovery engines and categorization patterns. Open an issue to discuss proposed improvements.

## 📄 License
URLShine is released under the **MIT License**.

<div align="center">

**Engineered for the Elite. Built for the Hunt.**

[⬆ Back to top](#-urlshine)
</div>
