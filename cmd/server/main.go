package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dictybase/dcr-mcp/pkg/tools/calculator"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Initialize MCP server with a name and version
	mcpServer := server.NewMCPServer("DCR-MCP Server", "1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	// Create calculator tool
	calcTool, err := calculator.NewCalculator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create calculator tool: %v", err)
		os.Exit(1)
	}

	// Register calculator tool
	mcpServer.AddTool(
		calcTool.GetTool(),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			paramsJSON, err := json.Marshal(request.Params.Arguments)
			if err != nil {
				return nil, err
			}

			result, err := calcTool.Execute(string(paramsJSON))
			if err != nil {
				return mcp.NewToolResultText("Error: " + err.Error()), nil
			}

			return mcp.NewToolResultText(result), nil
		},
	)

	// Create an SSE server for HTTP communication
	sseServer := server.NewSSEServer(mcpServer)

	// Start server
	port := os.Getenv("DCR_MCP_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	port = fmt.Sprintf(":%s", port)

	log.Printf("Starting MCP server on port %s", port)
	if err := sseServer.Start(port); err != nil {
		fmt.Fprintf(os.Stderr, "server error %v", err)
	}
}
