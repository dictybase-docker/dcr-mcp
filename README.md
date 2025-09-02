# DCR MCP Server

A Model Context Protocol (MCP) server providing development and research utilities: git analysis, literature search, markdown/PDF conversion, and email drafting.

## Table of Contents

- [Quick Start](#quick-start)
- [Installation Options](#installation-options)
- [Configuration](#configuration)
- [Tools Reference](#tools-reference)
  - [🔍 Git Summary](#-git-summary)
  - [🔬 Literature Search](#-literature-search)
  - [📝 Markdown Converter](#-markdown-converter)
  - [📄 PDF Generator](#-pdf-generator)
  - [✉️ Email Prompt](#️-email-prompt)
- [Development](#development)

## Quick Start

1. **Install**
   ```bash
   go install github.com/dictybase/dcr-mcp/cmd/server@latest
   ```

2. **Configure** - Add to your MCP configuration:
   ```json
   {
       "mcpServers": {
           "dcr-mcp": {
               "command": "server",
               "env": {
                   "OPENAI_API_KEY": "your-api-key-here"
               }
           }
       }
   }
   ```

3. **Use** - The server provides these tools:
   - **Git Summary** - Analyze commit messages with AI
   - **Literature** - Fetch research papers by PMID/DOI  
   - **Markdown** - Convert to HTML with GFM support
   - **PDF** - Generate PDFs from markdown
   - **Email** - Draft casual emails

## Installation Options

| Method | Command | Binary Name | Best For |
|--------|---------|-------------|----------|
| **Go Install** | `go install github.com/dictybase/dcr-mcp/cmd/server@latest` | `server` | Quick setup |
| **From Source** | `git clone && make build` | `dcr-mcp-server` | Development |
| **Manual Build** | `go build -o dcr-mcp-server ./cmd/server` | `dcr-mcp-server` | Custom builds |

**Requirements:** Go 1.23.8+, OpenAI API key (for git summaries) 

## Configuration

### MCP Setup

Use the correct command based on your installation method:

```json
{
    "mcpServers": {
        "dcr-mcp": {
            "command": "server",  // for go install
            // or "command": "dcr-mcp-server",  // for source build
            // or "command": "/full/path/to/binary",  // for custom path
            "env": {
                "OPENAI_API_KEY": "your-api-key-here"
            }
        }
    }
}
```

The server auto-starts with your MCP client and provides structured logging to stderr for debugging.


## Tools Reference

### 🔍 Git Summary

Analyzes git commits and generates AI-powered summaries categorized by type.

**Parameters:**
- `repo_url` (required) - Git repository URL
- `branch` (required) - Branch to analyze  
- `start_date` (required) - Start date for analysis
- `end_date` (optional) - End date (defaults to now)
- `author` (required) - Filter by author name
- `api_key` (required) - OpenAI API key

**Output:** Markdown summary with categorized bullet points (features, bugs, docs, etc.)

### 🔬 Literature Search

Fetches scientific papers using PubMed IDs or DOIs via dictyBase API.

**Parameters:**
- `id` (required) - PMID (e.g., "12345678") or DOI (e.g., "10.1038/nature12373")
- `id_type` (required) - "pmid" or "doi"
- `provider` (optional) - "pubmed" (default) or "europepmc"

**Output:** Formatted article metadata with title, authors, abstract, and citation data

### 📝 Markdown Converter  

Converts Markdown to HTML with GitHub Flavored Markdown support.

**Parameters:**
- `content` (required) - Markdown content to convert

**Features:** GFM, syntax highlighting, emoji, tables, YAML metadata

### 📄 PDF Generator

Converts Markdown to PDF with professional formatting.

**Parameters:**
- `content` (required) - Markdown content
- `filename` (optional) - Output filename (defaults to "output.pdf")

**Output:** PDF file with IBM Plex fonts and confirmation message

### ✉️ Email Prompt

Generates casual email drafts with customizable tone.

**Parameters:**
- `from` (required) - Sender's name or email
- `to` (required) - Recipient's name or email

**Output:** Complete email draft with subject line and friendly, professional tone

## Development

### Build Commands

| Command | Purpose |
|---------|---------|
| `make build` | Build binary to `bin/dcr-mcp-server` |
| `make test` | Run all tests |
| `make test-verbose` | Run tests with detailed output |
| `make fmt` | Format code with `gofumpt` |
| `make clean` | Remove build artifacts |

### Project Structure

```
dcr-mcp/
├── cmd/server/              # Main entry point
├── pkg/
│   ├── tools/              # Core tool implementations
│   │   ├── gitsummary/     # Git analysis + OpenAI
│   │   ├── literaturetool/ # PubMed/DOI fetching  
│   │   ├── markdowntool/   # Markdown → HTML
│   │   └── pdftool/        # Markdown → PDF
│   ├── prompts/            # Email prompt logic
│   └── markdown/           # Shared markdown utilities
└── bin/                    # Built binaries
```

### Contributing

1. Fork and create feature branch
2. Follow conventions in `CLAUDE.md`
3. Add tests for new functionality
4. Run: `make fmt && make test && golangci-lint run`
5. Submit pull request

### Debugging

The server logs to stderr with prefixed messages:
- `[git-summary]`, `[literature]`, `[markdown]`, `[pdf-tool]`, `[email-prompt]`

Enable stderr output in your MCP client for debugging.
