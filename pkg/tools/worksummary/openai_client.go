package worksummary

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

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

// SummaryClient is the interface for clients that can generate summaries.
type SummaryClient interface {
	SummarizeCommitMessages(
		ctx context.Context,
		commitMsgs string,
	) (string, error)
}

// OpenAIConfig holds configuration for OpenAI client
type OpenAIConfig struct {
	APIKey  string `validate:"required"`
	BaseURL string
	Model   string `validate:"required"`
}

// OpenAIClient implements SummaryClient using OpenAI API.
type OpenAIClient struct {
	client *openai.Client
	model  string
}

// NewOpenAIClient creates a new OpenAI client with the provided configuration.
func NewOpenAIClient(config OpenAIConfig) (*OpenAIClient, error) {
	// Validate the config
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	openaiConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		openaiConfig.BaseURL = config.BaseURL
	}

	return &OpenAIClient{
		client: openai.NewClientWithConfig(openaiConfig),
		model:  config.Model,
	}, nil
}

// SummarizeCommitMessages generates a summary of commit messages using OpenAI.
func (c *OpenAIClient) SummarizeCommitMessages(
	ctx context.Context,
	commitMsgs string,
) (string, error) {
	if c.client == nil {
		return "", errors.New("OpenAI client not configured")
	}

	// Validate input
	if err := validate.Var(commitMsgs, "required"); err != nil {
		return "", fmt.Errorf("commit messages cannot be empty: %w", err)
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
				return sb.String(), fmt.Errorf(
					"OpenAI stream recv error: %w",
					err,
				)
			}
			sb.WriteString(resp.Choices[0].Delta.Content)
		}
	}
}
