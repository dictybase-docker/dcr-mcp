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
	t.Parallel()
	requireHelper := require.New(t)
	logger := log.New(os.Stderr, "", 0)

	tool, err := NewMarkdownTool(logger)
	requireHelper.NoError(err, "NewMarkdownTool should not return an error")
	requireHelper.NotNil(tool, "Tool should not be nil")
	requireHelper.Equal("markdown", tool.GetName(), "Tool name should be 'markdown'")
	requireHelper.NotNil(tool.GetSchema(), "Tool schema should not be nil")
}

func TestHandler(t *testing.T) {
	t.Parallel()
	requireHelper := require.New(t)
	logger := log.New(os.Stderr, "", 0)

	tool, err := NewMarkdownTool(logger)
	requireHelper.NoError(err, "NewMarkdownTool should not return an error")

	// Create a test request
	request := mcp.CallToolRequest{}
	request.Params.Name = "markdown"
	request.Params.Arguments = map[string]interface{}{
		"content": "# Test Heading\n\nTest paragraph.",
	}

	// Call the handler
	result, err := tool.Handler(context.Background(), request)
	requireHelper.NoError(err, "Handler should not return an error")
	requireHelper.NotNil(result, "Result should not be nil")

	// Check the result exists and is not an error
	requireHelper.NotNil(result, "Result should not be nil")
	requireHelper.False(result.IsError, "Result should not be an error")

	// Simply assert that we get at least one content item
	requireHelper.NotEmpty(result.Content, "Result should have at least one content item")

	// Test validation error
	invalidRequest := mcp.CallToolRequest{}
	invalidRequest.Params.Name = "markdown"
	invalidRequest.Params.Arguments = map[string]interface{}{}

	_, err = tool.Handler(context.Background(), invalidRequest)
	requireHelper.Error(err, "Handler should return an error for invalid request")
}
