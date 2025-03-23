# DCR MCP Server

A simple MCP (Model Control Protocol) server implementation using [mcp-go](https://github.com/mark3labs/mcp-go).

## Features

- Basic MCP server implementation
- Simple calculator tool example

## Getting Started

### Prerequisites

- Go 1.20 or later

### Running the Server

```bash
go run cmd/server/main.go
```

By default, the server runs on port 8080. You can change the port by setting the `PORT` environment variable.

### Tools

#### Calculator Tool

The calculator tool supports basic arithmetic operations:

- Addition
- Subtraction
- Multiplication
- Division

Parameters:
- `operation`: The operation to perform (add, subtract, multiply, divide)
- `operandA`: The first operand
- `operandB`: The second operand

## Testing

Run the tests with:

```bash
go test ./...
```

Or using gotestum:

```bash
gotestum --format-hide-empty-pkg --format testdox --format-icons hivis
```