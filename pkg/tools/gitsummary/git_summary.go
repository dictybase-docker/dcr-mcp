package gitsummary

import (
	"context"
	"fmt"
	"log"

	"github.com/dictybase/dcr-mcp/pkg/tools/worksummary"
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
	RepoURL    string `validate:"required"`
	Branch     string `validate:"required"`
	StartDate  string `validate:"required"`
	EndDate    string
	Author     string
	APIKey     string `validate:"required"`
	ModelName  string `validate:"required"`
	OpenAIBase string
}

// GitSummaryResponse represents the response from the git summary request.
type GitSummaryResponse struct {
	Summary string `json:"summary"`
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
			mcp.Description("Filter commits by author name (optional)"),
		),
		mcp.WithString(
			"api_key",
			mcp.Description("OpenAI API key"),
			mcp.Required(),
		),
		mcp.WithString(
			"model_name",
			mcp.Description("OpenAI model name (e.g., 'gpt-4-turbo')"),
			mcp.Required(),
		),
		mcp.WithString(
			"openai_base",
			mcp.Description("Custom OpenAI API base URL (optional)"),
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
func (g *GitSummaryTool) Handler() func(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments
		params := GitSummaryRequest{
			RepoURL:    args["repo_url"].(string),
			Branch:     args["branch"].(string),
			StartDate:  args["start_date"].(string),
			EndDate:    args["end_date"].(string),
			Author:     args["author"].(string),
			APIKey:     args["api_key"].(string),
			ModelName:  args["model_name"].(string),
			OpenAIBase: args["openai_base"].(string),
		}
		if err := validate.Struct(params); err != nil {
			return nil, fmt.Errorf("Validation error: %v", err)
		}
		openaiConfig := worksummary.OpenAIConfig{
			APIKey:  params.APIKey,
			BaseURL: params.OpenAIBase,
			Model:   params.ModelName,
		}

		client, err := worksummary.NewOpenAIClient(openaiConfig)
		if err != nil {
			return nil, fmt.Errorf("Error initializing OpenAI client: %v", err)
		}
		summary, err := g.GenerateSummary(ctx, client, params)
		if err != nil {
			return nil, fmt.Errorf("Error generating summary: %v", err)
		}

		return mcp.NewToolResultText(summary), nil
	}
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
