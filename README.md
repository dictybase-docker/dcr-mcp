# DCR MCP Server

A Model Context Protocol (MCP) server providing development and research utilities: git analysis, literature search, markdown/PDF conversion, and email drafting.

## What is MCP?

Model Context Protocol (MCP) enables AI assistants like Claude to securely access external tools and data sources. This server provides development utilities that your AI assistant can use directly - no need to copy-paste commands or switch between tools.

**How it works:**
- Install this MCP server on your system
- Configure your AI assistant to connect to it
- Ask your AI assistant to use the tools naturally (e.g., "analyze recent git commits" or "find papers about CRISPR")
- The AI assistant calls the tools directly and presents results

**Compatible with:** Claude Desktop, Continue, and other MCP-compatible AI tools.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Installation Options](#installation-options)
- [Configuration](#configuration)
- [Tools Reference](#tools-reference)
  - [🔍 Git Summary](#-git-summary)
  - [🔬 Literature Search](#-literature-search)
  - [📝 Markdown Converter](#-markdown-converter)
  - [📄 PDF Generator](#-pdf-generator)
  - [✉️ Email Prompt](#️-email-prompt)
- [Troubleshooting](#troubleshooting)
- [Development](#development)

## Prerequisites

Before installing the DCR MCP Server, ensure you have:

### System Requirements
- **Go 1.23.8+** (only required if building from source)
- **Operating System:** macOS, Linux, or Windows
- **MCP-compatible AI client** (Claude Desktop, Continue, etc.)

### Required API Keys
- **OpenAI API key** (required for git summary features)
  - Create an account at [OpenAI](https://platform.openai.com/)
  - Generate an API key in your account settings
  - Note: Git analysis features will consume OpenAI tokens

### Optional Dependencies
- **Git** (if using git summary tool with local repositories)
- **Internet connection** (for literature search and external repository analysis)

### MCP Client Setup
You'll need an MCP-compatible client. Popular options:
- **Claude Desktop** - Download from [Anthropic](https://claude.ai/download)
- **Continue** - VS Code extension for code assistance
- **Other MCP clients** - See [MCP documentation](https://modelcontextprotocol.io/clients)

## Quick Start

Get DCR MCP Server running in 3 simple steps:

### 1. Install the Server

```bash
go install github.com/dictybase/dcr-mcp/cmd/server@latest
```

### 2. Configure Your AI Client

Add this to your MCP configuration file (typically in Claude Desktop settings):

```json
{
    "mcpServers": {
        "dcr-mcp": {
            "command": "server",
            "env": {
                "OPENAI_API_KEY": "your-openai-api-key"
            }
        }
    }
}
```

> **💡 Note:** Replace `"your-openai-api-key"` with your actual OpenAI API key from [platform.openai.com](https://platform.openai.com/api-keys)

### 3. Test the Connection

Restart your AI client, then ask:

**"What MCP tools do you have available?"**

You should see these 5 tools:
- **🔍 Git Summary** - Analyze commit messages with AI
- **🔬 Literature** - Fetch research papers by PMID/DOI  
- **📝 Markdown** - Convert to HTML with GFM support
- **📄 PDF** - Generate PDFs from markdown
- **✉️ Email** - Draft casual emails

### Quick Examples

Try these commands with your AI assistant:

```
"Analyze git commits from the React repository in the last 30 days"
"Find research papers about CRISPR gene editing"
"Convert this markdown to PDF: # Hello World\nThis is a test."
```

**Next steps:** See [Installation Options](#installation-options) for alternative installation methods or [Tools Reference](#tools-reference) for detailed usage.

## Installation Options

| Method | Command | Binary Name | Best For |
|--------|---------|-------------|----------|
| **Go Install** | `go install github.com/dictybase/dcr-mcp/cmd/server@latest` | `server` | Quick setup |
| **From Source** | `git clone && make build` | `dcr-mcp-server` | Development |
| **Manual Build** | `go build -o dcr-mcp-server ./cmd/server` | `dcr-mcp-server` | Custom builds |

**Requirements:** Go 1.23.8+, OpenAI API key (for git summaries)

> **Note:** Different installation methods produce different binary names. Use the binary name from the table above in your MCP configuration. 

## Configuration

### MCP Setup

```bash
git clone https://github.com/dictybase/dcr-mcp.git
cd dcr-mcp
make build
```

This creates `bin/dcr-mcp-server`. Copy it to a location in your PATH. 

### Running the Server

Create an MCP JSON configuration file to connect to this server from compatible tools:

```json
{
    "mcpServers": {
        "dcr-mcp": {
            "command": "server",  // Use "dcr-mcp-server" if built from source
            "env": {
                "OPENAI_API_KEY": "apixxxxx"
            }
        }
    }
}
```

Save this configuration as `mcp.json` (or any other name) and load it with your MCP-compatible client.

### Configuration Examples by Installation Method

**If installed via `go install`:**
```json
{
    "dcr-mcp": {
        "command": "server",
        "env": {
            "OPENAI_API_KEY": "your-api-key"
        }
    }
}
```

**If built from source:**
```json
{
    "dcr-mcp": {
        "command": "dcr-mcp-server",  // or full path: "/path/to/bin/dcr-mcp-server"
        "env": {
            "OPENAI_API_KEY": "your-api-key"
        }
    }
}
```

## Tools Reference

### 🔍 Git Summary

This MCP tool generates summaries of git commit messages using OpenAI. It analyzes commit messages within a specified date range and creates a concise, user-friendly summary organized by categories.

#### Features

- Clone any git repository by URL and branch
- Filter commits by date range
- Filter by author
- Generate human-readable summaries using OpenAI
- Format output as markdown with categorized bullet points

#### Usage

##### Parameters

- `repo_url` (required): The URL of the git repository to analyze
- `branch` (required): The branch to analyze
- `start_date` (required): The start date for commit analysis (in any standard format)
- `end_date` (optional): The end date for commit analysis (defaults to current date)
- `author` (required): Filter commits by author name (case-insensitive contains match)
- `api_key` (required): Your OpenAI API key (defaults to OPENAI_API_KEY environment variable)

##### Example Response

```markdown
# Work Summary
**Feature Enhancements**
- Added support for filtering commits by author name. Users can now specify an
optional author parameter to focus on contributions from specific team members.
**Bug Fixes**
- Fixed date parsing issues that were causing incorrect commit ranges. The
system now correctly handles various date formats and timezone
considerations.
**Documentation**
- Added comprehensive README with usage examples and parameter descriptions. New
users will find it easier to understand how to use the tool effectively.
```

### 🔬 Literature Search

This MCP tool fetches comprehensive scientific literature information using PMID (PubMed ID) or DOI identifiers via the dictyBase literature API. It provides access to both PubMed and EuropePMC databases with automatic fallback for optimal data retrieval.

#### Features

- **Multiple Provider Support** - Access PubMed and EuropePMC databases
- **Flexible Identifier Support** - Search by PMID (PubMed ID) or DOI with automatic format normalization
- **Smart Fallback Strategy** - EuropePMC first with PubMed fallback for maximum success rate
- **Rich Metadata Extraction** - Complete article information including authors, abstracts, journal details, citations, and MeSH headings
- **Enhanced Data for EuropePMC** - Additional metadata like open access status, PDF availability, license information, and citation counts
- **Automatic Format Validation** - Input validation and normalization for both PMID and DOI formats
- **Comprehensive Author Information** - Full names, ORCID IDs, and institutional affiliations (when available)
- **MeSH and Chemical Data** - Medical subject headings and chemical compound information
- **Grant Information** - Funding sources and grant details

#### Usage

##### Parameters

- `id` (required): The identifier to search for - either a PubMed ID (PMID) or DOI
  - **PMID Format**: Numeric string (e.g., "12345678")
  - **DOI Format**: Standard DOI format, accepts various prefixes:
    - `10.1038/nature12373`
    - `DOI:10.1038/nature12373` 
    - `https://doi.org/10.1038/nature12373`
- `id_type` (required): Type of identifier - must be either "pmid" or "doi"
- `provider` (optional): Literature provider preference - "pubmed" (default) or "europepmc"
  - For DOI searches, EuropePMC is automatically used regardless of this setting
  - For PMID searches, EuropePMC is tried first with PubMed fallback

##### Example Response

```markdown
## Literature Information

**Title:** CRISPR-Cas9 gene editing for sickle cell disease and β-thalassemia

**Authors:** Victoria T. Frangoul, David Altshuler, M. Domenica Cappellini, Yi-Shan Chen, Jennifer Domm, Brenda K. Eustace, Juergen Foell, Josu de la Fuente, Stephan Grupp, Rupert Handgretinger, Tony W. Ho, Antonis Kattamis, Andreas Kernytsky, Julie Lekstrom-Himes, Amanda M. Li, Franco Locatelli, Markus Y. Mapara, Mariane de Montalembert, Damiano Rondelli, Akshay Sharma, Sujit Sheth, Sandeep Soni, Martin H. Steinberg, Donna Wall, Angela Yen, Selim Corbacioglu

**Journal:** New England Journal of Medicine (2021)

**Abstract:** Sickle cell disease and β-thalassemia are genetic disorders caused by mutations in the β-globin gene that result in altered hemoglobin. Gene editing with clustered regularly interspaced short palindromic repeats and CRISPR-associated protein 9 (CRISPR-Cas9) to disrupt the BCL11A erythroid enhancer in hematopoietic stem and progenitor cells can potentially reactivate fetal hemoglobin (HbF) expression and ameliorate these β-hemoglobinopathies.

**PMID:** 33283989
**DOI:** 10.1056/NEJMoa2031054
**Citations:** 1247

---

**Raw JSON Data:**
```json
{
  "id": "PMC7722121",
  "source": "europepmc",
  "pmid": "33283989",
  "pmcid": "PMC7722121",
  "doi": "10.1056/NEJMoa2031054",
  "title": "CRISPR-Cas9 gene editing for sickle cell disease and β-thalassemia",
  "author_string": "Frangoul VT, Altshuler D, Cappellini MD, Chen YS, Domm J, Eustace BK, Foell J, de la Fuente J, Grupp SA, Handgretinger R, Ho TW, Kattamis A, Kernytsky A, Lekstrom-Himes J, Li AM, Locatelli F, Mapara MY, de Montalembert M, Rondelli D, Sharma A, Sheth S, Soni S, Steinberg MH, Wall DA, Yen A, Corbacioglu S.",
  "authors": [
    {
      "full_name": "Victoria T. Frangoul",
      "first_name": "Victoria T.",
      "last_name": "Frangoul",
      "initials": "VT",
      "orcid": "0000-0002-1234-5678",
      "affiliations": [
        {
          "affiliation": "Department of Pediatrics, Sarah Cannon Research Institute, Nashville, TN"
        }
      ]
    }
  ],
  "abstract": "Sickle cell disease and β-thalassemia are genetic disorders...",
  "journal": {
    "title": "New England Journal of Medicine",
    "medline_abbreviation": "N Engl J Med",
    "issn": "0028-4793",
    "volume": "384",
    "issue": "3"
  },
  "pub_year": "2021",
  "is_open_access": true,
  "has_pdf": true,
  "license": "CC BY-NC-SA",
  "cited_by_count": 1247,
  "language": "eng",
  "pub_types": ["Journal Article", "Clinical Trial"],
  "keywords": ["CRISPR-Cas9", "gene editing", "sickle cell disease", "β-thalassemia"]
}
```

#### Use Cases

- **Research Literature Review** - Quickly gather comprehensive metadata for scientific papers
- **Citation Management** - Extract properly formatted citation information with author details
- **Grant and Funding Analysis** - Identify funding sources and research grants for specific publications
- **Open Access Discovery** - Determine availability and licensing of research papers
- **Author Network Analysis** - Collect author information including ORCID IDs and institutional affiliations
- **Medical Research** - Access MeSH headings and chemical compound data for biomedical literature
- **Citation Analysis** - Track citation counts and research impact
- **Database Integration** - Retrieve structured literature data for research management systems

### 📝 Markdown Converter

This MCP tool converts Markdown content to HTML with GitHub Flavored Markdown (GFM) support.

#### Features

- Full CommonMark and GFM support
- Syntax highlighting for code blocks
- Table rendering
- Task list support
- Automatic link generation

#### Usage

##### Parameters

- `content` (required): The markdown content to convert to HTML

##### Example Response

```html
<h1 id="hello-world">Hello, World!</h1>
<p>This is a <strong>markdown</strong> example with <em>formatting</em>.</p>
```

### 📄 PDF Generator

This MCP tool converts Markdown content into a PDF document using the Goldmark markdown parser and the `goldmark-pdf` renderer. The generated PDF is saved to a file (defaulting to `output.pdf` or a user-specified name). The tool returns a confirmation message indicating the save location.

#### Features

- Converts standard Markdown syntax to PDF format
- Uses pre-configured fonts (IBM Plex Serif for headings, Open Sans for body, Inconsolata for code)
- Sets a specific link color
- Saves the generated PDF to a local file (`output.pdf` by default)
- Allows specifying a custom output filename
- Returns a confirmation message with the filename

#### Usage

##### Parameters

- `content` (required): The markdown content to convert to PDF
- `filename` (optional): The desired filename for the output PDF. If omitted, defaults to `output.pdf`

##### Example Response

The tool returns a text result confirming the file save operation.

```text
PDF successfully saved to output.pdf
```

Or, if a filename was provided:

```text
PDF successfully saved to my_document.pdf
```

### ✉️ Email Prompt

This MCP prompt generates a draft casual email, including the subject line, based on provided sender, recipient, and desired tone. It helps quickly compose informal emails.

#### Features

- Generates email drafts based on key information
- Allows specifying the desired tone (e.g., friendly, professional but casual)
- Uses placeholders for easy customization

#### Usage

##### Parameters

- `from` (required): The sender's email address or name
- `to` (required): The recipient's email address or name

##### Example Invocation

When invoked via an MCP client, you would provide the parameters, and the prompt would return a generated email draft. For example, providing `from: "Alice"`, `to: "Bob"`, might result in a draft like: "Subject: Quick Question\n\nHi Bob, Just wanted to ask a quick question about...". The subject line is generated by the LLM based on the requested content.

## Troubleshooting

### Common Issues

#### Installation Problems

**Problem:** `command not found: server` after go install
```bash
Error: server: command not found
```
**Solution:** 
- Ensure `$GOPATH/bin` is in your PATH: `echo $PATH | grep go`
- Check if binary was installed: `ls $(go env GOPATH)/bin/server`
- Add to PATH: `export PATH=$PATH:$(go env GOPATH)/bin`

**Problem:** Build fails with "module not found" 
```bash
go: module github.com/dictybase/dcr-mcp not found
```
**Solution:**
- Verify repository URL: `git remote -v`
- Ensure Go modules are enabled: `go env GO111MODULE` (should be "on")
- Clear module cache: `go clean -modcache`

#### Configuration Issues

**Problem:** MCP client can't find the server binary
```json
Error: Failed to start MCP server "dcr-mcp"
```
**Solution:**
- Use full path in configuration: `"command": "/full/path/to/dcr-mcp-server"`
- Verify binary exists: `which server` or `which dcr-mcp-server`
- Check binary permissions: `ls -la /path/to/binary`

**Problem:** JSON configuration syntax errors
```json
Error: Invalid JSON in MCP configuration
```
**Solution:**
- Validate JSON syntax online or with `json_pp`
- Ensure no trailing commas in JSON
- Use proper quotes around strings

#### API and Authentication Issues

**Problem:** OpenAI API key errors
```bash
Error: OpenAI API key not found or invalid
```
**Solution:**
- Verify API key is set: `echo $OPENAI_API_KEY`
- Check API key format (starts with "sk-")
- Test API key: `curl -H "Authorization: Bearer $OPENAI_API_KEY" https://api.openai.com/v1/models`

**Problem:** Rate limiting or quota errors
```bash
Error: Rate limit exceeded or quota exceeded
```
**Solution:**
- Check OpenAI usage dashboard
- Wait before retrying requests
- Consider upgrading OpenAI plan if needed

#### Tool-Specific Issues

**Problem:** Git summary tool fails to clone repository
```bash
Error: Failed to clone repository
```
**Solution:**
- Verify repository URL is accessible
- Check if repository is private (may need authentication)
- Ensure sufficient disk space
- Check network connectivity

**Problem:** Literature search returns no results
```bash
Error: No articles found for ID
```
**Solution:**
- Verify PMID/DOI format (PMID: numbers only, DOI: with prefix)
- Try alternative ID if available
- Check if article exists in PubMed/EuropePMC
- Wait a moment and retry (temporary service issues)

**Problem:** PDF generation fails
```bash
Error: Failed to generate PDF
```
**Solution:**
- Check markdown syntax is valid
- Verify file permissions in output directory
- Ensure sufficient disk space
- Try with simpler markdown content first

### Debug Mode

Enable detailed logging by setting debug environment variables:

```bash
# For detailed MCP server logs
DEBUG=1 dcr-mcp-server

# For Go application debugging
GODEBUG=gctrace=1 dcr-mcp-server
```

### Getting Help

If you encounter issues not covered here:

1. **Check server logs** - Look for stderr output from the MCP server
2. **Verify dependencies** - Ensure all prerequisites are met
3. **Test with minimal configuration** - Use simplest possible setup
4. **Check GitHub issues** - Search for similar problems
5. **Create an issue** - Include logs, configuration, and error messages

### Log Messages Reference

The server prefixes log messages by tool:
- `[git-summary]` - Git analysis tool messages
- `[literature]` - Literature search tool messages  
- `[markdown]` - Markdown conversion tool messages
- `[pdf-tool]` - PDF generation tool messages
- `[email-prompt]` - Email prompt tool messages

## Development

### Contributing

1. Fork and create feature branch
2. Follow conventions in `CLAUDE.md`
3. Add tests for new functionality
4. Run: `make fmt && make test && golangci-lint run`
5. Submit pull request

### Debugging

The server logs to stderr with prefixed messages:
- `[git-summary]`, `[literature]`, `[markdown]`, `[pdf-tool]`, `[email-prompt]`

### Testing

Run the tests with:

```bash
go test ./...
```

Or using gotestsum:

```bash
gotestsum --format-hide-empty-pkg --format testdox --format-icons hivis
```