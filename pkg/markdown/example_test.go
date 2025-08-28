package markdown_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/dictybase/dcr-mcp/pkg/markdown"
)

func getSampleMarkdownContent() string {
	return `---
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
}

func printSampleOutput(htmlOutput string) {
	// Print first few lines as an example
	lines := strings.Split(htmlOutput, "\n")
	for i := 0; i < 5 && i < len(lines); i++ {
		fmt.Println(lines[i])
	}
	fmt.Println("...")
}

func printMetadata(markdownParser *markdown.Parser) {
	// Get metadata from the document
	metadata := markdownParser.GetMetadata()
	fmt.Printf("\nMetadata:\n")
	for key, value := range metadata {
		fmt.Printf("  %s: %v\n", key, value)
	}
}

func demonstrateXHTMLOutput() {
	// Example with XHTML output
	xhtmlParser := markdown.NewParser(markdown.WithXHTML(), markdown.WithLineNumbers())
	xhtmlOutput, _ := xhtmlParser.ParseString(`<br>`)

	// Print XHTML output to show self-closing tags
	fmt.Printf("\nXHTML output (shows self-closing tags):\n%s\n", xhtmlOutput)
}

func Example() {
	// Create a new parser with default settings
	markdownParser := markdown.NewParser()

	// Parse the markdown content
	htmlOutput, err := markdownParser.ParseString(getSampleMarkdownContent())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing markdown: %v\n", err)
		return
	}

	printSampleOutput(htmlOutput)
	printMetadata(markdownParser)
	demonstrateXHTMLOutput()
}
