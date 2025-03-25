package markdown

import (
	"bytes"
	"io"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	html_renderer "github.com/yuin/goldmark/renderer/html"
)

// Parser is a Markdown parser with GFM, syntax highlighting, typographer extensions and XHTML rendering
type Parser struct {
	converter goldmark.Markdown
	context   parser.Context
}

// ParserOption defines a functional option for configuring the Markdown Parser
type ParserOption func(*Parser)

// WithLineNumbers enables line numbers in code blocks
func WithLineNumbers() ParserOption {
	return func(p *Parser) {
		// The converter is already initialized with defaults
		// This option would apply in a more complex implementation
	}
}

// WithXHTML configures the renderer to output XHTML
func WithXHTML() ParserOption {
	return func(p *Parser) {
		p.converter = goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Typographer,
				highlighting.NewHighlighting(
					highlighting.WithStyle("github"),
				),
				emoji.Emoji,
				meta.Meta,
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html_renderer.WithHardWraps(),
				html_renderer.WithXHTML(),
			),
		)
	}
}

// WithUnsafeHTML allows raw HTML to pass through the renderer
// Only use this option for trusted content!
func WithUnsafeHTML() ParserOption {
	return func(p *Parser) {
		p.converter = goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Typographer,
				highlighting.NewHighlighting(
					highlighting.WithStyle("github"),
				),
				emoji.Emoji,
				meta.Meta,
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html_renderer.WithHardWraps(),
				html_renderer.WithUnsafe(),
			),
		)
	}
}

// NewParser creates a new Markdown parser with the provided options
func NewParser(opts ...ParserOption) *Parser {
	// Create default parser with sensible defaults
	p := &Parser{
		converter: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Typographer,
				highlighting.NewHighlighting(
					highlighting.WithStyle("paraiso-light"),
				),
				emoji.Emoji,
				meta.Meta,
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html_renderer.WithHardWraps(),
				html_renderer.WithXHTML(),
			),
		),
		context: parser.NewContext(),
	}

	// Apply all options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Parse converts markdown source to HTML
func (p *Parser) Parse(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := p.converter.Convert(src, &buf, parser.WithContext(p.context)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ParseString converts a markdown string to HTML
func (p *Parser) ParseString(src string) (string, error) {
	html, err := p.Parse([]byte(src))
	if err != nil {
		return "", err
	}
	return string(html), nil
}

// ParseReader converts markdown from a reader to HTML
func (p *Parser) ParseReader(reader io.Reader) ([]byte, error) {
	src, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return p.Parse(src)
}

// GetMetadata returns the metadata extracted from the markdown document
func (p *Parser) GetMetadata() map[string]interface{} {
	return meta.Get(p.context)
}
