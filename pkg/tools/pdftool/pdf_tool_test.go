package pdftool

import (
	"context" // Add this line
	"log"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestNewPdfTool(t *testing.T) {
	r := require.New(t)
	logger := log.New(os.Stderr, "[pdf-test] ", 0)

	tool, err := NewPdfTool(logger)
	r.NoError(err, "NewPdfTool should not return an error")
	r.NotNil(tool, "Tool should not be nil")
	r.Equal(
		"markdown_to_pdf",
		tool.GetName(),
		"Tool name should be 'markdown_to_pdf'",
	)
	r.Equal(
		"Converts markdown content to a PDF document.",
		tool.GetDescription(),
		"Tool description mismatch",
	)
	r.NotNil(tool.GetSchema(), "Tool schema should not be nil")
	r.Equal("markdown_to_pdf", tool.GetTool().Name, "MCP Tool name mismatch")

	// Check schema details
	schema := tool.GetSchema()
	r.Equal("object", schema.Type)
	r.Contains(schema.Properties, "content")
	_, ok := schema.Properties["content"]
	r.True(ok, "Schema should have 'content' property")
	r.Contains(schema.Required, "content")
}

func TestHandler(t *testing.T) {
	r := require.New(t)
	// Use a logger that writes to stderr for visibility during tests
	logger := log.New(os.Stderr, "[pdf-test-handler] ", log.LstdFlags)

	tool, err := NewPdfTool(logger)
	r.NoError(err, "NewPdfTool should not return an error")

	// Create a test request
	request := mcp.CallToolRequest{}
	request.Params.Name = "markdown_to_pdf"
	request.Params.Arguments = map[string]interface{}{
		"content": "# Hello PDF\n\nThis is **bold** text.\n\n* And an italic list item.",
	}

	// Execute the handler
	result, err := tool.Handler(context.Background(), request)

	// Assertions
	r.NoError(err, "Handler should not return an error")
	r.NotNil(result, "Result should not be nil")
	r.NotEmpty(result.Content, "Result content slice should not be empty")
}

func TestHandlerMissingContent(t *testing.T) {
	r := require.New(t)
	logger := log.New(os.Stderr, "[pdf-test] ", 0)

	tool, err := NewPdfTool(logger)
	r.NoError(err, "NewPdfTool should not return an error")

	// Create a test request with missing content
	request := mcp.CallToolRequest{}
	request.Params.Name = "markdown_to_pdf"
	request.Params.Arguments = map[string]interface{}{} // Missing 'content'

	// Execute the handler
	result, err := tool.Handler(context.Background(), request)

	// Assertions
	r.Error(err, "Handler should return an error for missing content")
	r.Nil(result, "Result should be nil on error")
	r.Contains(err.Error(), "missing required parameter: content")
}
