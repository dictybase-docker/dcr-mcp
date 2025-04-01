# DCR MCP Server

A simple MCP (Model Control Protocol) server implementation using [mcp-go](https://github.com/mark3labs/mcp-go).

## Features

- Git Summary tool for analyzing commit messages
- Markdown tool for converting Markdown to HTML
- PDF tool for extracting text from PDF files
- Email Prompt for drafting casual emails

## Getting Started

### Prerequisites

- Go 1.23 or later

### Installation

You can install the DCR MCP server directly using the Go tool:

```bash
go install github.com/dictybase-docker/dcr-mcp/cmd/server@latest
```

This will install the executable in your `$GOPATH/bin` directory (or `$GOBIN` if set).

### Building from Source

Clone the repository:

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
            "command": "dcr-mcp-server",
            "env": {
                "OPENAI_API_KEY": "apixxxxx",
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

This MCP tool converts Markdown content to HTML using the Goldmark markdown
parser. It supports GitHub Flavored Markdown, syntax highlighting for code
blocks, and other extended features.

##### Features

- GitHub Flavored Markdown (GFM) support
- Syntax highlighting for code blocks
- Emoji support
- YAML metadata extraction
- Typographic extensions
- Table formatting

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
markdown parser and the `goldmark-pdf` renderer. The generated PDF binary data
is then Base64 encoded and returned as a text string.

##### Features

- Converts standard Markdown syntax to PDF format.
- Uses pre-configured fonts (IBM Plex Serif for headings, Open Sans for body, Inconsolata for code).
- Sets a specific link color.
- Returns the generated PDF content as a Base64 encoded string.

##### Usage

###### Parameters

- `content` (required): The markdown content to convert to PDF.

###### Example Response

The tool returns a text result containing the Base64 encoded binary data of the generated PDF.

```text
JVBERi0xLjQKJeLjz9MKMSAwIG9iago8PC9UeXBlL0NhdGFsb2cvUGFnZXMgMyAwIFI+P... (truncated Base64 string)
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
