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
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Fprintf(os.Stderr, "server error %v", err)
	}
}
