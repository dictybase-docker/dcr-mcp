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

// OpenAIClient implements SummaryClient using OpenAI API.
type OpenAIClient struct {
	client *openai.Client
	model  string
	config openai.ClientConfig
}

// OpenAIClientOption defines a functional option for configuring OpenAIClient.
type OpenAIClientOption func(*OpenAIClient)

// WithBaseURL sets a custom base URL for the OpenAI client.
func WithBaseURL(baseURL string) OpenAIClientOption {
	return func(c *OpenAIClient) {
		if baseURL != "" {
			c.config.BaseURL = baseURL
		}
	}
}

// WithModel sets a custom model for the OpenAI client.
func WithModel(model string) OpenAIClientOption {
	return func(c *OpenAIClient) {
		if model != "" {
			c.model = model
		}
	}
}

// NewOpenAIClient creates a new OpenAI client with the provided configuration.
// Uses functional option pattern, default value of BaseURL is
// https://openrouter.ai/api/v1.
func NewOpenAIClient(
	apiKey string,
	opts ...OpenAIClientOption,
) (*OpenAIClient, error) {
	if err := validate.Var(apiKey, "required"); err != nil {
		return nil, errors.New("API key is required")
	}
	llm := &OpenAIClient{
		model:  "google/gemini-2.5-flash-lite",
		config: openai.DefaultConfig(apiKey),
	}
	llm.config.BaseURL = "https://openrouter.ai/api/v1"
	// Apply all options
	for _, opt := range opts {
		opt(llm)
	}
	llm.client = openai.NewClientWithConfig(llm.config)

	return llm, nil
}

// SummarizeCommitMessages generates a summary of commit messages using OpenAI.
func (c *OpenAIClient) SummarizeCommitMessages(
	ctx context.Context,
	commitMsgs string,
) (string, error) {
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

	var stringBuilder strings.Builder
	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", fmt.Errorf("OpenAI stream error: %w", err)
	}
	defer stream.Close()

	for {
		select {
		case <-ctx.Done():
			return stringBuilder.String(), ctx.Err()
		default:
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return stringBuilder.String(), nil
			}
			if err != nil {
				return stringBuilder.String(), fmt.Errorf(
					"OpenAI stream recv error: %w",
					err,
				)
			}
			stringBuilder.WriteString(resp.Choices[0].Delta.Content)
		}
	}
}
