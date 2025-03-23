package worksummary

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	dps "github.com/markusmobius/go-dateparser"
	"github.com/markusmobius/go-dateparser/date"
	"github.com/sashabaranov/go-openai"
)

const (
	GitSummaryPrompt = `
    You are an expert in summarizing git commit messages. You will be given a
	collection of git commit messages that you will summarize by creating
	not more than four focused bullet points. Each bullet point should:
    1. Begin with a bold category that reflects the theme of the changes (like
       "**User Interface**" or "**Performance**")
    2. Contain multiple sentences that explain what was changed in plain language
    3. Avoid technical jargon when possible, or explain technical terms when they must be used
    4. Focus on the business value and user impact rather than implementation details

    Present the output in markdown format, with "Work Summary" as the main
	heading (H1). The summary should be easily understood by someone without
	technical background, focusing on what was accomplished rather than how
	it was done.
    `
)

// GitAnalyzer handles git repository analysis operations including cloning
// repositories, parsing dates, and retrieving commit histories within specified
// date ranges.
type GitAnalyzer struct {
	logger     *log.Logger
	dateConfig *dps.Configuration
}

// CommitRangeParams holds parameters for listing commits in a date range.
type CommitRangeParams struct {
	Repo   *git.Repository
	Start  time.Time
	End    time.Time
	Author string
}

// GitAnalyzerOption defines a functional option for configuring GitAnalyzer.
type GitAnalyzerOption func(*GitAnalyzer)

// WithLogger sets a custom logger for GitAnalyzer.
func WithLogger(logger *log.Logger) GitAnalyzerOption {
	return func(ga *GitAnalyzer) {
		ga.logger = logger
	}
}

// WithCurrentTime sets a custom current time for date parsing.
func WithCurrentTime(t time.Time) GitAnalyzerOption {
	return func(ga *GitAnalyzer) {
		ga.dateConfig.CurrentTime = t
	}
}

// WithTimeZone sets a custom timezone for date parsing.
func WithTimeZone(tz *time.Location) GitAnalyzerOption {
	return func(ga *GitAnalyzer) {
		ga.dateConfig.DefaultTimezone = tz
	}
}

// NewGitAnalyzer creates a new GitAnalyzer with the provided options.
func NewGitAnalyzer(opts ...GitAnalyzerOption) *GitAnalyzer {
	ga := &GitAnalyzer{
		logger: log.New(
			os.Stderr,
			"[git-commit-summary] ",
			log.LstdFlags|log.Lmsgprefix,
		),
		dateConfig: &dps.Configuration{
			DefaultTimezone: time.Local,
			CurrentTime:     time.Now(),
		},
	}

	// Apply all options
	for _, opt := range opts {
		opt(ga)
	}

	return ga
}

// SummaryClient is the interface for clients that can generate summaries.
type SummaryClient interface {
	SummarizeCommitMessages(ctx context.Context, commitMsgs string) (string, error)
}

// OpenAIClient implements SummaryClient using OpenAI API.
type OpenAIClient struct {
	client *openai.Client
	model  string
}

// NewOpenAIClient creates a new OpenAI client with the provided configuration.
func NewOpenAIClient(apiKey, baseURL, model string) *OpenAIClient {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	return &OpenAIClient{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

// SummarizeCommitMessages generates a summary of commit messages using OpenAI.
func (c *OpenAIClient) SummarizeCommitMessages(ctx context.Context, commitMsgs string) (string, error) {
	if c.client == nil {
		return "", errors.New("OpenAI client not configured")
	}

	req := openai.ChatCompletionRequest{
		Model:       c.model,
		Stream:      true,
		Temperature: 0.1, // Controls randomness in the response
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: GitSummaryPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: commitMsgs,
			},
		},
	}

	var sb strings.Builder
	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", fmt.Errorf("OpenAI stream error: %w", err)
	}
	defer stream.Close()

	for {
		select {
		case <-ctx.Done():
			return sb.String(), ctx.Err()
		default:
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return sb.String(), nil
			}
			if err != nil {
				return sb.String(), fmt.Errorf("OpenAI stream recv error: %w", err)
			}
			sb.WriteString(resp.Choices[0].Delta.Content)
		}
	}
}

func (ga *GitAnalyzer) parseStartDate(dateStr string) (date.Date, error) {
	parsedDate, err := dps.Parse(ga.dateConfig, dateStr)
	if err != nil || parsedDate.Time.IsZero() {
		return parsedDate, fmt.Errorf("could not parse date '%s'", dateStr)
	}
	return parsedDate, nil
}

func (ga *GitAnalyzer) parseEndDate(endDateStr string) (date.Date, error) {
	if len(endDateStr) == 0 {
		return date.Date{Time: ga.dateConfig.CurrentTime, Period: date.Day}, nil
	}
	parsedDate, err := dps.Parse(ga.dateConfig, endDateStr)
	if err != nil || parsedDate.Time.IsZero() {
		return parsedDate, fmt.Errorf("could not parse date '%s'", endDateStr)
	}
	return parsedDate, nil
}

// ParseAnalysisDates parses start and end date strings into date.Date objects.
func (ga *GitAnalyzer) ParseAnalysisDates(startDate, endDate string) (date.Date, date.Date, error) {
	start, err := ga.parseStartDate(startDate)
	if err != nil {
		return start, date.Date{}, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := ga.parseEndDate(endDate)
	if err != nil {
		return start, end, fmt.Errorf("invalid end date: %w", err)
	}
	return start, end, nil
}

// CloneAndCheckout clones a repository and checks out the specified branch.
func (ga *GitAnalyzer) CloneAndCheckout(
	ctx context.Context, repoURL, branchName string,
) (*git.Repository, error) {
	ga.logger.Printf("Analyzing repository: %s", repoURL)
	ga.logger.Printf("Cloning branch: %s", branchName)

	repo, err := git.CloneContext(
		ctx,
		memory.NewStorage(),
		nil,
		&git.CloneOptions{
			URL:           repoURL,
			ReferenceName: plumbing.NewBranchReferenceName(branchName),
			SingleBranch:  true,
			Progress:      os.Stdout,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error cloning repository: %w", err)
	}
	return repo, nil
}

// ListCommitsInRange retrieves commit messages from the repository within the specified date range.
func (ga *GitAnalyzer) ListCommitsInRange(
	ctx context.Context, params CommitRangeParams,
) (string, error) {
	if params.Repo == nil {
		return "", errors.New("repository cannot be nil")
	}

	ga.logger.Printf(
		"Date range: %s - %s",
		params.Start.Format("2006-01-02"),
		params.End.Format("2006-01-02"),
	)

	var buf strings.Builder
	commitIter, err := params.Repo.Log(
		&git.LogOptions{
			Since: &params.Start,
			Until: &params.End,
			Order: git.LogOrderCommitterTime,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get commit history: %w", err)
	}

	err = commitIter.ForEach(func(cmt *object.Commit) error {
		if strings.Contains(cmt.Author.Name, "dependabot[bot]") ||
			strings.Contains(cmt.Author.Name, "kodiakhq[bot]") {
			return nil
		}

		// Skip commits not from the specified author if author filter is provided
		if params.Author != "" && !strings.Contains(
			strings.ToLower(cmt.Author.Name),
			strings.ToLower(params.Author),
		) {
			return nil
		}

		buf.WriteString(cmt.Message)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error iterating commits: %w", err)
	}

	return buf.String(), nil
}
