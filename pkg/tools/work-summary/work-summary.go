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
// repositories, parsing dates, retrieving commit histories within specified
// date ranges, and generating summaries of commit messages using AI. It
// maintains its own logger for operation tracking and date parsing
// configuration.
type GitAnalyzer struct {
	logger     *log.Logger
	dateConfig *dps.Configuration
}

// CommitRangeParams holds parameters for listing commits in a date range
type CommitRangeParams struct {
	Repo   *git.Repository
	Start  time.Time
	End    time.Time
	Author string
}

func NewGitAnalyzer() *GitAnalyzer {
	return &GitAnalyzer{
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

func (ga *GitAnalyzer) ParseAnalysisDates(
	startDate, endDate string,
) (date.Date, date.Date, error) {
	start, err := ga.parseStartDate(startDate)
	if err != nil {
		return start, date.Date{}, fmt.Errorf("Invalid start date: %v", err)
	}
	end, err := ga.parseEndDate(endDate)
	if err != nil {
		return start, end, fmt.Errorf("Invalid end date: %v", err)
	}
	return start, end, nil
}

func (ga *GitAnalyzer) CloneAndCheckout(
	repoURL, branchName string,
) (*git.Repository, error) {
	ga.logger.Printf("Analyzing repository: %s", repoURL)
	ga.logger.Printf("Cloning branch: %s", branchName)
	repo, err := git.Clone(
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
		return nil, fmt.Errorf("error in cloning repository %w", err)
	}
	return repo, nil
}

func (ga *GitAnalyzer) ListCommitsInRange(
	params CommitRangeParams,
) (strings.Builder, error) {
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
		return buf, fmt.Errorf("Failed to get commit history: %v", err)
	}
	err = commitIter.ForEach(func(cmt *object.Commit) error {
		if strings.Contains(cmt.Author.Name, "dependabot[bot]") ||
			strings.Contains(cmt.Author.Name, "kodiakhq[bot]") {
			return nil
		}

		// Skip commits not from the specified author if author filter is provided
		if !strings.Contains(
			strings.ToLower(cmt.Author.Name),
			strings.ToLower(params.Author),
		) {
			return nil
		}

		buf.WriteString(cmt.Message)
		return nil
	})
	if err != nil {
		return buf, err
	}

	return buf, nil
}

func (ga *GitAnalyzer) SummarizeCommitMessages(
	apiKey, model, baseURL string,
	commitMsgs strings.Builder,
) (strings.Builder, error) {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	client := openai.NewClientWithConfig(config)

	req := openai.ChatCompletionRequest{
		Model:       model,
		Stream:      true,
		Temperature: 0.1, // Controls randomness in the response
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: GitSummaryPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: commitMsgs.String(),
			},
		},
	}

	var sb strings.Builder
	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		return sb, fmt.Errorf("OpenAI stream error: %v", err)
	}
	defer stream.Close()

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return sb, fmt.Errorf("OpenAI stream recv error: %v", err)
		}
		sb.WriteString(resp.Choices[0].Delta.Content)
	}
	return sb, nil
}
