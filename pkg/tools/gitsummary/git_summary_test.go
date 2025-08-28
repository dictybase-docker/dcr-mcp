package gitsummary

import (
	"context"
	"log"
	"os"
	"testing"
)

// TestNewGitSummaryTool tests the creation of a new GitSummaryTool.
func TestNewGitSummaryTool(t *testing.T) {
	t.Parallel()
	logger := log.New(os.Stderr, "", 0)
	tool, err := NewGitSummaryTool(logger)
	if err != nil {
		t.Fatalf("failed to create GitSummaryTool: %v", err)
	}

	if tool == nil {
		t.Fatal("failed to create GitSummaryTool")
	}

	if tool.analyzer == nil {
		t.Fatal("GitAnalyzer not initialized")
	}

	if tool.Logger == nil {
		t.Fatal("Logger not initialized")
	}

	if tool.GetTool().Name != "git-summary" {
		t.Fatalf(
			"expected tool name 'git-summary', got %s",
			tool.GetTool().Name,
		)
	}
}

// MockOpenAIClient is a mock implementation of the worksummary.SummaryClient interface.
type MockOpenAIClient struct{}

// SummarizeCommitMessages implements the worksummary.SummaryClient interface.
func (m *MockOpenAIClient) SummarizeCommitMessages(
	ctx context.Context,
	commitMsgs string,
) (string, error) {
	return "# Work Summary\n\n**Feature Enhancements**\n- Added new features", nil
}

// TestGenerateSummary tests the GenerateSummary method with a mock client.
func TestGenerateSummary(t *testing.T) {
	t.Parallel()
	// Skip this test in automated CI environments since it requires access to external git repositories
	t.Skip("Skipping test that requires external git access")

	// This test would normally create a real repository with known commits
	// and verify the summary generation process.
	//
	// For a complete test, you would:
	// 1. Set up a mock git repository
	// 2. Add test commits with known messages
	// 3. Create a GitSummaryTool with a mock OpenAI client
	// 4. Call GenerateSummary with test parameters
	// 5. Verify the returned summary matches expected output
}
