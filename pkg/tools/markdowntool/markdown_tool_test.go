package markdowntool

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestNewMarkdownTool(t *testing.T) {
	r := require.New(t)
	logger := log.New(os.Stderr, "", 0)

	tool, err := NewMarkdownTool(logger)
	r.NoError(err, "NewMarkdownTool should not return an error")
	r.NotNil(tool, "Tool should not be nil")
	r.Equal("markdown", tool.GetName(), "Tool name should be 'markdown'")
	r.NotNil(tool.GetSchema(), "Tool schema should not be nil")
}

func TestHandler(t *testing.T) {
	r := require.New(t)
	logger := log.New(os.Stderr, "", 0)

	tool, err := NewMarkdownTool(logger)
	r.NoError(err, "NewMarkdownTool should not return an error")

	// Create a test request
	request := mcp.CallToolRequest{}
	request.Params.Name = "markdown"
	request.Params.Arguments = map[string]interface{}{
		"content": "# Test Heading\n\nTest paragraph.",
	}

	// Call the handler
	result, err := tool.Handler(context.Background(), request)
	r.NoError(err, "Handler should not return an error")
	r.NotNil(result, "Result should not be nil")

	// Check the result exists and is not an error
	r.NotNil(result, "Result should not be nil")
	r.False(result.IsError, "Result should not be an error")

	// Simply assert that we get at least one content item
	r.NotEmpty(result.Content, "Result should have at least one content item")

	// Test validation error
	invalidRequest := mcp.CallToolRequest{}
	invalidRequest.Params.Name = "markdown"
	invalidRequest.Params.Arguments = map[string]interface{}{}

	_, err = tool.Handler(context.Background(), invalidRequest)
	r.Error(err, "Handler should return an error for invalid request")
}
