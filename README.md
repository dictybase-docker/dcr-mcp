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

```bash
git clone https://github.com/dictybase-docker/dcr-mcp.git
cd dcr-mcp
go build -o dcr-mcp ./cmd/server
```
Then copy the compiled binary to a location where it's path could be found by
the shell. 

### Running the Server

Create an MCP JSON configuration file to connect to this server from compatible tools:

```json
{
    "mcpServers": {
        "dcr-mcp": {
            "command": "server",  // for go install
            // or "command": "dcr-mcp-server",  // for source build
            // or "command": "/full/path/to/binary",  // for custom path
            "env": {
                "OPENAI_API_KEY": "apixxxxx",
            }
        }
    }
}
```
Save this configuration as `mcp.json` (or any other name) and load it with your MCP-compatible client.

To add to an existing configuration just add the `dcr-mcp` section only.

```json

{
        "dcr-mcp": {
            "command": "dcr-mcp-server",
            "env": {
                "OPENAI_API_KEY": "apixxxxx",
            }
        }
}

```


### Tools

#### Git Summary Tool

This MCP tool generates summaries of git commit messages using OpenAI. It
analyzes commit messages within a specified date range and creates a concise,
user-friendly summary organized by categories.

##### Features

- Clone any git repository by URL and branch
- Filter commits by date range
- Filter by author
- Generate human-readable summaries using OpenAI
- Format output as markdown with categorized bullet points

##### Usage

###### Parameters

- `repo_url` (required): The URL of the git repository to analyze
- `branch` (required): The branch to analyze
- `start_date` (required): The start date for commit analysis (in any standard format)
- `end_date` (optional): The end date for commit analysis (defaults to current date)
- `author` (required): Filter commits by author name (case-insensitive contains match)
- `api_key` (required): Your OpenAI API key (defaults to OPENAI_API_KEY environment variable)

###### Example Response

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
users will find it easier to understand how to use the tool effectively."
```

#### Markdown Tool

1. Fork and create feature branch
2. Follow conventions in `CLAUDE.md`
3. Add tests for new functionality
4. Run: `make fmt && make test && golangci-lint run`
5. Submit pull request

##### Features

The server logs to stderr with prefixed messages:
- `[git-summary]`, `[literature]`, `[markdown]`, `[pdf-tool]`, `[email-prompt]`

##### Usage

###### Parameters

- `content` (required): The markdown content to convert to HTML


###### Example Response

```html
<h1 id="hello-world">Hello, World!</h1>
<p>This is a <strong>markdown</strong> example with <em>formatting</em>.</p>
```

#### PDF Tool (Markdown to PDF)

This MCP tool converts Markdown content into a PDF document using the Goldmark
markdown parser and the `goldmark-pdf` renderer. The generated PDF is saved
to a file (defaulting to `output.pdf` or a user-specified name). The tool
returns a confirmation message indicating the save location.

##### Features

- Converts standard Markdown syntax to PDF format.
- Uses pre-configured fonts (IBM Plex Serif for headings, Open Sans for body, Inconsolata for code).
- Sets a specific link color.
- Saves the generated PDF to a local file (`output.pdf` by default).
- Allows specifying a custom output filename.
- Returns a confirmation message with the filename.

##### Usage

###### Parameters

- `content` (required): The markdown content to convert to PDF.
- `filename` (optional): The desired filename for the output PDF. If omitted, defaults to `output.pdf`.

###### Example Response

The tool returns a text result confirming the file save operation.

```text
PDF successfully saved to output.pdf
```

Or, if a filename was provided:

```text
PDF successfully saved to my_document.pdf
```

### Email Prompt

This MCP prompt generates a draft casual email, including the subject line,
based on provided sender, recipient, and desired tone. It helps quickly compose informal
emails.

#### Features

- Generates email drafts based on key information.
- Allows specifying the desired tone (e.g., friendly, professional but casual).
- Uses placeholders for easy customization.

#### Usage

###### Parameters

- `from` (required): The sender's email address or name.
- `to` (required): The recipient's email address or name.

###### Example Invocation (Conceptual)

When invoked via an MCP client, you would provide the parameters, and the prompt
would return a generated email draft. For example, providing `from: "Alice"`,
`to: "Bob"`, might result in a draft like: "Subject: Quick Question\n\nHi Bob,
Just wanted to ask a quick question about...". The subject line is generated by
the LLM based on the requested content.

## Testing

Run the tests with:

```bash
go test ./...
```

Or using gotestum:

```bash
gotestum --format-hide-empty-pkg --format testdox --format-icons hivis
```
