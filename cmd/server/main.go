package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dictybase/dcr-mcp/pkg/tools/gitsummary"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Initialize MCP server with a name and version
	mcpServer := server.NewMCPServer("DCR-MCP Server", "1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)
	// Create and register git-summary tool
	gitSummaryTool, err := gitsummary.NewGitSummaryTool(
		log.New(os.Stderr, "[git-summary] ", log.LstdFlags),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create git-summary tool: %v", err)
		os.Exit(1)
	}
	mcpServer.AddTool(
		gitSummaryTool.GetTool(),
		gitSummaryTool.Handler,
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
