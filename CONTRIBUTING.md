# Contributing to URLShine

**Welcome to the vanguard of reconnaissance.** 

URLShine is a high-performance framework designed for the elite. To maintain its surgical precision and operational speed, we hold all contributions to a professional, high-standard. If you are looking to sharpen this tool, read this guide.

## 🛡️ Operator Code of Conduct

- **Professionalism First:** Maintain an objective, technical, and professional tone in all communications.
- **Code over Ego:** Focus on the efficiency, security, and stability of the codebase.
- **Security Integrity:** Report any discovered vulnerabilities privately. Do not disclose them in public issues.
- **Precision:** Submit focused, well-tested, and documented changes.

## 🛠️ Getting Started

### Technical Prerequisites
- **Go:** 1.21+ (Core engine requirement)
- **Git:** For version control and collaboration
- **Arsenal:** All external reconnaissance tools must be installed via `install.sh` or `install.bat`.

### Development Environment Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR-USERNAME/UrlShine.git
   cd UrlShine
   ```

2. **Provision Dependencies**
   ```bash
   go mod download
   bash install.sh  # or install.bat on Windows
   ```

3. **Build and Validate**
   ```bash
   make build
   ./urlshine -h
   ```

## 📐 System Architecture

### The Engine Room
```
UrlShine/
├── cmd/                    # CLI Command Center
│   ├── root.go            # Primary entry point & orchestrator
│   └── doctor.go          # System health & dependency auditor
├── internal/
│   ├── banner/            # Tactical identity (ASCII)
│   ├── logger/            # Structured telemetry system
│   ├── collector/         # Enumeration engine execution
│   ├── merger/            # Dataset aggregation
│   ├── normalizer/        # Signal extraction & deduplication
│   ├── splitter/          # Attack vector categorization
│   ├── alive/             # Host verification & probing
│   ├── reporter/          # Executive & machine reporting
│   └── utils/             # Core utility functions
├── main.go                # Application bootstrapper
├── go.mod, go.sum         # Dependency management
├── Makefile               # Build & deployment targets
└── README.md              # Tactical documentation
```

## 💻 Engineering Standards

### Go Implementation Guidelines
- **Formatting:** Strict adherence to `gofmt`. Run `go fmt ./...` before every commit.
- **Verification:** Use `go vet ./...` to ensure code correctness.
- **Modularity:** Keep functions surgical—single purpose, minimal side effects.
- **Clarity:** Variable names must be descriptive. Avoid ambiguity.
- **Documentation:** All exported functions must have a clear, concise comment explaining the *why*, not the *what*.

### Documentation Requirements
```go
// Package collector orchestrates the concurrent execution of URL enumeration tools.
package collector

// Run executes the selected toolset against targets with a defined thread limit.
// It ensures high-concurrency execution while maintaining system stability.
func Run(targets []string, opts Options) error {
    // implementation
}
```

### Naming Conventions
- **Packages:** lowercase, single-word (e.g., `collector`, `normalizer`).
- **Functions:** CamelCase, action-oriented (e.g., `NormalizeFile`).
- **Variables:** Concise but meaningful (e.g., `urlCount` instead of `c`).
- **Constants:** PascalCase or UPPER_SNAKE_CASE for global constants.

### Error Handling
**Do not panic in production code.** Wrap errors with context to allow precise debugging.
```go
// Correct: Provide context for the failure
if err != nil {
    return fmt.Errorf("collector: failed to execute gau: %w", err)
}

// Incorrect: Generic failure
if err != nil {
    panic(err)
}
```

## 🚀 Operationalizing Changes

### Bug Fixes
1. **Isolate:** `git checkout -b fix/issue-description`
2. **Execute:** Implement a minimal, focused fix.
3. **Validate:** Test against the specific failure case.
4. **Document:** Update documentation if the fix alters behavior.
5. **Commit:** Use a clear, professional message: `fix: resolve URL deduplication race condition`.

### New Feature Integration
1. **Strategize:** Open an issue to discuss the tactical value of the feature.
2. **Develop:** `git checkout -b feature/feature-name`.
3. **Implement:** Ensure the code is clean, modular, and documented.
4. **Stress Test:** Verify under both ideal and adversarial (failure) conditions.
5. **Deploy:** Commit and submit a detailed Pull Request.

### Performance Optimization
1. **Benchmark:** Establish a baseline. You must prove the gain with data.
2. **Optimize:** Focus on reducing CPU/Memory overhead or increasing throughput.
3. **Verify:** Ensure zero regressions in functionality.
4. **Report:** Document the exact performance delta (e.g., "Reduced memory usage by 20%").

## 🧪 Testing Protocol

### Manual Validation
```bash
# Single target audit
./urlshine -a example.com

# Specific tool chain validation
./urlshine -gau -katana example.com

# Full pipeline stress test
./urlshine -a -c example.com

# Mass target verification
./urlshine -f targets.txt -a -c
```

### Automated Quality Assurance
```bash
# Code formatting check
go fmt ./...

# Static analysis
go vet ./...

# Unit & Integration test suite
go test ./...
```

## 📨 The Pull Request Process

1. **Synchronize:** `git pull origin main`
2. **Title:** Be concise and professional (e.g., "feat: integrate Shodan API collection").
3. **Context:** Explain the tactical necessity of the change.
4. **Evidence:** Provide logs, screenshots, or benchmark results proving the change works.
5. **Scope:** One PR = One logical change. Do not bundle unrelated fixes.

### PR Template
```markdown
## Tactical Summary
Briefly describe the impact of this change.

## Change Type
- [ ] Bug fix (Corrects existing behavior)
- [ ] New Feature (Expands the arsenal)
- [ ] Performance Optimization (Increases throughput)
- [ ] Documentation Update (Improves clarity)

## Validation Plan
How was this tested? Provide the command and the expected output.

## Evidence
Attach logs or screenshots of successful execution.

## Related Issues
Fixes #123 | Relates to #456
```

## 🎯 Priority Areas for Contribution

### 🔴 High Priority (Critical)
- **New Arsenal/Engines:** Integration of elite tools (e.g., Shodan, Chaos, IronSource).
- **Advanced Vector Patterns:** Improving the regex for identifying high-value endpoints.
- **Core Pipeline Optimization:** Reducing the overhead of the merger and normalizer.

### 🟡 Medium Priority (Enhancements)
- **Executive HTML Reports:** Transforming data into visual attack surface maps.
- **Advanced Filtering:** Adding exclusion patterns and custom include/exclude lists.
- **Network Proxy Support:** Integrating SOCKS5/HTTP proxies for stealth.

### 🔵 Documentation
- **Advanced Use-Case Guides:** Documenting complex recon workflows.
- **Troubleshooting Matrices:** Expanding the error resolution guide.

## 🐛 Bug Reporting

### Pre-Flight Check
- Search existing issues to avoid duplicates.
- Ensure you are running the latest commit from `main`.
- Isolate the bug to a single tool or stage.

### Report Format
**Environment**
- OS: [e.g., Ubuntu 22.04 / Windows 11]
- Go Version: [e.g., 1.22]
- URLShine Version: [e.g., 2.0.0]

**The Defect**
Clear description of the observed vs. expected behavior.

**Steps to Reproduce**
1. `urlshine -a <target>`
2. Observe `<error>`

**Telemetry**
Attach the full terminal output or log file.

## 🔒 Security Disclosures

**NEVER open a public issue for a security vulnerability in URLShine.**

To report a flaw:
1. Send an encrypted email to `security@urlshine.dev` (if available).
2. Use GitHub's private security advisory feature.
3. Provide a Proof-of-Concept (PoC) and a suggested remediation.

## 📜 Licensing
By contributing to URLShine, you agree that your work is licensed under the **MIT License**.

---

**Precision. Power. Performance. Help us build the ultimate recon framework.**
