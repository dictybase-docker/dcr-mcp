package pdftool

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestNewPdfTool(t *testing.T) {
	t.Parallel()
	requireHelper := require.New(t)
	logger := log.New(os.Stderr, "[pdf-test] ", 0)

	tool, err := NewPdfTool(logger)
	requireHelper.NoError(err, "NewPdfTool should not return an error")
	requireHelper.NotNil(tool, "Tool should not be nil")
	requireHelper.Equal(
		"markdown_to_pdf",
		tool.GetName(),
		"Tool name should be 'markdown_to_pdf'",
	)
	requireHelper.Equal(
		"Converts markdown content to a PDF document and saves it to a file.",
		tool.GetDescription(),
		"Tool description mismatch",
	)
	requireHelper.NotNil(tool.GetSchema(), "Tool schema should not be nil")
	requireHelper.Equal("markdown_to_pdf", tool.GetTool().Name, "MCP Tool name mismatch")

	// Check schema details
	schema := tool.GetSchema()
	requireHelper.Equal("object", schema.Type)
	requireHelper.Contains(schema.Properties, "content")
	_, contentExists := schema.Properties["content"] // Get content prop details
	requireHelper.True(contentExists, "Schema should have 'content' property")

	// Check for optional filename parameter
	requireHelper.Contains(
		schema.Properties,
		"filename",
	)
	filenameProp, filenameExists := schema.Properties["filename"]
	requireHelper.True(
		filenameExists,
		"Schema should have 'filename' property",
	)
	filenamePropMap, isValidMap := filenameProp.(map[string]any)
	requireHelper.True(isValidMap, "filename property should be a map")
	requireHelper.Equal(
		"string",
		filenamePropMap["type"],
	)
	requireHelper.Equal(
		"Optional filename for the output PDF. Defaults to 'output.pdf'.",
		filenamePropMap["description"],
	)

	requireHelper.Contains(schema.Required, "content")
	requireHelper.NotContains(
		schema.Required,
		"filename",
	) // Ensure filename is not required
}

func TestHandlerDefaultFilename(t *testing.T) {
	t.Parallel()
	requireHelper := require.New(t)
	// Use a logger that writes to stderr for visibility during tests
	logger := log.New(os.Stderr, "[pdf-test-handler-default] ", log.LstdFlags)

	tool, err := NewPdfTool(logger)
	requireHelper.NoError(err, "NewPdfTool should not return an error")

	defaultFilename := "output.pdf"
	// Ensure the default file does not exist before the test
	_ = os.Remove(defaultFilename)
	// Schedule cleanup after the test
	defer os.Remove(defaultFilename)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "markdown_to_pdf",
			Arguments: map[string]interface{}{
				"content": "# Default File Test\n\nContent here.",
			},
		},
	}

	result, err := tool.Handler(context.Background(), request)

	// Assertions for default filename case
	requireHelper.NoError(err, "Handler should not return an error (default filename)")
	requireHelper.NotNil(result, "Result should not be nil (default filename)")
	requireHelper.NotEmpty(
		result.Content,
		"Result content slice should not be empty (default filename)",
	)

	// Check the success message
	textContent, success := mcp.AsTextContent(result.Content[0])
	requireHelper.True(success, "First content should be text content")
	expectedMsg := fmt.Sprintf(
		"PDF successfully saved to %s",
		defaultFilename,
	)
	requireHelper.Equal(
		expectedMsg,
		textContent.Text,
		"Success message mismatch (default filename)",
	)

	// Check if the file was created
	_, err = os.Stat(defaultFilename)
	requireHelper.NoError(err, "Default output file '%s' should exist", defaultFilename)

	// Optional: Check file content (basic PDF magic bytes)
	pdfBytes, err := os.ReadFile(defaultFilename)
	requireHelper.NoError(err, "Failed to read created file %s", defaultFilename)
	requireHelper.Greater(
		len(pdfBytes),
		4,
		"PDF data too short in file %s",
		defaultFilename,
	)
	requireHelper.Equal(
		[]byte{0x25, 0x50, 0x44, 0x46, 0x2d},
		pdfBytes[:5],
		"File %s should start with PDF magic bytes",
		defaultFilename,
	)
}

func TestHandlerCustomFilename(t *testing.T) {
	t.Parallel()
	requireHelper := require.New(t)
	// Use a logger that writes to stderr for visibility during tests
	logger := log.New(os.Stderr, "[pdf-test-handler-custom] ", log.LstdFlags)

	tool, err := NewPdfTool(logger)
	requireHelper.NoError(err, "NewPdfTool should not return an error")

	customFilename := filepath.Join(
		t.TempDir(),
		"custom_test.pdf",
	) // Use temp dir for custom file
	// No need to remove beforehand, t.TempDir handles cleanup

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "markdown_to_pdf",
			Arguments: map[string]interface{}{
				"content":  "# Custom File Test\n\nMore content.",
				"filename": customFilename,
			},
		},
	}

	result, err := tool.Handler(context.Background(), request)

	// Assertions for custom filename case
	requireHelper.NoError(err, "Handler should not return an error (custom filename)")
	requireHelper.NotNil(result, "Result should not be nil (custom filename)")
	requireHelper.NotEmpty(
		result.Content,
		"Result content slice should not be empty (custom filename)",
	)

	// Check the success message
	textContent, success := mcp.AsTextContent(result.Content[0])
	requireHelper.True(success, "First content should be text content")
	expectedMsg := fmt.Sprintf(
		"PDF successfully saved to %s",
		customFilename,
	)
	requireHelper.Equal(
		expectedMsg,
		textContent.Text,
		"Success message mismatch (custom filename)",
	)

	// Check if the file was created
	_, err = os.Stat(customFilename)
	requireHelper.NoError(err, "Custom output file '%s' should exist", customFilename)

	// Optional: Check file content (basic PDF magic bytes)
	pdfBytes, err := os.ReadFile(customFilename)
	requireHelper.NoError(err, "Failed to read created file %s", customFilename)
	requireHelper.Greater(
		len(pdfBytes),
		4,
		"PDF data too short in file %s",
		customFilename,
	)
	requireHelper.Equal(
		[]byte{0x25, 0x50, 0x44, 0x46, 0x2d},
		pdfBytes[:5],
		"File %s should start with PDF magic bytes",
		customFilename,
	)
}

func TestHandlerMissingContent(t *testing.T) {
	t.Parallel()
	requireHelper := require.New(t)
	logger := log.New(os.Stderr, "[pdf-test] ", 0)

	tool, err := NewPdfTool(logger)
	requireHelper.NoError(err, "NewPdfTool should not return an error")

	// Create a test request with missing content
	request := mcp.CallToolRequest{}
	request.Params.Name = "markdown_to_pdf"
	request.Params.Arguments = map[string]interface{}{} // Missing 'content'

	// Execute the handler
	result, err := tool.Handler(context.Background(), request)

	// Assertions
	requireHelper.Error(err, "Handler should return an error for missing content")
	requireHelper.Nil(result, "Result should be nil on error")
	requireHelper.Contains(err.Error(), "missing required parameter: content")
}
