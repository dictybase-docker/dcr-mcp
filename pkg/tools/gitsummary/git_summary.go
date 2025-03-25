package gitsummary

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dictybase/dcr-mcp/pkg/worksummary"
	"github.com/go-playground/validator/v10"
	"github.com/mark3labs/mcp-go/mcp"
)

// Initialize validator
var validate = validator.New()

// GitSummaryTool is a tool that summarizes git commit messages.
type GitSummaryTool struct {
	Name        string
	Description string
	Tool        mcp.Tool
	analyzer    *worksummary.GitAnalyzer
	Logger      *log.Logger
}

// GitSummaryRequest represents the parameters for the git summary request.
type GitSummaryRequest struct {
	RepoURL   string `validate:"required"`
	Branch    string `validate:"required"`
	StartDate string `validate:"required"`
	EndDate   string
	Author    string `validate:"required"`
	APIKey    string `validate:"required"`
}

// NewGitSummaryTool creates a new GitSummaryTool instance.
func NewGitSummaryTool(logger *log.Logger) (*GitSummaryTool, error) {
	// Create the tool with proper schema
	tool := mcp.NewTool(
		"git-summary",
		mcp.WithDescription(
			"Summarizes git commit messages within a date range using OpenAI",
		),
		mcp.WithString(
			"repo_url",
			mcp.Description("The URL of the git repository"),
			mcp.Required(),
		),
		mcp.WithString(
			"branch",
			mcp.Description("The branch to analyze"),
			mcp.Required(),
		),
		mcp.WithString(
			"start_date",
			mcp.Description("The start date for commit analysis"),
			mcp.Required(),
		),
		mcp.WithString(
			"end_date",
			mcp.Description(
				"The end date for commit analysis (optional, defaults to today)",
			),
		),
		mcp.WithString(
			"author",
			mcp.Description("Filter commits by author name"),
			mcp.Required(),
		),
		mcp.WithString(
			"api_key",
			mcp.Description(
				"OpenAI API key (optional, defaults to OPENAI_API_KEY environment variable)",
			),
		),
	)

	analyzer := worksummary.NewGitAnalyzer(
		worksummary.WithLogger(logger),
	)

	return &GitSummaryTool{
		Name:        "git-summary",
		Description: "Summarizes git commit messages within a date range using OpenAI",
		Tool:        tool,
		analyzer:    analyzer,
		Logger:      logger,
	}, nil
}

// GetName returns the name of the tool
func (g *GitSummaryTool) GetName() string {
	return g.Name
}

// GetDescription returns the description of the tool
func (g *GitSummaryTool) GetDescription() string {
	return g.Description
}

// GetSchema returns the JSON schema for the tool's parameters
func (g *GitSummaryTool) GetSchema() mcp.ToolInputSchema {
	return g.Tool.InputSchema
}

// GetTool returns the MCP Tool
func (g *GitSummaryTool) GetTool() mcp.Tool {
	return g.Tool
}

// Handler returns a function that handles tool execution requests
func (g *GitSummaryTool) Handler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	
	// Create request with required parameters
	params := GitSummaryRequest{
		RepoURL:   args["repo_url"].(string),
		Branch:    args["branch"].(string),
		StartDate: args["start_date"].(string),
		Author:    args["author"].(string),
		APIKey:    os.Getenv("OPENAI_API_KEY"),
	}
	
	// Only add end_date if it was provided in the arguments
	if endDate, ok := args["end_date"].(string); ok && endDate != "" {
		params.EndDate = endDate
	}
	if err := validate.Struct(params); err != nil {
		return nil, fmt.Errorf("Validation error: %v", err)
	}

	client, err := worksummary.NewOpenAIClient(params.APIKey)
	if err != nil {
		return nil, fmt.Errorf("Error initializing OpenAI client: %v", err)
	}
	summary, err := g.GenerateSummary(ctx, client, params)
	if err != nil {
		return nil, fmt.Errorf("Error generating summary: %v", err)
	}

	return mcp.NewToolResultText(summary), nil
}

// GenerateSummary generates a summary of git commit messages.
func (g *GitSummaryTool) GenerateSummary(
	ctx context.Context,
	client *worksummary.OpenAIClient,
	req GitSummaryRequest,
) (string, error) {
	// Clone the repository
	repo, err := g.analyzer.CloneAndCheckout(ctx, req.RepoURL, req.Branch)
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	// Parse dates
	startDate, endDate, err := g.analyzer.ParseAnalysisDates(
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		return "", fmt.Errorf("failed to parse dates: %w", err)
	}

	// Create commit range parameters
	params := worksummary.CommitRangeParams{
		Repo:   repo,
		Start:  startDate.Time,
		End:    endDate.Time,
		Author: req.Author,
	}

	// Get commit messages
	commitMsgs, err := g.analyzer.ListCommitsInRange(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to list commits: %w", err)
	}

	// No commits found
	if commitMsgs == "" {
		return "No commits found in the specified date range.", nil
	}

	// Generate summary using OpenAI
	summary, err := client.SummarizeCommitMessages(ctx, commitMsgs)
	if err != nil {
		return "", fmt.Errorf("failed to summarize commit messages: %w", err)
	}

	return summary, nil
}
