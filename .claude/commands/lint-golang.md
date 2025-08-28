---
allowed-tools: Bash(golangci-lint run:*), Bash(golangci-lint:*), Bash(gofumpt:*), Bash(find:*), Bash(wc:*), Bash(grep:*), Bash(awk:*), TodoWrite
argument-hint: [path]
description: Run golangci-lint on Go code and provide a comprehensive analysis report (requires golangci-lint v2+)
model: claude-3-5-haiku@20241022
---

I'll run the golangci-lint tool to analyze your Go codebase and provide a concise summary report.

**Prerequisites:**
- Requires golangci-lint version 2.0 or higher
- First step: Verify version with `golangci-lint --version`

**Lint Analysis Strategy:**
1. Check golangci-lint version (must be v2+)
2. Run basic lint check: `golangci-lint run [path]`
3. For comprehensive analysis: `golangci-lint run --verbose --show-stats [path]`
4. Check available linters: `golangci-lint linters`
5. Parse results and generate detailed report

**Validated Command Flags:**
- `golangci-lint run` - Basic linting (primary command)
- `--verbose` - Detailed execution information
- `--show-stats` - Performance and memory usage statistics
- `golangci-lint linters` - List all available linters and their status

**Report Components:**
- **Executive Summary**: Total issues found with pass/fail status
- **Analysis Details**: Processing time, memory usage, linter count
- **Linter Coverage**: Comprehensive checks for code quality, security, performance, and correctness
- **Configuration**: Details about .golangci.yml config and enabled linters
- **Insights**: Educational explanations about linting results and exclusion patterns
- **Recommendations**: Next steps for maintaining code quality

**Configuration Detection:**
The tool automatically detects and uses `.golangci.yml` configuration files, respecting project-specific linting rules and exclusions for directories like third_party, builtin, and examples.

**Performance Metrics:**
Reports include execution time, memory usage (average and peak), and the number of active linters to help understand the thoroughness of the analysis.

**Version Compatibility:**
- **Minimum Required**: golangci-lint v2.0+
- **Tested Version**: v2.1.6 (confirmed working)
- **Version Check**: Always runs `golangci-lint --version` first to ensure compatibility
- **Fallback**: If version < 2.0, the command will provide guidance on upgrading

**Installation Notes:**
If golangci-lint is not installed or version is insufficient:
```bash
# Install latest version
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Or via curl (Linux/macOS)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```
