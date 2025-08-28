package pdftool

import (
	"context"
	"fmt" // Add this import
	"log"
	"os"
	"path/filepath" // Add this import
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
	_, ok := schema.Properties["content"] // Get content prop details
	r.True(ok, "Schema should have 'content' property")

	// Check for optional filename parameter
	r.Contains(
		schema.Properties,
		"filename",
	) // Add this line
	filenameProp, ok := schema.Properties["filename"] // Add this line
	r.True(
		ok,
		"Schema should have 'filename' property",
	) // Add this line
	r.Equal(
		"string",
		filenameProp.Type,
	) // Add this line
	r.Equal(
		"Optional filename for the output PDF. Defaults to 'output.pdf'.",
		filenameProp.Description,
	) // Add this line

	r.Contains(schema.Required, "content")
	r.NotContains(
		schema.Required,
		"filename",
	) // Ensure filename is not required
}

func TestHandler(t *testing.T) {
	r := require.New(t)
	// Use a logger that writes to stderr for visibility during tests
	logger := log.New(os.Stderr, "[pdf-test-handler] ", log.LstdFlags)

	tool, err := NewPdfTool(logger)
	r.NoError(err, "NewPdfTool should not return an error")

	// --- Test Case 1: Default filename ---
	t.Run("DefaultFilename", func(t *testing.T) {
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
		r.NoError(err, "Handler should not return an error (default filename)")
		r.NotNil(result, "Result should not be nil (default filename)")
		r.NotEmpty(
			result.Content,
			"Result content slice should not be empty (default filename)",
		)

		// Check the success message
		textContent := result.Content[0].Text
		expectedMsg := fmt.Sprintf(
			"PDF successfully saved to %s",
			defaultFilename,
		)
		r.Equal(
			expectedMsg,
			textContent,
			"Success message mismatch (default filename)",
		)

		// Check if the file was created
		_, err = os.Stat(defaultFilename)
		r.NoError(err, "Default output file '%s' should exist", defaultFilename)

		// Optional: Check file content (basic PDF magic bytes)
		pdfBytes, err := os.ReadFile(defaultFilename)
		r.NoError(err, "Failed to read created file %s", defaultFilename)
		r.True(
			len(pdfBytes) > 4,
			"PDF data too short in file %s",
			defaultFilename,
		)
		r.Equal(
			[]byte{0x25, 0x50, 0x44, 0x46, 0x2d},
			pdfBytes[:5],
			"File %s should start with PDF magic bytes",
			defaultFilename,
		)
	})

	// --- Test Case 2: Custom filename ---
	t.Run("CustomFilename", func(t *testing.T) {
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
		r.NoError(err, "Handler should not return an error (custom filename)")
		r.NotNil(result, "Result should not be nil (custom filename)")
		r.NotEmpty(
			result.Content,
			"Result content slice should not be empty (custom filename)",
		)

		// Check the success message
		textContent := result.Content[0].Text
		expectedMsg := fmt.Sprintf(
			"PDF successfully saved to %s",
			customFilename,
		)
		r.Equal(
			expectedMsg,
			textContent,
			"Success message mismatch (custom filename)",
		)

		// Check if the file was created
		_, err = os.Stat(customFilename)
		r.NoError(err, "Custom output file '%s' should exist", customFilename)

		// Optional: Check file content (basic PDF magic bytes)
		pdfBytes, err := os.ReadFile(customFilename)
		r.NoError(err, "Failed to read created file %s", customFilename)
		r.True(
			len(pdfBytes) > 4,
			"PDF data too short in file %s",
			customFilename,
		)
		r.Equal(
			[]byte{0x25, 0x50, 0x44, 0x46, 0x2d},
			pdfBytes[:5],
			"File %s should start with PDF magic bytes",
			customFilename,
		)
	})
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
