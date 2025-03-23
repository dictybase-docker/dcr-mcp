package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/dictybase/dcr-mcp/pkg/tools/calculator"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Initialize MCP server with a name and version
	mcpServer := server.NewMCPServer("DCR-MCP Server", "0.1.0", 
		server.WithToolCapabilities(true),
	)

	// Create calculator tool
	calcTool, err := calculator.NewCalculator()
	if err != nil {
		log.Fatalf("failed to create calculator tool: %v", err)
	}

	// Register calculator tool
	mcpServer.AddTool(calcTool.GetTool(), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paramsJSON, err := json.Marshal(request.Params.Arguments)
		if err != nil {
			return nil, err
		}
		
		result, err := calcTool.Execute(string(paramsJSON))
		if err != nil {
			return mcp.NewToolResultText("Error: " + err.Error()), nil
		}
		
		return mcp.NewToolResultText(result), nil
	})

	// Create an SSE server for HTTP communication
	sseServer := server.NewSSEServer(mcpServer)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting MCP server on port %s", port)
	if err := sseServer.Start(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}