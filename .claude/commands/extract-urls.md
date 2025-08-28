---
description: Extract and catalog all URLs and links from codebase files
category: workflow
allowed-tools: Read, Write, Grep, Glob, Bash(find:*), Bash(grep:*), Bash(wc:*)
argument-hint: "[directory-path]"
model: sonnet
---

# Extract URLs

Comprehensively scan codebase files and create an inventory of all URLs and links found.

## Usage

```bash
# Scan current directory
/extract-urls

# Scan specific directory
/extract-urls src/

# Scan specific file types in directory
/extract-urls --type=js,html,md
```

## What gets extracted

### URL Types
- Absolute URLs (https://example.com)
- Protocol-relative URLs (//example.com)
- Root-relative URLs (/path/to/page)
- Relative URLs (../images/logo.png)
- API endpoints and fetch URLs
- Asset references (images, scripts, stylesheets)
- Social media links
- Email links (mailto:)
- Tel links (tel:)
- Anchor links (#section)

### File Types Scanned
- HTML, JavaScript, TypeScript
- CSS, SCSS, Sass
- Markdown, MDX
- JSON, YAML configuration files
- Environment files (.env)
- Package.json, composer.json, etc.

### Extraction Contexts
- HTML attributes (href, src, action, data-*)
- JavaScript strings and template literals
- CSS url() functions
- Markdown link syntax [text](url)
- Configuration values
- Comments containing URLs

## Output Format

Creates a comprehensive report with:

### Summary Statistics
- Total URLs found
- Unique URLs count
- Internal vs external ratio
- File type breakdown

### Categorized Inventory
- **Internal Links**: Site-relative URLs
- **External Links**: Third-party domains
- **API Endpoints**: Backend service URLs
- **Assets**: Images, stylesheets, scripts
- **Social Media**: Social platform links
- **Contact**: Email and phone links

### Detailed Listing
For each URL:
- Full URL
- File path and line number
- Context (navigation, asset, API, etc.)
- Anchor text or surrounding context

### Issue Identification
- Hardcoded localhost URLs
- Mixed HTTP/HTTPS protocols
- Potential broken patterns
- Duplicate URLs across files
- Suspicious or problematic links

## Example Output

```
📊 URL Extraction Report
========================

Summary:
- Total URLs: 127
- Unique URLs: 89
- Internal: 45 (51%)
- External: 44 (49%)

🔗 External Links (44):
- https://api.github.com/repos/... (src/api/github.js:23)
- https://docs.example.com/guide (README.md:15)
- https://cdn.jsdelivr.net/npm/... (index.html:8)

🏠 Internal Links (45):
- /api/users (src/components/UserList.tsx:12)
- ../assets/logo.png (src/Header.tsx:8)
- #navigation (docs/guide.md:45)

⚠️  Issues Found:
- Hardcoded localhost: http://localhost:3000 (src/config.js:5)
- Mixed protocols: 3 HTTP links in HTTPS context
```

This creates an immediately useful inventory for link validation, domain migration, SEO audits, or security reviews.