# Contributing to URLShine

Thank you for your interest in contributing to URLShine! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

- Be respectful and professional
- Focus on the code, not the person
- Help others learn and grow
- Report security issues privately

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Git
- External reconnaissance tools (installed via `install.sh` or `install.bat`)

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR-USERNAME/UrlShine.git
   cd UrlShine
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   bash install.sh  # or install.bat on Windows
   ```

3. **Build and Test**
   ```bash
   make build
   ./urlshine -h
   ```

## Project Structure

```
UrlShine/
├── cmd/                    # CLI command definitions
│   ├── root.go            # Main command entry point
│   └── doctor.go          # Dependency checker
├── internal/
│   ├── banner/            # ASCII art header
│   ├── logger/            # Terminal output formatting
│   ├── collector/         # URL collection tools
│   ├── merger/            # URL aggregation
│   ├── normalizer/        # URL cleaning & deduplication
│   ├── splitter/          # URL categorization
│   ├── alive/             # Live host verification
│   ├── reporter/          # Report generation
│   └── utils/             # Utility functions
├── main.go                # Application entry point
├── go.mod, go.sum         # Dependencies
├── Makefile               # Build targets
└── README.md              # User documentation
```

## Code Style & Standards

### Go Best Practices
- Follow `gofmt` formatting: `go fmt ./...`
- Use `go vet`: `go vet ./...`
- Keep functions small and focused
- Write meaningful variable names
- Add comments for exported functions and packages

### Documentation Requirements
```go
// Package collector executes URL enumeration tools concurrently.
package collector

// Run collects URLs using the configured tools.
// It manages concurrent execution with configurable thread limits.
func Run(targets []string, opts Options) error {
    // implementation
}
```

### Naming Conventions
- **Packages**: lowercase, single word when possible
- **Functions**: CamelCase, descriptive names
- **Variables**: short but meaningful (not `x`, use `count`)
- **Constants**: PascalCase or UPPER_CASE

### Error Handling
```go
// Good
if err != nil {
    return fmt.Errorf("operation name: %w", err)
}

// Bad
if err != nil {
    panic(err)
}
```

## Making Changes

### For Bug Fixes

1. **Create a branch**
   ```bash
   git checkout -b fix/issue-description
   ```

2. **Make your changes**
   - Keep changes focused and minimal
   - Add tests if applicable
   - Update documentation if needed

3. **Test your changes**
   ```bash
   ./urlshine -a test-domain.com
   make build
   ```

4. **Commit with clear messages**
   ```bash
   git commit -m "fix: description of fix

   - Specific change 1
   - Specific change 2"
   ```

### For New Features

1. **Discuss first** — Open an issue to discuss the feature
2. **Create a branch** — `git checkout -b feature/feature-name`
3. **Implement** — Keep code clean and well-documented
4. **Test thoroughly** — Test both success and failure cases
5. **Update documentation** — README, comments, etc.
6. **Commit** — Use clear commit messages

### For Performance Improvements

1. **Measure first** — Benchmark before and after
2. **Document impact** — Show performance gains
3. **Maintain compatibility** — Don't break existing features
4. **Test edge cases** — Ensure no regressions

## Testing

### Manual Testing
```bash
# Test single domain
./urlshine -a example.com

# Test with specific tools
./urlshine -gau -katana example.com

# Test full pipeline
./urlshine -a -c example.com

# Test multiple targets
./urlshine -f targets.txt -a -c

# Check dependencies
./urlshine doctor
```

### Code Quality Checks
```bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Check for issues
go test ./...
```

## Pull Request Process

1. **Update your branch** — `git pull origin main`
2. **Write descriptive PR title** — "Add X feature" or "Fix Y bug"
3. **Describe changes** — Explain what and why
4. **Link issues** — Reference related issues (`Fixes #123`)
5. **Test thoroughly** — Show testing evidence
6. **Keep it focused** — One feature or fix per PR

### PR Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature  
- [ ] Performance improvement
- [ ] Documentation update

## How to Test
Steps to verify the changes work correctly

## Testing Evidence
Screenshots, output examples, test results

## Related Issues
Fixes #123
Related to #456
```

## Areas for Contribution

### High Priority
- [ ] Additional collection tools (IronSource, Shodan, etc.)
- [ ] Improved categorization patterns
- [ ] Performance optimizations
- [ ] Better error messages

### Medium Priority
- [ ] HTML report generation
- [ ] Advanced filtering options
- [ ] Custom output formats
- [ ] Proxy support

### Documentation
- [ ] Additional examples
- [ ] Troubleshooting guides
- [ ] Video tutorials
- [ ] Use case studies

## Reporting Bugs

### Before Reporting
- Check existing issues
- Try latest version
- Include full error messages
- Test on different OS if possible

### Bug Report Template
```markdown
**Environment**
- OS: [e.g., Linux/Windows/macOS]
- Go Version: [e.g., 1.21]
- URLShine Version: [e.g., 2.0.0]

**Description**
Brief description of the bug

**Steps to Reproduce**
1. Run command: `urlshine ...`
2. Observe behavior
3. Expected vs actual

**Error Messages**
Full error output or logs

**Screenshots**
If applicable, terminal output or screenshots
```

## Security Issues

**Do not open public issues for security vulnerabilities.**

Instead:
1. Email: security@urlshine.dev (if applicable)
2. Or open a private security advisory via GitHub
3. Allow reasonable time for response and fix

## Documentation Guidelines

### README Updates
- Keep examples current
- Test all examples before committing
- Use proper markdown formatting
- Add table of contents for long docs

### Code Comments
- Explain the "why", not the "what"
- Use proper grammar and punctuation
- Keep comments current with code
- Avoid obvious comments

### Commit Messages
```
feat: add new collection tool support
fix: resolve URL deduplication issue
docs: update installation guide
perf: optimize URL normalization
chore: update dependencies
```

## Recognition

Contributors are recognized in:
- GitHub contributors page
- CONTRIBUTORS.md file (when created)
- Release notes
- GitHub Sponsors (if applicable)

## Questions?

- **General questions**: Open a discussion on GitHub
- **Bug reports**: Create an issue with details
- **Feature requests**: Open an issue with `[FEATURE]` prefix
- **Documentation**: Check README first, then open an issue

## License

By contributing to URLShine, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to URLShine and helping make it better! 🙏
