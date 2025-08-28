package pdftool

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"

	// Add this line
	"github.com/mark3labs/mcp-go/mcp"
	pdf "github.com/stephenafamo/goldmark-pdf" // pdf renderer
	"github.com/yuin/goldmark"
)

// PdfTool is a tool that converts markdown to PDF.
type PdfTool struct {
	Name        string
	Description string
	Tool        mcp.Tool
	Logger      *log.Logger
}

// NewPdfTool creates a new PdfTool instance.
func NewPdfTool(logger *log.Logger) (*PdfTool, error) {
	// Create the tool with proper schema
	// Create the tool with proper schema
	tool := mcp.NewTool(
		"markdown_to_pdf",
		mcp.WithDescription(
			"Converts markdown content to a PDF document and saves it to a file.", // Updated description
		),
		mcp.WithString(
			"content",
			mcp.Description("The markdown content to convert to PDF"),
			mcp.Required(),
		),
		// Add optional filename parameter
		mcp.WithString( // Add this block
			"filename",
			mcp.Description(
				"Optional filename for the output PDF. Defaults to 'output.pdf'.",
			),
			// Not required
		),
	)
	return &PdfTool{
		Name:        "markdown_to_pdf",
		Description: "Converts markdown content to a PDF document and saves it to a file.", // Updated description
		Tool:        tool,
		Logger:      logger,
	}, nil
}

// GetName returns the name of the tool
func (pt *PdfTool) GetName() string {
	return pt.Name
}

// GetDescription returns the description of the tool
func (pt *PdfTool) GetDescription() string {
	return pt.Description
}

// GetSchema returns the JSON schema for the tool's parameters
func (pt *PdfTool) GetSchema() mcp.ToolInputSchema {
	return pt.Tool.InputSchema
}

// GetTool returns the MCP Tool
func (pt *PdfTool) GetTool() mcp.Tool {
	return pt.Tool
}

// Handler returns a function that handles tool execution requests
func (pt *PdfTool) Handler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	contentVal, ok := args["content"].(string)
	if !ok {
		return nil, errors.New("missing required parameter: content")
	}
	// --- Determine output filename ---
	outputFilename := "output.pdf" // Default filename
	if fname, ok := args["filename"].(string); ok &&
		fname != "" {
		outputFilename = fname
	}
	pdfFile, err := os.Create(outputFilename)
	if err != nil {
		return nil, fmt.Errorf(
			"error creating file %s %w", outputFilename, err,
		)
	}
	defer pdfFile.Close()

	markdown := goldmark.New(
		goldmark.WithRenderer(pdf.New(
			pdf.WithContext(
				context.Background(),
			),
			pdf.WithLinkColor(
				color.RGBA{R: 204, G: 69, B: 120, A: 255},
			),
			pdf.WithImageFS(
				http.FS(os.DirFS(".")),
			), // Consider security implications of reading local files
			pdf.WithHeadingFont(
				pdf.GetTextFont(
					"IBM Plex Serif", pdf.FontLora,
				),
			),
			pdf.WithBodyFont(
				pdf.GetTextFont("Open Sans", pdf.FontRoboto)),
			pdf.WithCodeFont(
				pdf.GetCodeFont("Inconsolata", pdf.FontRobotoMono),
			),
		)),
	)
	err = markdown.Convert([]byte(contentVal), pdfFile)
	if err != nil {
		pt.Logger.Printf("Error converting markdown to PDF: %v", err)
		return nil, fmt.Errorf("failed to convert markdown to PDF: %w", err)
	}
	pt.Logger.Println(
		"Successfully converted markdown to PDF",
	)
	pt.Logger.Printf("Saved PDF to file: %s", outputFilename)
	return mcp.NewToolResultText(
		fmt.Sprintf("PDF successfully saved to %s", outputFilename),
	), nil
}
