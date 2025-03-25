package markdowntool

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/dictybase/dcr-mcp/pkg/markdown"
	"github.com/mark3labs/mcp-go/mcp"
)

// MarkdownTool is a tool that converts markdown to HTML.
type MarkdownTool struct {
	Name        string
	Description string
	Tool        mcp.Tool
	Logger      *log.Logger
}

// NewMarkdownTool creates a new MarkdownTool instance.
func NewMarkdownTool(logger *log.Logger) (*MarkdownTool, error) {
	// Create the tool with proper schema
	tool := mcp.NewTool(
		"markdown",
		mcp.WithDescription(
			"Converts markdown to HTML with support for GFM, syntax highlighting, and more",
		),
		mcp.WithString(
			"content",
			mcp.Description("The markdown content to convert to HTML"),
			mcp.Required(),
		),
	)
	return &MarkdownTool{
		Name:        "markdown",
		Description: "Converts markdown to HTML with support for GFM, syntax highlighting, and more",
		Tool:        tool,
		Logger:      logger,
	}, nil
}

// GetName returns the name of the tool
func (m *MarkdownTool) GetName() string {
	return m.Name
}

// GetDescription returns the description of the tool
func (m *MarkdownTool) GetDescription() string {
	return m.Description
}

// GetSchema returns the JSON schema for the tool's parameters
func (m *MarkdownTool) GetSchema() mcp.ToolInputSchema {
	return m.Tool.InputSchema
}

// GetTool returns the MCP Tool
func (m *MarkdownTool) GetTool() mcp.Tool {
	return m.Tool
}

// Handler returns a function that handles tool execution requests
func (m *MarkdownTool) Handler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	contentVal, ok := request.Params.Arguments["content"].(string)
	if !ok {
		return nil, errors.New("missing required parameter: content")
	}
	parser := markdown.NewParser()
	html, err := parser.ParseString(contentVal)
	if err != nil {
		return nil, fmt.Errorf("failed to parse markdown: %w", err)
	}
	return mcp.NewToolResultText(html), nil
}
