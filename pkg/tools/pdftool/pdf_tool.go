package pdftool

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/color"
	"log"

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
	tool := mcp.NewTool(
		"markdown_to_pdf",
		mcp.WithDescription(
			"Converts markdown content to a PDF document.",
		),
		mcp.WithString(
			"content",
			mcp.Description("The markdown content to convert to PDF"),
			mcp.Required(),
		),
	)
	return &PdfTool{
		Name:        "markdown_to_pdf",
		Description: "Converts markdown content to a PDF document.",
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
	contentVal, ok := request.Params.Arguments["content"].(string)
	if !ok {
		return nil, errors.New("missing required parameter: content")
	}
	markdown := goldmark.New(
		// Use the PDF renderer instead of the HTML one
		goldmark.WithRenderer(
			pdf.New(
				pdf.WithContext(context.Background()),
				pdf.WithLinkColor(color.RGBA{204, 69, 120, 255}),
				pdf.WithHeadingFont(
					pdf.GetTextFont("IBM Plex Serif", pdf.FontLora),
				),
				pdf.WithBodyFont(pdf.GetTextFont("Open Sans", pdf.FontRoboto)),
				pdf.WithCodeFont(
					pdf.GetCodeFont("Inconsolata", pdf.FontRobotoMono),
				),
			),
		),
	)

	var pdfBuffer bytes.Buffer
	err := markdown.Convert([]byte(contentVal), &pdfBuffer)
	if err != nil {
		pt.Logger.Printf("Error converting markdown to PDF: %v", err)
		return nil, fmt.Errorf("failed to convert markdown to PDF: %w", err)
	}
	pt.Logger.Printf(
		"Successfully converted markdown to PDF (%d bytes)",
		pdfBuffer.Len(),
	)
	return mcp.NewToolResultText(
		base64.StdEncoding.EncodeToString(pdfBuffer.Bytes()),
	), nil
}
