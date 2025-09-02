package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dictybase/dcr-mcp/pkg/prompts"
	"github.com/dictybase/dcr-mcp/pkg/tools/gitsummary"
	"github.com/dictybase/dcr-mcp/pkg/tools/literaturetool"
	"github.com/dictybase/dcr-mcp/pkg/tools/markdowntool"
	"github.com/dictybase/dcr-mcp/pkg/tools/pdftool"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := createMCPServer()

	registerTools(mcpServer)
	registerPrompts(mcpServer)

	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Fprintf(os.Stderr, "server error %v", err)
	}
}

// createMCPServer initializes the MCP server with capabilities.
func createMCPServer() *server.MCPServer {
	return server.NewMCPServer("DCR-MCP Server", "1.0.0",
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
}

// registerTools creates and registers all tools with the MCP server.
func registerTools(mcpServer *server.MCPServer) {
	registerGitSummaryTool(mcpServer)
	registerMarkdownTool(mcpServer)
	registerPdfTool(mcpServer)
	registerLiteratureTool(mcpServer)
}

// registerGitSummaryTool creates and registers the git summary tool.
func registerGitSummaryTool(mcpServer *server.MCPServer) {
	gitSummaryTool, err := gitsummary.NewGitSummaryTool(
		log.New(os.Stderr, "[git-summary] ", log.LstdFlags),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create git-summary tool: %v", err)
		os.Exit(1)
	}
	mcpServer.AddTool(gitSummaryTool.GetTool(), gitSummaryTool.Handler)
}

// registerMarkdownTool creates and registers the markdown tool.
func registerMarkdownTool(mcpServer *server.MCPServer) {
	markdownTool, err := markdowntool.NewMarkdownTool(
		log.New(os.Stderr, "[markdown] ", log.LstdFlags),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create markdown tool: %v", err)
		os.Exit(1)
	}
	mcpServer.AddTool(markdownTool.GetTool(), markdownTool.Handler)
}

// registerPdfTool creates and registers the PDF tool.
func registerPdfTool(mcpServer *server.MCPServer) {
	pdfTool, err := pdftool.NewPdfTool(
		log.New(os.Stderr, "[pdf-tool] ", log.LstdFlags),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create pdf tool: %v", err)
		os.Exit(1)
	}
	mcpServer.AddTool(pdfTool.GetTool(), pdfTool.Handler)
}

// registerLiteratureTool creates and registers the literature tool.
func registerLiteratureTool(mcpServer *server.MCPServer) {
	literatureTool, err := literaturetool.NewLiteratureTool(
		log.New(os.Stderr, "[literature] ", log.LstdFlags),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create literature tool: %v", err)
		os.Exit(1)
	}
	mcpServer.AddTool(literatureTool.GetTool(), literatureTool.Handler)
}

// registerPrompts creates and registers all prompts with the MCP server.
func registerPrompts(mcpServer *server.MCPServer) {
	emailPrompt, err := prompts.NewEmailPrompt(
		log.New(os.Stderr, "[email-prompt] ", log.LstdFlags),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create email prompt: %v", err)
		os.Exit(1)
	}
	mcpServer.AddPrompt(emailPrompt.GetPrompt(), emailPrompt.Handler)
}
