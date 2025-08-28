---
allowed-tools: Bash(find:*), Bash(grep:*), Bash(curl:*), Bash(awk:*), Bash(head:*), Bash(wc:*), Read, Task, TodoWrite
argument-hint: [--dir=DIR] [--links] [--mermaid] [--syntax] [--format] [--server=URL]
description: Test markdown files for rendering issues, broken links, diagrams, and content validation in web applications
model: claude-3-5-sonnet@20241022
---

I'll help you validate markdown files for rendering quality, broken links, diagram syntax, and content issues in your markdown-based web application.

**Prerequisites:**
- Markdown files with supported extensions (.md, .markdown, .mdown, etc.)
- Optional: Running web application server for live content testing
- Network access for external link validation

**Validation Categories:**
1. **Content Structure**: Headers, lists, tables, footnotes
2. **Link Validation**: Internal links, external URLs, anchor links
3. **Mermaid Diagrams**: Syntax validation and rendering verification
4. **Code Syntax**: Fenced code blocks and language specifications
5. **Format Compliance**: GitHub Flavored Markdown features

**Test Types:**

**Content Structure Tests:**
- Header hierarchy validation (h1 → h2 → h3 progression)
- Table syntax and formatting verification
- Task list checkbox syntax checking
- Footnote reference and definition matching
- Definition list syntax validation

**Link Validation Tests:**
- Internal markdown file references
- Relative path resolution
- Anchor link target verification
- External URL accessibility (HTTP status codes)
- Image reference validation

**Mermaid Diagram Tests:**
- Syntax parsing for all supported diagram types
- Flowchart node and edge validation
- Sequence diagram participant verification
- Pie chart data format checking
- Class diagram syntax compliance

**Code Block Tests:**
- Fenced code block syntax verification
- Language identifier validation
- Syntax highlighting compatibility
- Code block closure matching
- Triple backtick formatting

**Command Options:**
- **Default**: Run all validation tests on current directory
- **Specific Directory**: `--dir=/path/to/docs` - Test specified directory
- **Link Testing**: `--links` - Focus on link validation only
- **Mermaid Testing**: `--mermaid` - Focus on Mermaid diagram validation
- **Syntax Testing**: `--syntax` - Focus on code syntax validation
- **Format Testing**: `--format` - Focus on markdown format compliance
- **Server Testing**: `--server=http://localhost:PORT` - Test against live server

**Validated Test Patterns:**
- `find . -name "*.md" -type f` - Discover markdown files
- `grep -n "^\[.*\]:" file.md` - Find reference-style links
- `grep -n "^```mermaid" file.md` - Locate Mermaid diagrams
- `curl -I -s http://example.com` - Check external link status
- `grep -n "^#{1,6} " file.md` - Find headers for hierarchy check

**Error Detection:**
- **Broken Links**: 404 errors, unreachable URLs, missing internal files
- **Malformed Mermaid**: Syntax errors, unsupported diagram types
- **Header Issues**: Skipped levels, missing h1, duplicate anchors
- **Table Problems**: Misaligned columns, missing separators
- **Code Block Issues**: Unclosed blocks, invalid language tags

**Report Sections:**
1. **File Discovery**: Count and list of markdown files found
2. **Structure Analysis**: Header hierarchy, table count, list validation
3. **Link Report**: Internal/external link status with error details
4. **Mermaid Validation**: Diagram syntax and type verification
5. **Code Analysis**: Language tags and syntax compliance
6. **Format Compliance**: GFM feature usage and correctness

**Integration with Markdown Web Applications:**
- File matching and routing compatibility
- Markdown processor extension verification
- CSS framework class usage validation
- Client-side rendering readiness (diagrams, etc.)
- Server endpoint accessibility testing

**Usage Examples:**

```bash
# Complete validation of current directory
/markdown-test

# Test specific documentation directory
/markdown-test --dir=./docs

# Focus on link validation only
/markdown-test --links

# Validate Mermaid diagrams
/markdown-test --mermaid

# Test against running web application server
/markdown-test --server=http://localhost:8080

# Combined testing with custom directory and server
/markdown-test --dir=./content --server=http://localhost:3000 --links
```

**Performance Optimization:**
- Parallel processing of multiple files
- Efficient regex patterns for content parsing
- Cached external link checking
- Incremental validation for large document sets
- Smart file filtering based on modification time

**Error Reporting Format:**
```
📋 Markdown Validation Report
═══════════════════════════════

📁 Files Scanned: 23 markdown files in ./docs
⚠️  Issues Found: 5 problems across 3 files

🔗 Link Validation:
   ✅ Internal Links: 45/45 valid
   ❌ External Links: 2/12 broken
      docs/api.md:15 - https://example.com/api (404 Not Found)
      README.md:67 - https://oldsite.com (Connection timeout)

📊 Mermaid Diagrams:
   ✅ Valid Diagrams: 8/8
   ✅ Supported Types: flowchart(3), sequence(2), pie(2), class(1)

🔤 Code Blocks:
   ✅ Syntax Valid: 34/34
   ⚠️  Language Tags: 2 unknown languages found
      guide.md:23 - Language 'pseudocode' not recognized
      tutorial.md:89 - Language 'custom-lang' not recognized

📝 Format Compliance:
   ✅ Headers: Valid hierarchy in all files
   ✅ Tables: 12/12 properly formatted
   ✅ Task Lists: 8/8 valid checkbox syntax
```

**Advanced Features:**
- Custom validation rules configuration
- Integration with CI/CD pipelines
- Markdown linting rule enforcement
- Content freshness checking (last modified dates)
- Cross-reference validation between documents

**Integration Capabilities:**
- Works with `/serve-dev` for live testing
- Compatible with `/cleanup` for test artifact management
- Integrates with existing lint workflows
- Supports batch processing for documentation updates

This command ensures your markdown content is high-quality, properly formatted, and renders correctly in your web application, preventing user-facing issues and maintaining documentation standards.