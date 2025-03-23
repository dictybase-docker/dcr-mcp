package calculator

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// Calculator represents a simple calculator tool
type Calculator struct {
	Name        string
	Description string
	Tool        mcp.Tool
}

// CalculateParams defines the parameters for the calculator
type CalculateParams struct {
	Operation string  `json:"operation"`
	OperandA  float64 `json:"operandA"`
	OperandB  float64 `json:"operandB"`
}

// CalculateResult defines the result structure
type CalculateResult struct {
	Result float64 `json:"result"`
}

// NewCalculator creates a new calculator tool
func NewCalculator() (*Calculator, error) {
	// Create the tool with proper schema
	tool := mcp.NewTool("calculate",
		mcp.WithDescription("A simple calculator tool that performs basic math operations"),
		mcp.WithString("operation", 
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
			mcp.Required(),
		),
		mcp.WithNumber("operandA", 
			mcp.Description("The first operand"),
			mcp.Required(),
		),
		mcp.WithNumber("operandB", 
			mcp.Description("The second operand"),
			mcp.Required(),
		),
	)

	return &Calculator{
		Name:        "calculate",
		Description: "A simple calculator tool that performs basic math operations",
		Tool:        tool,
	}, nil
}

// GetName returns the name of the tool
func (c *Calculator) GetName() string {
	return c.Name
}

// GetDescription returns the description of the tool
func (c *Calculator) GetDescription() string {
	return c.Description
}

// GetSchema returns the JSON schema for the tool's parameters
func (c *Calculator) GetSchema() mcp.ToolInputSchema {
	return c.Tool.InputSchema
}

// GetTool returns the MCP Tool
func (c *Calculator) GetTool() mcp.Tool {
	return c.Tool
}

// Execute performs the calculation based on provided parameters
func (c *Calculator) Execute(paramsJSON string) (string, error) {
	// Parse parameters
	var params CalculateParams
	if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
		return "", fmt.Errorf("failed to parse parameters: %w", err)
	}

	// Perform calculation
	var result float64
	switch params.Operation {
	case "add":
		result = params.OperandA + params.OperandB
	case "subtract":
		result = params.OperandA - params.OperandB
	case "multiply":
		result = params.OperandA * params.OperandB
	case "divide":
		if params.OperandB == 0 {
			return "", fmt.Errorf("division by zero not allowed")
		}
		result = params.OperandA / params.OperandB
	default:
		return "", fmt.Errorf("unsupported operation: %s", params.Operation)
	}

	// Return result
	response := CalculateResult{Result: result}
	resJSON, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resJSON), nil
}