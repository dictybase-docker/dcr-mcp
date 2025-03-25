package markdown_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/dictybase/dcr-mcp/pkg/markdown"
)

func Example() {
	// Create a new parser with default settings
	parser := markdown.NewParser()

	// Sample markdown content with metadata, GFM features, and code blocks
	markdownContent := `---
title: Example Document
author: John Doe
date: 2025-03-25
---

# Markdown Example

This is an example of **Markdown** with _formatting_.

## GitHub Flavored Markdown Features

* [ ] Task list item 1
* [x] Task list item 2 (completed)

## Table Example

| Name  | Age | Role          |
|-------|-----|---------------|
| Alice | 28  | Developer     |
| Bob   | 32  | Project Lead  |

## Code Example with Syntax Highlighting

` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Hello, Markdown!")
}
` + "```" + `

:smile: Emoji support is included!

Visit [Goldmark](https://github.com/yuin/goldmark) for more info.
`

	// Parse the markdown content
	htmlOutput, err := parser.ParseString(markdownContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing markdown: %v\n", err)
		return
	}

	// Print first few lines as an example
	lines := strings.Split(htmlOutput, "\n")
	for i := 0; i < 5 && i < len(lines); i++ {
		fmt.Println(lines[i])
	}
	fmt.Println("...")

	// Get metadata from the document
	metadata := parser.GetMetadata()
	fmt.Printf("\nMetadata:\n")
	for k, v := range metadata {
		fmt.Printf("  %s: %v\n", k, v)
	}

	// Example with XHTML output
	xhtmlParser := markdown.NewParser(markdown.WithXHTML(), markdown.WithLineNumbers())
	xhtmlOutput, _ := xhtmlParser.ParseString(`<br>`)
	
	// Print XHTML output to show self-closing tags
	fmt.Printf("\nXHTML output (shows self-closing tags):\n%s\n", xhtmlOutput)
}
