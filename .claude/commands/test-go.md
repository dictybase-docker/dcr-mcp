---
allowed-tools: Bash(gotestsum:*), Bash(gotestum:*), Bash(go test:*), Bash(go list:*), Bash(find:*), Bash(wc:*), Bash(grep:*), Bash(awk:*), TodoWrite
argument-hint: [test-pattern] [--verbose] [--specific]
description: Run Go unit tests with comprehensive reporting using gotestsum or fallback to go test
model: claude-3-5-haiku@20241022
---

I'll run your Go unit tests using the project's preferred testing tools and provide comprehensive test analysis.

**Prerequisites:**
- Prefers gotestsum for enhanced test output formatting
- Falls back to standard `go test` if gotestsum unavailable
- Automatically detects test files and packages

**Test Execution Strategy:**
1. Check for gotestsum availability and version
2. Discover Go test packages in the project
3. Run tests with appropriate formatter and options
4. Provide detailed test analysis and coverage insights
5. Report any test failures with actionable information

**Command Options:**
- **Default**: `gotestsum --format-hide-empty-pkg --format testdox --format-icons hivis`
- **Verbose**: `gotestsum --format-hide-empty-pkg --format standard-verbose --format-icons hivis`
- **Specific Test**: `gotestsum --format-hide-empty-pkg --format testdox --format-icons hivis -- -run TestName ./...`
- **Fallback**: `go test -v ./...` (when gotestsum unavailable)

**Validated Command Patterns:**
- `gotestsum --format-hide-empty-pkg --format testdox --format-icons hivis` - Standard test run
- `gotestsum --format-hide-empty-pkg --format testdox --format-icons hivis -- -run TestPattern ./...` - Specific test pattern
- `gotestsum --format-hide-empty-pkg --format standard-verbose --format-icons hivis` - Verbose output
- `go test -v ./...` - Fallback standard testing
- `go test -short ./...` - Quick test run (excludes long-running tests)

**Report Components:**
- **Test Summary**: Total tests run, passed, failed, skipped
- **Package Coverage**: Tests per package and execution time
- **Failure Analysis**: Detailed failure reports with file/line references
- **Performance Metrics**: Test execution time and resource usage
- **Test Discovery**: Automatic detection of test files and packages
- **Recommendations**: Suggestions for test improvements and coverage

**Test Pattern Matching:**
The command supports Go's test pattern matching:
- Specific test: `-run TestFindSimilar`
- Package pattern: `-run TestFindSimilar ./pkg/...`
- Method pattern: `-run TestUser/Create`
- Regex patterns: `-run "Test.*Integration"`

**gotestsum Format Options:**
- **testdox**: BDD-style readable test output (default)
- **standard-verbose**: Detailed Go test output with full logging
- **testname**: Simple test name listing
- **pkgname**: Package-focused output

**Fallback Behavior:**
If gotestsum is not available, the command will:
1. Inform about gotestsum installation benefits
2. Use standard `go test` with appropriate flags
3. Parse and format output for better readability
4. Provide installation instructions for gotestsum

**Installation Notes:**
To install gotestsum for enhanced test reporting:
```bash
go install gotest.tools/gotestsum@latest
```

**Performance Optimization:**
- Uses parallel test execution by default
- Supports test caching for faster repeated runs
- Can exclude vendor and third-party directories
- Optimized for CI/CD pipeline integration

**Integration Features:**
- Automatically respects Go module boundaries
- Detects and reports test coverage gaps
- Provides actionable failure diagnostics
- Supports build tags and test constraints
