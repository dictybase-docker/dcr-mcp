package worksummary

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	validator "github.com/go-playground/validator/v10"
	dps "github.com/markusmobius/go-dateparser"
	"github.com/markusmobius/go-dateparser/date"
)

// Global validator instance.
var validate = validator.New()

// GitAnalyzer handles git repository analysis operations including cloning
// repositories, parsing dates, and retrieving commit histories within specified
// date ranges.
type GitAnalyzer struct {
	logger     *log.Logger
	dateConfig *dps.Configuration
}

// CommitRangeParams holds parameters for listing commits in a date range.
type CommitRangeParams struct {
	Repo   *git.Repository `validate:"required"`
	Start  time.Time       `validate:"required"`
	End    time.Time       `validate:"required"`
	Author string          `validate:"required"`
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
	gitAnalyzer := &GitAnalyzer{
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
		opt(gitAnalyzer)
	}

	return gitAnalyzer
}

func (ga *GitAnalyzer) parseStartDate(dateStr string) (date.Date, error) {
	// Validate input
	if err := validate.Var(dateStr, "required"); err != nil {
		return date.Date{}, fmt.Errorf("start date cannot be empty: %w", err)
	}

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
func (ga *GitAnalyzer) ParseAnalysisDates(
	startDate, endDate string,
) (date.Date, date.Date, error) {
	// Validate startDate
	if err := validate.Var(startDate, "required"); err != nil {
		return date.Date{}, date.Date{}, fmt.Errorf(
			"start date cannot be empty: %w",
			err,
		)
	}

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
	// Validate inputs
	if err := validate.Var(repoURL, "required"); err != nil {
		return nil, fmt.Errorf("repository URL cannot be empty: %w", err)
	}
	if err := validate.Var(branchName, "required"); err != nil {
		return nil, fmt.Errorf("branch name cannot be empty: %w", err)
	}

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
	// Validate params using validator
	if err := validate.Struct(params); err != nil {
		return "", fmt.Errorf("invalid commit range parameters: %w", err)
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
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

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
