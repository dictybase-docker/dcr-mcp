# Literature Tool

An MCP tool for fetching scientific literature information using PubMed or DOI IDs via the [dictyBase literature API](https://github.com/dictybase/literature).

## Features

- **Smart Provider Selection**: Automatically chooses the best data source:
  - For DOI: Uses EuropePMC (better DOI support)
  - For PMID: Uses EuropePMC first with PubMed fallback
- **Comprehensive Validation**: Validates and normalizes both PMID and DOI inputs
- **Rich Metadata**: Returns detailed article information including authors, abstracts, citations, MeSH headings, and more
- **Flexible Input**: Handles various ID formats (with/without prefixes, URLs, etc.)
- **Structured Output**: Returns both formatted text and raw JSON data

## Usage

### Basic Examples

```json
{
  "name": "literature-fetch",
  "arguments": {
    "id": "12345678",
    "id_type": "pmid"
  }
}
```

```json
{
  "name": "literature-fetch", 
  "arguments": {
    "id": "10.1038/nature12373",
    "id_type": "doi"
  }
}
```

### With Provider Selection

```json
{
  "name": "literature-fetch",
  "arguments": {
    "id": "PMID:12345678", 
    "id_type": "pmid",
    "provider": "europepmc"
  }
}
```

## Parameters

| Parameter | Type | Required | Description | Valid Values |
|-----------|------|----------|-------------|--------------|
| `id` | string | Yes | The identifier (PMID or DOI) | Any valid PMID or DOI |
| `id_type` | string | Yes | Type of identifier | `"pmid"`, `"doi"` |
| `provider` | string | No | Preferred provider (auto-selected if not specified) | `"pubmed"`, `"europepmc"` |

## Input Normalization

The tool automatically normalizes various input formats:

### PMID Examples
- `12345678` → `12345678`
- `PMID:12345678` → `12345678`
- `pmid:12345678` → `12345678`
- `  12345678  ` → `12345678`

### DOI Examples  
- `10.1038/nature12373` → `10.1038/nature12373`
- `DOI:10.1038/nature12373` → `10.1038/nature12373`
- `https://doi.org/10.1038/nature12373` → `10.1038/nature12373`
- `http://doi.org/10.1038/nature12373` → `10.1038/nature12373`

## Output Format

The tool returns formatted text that includes:

1. **Human-readable summary** with:
   - Title
   - Authors
   - Journal and publication year
   - Abstract
   - PMID/DOI
   - Citation count (if available)

2. **Raw JSON data** with complete metadata including:
   - Author details with affiliations and ORCIDs
   - MeSH headings and qualifiers
   - Chemical substances
   - Grant information
   - Publication dates and revision history

## Error Handling

The tool provides detailed error messages for:
- Missing required parameters
- Invalid ID formats
- Unsupported ID types
- Article not found
- API communication errors

## Implementation Details

### Provider Strategy

1. **For DOI requests**: Uses EuropePMC exclusively (better DOI support)
2. **For PMID requests**: Tries EuropePMC first, falls back to PubMed if needed

### Data Sources

- **PubMed (NCBI eUtils)**: Authoritative biomedical literature database
- **EuropePMC**: Enhanced metadata, citation analytics, European content focus

## Testing

Run the comprehensive test suite:

```bash
gotestsum --format-hide-empty-pkg --format testdox --format-icons hivis -- ./pkg/tools/literaturetool/...
```

Tests cover:
- Input validation and normalization
- Error handling
- Output formatting
- Tool registration and MCP integration