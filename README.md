# DCR MCP Server

A simple MCP (Model Control Protocol) server implementation using [mcp-go](https://github.com/mark3labs/mcp-go).

## Features

- Basic MCP server implementation
- Git Summary tool for analyzing commit messages

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

## Testing

Run the tests with:

```bash
go test ./...
```

Or using gotestum:

```bash
gotestum --format-hide-empty-pkg --format testdox --format-icons hivis
```
