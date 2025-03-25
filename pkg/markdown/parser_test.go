package markdown

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
		options  []ParserOption
	}{
		{
			name:     "basic formatting",
			markdown: "**Bold** and *italic*",
			want:     "<p><strong>Bold</strong> and <em>italic</em></p>",
			options:  nil,
		},
		{
			name:     "code highlighting",
			markdown: "```go\nfunc main() {}\n```",
			want:     "<pre",
			options:  nil,
		},
		{
			name:     "gfm tables",
			markdown: "| a | b |\n|---|---|\n| 1 | 2 |",
			want:     "<table>",
			options:  nil,
		},
		{
			name:     "emoji",
			markdown: ":smile:",
			want:     "&#x1f604;", // HTML entity for smile emoji
			options:  nil,
		},
		{
			name:     "xhtml output",
			markdown: "<br>",
			want:     "<!-- raw HTML omitted -->", // Goldmark sanitizes raw HTML by default
			options:  []ParserOption{WithXHTML()},
		},
		{
			name:     "unsafe html",
			markdown: "<div class=\"custom\">Raw HTML</div>",
			want:     "<div class=\"custom\">Raw HTML</div>",
			options:  []ParserOption{WithUnsafeHTML()},
		},
		{
			name:     "metadata extraction",
			markdown: "---\ntitle: Test\n---\n# Content",
			want:     "<h1",
			options:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			
			p := NewParser(tt.options...)
			got, err := p.ParseString(tt.markdown)
			
			r.NoError(err, "Parser.ParseString() should not return an error")
			r.Contains(got, tt.want, "Output should contain expected HTML")
		})
	}
}

func TestParserReader(t *testing.T) {
	r := require.New(t)
	markdown := "# Heading"
	reader := bytes.NewReader([]byte(markdown))
	
	p := NewParser()
	got, err := p.ParseReader(reader)
	
	r.NoError(err, "Parser.ParseReader() should not return an error")
	r.Contains(string(got), "<h1", "Output should contain h1 heading")
}

func TestMetadata(t *testing.T) {
	r := require.New(t)
	markdown := `---
title: Test Document
author: John Doe
---
# Content`

	p := NewParser()
	_, err := p.ParseString(markdown)
	
	r.NoError(err, "Parser.ParseString() should not return an error")
	
	meta := p.GetMetadata()
	r.Equal("Test Document", meta["title"], "Metadata should contain correct title")
	r.Equal("John Doe", meta["author"], "Metadata should contain correct author")
}
