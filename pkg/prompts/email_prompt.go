package prompts

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
)

// EmailPrompt defines the structure for the email generation prompt.
type EmailPrompt struct {
	Name        string
	Description string
	Prompt      mcp.Prompt
	Logger      *log.Logger
}

// NewEmailPrompt creates a new EmailPrompt instance.
func NewEmailPrompt(logger *log.Logger) (*EmailPrompt, error) {
	// Define the dynamic email prompt template
	prompt := mcp.NewPrompt(
		"generate_casual_email", // Unique name for the prompt
		mcp.WithPromptDescription(
			"Generates a draft casual email based on sender and recipient.",
		),
		// Define the 'from' argument
		mcp.WithArgument("from",
			mcp.ArgumentDescription("The sender's email address or name."),
			mcp.RequiredArgument(), // Make 'from' mandatory
		),
		// Define the 'to' argument
		mcp.WithArgument("to",
			mcp.ArgumentDescription("The recipient's email address or name."),
			mcp.RequiredArgument(), // Make 'to' mandatory
		),
	)

	return &EmailPrompt{
		Name:        "generate_email",
		Description: "Generates a draft email based on sender and recipient.",
		Prompt:      prompt,
		Logger:      logger,
	}, nil
}

// GetName returns the name of the prompt.
func (ep *EmailPrompt) GetName() string {
	return ep.Name
}

// GetDescription returns the description of the prompt.
func (ep *EmailPrompt) GetDescription() string {
	return ep.Description
}

// GetPrompt returns the MCP Prompt definition.
func (ep *EmailPrompt) GetPrompt() mcp.Prompt {
	return ep.Prompt
}

// Handler generates the prompt content based on the request arguments.
func (ep *EmailPrompt) Handler(
	ctx context.Context,
	request mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	fromArg, fromOk := request.Params.Arguments["from"]
	if !fromOk {
		return nil, fmt.Errorf("required argument 'from' is missing")
	}
	toArg, toOk := request.Params.Arguments["to"]
	if !toOk {
		return nil, fmt.Errorf("required argument 'to' is missing")
	}

	// Construct the dynamic prompt message content
	// This prompt instructs the LLM on how to assist the user (fromArg)
	// in writing a casual email to toArg.
	promptContent := fmt.Sprintf(
		`You are a helpful assistant aiding %s in drafting a casual and friendly email to %s.
		%s will provide a brief idea of what they want to write. Your task is to help them flesh out the content.

		Here's how you should respond:
			1. Suggest a suitable subject line.
			2. Suggest body paragraphs based on the idea provided.
			3. Include relevant details, potentially suggesting
		colloquial expressions, emojis, or other informal language
		appropriate for an email to a friend (%s).
			4. **Crucially:** If the initial idea is unclear or
		ambiguous, ask clarifying questions to get the necessary details
		before suggesting content. For example, if %s says 'I want to
		invite %s to a concert', you might ask 'Cool! What's the
		band/artist? Got the date, time, and place handy?'.
			5. Maintain a relaxed, friendly, and conversational tone
		throughout your response.`,
		fromArg,
		toArg,
		fromArg, // User providing the idea
		toArg,   // Friend receiving the email
		fromArg, // User asking about concert
		toArg,   // Friend invited to concert
	)

	// Create the prompt result structure
	// We use RoleAssistant here to provide the initial instruction/template.
	// RoleSystem could also be used for meta-instructions.
	result := mcp.NewGetPromptResult(
		"Email Draft Request", // Title for the prompt result
		[]mcp.PromptMessage{ // List of messages forming the prompt
			mcp.NewPromptMessage(
				mcp.RoleAssistant,                 // The role providing the prompt content
				mcp.NewTextContent(promptContent), // The actual text content
			),
		},
	)

	return result, nil
}
