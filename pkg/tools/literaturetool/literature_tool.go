package literaturetool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mark3labs/mcp-go/mcp"
)

// Initialize validator.
var validate = validator.New()

// DOI regex pattern to match and extract DOI from various formats.
// Handles optional prefixes: doi:, DOI:, https://doi.org/, http://doi.org/
// Captures the actual DOI part (10.xxxx/yyyy) with whitespace trimming.
var doiRegex = regexp.MustCompile(
	`(?i)^(?:(?:https?://)?doi\.org/|doi:)?\s*(10\.\S+/\S+)\s*$`,
)

// PMID regex pattern to validate and extract PMID (positive integers only).
var pmidRegex = regexp.MustCompile(`^\d+$`)

// LiteratureTool is a tool that fetches literature information using PubMed or DOI IDs.
type LiteratureTool struct {
	Name        string
	Description string
	Tool        mcp.Tool
	client      *LiteratureClient
	Logger      *log.Logger
}

// LiteratureRequest represents the parameters for the literature fetch request.
type LiteratureRequest struct {
	ID       string `validate:"required"                         json:"id"`
	IDType   string `validate:"required,oneof=pmid doi"          json:"id_type"`
	Provider string `validate:"omitempty,oneof=pubmed europepmc" json:"provider"`
}

// fetchArticle retrieves article information using the recommended strategy:
// - For DOI: Try EuropePMC
// - For PMID: Try EuropePMC first, fallback to NCBI/PubMed.
func (l *LiteratureTool) fetchArticle(
	ctx context.Context,
	params LiteratureRequest,
) (*Article, error) {
	if params.IDType == IDTypeDOI {
		// For DOI, only use EuropePMC as it has better DOI support
		l.Logger.Printf(
			"Fetching article for DOI %s using EuropePMC",
			params.ID,
		)
		return l.client.GetArticleFromEuropePMC(ctx, params.ID, params.IDType)
	}

	// For PMID, use EuropePMC first with PubMed fallback
	l.Logger.Printf(
		"Fetching article for PMID %s using EuropePMC with PubMed fallback",
		params.ID,
	)
	return l.client.GetArticleWithFallback(ctx, params.ID, params.IDType)
}

// NewLiteratureTool creates a new LiteratureTool instance.
func NewLiteratureTool(logger *log.Logger) (*LiteratureTool, error) {
	// Create the tool with proper schema
	tool := mcp.NewTool(
		"literature-fetch",
		mcp.WithDescription(
			"Fetches scientific literature information using PubMed or DOI IDs via the dictyBase literature API",
		),
		mcp.WithString(
			"id",
			mcp.Description("The PubMed ID (PMID) or DOI identifier"),
			mcp.Required(),
		),
		mcp.WithString(
			"id_type",
			mcp.Description(
				"Type of identifier: 'pmid' for PubMed IDs or 'doi' for DOI",
			),
			mcp.Required(),
			mcp.Enum("pmid", "doi"),
		),
		mcp.WithString(
			"provider",
			mcp.Description(
				"Literature provider: 'pubmed' (default) or 'europepmc' for enhanced metadata",
			),
			mcp.Enum("pubmed", "europepmc"),
		),
	)

	client, err := NewLiteratureClient(WithLogger(logger))
	if err != nil {
		return nil, fmt.Errorf("failed to create literature client: %w", err)
	}

	return &LiteratureTool{
		Name:        "literature-fetch",
		Description: "Fetches scientific literature information using PubMed or DOI IDs via the dictyBase literature API",
		Tool:        tool,
		client:      client,
		Logger:      logger,
	}, nil
}

// GetName returns the name of the tool.
func (l *LiteratureTool) GetName() string {
	return l.Name
}

// GetDescription returns the description of the tool.
func (l *LiteratureTool) GetDescription() string {
	return l.Description
}

// GetSchema returns the JSON schema for the tool's parameters.
func (l *LiteratureTool) GetSchema() mcp.ToolInputSchema {
	return l.Tool.InputSchema
}

// GetTool returns the MCP Tool.
func (l *LiteratureTool) GetTool() mcp.Tool {
	return l.Tool
}

// Handler returns a function that handles tool execution requests.
func (l *LiteratureTool) Handler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	// Create request with required parameters
	identifier, idOk := args["id"].(string)
	idType, idTypeOk := args["id_type"].(string)

	if !idOk || !idTypeOk {
		return nil, fmt.Errorf("missing required parameters: id and id_type")
	}

	params := LiteratureRequest{
		ID:     identifier,
		IDType: idType,
	}

	// Set default provider if not specified
	if provider, ok := args["provider"].(string); ok && provider != "" {
		params.Provider = provider
	} else {
		params.Provider = "pubmed" // Default to PubMed
	}

	// Validate parameters
	if err := validate.Struct(params); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Normalize ID based on type
	normalizedID, err := l.normalizeID(params.ID, params.IDType)
	if err != nil {
		return nil, fmt.Errorf("invalid %s format: %w", params.IDType, err)
	}
	params.ID = normalizedID

	// Fetch literature information
	article, err := l.fetchArticle(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch literature: %w", err)
	}

	// Format and return the result
	result, err := l.formatArticleResult(article)
	if err != nil {
		return nil, fmt.Errorf("failed to format result: %w", err)
	}

	return mcp.NewToolResultText(result), nil
}

// normalizeID validates and normalizes the identifier based on its type.
func (l *LiteratureTool) normalizeID(id, idType string) (string, error) {
	switch idType {
	case IDTypePMID:
		return l.normalizePMID(id)
	case IDTypeDOI:
		return l.normalizeDOI(id)
	default:
		return "", fmt.Errorf("unsupported ID type: %s", idType)
	}
}

// normalizePMID validates and normalizes a PubMed ID.
func (l *LiteratureTool) normalizePMID(pmid string) (string, error) {
	pid := strings.TrimSpace(pmid)
	if len(pid) == 0 {
		return "", fmt.Errorf("PMID cannot be empty")
	}
	if !pmidRegex.MatchString(pid) {
		return "", fmt.Errorf(
			"PMID must contain only digits, got: %s",
			pmid,
		)
	}

	return pid, nil
}

// normalizeDOI validates and normalizes a DOI.
func (l *LiteratureTool) normalizeDOI(doi string) (string, error) {
	// Use regex to match and extract the DOI from various formats
	matches := doiRegex.FindStringSubmatch(doi)
	if len(matches) < 2 {
		return "", fmt.Errorf(
			"invalid DOI format, expected '10.xxxx/yyyy', got: %s",
			doi,
		)
	}

	// The captured group contains the normalized DOI
	normalizedDOI := matches[1]

	// Validate that the suffix after "/" is not empty
	parts := strings.SplitN(normalizedDOI, "/", 2)
	if len(parts) != 2 || parts[1] == "" {
		return "", fmt.Errorf(
			"invalid DOI format, expected '10.xxxx/yyyy', got: %s",
			doi,
		)
	}

	return normalizedDOI, nil
}

// formatArticleResult formats the article information for display.
func (l *LiteratureTool) formatArticleResult(article *Article) (string, error) {
	if article == nil {
		return "No article found", nil
	}

	jsonData, err := json.MarshalIndent(article, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal article data: %w", err)
	}

	var result strings.Builder
	result.WriteString("## Literature Information\n\n")

	l.formatBasicInfo(&result, article)
	l.formatMetadata(&result, article)
	l.formatJSONData(&result, jsonData)

	return result.String(), nil
}

// formatBasicInfo formats title, authors, and journal information.
func (l *LiteratureTool) formatBasicInfo(result *strings.Builder, article *Article) {
	if article.Title != "" {
		fmt.Fprintf(result, "**Title:** %s\n\n", article.Title)
	}

	if len(article.Authors) > 0 {
		result.WriteString("**Authors:** ")
		for index, author := range article.Authors {
			if index > 0 {
				result.WriteString(", ")
			}
			result.WriteString(author.FullName)
		}
		result.WriteString("\n\n")
	}

	if article.Journal.Title != "" {
		fmt.Fprintf(result, "**Journal:** %s", article.Journal.Title)
		if article.PubYear != "" {
			fmt.Fprintf(result, " (%s)", article.PubYear)
		}
		result.WriteString("\n\n")
	}

	if article.Abstract != "" {
		fmt.Fprintf(result, "**Abstract:** %s\n\n", article.Abstract)
	}
}

// formatMetadata formats PMID, DOI, and citation information.
func (l *LiteratureTool) formatMetadata(result *strings.Builder, article *Article) {
	if article.PMID != "" {
		fmt.Fprintf(result, "**PMID:** %s\n", article.PMID)
	}

	if article.DOI != "" {
		fmt.Fprintf(result, "**DOI:** %s\n", article.DOI)
	}

	if article.CitedByCount > 0 {
		fmt.Fprintf(result, "**Citations:** %d\n", article.CitedByCount)
	}
}

// formatJSONData appends the raw JSON data section.
func (l *LiteratureTool) formatJSONData(result *strings.Builder, jsonData []byte) {
	result.WriteString("\n---\n\n")
	result.WriteString("**Raw JSON Data:**\n```json\n")
	result.WriteString(string(jsonData))
	result.WriteString("\n```")
}
