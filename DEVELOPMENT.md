# URLShine Engineering & Development Manual

This document serves as the technical blueprint for URLShine developers and maintainers. It outlines the system architecture, data flow, and engineering standards required to maintain this high-performance reconnaissance framework.

## 🏗️ System Architecture

### The Reconnaissance Pipeline
URLShine implements a linear, 5-stage processing pipeline designed for maximum throughput and minimal noise.

```
┌─────────────────────────────────────────────────────────────┐
│ STAGE 1: COLLECTION (The Arsenal)                           │
│ - Concurrent execution of 9 elite URL enumeration tools      │
│ - Raw data streamed to the `raw/` evidence directory         │
│ - Per-tool isolation: gau.txt, katana.txt, etc.             │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│ STAGE 2: AGGREGATION (The Merger)                           │
│ - Fusing disparate tool outputs into a single unified stream  │
│ - Result: `merged_raw.txt`                                   │
│ - Preserves all raw data for auditing and reprocessing      │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│ STAGE 3: SIGNAL EXTRACTION (The Normalizer)                 │
│ - Industrial noise reduction: strips static assets (.jpg, .css)│
│ - Protocol cleanup: removes redundant ports (:80, :443)      │
│ - Hash-based deduplication: zero-redundancy filtering        │
│ - Format validation: ensures strict URL compliance          │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│ STAGE 4: VECTOR CATEGORIZATION (The Splitter)               │
│ - Regex-driven classification of the attack surface         │
│ - High-value target groups:                                 │
│   * API Endpoints                                           │
│   * Auth & Admin Portals                                     │
│   * Parameterized URLs                                      │
│   * JS & Config Leaks                                        │
│   * Infrastructure Paths                                    │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│ STAGE 5: LIVE VERIFICATION (The Prober)                      │
│ - High-speed HTTP probing via `httpx` or Go-native fallback  │
│ - Metadata extraction: status codes, titles, technology stack│
│ - Per-vector alive verification                               │
└─────────────────────────────────────────────────────────────┘
```

## 📦 Core Module Specifications

### `cmd/` — Operational Command Center
- **root.go**: The orchestrator. Uses the Cobra framework to handle CLI flags, validate targets, and trigger the pipeline.
- **doctor.go**: The system auditor. Verifies the presence and version of external dependencies to ensure operational readiness.

### `internal/banner/` — Tactical Identity
- Manages the professional ASCII brand presence.
- Provides system-aware headers showing version, environment, and loaded toolsets.

### `internal/logger/` — Telemetry System
- A structured, thread-safe logging framework.
- **Levels:** `Debug` (Telemetry), `Info` (Standard), `Warn` (Caution), `Error` (Failure), `Silent` (Stealth).
- Implements color-coded output for rapid visual parsing of operational status.

### `internal/collector/` — The Enumeration Engine
- Executes the external toolset (GAU, Katana, GoSpider, etc.) using a semaphore-based worker pool.
- **Concurrency Model:** Configurable thread limits (1-500) to prevent system exhaustion.
- **Isolation:** Each tool writes to its own file in the `raw/` directory to allow for granular reprocessing.

### `internal/merger/` — Dataset Aggregation
- Performs a high-speed concatenation of all raw tool outputs.
- Focuses on raw throughput; deduplication is deferred to the Normalizer stage for better performance.

### `internal/normalizer/` — Signal Processing
- The "cleaning" phase of the pipeline.
- Implements a hash-map based deduplication algorithm to ensure each URL is processed exactly once.
- Aggressively filters non-interesting file extensions to reduce the final attack surface.

### `internal/splitter/` — Vector Classification
- Uses a library of regex patterns to categorize URLs into high-value attack vectors.
- Outputs specialized files (e.g., `api_endpoints.txt`) to allow operators to target specific vulnerabilities.

### `internal/alive/` — Host Verification
- Validates the reachability of discovered endpoints.
- Leverages `httpx` for professional-grade probing (status codes, titles, tech stack).
- Separates the "discovered" surface from the "reachable" surface.

### `internal/reporter/` — Evidence Generation
- Transforms raw statistics into an executive summary.
- **Outputs:**
    - Terminal: High-density ASCII summary.
    - JSON: Machine-readable data for downstream automation.
    - Markdown: Human-readable report for documentation.

### `internal/runner/` — Pipeline Orchestration
- The "brain" of URLShine. Coordinates the transition between the 5 stages.
- Manages overall timing, statistics collection, and error propagation across the pipeline.

---

## 🛠️ Engineering Standards

### Performance & Data Structures
- **Deduplication:** Implemented via `map[string]bool` for $\mathcal{O}(1)$ lookup time.
- **Concurrency:** Semaphore-based pools prevent "thundering herd" issues when calling external tools.
- **I/O:** All file operations are streamed where possible to handle multi-million URL datasets without OOM (Out-of-Memory) crashes.

### Implementation Guidelines
- **Modularity:** Each package must have a single, well-defined responsibility.
- **Error Wrapping:** Use `fmt.Errorf("...: %w", err)` to maintain a clear stack trace of failures.
- **Resource Management:** Always close file handles and network connections using `defer`.

---

## 📉 Operational Tuning

### Performance Baselines
- **Threads:** 50 (Default) | 100 (Balanced) | 200+ (Aggressive/Enterprise).
- **Depth:** 3 (Standard) | 5 (Thorough).
- **Memory:** Optimized for $\approx 500\text{MB}$ per 1 million URLs.

### Optimization Strategies
1. **Throughput:** Increase `-t` for high-bandwidth connections and powerful hardware.
2. **Speed:** Use `--no-alive` to bypass the most time-consuming stage of the pipeline.
3. **Precision:** Adjust regex patterns in `internal/splitter/` to better match specific target environments.

---

## 🧪 Testing & Validation

### The Testing Suite
- **Unit Tests:** `go test ./...` (Core logic validation).
- **Static Analysis:** `go vet ./...` (Code correctness).
- **Integration:** `urlshine -a -c <target>` (Full pipeline validation).

### Performance Profiling
Use the Go pprof tool to identify bottlenecks in the normalizer or merger:
```bash
go tool pprof ./urlshine
```

---

## 🚩 Release Protocol

### Pre-Release Checklist
- [ ] **Version Update:** Update `banner.go` with the new version string.
- [ ] **Integrity Check:** Run the full test suite (`go test ./...`).
- [ ] **Build All:** Generate binaries for all supported platforms via `make build-all`.
- [ ] **Documentation:** Update the README and Development Manual.
- [ ] **Tagging:** Create a semantic version tag (`git tag v2.x.x`).

---

## 🤝 Contribution Workflow
1. **Fork** $\rightarrow$ 2. **Branch** $\rightarrow$ 3. **Code** $\rightarrow$ 4. **Test** $\rightarrow$ 5. **PR**.
Detailed guidelines are in `CONTRIBUTING.md`.

**Engineered for Precision. Built for Impact.**
