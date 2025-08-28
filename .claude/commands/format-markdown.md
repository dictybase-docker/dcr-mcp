---
description: Format and fix markdown syntax in files or content
category: workflow
allowed-tools: Read, Write, Edit
argument-hint: "[file-path or content]"
model: sonnet
---

# Format Markdown

Fix markdown syntax and ensure consistent document structure following CommonMark and GitHub Flavored Markdown specifications.

## Usage

You can format markdown in two ways:
1. **File formatting**: `/format-markdown path/to/file.md`
2. **Content formatting**: `/format-markdown "# My content to format"`

## What gets formatted

### Document Structure
- Proper heading hierarchy (# for H1, ## for H2, etc.)
- Consistent heading spacing (blank lines before/after)
- Logical progression without skipping levels

### Lists
- Consistent list markers (- for unordered lists)
- Proper indentation (2 spaces for nested items)
- Blank lines before and after list blocks
- Convert numbered sequences to ordered lists

### Code Formatting
- Triple backticks (```) for multi-line code blocks
- Language identifiers when detectable (```python, ```javascript)
- Single backticks for inline code
- Preserved indentation within blocks

### Text Emphasis
- **Double asterisks** for bold text
- *Single asterisks* for italic text
- `Backticks` for code and technical terms
- Proper link formatting [text](url)

### Content Preservation
- Maintains original document flow and structure
- Keeps all content intact while improving formatting
- Respects existing correct markdown
- Adds horizontal rules (---) for major section breaks

## Process

For file paths:
1. Read the specified markdown file
2. Analyze structure and identify formatting issues
3. Apply consistent markdown syntax
4. Write the improved version back to the file

For content:
1. Process the provided text content
2. Apply markdown formatting rules
3. Return the formatted content

## Examples

```bash
# Format a specific file
/format-markdown README.md

# Format inline content
/format-markdown "# HEADING\n- bullet point\n- another point"

# Format documentation files
/format-markdown docs/api.md
```

The formatting ensures your markdown renders correctly in any standard parser while maintaining readability and proper document structure.