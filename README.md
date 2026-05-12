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
| **Hakrawler** | Active Enumeration | Fast content crawling with header analysis |
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
```

### 🩺 System Validation
Ensure your environment is battle-ready with the internal health check:
```bash
urlshine doctor
```
The `doctor` command performs a full dependency audit and provides immediate remediation steps for any missing engines.

---

## 🕹️ CLI Reference

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
