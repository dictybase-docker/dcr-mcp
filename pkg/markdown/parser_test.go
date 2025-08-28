package markdown

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

type parserTestCase struct {
	name     string
	markdown string
	want     string
	options  []ParserOption
}

func getParserTestCases() []parserTestCase {
	return []parserTestCase{
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
}

func TestParser(t *testing.T) {
	t.Parallel()
	testCases := getParserTestCases()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			requireHelper := require.New(t)

			markdownParser := NewParser(testCase.options...)
			gotResult, err := markdownParser.ParseString(testCase.markdown)

			requireHelper.NoError(err, "Parser.ParseString() should not return an error")
			requireHelper.Contains(gotResult, testCase.want, "Output should contain expected HTML")
		})
	}
}

func TestParserReader(t *testing.T) {
	t.Parallel()
	requireHelper := require.New(t)
	markdown := "# Heading"
	reader := bytes.NewReader([]byte(markdown))

	markdownParser := NewParser()
	got, err := markdownParser.ParseReader(reader)

	requireHelper.NoError(err, "Parser.ParseReader() should not return an error")
	requireHelper.Contains(string(got), "<h1", "Output should contain h1 heading")
}

func TestMetadata(t *testing.T) {
	t.Parallel()
	requireHelper := require.New(t)
	markdown := `---
title: Test Document
author: John Doe
---
# Content`

	markdownParser := NewParser()
	_, err := markdownParser.ParseString(markdown)

	requireHelper.NoError(err, "Parser.ParseString() should not return an error")

	meta := markdownParser.GetMetadata()
	requireHelper.Equal("Test Document", meta["title"], "Metadata should contain correct title")
	requireHelper.Equal("John Doe", meta["author"], "Metadata should contain correct author")
}
