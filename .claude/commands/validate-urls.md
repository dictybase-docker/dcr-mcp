---
description: Validate URLs for functionality and contextual appropriateness
category: workflow
allowed-tools: Read, Write, WebFetch, Grep, Glob
argument-hint: "[file-path or url-list]"
model: sonnet
---

# Validate URLs

Validate URLs both technically (functionality) and contextually (appropriateness and relevance).

## Usage

```bash
# Validate URLs in a file
/validate-urls README.md

# Validate URLs in directory
/validate-urls docs/

# Validate specific URLs
/validate-urls "https://example.com, https://github.com/user/repo"

# Extract and validate all URLs in project
/validate-urls --extract-all
```

## Validation Types

### Technical Validation
- HTTP status codes (200, 301, 302, 404, 500, etc.)
- Redirect chains and final destinations
- Response times and timeout detection
- SSL certificate validity for HTTPS
- Malformed URL syntax detection

### Contextual Analysis
- Semantic alignment with surrounding text
- Anchor text vs destination content matching
- Content relevance and appropriateness
- Publication date relevance
- Authority and quality assessment

### Security Checks
- Mixed HTTP/HTTPS protocol issues
- Suspicious or potentially malicious domains
- Authentication-protected links
- Regional access restrictions

## Report Format

### Executive Summary
- Total URLs checked
- Functional status breakdown
- Critical issues count
- Recommendations summary

### Detailed Results

For each URL:
- **Status**: Working, Dead, Redirect, Suspicious
- **Response**: HTTP code and response time
- **Context**: Surrounding text and anchor text
- **Relevance**: Highly relevant, Questionable, Misaligned
- **Issues**: Specific problems found
- **Recommendation**: Keep, Update, Replace, Remove

### Issue Categories

**🔴 Critical Issues** (Must Fix):
- Dead links (404, 500 errors)
- Security vulnerabilities
- Malformed URLs
- Mixed protocol violations

**🟠 High Priority** (Should Fix):
- Slow-loading links (>5s response)
- Redirect chains (>3 redirects)
- Contextual mismatches
- Outdated content links

**🟡 Medium Priority** (Consider Fixing):
- Suboptimal link text
- Non-HTTPS links where HTTPS available
- Links to homepages vs specific pages
- Multiple links to same content

**🟢 Low Priority** (Optional):
- Minor optimization opportunities
- Style consistency improvements

### Recommended Actions

For problematic URLs:
- **Specific replacement suggestions**
- **Alternative authoritative sources**
- **Update strategies**
- **Removal justifications**

## Example Output

```
🔗 URL Validation Report
========================

Executive Summary:
- URLs Checked: 45
- Working: 38 (84%)
- Dead: 4 (9%)
- Redirects: 3 (7%)
- Critical Issues: 7

🔴 Critical Issues:
- https://old-docs.example.com/guide (404) in README.md:15
  → Replace with: https://docs.example.com/new-guide
  
- http://api.service.com/docs (Mixed protocol) in api.md:8
  → Update to: https://api.service.com/docs

🟠 High Priority:
- https://example.com (Homepage link) in tutorial.md:23
  → Context suggests specific page: https://example.com/getting-started
  → Anchor text "installation guide" doesn't match homepage content

🟡 Medium Priority:
- Link text "click here" is non-descriptive (accessibility issue)
- Multiple links to same GitHub repo could be consolidated
```

## Special Handling

### Edge Cases
- Authentication-required links (noted but not fully validated)
- Dynamic content URLs (flagged for manual review)  
- Time-sensitive or event-specific links
- Regional restrictions and geo-blocking

### Performance
- Batch processing for multiple URLs
- Timeout handling for slow responses
- Rate limiting to avoid being blocked
- Caching for repeated validation runs

This creates actionable reports prioritizing the most critical link issues while providing specific guidance for fixes.