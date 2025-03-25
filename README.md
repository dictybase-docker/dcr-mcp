# DCR MCP Server

A simple MCP (Model Control Protocol) server implementation using [mcp-go](https://github.com/mark3labs/mcp-go).

## Features

- Basic MCP server implementation
- Git Summary tool for analyzing commit messages

## Getting Started

### Prerequisites

- Go 1.23 or later

### Running the Server

```bash
go run cmd/server/main.go
```

By default, the server runs on port 8080. You can change the port by setting the `DCR_MCP_PORT` environment variable.

### Tools

#### Git Summary Tool

This MCP tool generates summaries of git commit messages using OpenAI. It analyzes commit messages within a specified date range and creates a concise, user-friendly summary organized by categories.

##### Features

- Clone any git repository by URL and branch
- Filter commits by date range
- Optionally filter by author
- Generate human-readable summaries using OpenAI
- Format output as markdown with categorized bullet points

##### Usage

The tool accepts the following parameters:

```json
{
  "repo_url": "https://github.com/username/repo",
  "branch": "main",
  "start_date": "2023-01-01",
  "end_date": "2023-12-31",
  "author": "required-author-name",
  "api_key": "your-openai-api-key",
  "model_name": "google/gemini-2.0-flash-001",
  "openai_base": "https://openrouter.ai/api/v1" 
}
```

###### Parameters

- `repo_url` (required): The URL of the git repository to analyze
- `branch` (required): The branch to analyze
- `start_date` (required): The start date for commit analysis (in any standard format)
- `end_date` (optional): The end date for commit analysis (defaults to current date)
- `author` (required): Filter commits by author name (case-insensitive contains match)
- `api_key` (optional): Your OpenAI API key (defaults to OPENAI_API_KEY environment variable)
- `model_name` (optional): The OpenAI model to use (defaults to "google/gemini-2.0-flash-001")
- `openai_base` (optional): Custom OpenAI API base URL (defaults to official API)

###### Example Response

"# Work Summary\n\n**Feature Enhancements**\n- Added support for filtering
commits by author name. Users can now specify an optional author parameter to
focus on contributions from specific team members.\n\n**Bug Fixes**\n- Fixed
date parsing issues that were causing incorrect commit ranges. The system now
correctly handles various date formats and timezone
considerations.\n\n**Documentation**\n- Added comprehensive README with usage
examples and parameter descriptions. New users will find it easier to understand
how to use the tool effectively."

## Testing

Run the tests with:

```bash
go test ./...
```

Or using gotestum:

```bash
gotestum --format-hide-empty-pkg --format testdox --format-icons hivis
```
